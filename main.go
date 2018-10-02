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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd *cobra.Command

func main() {
	rootCmd := &cobra.Command{
		Use: "updatehub-ose-server",
	}

	rootCmd.PersistentFlags().StringP("db", "", "updatehub.db", "Database file")
	rootCmd.PersistentFlags().StringP("username", "", "admin", "Admin username")
	rootCmd.PersistentFlags().StringP("password", "", "admin", "Admin password")
	rootCmd.PersistentFlags().IntP("port", "", 8080, "Port")

	viper.SetEnvPrefix("")
	viper.BindEnv("db")
	viper.BindEnv("username")
	viper.BindEnv("password")
	viper.BindEnv("port")
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	db, err := storm.Open(viper.GetString("db"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/ui/")
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

		if login.Username == viper.GetString("username") && login.Password == viper.GetString("password") {
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

	log.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("port"))))
}
