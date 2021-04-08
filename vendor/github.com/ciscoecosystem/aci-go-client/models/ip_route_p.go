package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const IproutepClassName = "ipRouteP"

type L3outStaticRoute struct {
	BaseAttributes
	L3outStaticRouteAttributes
}

type L3outStaticRouteAttributes struct {
	Ip string `json:",omitempty"`

	Aggregate string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Pref string `json:",omitempty"`

	RtCtrl string `json:",omitempty"`
}

func NewL3outStaticRoute(ipRoutePRn, parentDn, description string, ipRoutePattr L3outStaticRouteAttributes) *L3outStaticRoute {
	dn := fmt.Sprintf("%s/%s", parentDn, ipRoutePRn)
	return &L3outStaticRoute{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IproutepClassName,
			Rn:                ipRoutePRn,
		},

		L3outStaticRouteAttributes: ipRoutePattr,
	}
}

func (ipRouteP *L3outStaticRoute) ToMap() (map[string]string, error) {
	ipRoutePMap, err := ipRouteP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ipRoutePMap, "ip", ipRouteP.Ip)

	A(ipRoutePMap, "aggregate", ipRouteP.Aggregate)

	A(ipRoutePMap, "annotation", ipRouteP.Annotation)

	A(ipRoutePMap, "nameAlias", ipRouteP.NameAlias)

	A(ipRoutePMap, "pref", ipRouteP.Pref)

	A(ipRoutePMap, "rtCtrl", ipRouteP.RtCtrl)

	return ipRoutePMap, err
}

func L3outStaticRouteFromContainerList(cont *container.Container, index int) *L3outStaticRoute {

	L3outStaticRouteCont := cont.S("imdata").Index(index).S(IproutepClassName, "attributes")
	return &L3outStaticRoute{
		BaseAttributes{
			DistinguishedName: G(L3outStaticRouteCont, "dn"),
			Description:       G(L3outStaticRouteCont, "descr"),
			Status:            G(L3outStaticRouteCont, "status"),
			ClassName:         IproutepClassName,
			Rn:                G(L3outStaticRouteCont, "rn"),
		},

		L3outStaticRouteAttributes{

			Ip: G(L3outStaticRouteCont, "ip"),

			Aggregate: G(L3outStaticRouteCont, "aggregate"),

			Annotation: G(L3outStaticRouteCont, "annotation"),

			NameAlias: G(L3outStaticRouteCont, "nameAlias"),

			Pref: G(L3outStaticRouteCont, "pref"),

			RtCtrl: G(L3outStaticRouteCont, "rtCtrl"),
		},
	}
}

func L3outStaticRouteFromContainer(cont *container.Container) *L3outStaticRoute {

	return L3outStaticRouteFromContainerList(cont, 0)
}

func L3outStaticRouteListFromContainer(cont *container.Container) []*L3outStaticRoute {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outStaticRoute, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outStaticRouteFromContainerList(cont, i)
	}

	return arr
}
