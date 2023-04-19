package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetAddComm        = "uni/tn-%s/attr-%s/saddcomm-%s"
	RnrtctrlSetAddComm        = "saddcomm-%s"
	ParentDnrtctrlSetAddComm  = "uni/tn-%s/attr-%s"
	RtctrlsetaddcommClassName = "rtctrlSetAddComm"
)

type RtctrlSetAddComm struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetAddCommAttributes
}

type RtctrlSetAddCommAttributes struct {
	Annotation  string `json:",omitempty"`
	Community   string `json:",omitempty"`
	Name        string `json:",omitempty"`
	SetCriteria string `json:",omitempty"`
	Type        string `json:",omitempty"`
}

func NewRtctrlSetAddComm(rtctrlSetAddCommRn, parentDn, description, nameAlias string, rtctrlSetAddCommAttr RtctrlSetAddCommAttributes) *RtctrlSetAddComm {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetAddCommRn)
	return &RtctrlSetAddComm{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetaddcommClassName,
			Rn:                rtctrlSetAddCommRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetAddCommAttributes: rtctrlSetAddCommAttr,
	}
}

func (rtctrlSetAddComm *RtctrlSetAddComm) ToMap() (map[string]string, error) {
	rtctrlSetAddCommMap, err := rtctrlSetAddComm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetAddComm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetAddCommMap, key, value)
	}

	A(rtctrlSetAddCommMap, "annotation", rtctrlSetAddComm.Annotation)
	A(rtctrlSetAddCommMap, "community", rtctrlSetAddComm.Community)
	A(rtctrlSetAddCommMap, "name", rtctrlSetAddComm.Name)
	A(rtctrlSetAddCommMap, "setCriteria", rtctrlSetAddComm.SetCriteria)
	A(rtctrlSetAddCommMap, "type", rtctrlSetAddComm.Type)
	return rtctrlSetAddCommMap, err
}

func RtctrlSetAddCommFromContainerList(cont *container.Container, index int) *RtctrlSetAddComm {
	RtctrlSetAddCommCont := cont.S("imdata").Index(index).S(RtctrlsetaddcommClassName, "attributes")
	return &RtctrlSetAddComm{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetAddCommCont, "dn"),
			Description:       G(RtctrlSetAddCommCont, "descr"),
			Status:            G(RtctrlSetAddCommCont, "status"),
			ClassName:         RtctrlsetaddcommClassName,
			Rn:                G(RtctrlSetAddCommCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetAddCommCont, "nameAlias"),
		},
		RtctrlSetAddCommAttributes{
			Annotation:  G(RtctrlSetAddCommCont, "annotation"),
			Community:   G(RtctrlSetAddCommCont, "community"),
			Name:        G(RtctrlSetAddCommCont, "name"),
			SetCriteria: G(RtctrlSetAddCommCont, "setCriteria"),
			Type:        G(RtctrlSetAddCommCont, "type"),
		},
	}
}

func RtctrlSetAddCommFromContainer(cont *container.Container) *RtctrlSetAddComm {
	return RtctrlSetAddCommFromContainerList(cont, 0)
}

func RtctrlSetAddCommListFromContainer(cont *container.Container) []*RtctrlSetAddComm {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetAddComm, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetAddCommFromContainerList(cont, i)
	}

	return arr
}
