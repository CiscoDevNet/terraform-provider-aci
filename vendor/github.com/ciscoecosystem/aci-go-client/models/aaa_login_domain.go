package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaLoginDomain        = "uni/userext/logindomain-%s"
	RnaaaLoginDomain        = "logindomain-%s"
	ParentDnaaaLoginDomain  = "uni/userext"
	AaalogindomainClassName = "aaaLoginDomain"
)

type LoginDomain struct {
	BaseAttributes
	NameAliasAttribute
	LoginDomainAttributes
}

type LoginDomainAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewLoginDomain(aaaLoginDomainRn, parentDn, description, nameAlias string, aaaLoginDomainAttr LoginDomainAttributes) *LoginDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaLoginDomainRn)
	return &LoginDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaalogindomainClassName,
			Rn:                aaaLoginDomainRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LoginDomainAttributes: aaaLoginDomainAttr,
	}
}

func (aaaLoginDomain *LoginDomain) ToMap() (map[string]string, error) {
	aaaLoginDomainMap, err := aaaLoginDomain.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaLoginDomain.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaLoginDomainMap, key, value)
	}
	A(aaaLoginDomainMap, "annotation", aaaLoginDomain.Annotation)
	A(aaaLoginDomainMap, "name", aaaLoginDomain.Name)
	return aaaLoginDomainMap, err
}

func LoginDomainFromContainerList(cont *container.Container, index int) *LoginDomain {
	LoginDomainCont := cont.S("imdata").Index(index).S(AaalogindomainClassName, "attributes")
	return &LoginDomain{
		BaseAttributes{
			DistinguishedName: G(LoginDomainCont, "dn"),
			Description:       G(LoginDomainCont, "descr"),
			Status:            G(LoginDomainCont, "status"),
			ClassName:         AaalogindomainClassName,
			Rn:                G(LoginDomainCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LoginDomainCont, "nameAlias"),
		},
		LoginDomainAttributes{
			Annotation: G(LoginDomainCont, "annotation"),
			Name:       G(LoginDomainCont, "name"),
		},
	}
}

func LoginDomainFromContainer(cont *container.Container) *LoginDomain {
	return LoginDomainFromContainerList(cont, 0)
}

func LoginDomainListFromContainer(cont *container.Container) []*LoginDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LoginDomain, length)
	for i := 0; i < length; i++ {
		arr[i] = LoginDomainFromContainerList(cont, i)
	}
	return arr
}
