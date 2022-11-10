package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetPref        = "uni/tn-%s/attr-%s/spref"
	RnrtctrlSetPref        = "spref"
	ParentDnrtctrlSetPref  = "uni/tn-%s/attr-%s"
	RtctrlsetprefClassName = "rtctrlSetPref"
)

type RtctrlSetPref struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetPrefAttributes
}

type RtctrlSetPrefAttributes struct {
	Annotation string `json:",omitempty"`
	LocalPref  string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRtctrlSetPref(rtctrlSetPrefRn, parentDn, description, nameAlias string, rtctrlSetPrefAttr RtctrlSetPrefAttributes) *RtctrlSetPref {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetPrefRn)
	return &RtctrlSetPref{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetprefClassName,
			Rn:                rtctrlSetPrefRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetPrefAttributes: rtctrlSetPrefAttr,
	}
}

func (rtctrlSetPref *RtctrlSetPref) ToMap() (map[string]string, error) {
	rtctrlSetPrefMap, err := rtctrlSetPref.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetPref.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetPrefMap, key, value)
	}

	A(rtctrlSetPrefMap, "annotation", rtctrlSetPref.Annotation)
	A(rtctrlSetPrefMap, "localPref", rtctrlSetPref.LocalPref)
	A(rtctrlSetPrefMap, "name", rtctrlSetPref.Name)
	A(rtctrlSetPrefMap, "type", rtctrlSetPref.Type)
	return rtctrlSetPrefMap, err
}

func RtctrlSetPrefFromContainerList(cont *container.Container, index int) *RtctrlSetPref {
	RtctrlSetPrefCont := cont.S("imdata").Index(index).S(RtctrlsetprefClassName, "attributes")
	return &RtctrlSetPref{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetPrefCont, "dn"),
			Description:       G(RtctrlSetPrefCont, "descr"),
			Status:            G(RtctrlSetPrefCont, "status"),
			ClassName:         RtctrlsetprefClassName,
			Rn:                G(RtctrlSetPrefCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetPrefCont, "nameAlias"),
		},
		RtctrlSetPrefAttributes{
			Annotation: G(RtctrlSetPrefCont, "annotation"),
			LocalPref:  G(RtctrlSetPrefCont, "localPref"),
			Name:       G(RtctrlSetPrefCont, "name"),
			Type:       G(RtctrlSetPrefCont, "type"),
		},
	}
}

func RtctrlSetPrefFromContainer(cont *container.Container) *RtctrlSetPref {
	return RtctrlSetPrefFromContainerList(cont, 0)
}

func RtctrlSetPrefListFromContainer(cont *container.Container) []*RtctrlSetPref {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetPref, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetPrefFromContainerList(cont, i)
	}

	return arr
}
