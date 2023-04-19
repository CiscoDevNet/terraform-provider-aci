package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnleakTo        = "to-[%s]-[%s]"
	LeaktoClassName = "leakTo"
)

type TenantandVRFdestinationforInterVRFLeakedRoutes struct {
	BaseAttributes
	NameAliasAttribute
	TenantandVRFdestinationforInterVRFLeakedRoutesAttributes
}

type TenantandVRFdestinationforInterVRFLeakedRoutesAttributes struct {
	Annotation            string `json:",omitempty"`
	DestinationCtxName    string `json:",omitempty"`
	Name                  string `json:",omitempty"`
	Scope                 string `json:",omitempty"`
	DestinationTenantName string `json:",omitempty"`
	ToCtxDn               string `json:",omitempty"`
}

func NewTenantandVRFdestinationforInterVRFLeakedRoutes(leakToRn, parentDn, description, nameAlias string, leakToAttr TenantandVRFdestinationforInterVRFLeakedRoutesAttributes) *TenantandVRFdestinationforInterVRFLeakedRoutes {
	dn := fmt.Sprintf("%s/%s", parentDn, leakToRn)
	return &TenantandVRFdestinationforInterVRFLeakedRoutes{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LeaktoClassName,
			Rn:                leakToRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TenantandVRFdestinationforInterVRFLeakedRoutesAttributes: leakToAttr,
	}
}

func (leakTo *TenantandVRFdestinationforInterVRFLeakedRoutes) ToMap() (map[string]string, error) {
	leakToMap, err := leakTo.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := leakTo.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(leakToMap, key, value)
	}

	A(leakToMap, "ctxName", leakTo.DestinationCtxName)
	A(leakToMap, "toCtxDn", leakTo.ToCtxDn)
	A(leakToMap, "name", leakTo.Name)
	A(leakToMap, "scope", leakTo.Scope)
	A(leakToMap, "tenantName", leakTo.DestinationTenantName)
	return leakToMap, err
}

func TenantandVRFdestinationforInterVRFLeakedRoutesFromContainerList(cont *container.Container, index int) *TenantandVRFdestinationforInterVRFLeakedRoutes {
	TenantandVRFdestinationforInterVRFLeakedRoutesCont := cont.S("imdata").Index(index).S(LeaktoClassName, "attributes")
	return &TenantandVRFdestinationforInterVRFLeakedRoutes{
		BaseAttributes{
			DistinguishedName: G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "dn"),
			Description:       G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "descr"),
			Status:            G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "status"),
			ClassName:         LeaktoClassName,
			Rn:                G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "nameAlias"),
		},
		TenantandVRFdestinationforInterVRFLeakedRoutesAttributes{
			DestinationCtxName:    G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "ctxName"),
			ToCtxDn:               G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "toCtxDn"),
			Name:                  G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "name"),
			Scope:                 G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "scope"),
			DestinationTenantName: G(TenantandVRFdestinationforInterVRFLeakedRoutesCont, "tenantName"),
		},
	}
}

func TenantandVRFdestinationforInterVRFLeakedRoutesFromContainer(cont *container.Container) *TenantandVRFdestinationforInterVRFLeakedRoutes {
	return TenantandVRFdestinationforInterVRFLeakedRoutesFromContainerList(cont, 0)
}

func TenantandVRFdestinationforInterVRFLeakedRoutesListFromContainer(cont *container.Container) []*TenantandVRFdestinationforInterVRFLeakedRoutes {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TenantandVRFdestinationforInterVRFLeakedRoutes, length)

	for i := 0; i < length; i++ {
		arr[i] = TenantandVRFdestinationforInterVRFLeakedRoutesFromContainerList(cont, i)
	}

	return arr
}
