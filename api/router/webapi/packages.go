// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package webapi

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/UpdateHub/updatehub/installmodes/copy"
	_ "github.com/UpdateHub/updatehub/installmodes/flash"
	_ "github.com/UpdateHub/updatehub/installmodes/imxkobs"
	_ "github.com/UpdateHub/updatehub/installmodes/mender"
	_ "github.com/UpdateHub/updatehub/installmodes/raw"
	_ "github.com/UpdateHub/updatehub/installmodes/tarball"
	_ "github.com/UpdateHub/updatehub/installmodes/ubifs"
	_ "github.com/UpdateHub/updatehub/installmodes/zephyr"

	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/UpdateHub/updatehub/libarchive"
	"github.com/UpdateHub/updatehub/metadata"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
)

const (
	GetAllPackagesUrl = "/packages"
	GetPackageUrl     = "/packages/:uid"
	UploadPackageUrl  = "/packages"
)

type PackagesAPI struct {
	db  *storm.DB
	dir string
}

func NewPackagesAPI(db *storm.DB, dir string) *PackagesAPI {
	return &PackagesAPI{db: db, dir: dir}
}

func (api *PackagesAPI) GetAllPackages(c echo.Context) error {
	var packages []models.Package
	if err := api.db.All(&packages); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, packages)
}

func (api *PackagesAPI) GetPackage(c echo.Context) error {
	var pkg models.Package
	if err := api.db.One("UID", c.Param("uid"), &pkg); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pkg)
}

func (api *PackagesAPI) UploadPackage(c echo.Context) error {
	c.Request().ParseMultipartForm(0)

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	metadata, rawMetadata, signature, err := parsePackage(src.(*os.File).Name())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse package file")
	}

	uid := sha256.Sum256(rawMetadata)

	dst, err := os.Create(fmt.Sprintf("%s/%x", api.dir, uid))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	var supportedHardware []string
	switch t := metadata.SupportedHardware.(type) {
	case []interface{}:
		for _, s := range t {
			supportedHardware = append(supportedHardware, s.(string))
		}
	case interface{}:
		supportedHardware = append(supportedHardware, t.(string))
	}

	pkg := &models.Package{
		UID:               fmt.Sprintf("%x", uid),
		Version:           metadata.Version,
		SupportedHardware: supportedHardware,
		Signature:         signature,
		Metadata:          rawMetadata,
	}

	if err := api.db.Save(pkg); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pkg)
}

func parsePackage(file string) (*metadata.UpdateMetadata, []byte, []byte, error) {
	la := &libarchive.LibArchive{}

	reader, err := libarchive.NewReader(la, file, 10240)
	if err != nil {
		return nil, nil, nil, err
	}
	defer reader.Free()

	data := bytes.NewBuffer(nil)
	err = reader.ExtractFile("metadata", data)
	if err != nil {
		return nil, nil, nil, err
	}

	metadata, err := metadata.NewUpdateMetadata(data.Bytes())
	if err != nil {
		return nil, nil, nil, err
	}

	reader, err = libarchive.NewReader(la, file, 10240)
	if err != nil {
		return metadata, data.Bytes(), nil, err
	}

	signature := bytes.NewBuffer(nil)
	err = reader.ExtractFile("signature", signature)
	if err != nil {
		return metadata, data.Bytes(), nil, err
	}

	return metadata, data.Bytes(), signature.Bytes(), nil
}
