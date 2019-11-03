package main

import (
	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/labstack/echo/v4"
	"github.com/mdaliyan/excel-importer/src/lib/database"
	"net/http"
)

func main() {

	database.ConnectToInfluxDB()

	Router := echo.New()

	Router.GET("/", func(c echo.Context) error {

		q := client.NewQuery("SELECT count(value) FROM cpu_load", "mydb", "")

		r, err := database.InfluxDClient().Query(q)
		if err != nil {
			return c.String(http.StatusNotAcceptable, err.Error())
		}

		return c.JSONPretty(http.StatusOK, r, "   ")
	})

	Router.Logger.Fatal(Router.Start(":4044"))

}
