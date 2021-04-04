package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpprotpClassName = "bgpProtP"

type L3outBGPProtocolProfile struct {
	BaseAttributes
	L3outBGPProtocolProfileAttributes
}

type L3outBGPProtocolProfileAttributes struct {
	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outBGPProtocolProfile(bgpProtPRn, parentDn string, bgpProtPattr L3outBGPProtocolProfileAttributes) *L3outBGPProtocolProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpProtPRn)
	return &L3outBGPProtocolProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         BgpprotpClassName,
			Rn:                bgpProtPRn,
		},

		L3outBGPProtocolProfileAttributes: bgpProtPattr,
	}
}

func (bgpProtP *L3outBGPProtocolProfile) ToMap() (map[string]string, error) {
	bgpProtPMap, err := bgpProtP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpProtPMap, "annotation", bgpProtP.Annotation)

	A(bgpProtPMap, "nameAlias", bgpProtP.NameAlias)

	return bgpProtPMap, err
}

func L3outBGPProtocolProfileFromContainerList(cont *container.Container, index int) *L3outBGPProtocolProfile {

	L3outBGPProtocolProfileCont := cont.S("imdata").Index(index).S(BgpprotpClassName, "attributes")
	return &L3outBGPProtocolProfile{
		BaseAttributes{
			DistinguishedName: G(L3outBGPProtocolProfileCont, "dn"),
			Status:            G(L3outBGPProtocolProfileCont, "status"),
			ClassName:         BgpprotpClassName,
			Rn:                G(L3outBGPProtocolProfileCont, "rn"),
		},

		L3outBGPProtocolProfileAttributes{

			Annotation: G(L3outBGPProtocolProfileCont, "annotation"),

			NameAlias: G(L3outBGPProtocolProfileCont, "nameAlias"),
		},
	}
}

func L3outBGPProtocolProfileFromContainer(cont *container.Container) *L3outBGPProtocolProfile {

	return L3outBGPProtocolProfileFromContainerList(cont, 0)
}

func L3outBGPProtocolProfileListFromContainer(cont *container.Container) []*L3outBGPProtocolProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outBGPProtocolProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outBGPProtocolProfileFromContainerList(cont, i)
	}

	return arr
}
