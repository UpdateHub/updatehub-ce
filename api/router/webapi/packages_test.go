package webapi

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	ValidJSONMetadata = `{
			  "product-uid": "548873c4d4e8e751fdd46c38a3d5a8656cf87bf27a404f346ad58086f627a4ea",
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

	InvalidJSONMetadata = `{
			  "product-
			}`
)

type PackageTestSuite struct {
	suite.Suite
	db *storm.DB
}

func (suite *PackageTestSuite) SetupTest() {
	viper.Set("dir", "/tmp")
	db, err := storm.Open("updatehub_test.db")
	if err != nil {
		log.Fatal(err)
	}

	suite.db = db
}

func (suite *PackageTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *PackageTestSuite) TestGetPackage() {
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

	expected, err := json.Marshal(pack)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/packages/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *PackageTestSuite) TestGetPackageWithoutPackage() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/packages/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *PackageTestSuite) TestGetPackageNoIntegerArgument() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/packages/id", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
}

func (suite *PackageTestSuite) TestGetDevices() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	pack1 := models.Package{
		UID:               "1",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte("000"),
		Version:           "1",
	}

	pack2 := models.Package{
		UID:               "2",
		SupportedHardware: []string{"x86"},
		Signature:         []byte("SIGNATURE"),
		Metadata:          []byte("000"),
		Version:           "2",
	}

	err := suite.db.Save(&pack1)
	assert.NoError(suite.T(), err)

	err = suite.db.Save(&pack2)
	assert.NoError(suite.T(), err)

	expected, err := json.Marshal([]models.Package{pack1, pack2})
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/packages", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string(expected), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *PackageTestSuite) TestGetDevicesWithoutPackages() {
	defer os.RemoveAll(suite.db.Bolt.Path())

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodGet, "/packages", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), string("[]"), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (suite *RolloutTestSuite) TestUploadPackage() {
	defer os.RemoveAll(suite.db.Bolt.Path())
	memFs := afero.NewOsFs()

	bodyBuf, contentType, err := generateUhupkg(memFs, ValidJSONMetadata, true, true)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/packages", bodyBuf)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func (suite *RolloutTestSuite) TestUploadPackageNoSignature() {
	defer os.RemoveAll(suite.db.Bolt.Path())
	memFs := afero.NewOsFs()

	bodyBuf, contentType, err := generateUhupkg(memFs, ValidJSONMetadata, true, false)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/packages", bodyBuf)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *RolloutTestSuite) TestUploadPackageNoMetadata() {
	defer os.RemoveAll(suite.db.Bolt.Path())
	memFs := afero.NewOsFs()

	bodyBuf, contentType, err := generateUhupkg(memFs, ValidJSONMetadata, false, true)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/packages", bodyBuf)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *RolloutTestSuite) TestUploadPackageInvalidMetadata() {
	defer os.RemoveAll(suite.db.Bolt.Path())
	memFs := afero.NewOsFs()

	bodyBuf, contentType, err := generateUhupkg(memFs, InvalidJSONMetadata, false, true)
	assert.NoError(suite.T(), err)

	e := echo.New()
	SetupRoutes(e.Group(""), suite.db)
	req := httptest.NewRequest(http.MethodPost, "/packages", bodyBuf)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func TestPackageSuite(t *testing.T) {
	suite.Run(t, new(PackageTestSuite))
}

func generateUhupkg(fsBackend afero.Fs, metadata string, withMetadata, withSignature bool) (*bytes.Buffer, string, error) {
	zipPath := "/tmp/output.zip"

	file, err := fsBackend.Create(zipPath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	zw := zip.NewWriter(file)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}

	sha256sum := sha256.Sum256([]byte(metadata))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sha256sum[:])
	if err != nil {
		return nil, "", err
	}

	var files = []struct {
		Name, Body string
	}{
		{"d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa", "content1"},
		{"metadata", metadata},
		{"signature", string(signature)},
	}

	for _, file := range files {
		if !withMetadata && file.Name == "metadata" || !withSignature && file.Name == "signature" {
			// if it's an invalid file, write only the metadata
			continue
		}

		f, err := zw.Create(file.Name)
		if err != nil {
			return nil, "", err
		}

		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return nil, "", err
		}
	}

	err = zw.Close()
	if err != nil {
		return nil, "", err
	}

	testPath, err := afero.TempDir(fsBackend, "", "package-test")
	if err != nil {
		return nil, "", err
	}
	defer fsBackend.RemoveAll(testPath)

	pkgpath := path.Join(testPath, "test.pkg")

	os.Link(zipPath, pkgpath)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", "test.pkg")
	if err != nil {
		return nil, "", err
	}
	fh, err := os.Open(pkgpath)
	if err != nil {
		return nil, "", err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	return bodyBuf, contentType, nil
}
