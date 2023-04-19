package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetNhUnchanged        = "uni/tn-%s/attr-%s/nhunchanged"
	RnrtctrlSetNhUnchanged        = "nhunchanged"
	ParentDnrtctrlSetNhUnchanged  = "uni/tn-%s/attr-%s"
	RtctrlsetnhunchangedClassName = "rtctrlSetNhUnchanged"
)

type NexthopUnchangedAction struct {
	BaseAttributes
	NameAliasAttribute
	NexthopUnchangedActionAttributes
}

type NexthopUnchangedActionAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewNexthopUnchangedAction(rtctrlSetNhUnchangedRn, parentDn, description, nameAlias string, rtctrlSetNhUnchangedAttr NexthopUnchangedActionAttributes) *NexthopUnchangedAction {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetNhUnchangedRn)
	return &NexthopUnchangedAction{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetnhunchangedClassName,
			Rn:                rtctrlSetNhUnchangedRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		NexthopUnchangedActionAttributes: rtctrlSetNhUnchangedAttr,
	}
}

func (rtctrlSetNhUnchanged *NexthopUnchangedAction) ToMap() (map[string]string, error) {
	rtctrlSetNhUnchangedMap, err := rtctrlSetNhUnchanged.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetNhUnchanged.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetNhUnchangedMap, key, value)
	}

	A(rtctrlSetNhUnchangedMap, "annotation", rtctrlSetNhUnchanged.Annotation)
	A(rtctrlSetNhUnchangedMap, "name", rtctrlSetNhUnchanged.Name)
	A(rtctrlSetNhUnchangedMap, "type", rtctrlSetNhUnchanged.Type)
	return rtctrlSetNhUnchangedMap, err
}

func NexthopUnchangedActionFromContainerList(cont *container.Container, index int) *NexthopUnchangedAction {
	NexthopUnchangedActionCont := cont.S("imdata").Index(index).S(RtctrlsetnhunchangedClassName, "attributes")
	return &NexthopUnchangedAction{
		BaseAttributes{
			DistinguishedName: G(NexthopUnchangedActionCont, "dn"),
			Description:       G(NexthopUnchangedActionCont, "descr"),
			Status:            G(NexthopUnchangedActionCont, "status"),
			ClassName:         RtctrlsetnhunchangedClassName,
			Rn:                G(NexthopUnchangedActionCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(NexthopUnchangedActionCont, "nameAlias"),
		},
		NexthopUnchangedActionAttributes{
			Annotation: G(NexthopUnchangedActionCont, "annotation"),
			Name:       G(NexthopUnchangedActionCont, "name"),
			Type:       G(NexthopUnchangedActionCont, "type"),
		},
	}
}

func NexthopUnchangedActionFromContainer(cont *container.Container) *NexthopUnchangedAction {
	return NexthopUnchangedActionFromContainerList(cont, 0)
}

func NexthopUnchangedActionListFromContainer(cont *container.Container) []*NexthopUnchangedAction {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*NexthopUnchangedAction, length)

	for i := 0; i < length; i++ {
		arr[i] = NexthopUnchangedActionFromContainerList(cont, i)
	}

	return arr
}
