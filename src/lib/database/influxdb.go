package database

import (
	"fmt"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/mdaliyan/excel-importer/src/app"
)

var influxClient *influx.Client

func InfluxDClient() (*influx.Client, bool) {
	if influxClient == nil {
		return nil, false
	}
	ConnectToInfluxDB()
	return influxClient, true
}

func ConnectToInfluxDB() {
	if influxClient == nil {
		const connErr = "connection to influxDB on"
		defer fmt.Println(connErr, app.Config.InfluxDB.Url, ", error:", refreshInfluxDBConnection())
	}
}

func refreshInfluxDBConnection() (err error) {
	if influxClient != nil {
		return
	}
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     app.Config.InfluxDB.Url,
		Username: app.Config.InfluxDB.User,
		Password: app.Config.InfluxDB.Pass,
	})
	influxClient = &c
	return err
}
