package dci

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/mtools"
)

const (
	AppName = "dci"
)

type Service interface {
	//冲突检测
	ConflictCheck(context.Context, *CreateDCIRequest) (ConflictStatus, error)
	//配置下发与配置回收
	ConfigDeployment(context.Context, *CreateDCIRequest) (*DCI, error)
}

// 网内vxlan设备主要为h3c，因此只实现h3c的校验逻辑
func (c *CreateDCIRequest) BasicCheck(cfi *rcdevice.ConfigInfo) error {
	//Vlan check
	vlancheck := fmt.Sprintf("display vlan %v\n", c.Vlan)
	//Vni check
	vnicheck := fmt.Sprintf("display current-configuration | include %v\n", c.Vni)
	//Vsi check
	vsicheck := fmt.Sprintf("display l2vpn vsi name %v\n", c.Vname)
	//Brige Port check
	brdgcheck := fmt.Sprintf("display interface Bridge-Aggregation %v\n", c.BridgePort)
	//Tunnel Check
	tunnelcheck := fmt.Sprintf("display interface brief description | include Tun.*%v\n", c.DestDevName)
	//number of port Check
	portnumcheck := "display interface brief down\n"
	mtools.CommandGenerator(
		cfi.Configfile,
		"screen-length disable\n",
		vlancheck,
		vnicheck,
		vsicheck,
		brdgcheck,
		tunnelcheck,
		portnumcheck,
		"exit\n",
	)
	err := mtools.Regexper(
		cfi.Recordfile,
		0,
		`VLAN ID: \d+`,
		` vxlan \d+`,
		`Total number of VSIs: 1,`,
		`Current state: `,
		`Tun\d+\s+\w+\s+\w+\s+.*T`,
	)

	return nil

}
