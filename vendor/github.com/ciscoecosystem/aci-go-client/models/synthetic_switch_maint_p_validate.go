package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const SyntheticSwitchMaintPValidateClassName = "syntheticSwitchMaintPValidate"

type SwitchMaintPValidate struct {
	BaseAttributes
	SwitchMaintPValidateAttributes
}

type SwitchMaintPValidateAttributes struct {
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
}

func NewSwitchMaintPValidate(syntheticSwitchMaintPValidateRn, parentDn, description string, syntheticSwitchMaintPValidateattr SwitchMaintPValidateAttributes) *SwitchMaintPValidate {
	dn := fmt.Sprintf("%s/%s", parentDn, syntheticSwitchMaintPValidateRn)
	return &SwitchMaintPValidate{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         SyntheticSwitchMaintPValidateClassName,
			Rn:                syntheticSwitchMaintPValidateRn,
		},

		SwitchMaintPValidateAttributes: syntheticSwitchMaintPValidateattr,
	}
}

func (syntheticSwitchMaintPValidate *SwitchMaintPValidate) ToMap() (map[string]string, error) {
	syntheticSwitchMaintPValidateMap, err := syntheticSwitchMaintPValidate.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(syntheticSwitchMaintPValidateMap, "annotation", syntheticSwitchMaintPValidate.Annotation)
	A(syntheticSwitchMaintPValidateMap, "childAction", syntheticSwitchMaintPValidate.ChildAction)
	A(syntheticSwitchMaintPValidateMap, "criticality", syntheticSwitchMaintPValidate.Criticality)
	A(syntheticSwitchMaintPValidateMap, "extmngdby", syntheticSwitchMaintPValidate.ExtMngdBy)
	A(syntheticSwitchMaintPValidateMap, "lcOwn", syntheticSwitchMaintPValidate.LcOwn)
	A(syntheticSwitchMaintPValidateMap, "modTs", syntheticSwitchMaintPValidate.ModTs)
	A(syntheticSwitchMaintPValidateMap, "name", syntheticSwitchMaintPValidate.Name)
	A(syntheticSwitchMaintPValidateMap, "reason", syntheticSwitchMaintPValidate.Reason)
	A(syntheticSwitchMaintPValidateMap, "recommended_action", syntheticSwitchMaintPValidate.Recommended_action)
	A(syntheticSwitchMaintPValidateMap, "result", syntheticSwitchMaintPValidate.Result)

	return syntheticSwitchMaintPValidateMap, err
}

func SwitchMaintPValidateFromContainerList(cont *container.Container, index int) *SwitchMaintPValidate {

	SwitchMaintPValidateCont := cont.S("imdata").Index(index).S(SyntheticSwitchMaintPValidateClassName, "attributes")
	return &SwitchMaintPValidate{
		BaseAttributes{
			DistinguishedName: G(SwitchMaintPValidateCont, "dn"),
			Status:            G(SwitchMaintPValidateCont, "status"),
			ClassName:         SyntheticSwitchMaintPValidateClassName,
			Rn:                G(SwitchMaintPValidateCont, "rn"),
		},

		SwitchMaintPValidateAttributes{
			Annotation:         G(SwitchMaintPValidateCont, "annotation"),
			ChildAction:        G(SwitchMaintPValidateCont, "childAction"),
			Criticality:        G(SwitchMaintPValidateCont, "criticality"),
			ExtMngdBy:          G(SwitchMaintPValidateCont, "extmngdby"),
			LcOwn:              G(SwitchMaintPValidateCont, "lcOwn"),
			ModTs:              G(SwitchMaintPValidateCont, "modTs"),
			Name:               G(SwitchMaintPValidateCont, "name"),
			Reason:             G(SwitchMaintPValidateCont, "reason"),
			Recommended_action: G(SwitchMaintPValidateCont, "recommended_action"),
			Result:             G(SwitchMaintPValidateCont, "result"),
		},
	}
}

func SwitchMaintPValidateFromContainer(cont *container.Container) *SwitchMaintPValidate {

	return SwitchMaintPValidateFromContainerList(cont, 0)
}

func SwitchMaintPValidateListFromContainer(cont *container.Container) []*SwitchMaintPValidate {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SwitchMaintPValidate, length)

	for i := 0; i < length; i++ {

		arr[i] = SwitchMaintPValidateFromContainerList(cont, i)
	}

	return arr
}
