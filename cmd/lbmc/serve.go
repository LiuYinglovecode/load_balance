// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"code.htres.cn/casicloud/alb/center"
	"code.htres.cn/casicloud/alb/center/common"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		config, err := parseConfig()
		if err != nil {
			fmt.Println("error load config file, reason: ", err)
		}

		common.GlobalConfig = *config

		err = common.InitLogger(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = center.StartServer()
		if err != nil {
			fmt.Println("error start server, reason: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().String("port", "8080", "Port to run lbmc server")
	viper.BindPFlags(serveCmd.Flags())

	dir, err := os.Getwd()
	if err == nil {
		viper.SetDefault("WorkDir", dir)
	}

	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("Port", "8080")

	viper.SetDefault("AuditLogPath", filepath.Join(dir, "logs", "audit.log"))
	viper.SetDefault("SysLogPath", filepath.Join(dir, "logs", "sys.log"))

}

// parseConfig 返回config
func parseConfig() (*common.Config, error) {
	var c common.Config

	c.WorkDir = viper.GetString("WorkDir")

	c.LogLevel = viper.GetInt("LogLevel")
	c.AuditLogPath = viper.GetString("AuditLogPath")
	c.SysLogPath = viper.GetString("SysLogPath")

	timeout := viper.GetInt64("AgentTimeout")
	c.AgentTimeout = time.Duration(timeout) * time.Millisecond

	c.Port = viper.GetInt("Port")
	c.DBArgs = viper.GetString("DBArgs")
	c.Dialect = viper.GetString("Dialect")
	c.EtcdEndpoints = viper.GetStringSlice("EtcdEndpoints")

	c.EtcdCAPath = viper.GetString("EtcdTLS.ca")
	c.EtcdCertPath = viper.GetString("EtcdTLS.cert")
	c.EtcdKeyPath = viper.GetString("EtcdTLS.key")
	return &c, nil
}
