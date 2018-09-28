package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"io"

	"github.com/asdine/storm"
	jwt "github.com/dgrijalva/jwt-go"
	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/updatehub/updatehub/libarchive"
	"github.com/updatehub/updatehub/metadata"
)

type User struct {
	ID       int    `storm:"id,increment"`
	Username string `storm:"index"`
	Password string
}

type Package struct {
	UID               string   `storm:"id" json:"uid"`
	Version           string   `json:"version"`
	SupportedHardware []string `json:"supported_hardware"`
	Signature         []byte   `json:"signature"`
	Metadata          []byte   `json:"metadata"`
}

func main() {
	db, err := storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Static("/", "public/")

	ui, _ := url.Parse("http://localhost:1314/")

	e.Group("/ui", middleware.Proxy(middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{{URL: ui}},
	)))

	e.POST("/", func(c echo.Context) error {
		count, err := db.Count(&User{})
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if count > 0 {
			return echo.ErrUnauthorized
		}

		return db.Save(&User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		})
	})

	e.POST("/login", func(c echo.Context) error {
		var login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		c.Bind(&login)

		if login.Username == "" {
			return echo.ErrUnauthorized
		}

		var user User
		if err = db.One("Username", login.Username, &user); err != nil {
			return echo.ErrUnauthorized
		}

		fmt.Println(user)

		if user.Password == login.Password {
			token := jwt.New(jwt.SigningMethodHS256)

			claims := token.Claims.(jwt.MapClaims)
			claims["name"] = "root"
			claims["admin"] = true
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"name":  "Gustavo",
				"token": t,
			})
		}

		return echo.ErrUnauthorized
	})

	e.POST("/report", func(c echo.Context) error {
		var body struct {
			metadata.FirmwareMetadata
			ErrorMessage  string `json:"error-message"`
			PreviousState string `json:"previous-state"`
			Status        string `json:"status"`
		}

		if err := c.Bind(&body); err != nil {
			return err
		}

		s, _ := prettyjson.Marshal(body)
		fmt.Println(string(s))

		deviceIdentity, err := json.Marshal(body.DeviceIdentity)
		if err != nil {
			return err
		}

		uid := sha256.Sum256(deviceIdentity)

		var device Device
		if err := db.One("UID", fmt.Sprintf("%x", uid), &device); err != nil {
			return err
		}

		fmt.Println("setando o corno", body.Status)
		device.Status = body.Status

		if err := db.Save(&device); err != nil {
			return err
		}

		fmt.Println(device)

		return c.NoContent(http.StatusOK)
	})

	e.GET("/products/:product/packages/:package/objects/:object", func(c echo.Context) error {
		objectUID := c.Param("object")
		packageUID := c.Param("package")

		reader, err := libarchive.NewReader(&libarchive.LibArchive{}, packageUID, 10240)
		if err != nil {
			return err
		}
		defer reader.Free()

		if err := reader.ExtractFile(objectUID, c.Response()); err != nil {
			return err
		}

		return nil
	})

	e.POST("/upgrades", func(c echo.Context) error {
		var metadata struct {
			Retries int `json:"retries"`
			metadata.FirmwareMetadata
			LastInstalledPackage string `json:"last-installed-package,omitempty"`
		}

		c.Bind(&metadata)

		deviceIdentity, err := json.Marshal(metadata.DeviceIdentity)
		if err != nil {
			return err
		}

		uid := sha256.Sum256(deviceIdentity)

		device := &Device{
			UID:              fmt.Sprintf("%x", uid),
			Version:          metadata.Version,
			Hardware:         metadata.Hardware,
			ProductUID:       metadata.ProductUID,
			DeviceIdentity:   metadata.DeviceIdentity,
			DeviceAttributes: metadata.DeviceAttributes,
		}

		if metadata.LastInstalledPackage != "" {
			device.Status = "updated"
		}

		if err := db.Save(device); err != nil {
			return err
		}

		var rollouts []Rollout
		if err = db.All(&rollouts); err != nil {
			return err
		}

		var rollout *Rollout

		for _, r := range rollouts {
			for _, d := range r.Devices {
				if d == device.UID {
					rollout = &r
					break
				}
			}

			if rollout != nil {
				break
			}
		}

		if rollout == nil || !rollout.Running {
			return c.JSON(http.StatusNotFound, nil)
		}

		if device.Status == "updated" {
			if rollout != nil {
				finished, err := rollout.IsFinished(db)
				if err != nil {
					return err
				}

				if !finished {
					rollout.FinishedAt = time.Now()
					rollout.Running = false

					if err = db.Save(rollout); err != nil {
						return err
					}
				}
			}

			return c.JSON(http.StatusNotFound, nil)
		}

		var pkg Package
		if err := db.One("UID", rollout.Package, &pkg); err != nil {
			return err
		}

		c.Response().Header().Set("UH-Signature", string(pkg.Signature))
		c.Response().WriteHeader(http.StatusOK)

		_, err = io.Copy(c.Response(), bytes.NewReader(pkg.Metadata))
		return err
	})

	api := e.Group("/api")
	//	api.Use(middleware.JWT([]byte("secret")))
	api.GET("/devices", func(c echo.Context) error {
		var devices []Device
		if err = db.All(&devices); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, devices)
	})

	api.GET("/devices/:uid", func(c echo.Context) error {
		var device Device
		if err = db.One("UID", c.Param("uid"), &device); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, device)
	})

	e.POST("/devices", func(c echo.Context) error {
		d := new(Device)

		if err = c.Bind(d); err != nil {
			fmt.Println(err)
			return err
		}

		err := db.Save(d)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, d)
	})

	api.GET("/packages", func(c echo.Context) error {
		var packages []Package
		if err = db.All(&packages); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, packages)
	})

	api.GET("/packages/:uid", func(c echo.Context) error {
		var pkg Package
		if err = db.One("UID", c.Param("uid"), &pkg); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, pkg)
	})

	api.POST("/packages", func(c echo.Context) error {
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

		dst, err := os.Create(fmt.Sprintf("%x", uid))
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

		err = db.Save(&Package{
			UID:               fmt.Sprintf("%x", uid),
			Version:           metadata.Version,
			SupportedHardware: supportedHardware,
			Signature:         signature,
			Metadata:          rawMetadata,
		})

		fmt.Println(err)

		return err
	})

	api.GET("/rollouts", func(c echo.Context) error {
		var rollouts []Rollout
		if err = db.All(&rollouts); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, rollouts)
	})

	api.GET("/rollouts/:id/statistics", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		var rollout Rollout
		if err = db.One("ID", id, &rollout); err != nil {
			fmt.Println(err)
			return err
		}

		var statistics struct {
			Pending  int `json:"pending"`
			Updating int `json:"updating"`
			Updated  int `json:"updated"`
			Failed   int `json:"failed"`
		}

		for _, uid := range rollout.Devices {
			var d Device
			if err = db.One("UID", uid, &d); err != nil {
				continue
			}

			switch d.Status {
			case "pending":
				statistics.Pending = statistics.Pending + 1
			case "downloading", "downloaded", "installing", "installed", "rebooting":
				statistics.Updating = statistics.Updating + 1
			case "updated":
				statistics.Updated = statistics.Updated + 1
			case "error":
				statistics.Failed = statistics.Failed + 1
			}
		}

		return c.JSON(http.StatusOK, statistics)
	})

	api.GET("/rollouts/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		var rollout Rollout
		if err = db.One("ID", id, &rollout); err != nil {
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, rollout)
	})

	api.GET("/rollouts", func(c echo.Context) error {
		var rollouts []Rollout
		if err = db.All(&rollouts); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, rollouts)
	})

	e.POST("/api/rollouts", func(c echo.Context) error {
		var body struct {
			Package string   `json:"package"`
			Devices []string `json:"devices"`
			Running bool     `json:"running"`
		}

		c.Bind(&body)

		for _, uid := range body.Devices {
			var device Device
			fmt.Println(uid)
			if err = db.One("UID", uid, &device); err != nil {
				fmt.Println(err)
				return err
			}

			r, err := device.ActiveRollout(db)
			if err != nil {
				return err
			}

			if r != nil {
				return echo.NewHTTPError(http.StatusNotAcceptable)
			}
		}

		var pkg Package
		if err = db.One("UID", body.Package, &pkg); err != nil {
			return err
		}

		rollout := Rollout{
			Package: body.Package,
			Devices: body.Devices,
			Running: body.Running,
		}

		if rollout.Running {
			rollout.StartedAt = time.Now()
		}

		if err := db.Save(&rollout); err != nil {
			return err
		}

		for _, uid := range rollout.Devices {
			var d Device
			if err = db.One("UID", uid, &d); err != nil {
				continue
			}

			d.Status = "pending"

			if err = db.Save(&d); err != nil {
				fmt.Println(err)
			}
		}

		return c.JSON(http.StatusOK, rollout)
	})

	go func() {
		cmd := exec.Command("npm", "run", "serve", "--", "--port", "1314")
		cmd.Dir = "ui/"
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(e.Start(":1313"))
}
