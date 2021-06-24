package agentapi

import (
	"archive/zip"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AgentApiTestSuite struct {
	suite.Suite
	db *storm.DB
}

func (suite *AgentApiTestSuite) SetupTest() {
	db, err := storm.Open("updatehub_test.db")
	assert.NoError(suite.T(), err)

	suite.db = db
}

func (suite *AgentApiTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *AgentApiTestSuite) TestGetRolloutForDevice() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 8, 14, 12, 0, 0, 0, time.UTC)
	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	expectReturn := "Returned_Metadata"

	pack := models.Package{
		UID:               "12",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte(expectReturn),
		Version:           "2",
	}

	err := suite.db.Save(&pack)
	assert.NoError(suite.T(), err)

	hash := sha256.New()
	hash.Write([]byte(`{"serial":"123"}`))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	rollout := models.Rollout{
		ID:        1,
		Package:   "12",
		Devices:   []string{mdStr},
		Running:   true,
		StartedAt: startedAt,
	}

	err = suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	dev := models.Device{
		UID:            mdStr,
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
	}

	err = suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	userPostJson := `{
		"product-uid":"000000000",
		"version":"1",
		"hardware":"x86",
		"device-identity":{
		   "serial":"123"
		}
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/upgrades", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expectReturn), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *AgentApiTestSuite) TestGetRolloutForDeviceNoUpdate() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	userPostJson := `{
		"product-uid":"000000000",
		"version":"1",
		"hardware":"x86",
		"device-identity":{
		   "serial":"123"
		}
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/upgrades", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *AgentApiTestSuite) TestGetRolloutForDeviceWrongSended() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	userPostJson := `{
		"product-uid":"000000000",
		"versi
		"hardware
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/upgrades", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *AgentApiTestSuite) TestReportDeviceState() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	identity := make(map[string]string)
	identity["serial"] = "123"
	endDate := time.Date(2018, 8, 14, 12, 0, 0, 0, time.UTC)
	startedAt := time.Date(2018, 7, 14, 12, 0, 0, 0, time.UTC)
	expectReturn := "Returned_Metadata"

	pack := models.Package{
		UID:               "12",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte(expectReturn),
		Version:           "2",
	}

	err := suite.db.Save(&pack)
	assert.NoError(suite.T(), err)

	hash := sha256.New()
	hash.Write([]byte(`{"serial":"123"}`))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	rollout := models.Rollout{
		ID:        1,
		Package:   "12",
		Devices:   []string{mdStr},
		Running:   true,
		StartedAt: startedAt,
	}

	err = suite.db.Save(&rollout)
	assert.NoError(suite.T(), err)

	dev := models.Device{
		UID:            mdStr,
		Hardware:       "x86",
		DeviceIdentity: identity,
		LastSeen:       endDate,
		ProductUID:     "000000000",
		Version:        "1",
	}

	err = suite.db.Save(&dev)
	assert.NoError(suite.T(), err)

	userPostJson := `{
		"product-uid":"000000000",
		"version":"1",
		"hardware":"x86",
		"device-identity":{
		   "serial":"123"
		},
		"error-message":"mistake",
		"previous-state":"dowloading",
		"status":"error"
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/report", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func (suite *AgentApiTestSuite) TestReportDeviceStateError() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	userPostJson := `{
		"product-uid":"000000000",
		"version":"1",
		"hardware":"x86",
		"device-identity":{
		   "serial":"123"
		},
		"error-message":"mistake",
		"previous-state":"dowloading",
		"status":"error"
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/report", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *AgentApiTestSuite) TestReportDeviceStateWrongSended() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	userPostJson := `{
		"product-":"000000000",
		"version":"1",
		},
		"stat
	 }`

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/report", strings.NewReader(userPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *AgentApiTestSuite) TestGetObjectFromPackage() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	memFs := afero.NewOsFs()

	tarballPath, err := generateArchive(memFs)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/products/000000000/packages/output.zip/objects/d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	err = os.Remove(tarballPath)
	assert.NoError(suite.T(), err)
}

func (suite *AgentApiTestSuite) TestGetObjectFromPackageWrongFile() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	memFs := afero.NewOsFs()

	tarballPath, err := generateArchive(memFs)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/products/000000000/packages//objects/d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)

	err = os.Remove(tarballPath)
	assert.NoError(suite.T(), err)
}

func (suite *AgentApiTestSuite) TestGetObjectFromPackageWrongObject() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	memFs := afero.NewOsFs()

	tarballPath, err := generateArchive(memFs)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/products/000000000/packages/output.zip/objects/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusNotFound, rec.Code)

	err = os.Remove(tarballPath)
	assert.NoError(suite.T(), err)
}

func TestAgentApiSuite(t *testing.T) {
	suite.Run(t, new(AgentApiTestSuite))
}

func generateArchive(fsBackend afero.Fs) (string, error) {
	archivePath := "output.zip"
	file, err := fsBackend.Create(archivePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	metadata := `{
		"product-uid": "000000000",
		"supported-hardware": [
		  "hardware1-revA",
		  "hardware2-revB"
		],
		"objects": [
		  [
			{ "mode": "zephyr", "sha256sum": "d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa" }
		  ]
		]
	  }`

	sha256sum := sha256.Sum256([]byte(metadata))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sha256sum[:])
	if err != nil {
		return "", err
	}

	var files = []struct {
		Name, Body string
	}{
		{"d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa", "content1"},
		{"metadata", metadata},
		{"signature", string(signature)},
	}

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return "", err
		}

		if _, err := f.Write([]byte(file.Body)); err != nil {
			return "", err
		}
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return archivePath, nil
}
