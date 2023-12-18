package cmd

import (
	"product-es-migration/elastic"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate elastic db ",
	Long:  `Migrate elastic db `,
	Run: func(cmd *cobra.Command, args []string) {
		elastic.Migrate()

	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
