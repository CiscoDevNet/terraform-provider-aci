package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlSetDamp        = "uni/tn-%s/attr-%s/sdamp"
	RnrtctrlSetDamp        = "sdamp"
	ParentDnrtctrlSetDamp  = "uni/tn-%s/attr-%s"
	RtctrlsetdampClassName = "rtctrlSetDamp"
)

type RtctrlSetDamp struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetDampAttributes
}

type RtctrlSetDampAttributes struct {
	Annotation      string `json:",omitempty"`
	HalfLife        string `json:",omitempty"`
	MaxSuppressTime string `json:",omitempty"`
	Name            string `json:",omitempty"`
	Reuse           string `json:",omitempty"`
	Suppress        string `json:",omitempty"`
	Type            string `json:",omitempty"`
}

func NewRtctrlSetDamp(rtctrlSetDampRn, parentDn, description, nameAlias string, rtctrlSetDampAttr RtctrlSetDampAttributes) *RtctrlSetDamp {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetDampRn)
	return &RtctrlSetDamp{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetdampClassName,
			Rn:                rtctrlSetDampRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetDampAttributes: rtctrlSetDampAttr,
	}
}

func (rtctrlSetDamp *RtctrlSetDamp) ToMap() (map[string]string, error) {
	rtctrlSetDampMap, err := rtctrlSetDamp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetDamp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetDampMap, key, value)
	}

	A(rtctrlSetDampMap, "annotation", rtctrlSetDamp.Annotation)
	A(rtctrlSetDampMap, "halfLife", rtctrlSetDamp.HalfLife)
	A(rtctrlSetDampMap, "maxSuppressTime", rtctrlSetDamp.MaxSuppressTime)
	A(rtctrlSetDampMap, "name", rtctrlSetDamp.Name)
	A(rtctrlSetDampMap, "reuse", rtctrlSetDamp.Reuse)
	A(rtctrlSetDampMap, "suppress", rtctrlSetDamp.Suppress)
	A(rtctrlSetDampMap, "type", rtctrlSetDamp.Type)
	return rtctrlSetDampMap, err
}

func RtctrlSetDampFromContainerList(cont *container.Container, index int) *RtctrlSetDamp {
	RtctrlSetDampCont := cont.S("imdata").Index(index).S(RtctrlsetdampClassName, "attributes")
	return &RtctrlSetDamp{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetDampCont, "dn"),
			Description:       G(RtctrlSetDampCont, "descr"),
			Status:            G(RtctrlSetDampCont, "status"),
			ClassName:         RtctrlsetdampClassName,
			Rn:                G(RtctrlSetDampCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetDampCont, "nameAlias"),
		},
		RtctrlSetDampAttributes{
			Annotation:      G(RtctrlSetDampCont, "annotation"),
			HalfLife:        G(RtctrlSetDampCont, "halfLife"),
			MaxSuppressTime: G(RtctrlSetDampCont, "maxSuppressTime"),
			Name:            G(RtctrlSetDampCont, "name"),
			Reuse:           G(RtctrlSetDampCont, "reuse"),
			Suppress:        G(RtctrlSetDampCont, "suppress"),
			Type:            G(RtctrlSetDampCont, "type"),
		},
	}
}

func RtctrlSetDampFromContainer(cont *container.Container) *RtctrlSetDamp {
	return RtctrlSetDampFromContainerList(cont, 0)
}

func RtctrlSetDampListFromContainer(cont *container.Container) []*RtctrlSetDamp {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetDamp, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetDampFromContainerList(cont, i)
	}

	return arr
}
