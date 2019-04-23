package webapi

import (
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func SetupRoutes(e *echo.Group, db *storm.DB) {
	devicesEndpoint := NewDevicesAPI(db)
	e.GET(GetAllDevicesUrl, devicesEndpoint.GetAllDevices)
	e.GET(GetDeviceUrl, devicesEndpoint.GetDevice)
	e.GET(GetDeviceRolloutReportsUrl, devicesEndpoint.GetDeviceRolloutReports)

	packagesEndpoint := NewPackagesAPI(db, viper.GetString("dir"))
	e.GET(GetAllPackagesUrl, packagesEndpoint.GetAllPackages)
	e.GET(GetPackageUrl, packagesEndpoint.GetPackage)
	e.DELETE(DeletePackageUrl, packagesEndpoint.DeletePackage)
	e.POST(UploadPackageUrl, packagesEndpoint.UploadPackage)

	rolloutsEndpoint := NewRolloutsAPI(db)
	e.GET(GetAllRolloutsUrl, rolloutsEndpoint.GetAllRollouts)
	e.GET(GetRolloutUrl, rolloutsEndpoint.GetRollout)
	e.GET(GetRolloutStatisticsUrl, rolloutsEndpoint.GetRolloutStatistics)
	e.GET(GetRolloutDevicesUrl, rolloutsEndpoint.GetRolloutDevices)
	e.POST(CreateRolloutUrl, rolloutsEndpoint.CreateRollout)
	e.PUT(StopRolloutUrl, rolloutsEndpoint.StopRollout)
}
