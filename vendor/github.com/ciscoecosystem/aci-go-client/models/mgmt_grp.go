package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MgmtgrpClassName = "mgmtGrp"

type ManagedNodeConnectivityGroup struct {
	BaseAttributes
	ManagedNodeConnectivityGroupAttributes
}

type ManagedNodeConnectivityGroupAttributes struct {
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
}

func NewManagedNodeConnectivityGroup(mgmtGrpRn, parentDn string, mgmtGrpattr ManagedNodeConnectivityGroupAttributes) *ManagedNodeConnectivityGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtGrpRn)
	return &ManagedNodeConnectivityGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         MgmtgrpClassName,
			Rn:                mgmtGrpRn,
		},

		ManagedNodeConnectivityGroupAttributes: mgmtGrpattr,
	}
}

func (mgmtGrp *ManagedNodeConnectivityGroup) ToMap() (map[string]string, error) {
	mgmtGrpMap, err := mgmtGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(mgmtGrpMap, "name", mgmtGrp.Name)

	A(mgmtGrpMap, "annotation", mgmtGrp.Annotation)

	return mgmtGrpMap, err
}

func ManagedNodeConnectivityGroupFromContainerList(cont *container.Container, index int) *ManagedNodeConnectivityGroup {

	ManagedNodeConnectivityGroupCont := cont.S("imdata").Index(index).S(MgmtgrpClassName, "attributes")
	return &ManagedNodeConnectivityGroup{
		BaseAttributes{
			DistinguishedName: G(ManagedNodeConnectivityGroupCont, "dn"),
			Status:            G(ManagedNodeConnectivityGroupCont, "status"),
			ClassName:         MgmtgrpClassName,
			Rn:                G(ManagedNodeConnectivityGroupCont, "rn"),
		},

		ManagedNodeConnectivityGroupAttributes{

			Name:       G(ManagedNodeConnectivityGroupCont, "name"),
			Annotation: G(ManagedNodeConnectivityGroupCont, "annotation"),
		},
	}
}

func ManagedNodeConnectivityGroupFromContainer(cont *container.Container) *ManagedNodeConnectivityGroup {

	return ManagedNodeConnectivityGroupFromContainerList(cont, 0)
}

func ManagedNodeConnectivityGroupListFromContainer(cont *container.Container) []*ManagedNodeConnectivityGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ManagedNodeConnectivityGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = ManagedNodeConnectivityGroupFromContainerList(cont, i)
	}

	return arr
}
