// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package agentapi

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/UpdateHub/updatehub-ce/metadata"
	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
)

const (
	GetRolloutForDeviceUrl  = "/upgrades"
	ReportDeviceStateUrl    = "/report"
	GetObjectFromPackageUrl = "/products/:product/packages/:package/objects/:object"
)

type AgentAPI struct {
	db *storm.DB
}

func NewAgentAPI(db *storm.DB) *AgentAPI {
	return &AgentAPI{db: db}
}

func (api *AgentAPI) GetRolloutForDevice(c echo.Context) error {
	var metadata struct {
		Retries int `json:"retries"`
		metadata.FirmwareMetadata
		LastInstalledPackage string `json:"last-installed-package,omitempty"`
	}

	c.Bind(&metadata)

	deviceIdentity, err := json.Marshal(metadata.DeviceIdentity)
	if err != nil {
		return err
	}

	uid := sha256.Sum256(deviceIdentity)

	var device models.Device
	if err = api.db.One("UID", uid, &device); err != nil && err != storm.ErrNotFound {
		return err
	}

	if device.UID == "" {
		device.UID = fmt.Sprintf("%x", uid)
	}

	device.Version = metadata.Version
	device.Hardware = metadata.Hardware
	device.ProductUID = metadata.ProductUID
	device.DeviceIdentity = metadata.DeviceIdentity
	device.DeviceAttributes = metadata.DeviceAttributes
	device.LastSeen = time.Now()

	var rollouts []models.Rollout
	if err = api.db.All(&rollouts); err != nil {
		return err
	}

	var rollout *models.Rollout

	for _, r := range rollouts {
		for _, d := range r.Devices {
			if d == device.UID && r.Running {
				rollout = &r
				break
			}
		}

		if rollout != nil {
			break
		}
	}

	if rollout == nil || !rollout.Running {
		device.Version = metadata.Version

		if err := api.db.Save(&device); err != nil {
			return err
		}

		return c.NoContent(http.StatusNotFound)
	}

	var pkg models.Package
	if err := api.db.One("UID", rollout.Package, &pkg); err != nil {
		return err
	}

	if pkg.Version == metadata.Version {
		report := models.Report{
			Device:    device.UID,
			Rollout:   rollout.ID,
			Status:    "updated",
			Timestamp: time.Now(),
			IsError:   false,
			Virtual:   true,
		}

		if err := api.db.Save(&report); err != nil {
			return err
		}

		device.Status = "updated"

		if err := api.db.Save(&device); err != nil {
			return err
		}

		finished, err := rollout.IsFinished(api.db)
		if err != nil {
			return err
		}

		if !finished {
			rollout.FinishedAt = time.Now()
			rollout.Running = false

			if err = api.db.Save(rollout); err != nil {
				return err
			}
		}

		return c.NoContent(http.StatusNotFound)
	}

	c.Response().Header().Set("UH-Signature", string(pkg.Signature))
	c.Response().WriteHeader(http.StatusOK)

	_, err = io.Copy(c.Response(), bytes.NewReader(pkg.Metadata))
	return err
}

func (api *AgentAPI) ReportDeviceState(c echo.Context) error {
	var body struct {
		metadata.FirmwareMetadata
		ErrorMessage  string `json:"error-message"`
		PreviousState string `json:"previous-state"`
		Status        string `json:"status"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	deviceIdentity, err := json.Marshal(body.DeviceIdentity)
	if err != nil {
		return err
	}

	uid := sha256.Sum256(deviceIdentity)

	var device models.Device
	if err := api.db.One("UID", fmt.Sprintf("%x", uid), &device); err != nil {
		return err
	}

	rollout, err := device.ActiveRollout(api.db)
	if err != nil {
		return err
	}

	if rollout == nil || !rollout.Running {
		return c.NoContent(http.StatusOK)
	}

	device.Status = body.Status

	if err := api.db.Save(&device); err != nil {
		return err
	}

	if body.Status == "error" {
		rollout.Running = false
		rollout.FinishedAt = time.Now()

		if err := api.db.Save(rollout); err != nil {
			return err
		}
	}

	report := &models.Report{
		Rollout:   rollout.ID,
		Device:    device.UID,
		Message:   body.ErrorMessage,
		Status:    body.Status,
		Timestamp: time.Now(),
		IsError:   false,
		Virtual:   false,
	}

	if body.Status == "error" {
		report.Status = body.PreviousState
		report.IsError = true
	}

	if err := api.db.Save(report); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *AgentAPI) GetObjectFromPackage(c echo.Context) error {
	objectUID := c.Param("object")
	packageUID := c.Param("package")

	lower, upper, err := parseRangeHeader(c.Request().Header["Range"])
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(packageUID)

	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name == objectUID {
			f, err := f.Open()
			if err != nil {
				return err
			}
			defer f.Close()

			rd, _ := f.(io.Reader)

			if lower != nil {
				step := 128 * 1024
				read := 0
				for read < *lower {
					buf_size := min((*lower - read), step)
					buf := make([]byte, buf_size)

					n, err := io.ReadFull(f, buf)
					if err != nil {
						return err
					}
					read += n
				}

				if upper != nil {
					rd = io.LimitReader(f, int64(*upper-*lower)+1)
				}
			}

			_, err = io.Copy(c.Response(), rd)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Package does not contain object %q", objectUID)
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func parseRangeHeader(headers []string) (*int, *int, error) {
	if len(headers) < 1 {
		return nil, nil, nil
	}
	header := headers[0]
	re := regexp.MustCompile("bytes=([\\d]*)-([\\d]*)")

	ms := re.FindAllStringSubmatch(header, -1)
	if len(ms) < 1 {
		return nil, nil, errors.New("Failed to parse byte rage header")
	}

	var lower *int
	var upper *int

	if len(ms[0][1]) > 0 {
		l, err := strconv.Atoi(ms[0][1])
		if err != nil {
			return nil, nil, errors.New("Failed to parse lower bound from rage header")
		}
		lower = &l
	}

	if len(ms[0][2]) > 0 {
		u, err := strconv.Atoi(ms[0][2])
		if err != nil {
			return nil, nil, errors.New("Failed to parse upper bound from rage header")
		}
		upper = &u
	}

	return lower, upper, nil
}
