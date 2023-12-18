package cmd

import (
	"fmt"
	"product-es-migration/elastic"

	"github.com/spf13/cobra"
)

var checkIndexExistCmd = &cobra.Command{
	Use:   "check-index",
	Short: "Check index is exist or not",
	Long:  `Create index is exist or not.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := elastic.CheckIndex()
		if err != nil {
			fmt.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkIndexExistCmd)
}
