package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgplocalasnpClassName = "bgpLocalAsnP"

type LocalAutonomousSystemProfile struct {
	BaseAttributes
	LocalAutonomousSystemProfileAttributes
}

type LocalAutonomousSystemProfileAttributes struct {
	Annotation string `json:",omitempty"`

	AsnPropagate string `json:",omitempty"`

	LocalAsn string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLocalAutonomousSystemProfile(bgpLocalAsnPRn, parentDn, description string, bgpLocalAsnPattr LocalAutonomousSystemProfileAttributes) *LocalAutonomousSystemProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpLocalAsnPRn)
	return &LocalAutonomousSystemProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgplocalasnpClassName,
			Rn:                bgpLocalAsnPRn,
		},

		LocalAutonomousSystemProfileAttributes: bgpLocalAsnPattr,
	}
}

func (bgpLocalAsnP *LocalAutonomousSystemProfile) ToMap() (map[string]string, error) {
	bgpLocalAsnPMap, err := bgpLocalAsnP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpLocalAsnPMap, "annotation", bgpLocalAsnP.Annotation)

	A(bgpLocalAsnPMap, "asnPropagate", bgpLocalAsnP.AsnPropagate)

	A(bgpLocalAsnPMap, "localAsn", bgpLocalAsnP.LocalAsn)

	A(bgpLocalAsnPMap, "nameAlias", bgpLocalAsnP.NameAlias)

	return bgpLocalAsnPMap, err
}

func LocalAutonomousSystemProfileFromContainerList(cont *container.Container, index int) *LocalAutonomousSystemProfile {

	LocalAutonomousSystemProfileCont := cont.S("imdata").Index(index).S(BgplocalasnpClassName, "attributes")
	return &LocalAutonomousSystemProfile{
		BaseAttributes{
			DistinguishedName: G(LocalAutonomousSystemProfileCont, "dn"),
			Description:       G(LocalAutonomousSystemProfileCont, "descr"),
			Status:            G(LocalAutonomousSystemProfileCont, "status"),
			ClassName:         BgplocalasnpClassName,
			Rn:                G(LocalAutonomousSystemProfileCont, "rn"),
		},

		LocalAutonomousSystemProfileAttributes{

			Annotation: G(LocalAutonomousSystemProfileCont, "annotation"),

			AsnPropagate: G(LocalAutonomousSystemProfileCont, "asnPropagate"),

			LocalAsn: G(LocalAutonomousSystemProfileCont, "localAsn"),

			NameAlias: G(LocalAutonomousSystemProfileCont, "nameAlias"),
		},
	}
}

func LocalAutonomousSystemProfileFromContainer(cont *container.Container) *LocalAutonomousSystemProfile {

	return LocalAutonomousSystemProfileFromContainerList(cont, 0)
}

func LocalAutonomousSystemProfileListFromContainer(cont *container.Container) []*LocalAutonomousSystemProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LocalAutonomousSystemProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LocalAutonomousSystemProfileFromContainerList(cont, i)
	}

	return arr
}
