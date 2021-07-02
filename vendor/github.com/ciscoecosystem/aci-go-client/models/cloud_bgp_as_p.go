package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudbgpaspClassName = "cloudBgpAsP"

type AutonomousSystemProfile struct {
	BaseAttributes
	AutonomousSystemProfileAttributes
}

type AutonomousSystemProfileAttributes struct {
	Annotation string `json:",omitempty"`

	Asn string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewAutonomousSystemProfile(cloudBgpAsPRn, parentDn, description string, cloudBgpAsPattr AutonomousSystemProfileAttributes) *AutonomousSystemProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudBgpAsPRn)
	return &AutonomousSystemProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudbgpaspClassName,
			Rn:                cloudBgpAsPRn,
		},

		AutonomousSystemProfileAttributes: cloudBgpAsPattr,
	}
}

func (cloudBgpAsP *AutonomousSystemProfile) ToMap() (map[string]string, error) {
	cloudBgpAsPMap, err := cloudBgpAsP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudBgpAsPMap, "annotation", cloudBgpAsP.Annotation)

	A(cloudBgpAsPMap, "asn", cloudBgpAsP.Asn)

	A(cloudBgpAsPMap, "nameAlias", cloudBgpAsP.NameAlias)

	return cloudBgpAsPMap, err
}

func AutonomousSystemProfileFromContainerList(cont *container.Container, index int) *AutonomousSystemProfile {

	AutonomousSystemProfileCont := cont.S("imdata").Index(index).S(CloudbgpaspClassName, "attributes")
	return &AutonomousSystemProfile{
		BaseAttributes{
			DistinguishedName: G(AutonomousSystemProfileCont, "dn"),
			Description:       G(AutonomousSystemProfileCont, "descr"),
			Status:            G(AutonomousSystemProfileCont, "status"),
			ClassName:         CloudbgpaspClassName,
			Rn:                G(AutonomousSystemProfileCont, "rn"),
		},

		AutonomousSystemProfileAttributes{

			Annotation: G(AutonomousSystemProfileCont, "annotation"),

			Asn: G(AutonomousSystemProfileCont, "asn"),

			NameAlias: G(AutonomousSystemProfileCont, "nameAlias"),
		},
	}
}

func AutonomousSystemProfileFromContainer(cont *container.Container) *AutonomousSystemProfile {

	return AutonomousSystemProfileFromContainerList(cont, 0)
}

func AutonomousSystemProfileListFromContainer(cont *container.Container) []*AutonomousSystemProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AutonomousSystemProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = AutonomousSystemProfileFromContainerList(cont, i)
	}

	return arr
}
