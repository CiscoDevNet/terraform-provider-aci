package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfragenericClassName = "infraGeneric"

type AccessGeneric struct {
	BaseAttributes
	AccessGenericAttributes
}

type AccessGenericAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewAccessGeneric(infraGenericRn, parentDn, description string, infraGenericattr AccessGenericAttributes) *AccessGeneric {
	dn := fmt.Sprintf("%s/%s", parentDn, infraGenericRn)
	return &AccessGeneric{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfragenericClassName,
			Rn:                infraGenericRn,
		},

		AccessGenericAttributes: infraGenericattr,
	}
}

func (infraGeneric *AccessGeneric) ToMap() (map[string]string, error) {
	infraGenericMap, err := infraGeneric.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraGenericMap, "name", infraGeneric.Name)

	A(infraGenericMap, "annotation", infraGeneric.Annotation)

	A(infraGenericMap, "nameAlias", infraGeneric.NameAlias)

	return infraGenericMap, err
}

func AccessGenericFromContainerList(cont *container.Container, index int) *AccessGeneric {

	AccessGenericCont := cont.S("imdata").Index(index).S(InfragenericClassName, "attributes")
	return &AccessGeneric{
		BaseAttributes{
			DistinguishedName: G(AccessGenericCont, "dn"),
			Description:       G(AccessGenericCont, "descr"),
			Status:            G(AccessGenericCont, "status"),
			ClassName:         InfragenericClassName,
			Rn:                G(AccessGenericCont, "rn"),
		},

		AccessGenericAttributes{

			Name: G(AccessGenericCont, "name"),

			Annotation: G(AccessGenericCont, "annotation"),

			NameAlias: G(AccessGenericCont, "nameAlias"),
		},
	}
}

func AccessGenericFromContainer(cont *container.Container) *AccessGeneric {

	return AccessGenericFromContainerList(cont, 0)
}

func AccessGenericListFromContainer(cont *container.Container) []*AccessGeneric {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AccessGeneric, length)

	for i := 0; i < length; i++ {

		arr[i] = AccessGenericFromContainerList(cont, i)
	}

	return arr
}
