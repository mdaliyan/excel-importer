package elastic

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd *cobra.Command

var url, user, pass string

func init() {

	Cmd = &cobra.Command{
		Use:   "elastic",
		Short: "import data to elasticsearch",
		Long: `echo things multiple times back to the user by providing a count and a string.`,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(cmd.Flags())
			fmt.Println(args)

		},
	}



}
