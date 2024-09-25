package cmd

import (
	"fmt"

	_ "github.com/fubieliangpu/WorkOrderDeployment/apps"
	initCmd "github.com/fubieliangpu/WorkOrderDeployment/cmd/init"
	"github.com/fubieliangpu/WorkOrderDeployment/cmd/start"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

var RootCmd = &cobra.Command{
	Use:   "wod",
	Short: "wod service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if args[0] == "version" {
				fmt.Println("v1.0.0")
			}
		} else {
			cmd.Help()
		}
	},
}

func Excute() error {
	//初始化需要执行
	cobra.OnInitialize(func() {
		//加载配置
		cobra.CheckErr(conf.LoadConfigFromYaml(configPath))
		//初始化IOC
		cobra.CheckErr(ioc.Controller.Init())
		//初始化Api
		cobra.CheckErr(ioc.Api.Init())
	})
	return RootCmd.Execute()
}

func init() {
	// --config
	RootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "etc/application.yaml", "the service config file")
	// Root --> init
	RootCmd.AddCommand(initCmd.Cmd)
	// Root --> start
	RootCmd.AddCommand(start.Cmd)
}
