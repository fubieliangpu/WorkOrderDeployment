package start

import (
	"os"

	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/spf13/cobra"
)

var (
	testParam string
)

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "start Workorderdeployment api server",
	Run: func(cmd *cobra.Command, args []string) {
		//加载配置,首先从环境变量加载
		configPath := os.Getenv("WODCONFIG_PATH")
		if configPath == "" {
			configPath = "etc/application.yaml"
		}
		//启动服务
		cobra.CheckErr(conf.C().Application.Start())
	},
}

func init() {
	Cmd.Flags().StringVarP(&testParam, "test", "t", "test", "test flag")
}
