package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlSetNh        = "uni/tn-%s/attr-%s/nh"
	RnrtctrlSetNh        = "nh"
	ParentDnrtctrlSetNh  = "uni/tn-%s/attr-%s"
	RtctrlsetnhClassName = "rtctrlSetNh"
)

type RtctrlSetNh struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetNhAttributes
}

type RtctrlSetNhAttributes struct {
	Annotation string `json:",omitempty"`
	Addr       string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRtctrlSetNh(rtctrlSetNhRn, parentDn, description, nameAlias string, rtctrlSetNhAttr RtctrlSetNhAttributes) *RtctrlSetNh {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetNhRn)
	return &RtctrlSetNh{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetnhClassName,
			Rn:                rtctrlSetNhRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetNhAttributes: rtctrlSetNhAttr,
	}
}

func (rtctrlSetNh *RtctrlSetNh) ToMap() (map[string]string, error) {
	rtctrlSetNhMap, err := rtctrlSetNh.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetNh.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetNhMap, key, value)
	}

	A(rtctrlSetNhMap, "addr", rtctrlSetNh.Addr)
	A(rtctrlSetNhMap, "annotation", rtctrlSetNh.Annotation)
	A(rtctrlSetNhMap, "name", rtctrlSetNh.Name)
	A(rtctrlSetNhMap, "type", rtctrlSetNh.Type)
	return rtctrlSetNhMap, err
}

func RtctrlSetNhFromContainerList(cont *container.Container, index int) *RtctrlSetNh {
	RtctrlSetNhCont := cont.S("imdata").Index(index).S(RtctrlsetnhClassName, "attributes")
	return &RtctrlSetNh{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetNhCont, "dn"),
			Description:       G(RtctrlSetNhCont, "descr"),
			Status:            G(RtctrlSetNhCont, "status"),
			ClassName:         RtctrlsetnhClassName,
			Rn:                G(RtctrlSetNhCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetNhCont, "nameAlias"),
		},
		RtctrlSetNhAttributes{
			Addr:       G(RtctrlSetNhCont, "addr"),
			Annotation: G(RtctrlSetNhCont, "annotation"),
			Name:       G(RtctrlSetNhCont, "name"),
			Type:       G(RtctrlSetNhCont, "type"),
		},
	}
}

func RtctrlSetNhFromContainer(cont *container.Container) *RtctrlSetNh {
	return RtctrlSetNhFromContainerList(cont, 0)
}

func RtctrlSetNhListFromContainer(cont *container.Container) []*RtctrlSetNh {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetNh, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetNhFromContainerList(cont, i)
	}

	return arr
}
