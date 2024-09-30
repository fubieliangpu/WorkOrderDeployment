package internet

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
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
		//判断查找到的核心设备品牌
		if brand == common.H3C || brand == common.Huawei_CE {
			cfi.Ip = device.ServerAddr
			cfi.Port = device.Port
			cfi.Configfile = "H3CHWDisCommand.txt"
			cfi.Recordfile = "H3CHWCoreVRRPVpn.txt"
			//首先查看是否存在对应运营商的VPN-INSTANCE，如果存在则查路由表时需要带上对应参数
			mtools.CommandGenerator(cfi.Configfile, "display ip vpn-instance\nexit\n")
			rcdevice.SshConfigTool(cfi)
			err := mtools.Regexper(cfi.Recordfile, 0, d.Operators)

			//检查路由表
			//存在对应运营商的VPN-INSTANCE是如何判断路由表
			if err == nil {
				command := fmt.Sprintf(
					"display ip routing-table vpn-instance %v %v %v\ndisplay ip routing-table vpn-instance %v %v %v\nexit\n",
					d.Operators,
					d.IpAddr,
					d.IpMask,
					d.Operators,
					d.NeighborIp,
					d.NeighborMask,
				)
				mtools.CommandGenerator(cfi.Configfile, command)
				cfi.Recordfile = "H3CHWVPNRoutingtable.txt"
				rcdevice.SshConfigTool(cfi)
				//只要匹配到了要下发的地址，则存在冲突，不再继续判断,采用or的匹配模式，即只要有一则正则匹配到就算匹配成功，返回nil
				err = mtools.Regexper(cfi.Recordfile, 0, fmt.Sprintf("^%v", d.IpAddr), fmt.Sprintf("^%v", d.NeighborIp))
				if err == nil {
					return ErrRouteConflict
				}
			}
			//不存在对应运营商的VPN-INSTANCE是如何判断路由表
			if err == exception.ErrRegularMatchFailed("Regular expression matching failed!") {
				command := fmt.Sprintf(
					"display ip routing-table %v %v\ndisplay ip routing-table %v %v\nexit\n",
					d.IpAddr,
					d.IpMask,
					d.NeighborIp,
					d.NeighborMask,
				)
				mtools.CommandGenerator(cfi.Configfile, command)
				cfi.Recordfile = "H3CHWRoutingtable.txt"
				rcdevice.SshConfigTool(cfi)
				//只要匹配到了要下发的地址，则存在冲突，不再继续判断,采用or的匹配模式，即只要有一则正则匹配到就算匹配成功，返回nil
				err = mtools.Regexper(cfi.Recordfile, 0, fmt.Sprintf("^%v", d.IpAddr), fmt.Sprintf("^%v", d.NeighborIp))
				if err == nil {
					return ErrRouteConflict
				}
			}

			//然后检查配置,主要针对是否有之前遗留的静态路由
			command := fmt.Sprintf(
				"display current-configuration | include %v\nexit\n",
				d.IpAddr,
			)
			mtools.CommandGenerator(cfi.Configfile, command)
			cfi.Recordfile = "H3CHWRouteConfig.txt"
			rcdevice.SshConfigTool(cfi)
			err = mtools.Regexper(
				cfi.Recordfile,
				0,
				fmt.Sprintf("ip route-static.*%v %v.*", d.IpAddr, d.IpMask),
			)
			//只要匹配到了，则存在冲突，不再继续判断
			if err == nil {
				return ErrRouteConflict
			}

		} else {
			return ErrBrandNotSupport
		}
	} else if d.AccessDeviceLevel == common.CONVERGE {

	}
	return nil
}
