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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ip75/meteostation/config"
	"github.com/ip75/meteostation/server"
	"github.com/ip75/meteostation/storage"
	"github.com/sirupsen/logrus"
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
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:             false,
			FullTimestamp:             true,
			ForceColors:               true,
			DisableTimestamp:          false,
			EnvironmentOverrideColors: false,
		})

		logrus.SetOutput(os.Stdout)

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		storage.PG.PlayMigrations()

		// launch routine to receive data from redis and put them to postgres
		storage.RDB.Init()
		if err := storage.PG.Init(); err != nil {
			logrus.Error("initialize postgres failed: ", err)
		}

		sdChannel := make(chan []storage.SensorDataDatabase)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// start backend
		logrus.Print("start backend...")
		go func() { server.StartHosting(ctx) }()
		go func() { server.StartGrpcService(ctx) }()
		go func() { server.StartGateway(ctx) }()

		// routine to catch data from channel and drop it to database
		go func() {
			for {
				if err := storage.PG.StoreSensorData(<-sdChannel); err != nil {
					logrus.Error("store sensor data error: ", err)
				}
			}
		}()

		// routine to pull sensor data from redis to channel
		logrus.Print("start receiving data from sensors...")
		for {
			select {
			case sdChannel <- storage.RDB.Pull():
			case <-term:
				logrus.Info("Got signal SIGTERM and exit", term)
				os.Exit(0)
				// default:
				//	logrus.Print("waiting for data from channel...")
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
		logrus.Info("Using config file: ", viper.ConfigFileUsed())
	} else {
		logrus.Error("config failed: ", err)
	}

	err := viper.Unmarshal(&config.C)
	if err != nil {
		logrus.Error("unmarshal config file: ", err)
	}
}
