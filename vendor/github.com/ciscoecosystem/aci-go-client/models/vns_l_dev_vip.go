package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvnsLDevVip        = "uni/tn-%s/lDevVip-%s"
	RnvnsLDevVip        = "lDevVip-%s"
	ParentDnvnsLDevVip  = "uni/tn-%s"
	VnsldevvipClassName = "vnsLDevVip"
)

type L4ToL7Devices struct {
	BaseAttributes
	NameAliasAttribute
	L4ToL7DevicesAttributes
}

type L4ToL7DevicesAttributes struct {
	Active       string `json:",omitempty"`
	Annotation   string `json:",omitempty"`
	ContextAware string `json:",omitempty"`
	Devtype      string `json:",omitempty"`
	FuncType     string `json:",omitempty"`
	IsCopy       string `json:",omitempty"`
	Managed      string `json:",omitempty"`
	Mode         string `json:",omitempty"`
	Name         string `json:",omitempty"`
	PackageModel string `json:",omitempty"`
	PromMode     string `json:",omitempty"`
	SvcType      string `json:",omitempty"`
	Trunking     string `json:",omitempty"`
}

func NewL4ToL7Devices(vnsLDevVipRn, parentDn, nameAlias string, vnsLDevVipAttr L4ToL7DevicesAttributes) *L4ToL7Devices {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsLDevVipRn)
	return &L4ToL7Devices{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VnsldevvipClassName,
			Rn:                vnsLDevVipRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		L4ToL7DevicesAttributes: vnsLDevVipAttr,
	}
}

func (vnsLDevVip *L4ToL7Devices) ToMap() (map[string]string, error) {
	vnsLDevVipMap, err := vnsLDevVip.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsLDevVip.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsLDevVipMap, key, value)
	}

	A(vnsLDevVipMap, "activeActive", vnsLDevVip.Active)
	A(vnsLDevVipMap, "annotation", vnsLDevVip.Annotation)
	A(vnsLDevVipMap, "contextAware", vnsLDevVip.ContextAware)
	A(vnsLDevVipMap, "devtype", vnsLDevVip.Devtype)
	A(vnsLDevVipMap, "funcType", vnsLDevVip.FuncType)
	A(vnsLDevVipMap, "isCopy", vnsLDevVip.IsCopy)
	A(vnsLDevVipMap, "managed", vnsLDevVip.Managed)
	A(vnsLDevVipMap, "mode", vnsLDevVip.Mode)
	A(vnsLDevVipMap, "name", vnsLDevVip.Name)
	A(vnsLDevVipMap, "packageModel", vnsLDevVip.PackageModel)
	A(vnsLDevVipMap, "promMode", vnsLDevVip.PromMode)
	A(vnsLDevVipMap, "svcType", vnsLDevVip.SvcType)
	A(vnsLDevVipMap, "trunking", vnsLDevVip.Trunking)
	return vnsLDevVipMap, err
}

func L4ToL7DevicesFromContainerList(cont *container.Container, index int) *L4ToL7Devices {
	L4ToL7DevicesCont := cont.S("imdata").Index(index).S(VnsldevvipClassName, "attributes")
	return &L4ToL7Devices{
		BaseAttributes{
			DistinguishedName: G(L4ToL7DevicesCont, "dn"),
			Status:            G(L4ToL7DevicesCont, "status"),
			ClassName:         VnsldevvipClassName,
			Rn:                G(L4ToL7DevicesCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(L4ToL7DevicesCont, "nameAlias"),
		},
		L4ToL7DevicesAttributes{
			Active:       G(L4ToL7DevicesCont, "activeActive"),
			Annotation:   G(L4ToL7DevicesCont, "annotation"),
			ContextAware: G(L4ToL7DevicesCont, "contextAware"),
			Devtype:      G(L4ToL7DevicesCont, "devtype"),
			FuncType:     G(L4ToL7DevicesCont, "funcType"),
			IsCopy:       G(L4ToL7DevicesCont, "isCopy"),
			Managed:      G(L4ToL7DevicesCont, "managed"),
			Mode:         G(L4ToL7DevicesCont, "mode"),
			Name:         G(L4ToL7DevicesCont, "name"),
			PackageModel: G(L4ToL7DevicesCont, "packageModel"),
			PromMode:     G(L4ToL7DevicesCont, "promMode"),
			SvcType:      G(L4ToL7DevicesCont, "svcType"),
			Trunking:     G(L4ToL7DevicesCont, "trunking"),
		},
	}
}

func L4ToL7DevicesFromContainer(cont *container.Container) *L4ToL7Devices {
	return L4ToL7DevicesFromContainerList(cont, 0)
}

func L4ToL7DevicesListFromContainer(cont *container.Container) []*L4ToL7Devices {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*L4ToL7Devices, length)

	for i := 0; i < length; i++ {
		arr[i] = L4ToL7DevicesFromContainerList(cont, i)
	}

	return arr
}
