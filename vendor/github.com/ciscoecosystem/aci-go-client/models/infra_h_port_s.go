package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrahportsClassName = "infraHPortS"

type AccessPortSelector struct {
	BaseAttributes
	AccessPortSelectorAttributes
}

type AccessPortSelectorAttributes struct {
	Name string `json:",omitempty"`

	AccessPortSelector_type string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewAccessPortSelector(infraHPortSRn, parentDn, description string, infraHPortSattr AccessPortSelectorAttributes) *AccessPortSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, infraHPortSRn)
	return &AccessPortSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrahportsClassName,
			Rn:                infraHPortSRn,
		},

		AccessPortSelectorAttributes: infraHPortSattr,
	}
}

func (infraHPortS *AccessPortSelector) ToMap() (map[string]string, error) {
	infraHPortSMap, err := infraHPortS.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraHPortSMap, "name", infraHPortS.Name)

	A(infraHPortSMap, "type", infraHPortS.AccessPortSelector_type)

	A(infraHPortSMap, "annotation", infraHPortS.Annotation)

	A(infraHPortSMap, "nameAlias", infraHPortS.NameAlias)

	return infraHPortSMap, err
}

func AccessPortSelectorFromContainerList(cont *container.Container, index int) *AccessPortSelector {

	AccessPortSelectorCont := cont.S("imdata").Index(index).S(InfrahportsClassName, "attributes")
	return &AccessPortSelector{
		BaseAttributes{
			DistinguishedName: G(AccessPortSelectorCont, "dn"),
			Description:       G(AccessPortSelectorCont, "descr"),
			Status:            G(AccessPortSelectorCont, "status"),
			ClassName:         InfrahportsClassName,
			Rn:                G(AccessPortSelectorCont, "rn"),
		},

		AccessPortSelectorAttributes{

			Name: G(AccessPortSelectorCont, "name"),

			AccessPortSelector_type: G(AccessPortSelectorCont, "type"),

			Annotation: G(AccessPortSelectorCont, "annotation"),

			NameAlias: G(AccessPortSelectorCont, "nameAlias"),
		},
	}
}

func AccessPortSelectorFromContainer(cont *container.Container) *AccessPortSelector {

	return AccessPortSelectorFromContainerList(cont, 0)
}

func AccessPortSelectorListFromContainer(cont *container.Container) []*AccessPortSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AccessPortSelector, length)

	for i := 0; i < length; i++ {

		arr[i] = AccessPortSelectorFromContainerList(cont, i)
	}

	return arr
}
