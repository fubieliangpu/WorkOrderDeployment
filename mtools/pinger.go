package mtools

import (
	"log"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/go-ping/ping"
)

//定义一个ping测工具

func PingTool(pingnum int, ipaddr string) (bool, error) {
	//新建一个pinger
	pinger, err := ping.NewPinger(ipaddr)
	if err != nil {
		log.Fatal(err)
	}
	pinger.Count = pingnum
	//启用ping，未完成时会阻塞该进程
	err = pinger.Run()
	if err != nil {
		log.Fatal(err)
	}
	stats := pinger.Statistics()
	if stats.PacketLoss != 0 {
		//后面替换为自定义异常类型
		return false, exception.ErrPingLoss("Ping %d packets,loss %d packets", stats.PacketsSent, stats.PacketsSent-stats.PacketsRecv)
	}
	//后面替换为自定义异常类型
	return true, nil
}
