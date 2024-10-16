package impl

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/dci"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/mtools"
)

func (i *DCIDeplImpl) ConflictCheck(ctx context.Context, in *dci.CreateDCIRequest) error {
	if err := in.Validate(); err != nil {
		return err
	}
	creq := rcdevice.NewChangeDeviceConfigRequest(in.SrcDevName)
	creq.DeviceConfigFile = "DCIConflictCheckConfig.txt"
	creq.DeploymentRecord = "DCICheckRecord.txt"
	creq.UserFile = "user.yaml"

	//生成冲突检查命令
	mtools.CommandGenerator(
		creq.DeviceConfigFile,
		"screen-length disable\n",
		//Vlan check
		fmt.Sprintf("display vlan %v\n", in.Vlan),
		//Vni check
		fmt.Sprintf("display current-configuration | include %v\n", in.Vni),
		//Vsi check
		fmt.Sprintf("display l2vpn vsi name %v\n", in.Vname),
		//Brige Port check
		fmt.Sprintf("display interface Bridge-Aggregation %v\n", in.BridgePort),
		//Tunnel Check
		fmt.Sprintf("display interface brief description | include Tun.*%v\n", in.DestDevName),
		"exit\n",
	)

	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, creq); err != nil {
		return err
	}

	//根据命令回显结果，正则匹配回显检查
	err := in.BasicCheck(creq.DeploymentRecord)
	return err

}

func (i *DCIDeplImpl) ConfigDeployment(ctx context.Context, in *dci.CreateDCIRequest) (*dci.DCI, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	//根据设备名查找本端设备，并部署配置
	creq := rcdevice.NewChangeDeviceConfigRequest(in.SrcDevName)
	creq.DeviceConfigFile = "DCISrcConfig.txt"
	creq.DeploymentRecord = "DCISrcDeploymentRecord.txt"
	creq.UserFile = "user.yaml"
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, creq); err != nil {
		return nil, err
	}

	//根据设备名查找远端设备，并部署配置
	creq.DeviceName = in.DestDevName
	creq.DeviceConfigFile = "DCIDstConfig.txt"
	creq.DeploymentRecord = "DCIDstDeploymentRecord.txt"
	if _, err := i.ctldevice.ChangeDeviceConfig(ctx, creq); err != nil {
		return nil, err
	}

	//返回dci.DCI类型方便后面扩展
	ddci := dci.NewDCI()
	ddci.CreateDCIRequest = in
	return ddci, nil
}
