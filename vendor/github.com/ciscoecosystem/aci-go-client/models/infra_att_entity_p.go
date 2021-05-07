package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraattentitypClassName = "infraAttEntityP"

type AttachableAccessEntityProfile struct {
	BaseAttributes
	AttachableAccessEntityProfileAttributes
}

type AttachableAccessEntityProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewAttachableAccessEntityProfile(infraAttEntityPRn, parentDn, description string, infraAttEntityPattr AttachableAccessEntityProfileAttributes) *AttachableAccessEntityProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAttEntityPRn)
	return &AttachableAccessEntityProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraattentitypClassName,
			Rn:                infraAttEntityPRn,
		},

		AttachableAccessEntityProfileAttributes: infraAttEntityPattr,
	}
}

func (infraAttEntityP *AttachableAccessEntityProfile) ToMap() (map[string]string, error) {
	infraAttEntityPMap, err := infraAttEntityP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraAttEntityPMap, "name", infraAttEntityP.Name)

	A(infraAttEntityPMap, "annotation", infraAttEntityP.Annotation)

	A(infraAttEntityPMap, "nameAlias", infraAttEntityP.NameAlias)

	return infraAttEntityPMap, err
}

func AttachableAccessEntityProfileFromContainerList(cont *container.Container, index int) *AttachableAccessEntityProfile {

	AttachableAccessEntityProfileCont := cont.S("imdata").Index(index).S(InfraattentitypClassName, "attributes")
	return &AttachableAccessEntityProfile{
		BaseAttributes{
			DistinguishedName: G(AttachableAccessEntityProfileCont, "dn"),
			Description:       G(AttachableAccessEntityProfileCont, "descr"),
			Status:            G(AttachableAccessEntityProfileCont, "status"),
			ClassName:         InfraattentitypClassName,
			Rn:                G(AttachableAccessEntityProfileCont, "rn"),
		},

		AttachableAccessEntityProfileAttributes{

			Name: G(AttachableAccessEntityProfileCont, "name"),

			Annotation: G(AttachableAccessEntityProfileCont, "annotation"),

			NameAlias: G(AttachableAccessEntityProfileCont, "nameAlias"),
		},
	}
}

func AttachableAccessEntityProfileFromContainer(cont *container.Container) *AttachableAccessEntityProfile {

	return AttachableAccessEntityProfileFromContainerList(cont, 0)
}

func AttachableAccessEntityProfileListFromContainer(cont *container.Container) []*AttachableAccessEntityProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AttachableAccessEntityProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = AttachableAccessEntityProfileFromContainerList(cont, i)
	}

	return arr
}
