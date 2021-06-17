package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FcdompClassName = "fcDomP"

type FCDomain struct {
	BaseAttributes
	FCDomainAttributes
}

type FCDomainAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewFCDomain(fcDomPRn, parentDn, description string, fcDomPattr FCDomainAttributes) *FCDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, fcDomPRn)
	return &FCDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FcdompClassName,
			Rn:                fcDomPRn,
		},

		FCDomainAttributes: fcDomPattr,
	}
}

func (fcDomP *FCDomain) ToMap() (map[string]string, error) {
	fcDomPMap, err := fcDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fcDomPMap, "name", fcDomP.Name)

	A(fcDomPMap, "annotation", fcDomP.Annotation)

	A(fcDomPMap, "nameAlias", fcDomP.NameAlias)

	return fcDomPMap, err
}

func FCDomainFromContainerList(cont *container.Container, index int) *FCDomain {

	FCDomainCont := cont.S("imdata").Index(index).S(FcdompClassName, "attributes")
	return &FCDomain{
		BaseAttributes{
			DistinguishedName: G(FCDomainCont, "dn"),
			Description:       G(FCDomainCont, "descr"),
			Status:            G(FCDomainCont, "status"),
			ClassName:         FcdompClassName,
			Rn:                G(FCDomainCont, "rn"),
		},

		FCDomainAttributes{

			Name: G(FCDomainCont, "name"),

			Annotation: G(FCDomainCont, "annotation"),

			NameAlias: G(FCDomainCont, "nameAlias"),
		},
	}
}

func FCDomainFromContainer(cont *container.Container) *FCDomain {

	return FCDomainFromContainerList(cont, 0)
}

func FCDomainListFromContainer(cont *container.Container) []*FCDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FCDomain, length)

	for i := 0; i < length; i++ {

		arr[i] = FCDomainFromContainerList(cont, i)
	}

	return arr
}
