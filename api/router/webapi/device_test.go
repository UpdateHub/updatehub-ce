package webapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeviceTestSuite struct {
	suite.Suite
	db *storm.DB
}

func (suite *DeviceTestSuite) SetupTest() {
	db, err := storm.Open("updatehub_test.db")
	if err != nil {
		log.Fatal(err)
	}

	suite.db = db
}

func (suite *DeviceTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *DeviceTestSuite) TestGetDevice() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	dev := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err := suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	expected, err := json.Marshal(dev)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetDeviceNoDevice() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *DeviceTestSuite) TestGetDevices() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	dev1 := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	dev2 := models.Device{
		UID:            "2",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err := suite.db.Save(&dev1)
	assert.NoError(suite.T(), err)

	err = suite.db.Save(&dev2)
	assert.NoError(suite.T(), err)

	expected, err := json.Marshal([]models.Device{dev1, dev2})
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetDevicesWithoutDevices() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string("[]"), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReports() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	rollout := models.Rollout{
		ID:        1,
		Package:   "1",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	dev := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err = suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	report := models.Report{
		ID:      1,
		Device:  "1",
		Rollout: 1,
	}

	err = suite.db.Save(&report)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1/rollouts/1/reports", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	expected := `[{"id":1,"device":"1","rollout":1,"status":"","error":false,"message":"","timestamp":"0001-01-01T00:00:00Z","virtual":false}]`
	assert.Equal(suite.T(), expected, strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReportsNoDevice() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1/rollouts/1/reports", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReportsWithoutRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	dev := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err := suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1/rollouts/1/reports", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReportsWithoutReport() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	rollout := models.Rollout{
		ID:        1,
		Package:   "1",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	dev := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err = suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1/rollouts/1/reports", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReportsFirstArgumentNoInterger() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/id/rollouts/1/reports", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *DeviceTestSuite) TestGetDeviceRolloutReportsSecondArgumentNoInterger() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	rollout := models.Rollout{
		ID:        1,
		Package:   "1",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	dev := models.Device{
		UID:            "1",
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
		Status:         "pending",
	}

	err = suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/devices/1/rollouts/id/reports", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}
