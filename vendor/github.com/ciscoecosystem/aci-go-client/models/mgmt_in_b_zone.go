package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnmgmtInBZone        = "uni/infra/funcprof/grp-%s/inbzone"
	RnmgmtInBZone        = "inbzone"
	ParentDnmgmtInBZone  = "uni/infra/funcprof/grp-%s"
	MgmtinbzoneClassName = "mgmtInBZone"
)

type InBManagedNodesZone struct {
	BaseAttributes
	NameAliasAttribute
	InBManagedNodesZoneAttributes
}

type InBManagedNodesZoneAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewInBManagedNodesZone(mgmtInBZoneRn, parentDn, description, nameAlias string, mgmtInBZoneAttr InBManagedNodesZoneAttributes) *InBManagedNodesZone {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtInBZoneRn)
	return &InBManagedNodesZone{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtinbzoneClassName,
			Rn:                mgmtInBZoneRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		InBManagedNodesZoneAttributes: mgmtInBZoneAttr,
	}
}

func (mgmtInBZone *InBManagedNodesZone) ToMap() (map[string]string, error) {
	mgmtInBZoneMap, err := mgmtInBZone.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := mgmtInBZone.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(mgmtInBZoneMap, key, value)
	}
	A(mgmtInBZoneMap, "annotation", mgmtInBZone.Annotation)
	A(mgmtInBZoneMap, "name", mgmtInBZone.Name)
	return mgmtInBZoneMap, err
}

func InBManagedNodesZoneFromContainerList(cont *container.Container, index int) *InBManagedNodesZone {
	InBManagedNodesZoneCont := cont.S("imdata").Index(index).S(MgmtinbzoneClassName, "attributes")
	return &InBManagedNodesZone{
		BaseAttributes{
			DistinguishedName: G(InBManagedNodesZoneCont, "dn"),
			Description:       G(InBManagedNodesZoneCont, "descr"),
			Status:            G(InBManagedNodesZoneCont, "status"),
			ClassName:         MgmtinbzoneClassName,
			Rn:                G(InBManagedNodesZoneCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(InBManagedNodesZoneCont, "nameAlias"),
		},
		InBManagedNodesZoneAttributes{
			Annotation: G(InBManagedNodesZoneCont, "annotation"),
			Name:       G(InBManagedNodesZoneCont, "name"),
		},
	}
}

func InBManagedNodesZoneFromContainer(cont *container.Container) *InBManagedNodesZone {
	return InBManagedNodesZoneFromContainerList(cont, 0)
}

func InBManagedNodesZoneListFromContainer(cont *container.Container) []*InBManagedNodesZone {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*InBManagedNodesZone, length)
	for i := 0; i < length; i++ {
		arr[i] = InBManagedNodesZoneFromContainerList(cont, i)
	}
	return arr
}
