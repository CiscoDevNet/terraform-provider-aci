package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpaspClassName = "bgpAsP"

type BgpAutonomousSystemProfile struct {
	BaseAttributes
	BgpAutonomousSystemProfileAttributes
}

type BgpAutonomousSystemProfileAttributes struct {
	Annotation string `json:",omitempty"`

	Asn string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewBgpAutonomousSystemProfile(bgpAsPRn, parentDn, description string, bgpAsPattr BgpAutonomousSystemProfileAttributes) *BgpAutonomousSystemProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpAsPRn)
	return &BgpAutonomousSystemProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgpaspClassName,
			Rn:                bgpAsPRn,
		},

		BgpAutonomousSystemProfileAttributes: bgpAsPattr,
	}
}

func (bgpAsP *BgpAutonomousSystemProfile) ToMap() (map[string]string, error) {
	bgpAsPMap, err := bgpAsP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpAsPMap, "annotation", bgpAsP.Annotation)

	A(bgpAsPMap, "asn", bgpAsP.Asn)

	A(bgpAsPMap, "nameAlias", bgpAsP.NameAlias)

	return bgpAsPMap, err
}

func BgpAutonomousSystemProfileFromContainerList(cont *container.Container, index int) *BgpAutonomousSystemProfile {

	BgpAutonomousSystemProfileCont := cont.S("imdata").Index(index).S(BgpaspClassName, "attributes")
	return &BgpAutonomousSystemProfile{
		BaseAttributes{
			DistinguishedName: G(BgpAutonomousSystemProfileCont, "dn"),
			Description:       G(BgpAutonomousSystemProfileCont, "descr"),
			Status:            G(BgpAutonomousSystemProfileCont, "status"),
			ClassName:         BgpaspClassName,
			Rn:                G(BgpAutonomousSystemProfileCont, "rn"),
		},

		BgpAutonomousSystemProfileAttributes{

			Annotation: G(BgpAutonomousSystemProfileCont, "annotation"),

			Asn: G(BgpAutonomousSystemProfileCont, "asn"),

			NameAlias: G(BgpAutonomousSystemProfileCont, "nameAlias"),
		},
	}
}

func BgpAutonomousSystemProfileFromContainer(cont *container.Container) *BgpAutonomousSystemProfile {

	return BgpAutonomousSystemProfileFromContainerList(cont, 0)
}

func BgpAutonomousSystemProfileListFromContainer(cont *container.Container) []*BgpAutonomousSystemProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BgpAutonomousSystemProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = BgpAutonomousSystemProfileFromContainerList(cont, i)
	}

	return arr
}
