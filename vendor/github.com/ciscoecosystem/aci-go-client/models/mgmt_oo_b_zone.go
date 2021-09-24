package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnmgmtOoBZone        = "uni/infra/funcprof/grp-%s/oobzone"
	RnmgmtOoBZone        = "oobzone"
	ParentDnmgmtOoBZone  = "uni/infra/funcprof/grp-%s"
	MgmtoobzoneClassName = "mgmtOoBZone"
)

type OOBManagedNodesZone struct {
	BaseAttributes
	NameAliasAttribute
	OOBManagedNodesZoneAttributes
}

type OOBManagedNodesZoneAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewOOBManagedNodesZone(mgmtOoBZoneRn, parentDn, description, nameAlias string, mgmtOoBZoneAttr OOBManagedNodesZoneAttributes) *OOBManagedNodesZone {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtOoBZoneRn)
	return &OOBManagedNodesZone{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtoobzoneClassName,
			Rn:                mgmtOoBZoneRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		OOBManagedNodesZoneAttributes: mgmtOoBZoneAttr,
	}
}

func (mgmtOoBZone *OOBManagedNodesZone) ToMap() (map[string]string, error) {
	mgmtOoBZoneMap, err := mgmtOoBZone.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := mgmtOoBZone.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(mgmtOoBZoneMap, key, value)
	}
	A(mgmtOoBZoneMap, "annotation", mgmtOoBZone.Annotation)
	A(mgmtOoBZoneMap, "name", mgmtOoBZone.Name)
	return mgmtOoBZoneMap, err
}

func OOBManagedNodesZoneFromContainerList(cont *container.Container, index int) *OOBManagedNodesZone {
	OOBManagedNodesZoneCont := cont.S("imdata").Index(index).S(MgmtoobzoneClassName, "attributes")
	return &OOBManagedNodesZone{
		BaseAttributes{
			DistinguishedName: G(OOBManagedNodesZoneCont, "dn"),
			Description:       G(OOBManagedNodesZoneCont, "descr"),
			Status:            G(OOBManagedNodesZoneCont, "status"),
			ClassName:         MgmtoobzoneClassName,
			Rn:                G(OOBManagedNodesZoneCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(OOBManagedNodesZoneCont, "nameAlias"),
		},
		OOBManagedNodesZoneAttributes{
			Annotation: G(OOBManagedNodesZoneCont, "annotation"),
			Name:       G(OOBManagedNodesZoneCont, "name"),
		},
	}
}

func OOBManagedNodesZoneFromContainer(cont *container.Container) *OOBManagedNodesZone {
	return OOBManagedNodesZoneFromContainerList(cont, 0)
}

func OOBManagedNodesZoneListFromContainer(cont *container.Container) []*OOBManagedNodesZone {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*OOBManagedNodesZone, length)
	for i := 0; i < length; i++ {
		arr[i] = OOBManagedNodesZoneFromContainerList(cont, i)
	}
	return arr
}
