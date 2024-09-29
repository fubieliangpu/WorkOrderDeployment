package impl

import (
	"context"
	"io"
	"os"

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
		return 1, exception.ErrValidateFailed(err.Error())
	}
	//ping测即将要配置的IP是否通，看是否有IP冲突
	//首先要查询到指定IDC、属于指定设备层级的设备并登录设备
	req := rcdevice.NewQueryDeviceListRequest()
	req.IDC = in.Idc
	req.DeviceLevel = &in.AccessDeviceLevel
	deviceset, err := i.ctldevice.QueryDeviceList(ctx, req)
	if err != nil {
		return 1, exception.ErrServerInternal(err.Error())
	}
	if deviceset.Total == 0 {
		return 1, internet.ErrNoDeviceInIdc
	}
	//当请求部署的位置在核心层，需要在对应的vpn-Instance下检查路由条目及在其下ping测，作为第一阶段检验
	//根据不同的接入方式及准备接入的设备层进行冲突检测，需要再添加一个辅助函数判定设备接入层，不同接入层有不同的部署逻辑，所以检测机制不同
	switch in.ConnectMethod {
	case internet.VRRP:
		if in.AccessDeviceLevel == common.CORE {
			//判断查找到的第一台核心的设备品牌
			if deviceset.Items[0].Brand == common.H3C {
				//首先查看是否存在对应运营商的VPN-INSTANCE，如果存在则查路由表时需要带上对应参数
				discommand := "display ip vpn-instance\n"
				disrecordfile, _ := os.OpenFile("displaycommandfile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
				defer disrecordfile.Close()
				io.WriteString(disrecordfile, discommand)
				mtools.Regexper("displayrecordfile.txt", "")
				//检查路由表

				//然后检查配置
				//然后在设备上ping测
			} else if deviceset.Items[0].Brand == common.Huawei_CE {

			}
		}
	case internet.SINGLE:
	case internet.STATIC_MASTER_BACKUP:
	case internet.SHAREGATEWAY:
	case internet.STATIC_LOADBALANCE:
	default:
	}

	return 1, nil
}

// 业务配置下发
func (i *NetProdDeplImpl) ConfigDeployment(ctx context.Context, in *internet.DeploymentNetworkProductRequest) (*internet.NetProd, error) {
	return nil, nil
}

// 业务配置回收
func (i *NetProdDeplImpl) ConfigRevoke(ctx context.Context, in *internet.UndoDeviceConfigRequest) (*internet.NetProd, error) {
	return nil, nil
}
