package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaRadiusProvider        = "uni/userext/duoext/radiusprovider-%s"
	RnaaaRadiusProvider        = "radiusprovider-%s"
	ParentDnaaaRadiusProvider  = "uni/userext/duoext"
	AaaradiusproviderClassName = "aaaRadiusProvider"
)

type RADIUSProvider struct {
	BaseAttributes
	NameAliasAttribute
	RADIUSProviderAttributes
}

type RADIUSProviderAttributes struct {
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

func NewRADIUSProvider(aaaRadiusProviderRn, parentDn, description, nameAlias string, aaaRadiusProviderAttr RADIUSProviderAttributes) *RADIUSProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaRadiusProviderRn)
	return &RADIUSProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaradiusproviderClassName,
			Rn:                aaaRadiusProviderRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RADIUSProviderAttributes: aaaRadiusProviderAttr,
	}
}

func (aaaRadiusProvider *RADIUSProvider) ToMap() (map[string]string, error) {
	aaaRadiusProviderMap, err := aaaRadiusProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaRadiusProvider.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaRadiusProviderMap, key, value)
	}
	A(aaaRadiusProviderMap, "annotation", aaaRadiusProvider.Annotation)
	A(aaaRadiusProviderMap, "authPort", aaaRadiusProvider.AuthPort)
	A(aaaRadiusProviderMap, "authProtocol", aaaRadiusProvider.AuthProtocol)
	A(aaaRadiusProviderMap, "key", aaaRadiusProvider.Key)
	A(aaaRadiusProviderMap, "monitorServer", aaaRadiusProvider.MonitorServer)
	A(aaaRadiusProviderMap, "monitoringPassword", aaaRadiusProvider.MonitoringPassword)
	A(aaaRadiusProviderMap, "monitoringUser", aaaRadiusProvider.MonitoringUser)
	A(aaaRadiusProviderMap, "name", aaaRadiusProvider.Name)
	A(aaaRadiusProviderMap, "retries", aaaRadiusProvider.Retries)
	A(aaaRadiusProviderMap, "timeout", aaaRadiusProvider.Timeout)
	return aaaRadiusProviderMap, err
}

func RADIUSProviderFromContainerList(cont *container.Container, index int) *RADIUSProvider {
	RADIUSProviderCont := cont.S("imdata").Index(index).S(AaaradiusproviderClassName, "attributes")
	return &RADIUSProvider{
		BaseAttributes{
			DistinguishedName: G(RADIUSProviderCont, "dn"),
			Description:       G(RADIUSProviderCont, "descr"),
			Status:            G(RADIUSProviderCont, "status"),
			ClassName:         AaaradiusproviderClassName,
			Rn:                G(RADIUSProviderCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RADIUSProviderCont, "nameAlias"),
		},
		RADIUSProviderAttributes{
			Annotation:         G(RADIUSProviderCont, "annotation"),
			AuthPort:           G(RADIUSProviderCont, "authPort"),
			AuthProtocol:       G(RADIUSProviderCont, "authProtocol"),
			Key:                G(RADIUSProviderCont, "key"),
			MonitorServer:      G(RADIUSProviderCont, "monitorServer"),
			MonitoringPassword: G(RADIUSProviderCont, "monitoringPassword"),
			MonitoringUser:     G(RADIUSProviderCont, "monitoringUser"),
			Name:               G(RADIUSProviderCont, "name"),
			Retries:            G(RADIUSProviderCont, "retries"),
			Timeout:            G(RADIUSProviderCont, "timeout"),
		},
	}
}

func RADIUSProviderFromContainer(cont *container.Container) *RADIUSProvider {
	return RADIUSProviderFromContainerList(cont, 0)
}

func RADIUSProviderListFromContainer(cont *container.Container) []*RADIUSProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RADIUSProvider, length)
	for i := 0; i < length; i++ {
		arr[i] = RADIUSProviderFromContainerList(cont, i)
	}
	return arr
}
