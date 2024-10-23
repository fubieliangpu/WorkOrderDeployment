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

// 冲突检测
func (i *NetProdDeplImpl) ConflictCheck(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (internet.ConfigConflictStatus, error) {
	//校验请求完整性和合法性
	if err := in.Validate(); err != nil {
		return internet.CONFLICT, exception.ErrValidateFailed(err.Error())
	}
	//ping测即将要配置的IP是否通，看是否有IP冲突
	//首先要查询到指定IDC、属于指定设备层级的设备并登录设备
	req := rcdevice.NewQueryDeviceListRequest()
	req.IDC = in.Idc
	req.DeviceLevel = &in.AccessDeviceLevel
	deviceset, err := i.ctldevice.QueryDeviceList(ctx, req)
	if err != nil {
		return internet.CONFLICT, exception.ErrServerInternal(err.Error())
	}
	if len(deviceset.Items) == 0 {
		return internet.CONFLICT, internet.ErrNoDeviceInIdc
	}
	//当请求部署的位置在核心层，需要在对应的vpn-Instance下检查路由条目及在其下ping测，作为第一阶段检验
	//根据不同的接入方式及准备接入的设备层进行冲突检测，需要再添加一个辅助函数判定设备接入层，不同接入层有不同的部署逻辑，所以检测机制不同
	switch in.ConnectMethod {
	case internet.VRRP:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.SINGLE:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	//暂不支持跨机房
	case internet.STATIC_MASTER_BACKUP:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.SHAREGATEWAY:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	case internet.STATIC_LOADBALANCE:
		err := in.BasicCheck(deviceset.Items[0])
		if err != nil {
			return internet.CONFLICT, err
		}
		err = in.BasicCheck(deviceset.Items[1])
		if err != nil {
			return internet.CONFLICT, err
		}
		return internet.FIT, nil

	default:
		return internet.CONFLICT, exception.ErrDeviceAccessMothed("The ConnectMethod is not supported")
	}

}

// 业务配置下发与业务配置回收
func (i *NetProdDeplImpl) ConfigDeployment(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (*internet.NetProd, error) {
	//校验与查询设备
	if err := in.Validate(); err != nil {
		return nil, exception.ErrValidateFailed(err.Error())
	}
	req := rcdevice.NewQueryDeviceListRequest()
	req.IDC = in.Idc
	req.DeviceLevel = &in.AccessDeviceLevel
	deviceset, err := i.ctldevice.QueryDeviceList(ctx, req)
	if err != nil {
		return nil, exception.ErrServerInternal(err.Error())
	}
	if len(deviceset.Items) == 0 {
		return nil, internet.ErrNoDeviceInIdc
	}
	//由于先前已经判断过冲突问题，此处其他判断如分配哪个端口，是否配置携带vpn-instance由用户脚本中体现，同时会存在临时变化，程序不做判断，
	//安全起见,配置前需要人工审核脚本,脚本名称固定为configdeploymentA.txt configdeploymentA.txt,对应主备设备，如果没有备设备则只有configdeploymentA.txt
	configreq := rcdevice.NewChangeDeviceConfigRequest(deviceset.Items[0].Name)
	configreq.UserFile = "user.yaml"
	configreq.DeploymentRecord = "configrecord.txt"
	res := internet.NewNetProd()
	res.DeploymentNetworkProductRequest = in
	if in.ConnectMethod == internet.STATIC_LOADBALANCE || in.ConnectMethod == internet.STATIC_MASTER_BACKUP || in.ConnectMethod == internet.VRRP {
		configreq.DeviceConfigFile = "configdeploymentA.txt"
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
		configreq.DeviceConfigFile = "configdeploymentB.txt"
		configreq.DeviceName = deviceset.Items[1].Name
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
	} else if in.ConnectMethod == internet.SHAREGATEWAY || in.ConnectMethod == internet.SINGLE {
		configreq.DeviceConfigFile = "configdeploymentA.txt"
		if _, err := i.ctldevice.ChangeDeviceConfig(ctx, configreq); err != nil {
			res.Status = internet.FAIL
			return res, err
		}
	} else {
		res.Status = internet.FAIL
		return res, exception.ErrDeviceAccessMothed("The ConnectMethod is not supported")
	}
	res.Status = internet.SUCCESS
	return res, nil
}

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
