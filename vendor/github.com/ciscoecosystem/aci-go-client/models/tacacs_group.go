package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DntacacsGroup        = "uni/fabric/tacacsgroup-%s"
	RntacacsGroup        = "tacacsgroup-%s"
	ParentDntacacsGroup  = "uni/fabric"
	TacacsgroupClassName = "tacacsGroup"
)

type TACACSMonitoringDestinationGroup struct {
	BaseAttributes
	NameAliasAttribute
	TACACSMonitoringDestinationGroupAttributes
}

type TACACSMonitoringDestinationGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewTACACSMonitoringDestinationGroup(tacacsGroupRn, parentDn, description, nameAlias string, tacacsGroupAttr TACACSMonitoringDestinationGroupAttributes) *TACACSMonitoringDestinationGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, tacacsGroupRn)
	return &TACACSMonitoringDestinationGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         TacacsgroupClassName,
			Rn:                tacacsGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TACACSMonitoringDestinationGroupAttributes: tacacsGroupAttr,
	}
}

func (tacacsGroup *TACACSMonitoringDestinationGroup) ToMap() (map[string]string, error) {
	tacacsGroupMap, err := tacacsGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := tacacsGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(tacacsGroupMap, key, value)
	}
	A(tacacsGroupMap, "annotation", tacacsGroup.Annotation)
	A(tacacsGroupMap, "name", tacacsGroup.Name)
	return tacacsGroupMap, err
}

func TACACSMonitoringDestinationGroupFromContainerList(cont *container.Container, index int) *TACACSMonitoringDestinationGroup {
	TACACSMonitoringDestinationGroupCont := cont.S("imdata").Index(index).S(TacacsgroupClassName, "attributes")
	return &TACACSMonitoringDestinationGroup{
		BaseAttributes{
			DistinguishedName: G(TACACSMonitoringDestinationGroupCont, "dn"),
			Description:       G(TACACSMonitoringDestinationGroupCont, "descr"),
			Status:            G(TACACSMonitoringDestinationGroupCont, "status"),
			ClassName:         TacacsgroupClassName,
			Rn:                G(TACACSMonitoringDestinationGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TACACSMonitoringDestinationGroupCont, "nameAlias"),
		},
		TACACSMonitoringDestinationGroupAttributes{
			Annotation: G(TACACSMonitoringDestinationGroupCont, "annotation"),
			Name:       G(TACACSMonitoringDestinationGroupCont, "name"),
		},
	}
}

func TACACSMonitoringDestinationGroupFromContainer(cont *container.Container) *TACACSMonitoringDestinationGroup {
	return TACACSMonitoringDestinationGroupFromContainerList(cont, 0)
}

func TACACSMonitoringDestinationGroupListFromContainer(cont *container.Container) []*TACACSMonitoringDestinationGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TACACSMonitoringDestinationGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = TACACSMonitoringDestinationGroupFromContainerList(cont, i)
	}
	return arr
}
