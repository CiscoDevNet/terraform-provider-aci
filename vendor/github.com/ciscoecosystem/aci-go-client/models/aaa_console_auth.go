package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaConsoleAuth        = "uni/userext/authrealm/consoleauth"
	RnaaaConsoleAuth        = "consoleauth"
	ParentDnaaaConsoleAuth  = "uni/userext/authrealm"
	AaaconsoleauthClassName = "aaaConsoleAuth"
)

type ConsoleAuthenticationMethod struct {
	BaseAttributes
	NameAliasAttribute
	ConsoleAuthenticationMethodAttributes
}

type ConsoleAuthenticationMethodAttributes struct {
	Annotation    string `json:",omitempty"`
	Name          string `json:",omitempty"`
	ProviderGroup string `json:",omitempty"`
	Realm         string `json:",omitempty"`
	RealmSubType  string `json:",omitempty"`
}

func NewConsoleAuthenticationMethod(aaaConsoleAuthRn, parentDn, description, nameAlias string, aaaConsoleAuthAttr ConsoleAuthenticationMethodAttributes) *ConsoleAuthenticationMethod {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaConsoleAuthRn)
	return &ConsoleAuthenticationMethod{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaconsoleauthClassName,
			Rn:                aaaConsoleAuthRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ConsoleAuthenticationMethodAttributes: aaaConsoleAuthAttr,
	}
}

func (aaaConsoleAuth *ConsoleAuthenticationMethod) ToMap() (map[string]string, error) {
	aaaConsoleAuthMap, err := aaaConsoleAuth.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaConsoleAuth.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaConsoleAuthMap, key, value)
	}
	A(aaaConsoleAuthMap, "annotation", aaaConsoleAuth.Annotation)
	A(aaaConsoleAuthMap, "name", aaaConsoleAuth.Name)
	A(aaaConsoleAuthMap, "providerGroup", aaaConsoleAuth.ProviderGroup)
	A(aaaConsoleAuthMap, "realm", aaaConsoleAuth.Realm)
	A(aaaConsoleAuthMap, "realmSubType", aaaConsoleAuth.RealmSubType)
	return aaaConsoleAuthMap, err
}

func ConsoleAuthenticationMethodFromContainerList(cont *container.Container, index int) *ConsoleAuthenticationMethod {
	ConsoleAuthenticationMethodCont := cont.S("imdata").Index(index).S(AaaconsoleauthClassName, "attributes")
	return &ConsoleAuthenticationMethod{
		BaseAttributes{
			DistinguishedName: G(ConsoleAuthenticationMethodCont, "dn"),
			Description:       G(ConsoleAuthenticationMethodCont, "descr"),
			Status:            G(ConsoleAuthenticationMethodCont, "status"),
			ClassName:         AaaconsoleauthClassName,
			Rn:                G(ConsoleAuthenticationMethodCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ConsoleAuthenticationMethodCont, "nameAlias"),
		},
		ConsoleAuthenticationMethodAttributes{
			Annotation:    G(ConsoleAuthenticationMethodCont, "annotation"),
			Name:          G(ConsoleAuthenticationMethodCont, "name"),
			ProviderGroup: G(ConsoleAuthenticationMethodCont, "providerGroup"),
			Realm:         G(ConsoleAuthenticationMethodCont, "realm"),
			RealmSubType:  G(ConsoleAuthenticationMethodCont, "realmSubType"),
		},
	}
}

func ConsoleAuthenticationMethodFromContainer(cont *container.Container) *ConsoleAuthenticationMethod {
	return ConsoleAuthenticationMethodFromContainerList(cont, 0)
}

func ConsoleAuthenticationMethodListFromContainer(cont *container.Container) []*ConsoleAuthenticationMethod {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ConsoleAuthenticationMethod, length)
	for i := 0; i < length; i++ {
		arr[i] = ConsoleAuthenticationMethodFromContainerList(cont, i)
	}
	return arr
}
