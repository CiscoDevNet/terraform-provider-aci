package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L2extdompClassName = "l2extDomP"

type L2Domain struct {
	BaseAttributes
	L2DomainAttributes
}

type L2DomainAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL2Domain(l2extDomPRn, parentDn string, l2extDomPattr L2DomainAttributes) *L2Domain {
	dn := fmt.Sprintf("%s/%s", parentDn, l2extDomPRn)
	return &L2Domain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         L2extdompClassName,
			Rn:                l2extDomPRn,
		},

		L2DomainAttributes: l2extDomPattr,
	}
}

func (l2extDomP *L2Domain) ToMap() (map[string]string, error) {
	l2extDomPMap, err := l2extDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l2extDomPMap, "name", l2extDomP.Name)

	A(l2extDomPMap, "annotation", l2extDomP.Annotation)

	A(l2extDomPMap, "nameAlias", l2extDomP.NameAlias)

	return l2extDomPMap, err
}

func L2DomainFromContainerList(cont *container.Container, index int) *L2Domain {

	L2DomainCont := cont.S("imdata").Index(index).S(L2extdompClassName, "attributes")
	return &L2Domain{
		BaseAttributes{
			DistinguishedName: G(L2DomainCont, "dn"),
			Status:            G(L2DomainCont, "status"),
			ClassName:         L2extdompClassName,
			Rn:                G(L2DomainCont, "rn"),
		},

		L2DomainAttributes{

			Name: G(L2DomainCont, "name"),

			Annotation: G(L2DomainCont, "annotation"),

			NameAlias: G(L2DomainCont, "nameAlias"),
		},
	}
}

func L2DomainFromContainer(cont *container.Container) *L2Domain {

	return L2DomainFromContainerList(cont, 0)
}

func L2DomainListFromContainer(cont *container.Container) []*L2Domain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L2Domain, length)

	for i := 0; i < length; i++ {

		arr[i] = L2DomainFromContainerList(cont, i)
	}

	return arr
}
