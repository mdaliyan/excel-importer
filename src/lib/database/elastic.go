package database

import (
	"fmt"
	"github.com/mdaliyan/excel-importer/src/app"

	"github.com/olivere/elastic/v7"
)

var elasticClient *elastic.Client

func ElasticSearchClient() (*elastic.Client, bool) {
	if elasticClient == nil {
		return nil, false
	}
	ConnectToElastic()
	return elasticClient, true
}

func ConnectToElastic() {
	if influxClient != nil {
		return
	}
	const connErr = "connection to elasticsearch on"
	fmt.Println(connErr, app.Config.Elasticsearch.Url, ", error:", refreshElasticConnection())
}

func refreshElasticConnection() (err error) {

	// Get a client to the local Elasticsearch instance.
	if app.Config.Elasticsearch.User == "" {
		elasticClient, err = elastic.NewSimpleClient(elastic.SetURL(app.Config.Elasticsearch.Url))

		return
	}

	elasticClient, err = elastic.NewSimpleClient(
		elastic.SetURL(app.Config.Elasticsearch.Url),
		elastic.SetBasicAuth(app.Config.Elasticsearch.User, app.Config.Elasticsearch.Pass),
		elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetMaxRetries(60),
	)

	return
}
