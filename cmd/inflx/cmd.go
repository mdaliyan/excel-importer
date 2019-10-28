package inflx

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd *cobra.Command

var url, user, pass string

func init() {

	Cmd = &cobra.Command{
		Use:   "influx",
		Short: "import data to influxDB",
		Long:  `echo things multiple times back to the user by providing a count and a string.`,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(cmd.Flags())
			fmt.Println(args)

		},
	}

	Cmd.Flags().StringVarP(&url, "url", "", "127.0.0.1:9200", "influx url and port")
	Cmd.Flags().StringVarP(&user, "user", "", "", "influx username")
	Cmd.Flags().StringVarP(&pass, "pass", "", "", "influx password")
	_ = Cmd.MarkFlagRequired("url")
}
