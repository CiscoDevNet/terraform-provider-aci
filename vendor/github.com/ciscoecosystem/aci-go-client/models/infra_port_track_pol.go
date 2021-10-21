package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DninfraPortTrackPol        = "uni/infra/trackEqptFabP-%s"
	RninfraPortTrackPol        = "trackEqptFabP-%s"
	ParentDninfraPortTrackPol  = "uni/infra"
	InfraporttrackpolClassName = "infraPortTrackPol"
)

type PortTracking struct {
	BaseAttributes
	NameAliasAttribute
	PortTrackingAttributes
}

type PortTrackingAttributes struct {
	AdminSt          string `json:",omitempty"`
	Annotation       string `json:",omitempty"`
	Delay            string `json:",omitempty"`
	IncludeApicPorts string `json:",omitempty"`
	Minlinks         string `json:",omitempty"`
	Name             string `json:",omitempty"`
}

func NewPortTracking(infraPortTrackPolRn, parentDn, description, nameAlias string, infraPortTrackPolAttr PortTrackingAttributes) *PortTracking {
	dn := fmt.Sprintf("%s/%s", parentDn, infraPortTrackPolRn)
	return &PortTracking{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraporttrackpolClassName,
			Rn:                infraPortTrackPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		PortTrackingAttributes: infraPortTrackPolAttr,
	}
}

func (infraPortTrackPol *PortTracking) ToMap() (map[string]string, error) {
	infraPortTrackPolMap, err := infraPortTrackPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := infraPortTrackPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(infraPortTrackPolMap, key, value)
	}
	A(infraPortTrackPolMap, "adminSt", infraPortTrackPol.AdminSt)
	A(infraPortTrackPolMap, "annotation", infraPortTrackPol.Annotation)
	A(infraPortTrackPolMap, "delay", infraPortTrackPol.Delay)
	A(infraPortTrackPolMap, "includeApicPorts", infraPortTrackPol.IncludeApicPorts)
	A(infraPortTrackPolMap, "minlinks", infraPortTrackPol.Minlinks)
	A(infraPortTrackPolMap, "name", infraPortTrackPol.Name)
	return infraPortTrackPolMap, err
}

func PortTrackingFromContainerList(cont *container.Container, index int) *PortTracking {
	PortTrackingCont := cont.S("imdata").Index(index).S(InfraporttrackpolClassName, "attributes")
	return &PortTracking{
		BaseAttributes{
			DistinguishedName: G(PortTrackingCont, "dn"),
			Description:       G(PortTrackingCont, "descr"),
			Status:            G(PortTrackingCont, "status"),
			ClassName:         InfraporttrackpolClassName,
			Rn:                G(PortTrackingCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(PortTrackingCont, "nameAlias"),
		},
		PortTrackingAttributes{
			AdminSt:          G(PortTrackingCont, "adminSt"),
			Annotation:       G(PortTrackingCont, "annotation"),
			Delay:            G(PortTrackingCont, "delay"),
			IncludeApicPorts: G(PortTrackingCont, "includeApicPorts"),
			Minlinks:         G(PortTrackingCont, "minlinks"),
			Name:             G(PortTrackingCont, "name"),
		},
	}
}

func PortTrackingFromContainer(cont *container.Container) *PortTracking {
	return PortTrackingFromContainerList(cont, 0)
}

func PortTrackingListFromContainer(cont *container.Container) []*PortTracking {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PortTracking, length)
	for i := 0; i < length; i++ {
		arr[i] = PortTrackingFromContainerList(cont, i)
	}
	return arr
}
