package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extsubnetClassName = "l3extSubnet"

type L3ExtSubnet struct {
	BaseAttributes
	L3ExtSubnetAttributes
}

type L3ExtSubnetAttributes struct {
	Ip string `json:",omitempty"`

	Aggregate string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Scope string `json:",omitempty"`
}

func NewL3ExtSubnet(l3extSubnetRn, parentDn, description string, l3extSubnetattr L3ExtSubnetAttributes) *L3ExtSubnet {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extSubnetRn)
	return &L3ExtSubnet{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extsubnetClassName,
			Rn:                l3extSubnetRn,
		},

		L3ExtSubnetAttributes: l3extSubnetattr,
	}
}

func (l3extSubnet *L3ExtSubnet) ToMap() (map[string]string, error) {
	l3extSubnetMap, err := l3extSubnet.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extSubnetMap, "ip", l3extSubnet.Ip)

	A(l3extSubnetMap, "aggregate", l3extSubnet.Aggregate)

	A(l3extSubnetMap, "annotation", l3extSubnet.Annotation)

	A(l3extSubnetMap, "nameAlias", l3extSubnet.NameAlias)

	A(l3extSubnetMap, "scope", l3extSubnet.Scope)

	return l3extSubnetMap, err
}

func L3ExtSubnetFromContainerList(cont *container.Container, index int) *L3ExtSubnet {

	SubnetCont := cont.S("imdata").Index(index).S(L3extsubnetClassName, "attributes")
	return &L3ExtSubnet{
		BaseAttributes{
			DistinguishedName: G(SubnetCont, "dn"),
			Description:       G(SubnetCont, "descr"),
			Status:            G(SubnetCont, "status"),
			ClassName:         L3extsubnetClassName,
			Rn:                G(SubnetCont, "rn"),
		},

		L3ExtSubnetAttributes{

			Ip: G(SubnetCont, "ip"),

			Aggregate: G(SubnetCont, "aggregate"),

			Annotation: G(SubnetCont, "annotation"),

			NameAlias: G(SubnetCont, "nameAlias"),

			Scope: G(SubnetCont, "scope"),
		},
	}
}

func L3ExtSubnetFromContainer(cont *container.Container) *L3ExtSubnet {

	return L3ExtSubnetFromContainerList(cont, 0)
}

func L3ExtSubnetListFromContainer(cont *container.Container) []*L3ExtSubnet {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3ExtSubnet, length)

	for i := 0; i < length; i++ {

		arr[i] = L3ExtSubnetFromContainerList(cont, i)
	}

	return arr
}
