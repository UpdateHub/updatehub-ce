package webapi

import (
	"encoding/json"
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

type RolloutTestSuite struct {
	suite.Suite
	db *storm.DB
}

func (suite *RolloutTestSuite) SetupTest() {
	db, err := storm.Open("updatehub_test.db")
	assert.NoError(suite.T(), err)
	suite.db = db
}

func (suite *RolloutTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *RolloutTestSuite) TestGetRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	rollout := models.Rollout{
		ID:        1,
		Package:   "0000000",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	expected, err := json.Marshal(rollout)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *RolloutTestSuite) TestGetRolloutWithNoIntegerElement() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/id", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestGetRolloutWithNoRollOut() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestGetRollouts() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	startedAt := time.Date(2018, 8, 14, 12, 0, 0, 0, time.UTC)
	rollout1 := models.Rollout{
		ID:        1,
		Package:   "0000000",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	rollout2 := models.Rollout{
		ID:        2,
		Package:   "0000000",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout1)
	assert.NoError(suite.T(), err)

	err = suite.db.Save(&rollout2)
	assert.NoError(suite.T(), err)

	expected, err := json.Marshal([]models.Rollout{rollout1, rollout2})
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *RolloutTestSuite) TestGetRolloutsNoRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), "[]", strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetRolloutDevices() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	pack := models.Package{
		UID:               "1",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte("000"),
		Version:           "1",
	}

	err := suite.db.Save(&pack)
	assert.NoError(suite.T(), err)

	rollout := models.Rollout{
		ID:        1,
		Package:   "1",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err = suite.db.Save(&rollout)
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
	req := httptest.NewRequest(http.MethodGet, "/rollouts/1/devices", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), "null", strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *DeviceTestSuite) TestGetRolloutDevicesWrongDeviceID() {
	defer os.RemoveAll(suite.db.Bolt.Path())

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

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/2/devices", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestGetRolloutStatistcs() {
	defer os.RemoveAll(suite.db.Bolt.Path())
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

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/1/statistics", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	expected := `{"status":"running","statuses":{"pending":1,"updating":0,"updated":0,"failed":0}}`
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *RolloutTestSuite) TestGetRolloutStatistcsWrongIDRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())
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

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/rollouts/2/statistics", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestCreateRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)

	pack := models.Package{
		UID:               "1",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte("000"),
		Version:           "1",
	}

	err := suite.db.Save(&pack)
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

	postUserJSON := `{"package": "1", "devices":["1"]}`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/rollouts", strings.NewReader(postUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	rollout := &models.Rollout{}
	err = json.Unmarshal(rec.Body.Bytes(), rollout)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), bool(true), rollout.Running)
}

func (suite *RolloutTestSuite) TestCreateRolloutNoDevice() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	pack := models.Package{
		UID:               "1",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte("000"),
		Version:           "1",
	}

	err := suite.db.Save(&pack)
	assert.NoError(suite.T(), err)

	postUserJSON := `{"package": "1", "devices":["1"]}`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/rollouts", strings.NewReader(postUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestCreateRolloutNoPackage() {
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

	postUserJSON := `{"package": "1", "devices":["1"]}`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/rollouts", strings.NewReader(postUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *RolloutTestSuite) TestStopRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	rollout := models.Rollout{
		ID:        1,
		Package:   "0000000",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPut, "/rollouts/1/stop", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func (suite *RolloutTestSuite) TestStopRolloutWrongIDRollout() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	rollout := models.Rollout{
		ID:        1,
		Package:   "0000000",
		Devices:   []string{"x86"},
		Running:   true,
		StartedAt: startedAt,
	}

	err := suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPut, "/rollouts/2/stop", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func TestRolloutSuite(t *testing.T) {
	suite.Run(t, new(RolloutTestSuite))
}
