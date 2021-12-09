package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlScope        = "uni/tn-%s/prof-%s/ctx-%s/scp"
	RnrtctrlScope        = "scp"
	ParentDnrtctrlScope  = "uni/tn-%s/prof-%s/ctx-%s"
	RtctrlscopeClassName = "rtctrlScope"
)

type RouteContextScope struct {
	BaseAttributes
	NameAliasAttribute
	RouteContextScopeAttributes
}

type RouteContextScopeAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewRouteContextScope(rtctrlScopeRn, parentDn, description, nameAlias string, rtctrlScopeAttr RouteContextScopeAttributes) *RouteContextScope {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlScopeRn)
	return &RouteContextScope{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlscopeClassName,
			Rn:                rtctrlScopeRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RouteContextScopeAttributes: rtctrlScopeAttr,
	}
}

func (rtctrlScope *RouteContextScope) ToMap() (map[string]string, error) {
	rtctrlScopeMap, err := rtctrlScope.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := rtctrlScope.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(rtctrlScopeMap, key, value)
	}

	A(rtctrlScopeMap, "annotation", rtctrlScope.Annotation)
	A(rtctrlScopeMap, "name", rtctrlScope.Name)
	return rtctrlScopeMap, err
}

func RouteContextScopeFromContainerList(cont *container.Container, index int) *RouteContextScope {
	RouteContextScopeCont := cont.S("imdata").Index(index).S(RtctrlscopeClassName, "attributes")
	return &RouteContextScope{
		BaseAttributes{
			DistinguishedName: G(RouteContextScopeCont, "dn"),
			Description:       G(RouteContextScopeCont, "descr"),
			Status:            G(RouteContextScopeCont, "status"),
			ClassName:         RtctrlscopeClassName,
			Rn:                G(RouteContextScopeCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RouteContextScopeCont, "nameAlias"),
		},
		RouteContextScopeAttributes{

			Annotation: G(RouteContextScopeCont, "annotation"),
			Name:       G(RouteContextScopeCont, "name"),
		},
	}
}

func RouteContextScopeFromContainer(cont *container.Container) *RouteContextScope {
	return RouteContextScopeFromContainerList(cont, 0)
}

func RouteContextScopeListFromContainer(cont *container.Container) []*RouteContextScope {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RouteContextScope, length)
	for i := 0; i < length; i++ {
		arr[i] = RouteContextScopeFromContainerList(cont, i)
	}
	return arr
}
