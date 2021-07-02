package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraaccportpClassName = "infraAccPortP"

type LeafInterfaceProfile struct {
	BaseAttributes
	LeafInterfaceProfileAttributes
}

type LeafInterfaceProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLeafInterfaceProfile(infraAccPortPRn, parentDn, description string, infraAccPortPattr LeafInterfaceProfileAttributes) *LeafInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAccPortPRn)
	return &LeafInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraaccportpClassName,
			Rn:                infraAccPortPRn,
		},

		LeafInterfaceProfileAttributes: infraAccPortPattr,
	}
}

func (infraAccPortP *LeafInterfaceProfile) ToMap() (map[string]string, error) {
	infraAccPortPMap, err := infraAccPortP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraAccPortPMap, "name", infraAccPortP.Name)

	A(infraAccPortPMap, "annotation", infraAccPortP.Annotation)

	A(infraAccPortPMap, "nameAlias", infraAccPortP.NameAlias)

	return infraAccPortPMap, err
}

func LeafInterfaceProfileFromContainerList(cont *container.Container, index int) *LeafInterfaceProfile {

	LeafInterfaceProfileCont := cont.S("imdata").Index(index).S(InfraaccportpClassName, "attributes")
	return &LeafInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(LeafInterfaceProfileCont, "dn"),
			Description:       G(LeafInterfaceProfileCont, "descr"),
			Status:            G(LeafInterfaceProfileCont, "status"),
			ClassName:         InfraaccportpClassName,
			Rn:                G(LeafInterfaceProfileCont, "rn"),
		},

		LeafInterfaceProfileAttributes{

			Name: G(LeafInterfaceProfileCont, "name"),

			Annotation: G(LeafInterfaceProfileCont, "annotation"),

			NameAlias: G(LeafInterfaceProfileCont, "nameAlias"),
		},
	}
}

func LeafInterfaceProfileFromContainer(cont *container.Container) *LeafInterfaceProfile {

	return LeafInterfaceProfileFromContainerList(cont, 0)
}

func LeafInterfaceProfileListFromContainer(cont *container.Container) []*LeafInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LeafInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LeafInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
