package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const FirmwareCtrlrFwPClassName = "firmwareCtrlrFwP"

type CtrlrFwP struct {
	BaseAttributes
	CtrlrFwPAttributes
}

type CtrlrFwPAttributes struct {
	Name         string `json:",omitempty"`
	Annotation   string `json:",omitempty"`
	IgnoreCompat string `json:",omitempty"`
	Version      string `json:",omitempty"`
}

func NewCtrlrFwP(firmwareCtrlrFwPRn, parentDn, description string, firmwareCtrlrFwPAttr CtrlrFwPAttributes) *CtrlrFwP {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareCtrlrFwPRn)
	return &CtrlrFwP{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "modified",
			ClassName:         FirmwareCtrlrFwPClassName,
			Rn:                firmwareCtrlrFwPRn,
		},
		CtrlrFwPAttributes: firmwareCtrlrFwPAttr,
	}
}

func (firmwareCtrlrFwP *CtrlrFwP) ToMap() (map[string]string, error) {
	firmwareCtrlrFwPMap, err := firmwareCtrlrFwP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareCtrlrFwPMap, "name", firmwareCtrlrFwP.Name)
	A(firmwareCtrlrFwPMap, "annotation", firmwareCtrlrFwP.Annotation)
	A(firmwareCtrlrFwPMap, "ignoreCompat", firmwareCtrlrFwP.IgnoreCompat)
	A(firmwareCtrlrFwPMap, "version", firmwareCtrlrFwP.Version)

	return firmwareCtrlrFwPMap, err
}

func CtrlrFwPFromContainerList(cont *container.Container, index int) *CtrlrFwP {

	CtrlrFwPCont := cont.S("imdata").Index(index).S(FirmwareCtrlrFwPClassName, "attributes")
	return &CtrlrFwP{
		BaseAttributes{
			DistinguishedName: G(CtrlrFwPCont, "dn"),
			Description:       G(CtrlrFwPCont, "descr"),
			Status:            G(CtrlrFwPCont, "status"),
			ClassName:         FirmwareCtrlrFwPClassName,
			Rn:                G(CtrlrFwPCont, "rn"),
		},

		CtrlrFwPAttributes{
			Name:         G(CtrlrFwPCont, "name"),
			Annotation:   G(CtrlrFwPCont, "annotation"),
			IgnoreCompat: G(CtrlrFwPCont, "ignoreCompat"),
			Version:      G(CtrlrFwPCont, "version"),
		},
	}
}

func CtrlrFwPFromContainer(cont *container.Container) *CtrlrFwP {

	return CtrlrFwPFromContainerList(cont, 0)
}

func CtrlrFwPListFromContainer(cont *container.Container) []*CtrlrFwP {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CtrlrFwP, length)

	for i := 0; i < length; i++ {

		arr[i] = CtrlrFwPFromContainerList(cont, i)
	}

	return arr
}
