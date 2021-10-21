package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaLdapProvider        = "uni/userext/duoext/ldapprovider-%s"
	RnaaaLdapProvider        = "ldapprovider-%s"
	ParentDnaaaLdapProvider  = "uni/userext/duoext"
	AaaldapproviderClassName = "aaaLdapProvider"
)

type LDAPProvider struct {
	BaseAttributes
	NameAliasAttribute
	LDAPProviderAttributes
}

type LDAPProviderAttributes struct {
	SSLValidationLevel string `json:",omitempty"`
	Annotation         string `json:",omitempty"`
	Attribute          string `json:",omitempty"`
	Basedn             string `json:",omitempty"`
	EnableSSL          string `json:",omitempty"`
	Filter             string `json:",omitempty"`
	Key                string `json:",omitempty"`
	MonitorServer      string `json:",omitempty"`
	MonitoringPassword string `json:",omitempty"`
	MonitoringUser     string `json:",omitempty"`
	Name               string `json:",omitempty"`
	Port               string `json:",omitempty"`
	Retries            string `json:",omitempty"`
	Rootdn             string `json:",omitempty"`
	Timeout            string `json:",omitempty"`
}

func NewLDAPProvider(aaaLdapProviderRn, parentDn, description, nameAlias string, aaaLdapProviderAttr LDAPProviderAttributes) *LDAPProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaLdapProviderRn)
	return &LDAPProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaldapproviderClassName,
			Rn:                aaaLdapProviderRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LDAPProviderAttributes: aaaLdapProviderAttr,
	}
}

func (aaaLdapProvider *LDAPProvider) ToMap() (map[string]string, error) {
	aaaLdapProviderMap, err := aaaLdapProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaLdapProvider.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaLdapProviderMap, key, value)
	}
	A(aaaLdapProviderMap, "SSLValidationLevel", aaaLdapProvider.SSLValidationLevel)
	A(aaaLdapProviderMap, "annotation", aaaLdapProvider.Annotation)
	A(aaaLdapProviderMap, "attribute", aaaLdapProvider.Attribute)
	A(aaaLdapProviderMap, "basedn", aaaLdapProvider.Basedn)
	A(aaaLdapProviderMap, "enableSSL", aaaLdapProvider.EnableSSL)
	A(aaaLdapProviderMap, "filter", aaaLdapProvider.Filter)
	A(aaaLdapProviderMap, "key", aaaLdapProvider.Key)
	A(aaaLdapProviderMap, "monitorServer", aaaLdapProvider.MonitorServer)
	A(aaaLdapProviderMap, "monitoringPassword", aaaLdapProvider.MonitoringPassword)
	A(aaaLdapProviderMap, "monitoringUser", aaaLdapProvider.MonitoringUser)
	A(aaaLdapProviderMap, "name", aaaLdapProvider.Name)
	A(aaaLdapProviderMap, "port", aaaLdapProvider.Port)
	A(aaaLdapProviderMap, "retries", aaaLdapProvider.Retries)
	A(aaaLdapProviderMap, "rootdn", aaaLdapProvider.Rootdn)
	A(aaaLdapProviderMap, "timeout", aaaLdapProvider.Timeout)
	return aaaLdapProviderMap, err
}

func LDAPProviderFromContainerList(cont *container.Container, index int) *LDAPProvider {
	LDAPProviderCont := cont.S("imdata").Index(index).S(AaaldapproviderClassName, "attributes")
	return &LDAPProvider{
		BaseAttributes{
			DistinguishedName: G(LDAPProviderCont, "dn"),
			Description:       G(LDAPProviderCont, "descr"),
			Status:            G(LDAPProviderCont, "status"),
			ClassName:         AaaldapproviderClassName,
			Rn:                G(LDAPProviderCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LDAPProviderCont, "nameAlias"),
		},
		LDAPProviderAttributes{
			SSLValidationLevel: G(LDAPProviderCont, "SSLValidationLevel"),
			Annotation:         G(LDAPProviderCont, "annotation"),
			Attribute:          G(LDAPProviderCont, "attribute"),
			Basedn:             G(LDAPProviderCont, "basedn"),
			EnableSSL:          G(LDAPProviderCont, "enableSSL"),
			Filter:             G(LDAPProviderCont, "filter"),
			Key:                G(LDAPProviderCont, "key"),
			MonitorServer:      G(LDAPProviderCont, "monitorServer"),
			MonitoringPassword: G(LDAPProviderCont, "monitoringPassword"),
			MonitoringUser:     G(LDAPProviderCont, "monitoringUser"),
			Name:               G(LDAPProviderCont, "name"),
			Port:               G(LDAPProviderCont, "port"),
			Retries:            G(LDAPProviderCont, "retries"),
			Rootdn:             G(LDAPProviderCont, "rootdn"),
			Timeout:            G(LDAPProviderCont, "timeout"),
		},
	}
}

func LDAPProviderFromContainer(cont *container.Container) *LDAPProvider {
	return LDAPProviderFromContainerList(cont, 0)
}

func LDAPProviderListFromContainer(cont *container.Container) []*LDAPProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LDAPProvider, length)
	for i := 0; i < length; i++ {
		arr[i] = LDAPProviderFromContainerList(cont, i)
	}
	return arr
}
