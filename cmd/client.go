package cmd

import (
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs a new client",
	Run: func(cmd *cobra.Command, args []string) {
		println("client")
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs a new client",
	Run: func(cmd *cobra.Command, args []string) {
		println("client test")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(testCmd)
}
