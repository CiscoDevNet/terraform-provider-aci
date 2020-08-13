package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraspinesClassName = "infraSpineS"

type SwitchSpineAssociation struct {
	BaseAttributes
	SwitchSpineAssociationAttributes
}

type SwitchSpineAssociationAttributes struct {
	Name string `json:",omitempty"`

	SwitchAssociationType string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSwitchSpineAssociation(infraSpineSRn, parentDn, description string, infraSpineSattr SwitchSpineAssociationAttributes) *SwitchSpineAssociation {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSpineSRn)
	return &SwitchSpineAssociation{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraspinesClassName,
			Rn:                infraSpineSRn,
		},

		SwitchSpineAssociationAttributes: infraSpineSattr,
	}
}

func (infraSpineS *SwitchSpineAssociation) ToMap() (map[string]string, error) {
	infraSpineSMap, err := infraSpineS.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraSpineSMap, "name", infraSpineS.Name)

	A(infraSpineSMap, "type", infraSpineS.SwitchAssociationType)

	A(infraSpineSMap, "annotation", infraSpineS.Annotation)

	A(infraSpineSMap, "nameAlias", infraSpineS.NameAlias)

	return infraSpineSMap, err
}

func SwitchSpineAssociationFromContainerList(cont *container.Container, index int) *SwitchSpineAssociation {

	SwitchAssociationCont := cont.S("imdata").Index(index).S(InfraspinesClassName, "attributes")
	return &SwitchSpineAssociation{
		BaseAttributes{
			DistinguishedName: G(SwitchAssociationCont, "dn"),
			Description:       G(SwitchAssociationCont, "descr"),
			Status:            G(SwitchAssociationCont, "status"),
			ClassName:         InfraspinesClassName,
			Rn:                G(SwitchAssociationCont, "rn"),
		},

		SwitchSpineAssociationAttributes{

			Name: G(SwitchAssociationCont, "name"),

			SwitchAssociationType: G(SwitchAssociationCont, "type"),

			Annotation: G(SwitchAssociationCont, "annotation"),

			NameAlias: G(SwitchAssociationCont, "nameAlias"),
		},
	}
}

func SwitchSpineAssociationFromContainer(cont *container.Container) *SwitchSpineAssociation {

	return SwitchSpineAssociationFromContainerList(cont, 0)
}

func SwitchSpineAssociationListFromContainer(cont *container.Container) []*SwitchSpineAssociation {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SwitchSpineAssociation, length)

	for i := 0; i < length; i++ {

		arr[i] = SwitchSpineAssociationFromContainerList(cont, i)
	}

	return arr
}
