package cmd

import (
	"github.com/willschroeder/fingerprint/pkg/db"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans any expired records from the database",
	Run: func(cmd *cobra.Command, args []string) {
		db.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
