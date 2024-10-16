package dci

import (
	"context"

	"github.com/fubieliangpu/WorkOrderDeployment/mtools"
)

const (
	AppName = "dci"
)

type Service interface {
	//冲突检测
	ConflictCheck(context.Context, *CreateDCIRequest) error
	//配置下发与配置回收
	ConfigDeployment(context.Context, *CreateDCIRequest) (*DCI, error)
}

// 网内vxlan设备主要为h3c，因此只实现h3c的校验逻辑
func (c *CreateDCIRequest) BasicCheck(recordfile string) error {

	//vlan冲突时的反馈
	if err := mtools.Regexper(
		recordfile,
		0,
		`VLAN ID: \d+`,
	); err == nil {
		return ErrDCIVlanconflict
	}

	//vni冲突时的反馈
	if err := mtools.Regexper(
		recordfile,
		0,
		` vxlan \d+`,
	); err == nil {
		return ErrDCIVNIconflict
	}

	//vsi名称冲突时的反馈
	if err := mtools.Regexper(
		recordfile,
		0,
		`Total number of VSIs: 1,`,
	); err == nil {
		return ErrDCIVSIconflict
	}

	//聚合端口号冲突时的反馈
	if err := mtools.Regexper(
		recordfile,
		0,
		`Current state: `,
	); err == nil {
		return ErrDCIBRconflict
	}

	//Tunnel口不存在时的反馈,不存在时需要手动创建Tunnel口
	if err := mtools.Regexper(
		recordfile,
		0,
		`Tun\d+\s+\w+\s+\w+\s+.*T`,
	); err != nil {
		return ErrDCINoTun
	} else {
		return nil
	}
}
