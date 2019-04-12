// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package agentapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/UpdateHub/updatehub/libarchive"
	"github.com/UpdateHub/updatehub/metadata"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
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

	reader, err := libarchive.NewReader(&libarchive.LibArchive{}, packageUID, 10240)
	if err != nil {
		return err
	}
	defer reader.Free()

	if err := reader.ExtractFile(objectUID, c.Response()); err != nil {
		return err
	}

	return nil
}
