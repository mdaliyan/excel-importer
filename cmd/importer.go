package main

import (
	"fmt"
	"github.com/mdaliyan/excel-importer/cmd/elastic"
	"github.com/mdaliyan/excel-importer/cmd/inflx"
	"github.com/spf13/cobra"
)

var listen string
var DBHost, user, pass string

func main() {

	var influxDB *cobra.Command

	var elacticDB = &cobra.Command{
		Use:   "elastic",
		Short: "import data to elasticsearch",
		Long:  `echo things multiple times back to the user by providing a count and a string.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Flags())
			fmt.Println(args)
		},
	}
	elacticDB.Flags().StringVarP(&listen, "listen", "", ":8099", "listen on http host")
	elacticDB.Flags().StringVarP(&listen, "url", "", "127.0.0.1:9200", "elastic url and port")
	elacticDB.Flags().StringVarP(&user, "user", "", "", "elastic username")
	elacticDB.Flags().StringVarP(&pass, "pass", "", "", "elastic password")
	_ = elacticDB.MarkFlagRequired("listen")
	_ = elacticDB.MarkFlagRequired("host")

	var rootCmd = &cobra.Command{
		Use:  "importer",
		Args: cobra.MinimumNArgs(1),
	}

	rootCmd.AddCommand(elastic.Cmd, inflx.Cmd)

	rootCmd.Execute()
}
