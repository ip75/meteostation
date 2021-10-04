/*
Copyright © 2021 NAME HERE <Igor.Polovykh@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ip75/meteostation/config"
	"github.com/ip75/meteostation/server"
	"github.com/ip75/meteostation/storage"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meteostation",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "run root:", cmd.Use)

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		storage.PG.PlayMigrations()

		// launch routine to receive data from redis and put them to postgres

		storage.RDB.Init()
		storage.PG.Init()

		sdChannel := make(chan []storage.SensorDataDatabase)

		go func() {
			server.Start()
		}()

		go func() {
			for {
				storage.PG.StoreSensorData(<-sdChannel)
			}
		}()

		for {
			select {
			case sdChannel <- storage.RDB.Pull():
			case <-term:
				log.Print("Got signal SIGTERM and exit", term)
				os.Exit(0)
				//default:
				//	log.Print("waiting for data from channel...")
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.meteostation.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType("json")

	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".meteostation" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/")
		viper.SetConfigName(".meteostation") // name of config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file: ", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config.C)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to decode into struct, ", err)
	}
}
