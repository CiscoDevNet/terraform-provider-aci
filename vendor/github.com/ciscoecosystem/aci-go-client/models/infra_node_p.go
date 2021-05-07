package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfranodepClassName = "infraNodeP"

type LeafProfile struct {
	BaseAttributes
	LeafProfileAttributes
}

type LeafProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLeafProfile(infraNodePRn, parentDn, description string, infraNodePattr LeafProfileAttributes) *LeafProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraNodePRn)
	return &LeafProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfranodepClassName,
			Rn:                infraNodePRn,
		},

		LeafProfileAttributes: infraNodePattr,
	}
}

func (infraNodeP *LeafProfile) ToMap() (map[string]string, error) {
	infraNodePMap, err := infraNodeP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraNodePMap, "name", infraNodeP.Name)

	A(infraNodePMap, "annotation", infraNodeP.Annotation)

	A(infraNodePMap, "nameAlias", infraNodeP.NameAlias)

	return infraNodePMap, err
}

func LeafProfileFromContainerList(cont *container.Container, index int) *LeafProfile {

	LeafProfileCont := cont.S("imdata").Index(index).S(InfranodepClassName, "attributes")
	return &LeafProfile{
		BaseAttributes{
			DistinguishedName: G(LeafProfileCont, "dn"),
			Description:       G(LeafProfileCont, "descr"),
			Status:            G(LeafProfileCont, "status"),
			ClassName:         InfranodepClassName,
			Rn:                G(LeafProfileCont, "rn"),
		},

		LeafProfileAttributes{

			Name: G(LeafProfileCont, "name"),

			Annotation: G(LeafProfileCont, "annotation"),

			NameAlias: G(LeafProfileCont, "nameAlias"),
		},
	}
}

func LeafProfileFromContainer(cont *container.Container) *LeafProfile {

	return LeafProfileFromContainerList(cont, 0)
}

func LeafProfileListFromContainer(cont *container.Container) []*LeafProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LeafProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LeafProfileFromContainerList(cont, i)
	}

	return arr
}
