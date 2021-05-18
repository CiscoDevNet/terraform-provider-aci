package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const HsrpifpClassName = "hsrpIfP"

type L3outHSRPInterfaceProfile struct {
	BaseAttributes
	L3outHSRPInterfaceProfileAttributes
}

type L3outHSRPInterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Version string `json:",omitempty"`
}

func NewL3outHSRPInterfaceProfile(hsrpIfPRn, parentDn, description string, hsrpIfPattr L3outHSRPInterfaceProfileAttributes) *L3outHSRPInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, hsrpIfPRn)
	return &L3outHSRPInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         HsrpifpClassName,
			Rn:                hsrpIfPRn,
		},

		L3outHSRPInterfaceProfileAttributes: hsrpIfPattr,
	}
}

func (hsrpIfP *L3outHSRPInterfaceProfile) ToMap() (map[string]string, error) {
	hsrpIfPMap, err := hsrpIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(hsrpIfPMap, "annotation", hsrpIfP.Annotation)

	A(hsrpIfPMap, "nameAlias", hsrpIfP.NameAlias)

	A(hsrpIfPMap, "version", hsrpIfP.Version)

	return hsrpIfPMap, err
}

func L3outHSRPInterfaceProfileFromContainerList(cont *container.Container, index int) *L3outHSRPInterfaceProfile {

	L3outHSRPInterfaceProfileCont := cont.S("imdata").Index(index).S(HsrpifpClassName, "attributes")
	return &L3outHSRPInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(L3outHSRPInterfaceProfileCont, "dn"),
			Description:       G(L3outHSRPInterfaceProfileCont, "descr"),
			Status:            G(L3outHSRPInterfaceProfileCont, "status"),
			ClassName:         HsrpifpClassName,
			Rn:                G(L3outHSRPInterfaceProfileCont, "rn"),
		},

		L3outHSRPInterfaceProfileAttributes{

			Annotation: G(L3outHSRPInterfaceProfileCont, "annotation"),

			NameAlias: G(L3outHSRPInterfaceProfileCont, "nameAlias"),

			Version: G(L3outHSRPInterfaceProfileCont, "version"),
		},
	}
}

func L3outHSRPInterfaceProfileFromContainer(cont *container.Container) *L3outHSRPInterfaceProfile {

	return L3outHSRPInterfaceProfileFromContainerList(cont, 0)
}

func L3outHSRPInterfaceProfileListFromContainer(cont *container.Container) []*L3outHSRPInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outHSRPInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outHSRPInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
