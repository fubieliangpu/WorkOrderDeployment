package test

import (
	_ "github.com/fubieliangpu/WorkOrderDeployment/apps"
	"github.com/fubieliangpu/WorkOrderDeployment/conf"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/spf13/cobra"
)

func DevelopmentSetup() {
	// 加载配置, 单元测试 通过环境变量读取, vscode 传递进来的
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}
	//init ioc
	cobra.CheckErr(ioc.Controller.Init())
}
