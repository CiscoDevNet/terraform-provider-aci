package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extdompClassName = "l3extDomP"

type L3DomainProfile struct {
	BaseAttributes
	L3DomainProfileAttributes
}

type L3DomainProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3DomainProfile(l3extDomPRn, parentDn, description string, l3extDomPattr L3DomainProfileAttributes) *L3DomainProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extDomPRn)
	return &L3DomainProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extdompClassName,
			Rn:                l3extDomPRn,
		},

		L3DomainProfileAttributes: l3extDomPattr,
	}
}

func (l3extDomP *L3DomainProfile) ToMap() (map[string]string, error) {
	l3extDomPMap, err := l3extDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extDomPMap, "name", l3extDomP.Name)

	A(l3extDomPMap, "annotation", l3extDomP.Annotation)

	A(l3extDomPMap, "nameAlias", l3extDomP.NameAlias)

	return l3extDomPMap, err
}

func L3DomainProfileFromContainerList(cont *container.Container, index int) *L3DomainProfile {

	L3DomainProfileCont := cont.S("imdata").Index(index).S(L3extdompClassName, "attributes")
	return &L3DomainProfile{
		BaseAttributes{
			DistinguishedName: G(L3DomainProfileCont, "dn"),
			Description:       G(L3DomainProfileCont, "descr"),
			Status:            G(L3DomainProfileCont, "status"),
			ClassName:         L3extdompClassName,
			Rn:                G(L3DomainProfileCont, "rn"),
		},

		L3DomainProfileAttributes{

			Name: G(L3DomainProfileCont, "name"),

			Annotation: G(L3DomainProfileCont, "annotation"),

			NameAlias: G(L3DomainProfileCont, "nameAlias"),
		},
	}
}

func L3DomainProfileFromContainer(cont *container.Container) *L3DomainProfile {

	return L3DomainProfileFromContainerList(cont, 0)
}

func L3DomainProfileListFromContainer(cont *container.Container) []*L3DomainProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3DomainProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = L3DomainProfileFromContainerList(cont, i)
	}

	return arr
}
