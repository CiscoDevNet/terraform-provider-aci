package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MaintmaintpClassName = "maintMaintP"

type MaintenancePolicy struct {
	BaseAttributes
	MaintenancePolicyAttributes
}

type MaintenancePolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Graceful string `json:",omitempty"`

	IgnoreCompat string `json:",omitempty"`

	InternalLabel string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	NotifCond string `json:",omitempty"`

	RunMode string `json:",omitempty"`

	Version string `json:",omitempty"`

	VersionCheckOverride string `json:",omitempty"`
}

func NewMaintenancePolicy(maintMaintPRn, parentDn, description string, maintMaintPattr MaintenancePolicyAttributes) *MaintenancePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, maintMaintPRn)
	return &MaintenancePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MaintmaintpClassName,
			Rn:                maintMaintPRn,
		},

		MaintenancePolicyAttributes: maintMaintPattr,
	}
}

func (maintMaintP *MaintenancePolicy) ToMap() (map[string]string, error) {
	maintMaintPMap, err := maintMaintP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(maintMaintPMap, "name", maintMaintP.Name)

	A(maintMaintPMap, "adminSt", maintMaintP.AdminSt)

	A(maintMaintPMap, "annotation", maintMaintP.Annotation)

	A(maintMaintPMap, "graceful", maintMaintP.Graceful)

	A(maintMaintPMap, "ignoreCompat", maintMaintP.IgnoreCompat)

	A(maintMaintPMap, "internalLabel", maintMaintP.InternalLabel)

	A(maintMaintPMap, "nameAlias", maintMaintP.NameAlias)

	A(maintMaintPMap, "notifCond", maintMaintP.NotifCond)

	A(maintMaintPMap, "runMode", maintMaintP.RunMode)

	A(maintMaintPMap, "version", maintMaintP.Version)

	A(maintMaintPMap, "versionCheckOverride", maintMaintP.VersionCheckOverride)

	return maintMaintPMap, err
}

func MaintenancePolicyFromContainerList(cont *container.Container, index int) *MaintenancePolicy {

	MaintenancePolicyCont := cont.S("imdata").Index(index).S(MaintmaintpClassName, "attributes")
	return &MaintenancePolicy{
		BaseAttributes{
			DistinguishedName: G(MaintenancePolicyCont, "dn"),
			Description:       G(MaintenancePolicyCont, "descr"),
			Status:            G(MaintenancePolicyCont, "status"),
			ClassName:         MaintmaintpClassName,
			Rn:                G(MaintenancePolicyCont, "rn"),
		},

		MaintenancePolicyAttributes{

			Name: G(MaintenancePolicyCont, "name"),

			AdminSt: G(MaintenancePolicyCont, "adminSt"),

			Annotation: G(MaintenancePolicyCont, "annotation"),

			Graceful: G(MaintenancePolicyCont, "graceful"),

			IgnoreCompat: G(MaintenancePolicyCont, "ignoreCompat"),

			InternalLabel: G(MaintenancePolicyCont, "internalLabel"),

			NameAlias: G(MaintenancePolicyCont, "nameAlias"),

			NotifCond: G(MaintenancePolicyCont, "notifCond"),

			RunMode: G(MaintenancePolicyCont, "runMode"),

			Version: G(MaintenancePolicyCont, "version"),

			VersionCheckOverride: G(MaintenancePolicyCont, "versionCheckOverride"),
		},
	}
}

func MaintenancePolicyFromContainer(cont *container.Container) *MaintenancePolicy {

	return MaintenancePolicyFromContainerList(cont, 0)
}

func MaintenancePolicyListFromContainer(cont *container.Container) []*MaintenancePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MaintenancePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = MaintenancePolicyFromContainerList(cont, i)
	}

	return arr
}
