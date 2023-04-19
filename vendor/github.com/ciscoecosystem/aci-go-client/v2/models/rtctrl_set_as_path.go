package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetASPath        = "uni/tn-%s/attr-%s/saspath-%s"
	RnrtctrlSetASPath        = "saspath-%s"
	ParentDnrtctrlSetASPath  = "uni/tn-%s/attr-%s"
	RtctrlsetaspathClassName = "rtctrlSetASPath"
)

type SetASPath struct {
	BaseAttributes
	NameAliasAttribute
	SetASPathAttributes
}

type SetASPathAttributes struct {
	Annotation string `json:",omitempty"`
	Criteria   string `json:",omitempty"`
	Lastnum    string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewSetASPath(rtctrlSetASPathRn, parentDn, description, nameAlias string, rtctrlSetASPathAttr SetASPathAttributes) *SetASPath {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetASPathRn)
	return &SetASPath{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetaspathClassName,
			Rn:                rtctrlSetASPathRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SetASPathAttributes: rtctrlSetASPathAttr,
	}
}

func (rtctrlSetASPath *SetASPath) ToMap() (map[string]string, error) {
	rtctrlSetASPathMap, err := rtctrlSetASPath.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetASPath.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetASPathMap, key, value)
	}

	A(rtctrlSetASPathMap, "annotation", rtctrlSetASPath.Annotation)
	A(rtctrlSetASPathMap, "criteria", rtctrlSetASPath.Criteria)
	A(rtctrlSetASPathMap, "lastnum", rtctrlSetASPath.Lastnum)
	A(rtctrlSetASPathMap, "name", rtctrlSetASPath.Name)
	A(rtctrlSetASPathMap, "type", rtctrlSetASPath.Type)
	return rtctrlSetASPathMap, err
}

func SetASPathFromContainerList(cont *container.Container, index int) *SetASPath {
	SetASPathCont := cont.S("imdata").Index(index).S(RtctrlsetaspathClassName, "attributes")
	return &SetASPath{
		BaseAttributes{
			DistinguishedName: G(SetASPathCont, "dn"),
			Description:       G(SetASPathCont, "descr"),
			Status:            G(SetASPathCont, "status"),
			ClassName:         RtctrlsetaspathClassName,
			Rn:                G(SetASPathCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SetASPathCont, "nameAlias"),
		},
		SetASPathAttributes{
			Annotation: G(SetASPathCont, "annotation"),
			Criteria:   G(SetASPathCont, "criteria"),
			Lastnum:    G(SetASPathCont, "lastnum"),
			Name:       G(SetASPathCont, "name"),
			Type:       G(SetASPathCont, "type"),
		},
	}
}

func SetASPathFromContainer(cont *container.Container) *SetASPath {
	return SetASPathFromContainerList(cont, 0)
}

func SetASPathListFromContainer(cont *container.Container) []*SetASPath {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SetASPath, length)

	for i := 0; i < length; i++ {
		arr[i] = SetASPathFromContainerList(cont, i)
	}

	return arr
}
