package internet

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/mtools"
)

const (
	AppName = "internet"
)

type Service interface {
	//冲突检测
	ConflictCheck(context.Context, *DeploymentNetworkProductRequest) (ConfigConflictStatus, error)
	//业务配置下发
	ConfigDeployment(context.Context, *DeploymentNetworkProductRequest) (*NetProd, error)
	//业务配置回收
	ConfigRevoke(context.Context, *UndoDeviceConfigRequest) (*NetProd, error)
}

func (d *DeploymentNetworkProductRequest) CheckVRRP(brand common.Brand, device *rcdevice.Device) error {

	cfi := rcdevice.NewConfigInfo()
	cfi.UserInfo = rcdevice.NewDeviceUserInfo()
	cfi.UserInfo, _ = rcdevice.LoadUsernmPasswdFromYaml("user.yaml", cfi.UserInfo)

	if d.AccessDeviceLevel == common.CORE {
		//判断查找到的第一台核心的设备品牌
		if brand == common.H3C {
			//首先查看是否存在对应运营商的VPN-INSTANCE，如果存在则查路由表时需要带上对应参数
			if err := mtools.CommandGenerator("H3CDisCommand.txt", "display ip vpn-instance\nexit\n"); err != nil {
				return err
			}
			cfi.Ip = device.ServerAddr
			cfi.Port = device.Port
			cfi.Configfile = "H3CDisCommand.txt"
			cfi.Recordfile = "H3CCoreVRRPVpn.txt"
			rcdevice.SshConfigTool(cfi)
			if err := mtools.Regexper(cfi.Recordfile, d.Operators); err != nil {
				return err
			}

			//检查路由表
			command := fmt.Sprintf("display ip routing-table vpn-instance %v %v\nexit\n", d.Operators, d.IpAddr)
			if err := mtools.CommandGenerator("H3CDisCommand.txt", command); err != nil {
				return err
			}
			//然后检查配置
			//然后在设备上ping测
		} else if device.Brand == common.Huawei_CE {

		}
	}
	return nil
}
