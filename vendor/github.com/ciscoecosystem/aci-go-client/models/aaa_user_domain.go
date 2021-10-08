package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaUserDomain        = "uni/userext/user-%s/userdomain-%s"
	RnaaaUserDomain        = "userdomain-%s"
	ParentDnaaaUserDomain  = "uni/userext/user-%s"
	AaauserdomainClassName = "aaaUserDomain"
)

type UserDomain struct {
	BaseAttributes
	NameAliasAttribute
	UserDomainAttributes
}

type UserDomainAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewUserDomain(aaaUserDomainRn, parentDn, description, nameAlias string, aaaUserDomainAttr UserDomainAttributes) *UserDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaUserDomainRn)
	return &UserDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaauserdomainClassName,
			Rn:                aaaUserDomainRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		UserDomainAttributes: aaaUserDomainAttr,
	}
}

func (aaaUserDomain *UserDomain) ToMap() (map[string]string, error) {
	aaaUserDomainMap, err := aaaUserDomain.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaUserDomain.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaUserDomainMap, key, value)
	}
	A(aaaUserDomainMap, "annotation", aaaUserDomain.Annotation)
	A(aaaUserDomainMap, "name", aaaUserDomain.Name)
	return aaaUserDomainMap, err
}

func UserDomainFromContainerList(cont *container.Container, index int) *UserDomain {
	UserDomainCont := cont.S("imdata").Index(index).S(AaauserdomainClassName, "attributes")
	return &UserDomain{
		BaseAttributes{
			DistinguishedName: G(UserDomainCont, "dn"),
			Description:       G(UserDomainCont, "descr"),
			Status:            G(UserDomainCont, "status"),
			ClassName:         AaauserdomainClassName,
			Rn:                G(UserDomainCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(UserDomainCont, "nameAlias"),
		},
		UserDomainAttributes{
			Annotation: G(UserDomainCont, "annotation"),
			Name:       G(UserDomainCont, "name"),
		},
	}
}

func UserDomainFromContainer(cont *container.Container) *UserDomain {
	return UserDomainFromContainerList(cont, 0)
}

func UserDomainListFromContainer(cont *container.Container) []*UserDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*UserDomain, length)
	for i := 0; i < length; i++ {
		arr[i] = UserDomainFromContainerList(cont, i)
	}
	return arr
}
