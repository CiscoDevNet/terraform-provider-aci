package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaRsaProvider        = "uni/userext/rsaext/rsaprovider-%s"
	RnaaaRsaProvider        = "rsaprovider-%s"
	ParentDnaaaRsaProvider  = "uni/userext/rsaext"
	AaarsaproviderClassName = "aaaRsaProvider"
)

type RSAProvider struct {
	BaseAttributes
	NameAliasAttribute
	RSAProviderAttributes
}

type RSAProviderAttributes struct {
	Annotation         string `json:",omitempty"`
	AuthPort           string `json:",omitempty"`
	AuthProtocol       string `json:",omitempty"`
	Key                string `json:",omitempty"`
	MonitorServer      string `json:",omitempty"`
	MonitoringPassword string `json:",omitempty"`
	MonitoringUser     string `json:",omitempty"`
	Name               string `json:",omitempty"`
	Retries            string `json:",omitempty"`
	Timeout            string `json:",omitempty"`
}

func NewRSAProvider(aaaRsaProviderRn, parentDn, description, nameAlias string, aaaRsaProviderAttr RSAProviderAttributes) *RSAProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaRsaProviderRn)
	return &RSAProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaarsaproviderClassName,
			Rn:                aaaRsaProviderRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RSAProviderAttributes: aaaRsaProviderAttr,
	}
}

func (aaaRsaProvider *RSAProvider) ToMap() (map[string]string, error) {
	aaaRsaProviderMap, err := aaaRsaProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaRsaProvider.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaRsaProviderMap, key, value)
	}
	A(aaaRsaProviderMap, "annotation", aaaRsaProvider.Annotation)
	A(aaaRsaProviderMap, "authPort", aaaRsaProvider.AuthPort)
	A(aaaRsaProviderMap, "authProtocol", aaaRsaProvider.AuthProtocol)
	A(aaaRsaProviderMap, "key", aaaRsaProvider.Key)
	A(aaaRsaProviderMap, "monitorServer", aaaRsaProvider.MonitorServer)
	A(aaaRsaProviderMap, "monitoringPassword", aaaRsaProvider.MonitoringPassword)
	A(aaaRsaProviderMap, "monitoringUser", aaaRsaProvider.MonitoringUser)
	A(aaaRsaProviderMap, "name", aaaRsaProvider.Name)
	A(aaaRsaProviderMap, "retries", aaaRsaProvider.Retries)
	A(aaaRsaProviderMap, "timeout", aaaRsaProvider.Timeout)
	return aaaRsaProviderMap, err
}

func RSAProviderFromContainerList(cont *container.Container, index int) *RSAProvider {
	RSAProviderCont := cont.S("imdata").Index(index).S(AaarsaproviderClassName, "attributes")
	return &RSAProvider{
		BaseAttributes{
			DistinguishedName: G(RSAProviderCont, "dn"),
			Description:       G(RSAProviderCont, "descr"),
			Status:            G(RSAProviderCont, "status"),
			ClassName:         AaarsaproviderClassName,
			Rn:                G(RSAProviderCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RSAProviderCont, "nameAlias"),
		},
		RSAProviderAttributes{
			Annotation:         G(RSAProviderCont, "annotation"),
			AuthPort:           G(RSAProviderCont, "authPort"),
			AuthProtocol:       G(RSAProviderCont, "authProtocol"),
			Key:                G(RSAProviderCont, "key"),
			MonitorServer:      G(RSAProviderCont, "monitorServer"),
			MonitoringPassword: G(RSAProviderCont, "monitoringPassword"),
			MonitoringUser:     G(RSAProviderCont, "monitoringUser"),
			Name:               G(RSAProviderCont, "name"),
			Retries:            G(RSAProviderCont, "retries"),
			Timeout:            G(RSAProviderCont, "timeout"),
		},
	}
}

func RSAProviderFromContainer(cont *container.Container) *RSAProvider {
	return RSAProviderFromContainerList(cont, 0)
}

func RSAProviderListFromContainer(cont *container.Container) []*RSAProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RSAProvider, length)
	for i := 0; i < length; i++ {
		arr[i] = RSAProviderFromContainerList(cont, i)
	}
	return arr
}
