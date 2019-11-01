package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mdaliyan/excel-importer/src/lib/database"
	"net/http"
)

func main() {

	database.ConnectToElastic()
	database.ConnectToInfluxDB()

	fmt.Println("Hello World")

	Router := echo.New()

	Router.GET("/", func(c echo.Context) error {
		el, ok := database.ElasticSearchClient()

		if !ok {
			return c.String(http.StatusGone, "elasticsearch is gone")
		}

		r, err := el.NodesInfo().Do(context.Background())

		if err != nil {
			return c.String(http.StatusGone, "elasticsearch is gone")
		}

		return c.JSONPretty(http.StatusOK, r, "   ")
	})

	Router.Logger.Fatal(Router.Start(":4040"))

}
