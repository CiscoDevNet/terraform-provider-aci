package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrarsdompClassName = "infraRsDomP"

type Domain struct {
	BaseAttributes
	DomainAttributes
}

type DomainAttributes struct {
	TDn string `json:",omitempty"`

	Annotation string `json:",omitempty"`
}

func NewDomain(infraRsDomPRn, parentDn, description string, infraRsDomPattr DomainAttributes) *Domain {
	dn := fmt.Sprintf("%s/%s", parentDn, infraRsDomPRn)
	return &Domain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrarsdompClassName,
			Rn:                infraRsDomPRn,
		},

		DomainAttributes: infraRsDomPattr,
	}
}

func (infraRsDomP *Domain) ToMap() (map[string]string, error) {
	infraRsDomPMap, err := infraRsDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraRsDomPMap, "tDn", infraRsDomP.TDn)

	A(infraRsDomPMap, "annotation", infraRsDomP.Annotation)

	return infraRsDomPMap, err
}

func DomainFromContainerList(cont *container.Container, index int) *Domain {

	DomainCont := cont.S("imdata").Index(index).S(InfrarsdompClassName, "attributes")
	return &Domain{
		BaseAttributes{
			DistinguishedName: G(DomainCont, "dn"),
			Description:       G(DomainCont, "descr"),
			Status:            G(DomainCont, "status"),
			ClassName:         InfrarsdompClassName,
			Rn:                G(DomainCont, "rn"),
		},

		DomainAttributes{

			TDn: G(DomainCont, "tDn"),

			Annotation: G(DomainCont, "annotation"),
		},
	}
}

func DomainFromContainer(cont *container.Container) *Domain {

	return DomainFromContainerList(cont, 0)
}

func DomainListFromContainer(cont *container.Container) []*Domain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Domain, length)

	for i := 0; i < length; i++ {

		arr[i] = DomainFromContainerList(cont, i)
	}

	return arr
}
