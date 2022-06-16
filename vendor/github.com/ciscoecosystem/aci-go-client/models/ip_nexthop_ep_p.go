package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	RnipNexthopEpP        = "nh-[%s]"
	IpnexthopeppClassName = "ipNexthopEpP"
)

type NexthopEpPReachability struct {
	BaseAttributes
	NameAliasAttribute
	NexthopEpPReachabilityAttributes
}

type NexthopEpPReachabilityAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NhAddr     string `json:",omitempty"`
}

func NewNexthopEpPReachability(ipNexthopEpPRn, parentDn, description, nameAlias string, ipNexthopEpPAttr NexthopEpPReachabilityAttributes) *NexthopEpPReachability {
	dn := fmt.Sprintf("%s/%s", parentDn, ipNexthopEpPRn)
	return &NexthopEpPReachability{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IpnexthopeppClassName,
			Rn:                ipNexthopEpPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		NexthopEpPReachabilityAttributes: ipNexthopEpPAttr,
	}
}

func (ipNexthopEpP *NexthopEpPReachability) ToMap() (map[string]string, error) {
	ipNexthopEpPMap, err := ipNexthopEpP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := ipNexthopEpP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(ipNexthopEpPMap, key, value)
	}

	A(ipNexthopEpPMap, "name", ipNexthopEpP.Name)
	A(ipNexthopEpPMap, "nhAddr", ipNexthopEpP.NhAddr)
	return ipNexthopEpPMap, err
}

func NexthopEpPReachabilityFromContainerList(cont *container.Container, index int) *NexthopEpPReachability {
	NexthopEpPReachabilityCont := cont.S("imdata").Index(index).S(IpnexthopeppClassName, "attributes")
	return &NexthopEpPReachability{
		BaseAttributes{
			DistinguishedName: G(NexthopEpPReachabilityCont, "dn"),
			Description:       G(NexthopEpPReachabilityCont, "descr"),
			Status:            G(NexthopEpPReachabilityCont, "status"),
			ClassName:         IpnexthopeppClassName,
			Rn:                G(NexthopEpPReachabilityCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(NexthopEpPReachabilityCont, "nameAlias"),
		},
		NexthopEpPReachabilityAttributes{
			Name:   G(NexthopEpPReachabilityCont, "name"),
			NhAddr: G(NexthopEpPReachabilityCont, "nhAddr"),
		},
	}
}

func NexthopEpPReachabilityFromContainer(cont *container.Container) *NexthopEpPReachability {
	return NexthopEpPReachabilityFromContainerList(cont, 0)
}

func NexthopEpPReachabilityListFromContainer(cont *container.Container) []*NexthopEpPReachability {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*NexthopEpPReachability, length)

	for i := 0; i < length; i++ {
		arr[i] = NexthopEpPReachabilityFromContainerList(cont, i)
	}

	return arr
}
