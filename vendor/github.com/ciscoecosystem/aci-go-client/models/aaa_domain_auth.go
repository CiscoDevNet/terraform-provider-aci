package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaDomainAuth        = "uni/userext/logindomain-%s/domainauth"
	RnaaaDomainAuth        = "domainauth"
	ParentDnaaaDomainAuth  = "uni/userext/logindomain-%s"
	AaadomainauthClassName = "aaaDomainAuth"
)

type AuthenticationMethodfortheDomain struct {
	BaseAttributes
	NameAliasAttribute
	AuthenticationMethodfortheDomainAttributes
}

type AuthenticationMethodfortheDomainAttributes struct {
	Annotation    string `json:",omitempty"`
	ProviderGroup string `json:",omitempty"`
	Realm         string `json:",omitempty"`
	RealmSubType  string `json:",omitempty"`
	Name          string `json:",omitempty"`
}

func NewAuthenticationMethodfortheDomain(aaaDomainAuthRn, parentDn, description, nameAlias string, aaaDomainAuthAttr AuthenticationMethodfortheDomainAttributes) *AuthenticationMethodfortheDomain {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaDomainAuthRn)
	return &AuthenticationMethodfortheDomain{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaadomainauthClassName,
			Rn:                aaaDomainAuthRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AuthenticationMethodfortheDomainAttributes: aaaDomainAuthAttr,
	}
}

func (aaaDomainAuth *AuthenticationMethodfortheDomain) ToMap() (map[string]string, error) {
	aaaDomainAuthMap, err := aaaDomainAuth.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaDomainAuth.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaDomainAuthMap, key, value)
	}
	A(aaaDomainAuthMap, "annotation", aaaDomainAuth.Annotation)
	A(aaaDomainAuthMap, "providerGroup", aaaDomainAuth.ProviderGroup)
	A(aaaDomainAuthMap, "realm", aaaDomainAuth.Realm)
	A(aaaDomainAuthMap, "name", aaaDomainAuth.Name)
	A(aaaDomainAuthMap, "realmSubType", aaaDomainAuth.RealmSubType)
	return aaaDomainAuthMap, err
}

func AuthenticationMethodfortheDomainFromContainerList(cont *container.Container, index int) *AuthenticationMethodfortheDomain {
	AuthenticationMethodfortheDomainCont := cont.S("imdata").Index(index).S(AaadomainauthClassName, "attributes")
	return &AuthenticationMethodfortheDomain{
		BaseAttributes{
			DistinguishedName: G(AuthenticationMethodfortheDomainCont, "dn"),
			Description:       G(AuthenticationMethodfortheDomainCont, "descr"),
			Status:            G(AuthenticationMethodfortheDomainCont, "status"),
			ClassName:         AaadomainauthClassName,
			Rn:                G(AuthenticationMethodfortheDomainCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AuthenticationMethodfortheDomainCont, "nameAlias"),
		},
		AuthenticationMethodfortheDomainAttributes{
			Annotation:    G(AuthenticationMethodfortheDomainCont, "annotation"),
			ProviderGroup: G(AuthenticationMethodfortheDomainCont, "providerGroup"),
			Realm:         G(AuthenticationMethodfortheDomainCont, "realm"),
			RealmSubType:  G(AuthenticationMethodfortheDomainCont, "realmSubType"),
			Name:          G(AuthenticationMethodfortheDomainCont, "name"),
		},
	}
}

func AuthenticationMethodfortheDomainFromContainer(cont *container.Container) *AuthenticationMethodfortheDomain {
	return AuthenticationMethodfortheDomainFromContainerList(cont, 0)
}

func AuthenticationMethodfortheDomainListFromContainer(cont *container.Container) []*AuthenticationMethodfortheDomain {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AuthenticationMethodfortheDomain, length)
	for i := 0; i < length; i++ {
		arr[i] = AuthenticationMethodfortheDomainFromContainerList(cont, i)
	}
	return arr
}
