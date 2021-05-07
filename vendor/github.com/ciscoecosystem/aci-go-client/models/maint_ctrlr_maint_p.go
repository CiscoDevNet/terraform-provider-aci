package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MaintCtrlrMaintPName = "maintCtrlrMaintP"

type CtrlrMaintP struct {
	BaseAttributes
	CtrlrMaintPAttributes
}

type CtrlrMaintPAttributes struct {
	Name                 string `json:",omitempty"`
	Annotation           string `json:",omitempty"`
	AdminSt              string `json:",omitempty"`
	AdminState           string `json:",omitempty"`
	Graceful             string `json:",omitempty"`
	IgnoreCompat         string `json:",omitempty"`
	NotifyCond           string `json:",omitempty"`
	Parallel             string `json:",omitempty"`
	RunMode              string `json:",omitempty"`
	SrUpgrade            string `json:",omitempty"`
	SrVersion            string `json:",omitempty"`
	TriggerTime          string `json:",omitempty"`
	Version              string `json:",omitempty"`
	VersionCheckOverride string `json:",omitempty"`
}

func NewCtrlrMaintP(maintCtrlrMaintPRn, parentDn, description string, maintCtrlrMaintPAttr CtrlrMaintPAttributes) *CtrlrMaintP {
	dn := fmt.Sprintf("%s/%s", parentDn, maintCtrlrMaintPRn)
	return &CtrlrMaintP{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         MaintCtrlrMaintPName,
			Rn:                maintCtrlrMaintPRn,
		},
		CtrlrMaintPAttributes: maintCtrlrMaintPAttr,
	}
}

func (maintCtrlrMaintP *CtrlrMaintP) ToMap() (map[string]string, error) {
	maintCtrlrMaintPMap, err := maintCtrlrMaintP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(maintCtrlrMaintPMap, "name", maintCtrlrMaintP.Name)
	A(maintCtrlrMaintPMap, "annotation", maintCtrlrMaintP.Annotation)
	A(maintCtrlrMaintPMap, "adminSt", maintCtrlrMaintP.AdminSt)
	A(maintCtrlrMaintPMap, "adminState", maintCtrlrMaintP.AdminState)
	A(maintCtrlrMaintPMap, "graceful", maintCtrlrMaintP.Graceful)
	A(maintCtrlrMaintPMap, "ignoreCompat", maintCtrlrMaintP.IgnoreCompat)
	A(maintCtrlrMaintPMap, "notifyCond", maintCtrlrMaintP.NotifyCond)
	A(maintCtrlrMaintPMap, "parallel", maintCtrlrMaintP.Parallel)
	A(maintCtrlrMaintPMap, "runMode", maintCtrlrMaintP.RunMode)
	A(maintCtrlrMaintPMap, "srUpgrade", maintCtrlrMaintP.SrUpgrade)
	A(maintCtrlrMaintPMap, "srVesion", maintCtrlrMaintP.SrVersion)
	A(maintCtrlrMaintPMap, "triggerTime", maintCtrlrMaintP.TriggerTime)
	A(maintCtrlrMaintPMap, "version", maintCtrlrMaintP.Version)
	A(maintCtrlrMaintPMap, "versionCheckOverride", maintCtrlrMaintP.VersionCheckOverride)

	return maintCtrlrMaintPMap, err
}

func CtrlrMaintPFromContainerList(cont *container.Container, index int) *CtrlrMaintP {

	CtrlrMaintPCont := cont.S("imdata").Index(index).S(MaintCtrlrMaintPName, "attributes")
	return &CtrlrMaintP{
		BaseAttributes{
			DistinguishedName: G(CtrlrMaintPCont, "dn"),
			Description:       G(CtrlrMaintPCont, "descr"),
			Status:            G(CtrlrMaintPCont, "status"),
			ClassName:         MaintCtrlrMaintPName,
			Rn:                G(CtrlrMaintPCont, "rn"),
		},

		CtrlrMaintPAttributes{
			Name:                 G(CtrlrMaintPCont, "name"),
			Annotation:           G(CtrlrMaintPCont, "annotation"),
			AdminSt:              G(CtrlrMaintPCont, "adminSt"),
			AdminState:           G(CtrlrMaintPCont, "adminState"),
			Graceful:             G(CtrlrMaintPCont, "graceful"),
			IgnoreCompat:         G(CtrlrMaintPCont, "ignoreCompat"),
			NotifyCond:           G(CtrlrMaintPCont, "notifyCond"),
			Parallel:             G(CtrlrMaintPCont, "parallel"),
			RunMode:              G(CtrlrMaintPCont, "runMode"),
			SrUpgrade:            G(CtrlrMaintPCont, "srUpgrade"),
			SrVersion:            G(CtrlrMaintPCont, "srVesion"),
			TriggerTime:          G(CtrlrMaintPCont, "triggerTime"),
			Version:              G(CtrlrMaintPCont, "version"),
			VersionCheckOverride: G(CtrlrMaintPCont, "versionCheckOverride"),
		},
	}
}

func CtrlrMaintPFromContainer(cont *container.Container) *CtrlrMaintP {

	return CtrlrMaintPFromContainerList(cont, 0)
}

func CtrlrMaintPListFromContainer(cont *container.Container) []*CtrlrMaintP {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CtrlrMaintP, length)

	for i := 0; i < length; i++ {

		arr[i] = CtrlrMaintPFromContainerList(cont, i)
	}

	return arr
}
