package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetRedistMultipath        = "uni/tn-%s/attr-%s/redistmpath"
	RnrtctrlSetRedistMultipath        = "redistmpath"
	ParentDnrtctrlSetRedistMultipath  = "uni/tn-%s/attr-%s"
	RtctrlsetredistmultipathClassName = "rtctrlSetRedistMultipath"
)

type RedistributeMultipathAction struct {
	BaseAttributes
	NameAliasAttribute
	RedistributeMultipathActionAttributes
}

type RedistributeMultipathActionAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRedistributeMultipathAction(rtctrlSetRedistMultipathRn, parentDn, description, nameAlias string, rtctrlSetRedistMultipathAttr RedistributeMultipathActionAttributes) *RedistributeMultipathAction {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetRedistMultipathRn)
	return &RedistributeMultipathAction{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetredistmultipathClassName,
			Rn:                rtctrlSetRedistMultipathRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RedistributeMultipathActionAttributes: rtctrlSetRedistMultipathAttr,
	}
}

func (rtctrlSetRedistMultipath *RedistributeMultipathAction) ToMap() (map[string]string, error) {
	rtctrlSetRedistMultipathMap, err := rtctrlSetRedistMultipath.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetRedistMultipath.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetRedistMultipathMap, key, value)
	}

	A(rtctrlSetRedistMultipathMap, "annotation", rtctrlSetRedistMultipath.Annotation)
	A(rtctrlSetRedistMultipathMap, "name", rtctrlSetRedistMultipath.Name)
	A(rtctrlSetRedistMultipathMap, "type", rtctrlSetRedistMultipath.Type)
	return rtctrlSetRedistMultipathMap, err
}

func RedistributeMultipathActionFromContainerList(cont *container.Container, index int) *RedistributeMultipathAction {
	RedistributeMultipathActionCont := cont.S("imdata").Index(index).S(RtctrlsetredistmultipathClassName, "attributes")
	return &RedistributeMultipathAction{
		BaseAttributes{
			DistinguishedName: G(RedistributeMultipathActionCont, "dn"),
			Description:       G(RedistributeMultipathActionCont, "descr"),
			Status:            G(RedistributeMultipathActionCont, "status"),
			ClassName:         RtctrlsetredistmultipathClassName,
			Rn:                G(RedistributeMultipathActionCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RedistributeMultipathActionCont, "nameAlias"),
		},
		RedistributeMultipathActionAttributes{
			Annotation: G(RedistributeMultipathActionCont, "annotation"),
			Name:       G(RedistributeMultipathActionCont, "name"),
			Type:       G(RedistributeMultipathActionCont, "type"),
		},
	}
}

func RedistributeMultipathActionFromContainer(cont *container.Container) *RedistributeMultipathAction {
	return RedistributeMultipathActionFromContainerList(cont, 0)
}

func RedistributeMultipathActionListFromContainer(cont *container.Container) []*RedistributeMultipathAction {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RedistributeMultipathAction, length)

	for i := 0; i < length; i++ {
		arr[i] = RedistributeMultipathActionFromContainerList(cont, i)
	}

	return arr
}
