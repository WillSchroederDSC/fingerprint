package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willschroeder/fingerprint/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the fingerprint server",
	Run: func(cmd *cobra.Command, args []string) {
		server.NewServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
