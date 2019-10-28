package inflxdb

import (
	"errors"
	"fmt"
	"net/url"

	influx "github.com/influxdata/influxdb1-client/v2"
)

var Client influx.Client

var (
	MyDB     = app.Config.NewStatstics.Database
	Host     = app.Config.NewStatstics.Host
	Username = app.Config.NewStatstics.Username
	Password = app.Config.NewStatstics.Password
)

func GetNewClient() influx.Client {
	var err error
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     Host,
		Username: Username,
		Password: Password,
	})
	if err != nil {
		panic(errors.New("influx connection: " + err.Error()))
	}

	return c
}

func init() {
	var err error
	Client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     Host,
		Username: Username,
		Password: Password,
	})
	if err != nil {
		panic(errors.New("influx connection: " + err.Error()))
	}
}

func NewBatchPoints() (influx.BatchPoints, error) {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: MyDB})
	if err != nil {
		fmt.Println(err)
	}
	return bp, err
}

func GenerateTags(permutation *Permutation, item *Item, clientInfo *ClientInfo) (tags map[string]string) {
	tags = map[string]string{
		"publisher.user_uuid":      permutation.UserUUID,
		"publisher.id":             permutation.PublisherId,
		"publisher.position.id":    permutation.PositionId,
		"publisher.permutation.id": permutation.ID,

		"advertiser.user_uuid": item.GetString("user_uuid"),

		// "client.device": govert.String(clientInfo.IsMobile),

		"item.type":   item.GetString("type"), // campaign - anything else
		"item.module": permutation.Type,             // rbm - cbm - notif
	}

	tags["advertiser.id"] = item.GetString("advertiser_id")

	return
}

func NewStatisticPoint(TAGS map[string]string, fields map[string]interface{}, ) (*influx.Point, error) {
	t := fn.Now()
	fields["timestamp"] = t.UnixNano()
	return influx.NewPoint("statistics", TAGS, fields, t)
}

func NewBotPoint(ctx *fasthttp.RequestCtx, botKey string, fields map[string]interface{}, ) (*influx.Point, error) {
	t := fn.Now()
	fields["timestamp"] = t.UnixNano()
	fields["agent"] = string(ctx.UserAgent())
	tags := map[string]string{"key": botKey,}
	if URL, err := url.Parse(string(ctx.Referer())); err == nil {
		tags["publisher_domain"] = URL.Host
	}
	return influx.NewPoint("bots", tags, fields, t)
}
