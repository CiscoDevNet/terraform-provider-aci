package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const SyntheticMaintPValidateClassName = "syntheticMaintPValidate"

type MaintPValidate struct {
	BaseAttributes
	MaintPValidateAttributes
}

type MaintPValidateAttributes struct {
	Annotation         string `json:",omitempty"`
	ChildAction        string `json:",omitempty"`
	Criticality        string `json:",omitempty"`
	ExtMngdBy          string `json:",omitempty"`
	LcOwn              string `json:",omitempty"`
	ModTs              string `json:",omitempty"`
	Name               string `json:",omitempty"`
	Reason             string `json:",omitempty"`
	Recommended_action string `json:",omitempty"`
	Result             string `json:",omitempty"`
	Sub_reason         string `json:",omitempty"`
}

func NewMaintPValidate(syntheticMaintPValidateRn, parentDn, description string, syntheticMaintPValidateattr MaintPValidateAttributes) *MaintPValidate {
	dn := fmt.Sprintf("%s/%s", parentDn, syntheticMaintPValidateRn)
	return &MaintPValidate{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         SyntheticMaintPValidateClassName,
			Rn:                syntheticMaintPValidateRn,
		},

		MaintPValidateAttributes: syntheticMaintPValidateattr,
	}
}

func (syntheticMaintPValidate *MaintPValidate) ToMap() (map[string]string, error) {
	syntheticMaintPValidateMap, err := syntheticMaintPValidate.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(syntheticMaintPValidateMap, "annotation", syntheticMaintPValidate.Annotation)
	A(syntheticMaintPValidateMap, "childAction", syntheticMaintPValidate.ChildAction)
	A(syntheticMaintPValidateMap, "criticality", syntheticMaintPValidate.Criticality)
	A(syntheticMaintPValidateMap, "extmngdby", syntheticMaintPValidate.ExtMngdBy)
	A(syntheticMaintPValidateMap, "lcOwn", syntheticMaintPValidate.LcOwn)
	A(syntheticMaintPValidateMap, "modTs", syntheticMaintPValidate.ModTs)
	A(syntheticMaintPValidateMap, "name", syntheticMaintPValidate.Name)
	A(syntheticMaintPValidateMap, "reason", syntheticMaintPValidate.Reason)
	A(syntheticMaintPValidateMap, "recommended_action", syntheticMaintPValidate.Recommended_action)
	A(syntheticMaintPValidateMap, "result", syntheticMaintPValidate.Result)
	A(syntheticMaintPValidateMap, "sub_reason", syntheticMaintPValidate.Sub_reason)

	return syntheticMaintPValidateMap, err
}

func MaintPValidateFromContainerList(cont *container.Container, index int) *MaintPValidate {

	MaintPValidateCont := cont.S("imdata").Index(index).S(SyntheticMaintPValidateClassName, "attributes")
	return &MaintPValidate{
		BaseAttributes{
			DistinguishedName: G(MaintPValidateCont, "dn"),
			Status:            G(MaintPValidateCont, "status"),
			ClassName:         SyntheticMaintPValidateClassName,
			Rn:                G(MaintPValidateCont, "rn"),
		},

		MaintPValidateAttributes{
			Annotation:         G(MaintPValidateCont, "annotation"),
			ChildAction:        G(MaintPValidateCont, "childAction"),
			Criticality:        G(MaintPValidateCont, "criticality"),
			ExtMngdBy:          G(MaintPValidateCont, "extmngdby"),
			LcOwn:              G(MaintPValidateCont, "lcOwn"),
			ModTs:              G(MaintPValidateCont, "modTs"),
			Name:               G(MaintPValidateCont, "name"),
			Reason:             G(MaintPValidateCont, "reason"),
			Recommended_action: G(MaintPValidateCont, "recommended_action"),
			Result:             G(MaintPValidateCont, "result"),
			Sub_reason:         G(MaintPValidateCont, "sub_reason"),
		},
	}
}

func MaintPValidateFromContainer(cont *container.Container) *MaintPValidate {

	return MaintPValidateFromContainerList(cont, 0)
}

func MaintPValidateListFromContainer(cont *container.Container) []*MaintPValidate {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MaintPValidate, length)

	for i := 0; i < length; i++ {

		arr[i] = MaintPValidateFromContainerList(cont, i)
	}

	return arr
}
