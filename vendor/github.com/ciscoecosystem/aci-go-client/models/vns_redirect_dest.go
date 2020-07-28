package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsredirectdestClassName = "vnsRedirectDest"

type Destinationofredirectedtraffic struct {
	BaseAttributes
	DestinationofredirectedtrafficAttributes
}

type DestinationofredirectedtrafficAttributes struct {
	Ip string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	DestName string `json:",omitempty"`

	Ip2 string `json:",omitempty"`

	Mac string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PodId string `json:",omitempty"`
}

func NewDestinationofredirectedtraffic(vnsRedirectDestRn, parentDn, description string, vnsRedirectDestattr DestinationofredirectedtrafficAttributes) *Destinationofredirectedtraffic {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsRedirectDestRn)
	return &Destinationofredirectedtraffic{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsredirectdestClassName,
			Rn:                vnsRedirectDestRn,
		},

		DestinationofredirectedtrafficAttributes: vnsRedirectDestattr,
	}
}

func (vnsRedirectDest *Destinationofredirectedtraffic) ToMap() (map[string]string, error) {
	vnsRedirectDestMap, err := vnsRedirectDest.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsRedirectDestMap, "ip", vnsRedirectDest.Ip)

	A(vnsRedirectDestMap, "annotation", vnsRedirectDest.Annotation)

	A(vnsRedirectDestMap, "destName", vnsRedirectDest.DestName)

	A(vnsRedirectDestMap, "ip2", vnsRedirectDest.Ip2)

	A(vnsRedirectDestMap, "mac", vnsRedirectDest.Mac)

	A(vnsRedirectDestMap, "nameAlias", vnsRedirectDest.NameAlias)

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

		DestinationofredirectedtrafficAttributes{

			Ip: G(DestinationofredirectedtrafficCont, "ip"),

			Annotation: G(DestinationofredirectedtrafficCont, "annotation"),

			DestName: G(DestinationofredirectedtrafficCont, "destName"),

			Ip2: G(DestinationofredirectedtrafficCont, "ip2"),

			Mac: G(DestinationofredirectedtrafficCont, "mac"),

			NameAlias: G(DestinationofredirectedtrafficCont, "nameAlias"),

			PodId: G(DestinationofredirectedtrafficCont, "podId"),
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
