package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraleafsClassName = "infraLeafS"

type SwitchAssociation struct {
	BaseAttributes
	SwitchAssociationAttributes
}

type SwitchAssociationAttributes struct {
	Name string `json:",omitempty"`

	Switch_association_type string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSwitchAssociation(infraLeafSRn, parentDn, description string, infraLeafSattr SwitchAssociationAttributes) *SwitchAssociation {
	dn := fmt.Sprintf("%s/%s", parentDn, infraLeafSRn)
	return &SwitchAssociation{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraleafsClassName,
			Rn:                infraLeafSRn,
		},

		SwitchAssociationAttributes: infraLeafSattr,
	}
}

func (infraLeafS *SwitchAssociation) ToMap() (map[string]string, error) {
	infraLeafSMap, err := infraLeafS.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraLeafSMap, "name", infraLeafS.Name)

	A(infraLeafSMap, "type", infraLeafS.Switch_association_type)

	A(infraLeafSMap, "annotation", infraLeafS.Annotation)

	A(infraLeafSMap, "nameAlias", infraLeafS.NameAlias)

	return infraLeafSMap, err
}

func SwitchAssociationFromContainerList(cont *container.Container, index int) *SwitchAssociation {

	SwitchAssociationCont := cont.S("imdata").Index(index).S(InfraleafsClassName, "attributes")
	return &SwitchAssociation{
		BaseAttributes{
			DistinguishedName: G(SwitchAssociationCont, "dn"),
			Description:       G(SwitchAssociationCont, "descr"),
			Status:            G(SwitchAssociationCont, "status"),
			ClassName:         InfraleafsClassName,
			Rn:                G(SwitchAssociationCont, "rn"),
		},

		SwitchAssociationAttributes{

			Name: G(SwitchAssociationCont, "name"),

			Switch_association_type: G(SwitchAssociationCont, "type"),

			Annotation: G(SwitchAssociationCont, "annotation"),

			NameAlias: G(SwitchAssociationCont, "nameAlias"),
		},
	}
}

func SwitchAssociationFromContainer(cont *container.Container) *SwitchAssociation {

	return SwitchAssociationFromContainerList(cont, 0)
}

func SwitchAssociationListFromContainer(cont *container.Container) []*SwitchAssociation {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SwitchAssociation, length)

	for i := 0; i < length; i++ {

		arr[i] = SwitchAssociationFromContainerList(cont, i)
	}

	return arr
}
