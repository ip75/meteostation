/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

const configTemplate = `
{
    "general": {
        "log_level": 1,
        "log_to_syslog": false,
        "sensor_data_pool_size": 60,
        "grpc_service_port": 65000,
        "http_port": 64000,
        "web_client_dir": "./static"
    },

    "postgresql": {
        "dsn" : "postgres://meteostation:meteostation@pg:5432/meteostation?sslmode=disable",
        "max_open_connections" : 1000,
        "max_idle_connections" : 2000
    },

    "redis": {
        "url" : "redis:6379",
        "password" : "",
        "database" : 0,
        "queue" : "meteostation:bmp280"
    }
}
`

// getconfigCmd represents the getconfig command
var getconfigCmd = &cobra.Command{
	Use:   "getconfig",
	Short: "Put standard cofnig to stdout",
	Long:  `Put standard cofnig to stdout to make config file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(configTemplate)
	},
}

func init() {
	rootCmd.AddCommand(getconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
