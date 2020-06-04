package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const SyntheticMaintPSwitchDetailsClassName = "syntheticMaintPSwitchDetails"

type MaintPSwitchDetails struct {
	BaseAttributes
	MaintPSwitchDetailsAttributes
}

type MaintPSwitchDetailsAttributes struct {
	MaintPName    string `json:",omitempty"`
	NodeIds       string `json:",omitempty"`
	TargetVersion string `json:",omitempty"`
}

func NewMaintPSwitchDetails(syntheticMaintPSwitchDetailsRn, parentDn, description string, syntheticMaintPSwitchDetailsattr MaintPSwitchDetailsAttributes) *MaintPSwitchDetails {
	dn := fmt.Sprintf("%s/%s", parentDn, syntheticMaintPSwitchDetailsRn)
	return &MaintPSwitchDetails{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         SyntheticMaintPSwitchDetailsClassName,
			Rn:                syntheticMaintPSwitchDetailsRn,
		},

		MaintPSwitchDetailsAttributes: syntheticMaintPSwitchDetailsattr,
	}
}

func (syntheticMaintPSwitchDetails *MaintPSwitchDetails) ToMap() (map[string]string, error) {
	syntheticMaintPSwitchDetailsMap, err := syntheticMaintPSwitchDetails.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(syntheticMaintPSwitchDetailsMap, "maintPName", syntheticMaintPSwitchDetails.MaintPName)
	A(syntheticMaintPSwitchDetailsMap, "nodeIds", syntheticMaintPSwitchDetails.NodeIds)
	A(syntheticMaintPSwitchDetailsMap, "targetVersion", syntheticMaintPSwitchDetails.TargetVersion)

	return syntheticMaintPSwitchDetailsMap, err
}

func MaintPSwitchDetailsFromContainerList(cont *container.Container, index int) *MaintPSwitchDetails {

	MaintPSwitchDetailsCont := cont.S("imdata").Index(index).S(SyntheticMaintPSwitchDetailsClassName, "attributes")
	return &MaintPSwitchDetails{
		BaseAttributes{
			DistinguishedName: G(MaintPSwitchDetailsCont, "dn"),
			Status:            G(MaintPSwitchDetailsCont, "status"),
			ClassName:         SyntheticMaintPSwitchDetailsClassName,
			Rn:                G(MaintPSwitchDetailsCont, "rn"),
		},

		MaintPSwitchDetailsAttributes{
			MaintPName:    G(MaintPSwitchDetailsCont, "maintPName"),
			NodeIds:       G(MaintPSwitchDetailsCont, "nodeIds"),
			TargetVersion: G(MaintPSwitchDetailsCont, "targetVersion"),
		},
	}
}

func MaintPSwitchDetailsFromContainer(cont *container.Container) *MaintPSwitchDetails {

	return MaintPSwitchDetailsFromContainerList(cont, 0)
}

func MaintPSwitchDetailsListFromContainer(cont *container.Container) []*MaintPSwitchDetails {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MaintPSwitchDetails, length)

	for i := 0; i < length; i++ {

		arr[i] = MaintPSwitchDetailsFromContainerList(cont, i)
	}

	return arr
}
