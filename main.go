// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/UpdateHub/updatehub-ce/api/agentapi"
	"github.com/UpdateHub/updatehub-ce/api/webapi"
	"github.com/asdine/storm"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})

	rootCmd := &cobra.Command{
		Use: "updatehub-ce",
		Run: execute,
	}

	rootCmd.PersistentFlags().StringP("db", "", "updatehub.db", "Database file")
	rootCmd.PersistentFlags().StringP("username", "", "admin", "Admin username")
	rootCmd.PersistentFlags().StringP("password", "", "admin", "Admin password")
	rootCmd.PersistentFlags().IntP("http", "", 8080, "HTTP listen address")
	rootCmd.PersistentFlags().StringP("dir", "", "./", "Packages storage dir")
	rootCmd.PersistentFlags().IntP("coap", "", 5683, "Coap server listen port")
	rootCmd.PersistentFlags().StringP("secret", "", "secret", "JWT secret key")

	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("http", rootCmd.PersistentFlags().Lookup("http"))
	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("coap", rootCmd.PersistentFlags().Lookup("coap"))
	viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func execute(cmd *cobra.Command, args []string) {
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

			t, err := token.SignedString([]byte(viper.GetString("secret")))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
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
	api.Use(middleware.JWT([]byte(viper.GetString("secret"))))

	devicesEndpoint := webapi.NewDevicesAPI(db)
	api.GET(webapi.GetAllDevicesUrl, devicesEndpoint.GetAllDevices)
	api.GET(webapi.GetDeviceUrl, devicesEndpoint.GetDevice)
	api.GET(webapi.GetDeviceRolloutReportsUrl, devicesEndpoint.GetDeviceRolloutReports)

	packagesEndpoint := webapi.NewPackagesAPI(db, viper.GetString("dir"))
	api.GET(webapi.GetAllPackagesUrl, packagesEndpoint.GetAllPackages)
	api.GET(webapi.GetPackageUrl, packagesEndpoint.GetPackage)
	api.POST(webapi.UploadPackageUrl, packagesEndpoint.UploadPackage)

	rolloutsEndpoint := webapi.NewRolloutsAPI(db)
	api.GET(webapi.GetAllRolloutsUrl, rolloutsEndpoint.GetAllRollouts)
	api.GET(webapi.GetRolloutUrl, rolloutsEndpoint.GetRollout)
	api.GET(webapi.GetRolloutStatisticsUrl, rolloutsEndpoint.GetRolloutStatistics)
	api.GET(webapi.GetRolloutDevicesUrl, rolloutsEndpoint.GetRolloutDevices)
	api.POST(webapi.CreateRolloutUrl, rolloutsEndpoint.CreateRollout)
	api.PUT(webapi.StopRolloutUrl, rolloutsEndpoint.StopRollout)

	if os.Getenv("ENV") == "development" {
		ui, _ := url.Parse("http://localhost:1314/")
		e.Group("/ui", middleware.Proxy(middleware.NewRoundRobinBalancer(
			[]*middleware.ProxyTarget{{URL: ui}},
		)))

		tmpfile, err := ioutil.TempFile("", "gopid.*.js")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte("process.kill(process.env.GOPID, 'SIGUSR1')")); err != nil {
			log.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		logrus.Info("Starting Vue development server...")

		go func() {
			cmd := exec.Command("npm", "run", "serve", "--", "--open", "--port", "1314")
			cmd.Dir = "ui/"
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("BROWSER=%s", tmpfile.Name()))
			cmd.Env = append(cmd.Env, fmt.Sprintf("GOPID=%d", os.Getpid()))
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)
		_ = <-sigs

		os.Remove(tmpfile.Name())
	} else {
		box := packr.NewBox("./ui/dist")
		handler := echo.WrapHandler(http.StripPrefix("/ui", http.FileServer(box)))
		e.GET("/ui/*", handler)
		e.GET("/ui", handler)
	}

	go func() {
		log.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("http"))))
	}()

	go func() {
		log.Fatal(startCoapServer(viper.GetInt("coap"), viper.GetInt("http")))
	}()

	select {}
}
