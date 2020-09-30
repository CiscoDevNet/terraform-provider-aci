package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const AaadomainClassName = "aaaDomain"

type SecurityDomain struct {
	BaseAttributes
	SecurityDomainAttributes
}

type SecurityDomainAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSecurityDomain(aaaDomainRn, parentDn, description string, aaaDomainattr SecurityDomainAttributes) *SecurityDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaDomainRn)
	return &SecurityDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaadomainClassName,
			Rn:                aaaDomainRn,
		},

		SecurityDomainAttributes: aaaDomainattr,
	}
}

func (aaaDomain *SecurityDomain) ToMap() (map[string]string, error) {
	aaaDomainMap, err := aaaDomain.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(aaaDomainMap, "name", aaaDomain.Name)

	A(aaaDomainMap, "annotation", aaaDomain.Annotation)

	A(aaaDomainMap, "nameAlias", aaaDomain.NameAlias)

	return aaaDomainMap, err
}

func SecurityDomainFromContainerList(cont *container.Container, index int) *SecurityDomain {

	SecurityDomainCont := cont.S("imdata").Index(index).S(AaadomainClassName, "attributes")
	return &SecurityDomain{
		BaseAttributes{
			DistinguishedName: G(SecurityDomainCont, "dn"),
			Description:       G(SecurityDomainCont, "descr"),
			Status:            G(SecurityDomainCont, "status"),
			ClassName:         AaadomainClassName,
			Rn:                G(SecurityDomainCont, "rn"),
		},

		SecurityDomainAttributes{

			Name: G(SecurityDomainCont, "name"),

			Annotation: G(SecurityDomainCont, "annotation"),

			NameAlias: G(SecurityDomainCont, "nameAlias"),
		},
	}
}

func SecurityDomainFromContainer(cont *container.Container) *SecurityDomain {

	return SecurityDomainFromContainerList(cont, 0)
}

func SecurityDomainListFromContainer(cont *container.Container) []*SecurityDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SecurityDomain, length)

	for i := 0; i < length; i++ {

		arr[i] = SecurityDomainFromContainerList(cont, i)
	}

	return arr
}
