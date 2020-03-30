package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MaintmaintgrpClassName = "maintMaintGrp"

type PODMaintenanceGroup struct {
	BaseAttributes
	PODMaintenanceGroupAttributes
}

type PODMaintenanceGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Fwtype string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PODMaintenanceGroup_type string `json:",omitempty"`
}

func NewPODMaintenanceGroup(maintMaintGrpRn, parentDn, description string, maintMaintGrpattr PODMaintenanceGroupAttributes) *PODMaintenanceGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, maintMaintGrpRn)
	return &PODMaintenanceGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MaintmaintgrpClassName,
			Rn:                maintMaintGrpRn,
		},

		PODMaintenanceGroupAttributes: maintMaintGrpattr,
	}
}

func (maintMaintGrp *PODMaintenanceGroup) ToMap() (map[string]string, error) {
	maintMaintGrpMap, err := maintMaintGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(maintMaintGrpMap, "name", maintMaintGrp.Name)

	A(maintMaintGrpMap, "annotation", maintMaintGrp.Annotation)

	A(maintMaintGrpMap, "fwtype", maintMaintGrp.Fwtype)

	A(maintMaintGrpMap, "nameAlias", maintMaintGrp.NameAlias)

	A(maintMaintGrpMap, "type", maintMaintGrp.PODMaintenanceGroup_type)

	return maintMaintGrpMap, err
}

func PODMaintenanceGroupFromContainerList(cont *container.Container, index int) *PODMaintenanceGroup {

	PODMaintenanceGroupCont := cont.S("imdata").Index(index).S(MaintmaintgrpClassName, "attributes")
	return &PODMaintenanceGroup{
		BaseAttributes{
			DistinguishedName: G(PODMaintenanceGroupCont, "dn"),
			Description:       G(PODMaintenanceGroupCont, "descr"),
			Status:            G(PODMaintenanceGroupCont, "status"),
			ClassName:         MaintmaintgrpClassName,
			Rn:                G(PODMaintenanceGroupCont, "rn"),
		},

		PODMaintenanceGroupAttributes{

			Name: G(PODMaintenanceGroupCont, "name"),

			Annotation: G(PODMaintenanceGroupCont, "annotation"),

			Fwtype: G(PODMaintenanceGroupCont, "fwtype"),

			NameAlias: G(PODMaintenanceGroupCont, "nameAlias"),

			PODMaintenanceGroup_type: G(PODMaintenanceGroupCont, "type"),
		},
	}
}

func PODMaintenanceGroupFromContainer(cont *container.Container) *PODMaintenanceGroup {

	return PODMaintenanceGroupFromContainerList(cont, 0)
}

func PODMaintenanceGroupListFromContainer(cont *container.Container) []*PODMaintenanceGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*PODMaintenanceGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = PODMaintenanceGroupFromContainerList(cont, i)
	}

	return arr
}
