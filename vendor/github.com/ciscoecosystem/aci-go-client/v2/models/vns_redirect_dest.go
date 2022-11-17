package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnvnsRedirectDest        = "RedirectDest_ip-[%s]"
	VnsredirectdestClassName = "vnsRedirectDest"
)

type Destinationofredirectedtraffic struct {
	BaseAttributes
	NameAliasAttribute
	DestinationofredirectedtrafficAttributes
}

type DestinationofredirectedtrafficAttributes struct {
	Annotation string `json:",omitempty"`
	DestName   string `json:",omitempty"`
	Ip         string `json:",omitempty"`
	Ip2        string `json:",omitempty"`
	Mac        string `json:",omitempty"`
	Name       string `json:",omitempty"`
	PodId      string `json:",omitempty"`
}

func NewDestinationofredirectedtraffic(vnsRedirectDestRn, parentDn, description, nameAlias string, vnsRedirectDestAttr DestinationofredirectedtrafficAttributes) *Destinationofredirectedtraffic {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsRedirectDestRn)
	return &Destinationofredirectedtraffic{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsredirectdestClassName,
			Rn:                vnsRedirectDestRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		DestinationofredirectedtrafficAttributes: vnsRedirectDestAttr,
	}
}

func (vnsRedirectDest *Destinationofredirectedtraffic) ToMap() (map[string]string, error) {
	vnsRedirectDestMap, err := vnsRedirectDest.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsRedirectDest.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsRedirectDestMap, key, value)
	}

	A(vnsRedirectDestMap, "annotation", vnsRedirectDest.Annotation)
	A(vnsRedirectDestMap, "destName", vnsRedirectDest.DestName)
	A(vnsRedirectDestMap, "ip", vnsRedirectDest.Ip)
	A(vnsRedirectDestMap, "ip2", vnsRedirectDest.Ip2)
	A(vnsRedirectDestMap, "mac", vnsRedirectDest.Mac)
	A(vnsRedirectDestMap, "name", vnsRedirectDest.Name)
	A(vnsRedirectDestMap, "podId", vnsRedirectDest.PodId)
	return vnsRedirectDestMap, err
}

func DestinationofredirectedtrafficFromContainerList(cont *container.Container, index int) *Destinationofredirectedtraffic {
	DestinationofredirectedtrafficCont := cont.S("imdata").Index(index).S(VnsredirectdestClassName, "attributes")
	return &Destinationofredirectedtraffic{
		BaseAttributes{
			DistinguishedName: G(DestinationofredirectedtrafficCont, "dn"),
			Description:       G(DestinationofredirectedtrafficCont, "descr"),
			Status:            G(DestinationofredirectedtrafficCont, "status"),
			ClassName:         VnsredirectdestClassName,
			Rn:                G(DestinationofredirectedtrafficCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(DestinationofredirectedtrafficCont, "nameAlias"),
		},
		DestinationofredirectedtrafficAttributes{
			Annotation: G(DestinationofredirectedtrafficCont, "annotation"),
			DestName:   G(DestinationofredirectedtrafficCont, "destName"),
			Ip:         G(DestinationofredirectedtrafficCont, "ip"),
			Ip2:        G(DestinationofredirectedtrafficCont, "ip2"),
			Mac:        G(DestinationofredirectedtrafficCont, "mac"),
			Name:       G(DestinationofredirectedtrafficCont, "name"),
			PodId:      G(DestinationofredirectedtrafficCont, "podId"),
		},
	}
}

func DestinationofredirectedtrafficFromContainer(cont *container.Container) *Destinationofredirectedtraffic {
	return DestinationofredirectedtrafficFromContainerList(cont, 0)
}

func DestinationofredirectedtrafficListFromContainer(cont *container.Container) []*Destinationofredirectedtraffic {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*Destinationofredirectedtraffic, length)

	for i := 0; i < length; i++ {
		arr[i] = DestinationofredirectedtrafficFromContainerList(cont, i)
	}

	return arr
}
