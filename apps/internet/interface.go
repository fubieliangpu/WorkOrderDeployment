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

//	type Service interface {
//		//冲突检测
//		ConflictCheck(context.Context, *DeploymentNetworkProductRequest) (ConfigConflictStatus, error)
//		//业务配置下发与业务配置回收
//		ConfigDeployment(context.Context, *DeploymentNetworkProductRequest) (*NetProd, error)
//	}
type Service interface {
	//冲突检测
	VrrpConflictCheck(context.Context, *DeploymentVRRP) (ConfigConflictStatus, error)
	DoubleStaticConflictCheck(context.Context, *DeploymentDoubleStatic) (ConfigConflictStatus, error)
	SingleConflictCheck(context.Context, *DeploymentSingle) (ConfigConflictStatus, error)
	//配置下发与业务配置回收
	VrrpDeployment(context.Context, *DeploymentVRRP) (DeploymentResult, error)
	DoubleStaticDeployment(context.Context, *DeploymentDoubleStatic) (DeploymentResult, error)
	SingleDeployment(context.Context, *DeploymentSingle) (DeploymentResult, error)
}

// 检查基础冲突，不同接入层下的指定品牌设备的路由表检查，汇聚、接入层设备ping测检查
// 目前只支持H3C及Huawei_CE设备品牌
func (d *DeploymentNetworkProductRequest) BasicCheck(device *rcdevice.Device) error {
	var err error
	cfi := rcdevice.NewConfigInfo()
	cfi.UserInfo = rcdevice.NewDeviceUserInfo()
	cfi.UserInfo, err = rcdevice.LoadUsernmPasswdFromYaml("user.yaml", cfi.UserInfo)
	if err != nil {
		return err
	}
	if d.AccessDeviceLevel == common.CORE {
		//判断查找到的核心设备品牌

		if device.Brand == common.H3C || device.Brand == common.Huawei_CE {
			cfi.Ip = device.ServerAddr
			cfi.Port = device.Port
			cfi.Configfile = "H3CHWDisCommand.txt"
			cfi.Recordfile = "H3CHWCoreVpn.txt"
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
				err = mtools.Regexper(cfi.Recordfile, 0, fmt.Sprintf("%v/%v", d.IpAddr, d.IpMask), fmt.Sprintf("%v/%v", d.NeighborIp, d.NeighborMask))
				if err == nil {
					return ErrRouteConflict
				}
				//不存在对应运营商的VPN-INSTANCE是如何判断路由表
			} else if err.(*exception.ApiException).Code == 50444 {
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
				err = mtools.Regexper(cfi.Recordfile, 0, fmt.Sprintf("%v/%v", d.IpAddr, d.IpMask), fmt.Sprintf("%v/%v", d.NeighborIp, d.NeighborMask))
				if err == nil {
					return ErrRouteConflict
				}
			}

			//然后检查配置,主要针对是否有之前遗留的静态路由
			//有可能之前互联接口配置没有清空，但是路由表又查不到，这部分垃圾配置也需要找出来，如果有垃圾配置则存在冲突，需要人为去清理
			gips := mtools.GetGatewayByIpStr(d.IpAddr, d.NeighborIp)
			command := fmt.Sprintf(
				"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\nexit\n",
				d.IpAddr,
				gips[0],
				gips[1],
			)
			mtools.CommandGenerator(cfi.Configfile, command)
			cfi.Recordfile = "H3CHWRouteConfig.txt"
			rcdevice.SshConfigTool(cfi)
			err = mtools.Regexper(
				cfi.Recordfile,
				0,
				fmt.Sprintf("ip route-static.*%v %v.*", d.IpAddr, d.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
			)
			//只要匹配到了，则存在冲突，不再继续判断
			if err == nil {
				return ErrRouteConflict
			}
			//其他品牌后面开发
		} else {
			return ErrBrandNotSupport
		}
		//如果客户在汇聚层或接入层设备接入，由于汇聚层设备不存在VPN-INSTANCE，无需判断VPN-INSTANCE
	} else if d.AccessDeviceLevel == common.CONVERGE || d.AccessDeviceLevel == common.ACCESS {
		if device.Brand == common.H3C || device.Brand == common.Huawei_CE {
			cfi.Ip = device.ServerAddr
			cfi.Port = device.Port
			cfi.Configfile = "H3CHWDisCommand.txt"
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
			err := mtools.Regexper(cfi.Recordfile, 0, fmt.Sprintf("%v/%v", d.IpAddr, d.IpMask), fmt.Sprintf("%v/%v", d.NeighborIp, d.NeighborMask))
			if err == nil {
				return ErrRouteConflict
			}

			//然后检查配置,主要针对是否有之前遗留的静态路由
			//有可能之前互联接口配置没有清空，但是路由表又查不到，这部分垃圾配置也需要找出来，如果有垃圾配置则存在冲突，需要人为去清理
			gips := mtools.GetGatewayByIpStr(d.IpAddr, d.NeighborIp)
			command = fmt.Sprintf(
				"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\nexit\n",
				d.IpAddr,
				gips[0],
				gips[1],
			)
			mtools.CommandGenerator(cfi.Configfile, command)
			cfi.Recordfile = "H3CHWRouteConfig.txt"
			rcdevice.SshConfigTool(cfi)
			err = mtools.Regexper(
				cfi.Recordfile,
				0,
				fmt.Sprintf("ip route-static.*%v %v.*", d.IpAddr, d.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
			)
			//只要匹配到了，则存在冲突，不再继续判断
			if err == nil {
				return ErrRouteConflict
			}

			//因为是汇聚层或接入层，看不到全部路由，需要ping测，
			//这时候判断两种情况:
			//1.公网地址互联，本身网关配置在我司交换机上，则IP地址最后一位+1就是网关
			//2.私网地址互联，本身网关配置在我司交换机上，但是使用的时私网地址下一跳，那么客户侧公网地址未必可ping，需要额外测试ping一个互联地址，互联地址+1就是我司交换机端口地址
			//因此需要测试两个地址，公网地址段最后一位+1、私网地址段最后一位+1
			//分割业务IP，切片最后一个元素转为数字后+1再转回来，在组成字符串
			command = fmt.Sprintf(
				"ping %v\nping %v\nexit\n",
				gips[0],
				gips[1],
			)
			mtools.CommandGenerator(cfi.Configfile, command)
			cfi.Recordfile = "H3CHWPing.txt"
			rcdevice.SshConfigTool(cfi)
			err = mtools.Regexper(
				cfi.Recordfile,
				0,
				`, 0\.0% packet loss`,
			)
			//匹配到，则存在冲突
			if err == nil {
				return ErrRouteConflict
			}
		} else {
			return ErrBrandNotSupport
		}
	} else if d.AccessDeviceLevel != common.ACCESS && d.AccessDeviceLevel != common.CONVERGE && d.AccessDeviceLevel != common.CORE {
		return ErrAccessLevelNotFound
	}
	//未检测到基础路由冲突
	return nil
}

//
