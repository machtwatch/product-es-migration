package cmd

import (
	"fmt"
	"product-es-migration/elastic"

	"github.com/spf13/cobra"
)

var deleteIndexCmd = &cobra.Command{
	Use:   "delete-index",
	Short: "Delete index based on environment ",
	Long:  `Delete index based on environment. `,
	Run: func(cmd *cobra.Command, args []string) {
		err := elastic.DeleteIndex()
		if err != nil {
			fmt.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteIndexCmd)
}
