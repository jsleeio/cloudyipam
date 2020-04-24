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
package root

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/initialize"
	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/ping"
	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/subnet"
	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/zone"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudyipam",
	Short: "Commandline interface to CloudyIPAM service",
	Long:  "https://github.com/jsleeio/cloudyipam",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// config everywhere else should use viper, so don't leak these to package scope
	var user, name, host, password, port string
	var tls bool
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudyipam.yaml)")
	rootCmd.PersistentFlags().StringVar(&user, "db-user", "cloudyipam", "database username")
	rootCmd.PersistentFlags().StringVar(&password, "db-password", "", "database password")
	rootCmd.PersistentFlags().StringVar(&name, "db-name", "cloudyipam", "database name")
	rootCmd.PersistentFlags().StringVar(&port, "db-port", "", "database port")
	rootCmd.PersistentFlags().StringVar(&host, "db-host", "", "database hostname")
	rootCmd.PersistentFlags().BoolVar(&tls, "db-tls", false, "Enable TLS when connecting to database")
	rootCmd.AddCommand(initialize.Cmd)
	rootCmd.AddCommand(ping.Cmd)
	rootCmd.AddCommand(zone.Cmd)
	rootCmd.AddCommand(subnet.Cmd)
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db-tls", rootCmd.PersistentFlags().Lookup("db-tls"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cloudyipam" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cloudyipam")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
