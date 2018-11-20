package cmd

import (
	"github.com/willschroeder/fingerprint/pkg/client"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs a new client",
	Run: func(cmd *cobra.Command, args []string) {
		client.NewClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
