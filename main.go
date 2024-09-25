package main

import "github.com/fubieliangpu/WorkOrderDeployment/cmd"

func main() {
	if err := cmd.Excute(); err != nil {
		panic(err)
	}
}
