package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetComm        = "uni/tn-%s/attr-%s/scomm"
	RnrtctrlSetComm        = "scomm"
	ParentDnrtctrlSetComm  = "uni/tn-%s/attr-%s"
	RtctrlsetcommClassName = "rtctrlSetComm"
)

type RtctrlSetComm struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetCommAttributes
}

type RtctrlSetCommAttributes struct {
	Annotation  string `json:",omitempty"`
	Community   string `json:",omitempty"`
	Name        string `json:",omitempty"`
	SetCriteria string `json:",omitempty"`
	Type        string `json:",omitempty"`
}

func NewRtctrlSetComm(rtctrlSetCommRn, parentDn, description, nameAlias string, rtctrlSetCommAttr RtctrlSetCommAttributes) *RtctrlSetComm {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetCommRn)
	return &RtctrlSetComm{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetcommClassName,
			Rn:                rtctrlSetCommRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetCommAttributes: rtctrlSetCommAttr,
	}
}

func (rtctrlSetComm *RtctrlSetComm) ToMap() (map[string]string, error) {
	rtctrlSetCommMap, err := rtctrlSetComm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetComm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetCommMap, key, value)
	}

	A(rtctrlSetCommMap, "annotation", rtctrlSetComm.Annotation)
	A(rtctrlSetCommMap, "community", rtctrlSetComm.Community)
	A(rtctrlSetCommMap, "name", rtctrlSetComm.Name)
	A(rtctrlSetCommMap, "setCriteria", rtctrlSetComm.SetCriteria)
	A(rtctrlSetCommMap, "type", rtctrlSetComm.Type)
	return rtctrlSetCommMap, err
}

func RtctrlSetCommFromContainerList(cont *container.Container, index int) *RtctrlSetComm {
	RtctrlSetCommCont := cont.S("imdata").Index(index).S(RtctrlsetcommClassName, "attributes")
	return &RtctrlSetComm{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetCommCont, "dn"),
			Description:       G(RtctrlSetCommCont, "descr"),
			Status:            G(RtctrlSetCommCont, "status"),
			ClassName:         RtctrlsetcommClassName,
			Rn:                G(RtctrlSetCommCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetCommCont, "nameAlias"),
		},
		RtctrlSetCommAttributes{
			Annotation:  G(RtctrlSetCommCont, "annotation"),
			Community:   G(RtctrlSetCommCont, "community"),
			Name:        G(RtctrlSetCommCont, "name"),
			SetCriteria: G(RtctrlSetCommCont, "setCriteria"),
			Type:        G(RtctrlSetCommCont, "type"),
		},
	}
}

func RtctrlSetCommFromContainer(cont *container.Container) *RtctrlSetComm {
	return RtctrlSetCommFromContainerList(cont, 0)
}

func RtctrlSetCommListFromContainer(cont *container.Container) []*RtctrlSetComm {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetComm, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetCommFromContainerList(cont, i)
	}

	return arr
}
