package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlCtxP          = "uni/tn-%s/prof-%s/ctx-%s"
	DnrtctrlCtxPOut       = "uni/tn-%s/out-%s/prof-%s/ctx-%s"
	RnrtctrlCtxP          = "ctx-%s"
	ParentDnrtctrlCtxP    = "uni/tn-%s/prof-%s"
	ParentDnrtctrlCtxPOut = "uni/tn-%s/out-%s/prof-%s"
	RtctrlctxpClassName   = "rtctrlCtxP"
)

type RouteControlContext struct {
	BaseAttributes
	NameAliasAttribute
	RouteControlContextAttributes
}

type RouteControlContextAttributes struct {
	Action     string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Order      string `json:",omitempty"`
}

func NewRouteControlContext(rtctrlCtxPRn, parentDn, description, nameAlias string, rtctrlCtxPAttr RouteControlContextAttributes) *RouteControlContext {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlCtxPRn)
	return &RouteControlContext{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlctxpClassName,
			Rn:                rtctrlCtxPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RouteControlContextAttributes: rtctrlCtxPAttr,
	}
}

func (rtctrlCtxP *RouteControlContext) ToMap() (map[string]string, error) {
	rtctrlCtxPMap, err := rtctrlCtxP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := rtctrlCtxP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(rtctrlCtxPMap, key, value)
	}
	A(rtctrlCtxPMap, "action", rtctrlCtxP.Action)
	A(rtctrlCtxPMap, "annotation", rtctrlCtxP.Annotation)
	A(rtctrlCtxPMap, "name", rtctrlCtxP.Name)
	A(rtctrlCtxPMap, "order", rtctrlCtxP.Order)
	return rtctrlCtxPMap, err
}

func RouteControlContextFromContainerList(cont *container.Container, index int) *RouteControlContext {
	RouteControlContextCont := cont.S("imdata").Index(index).S(RtctrlctxpClassName, "attributes")
	return &RouteControlContext{
		BaseAttributes{
			DistinguishedName: G(RouteControlContextCont, "dn"),
			Description:       G(RouteControlContextCont, "descr"),
			Status:            G(RouteControlContextCont, "status"),
			ClassName:         RtctrlctxpClassName,
			Rn:                G(RouteControlContextCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RouteControlContextCont, "nameAlias"),
		},
		RouteControlContextAttributes{
			Action:     G(RouteControlContextCont, "action"),
			Annotation: G(RouteControlContextCont, "annotation"),
			Name:       G(RouteControlContextCont, "name"),
			Order:      G(RouteControlContextCont, "order"),
		},
	}
}

func RouteControlContextFromContainer(cont *container.Container) *RouteControlContext {
	return RouteControlContextFromContainerList(cont, 0)
}

func RouteControlContextListFromContainer(cont *container.Container) []*RouteControlContext {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RouteControlContext, length)
	for i := 0; i < length; i++ {
		arr[i] = RouteControlContextFromContainerList(cont, i)
	}
	return arr
}
