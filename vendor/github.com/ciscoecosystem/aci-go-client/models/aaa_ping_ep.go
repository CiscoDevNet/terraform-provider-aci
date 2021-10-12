package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaPingEp        = "uni/userext/pingext"
	RnaaaPingEp        = "pingext"
	ParentDnaaaPingEp  = "uni/userext"
	AaapingepClassName = "aaaPingEp"
)

type DefaultRadiusAuthenticationSettings struct {
	BaseAttributes
	NameAliasAttribute
	DefaultRadiusAuthenticationSettingsAttributes
}

type DefaultRadiusAuthenticationSettingsAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	PingCheck  string `json:",omitempty"`
	Retries    string `json:",omitempty"`
	Timeout    string `json:",omitempty"`
}

func NewDefaultRadiusAuthenticationSettings(aaaPingEpRn, parentDn, description, nameAlias string, aaaPingEpAttr DefaultRadiusAuthenticationSettingsAttributes) *DefaultRadiusAuthenticationSettings {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaPingEpRn)
	return &DefaultRadiusAuthenticationSettings{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaapingepClassName,
			Rn:                aaaPingEpRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		DefaultRadiusAuthenticationSettingsAttributes: aaaPingEpAttr,
	}
}

func (aaaPingEp *DefaultRadiusAuthenticationSettings) ToMap() (map[string]string, error) {
	aaaPingEpMap, err := aaaPingEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaPingEp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaPingEpMap, key, value)
	}
	A(aaaPingEpMap, "annotation", aaaPingEp.Annotation)
	A(aaaPingEpMap, "name", aaaPingEp.Name)
	A(aaaPingEpMap, "pingCheck", aaaPingEp.PingCheck)
	A(aaaPingEpMap, "retries", aaaPingEp.Retries)
	A(aaaPingEpMap, "timeout", aaaPingEp.Timeout)
	return aaaPingEpMap, err
}

func DefaultRadiusAuthenticationSettingsFromContainerList(cont *container.Container, index int) *DefaultRadiusAuthenticationSettings {
	DefaultRadiusAuthenticationSettingsCont := cont.S("imdata").Index(index).S(AaapingepClassName, "attributes")
	return &DefaultRadiusAuthenticationSettings{
		BaseAttributes{
			DistinguishedName: G(DefaultRadiusAuthenticationSettingsCont, "dn"),
			Description:       G(DefaultRadiusAuthenticationSettingsCont, "descr"),
			Status:            G(DefaultRadiusAuthenticationSettingsCont, "status"),
			ClassName:         AaapingepClassName,
			Rn:                G(DefaultRadiusAuthenticationSettingsCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(DefaultRadiusAuthenticationSettingsCont, "nameAlias"),
		},
		DefaultRadiusAuthenticationSettingsAttributes{
			Annotation: G(DefaultRadiusAuthenticationSettingsCont, "annotation"),
			Name:       G(DefaultRadiusAuthenticationSettingsCont, "name"),
			PingCheck:  G(DefaultRadiusAuthenticationSettingsCont, "pingCheck"),
			Retries:    G(DefaultRadiusAuthenticationSettingsCont, "retries"),
			Timeout:    G(DefaultRadiusAuthenticationSettingsCont, "timeout"),
		},
	}
}

func DefaultRadiusAuthenticationSettingsFromContainer(cont *container.Container) *DefaultRadiusAuthenticationSettings {
	return DefaultRadiusAuthenticationSettingsFromContainerList(cont, 0)
}

func DefaultRadiusAuthenticationSettingsListFromContainer(cont *container.Container) []*DefaultRadiusAuthenticationSettings {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*DefaultRadiusAuthenticationSettings, length)
	for i := 0; i < length; i++ {
		arr[i] = DefaultRadiusAuthenticationSettingsFromContainerList(cont, i)
	}
	return arr
}
