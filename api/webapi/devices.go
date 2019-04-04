// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package webapi

import (
	"net/http"
	"strconv"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/labstack/echo"
)

const (
	GetAllDevicesUrl           = "/devices"
	GetDeviceUrl               = "/devices/:uid"
	GetDeviceRolloutReportsUrl = "/devices/:uid/rollouts/:id/reports"
)

type DevicesAPI struct {
	db *storm.DB
}

func NewDevicesAPI(db *storm.DB) *DevicesAPI {
	return &DevicesAPI{db: db}
}

func (api *DevicesAPI) GetAllDevices(c echo.Context) error {
	var devices []models.Device
	if err := api.db.All(&devices); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, devices)
}

func (api *DevicesAPI) GetDevice(c echo.Context) error {
	var device models.Device
	if err := api.db.One("UID", c.Param("uid"), &device); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, device)
}

func (api *DevicesAPI) GetDeviceRolloutReports(c echo.Context) error {
	var device models.Device
	if err := api.db.One("UID", c.Param("uid"), &device); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var rollout models.Rollout
	if err := api.db.One("ID", id, &rollout); err != nil {
		return err
	}

	var reports []models.Report
	if err := api.db.Select(q.And(q.Eq("Device", device.UID), q.Eq("Rollout", rollout.ID))).OrderBy("Timestamp").Find(&reports); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reports)
}
