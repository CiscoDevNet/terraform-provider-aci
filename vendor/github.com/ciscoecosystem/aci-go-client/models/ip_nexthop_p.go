package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const IpnexthoppClassName = "ipNexthopP"

type L3outStaticRouteNextHop struct {
	BaseAttributes
	L3outStaticRouteNextHopAttributes
}

type L3outStaticRouteNextHopAttributes struct {
	NhAddr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Pref string `json:",omitempty"`

	NexthopProfile_type string `json:",omitempty"`
}

func NewL3outStaticRouteNextHop(ipNexthopPRn, parentDn, description string, ipNexthopPattr L3outStaticRouteNextHopAttributes) *L3outStaticRouteNextHop {
	dn := fmt.Sprintf("%s/%s", parentDn, ipNexthopPRn)
	return &L3outStaticRouteNextHop{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IpnexthoppClassName,
			Rn:                ipNexthopPRn,
		},

		L3outStaticRouteNextHopAttributes: ipNexthopPattr,
	}
}

func (ipNexthopP *L3outStaticRouteNextHop) ToMap() (map[string]string, error) {
	ipNexthopPMap, err := ipNexthopP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ipNexthopPMap, "nhAddr", ipNexthopP.NhAddr)

	A(ipNexthopPMap, "annotation", ipNexthopP.Annotation)

	A(ipNexthopPMap, "nameAlias", ipNexthopP.NameAlias)

	A(ipNexthopPMap, "pref", ipNexthopP.Pref)

	A(ipNexthopPMap, "type", ipNexthopP.NexthopProfile_type)

	return ipNexthopPMap, err
}

func L3outStaticRouteNextHopFromContainerList(cont *container.Container, index int) *L3outStaticRouteNextHop {

	L3outStaticRouteNextHopCont := cont.S("imdata").Index(index).S(IpnexthoppClassName, "attributes")
	return &L3outStaticRouteNextHop{
		BaseAttributes{
			DistinguishedName: G(L3outStaticRouteNextHopCont, "dn"),
			Description:       G(L3outStaticRouteNextHopCont, "descr"),
			Status:            G(L3outStaticRouteNextHopCont, "status"),
			ClassName:         IpnexthoppClassName,
			Rn:                G(L3outStaticRouteNextHopCont, "rn"),
		},

		L3outStaticRouteNextHopAttributes{

			NhAddr: G(L3outStaticRouteNextHopCont, "nhAddr"),

			Annotation: G(L3outStaticRouteNextHopCont, "annotation"),

			NameAlias: G(L3outStaticRouteNextHopCont, "nameAlias"),

			Pref: G(L3outStaticRouteNextHopCont, "pref"),

			NexthopProfile_type: G(L3outStaticRouteNextHopCont, "type"),
		},
	}
}

func L3outStaticRouteNextHopFromContainer(cont *container.Container) *L3outStaticRouteNextHop {

	return L3outStaticRouteNextHopFromContainerList(cont, 0)
}

func L3outStaticRouteNextHopListFromContainer(cont *container.Container) []*L3outStaticRouteNextHop {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outStaticRouteNextHop, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outStaticRouteNextHopFromContainerList(cont, i)
	}

	return arr
}
