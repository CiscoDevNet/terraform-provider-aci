package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetWeight        = "uni/tn-%s/attr-%s/sweight"
	RnrtctrlSetWeight        = "sweight"
	ParentDnrtctrlSetWeight  = "uni/tn-%s/attr-%s"
	RtctrlsetweightClassName = "rtctrlSetWeight"
)

type RtctrlSetWeight struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetWeightAttributes
}

type RtctrlSetWeightAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
	Weight     string `json:",omitempty"`
}

func NewRtctrlSetWeight(rtctrlSetWeightRn, parentDn, description, nameAlias string, rtctrlSetWeightAttr RtctrlSetWeightAttributes) *RtctrlSetWeight {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetWeightRn)
	return &RtctrlSetWeight{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetweightClassName,
			Rn:                rtctrlSetWeightRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetWeightAttributes: rtctrlSetWeightAttr,
	}
}

func (rtctrlSetWeight *RtctrlSetWeight) ToMap() (map[string]string, error) {
	rtctrlSetWeightMap, err := rtctrlSetWeight.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetWeight.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetWeightMap, key, value)
	}

	A(rtctrlSetWeightMap, "annotation", rtctrlSetWeight.Annotation)
	A(rtctrlSetWeightMap, "name", rtctrlSetWeight.Name)
	A(rtctrlSetWeightMap, "type", rtctrlSetWeight.Type)
	A(rtctrlSetWeightMap, "weight", rtctrlSetWeight.Weight)
	return rtctrlSetWeightMap, err
}

func RtctrlSetWeightFromContainerList(cont *container.Container, index int) *RtctrlSetWeight {
	RtctrlSetWeightCont := cont.S("imdata").Index(index).S(RtctrlsetweightClassName, "attributes")
	return &RtctrlSetWeight{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetWeightCont, "dn"),
			Description:       G(RtctrlSetWeightCont, "descr"),
			Status:            G(RtctrlSetWeightCont, "status"),
			ClassName:         RtctrlsetweightClassName,
			Rn:                G(RtctrlSetWeightCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetWeightCont, "nameAlias"),
		},
		RtctrlSetWeightAttributes{
			Annotation: G(RtctrlSetWeightCont, "annotation"),
			Name:       G(RtctrlSetWeightCont, "name"),
			Type:       G(RtctrlSetWeightCont, "type"),
			Weight:     G(RtctrlSetWeightCont, "weight"),
		},
	}
}

func RtctrlSetWeightFromContainer(cont *container.Container) *RtctrlSetWeight {
	return RtctrlSetWeightFromContainerList(cont, 0)
}

func RtctrlSetWeightListFromContainer(cont *container.Container) []*RtctrlSetWeight {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetWeight, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetWeightFromContainerList(cont, i)
	}

	return arr
}
