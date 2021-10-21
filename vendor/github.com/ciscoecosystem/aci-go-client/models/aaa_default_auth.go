package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaDefaultAuth        = "uni/userext/authrealm/defaultauth"
	RnaaaDefaultAuth        = "defaultauth"
	ParentDnaaaDefaultAuth  = "uni/userext/authrealm"
	AaadefaultauthClassName = "aaaDefaultAuth"
)

type DefaultAuthenticationMethodforallLogins struct {
	BaseAttributes
	NameAliasAttribute
	DefaultAuthenticationMethodforallLoginsAttributes
}

type DefaultAuthenticationMethodforallLoginsAttributes struct {
	Annotation    string `json:",omitempty"`
	FallbackCheck string `json:",omitempty"`
	Name          string `json:",omitempty"`
	ProviderGroup string `json:",omitempty"`
	Realm         string `json:",omitempty"`
	RealmSubType  string `json:",omitempty"`
}

func NewDefaultAuthenticationMethodforallLogins(aaaDefaultAuthRn, parentDn, description, nameAlias string, aaaDefaultAuthAttr DefaultAuthenticationMethodforallLoginsAttributes) *DefaultAuthenticationMethodforallLogins {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaDefaultAuthRn)
	return &DefaultAuthenticationMethodforallLogins{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaadefaultauthClassName,
			Rn:                aaaDefaultAuthRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		DefaultAuthenticationMethodforallLoginsAttributes: aaaDefaultAuthAttr,
	}
}

func (aaaDefaultAuth *DefaultAuthenticationMethodforallLogins) ToMap() (map[string]string, error) {
	aaaDefaultAuthMap, err := aaaDefaultAuth.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaDefaultAuth.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaDefaultAuthMap, key, value)
	}
	A(aaaDefaultAuthMap, "annotation", aaaDefaultAuth.Annotation)
	A(aaaDefaultAuthMap, "fallbackCheck", aaaDefaultAuth.FallbackCheck)
	A(aaaDefaultAuthMap, "name", aaaDefaultAuth.Name)
	A(aaaDefaultAuthMap, "providerGroup", aaaDefaultAuth.ProviderGroup)
	A(aaaDefaultAuthMap, "realm", aaaDefaultAuth.Realm)
	A(aaaDefaultAuthMap, "realmSubType", aaaDefaultAuth.RealmSubType)
	return aaaDefaultAuthMap, err
}

func DefaultAuthenticationMethodforallLoginsFromContainerList(cont *container.Container, index int) *DefaultAuthenticationMethodforallLogins {
	DefaultAuthenticationMethodforallLoginsCont := cont.S("imdata").Index(index).S(AaadefaultauthClassName, "attributes")
	return &DefaultAuthenticationMethodforallLogins{
		BaseAttributes{
			DistinguishedName: G(DefaultAuthenticationMethodforallLoginsCont, "dn"),
			Description:       G(DefaultAuthenticationMethodforallLoginsCont, "descr"),
			Status:            G(DefaultAuthenticationMethodforallLoginsCont, "status"),
			ClassName:         AaadefaultauthClassName,
			Rn:                G(DefaultAuthenticationMethodforallLoginsCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(DefaultAuthenticationMethodforallLoginsCont, "nameAlias"),
		},
		DefaultAuthenticationMethodforallLoginsAttributes{
			Annotation:    G(DefaultAuthenticationMethodforallLoginsCont, "annotation"),
			FallbackCheck: G(DefaultAuthenticationMethodforallLoginsCont, "fallbackCheck"),
			Name:          G(DefaultAuthenticationMethodforallLoginsCont, "name"),
			ProviderGroup: G(DefaultAuthenticationMethodforallLoginsCont, "providerGroup"),
			Realm:         G(DefaultAuthenticationMethodforallLoginsCont, "realm"),
			RealmSubType:  G(DefaultAuthenticationMethodforallLoginsCont, "realmSubType"),
		},
	}
}

func DefaultAuthenticationMethodforallLoginsFromContainer(cont *container.Container) *DefaultAuthenticationMethodforallLogins {
	return DefaultAuthenticationMethodforallLoginsFromContainerList(cont, 0)
}

func DefaultAuthenticationMethodforallLoginsListFromContainer(cont *container.Container) []*DefaultAuthenticationMethodforallLogins {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*DefaultAuthenticationMethodforallLogins, length)
	for i := 0; i < length; i++ {
		arr[i] = DefaultAuthenticationMethodforallLoginsFromContainerList(cont, i)
	}
	return arr
}
