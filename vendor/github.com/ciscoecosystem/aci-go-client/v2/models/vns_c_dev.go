package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnvnsCDev        = "%s/cDev-%s"
	RnvnsCDev        = "cDev-%s"
	VnscdevClassName = "vnsCDev"
)

type ConcreteDevice struct {
	BaseAttributes
	NameAliasAttribute
	ConcreteDeviceAttributes
}

type ConcreteDeviceAttributes struct {
	CloneCount       string `json:",omitempty"`
	DevCtxLbl        string `json:",omitempty"`
	Host             string `json:",omitempty"`
	IsCloneOperation string `json:",omitempty"`
	IsTemplate       string `json:",omitempty"`
	Name             string `json:",omitempty"`
	VcenterName      string `json:",omitempty"`
	VmName           string `json:",omitempty"`
	Annotation       string `json:",omitempty"`
}

func NewConcreteDevice(vnsCDevRn, parentDn, nameAlias string, vnsCDevAttr ConcreteDeviceAttributes) *ConcreteDevice {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsCDevRn)
	return &ConcreteDevice{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VnscdevClassName,
			Rn:                vnsCDevRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ConcreteDeviceAttributes: vnsCDevAttr,
	}
}

func (vnsCDev *ConcreteDevice) ToMap() (map[string]string, error) {
	vnsCDevMap, err := vnsCDev.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsCDev.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsCDevMap, key, value)
	}

	A(vnsCDevMap, "cloneCount", vnsCDev.CloneCount)
	A(vnsCDevMap, "annotation", vnsCDev.Annotation)
	A(vnsCDevMap, "devCtxLbl", vnsCDev.DevCtxLbl)
	A(vnsCDevMap, "host", vnsCDev.Host)
	A(vnsCDevMap, "isCloneOperation", vnsCDev.IsCloneOperation)
	A(vnsCDevMap, "isTemplate", vnsCDev.IsTemplate)
	A(vnsCDevMap, "name", vnsCDev.Name)
	A(vnsCDevMap, "vcenterName", vnsCDev.VcenterName)
	A(vnsCDevMap, "vmName", vnsCDev.VmName)
	return vnsCDevMap, err
}

func ConcreteDeviceFromContainerList(cont *container.Container, index int) *ConcreteDevice {
	ConcreteDeviceCont := cont.S("imdata").Index(index).S(VnscdevClassName, "attributes")
	return &ConcreteDevice{
		BaseAttributes{
			DistinguishedName: G(ConcreteDeviceCont, "dn"),
			Status:            G(ConcreteDeviceCont, "status"),
			ClassName:         VnscdevClassName,
			Rn:                G(ConcreteDeviceCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ConcreteDeviceCont, "nameAlias"),
		},
		ConcreteDeviceAttributes{
			CloneCount:       G(ConcreteDeviceCont, "cloneCount"),
			Annotation:       G(ConcreteDeviceCont, "annotation"),
			DevCtxLbl:        G(ConcreteDeviceCont, "devCtxLbl"),
			Host:             G(ConcreteDeviceCont, "host"),
			IsCloneOperation: G(ConcreteDeviceCont, "isCloneOperation"),
			IsTemplate:       G(ConcreteDeviceCont, "isTemplate"),
			Name:             G(ConcreteDeviceCont, "name"),
			VcenterName:      G(ConcreteDeviceCont, "vcenterName"),
			VmName:           G(ConcreteDeviceCont, "vmName"),
		},
	}
}

func ConcreteDeviceFromContainer(cont *container.Container) *ConcreteDevice {
	return ConcreteDeviceFromContainerList(cont, 0)
}

func ConcreteDeviceListFromContainer(cont *container.Container) []*ConcreteDevice {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ConcreteDevice, length)

	for i := 0; i < length; i++ {
		arr[i] = ConcreteDeviceFromContainerList(cont, i)
	}

	return arr
}
