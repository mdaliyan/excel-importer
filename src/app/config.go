package app

import (
	"github.com/spf13/viper"
	"math/rand"
	"runtime"
	"time"
)

var Config config

type config struct {
	Elasticsearch struct {
		Url, User, Pass string
	}
	InfluxDB struct {
		Url, User, Pass string
	}
}

func init() {
	if err := LoadConfig(); err != nil {
		panic(err)
	}
}

func LoadConfig() (err error) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	conf := viper.New()
	conf.SetConfigName("env")
	conf.SetConfigType("yml")
	conf.AddConfigPath("./")
	conf.AddConfigPath("../")
	conf.AddConfigPath("../../")
	conf.AddConfigPath("../../../")
	conf.AddConfigPath("../../../../")
	if err = conf.ReadInConfig(); err != nil {
		return
	}

	if err = conf.Unmarshal(&Config); err != nil {
		return
	}

	return
}
