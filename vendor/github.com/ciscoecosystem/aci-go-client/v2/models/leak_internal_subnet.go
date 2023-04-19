package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnleakInternalSubnet        = "uni/tn-%s/ctx-%s/leakroutes/leakintsubnet-[%s]"
	RnleakInternalSubnet        = "leakintsubnet-[%s]"
	ParentDnleakInternalSubnet  = "uni/tn-%s/ctx-%s/leakroutes"
	LeakinternalsubnetClassName = "leakInternalSubnet"
)

type LeakInternalSubnet struct {
	BaseAttributes
	NameAliasAttribute
	LeakInternalSubnetAttributes
}

type LeakInternalSubnetAttributes struct {
	Annotation string `json:",omitempty"`
	Ip         string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Scope      string `json:",omitempty"`
}

func NewLeakInternalSubnet(leakInternalSubnetRn, parentDn, description, nameAlias string, leakInternalSubnetAttr LeakInternalSubnetAttributes) *LeakInternalSubnet {
	dn := fmt.Sprintf("%s/%s", parentDn, leakInternalSubnetRn)
	return &LeakInternalSubnet{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LeakinternalsubnetClassName,
			Rn:                leakInternalSubnetRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LeakInternalSubnetAttributes: leakInternalSubnetAttr,
	}
}

func (leakInternalSubnet *LeakInternalSubnet) ToMap() (map[string]string, error) {
	leakInternalSubnetMap, err := leakInternalSubnet.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := leakInternalSubnet.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(leakInternalSubnetMap, key, value)
	}

	A(leakInternalSubnetMap, "annotation", leakInternalSubnet.Annotation)
	A(leakInternalSubnetMap, "ip", leakInternalSubnet.Ip)
	A(leakInternalSubnetMap, "name", leakInternalSubnet.Name)
	A(leakInternalSubnetMap, "scope", leakInternalSubnet.Scope)
	return leakInternalSubnetMap, err
}

func LeakInternalSubnetFromContainerList(cont *container.Container, index int) *LeakInternalSubnet {
	LeakInternalSubnetCont := cont.S("imdata").Index(index).S(LeakinternalsubnetClassName, "attributes")
	return &LeakInternalSubnet{
		BaseAttributes{
			DistinguishedName: G(LeakInternalSubnetCont, "dn"),
			Description:       G(LeakInternalSubnetCont, "descr"),
			Status:            G(LeakInternalSubnetCont, "status"),
			ClassName:         LeakinternalsubnetClassName,
			Rn:                G(LeakInternalSubnetCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LeakInternalSubnetCont, "nameAlias"),
		},
		LeakInternalSubnetAttributes{
			Annotation: G(LeakInternalSubnetCont, "annotation"),
			Ip:         G(LeakInternalSubnetCont, "ip"),
			Name:       G(LeakInternalSubnetCont, "name"),
			Scope:      G(LeakInternalSubnetCont, "scope"),
		},
	}
}

func LeakInternalSubnetFromContainer(cont *container.Container) *LeakInternalSubnet {
	return LeakInternalSubnetFromContainerList(cont, 0)
}

func LeakInternalSubnetListFromContainer(cont *container.Container) []*LeakInternalSubnet {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LeakInternalSubnet, length)

	for i := 0; i < length; i++ {
		arr[i] = LeakInternalSubnetFromContainerList(cont, i)
	}

	return arr
}
