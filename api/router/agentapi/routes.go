package agentapi

import (
	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Group, db *storm.DB) {
	api := NewAgentAPI(db)

	e.POST(GetRolloutForDeviceUrl, api.GetRolloutForDevice)
	e.POST(ReportDeviceStateUrl, api.ReportDeviceState)
	e.GET(GetObjectFromPackageUrl, api.GetObjectFromPackage)
}
