package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnleakRoutes        = "uni/tn-%s/ctx-%s/leakroutes"
	RnleakRoutes        = "leakroutes"
	ParentDnleakRoutes  = "uni/tn-%s/ctx-%s"
	LeakroutesClassName = "leakRoutes"
)

type InterVRFLeakedRoutesContainer struct {
	BaseAttributes
	NameAliasAttribute
	InterVRFLeakedRoutesContainerAttributes
}

type InterVRFLeakedRoutesContainerAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewInterVRFLeakedRoutesContainer(leakRoutesRn, parentDn, description, nameAlias string, leakRoutesAttr InterVRFLeakedRoutesContainerAttributes) *InterVRFLeakedRoutesContainer {
	dn := fmt.Sprintf("%s/%s", parentDn, leakRoutesRn)
	return &InterVRFLeakedRoutesContainer{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LeakroutesClassName,
			Rn:                leakRoutesRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		InterVRFLeakedRoutesContainerAttributes: leakRoutesAttr,
	}
}

func (leakRoutes *InterVRFLeakedRoutesContainer) ToMap() (map[string]string, error) {
	leakRoutesMap, err := leakRoutes.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := leakRoutes.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(leakRoutesMap, key, value)
	}

	A(leakRoutesMap, "annotation", leakRoutes.Annotation)
	A(leakRoutesMap, "name", leakRoutes.Name)
	return leakRoutesMap, err
}

func InterVRFLeakedRoutesContainerFromContainerList(cont *container.Container, index int) *InterVRFLeakedRoutesContainer {
	InterVRFLeakedRoutesContainerCont := cont.S("imdata").Index(index).S(LeakroutesClassName, "attributes")
	return &InterVRFLeakedRoutesContainer{
		BaseAttributes{
			DistinguishedName: G(InterVRFLeakedRoutesContainerCont, "dn"),
			Description:       G(InterVRFLeakedRoutesContainerCont, "descr"),
			Status:            G(InterVRFLeakedRoutesContainerCont, "status"),
			ClassName:         LeakroutesClassName,
			Rn:                G(InterVRFLeakedRoutesContainerCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(InterVRFLeakedRoutesContainerCont, "nameAlias"),
		},
		InterVRFLeakedRoutesContainerAttributes{
			Annotation: G(InterVRFLeakedRoutesContainerCont, "annotation"),
			Name:       G(InterVRFLeakedRoutesContainerCont, "name"),
		},
	}
}

func InterVRFLeakedRoutesContainerFromContainer(cont *container.Container) *InterVRFLeakedRoutesContainer {
	return InterVRFLeakedRoutesContainerFromContainerList(cont, 0)
}

func InterVRFLeakedRoutesContainerListFromContainer(cont *container.Container) []*InterVRFLeakedRoutesContainer {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*InterVRFLeakedRoutesContainer, length)

	for i := 0; i < length; i++ {
		arr[i] = InterVRFLeakedRoutesContainerFromContainerList(cont, i)
	}

	return arr
}
