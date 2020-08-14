package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraspaccportgrpClassName = "infraSpAccPortGrp"

type SpineAccessPortPolicyGroup struct {
	BaseAttributes
	SpineAccessPortPolicyGroupAttributes
}

type SpineAccessPortPolicyGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSpineAccessPortPolicyGroup(infraSpAccPortGrpRn, parentDn, description string, infraSpAccPortGrpattr SpineAccessPortPolicyGroupAttributes) *SpineAccessPortPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSpAccPortGrpRn)
	return &SpineAccessPortPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraspaccportgrpClassName,
			Rn:                infraSpAccPortGrpRn,
		},

		SpineAccessPortPolicyGroupAttributes: infraSpAccPortGrpattr,
	}
}

func (infraSpAccPortGrp *SpineAccessPortPolicyGroup) ToMap() (map[string]string, error) {
	infraSpAccPortGrpMap, err := infraSpAccPortGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraSpAccPortGrpMap, "name", infraSpAccPortGrp.Name)

	A(infraSpAccPortGrpMap, "annotation", infraSpAccPortGrp.Annotation)

	A(infraSpAccPortGrpMap, "nameAlias", infraSpAccPortGrp.NameAlias)

	return infraSpAccPortGrpMap, err
}

func SpineAccessPortPolicyGroupFromContainerList(cont *container.Container, index int) *SpineAccessPortPolicyGroup {

	SpineAccessPortPolicyGroupCont := cont.S("imdata").Index(index).S(InfraspaccportgrpClassName, "attributes")
	return &SpineAccessPortPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(SpineAccessPortPolicyGroupCont, "dn"),
			Description:       G(SpineAccessPortPolicyGroupCont, "descr"),
			Status:            G(SpineAccessPortPolicyGroupCont, "status"),
			ClassName:         InfraspaccportgrpClassName,
			Rn:                G(SpineAccessPortPolicyGroupCont, "rn"),
		},

		SpineAccessPortPolicyGroupAttributes{

			Name: G(SpineAccessPortPolicyGroupCont, "name"),

			Annotation: G(SpineAccessPortPolicyGroupCont, "annotation"),

			NameAlias: G(SpineAccessPortPolicyGroupCont, "nameAlias"),
		},
	}
}

func SpineAccessPortPolicyGroupFromContainer(cont *container.Container) *SpineAccessPortPolicyGroup {

	return SpineAccessPortPolicyGroupFromContainerList(cont, 0)
}

func SpineAccessPortPolicyGroupListFromContainer(cont *container.Container) []*SpineAccessPortPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SpineAccessPortPolicyGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = SpineAccessPortPolicyGroupFromContainerList(cont, i)
	}

	return arr
}
