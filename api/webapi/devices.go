package webapi

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gustavosbarreto/updatehub-server/models"
	"github.com/labstack/echo"
)

const (
	GetAllDevicesUrl = "/devices"
	GetDeviceUrl     = "/devices/:uid"
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
