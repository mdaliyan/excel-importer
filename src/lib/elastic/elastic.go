package rbmitems

import (
	"context"
	"strings"
	"time"

	"gitlab.com/sanjagh/core/src/lib/app"
	fn "gitlab.com/sanjagh/core/src/lib/functions"

	"github.com/olivere/elastic"

	"fmt"
	"github.com/mdaliyan/gabs"
)

var Client *elastic.Client
var Type = "campaign"
// "http://10.135.25.230:9200",

func ConnectDB() {
	if Client == nil {
		RefreshConnection()
	}
}

var RunningCampaignsQuery = elastic.NewBoolQuery().
	Must(
		elastic.NewTermQuery("enabled", true),
		elastic.NewTermQuery("approved", true),
		elastic.NewTermQuery("trashed", false),
	).
	Should(
		elastic.NewTermQuery("limits_daily_budget_limit", 0),
		elastic.NewScriptQuery(
			elastic.NewScriptInline("doc['limits_daily_budget_spent'].value < doc['limits_daily_budget_limit'].value").
				Lang("painless"),
		),
	)

func GetClient() *elastic.Client {
	ConnectDB()
	return Client
}

func GetGetService() *elastic.GetService {
	return GetClient().Get().Index(app.Config.RbmItems.Database)
}

func GetIndexService() *elastic.IndexService {
	return GetClient().Index().Index(app.Config.RbmItems.Database)
}

func GetSearchService() *elastic.SearchService {
	return GetClient().Search(app.Config.RbmItems.Database)
}

func GetUpdateService(TypeSlice ...string) *elastic.UpdateService {
	c := GetClient().Update().Index(app.Config.RbmItems.Database)
	if len(TypeSlice) != 0 {
		return c.Type(TypeSlice[0])
	}
	return c
}

func GetUpdateByQueryService() *elastic.UpdateByQueryService {
	return GetClient().UpdateByQuery().Index(app.Config.RbmItems.Database)
}

func GetDeleteService(TypeSlice ...string) *elastic.DeleteService {
	var T = Type
	if len(TypeSlice) != 0 {
		T = TypeSlice[0]
	}
	return GetClient().Delete().Index(app.Config.RbmItems.Database).Type(T)
}

func RefreshConnection() (err error) {
	fmt.Println(app.Config.RbmItems.Hosts)
	// Get a client to the local Elasticsearch instance.
	if app.Config.RbmItems.Username != "" {
		Client, err = elastic.NewSimpleClient(
			elastic.SetURL(app.Config.RbmItems.Hosts...),
			elastic.SetBasicAuth(app.Config.RbmItems.Username, app.Config.RbmItems.Password),
			elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetMaxRetries(60),
		)
	} else {
		Client, err = elastic.NewSimpleClient(elastic.SetURL(app.Config.RbmItems.Hosts...))
	}

	return
}

func MultiGet(index string, typ string, Ids []string) (searchResult *elastic.MgetResponse, err error) {
	query := gabs.New()
	query.Set(Ids, "ids")

	if index == "" {
		index = app.Config.RbmItems.Database
	}

	var retry int8

	elastic.NewSearchSource()
	ConnectDB()
	if index == "" {
		index = app.Config.RbmItems.Database
	}
	engine := GetClient().MultiGet()

	for _, id := range Ids {
		itemInfo := elastic.NewMultiGetItem().Index(index).Id(id)
		if typ != "" {
			itemInfo.Type(typ)
		}
		engine.Add(itemInfo)
	}
	for searchResult, err = engine.Do(context.Background()); err != nil; {
		if retry++; retry >= 10 {
			fn.LogError("RbmItems MultiGet() " + err.Error())
			return
		}

		time.Sleep(time.Second * 10)
		ConnectDB()
	}
	return
}

func Query(index string, typ string, query string) (searchResult *elastic.SearchResult, err error) {

	var retry int8

	elastic.NewSearchSource()
	ConnectDB()
	if index == "" {
		index = app.Config.RbmItems.Database
	}
	engine := GetClient().Search().Index(index)
	if typ != "" {
		engine.Type(typ)
	}
	for searchResult, err = engine.Source(query).Do(context.Background()); err != nil; {
		if retry++; retry >= 10 {
			fn.LogError("RbmItems Query() " + err.Error())
			return
		}
		fn.LogNotice(err)
		time.Sleep(time.Second * 10)
		ConnectDB()
	}
	return
}

func DeleteByQuery(index string, typ string, Query elastic.Query) (searchResult *elastic.BulkIndexByScrollResponse, err error) {

	var retry int8
	elastic.NewSearchSource()
	if index == "" {
		index = app.Config.RbmItems.Database
	}
	for searchResult, err = GetClient().DeleteByQuery().
		Index(index).Type(typ).Query(Query).
		Do(context.Background()); err != nil; {

		if retry++; retry >= 10 {
			fn.LogError("RbmItems DeleteByQuery()" + err.Error())
			return
		}

		time.Sleep(time.Second * 10)
		ConnectDB()
	}
	return
}

func DoSearch(Query *elastic.SearchService) (searchResult *elastic.SearchResult, err error) {
	var retry int8
	for searchResult, err = Query.Do(context.Background()); err != nil; {
		if strings.Contains(err.Error(), "no such index") {
			err = nil
			searchResult = &elastic.SearchResult{Hits: &elastic.SearchHits{}}
			return
		}
		if retry++; retry >= 10 {
			fn.LogError("RbmItems Query() " + err.Error())
			return
		}
		fn.LogNotice(err)
		time.Sleep(time.Second * 10)
		ConnectDB()
	}
	return
}
