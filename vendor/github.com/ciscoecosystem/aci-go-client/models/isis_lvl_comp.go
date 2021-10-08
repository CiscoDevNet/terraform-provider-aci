package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnisisLvlComp        = "uni/fabric/isisDomP-%s/lvl-%s"
	RnisisLvlComp        = "lvl-%s"
	ParentDnisisLvlComp  = "uni/fabric/isisDomP-%s"
	IsislvlcompClassName = "isisLvlComp"
)

type ISISLevel struct {
	BaseAttributes
	NameAliasAttribute
	ISISLevelAttributes
}

type ISISLevelAttributes struct {
	Annotation       string `json:",omitempty"`
	LspFastFlood     string `json:",omitempty"`
	LspGenInitIntvl  string `json:",omitempty"`
	LspGenMaxIntvl   string `json:",omitempty"`
	LspGenSecIntvl   string `json:",omitempty"`
	Name             string `json:",omitempty"`
	SpfCompInitIntvl string `json:",omitempty"`
	SpfCompMaxIntvl  string `json:",omitempty"`
	SpfCompSecIntvl  string `json:",omitempty"`
	ISISLevel_type   string `json:",omitempty"`
}

func NewISISLevel(isisLvlCompRn, parentDn, description, nameAlias string, isisLvlCompAttr ISISLevelAttributes) *ISISLevel {
	dn := fmt.Sprintf("%s/%s", parentDn, isisLvlCompRn)
	return &ISISLevel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IsislvlcompClassName,
			Rn:                isisLvlCompRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ISISLevelAttributes: isisLvlCompAttr,
	}
}

func (isisLvlComp *ISISLevel) ToMap() (map[string]string, error) {
	isisLvlCompMap, err := isisLvlComp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := isisLvlComp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(isisLvlCompMap, key, value)
	}
	A(isisLvlCompMap, "annotation", isisLvlComp.Annotation)
	A(isisLvlCompMap, "lspFastFlood", isisLvlComp.LspFastFlood)
	A(isisLvlCompMap, "lspGenInitIntvl", isisLvlComp.LspGenInitIntvl)
	A(isisLvlCompMap, "lspGenMaxIntvl", isisLvlComp.LspGenMaxIntvl)
	A(isisLvlCompMap, "lspGenSecIntvl", isisLvlComp.LspGenSecIntvl)
	A(isisLvlCompMap, "name", isisLvlComp.Name)
	A(isisLvlCompMap, "spfCompInitIntvl", isisLvlComp.SpfCompInitIntvl)
	A(isisLvlCompMap, "spfCompMaxIntvl", isisLvlComp.SpfCompMaxIntvl)
	A(isisLvlCompMap, "spfCompSecIntvl", isisLvlComp.SpfCompSecIntvl)
	A(isisLvlCompMap, "type", isisLvlComp.ISISLevel_type)
	return isisLvlCompMap, err
}

func ISISLevelFromContainerList(cont *container.Container, index int) *ISISLevel {
	ISISLevelCont := cont.S("imdata").Index(index).S(IsislvlcompClassName, "attributes")
	return &ISISLevel{
		BaseAttributes{
			DistinguishedName: G(ISISLevelCont, "dn"),
			Description:       G(ISISLevelCont, "descr"),
			Status:            G(ISISLevelCont, "status"),
			ClassName:         IsislvlcompClassName,
			Rn:                G(ISISLevelCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ISISLevelCont, "nameAlias"),
		},
		ISISLevelAttributes{
			Annotation:       G(ISISLevelCont, "annotation"),
			LspFastFlood:     G(ISISLevelCont, "lspFastFlood"),
			LspGenInitIntvl:  G(ISISLevelCont, "lspGenInitIntvl"),
			LspGenMaxIntvl:   G(ISISLevelCont, "lspGenMaxIntvl"),
			LspGenSecIntvl:   G(ISISLevelCont, "lspGenSecIntvl"),
			Name:             G(ISISLevelCont, "name"),
			SpfCompInitIntvl: G(ISISLevelCont, "spfCompInitIntvl"),
			SpfCompMaxIntvl:  G(ISISLevelCont, "spfCompMaxIntvl"),
			SpfCompSecIntvl:  G(ISISLevelCont, "spfCompSecIntvl"),
			ISISLevel_type:   G(ISISLevelCont, "type"),
		},
	}
}

func ISISLevelFromContainer(cont *container.Container) *ISISLevel {
	return ISISLevelFromContainerList(cont, 0)
}

func ISISLevelListFromContainer(cont *container.Container) []*ISISLevel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ISISLevel, length)
	for i := 0; i < length; i++ {
		arr[i] = ISISLevelFromContainerList(cont, i)
	}
	return arr
}
