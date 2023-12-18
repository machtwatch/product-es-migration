package cmd

import (
	"fmt"
	"product-es-migration/elastic"

	"github.com/spf13/cobra"
)

var createIndexCmd = &cobra.Command{
	Use:   "create-index",
	Short: "Create index based on environment",
	Long:  `Create index based on environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This function will be executed when the "subcommand" is called
		err := elastic.CreateIndex()
		if err != nil {
			fmt.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createIndexCmd)
}
