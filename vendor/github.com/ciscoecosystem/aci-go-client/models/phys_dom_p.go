package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const PhysdompClassName = "physDomP"

type PhysicalDomain struct {
	BaseAttributes
	PhysicalDomainAttributes
}

type PhysicalDomainAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewPhysicalDomain(physDomPRn, parentDn, description string, physDomPattr PhysicalDomainAttributes) *PhysicalDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, physDomPRn)
	return &PhysicalDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PhysdompClassName,
			Rn:                physDomPRn,
		},

		PhysicalDomainAttributes: physDomPattr,
	}
}

func (physDomP *PhysicalDomain) ToMap() (map[string]string, error) {
	physDomPMap, err := physDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(physDomPMap, "name", physDomP.Name)

	A(physDomPMap, "annotation", physDomP.Annotation)

	A(physDomPMap, "nameAlias", physDomP.NameAlias)

	return physDomPMap, err
}

func PhysicalDomainFromContainerList(cont *container.Container, index int) *PhysicalDomain {

	PhysicalDomainCont := cont.S("imdata").Index(index).S(PhysdompClassName, "attributes")
	return &PhysicalDomain{
		BaseAttributes{
			DistinguishedName: G(PhysicalDomainCont, "dn"),
			Description:       G(PhysicalDomainCont, "descr"),
			Status:            G(PhysicalDomainCont, "status"),
			ClassName:         PhysdompClassName,
			Rn:                G(PhysicalDomainCont, "rn"),
		},

		PhysicalDomainAttributes{

			Name: G(PhysicalDomainCont, "name"),

			Annotation: G(PhysicalDomainCont, "annotation"),

			NameAlias: G(PhysicalDomainCont, "nameAlias"),
		},
	}
}

func PhysicalDomainFromContainer(cont *container.Container) *PhysicalDomain {

	return PhysicalDomainFromContainerList(cont, 0)
}

func PhysicalDomainListFromContainer(cont *container.Container) []*PhysicalDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*PhysicalDomain, length)

	for i := 0; i < length; i++ {

		arr[i] = PhysicalDomainFromContainerList(cont, i)
	}

	return arr
}
