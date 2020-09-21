package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const FirmwareCtrlrRunningClassName = "firmwareCtrlrRunning"

type FirmwareCtrlrRunning struct {
	BaseAttributes
	FirmwareCtrlrRunningAttributes
}

type FirmwareCtrlrRunningAttributes struct {
	Mode     string `json:",omitempty"`
	TpmInUse string `json:",omitempty"`
	Ts       string `json:",omitempty"`
	Type     string `json:",omitempty"`
	Version  string `json:",omitempty"`
}

func NewFirmwareCtrlrRunning(firmwareCtrlrRunningRn, parentDn, description string, firmwareCtrlrRunningAttr FirmwareCtrlrRunningAttributes) *FirmwareCtrlrRunning {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareCtrlrRunningRn)
	return &FirmwareCtrlrRunning{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         FirmwareCtrlrRunningClassName,
			Rn:                firmwareCtrlrRunningRn,
		},
		FirmwareCtrlrRunningAttributes: firmwareCtrlrRunningAttr,
	}
}

func (firmwareCtrlrRunning *FirmwareCtrlrRunning) ToMap() (map[string]string, error) {
	firmwareCtrlrRunningMap, err := firmwareCtrlrRunning.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareCtrlrRunningMap, "mode", firmwareCtrlrRunning.Mode)
	A(firmwareCtrlrRunningMap, "tpmInUse", firmwareCtrlrRunning.TpmInUse)
	A(firmwareCtrlrRunningMap, "ts", firmwareCtrlrRunning.Ts)
	A(firmwareCtrlrRunningMap, "type", firmwareCtrlrRunning.Type)
	A(firmwareCtrlrRunningMap, "version", firmwareCtrlrRunning.Version)

	return firmwareCtrlrRunningMap, err
}

func FirmwareCtrlrRunningFromContainerList(cont *container.Container, index int) *FirmwareCtrlrRunning {

	FirmwareCtrlrRunningCont := cont.S("imdata").Index(index).S(FirmwareCtrlrRunningClassName, "attributes")
	return &FirmwareCtrlrRunning{
		BaseAttributes{
			DistinguishedName: G(FirmwareCtrlrRunningCont, "dn"),
			Description:       G(FirmwareCtrlrRunningCont, "descr"),
			Status:            G(FirmwareCtrlrRunningCont, "status"),
			ClassName:         FirmwareCtrlrRunningClassName,
			Rn:                G(FirmwareCtrlrRunningCont, "rn"),
		},

		FirmwareCtrlrRunningAttributes{
			Mode:     G(FirmwareCtrlrRunningCont, "mode"),
			TpmInUse: G(FirmwareCtrlrRunningCont, "tpmInUse"),
			Ts:       G(FirmwareCtrlrRunningCont, "ts"),
			Type:     G(FirmwareCtrlrRunningCont, "type"),
			Version:  G(FirmwareCtrlrRunningCont, "version"),
		},
	}
}

func FirmwareCtrlrRunningFromContainer(cont *container.Container) *FirmwareCtrlrRunning {

	return FirmwareCtrlrRunningFromContainerList(cont, 0)
}

func FirmwareCtrlrRunningListFromContainer(cont *container.Container) []*FirmwareCtrlrRunning {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FirmwareCtrlrRunning, length)

	for i := 0; i < length; i++ {

		arr[i] = FirmwareCtrlrRunningFromContainerList(cont, i)
	}

	return arr
}
