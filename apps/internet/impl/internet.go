package impl

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/internet"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/fubieliangpu/WorkOrderDeployment/mtools"
)

func (i *NetProdDeplImpl) VrrpConflictCheck(ctx context.Context, in *internet.DeploymentVRRP) (internet.ConfigConflictStatus, error) {
	if err := in.Validate(); err != nil {
		return internet.CONFLICT, exception.ErrValidateFailed(err.Error())
	}
	//主备设备查询
	for _, v := range []string{in.MasterDevName, in.BackupDevName} {
		querydev := rcdevice.NewDescribeDeviceRequest(v)
		dev, err := i.ctldevice.DescribeDevice(ctx, querydev)
		if err != nil {
			return internet.CONFLICT, err
		}
		if dev.Brand < common.H3C || dev.Brand > common.Huawei_CE {
			return internet.CONFLICT, internet.ErrBrandNotSupport
		}
		//生成设备登录配置
		cfi := rcdevice.NewConfigInfo()
		cfi.Configfile = "VRRPCheckConfig.cnf"
		cfi.Recordfile = "CheckRecord.log"
		cfi.UserInfo = rcdevice.NewDeviceUserInfo()
		if _, err := rcdevice.LoadUsernmPasswdFromYaml("user.yaml", cfi.UserInfo); err != nil {
			return internet.CONFLICT, err
		}
		cfi.Ip = dev.ServerAddr
		cfi.Port = dev.Port

		//生成互联网关地址
		gips := mtools.GetGatewayByIpStr(in.Detail.IpAddr, in.Detail.NeighborIp)

		//分设备接入层级判断
		if dev.DeviceLevel != common.CORE {
			mtools.CommandGenerator(
				cfi.Configfile,
				"screen-length disable\n",
				//Vrid check
				fmt.Sprintf("display current-configuration | include vrrp.vrid.%v .\n", in.Vrid),
				//routing-table check
				fmt.Sprintf(
					"display ip routing-table %v %v\ndisplay ip routing-table %v %v\n",
					in.Detail.IpAddr,
					in.Detail.IpMask,
					in.Detail.NeighborIp,
					in.Detail.NeighborMask,
				),
				//config check
				fmt.Sprintf(
					"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
					in.Detail.IpAddr,
					gips[0],
					gips[1],
				),
				//因为是汇聚层或接入层，看不到全部路由，需要ping测，
				//这时候判断两种情况:
				//1.公网地址互联，本身网关配置在我司交换机上，则IP地址最后一位+1就是网关
				//2.私网地址互联，本身网关配置在我司交换机上，但是使用的时私网地址下一跳，那么客户侧公网地址未必可ping，需要额外测试ping一个互联地址，互联地址+1就是我司交换机端口地址
				//因此需要测试两个地址，公网地址段最后一位+1、私网地址段最后一位+1
				//分割业务IP，切片最后一个元素转为数字后+1再转回来，在组成字符串
				fmt.Sprintf(
					"ping %v\nping %v\n",
					gips[0],
					gips[1],
				),
				//exit
				"exit\n",
			)
			//登录设备查看
			rcdevice.SshConfigTool(cfi)
			err := mtools.Regexper(
				cfi.Recordfile,
				0,
				//vrrp匹配
				fmt.Sprintf("vrrp vrid %v virtual-ip .*", in.Vrid),
				//路由匹配
				fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
				//设备配置匹配
				fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
				//ping测结果判断
				`, 0\.0% packet loss`,
			)
			//匹配到任意一个，就判定为冲突
			if err == nil {
				return internet.CONFLICT, internet.ErrRouteConflict
			}

		} else if dev.DeviceLevel == common.CORE {
			//判断vpn-instance
			mtools.CommandGenerator(
				cfi.Configfile,
				"display ip vpn-instance\nexit\n",
			)
			rcdevice.SshConfigTool(cfi)
			//有vpn-instance如何判断
			if err := mtools.Regexper(cfi.Recordfile, 0, in.Detail.Operators); err == nil {
				mtools.CommandGenerator(
					cfi.Configfile,
					"screen-length disable\n",
					fmt.Sprintf("display current-configuration | include vrrp.vrid.%v.*\n", in.Vrid),
					fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.NeighborIp, in.Detail.NeighborMask),
					fmt.Sprintf(
						"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
						in.Detail.IpAddr,
						gips[0],
						gips[1],
					),
					"exit\n",
				)
				rcdevice.SshConfigTool(cfi)
				if err := mtools.Regexper(
					cfi.Recordfile,
					0,
					//vrrp匹配
					fmt.Sprintf("vrrp vrid %v virtual-ip .*", in.Vrid),
					//路由匹配
					fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
					//设备配置匹配
					fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("ip address %v.*", gips[0]),
					fmt.Sprintf("ip address %v.*", gips[1]),
				); err == nil {
					return internet.CONFLICT, internet.ErrRouteConflict
				}
				//没有vpn-instance
			} else if err.(*exception.ApiException).Code == 50444 {
				mtools.CommandGenerator(
					cfi.Configfile,
					"screen-length disable\n",
					fmt.Sprintf("display current-configuration | include vrrp.vrid.%v.*\n", in.Vrid),
					fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.NeighborIp, in.Detail.NeighborMask),
					fmt.Sprintf(
						"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
						in.Detail.IpAddr,
						gips[0],
						gips[1],
					),
					"exit\n",
				)
				rcdevice.SshConfigTool(cfi)
				if err := mtools.Regexper(
					cfi.Recordfile,
					0,
					//vrrp匹配
					fmt.Sprintf("vrrp vrid %v virtual-ip .*", in.Vrid),
					//路由匹配
					fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
					//设备配置匹配
					fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("ip address %v.*", gips[0]),
					fmt.Sprintf("ip address %v.*", gips[1]),
				); err == nil {
					return internet.CONFLICT, internet.ErrRouteConflict
				}
			}

		}
	}
	return internet.FIT, nil
}

func (i *NetProdDeplImpl) DoubleStaticConflictCheck(ctx context.Context, in *internet.DeploymentDoubleStatic) (internet.ConfigConflictStatus, error) {
	if err := in.Validate(); err != nil {
		return internet.CONFLICT, exception.ErrValidateFailed(err.Error())
	}
	for _, v := range []string{in.FirstDevName, in.SecondDevName} {
		querydev := rcdevice.NewDescribeDeviceRequest(v)
		dev, err := i.ctldevice.DescribeDevice(ctx, querydev)
		if err != nil {
			return internet.CONFLICT, err
		}
		if dev.Brand < common.H3C || dev.Brand > common.Huawei_CE {
			return internet.CONFLICT, internet.ErrBrandNotSupport
		}
		//生成设备登录配置
		cfi := rcdevice.NewConfigInfo()
		cfi.Configfile = "DoubleCheckConfig.cnf"
		cfi.Recordfile = "CheckRecord.log"
		cfi.UserInfo = rcdevice.NewDeviceUserInfo()
		if _, err := rcdevice.LoadUsernmPasswdFromYaml("user.yaml", cfi.UserInfo); err != nil {
			return internet.CONFLICT, err
		}
		cfi.Ip = dev.ServerAddr
		cfi.Port = dev.Port

		//生成互联网关地址
		gips := mtools.GetGatewayByIpStr(in.Detail.IpAddr, in.Detail.NeighborIp)
		//分设备接入层级判断
		if dev.DeviceLevel != common.CORE {
			mtools.CommandGenerator(
				cfi.Configfile,
				"screen-length disable\n",
				//routing-table check
				fmt.Sprintf(
					"display ip routing-table %v %v\ndisplay ip routing-table %v %v\n",
					in.Detail.IpAddr,
					in.Detail.IpMask,
					in.Detail.NeighborIp,
					in.Detail.NeighborMask,
				),
				//config check
				fmt.Sprintf(
					"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
					in.Detail.IpAddr,
					gips[0],
					gips[1],
				),
				//因为是汇聚层或接入层，看不到全部路由，需要ping测，
				//这时候判断两种情况:
				//1.公网地址互联，本身网关配置在我司交换机上，则IP地址最后一位+1就是网关
				//2.私网地址互联，本身网关配置在我司交换机上，但是使用的时私网地址下一跳，那么客户侧公网地址未必可ping，需要额外测试ping一个互联地址，互联地址+1就是我司交换机端口地址
				//因此需要测试两个地址，公网地址段最后一位+1、私网地址段最后一位+1
				//分割业务IP，切片最后一个元素转为数字后+1再转回来，在组成字符串
				fmt.Sprintf(
					"ping %v\nping %v\n",
					gips[0],
					gips[1],
				),
				//exit
				"exit\n",
			)
			//登录设备查看
			rcdevice.SshConfigTool(cfi)
			err := mtools.Regexper(
				cfi.Recordfile,
				0,
				//路由匹配
				fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
				//设备配置匹配
				fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
				//ping测结果判断
				`, 0\.0% packet loss`,
			)
			//匹配到任意一个，就判定为冲突
			if err == nil {
				return internet.CONFLICT, internet.ErrRouteConflict
			}

		} else if dev.DeviceLevel == common.CORE {
			//判断vpn-instance
			mtools.CommandGenerator(
				cfi.Configfile,
				"display ip vpn-instance\nexit\n",
			)
			rcdevice.SshConfigTool(cfi)
			//有vpn-instance如何判断
			if err := mtools.Regexper(cfi.Recordfile, 0, in.Detail.Operators); err == nil {
				mtools.CommandGenerator(
					cfi.Configfile,
					"screen-length disable\n",
					fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.NeighborIp, in.Detail.NeighborMask),
					fmt.Sprintf(
						"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
						in.Detail.IpAddr,
						gips[0],
						gips[1],
					),
					"exit\n",
				)
				rcdevice.SshConfigTool(cfi)
				if err := mtools.Regexper(
					cfi.Recordfile,
					0,
					//路由匹配
					fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
					//设备配置匹配
					fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("ip address %v.*", gips[0]),
					fmt.Sprintf("ip address %v.*", gips[1]),
				); err == nil {
					return internet.CONFLICT, internet.ErrRouteConflict
				}
				//没有vpn-instance
			} else if err.(*exception.ApiException).Code == 50444 {
				mtools.CommandGenerator(
					cfi.Configfile,
					"screen-length disable\n",
					fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.NeighborIp, in.Detail.NeighborMask),
					fmt.Sprintf(
						"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
						in.Detail.IpAddr,
						gips[0],
						gips[1],
					),
					"exit\n",
				)
				rcdevice.SshConfigTool(cfi)
				if err := mtools.Regexper(
					cfi.Recordfile,
					0,
					//路由匹配
					fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
					//设备配置匹配
					fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
					fmt.Sprintf("ip address %v.*", gips[0]),
					fmt.Sprintf("ip address %v.*", gips[1]),
				); err == nil {
					return internet.CONFLICT, internet.ErrRouteConflict
				}
			}

		}
	}
	return internet.FIT, nil
}

func (i *NetProdDeplImpl) SingleConflictCheck(ctx context.Context, in *internet.DeploymentSingle) (internet.ConfigConflictStatus, error) {
	if err := in.Validate(); err != nil {
		return internet.CONFLICT, exception.ErrValidateFailed(err.Error())
	}
	//查询设备
	querydev := rcdevice.NewDescribeDeviceRequest(in.DevName)
	dev, err := i.ctldevice.DescribeDevice(ctx, querydev)
	if err != nil {
		return internet.CONFLICT, err
	}
	if dev.Brand < common.H3C || dev.Brand > common.Huawei_CE {
		return internet.CONFLICT, internet.ErrBrandNotSupport
	}
	//生成设备登录配置
	cfi := rcdevice.NewConfigInfo()
	cfi.Configfile = "SingleCheckConfig.cnf"
	cfi.Recordfile = "CheckRecord.log"
	cfi.UserInfo = rcdevice.NewDeviceUserInfo()
	if _, err := rcdevice.LoadUsernmPasswdFromYaml("user.yaml", cfi.UserInfo); err != nil {
		return internet.CONFLICT, err
	}
	cfi.Ip = dev.ServerAddr
	cfi.Port = dev.Port

	//生成互联网关地址
	gips := mtools.GetGatewayByIpStr(in.Detail.IpAddr, in.Detail.NeighborIp)
	//分设备接入层级判断
	if dev.DeviceLevel != common.CORE {
		mtools.CommandGenerator(
			cfi.Configfile,
			"screen-length disable\n",
			//routing-table check
			fmt.Sprintf(
				"display ip routing-table %v %v\ndisplay ip routing-table %v %v\n",
				in.Detail.IpAddr,
				in.Detail.IpMask,
				in.Detail.NeighborIp,
				in.Detail.NeighborMask,
			),
			//config check
			fmt.Sprintf(
				"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
				in.Detail.IpAddr,
				gips[0],
				gips[1],
			),
			//因为是汇聚层或接入层，看不到全部路由，需要ping测，
			//这时候判断两种情况:
			//1.公网地址互联，本身网关配置在我司交换机上，则IP地址最后一位+1就是网关
			//2.私网地址互联，本身网关配置在我司交换机上，但是使用的时私网地址下一跳，那么客户侧公网地址未必可ping，需要额外测试ping一个互联地址，互联地址+1就是我司交换机端口地址
			//因此需要测试两个地址，公网地址段最后一位+1、私网地址段最后一位+1
			//分割业务IP，切片最后一个元素转为数字后+1再转回来，在组成字符串
			fmt.Sprintf(
				"ping %v\nping %v\n",
				gips[0],
				gips[1],
			),
			//exit
			"exit\n",
		)
		//登录设备查看
		rcdevice.SshConfigTool(cfi)
		err := mtools.Regexper(
			cfi.Recordfile,
			0,
			//路由匹配
			fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
			fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
			//设备配置匹配
			fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
			fmt.Sprintf("ip address %v.*", gips[0]),
			fmt.Sprintf("ip address %v.*", gips[1]),
			//ping测结果判断
			`, 0\.0% packet loss`,
		)
		//匹配到任意一个，就判定为冲突
		if err == nil {
			return internet.CONFLICT, internet.ErrRouteConflict
		}

	} else if dev.DeviceLevel == common.CORE {
		//判断vpn-instance
		mtools.CommandGenerator(
			cfi.Configfile,
			"display ip vpn-instance\nexit\n",
		)
		rcdevice.SshConfigTool(cfi)
		//有vpn-instance如何判断
		if err := mtools.Regexper(cfi.Recordfile, 0, in.Detail.Operators); err == nil {
			mtools.CommandGenerator(
				cfi.Configfile,
				"screen-length disable\n",
				fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("display ip routing-table vpn-instance %v %v %v\n", in.Detail.Operators, in.Detail.NeighborIp, in.Detail.NeighborMask),
				fmt.Sprintf(
					"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
					in.Detail.IpAddr,
					gips[0],
					gips[1],
				),
				"exit\n",
			)
			rcdevice.SshConfigTool(cfi)
			if err := mtools.Regexper(
				cfi.Recordfile,
				0,
				//路由匹配
				fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
				//设备配置匹配
				fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
			); err == nil {
				return internet.CONFLICT, internet.ErrRouteConflict
			}
			//没有vpn-instance
		} else if err.(*exception.ApiException).Code == 50444 {
			mtools.CommandGenerator(
				cfi.Configfile,
				"screen-length disable\n",
				fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("display ip routing-table %v %v\n", in.Detail.NeighborIp, in.Detail.NeighborMask),
				fmt.Sprintf(
					"display current-configuration | include %v\ndisplay current-configuration | include %v\ndisplay current-configuration | include %v\n",
					in.Detail.IpAddr,
					gips[0],
					gips[1],
				),
				"exit\n",
			)
			rcdevice.SshConfigTool(cfi)
			if err := mtools.Regexper(
				cfi.Recordfile,
				0,
				//路由匹配
				fmt.Sprintf("%v/%v", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("%v/%v", in.Detail.NeighborIp, in.Detail.NeighborMask),
				//设备配置匹配
				fmt.Sprintf("ip route-static.*%v %v.*", in.Detail.IpAddr, in.Detail.IpMask),
				fmt.Sprintf("ip address %v.*", gips[0]),
				fmt.Sprintf("ip address %v.*", gips[1]),
			); err == nil {
				return internet.CONFLICT, internet.ErrRouteConflict
			}
		}

	}
	return internet.FIT, nil
}

func (i *NetProdDeplImpl) VrrpDeployment(ctx context.Context, in *internet.DeploymentVRRP) (internet.DeploymentResult, error) {
	//校验与查询设备
	if err := in.Validate(); err != nil {
		return internet.FAIL, exception.ErrValidateFailed(err.Error())
	}
	//创建配置请求
	configreq := rcdevice.NewChangeDeviceConfigRequest(in.MasterDevName)
	configreq.UserFile = "user.yaml"
	configreq.DeploymentRecord = "DeploymentMasterRecord.log"
	configreq.DeviceConfigFile = "MasterConfig.cnf"
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
		return internet.FAIL, err
	}
	configreq.DeploymentRecord = "DeploymentBackupRecord.log"
	configreq.DeviceConfigFile = "BackupConfig.cnf"
	configreq.DeviceName = in.BackupDevName
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
		return internet.FAIL, err
	}
	return internet.SUCCESS, nil
}

func (i *NetProdDeplImpl) DoubleStaticDeployment(ctx context.Context, in *internet.DeploymentDoubleStatic) (internet.DeploymentResult, error) {
	//校验与查询设备
	if err := in.Validate(); err != nil {
		return internet.FAIL, exception.ErrValidateFailed(err.Error())
	}
	//创建配置请求
	configreq := rcdevice.NewChangeDeviceConfigRequest(in.FirstDevName)
	configreq.UserFile = "user.yaml"
	configreq.DeploymentRecord = "DeploymentFirstRecord.log"
	configreq.DeviceConfigFile = "FirstConfig.cnf"
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
		return internet.FAIL, err
	}
	configreq.DeploymentRecord = "DeploymentSecondRecord.log"
	configreq.DeviceConfigFile = "SecondConfig.cnf"
	configreq.DeviceName = in.SecondDevName
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
		return internet.FAIL, err
	}
	return internet.SUCCESS, nil
}

func (i *NetProdDeplImpl) SingleDeployment(ctx context.Context, in *internet.DeploymentSingle) (internet.DeploymentResult, error) {
	//校验与查询设备
	if err := in.Validate(); err != nil {
		return internet.FAIL, exception.ErrValidateFailed(err.Error())
	}
	//创建配置请求
	configreq := rcdevice.NewChangeDeviceConfigRequest(in.DevName)
	configreq.UserFile = "user.yaml"
	configreq.DeploymentRecord = "DeploymentSingleRecord.log"
	configreq.DeviceConfigFile = "SingleConfig.cnf"
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
		return internet.FAIL, err
	}
	return internet.SUCCESS, nil
}
