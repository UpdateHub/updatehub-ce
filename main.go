package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/asdine/storm"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/packr"
	"github.com/gustavosbarreto/updatehub-server/api/agentapi"
	"github.com/gustavosbarreto/updatehub-server/api/webapi"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/ui/")
	})

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

	agentApi := agentapi.NewAgentAPI(db)
	e.POST(agentapi.GetRolloutForDeviceUrl, agentApi.GetRolloutForDevice)
	e.POST(agentapi.ReportDeviceStateUrl, agentApi.ReportDeviceState)
	e.GET(agentapi.GetObjectFromPackageUrl, agentApi.GetObjectFromPackage)

	api := e.Group("/api")
	api.Use(middleware.JWT([]byte("secret")))

	devicesEndpoint := webapi.NewDevicesAPI(db)
	api.GET(webapi.GetAllDevicesUrl, devicesEndpoint.GetAllDevices)
	api.GET(webapi.GetDeviceUrl, devicesEndpoint.GetDevice)

	packagesEndpoint := webapi.NewPackagesAPI(db)
	api.GET(webapi.GetAllPackagesUrl, packagesEndpoint.GetAllPackages)
	api.GET(webapi.GetPackageUrl, packagesEndpoint.GetPackage)
	api.POST(webapi.UploadPackageUrl, packagesEndpoint.UploadPackage)

	rolloutsEndpoint := webapi.NewRolloutsAPI(db)
	api.GET(webapi.GetAllRolloutsUrl, rolloutsEndpoint.GetAllRollouts)
	api.GET(webapi.GetRolloutUrl, rolloutsEndpoint.GetRollout)
	api.GET(webapi.GetRolloutStatisticsUrl, rolloutsEndpoint.GetRolloutStatistics)
	api.POST(webapi.CreateRolloutUrl, rolloutsEndpoint.CreateRollout)

	if os.Getenv("ENV") == "production" {
		box := packr.NewBox("./ui/dist")
		handler := echo.WrapHandler(http.StripPrefix("/ui", http.FileServer(box)))
		e.GET("/ui/*", handler)
		e.GET("/ui", handler)
	} else {
		ui, _ := url.Parse("http://localhost:1314/")
		e.Group("/ui", middleware.Proxy(middleware.NewRoundRobinBalancer(
			[]*middleware.ProxyTarget{{URL: ui}},
		)))

		go func() {
			cmd := exec.Command("npm", "run", "serve", "--", "--port", "1314")
			cmd.Dir = "ui/"
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	log.Fatal(e.Start(":1313"))
}
