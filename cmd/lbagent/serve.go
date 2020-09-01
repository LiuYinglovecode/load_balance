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

	"code.htres.cn/casicloud/alb/agent"
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
		config := parseConfig()
		// init log
		err := agent.InitLogger(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// create and load store
		store, err := agent.NewLBPolicyStore(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = store.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// create and start informer
		informer, err := agent.NewInformer(config, store)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = informer.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		rpcServer, err := agent.NewRPCServer(config, store)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		quit := make(chan int)
		go func() {
			err := rpcServer.Run()
			if err != nil {
				fmt.Println(err)
			}
			quit <- 0
		}()
		// processing...
		<-quit
		rpcServer.Stop()
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

	serveCmd.Flags().String("rpc", "0.0.0.0:6767", "Address to run rpc server")
	viper.BindPFlags(serveCmd.Flags())

	dir, err := os.Getwd()
	if err == nil {
		viper.SetDefault("WorkDir", dir)
	}

	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("RPC", "0.0.0.0:6767")

	viper.SetDefault("AuditLogPath", filepath.Join(dir, "logs", "agent_audit.log"))
	viper.SetDefault("SysLogPath", filepath.Join(dir, "logs", "agent_sys.log"))

	viper.SetDefault("Role", "master")
}

// parseConfig 读取配置文件
func parseConfig() *agent.Config {
	rpc := viper.GetString("RPC")
	path := viper.GetString("WorkDir")
	workDir, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	id := viper.GetInt64("AgentID")
	path = viper.GetString("AuditLogPath")
	auditLogpath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path = viper.GetString("SysLogPath")
	sysLogPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logLevel := viper.GetString("LogLevel")
	lbmcURL := viper.GetString("LBMCUrl")

	endopoints := viper.GetStringSlice("Endpoints")

	ca := viper.GetString("EtcdTLS.ca")
	cert := viper.GetString("EtcdTLS.cert")
	key := viper.GetString("EtcdTLS.key")

	// 读入keepalived 配置
	var kaCfg *agent.KeepalivedConfig
	on := viper.GetBool("Keepalived.TurnOn")
	if on == true {
		inet := viper.GetString("Keepalived.INet")
		vrID := viper.GetString("Keepalived.VirutalRouterID")
		stat := viper.GetString("Keepalived.State")
		pri := viper.GetInt("Keepalived.Priority")
		srcIP := viper.GetString("Keepalived.UnicastSrcIP")
		peers := viper.GetStringSlice("Keepalived.UnicastPeer")
		vIP := viper.GetString("Keepalived.VirtualIP")
		kaCfg = &agent.KeepalivedConfig{
			INet:            inet,
			VirutalRouterID: vrID,
			State:           stat,
			Priority:        pri,
			UnicastSrcIP:    srcIP,
			UnicastPeer:     peers,
			VirtualIP:       vIP,
		}
	}

	config := &agent.Config{
		RPC:           rpc,
		WorkDir:       workDir,
		LBMCUrl:       lbmcURL,
		LogLevel:      logLevel,
		AuditLogPath:  auditLogpath,
		SysLogPath:    sysLogPath,
		AgentID:       id,
		Endpoints:     endopoints,
		EtcdCAPath:    ca,
		EtcdCertPath:  cert,
		EtcdKeyPath:   key,
		KeepalivedCfg: kaCfg,
	}
	return config
}
