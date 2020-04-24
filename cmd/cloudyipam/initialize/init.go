/*
Copyright © 2020 John Slee <john@sleefamily.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package initialize

import (
	"fmt"
	"os"

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/shared"
	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/sqltext"
	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
	"github.com/lib/pq"
	"github.com/spf13/cobra"
)

func createDB() error {
	config := shared.DatabaseConfigFromViper()
	dbName := config.Database
	// for init we shouldn't have a database yet (and if we do,
	// we should abort)
	config.Database = ""
	ipam, err := cloudyipam.NewCloudyIPAM(config.DSN())
	if err != nil {
		return err
	}
	conn := ipam.DB()
	defer conn.Close()
	// can 't use Postgres native ordinal markers ($1) etc here because it's
	// an identifier rather than a value. See
	// https://stackoverflow.com/questions/37448982/what-am-i-not-getting-about-go-sql-query-with-variables/37449128
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
	return err
}

func createObjects() error {
	config := shared.DatabaseConfigFromViper()
	ipam, err := cloudyipam.NewCloudyIPAM(config.DSN())
	if err != nil {
		return err
	}
	conn := ipam.DB()
	defer conn.Close()
	transaction, err := conn.Begin()
	if err != nil {
		return err
	}
	_, err = transaction.Exec(sqltext.Text())
	if err != nil {
		transaction.Rollback()
		return err
	}
	err = transaction.Commit()
	return err
}

// Cmd represents the init command
var Cmd = &cobra.Command{
	Use:   "initialize",
	Short: "Deploy CloudyIPAM database structure and functions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := createDB(); err != nil {
			fmt.Printf("error creating database: %v\n", err)
			os.Exit(1)
		}
		if err := createObjects(); err != nil {
			fmt.Printf("error creating database objects: %v\n", err)
			os.Exit(1)
		}
	},
}
