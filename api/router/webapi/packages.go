// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package webapi

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	// Modules have to be included so the parser can support these install modes
	_ "github.com/UpdateHub/updatehub-ce/installmodes/copy"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/flash"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/imxkobs"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/mender"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/raw"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/raw_delta"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/tarball"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/ubifs"
	_ "github.com/UpdateHub/updatehub-ce/installmodes/zephyr"

	"github.com/UpdateHub/updatehub-ce/metadata"
	"github.com/UpdateHub/updatehub-ce/models"
	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
)

const (
	GetAllPackagesUrl = "/packages"
	GetPackageUrl     = "/packages/:uid"
	UploadPackageUrl  = "/packages"
	DeletePackageUrl  = "/packages/:uid/delete"
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

func (api *PackagesAPI) DeletePackage(c echo.Context) error {
	var pkg models.Package
	if err := api.db.One("UID", c.Param("uid"), &pkg); err != nil {
		return err
	}

	if err := api.db.DeleteStruct(&pkg); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
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

	pkg := &models.Package{
		UID:       fmt.Sprintf("%x", uid),
		Version:   metadata.Version,
		Signature: signature,
		Metadata:  rawMetadata,
	}

	switch t := metadata.SupportedHardware.(type) {
	case []interface{}:
		var supportedHardware []string
		for _, s := range t {
			supportedHardware = append(supportedHardware, s.(string))
		}
		pkg.SupportedHardware = supportedHardware
	case interface{}:
		pkg.SupportedHardware = t.(string)
	}

	if err := api.db.Save(pkg); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pkg)
}

func parsePackage(file string) (*metadata.UpdateMetadata, []byte, []byte, error) {
	reader, err := zip.OpenReader(file)
	if err != nil {
		return nil, nil, nil, err
	}
	defer reader.Close()

	var update_metadata *metadata.UpdateMetadata
	var signature []byte
	var rawMetadata []byte

	for _, f := range reader.File {
		if f.Name == "metadata" {
			f, err := f.Open()
			if err != nil {
				return nil, nil, nil, err
			}
			defer f.Close()
			data := bytes.NewBuffer(nil)

			_, err = io.Copy(data, f)
			if err != nil {
				return nil, nil, nil, err
			}

			update_metadata, err = metadata.NewUpdateMetadata(data.Bytes())
			if err != nil {
				return nil, nil, nil, err
			}

			rawMetadata = data.Bytes()
		}

		if f.Name == "signature" {
			f, err := f.Open()
			if err != nil {
				return nil, nil, nil, err
			}
			defer f.Close()
			data := bytes.NewBuffer(nil)

			_, err = io.Copy(data, f)
			if err != nil {
				return nil, nil, nil, err
			}

			signature = data.Bytes()
		}
	}

	if update_metadata == nil {
		return nil, nil, nil, errors.New("Package missing metadata file")
	}

	if signature == nil {
		return nil, nil, nil, errors.New("Package missing signature file")
	}

	return update_metadata, rawMetadata, signature, nil
}
