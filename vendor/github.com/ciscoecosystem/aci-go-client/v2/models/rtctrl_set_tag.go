package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetTag        = "uni/tn-%s/attr-%s/srttag"
	RnrtctrlSetTag        = "srttag"
	ParentDnrtctrlSetTag  = "uni/tn-%s/attr-%s"
	RtctrlsettagClassName = "rtctrlSetTag"
)

type RtctrlSetTag struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetTagAttributes
}

type RtctrlSetTagAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Tag        string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRtctrlSetTag(rtctrlSetTagRn, parentDn, description, nameAlias string, rtctrlSetTagAttr RtctrlSetTagAttributes) *RtctrlSetTag {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetTagRn)
	return &RtctrlSetTag{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsettagClassName,
			Rn:                rtctrlSetTagRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetTagAttributes: rtctrlSetTagAttr,
	}
}

func (rtctrlSetTag *RtctrlSetTag) ToMap() (map[string]string, error) {
	rtctrlSetTagMap, err := rtctrlSetTag.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetTag.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetTagMap, key, value)
	}

	A(rtctrlSetTagMap, "annotation", rtctrlSetTag.Annotation)
	A(rtctrlSetTagMap, "name", rtctrlSetTag.Name)
	A(rtctrlSetTagMap, "tag", rtctrlSetTag.Tag)
	A(rtctrlSetTagMap, "type", rtctrlSetTag.Type)
	return rtctrlSetTagMap, err
}

func RtctrlSetTagFromContainerList(cont *container.Container, index int) *RtctrlSetTag {
	RtctrlSetTagCont := cont.S("imdata").Index(index).S(RtctrlsettagClassName, "attributes")
	return &RtctrlSetTag{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetTagCont, "dn"),
			Description:       G(RtctrlSetTagCont, "descr"),
			Status:            G(RtctrlSetTagCont, "status"),
			ClassName:         RtctrlsettagClassName,
			Rn:                G(RtctrlSetTagCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetTagCont, "nameAlias"),
		},
		RtctrlSetTagAttributes{
			Annotation: G(RtctrlSetTagCont, "annotation"),
			Name:       G(RtctrlSetTagCont, "name"),
			Tag:        G(RtctrlSetTagCont, "tag"),
			Type:       G(RtctrlSetTagCont, "type"),
		},
	}
}

func RtctrlSetTagFromContainer(cont *container.Container) *RtctrlSetTag {
	return RtctrlSetTagFromContainerList(cont, 0)
}

func RtctrlSetTagListFromContainer(cont *container.Container) []*RtctrlSetTag {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetTag, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetTagFromContainerList(cont, i)
	}

	return arr
}
