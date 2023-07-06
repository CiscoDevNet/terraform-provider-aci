package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimIfP        = "pimifp"
	DnPimIfP        = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s/pimifp"
	ParentDnPimIfP  = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s"
	PimIfPClassName = "pimIfP"
)

type PIMInterfaceProfile struct {
	BaseAttributes
	PIMInterfaceProfileAttributes
}

type PIMInterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMInterfaceProfile(parentDn, description string, pimIfPAttr PIMInterfaceProfileAttributes) *PIMInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimIfP)
	return &PIMInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimIfPClassName,
			Rn:                RnPimIfP,
		},
		PIMInterfaceProfileAttributes: pimIfPAttr,
	}
}

func (pimIfP *PIMInterfaceProfile) ToMap() (map[string]string, error) {
	pimIfPMap, err := pimIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimIfPMap, "annotation", pimIfP.Annotation)
	A(pimIfPMap, "name", pimIfP.Name)
	A(pimIfPMap, "nameAlias", pimIfP.NameAlias)
	return pimIfPMap, err
}

func PIMInterfaceProfileFromContainerList(cont *container.Container, index int) *PIMInterfaceProfile {
	InterfaceProfileCont := cont.S("imdata").Index(index).S(PimIfPClassName, "attributes")
	return &PIMInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         PimIfPClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},
		PIMInterfaceProfileAttributes{
			Annotation: G(InterfaceProfileCont, "annotation"),
			Name:       G(InterfaceProfileCont, "name"),
			NameAlias:  G(InterfaceProfileCont, "nameAlias"),
		},
	}
}

func PIMInterfaceProfileFromContainer(cont *container.Container) *PIMInterfaceProfile {
	return PIMInterfaceProfileFromContainerList(cont, 0)
}

func PIMInterfaceProfileListFromContainer(cont *container.Container) []*PIMInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMInterfaceProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
