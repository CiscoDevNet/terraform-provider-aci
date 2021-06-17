package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrarsaccbasegrpClassName = "infraRsAccBaseGrp"

type AccessAccessGroup struct {
	BaseAttributes
	AccessAccessGroupAttributes
}

type AccessAccessGroupAttributes struct {
	Annotation string `json:",omitempty"`

	FexId string `json:",omitempty"`

	TDn string `json:",omitempty"`
}

func NewAccessAccessGroup(infraRsAccBaseGrpRn, parentDn, description string, infraRsAccBaseGrpattr AccessAccessGroupAttributes) *AccessAccessGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraRsAccBaseGrpRn)
	return &AccessAccessGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrarsaccbasegrpClassName,
			Rn:                infraRsAccBaseGrpRn,
		},

		AccessAccessGroupAttributes: infraRsAccBaseGrpattr,
	}
}

func (infraRsAccBaseGrp *AccessAccessGroup) ToMap() (map[string]string, error) {
	infraRsAccBaseGrpMap, err := infraRsAccBaseGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraRsAccBaseGrpMap, "annotation", infraRsAccBaseGrp.Annotation)

	A(infraRsAccBaseGrpMap, "fexId", infraRsAccBaseGrp.FexId)

	A(infraRsAccBaseGrpMap, "tDn", infraRsAccBaseGrp.TDn)

	return infraRsAccBaseGrpMap, err
}

func AccessAccessGroupFromContainerList(cont *container.Container, index int) *AccessAccessGroup {

	AccessAccessGroupCont := cont.S("imdata").Index(index).S(InfrarsaccbasegrpClassName, "attributes")
	return &AccessAccessGroup{
		BaseAttributes{
			DistinguishedName: G(AccessAccessGroupCont, "dn"),
			Description:       G(AccessAccessGroupCont, "descr"),
			Status:            G(AccessAccessGroupCont, "status"),
			ClassName:         InfrarsaccbasegrpClassName,
			Rn:                G(AccessAccessGroupCont, "rn"),
		},

		AccessAccessGroupAttributes{

			Annotation: G(AccessAccessGroupCont, "annotation"),

			FexId: G(AccessAccessGroupCont, "fexId"),

			TDn: G(AccessAccessGroupCont, "tDn"),
		},
	}
}

func AccessAccessGroupFromContainer(cont *container.Container) *AccessAccessGroup {

	return AccessAccessGroupFromContainerList(cont, 0)
}

func AccessAccessGroupListFromContainer(cont *container.Container) []*AccessAccessGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AccessAccessGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = AccessAccessGroupFromContainerList(cont, i)
	}

	return arr
}
