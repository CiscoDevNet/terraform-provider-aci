package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extvirtuallifpClassName = "l3extVirtualLIfP"

type VirtualLogicalInterfaceProfile struct {
	BaseAttributes
	VirtualLogicalInterfaceProfileAttributes
}

type VirtualLogicalInterfaceProfileAttributes struct {
	NodeDn string `json:",omitempty"`

	Encap string `json:",omitempty"`

	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Autostate string `json:",omitempty"`

	EncapScope string `json:",omitempty"`

	IfInstT string `json:",omitempty"`

	Ipv6Dad string `json:",omitempty"`

	LlAddr string `json:",omitempty"`

	Mac string `json:",omitempty"`

	Mode string `json:",omitempty"`

	Mtu string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`

	Userdom string `json:",omitempty"`
}

func NewVirtualLogicalInterfaceProfile(l3extVirtualLIfPRn, parentDn, description string, l3extVirtualLIfPattr VirtualLogicalInterfaceProfileAttributes) *VirtualLogicalInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extVirtualLIfPRn)
	return &VirtualLogicalInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extvirtuallifpClassName,
			Rn:                l3extVirtualLIfPRn,
		},

		VirtualLogicalInterfaceProfileAttributes: l3extVirtualLIfPattr,
	}
}

func (l3extVirtualLIfP *VirtualLogicalInterfaceProfile) ToMap() (map[string]string, error) {
	l3extVirtualLIfPMap, err := l3extVirtualLIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extVirtualLIfPMap, "nodeDn", l3extVirtualLIfP.NodeDn)

	A(l3extVirtualLIfPMap, "encap", l3extVirtualLIfP.Encap)

	A(l3extVirtualLIfPMap, "addr", l3extVirtualLIfP.Addr)

	A(l3extVirtualLIfPMap, "annotation", l3extVirtualLIfP.Annotation)

	A(l3extVirtualLIfPMap, "autostate", l3extVirtualLIfP.Autostate)

	A(l3extVirtualLIfPMap, "encapScope", l3extVirtualLIfP.EncapScope)

	A(l3extVirtualLIfPMap, "ifInstT", l3extVirtualLIfP.IfInstT)

	A(l3extVirtualLIfPMap, "ipv6Dad", l3extVirtualLIfP.Ipv6Dad)

	A(l3extVirtualLIfPMap, "llAddr", l3extVirtualLIfP.LlAddr)

	A(l3extVirtualLIfPMap, "mac", l3extVirtualLIfP.Mac)

	A(l3extVirtualLIfPMap, "mode", l3extVirtualLIfP.Mode)

	A(l3extVirtualLIfPMap, "mtu", l3extVirtualLIfP.Mtu)

	A(l3extVirtualLIfPMap, "targetDscp", l3extVirtualLIfP.TargetDscp)

	A(l3extVirtualLIfPMap, "userdom", l3extVirtualLIfP.Userdom)

	return l3extVirtualLIfPMap, err
}

func VirtualLogicalInterfaceProfileFromContainerList(cont *container.Container, index int) *VirtualLogicalInterfaceProfile {

	VirtualLogicalInterfaceProfileCont := cont.S("imdata").Index(index).S(L3extvirtuallifpClassName, "attributes")
	return &VirtualLogicalInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(VirtualLogicalInterfaceProfileCont, "dn"),
			Description:       G(VirtualLogicalInterfaceProfileCont, "descr"),
			Status:            G(VirtualLogicalInterfaceProfileCont, "status"),
			ClassName:         L3extvirtuallifpClassName,
			Rn:                G(VirtualLogicalInterfaceProfileCont, "rn"),
		},

		VirtualLogicalInterfaceProfileAttributes{

			NodeDn: G(VirtualLogicalInterfaceProfileCont, "nodeDn"),

			Encap: G(VirtualLogicalInterfaceProfileCont, "encap"),

			Addr: G(VirtualLogicalInterfaceProfileCont, "addr"),

			Annotation: G(VirtualLogicalInterfaceProfileCont, "annotation"),

			Autostate: G(VirtualLogicalInterfaceProfileCont, "autostate"),

			EncapScope: G(VirtualLogicalInterfaceProfileCont, "encapScope"),

			IfInstT: G(VirtualLogicalInterfaceProfileCont, "ifInstT"),

			Ipv6Dad: G(VirtualLogicalInterfaceProfileCont, "ipv6Dad"),

			LlAddr: G(VirtualLogicalInterfaceProfileCont, "llAddr"),

			Mac: G(VirtualLogicalInterfaceProfileCont, "mac"),

			Mode: G(VirtualLogicalInterfaceProfileCont, "mode"),

			Mtu: G(VirtualLogicalInterfaceProfileCont, "mtu"),

			TargetDscp: G(VirtualLogicalInterfaceProfileCont, "targetDscp"),

			Userdom: G(VirtualLogicalInterfaceProfileCont, "userdom"),
		},
	}
}

func VirtualLogicalInterfaceProfileFromContainer(cont *container.Container) *VirtualLogicalInterfaceProfile {

	return VirtualLogicalInterfaceProfileFromContainerList(cont, 0)
}

func VirtualLogicalInterfaceProfileListFromContainer(cont *container.Container) []*VirtualLogicalInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VirtualLogicalInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = VirtualLogicalInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
