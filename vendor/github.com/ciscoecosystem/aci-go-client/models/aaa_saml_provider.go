package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaSamlProvider        = "uni/userext/samlext/samlprovider-%s"
	RnaaaSamlProvider        = "samlprovider-%s"
	ParentDnaaaSamlProvider  = "uni/userext/samlext"
	AaasamlproviderClassName = "aaaSamlProvider"
)

type SAMLProvider struct {
	BaseAttributes
	NameAliasAttribute
	SAMLProviderAttributes
}

type SAMLProviderAttributes struct {
	Annotation              string `json:",omitempty"`
	EntityId                string `json:",omitempty"`
	GuiBannerMessage        string `json:",omitempty"`
	HttpsProxy              string `json:",omitempty"`
	IdP                     string `json:",omitempty"`
	Key                     string `json:",omitempty"`
	MetadataUrl             string `json:",omitempty"`
	MonitorServer           string `json:",omitempty"`
	MonitoringPassword      string `json:",omitempty"`
	MonitoringUser          string `json:",omitempty"`
	Name                    string `json:",omitempty"`
	Retries                 string `json:",omitempty"`
	SigAlg                  string `json:",omitempty"`
	Timeout                 string `json:",omitempty"`
	Tp                      string `json:",omitempty"`
	WantAssertionsEncrypted string `json:",omitempty"`
	WantAssertionsSigned    string `json:",omitempty"`
	WantRequestsSigned      string `json:",omitempty"`
	WantResponseSigned      string `json:",omitempty"`
}

func NewSAMLProvider(aaaSamlProviderRn, parentDn, description, nameAlias string, aaaSamlProviderAttr SAMLProviderAttributes) *SAMLProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaSamlProviderRn)
	return &SAMLProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaasamlproviderClassName,
			Rn:                aaaSamlProviderRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SAMLProviderAttributes: aaaSamlProviderAttr,
	}
}

func (aaaSamlProvider *SAMLProvider) ToMap() (map[string]string, error) {
	aaaSamlProviderMap, err := aaaSamlProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaSamlProvider.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaSamlProviderMap, key, value)
	}
	A(aaaSamlProviderMap, "annotation", aaaSamlProvider.Annotation)
	A(aaaSamlProviderMap, "entityId", aaaSamlProvider.EntityId)
	A(aaaSamlProviderMap, "guiBannerMessage", aaaSamlProvider.GuiBannerMessage)
	A(aaaSamlProviderMap, "httpsProxy", aaaSamlProvider.HttpsProxy)
	A(aaaSamlProviderMap, "idP", aaaSamlProvider.IdP)
	A(aaaSamlProviderMap, "key", aaaSamlProvider.Key)
	A(aaaSamlProviderMap, "metadataUrl", aaaSamlProvider.MetadataUrl)
	A(aaaSamlProviderMap, "monitorServer", aaaSamlProvider.MonitorServer)
	A(aaaSamlProviderMap, "monitoringPassword", aaaSamlProvider.MonitoringPassword)
	A(aaaSamlProviderMap, "monitoringUser", aaaSamlProvider.MonitoringUser)
	A(aaaSamlProviderMap, "name", aaaSamlProvider.Name)
	A(aaaSamlProviderMap, "retries", aaaSamlProvider.Retries)
	A(aaaSamlProviderMap, "sigAlg", aaaSamlProvider.SigAlg)
	A(aaaSamlProviderMap, "timeout", aaaSamlProvider.Timeout)
	A(aaaSamlProviderMap, "tp", aaaSamlProvider.Tp)
	A(aaaSamlProviderMap, "wantAssertionsEncrypted", aaaSamlProvider.WantAssertionsEncrypted)
	A(aaaSamlProviderMap, "wantAssertionsSigned", aaaSamlProvider.WantAssertionsSigned)
	A(aaaSamlProviderMap, "wantRequestsSigned", aaaSamlProvider.WantRequestsSigned)
	A(aaaSamlProviderMap, "wantResponseSigned", aaaSamlProvider.WantResponseSigned)
	return aaaSamlProviderMap, err
}

func SAMLProviderFromContainerList(cont *container.Container, index int) *SAMLProvider {
	SAMLProviderCont := cont.S("imdata").Index(index).S(AaasamlproviderClassName, "attributes")
	return &SAMLProvider{
		BaseAttributes{
			DistinguishedName: G(SAMLProviderCont, "dn"),
			Description:       G(SAMLProviderCont, "descr"),
			Status:            G(SAMLProviderCont, "status"),
			ClassName:         AaasamlproviderClassName,
			Rn:                G(SAMLProviderCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SAMLProviderCont, "nameAlias"),
		},
		SAMLProviderAttributes{
			Annotation:              G(SAMLProviderCont, "annotation"),
			EntityId:                G(SAMLProviderCont, "entityId"),
			GuiBannerMessage:        G(SAMLProviderCont, "guiBannerMessage"),
			HttpsProxy:              G(SAMLProviderCont, "httpsProxy"),
			IdP:                     G(SAMLProviderCont, "idP"),
			Key:                     G(SAMLProviderCont, "key"),
			MetadataUrl:             G(SAMLProviderCont, "metadataUrl"),
			MonitorServer:           G(SAMLProviderCont, "monitorServer"),
			MonitoringPassword:      G(SAMLProviderCont, "monitoringPassword"),
			MonitoringUser:          G(SAMLProviderCont, "monitoringUser"),
			Name:                    G(SAMLProviderCont, "name"),
			Retries:                 G(SAMLProviderCont, "retries"),
			SigAlg:                  G(SAMLProviderCont, "sigAlg"),
			Timeout:                 G(SAMLProviderCont, "timeout"),
			Tp:                      G(SAMLProviderCont, "tp"),
			WantAssertionsEncrypted: G(SAMLProviderCont, "wantAssertionsEncrypted"),
			WantAssertionsSigned:    G(SAMLProviderCont, "wantAssertionsSigned"),
			WantRequestsSigned:      G(SAMLProviderCont, "wantRequestsSigned"),
			WantResponseSigned:      G(SAMLProviderCont, "wantResponseSigned"),
		},
	}
}

func SAMLProviderFromContainer(cont *container.Container) *SAMLProvider {
	return SAMLProviderFromContainerList(cont, 0)
}

func SAMLProviderListFromContainer(cont *container.Container) []*SAMLProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SAMLProvider, length)
	for i := 0; i < length; i++ {
		arr[i] = SAMLProviderFromContainerList(cont, i)
	}
	return arr
}
