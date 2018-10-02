package webapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/gustavosbarreto/updatehub-server/models"
	"github.com/labstack/echo"
)

const (
	GetAllRolloutsUrl       = "/rollouts"
	GetRolloutUrl           = "/rollouts/:id"
	GetRolloutStatisticsUrl = "/rollouts/:id/statistics"
	CreateRolloutUrl        = "/rollouts"
)

type RolloutsAPI struct {
	db *storm.DB
}

func NewRolloutsAPI(db *storm.DB) *RolloutsAPI {
	return &RolloutsAPI{db: db}
}

func (api *RolloutsAPI) GetAllRollouts(c echo.Context) error {
	var rollouts []models.Rollout
	if err := api.db.All(&rollouts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rollouts)
}

func (api *RolloutsAPI) GetRollout(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var rollout models.Rollout
	if err = api.db.One("ID", id, &rollout); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rollout)
}

func (api *RolloutsAPI) GetRolloutStatistics(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var rollout models.Rollout
	if err = api.db.One("ID", id, &rollout); err != nil {
		return err
	}

	var statistics struct {
		Status   string `json:"status"`
		Statuses struct {
			Pending  int `json:"pending"`
			Updating int `json:"updating"`
			Updated  int `json:"updated"`
			Failed   int `json:"failed"`
		} `json:"statuses"`
	}

	for _, uid := range rollout.Devices {
		var d models.Device
		if err = api.db.One("UID", uid, &d); err != nil {
			continue
		}

		switch d.Status {
		case "pending":
			statistics.Statuses.Pending = statistics.Statuses.Pending + 1
		case "downloading", "downloaded", "installing", "installed", "rebooting":
			statistics.Statuses.Updating = statistics.Statuses.Updating + 1
		case "updated":
			statistics.Statuses.Updated = statistics.Statuses.Updated + 1
		case "error":
			statistics.Statuses.Failed = statistics.Statuses.Failed + 1
		}
	}

	if rollout.Running {
		statistics.Status = "running"
	} else {
		if rollout.FinishedAt.After(rollout.StartedAt) {
			if statistics.Statuses.Updated == len(rollout.Devices) {
				statistics.Status = "finished"
			} else if statistics.Statuses.Failed > 0 {
				statistics.Status = "failed"
			}
		} else {
			statistics.Status = "paused"
		}
	}

	return c.JSON(http.StatusOK, statistics)
}

func (api *RolloutsAPI) CreateRollout(c echo.Context) error {
	var body struct {
		Package string   `json:"package"`
		Devices []string `json:"devices"`
		Running bool     `json:"running"`
	}

	c.Bind(&body)

	for _, uid := range body.Devices {
		var device models.Device

		if err := api.db.One("UID", uid, &device); err != nil {
			return err
		}

		r, err := device.ActiveRollout(api.db)
		if err != nil {
			return err
		}

		if r != nil {
			return echo.NewHTTPError(http.StatusNotAcceptable)
		}
	}

	var pkg models.Package
	if err := api.db.One("UID", body.Package, &pkg); err != nil {
		return err
	}

	rollout := models.Rollout{
		Package: body.Package,
		Devices: body.Devices,
		Running: body.Running,
	}

	if rollout.Running {
		rollout.StartedAt = time.Now()
	}

	if err := api.db.Save(&rollout); err != nil {
		return err
	}

	for _, uid := range rollout.Devices {
		var d models.Device
		if err := api.db.One("UID", uid, &d); err != nil {
			continue
		}

		d.Status = "pending"

		if err := api.db.Save(&d); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, rollout)
}
