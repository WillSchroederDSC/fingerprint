// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	"github.com/willschroeder/fingerprint/pkg/db"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("db called")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		absPath, _ := filepath.Abs("./migrations")

		migrations := &migrate.FileMigrationSource{
			Dir: absPath,
		}

		db := db.ConnectToDatabase()
		defer db.Conn.Close()

		n, err := migrate.Exec(db.Conn, "postgres", migrations, migrate.Up)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Applied %d migrations\n", n)
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		absPath, _ := filepath.Abs("./migrations")

		migrations := &migrate.FileMigrationSource{
			Dir: absPath,
		}

		db := db.ConnectToDatabase()
		defer db.Conn.Close()

		n, err := migrate.Exec(db.Conn, "postgres", migrations, migrate.Down)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Rolled back %d migrations\n", n)
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(migrateCmd, rollbackCmd)
}
