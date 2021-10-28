package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaTacacsPlusProvider        = "uni/userext/tacacsext/tacacsplusprovider-%s"
	RnaaaTacacsPlusProvider        = "tacacsplusprovider-%s"
	ParentDnaaaTacacsPlusProvider  = "uni/userext/tacacsext"
	AaatacacsplusproviderClassName = "aaaTacacsPlusProvider"
)

type TACACSProvider struct {
	BaseAttributes
	NameAliasAttribute
	TACACSProviderAttributes
}
type TACACSProviderAttributes struct {
	Annotation         string `json:",omitempty"`
	AuthProtocol       string `json:",omitempty"`
	Key                string `json:",omitempty"`
	MonitorServer      string `json:",omitempty"`
	MonitoringPassword string `json:",omitempty"`
	MonitoringUser     string `json:",omitempty"`
	Name               string `json:",omitempty"`
	Port               string `json:",omitempty"`
	Retries            string `json:",omitempty"`
	Timeout            string `json:",omitempty"`
}

func NewTACACSProvider(aaaTacacsPlusProviderRn, parentDn, description, nameAlias string, aaaTacacsPlusProviderAttr TACACSProviderAttributes) *TACACSProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaTacacsPlusProviderRn)
	return &TACACSProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaatacacsplusproviderClassName,
			Rn:                aaaTacacsPlusProviderRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TACACSProviderAttributes: aaaTacacsPlusProviderAttr,
	}
}

func (aaaTacacsPlusProvider *TACACSProvider) ToMap() (map[string]string, error) {
	aaaTacacsPlusProviderMap, err := aaaTacacsPlusProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaTacacsPlusProvider.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaTacacsPlusProviderMap, key, value)
	}
	A(aaaTacacsPlusProviderMap, "annotation", aaaTacacsPlusProvider.Annotation)
	A(aaaTacacsPlusProviderMap, "authProtocol", aaaTacacsPlusProvider.AuthProtocol)
	A(aaaTacacsPlusProviderMap, "key", aaaTacacsPlusProvider.Key)
	A(aaaTacacsPlusProviderMap, "monitorServer", aaaTacacsPlusProvider.MonitorServer)
	A(aaaTacacsPlusProviderMap, "monitoringPassword", aaaTacacsPlusProvider.MonitoringPassword)
	A(aaaTacacsPlusProviderMap, "monitoringUser", aaaTacacsPlusProvider.MonitoringUser)
	A(aaaTacacsPlusProviderMap, "name", aaaTacacsPlusProvider.Name)
	A(aaaTacacsPlusProviderMap, "port", aaaTacacsPlusProvider.Port)
	A(aaaTacacsPlusProviderMap, "retries", aaaTacacsPlusProvider.Retries)
	A(aaaTacacsPlusProviderMap, "timeout", aaaTacacsPlusProvider.Timeout)
	return aaaTacacsPlusProviderMap, err
}

func TACACSProviderFromContainerList(cont *container.Container, index int) *TACACSProvider {
	TACACSProviderCont := cont.S("imdata").Index(index).S(AaatacacsplusproviderClassName, "attributes")
	return &TACACSProvider{
		BaseAttributes{
			DistinguishedName: G(TACACSProviderCont, "dn"),
			Description:       G(TACACSProviderCont, "descr"),
			Status:            G(TACACSProviderCont, "status"),
			ClassName:         AaatacacsplusproviderClassName,
			Rn:                G(TACACSProviderCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TACACSProviderCont, "nameAlias"),
		},
		TACACSProviderAttributes{
			Annotation:         G(TACACSProviderCont, "annotation"),
			AuthProtocol:       G(TACACSProviderCont, "authProtocol"),
			Key:                G(TACACSProviderCont, "key"),
			MonitorServer:      G(TACACSProviderCont, "monitorServer"),
			MonitoringPassword: G(TACACSProviderCont, "monitoringPassword"),
			MonitoringUser:     G(TACACSProviderCont, "monitoringUser"),
			Name:               G(TACACSProviderCont, "name"),
			Port:               G(TACACSProviderCont, "port"),
			Retries:            G(TACACSProviderCont, "retries"),
			Timeout:            G(TACACSProviderCont, "timeout"),
		},
	}
}

func TACACSProviderFromContainer(cont *container.Container) *TACACSProvider {
	return TACACSProviderFromContainerList(cont, 0)
}

func TACACSProviderListFromContainer(cont *container.Container) []*TACACSProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TACACSProvider, length)
	for i := 0; i < length; i++ {
		arr[i] = TACACSProviderFromContainerList(cont, i)
	}
	return arr
}
