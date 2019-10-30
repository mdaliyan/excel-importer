package main

import (
	"github.com/mdaliyan/excel-importer/src/lib/database"
)

func main() {

	database.ConnectToElastic()
	database.ConnectToInfluxDB()

}
