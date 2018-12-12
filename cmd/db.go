package cmd

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/willschroeder/fingerprint/pkg/db"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("db root called, add -h to see available commands")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply any un-applied migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db.MigrateUp()
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db.MigrateDown()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(migrateCmd, rollbackCmd)
}
