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

// START: Variable/Struct/Fuction Naming per ACI SDK Model Definitions
const MaintMaintGrpClassName = "maintMaintGrp"

type MaintGrp struct {
	BaseAttributes
	MaintGrpAttributes
}

type MaintGrpAttributes struct {
	Annotation string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	Fwtype     string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewMaintGrp(maintMaintGrpRn, parentDn, description string, maintMaintGrpattr MaintGrpAttributes) *MaintGrp {
	dn := fmt.Sprintf("%s/%s", parentDn, maintMaintGrpRn)
	return &MaintGrp{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         MaintMaintGrpClassName,
			Rn:                maintMaintGrpRn,
		},

		MaintGrpAttributes: maintMaintGrpattr,
	}
}

func (maintMaintGrp *MaintGrp) ToMap() (map[string]string, error) {
	maintMaintGrpMap, err := maintMaintGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(maintMaintGrpMap, "annotation", maintMaintGrp.Annotation)
	A(maintMaintGrpMap, "nameAlias", maintMaintGrp.NameAlias)
	A(maintMaintGrpMap, "fwtype", maintMaintGrp.Fwtype)
	A(maintMaintGrpMap, "type", maintMaintGrp.Type)

	return maintMaintGrpMap, err
}

func MaintGrpFromContainerList(cont *container.Container, index int) *MaintGrp {

	MaintGrpCont := cont.S("imdata").Index(index).S(MaintMaintGrpClassName, "attributes")
	return &MaintGrp{
		BaseAttributes{
			DistinguishedName: G(MaintGrpCont, "dn"),
			Description:       G(MaintGrpCont, "descr"),
			Status:            G(MaintGrpCont, "status"),
			ClassName:         MaintMaintGrpClassName,
			Rn:                G(MaintGrpCont, "rn"),
		},

		MaintGrpAttributes{
			Annotation: G(MaintGrpCont, "annotation"),
			NameAlias:  G(MaintGrpCont, "nameAlias"),
			Fwtype:     G(MaintGrpCont, "fwtype"),
			Type:       G(MaintGrpCont, "type"),
		},
	}
}

func MaintGrpFromContainer(cont *container.Container) *MaintGrp {

	return MaintGrpFromContainerList(cont, 0)
}

func MaintGrpListFromContainer(cont *container.Container) []*MaintGrp {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MaintGrp, length)

	for i := 0; i < length; i++ {

		arr[i] = MaintGrpFromContainerList(cont, i)
	}

	return arr
}
