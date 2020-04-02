package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const RtctrlattrpClassName = "rtctrlAttrP"

type ActionRuleProfile struct {
	BaseAttributes
	ActionRuleProfileAttributes
}

type ActionRuleProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewActionRuleProfile(rtctrlAttrPRn, parentDn, description string, rtctrlAttrPattr ActionRuleProfileAttributes) *ActionRuleProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlAttrPRn)
	return &ActionRuleProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlattrpClassName,
			Rn:                rtctrlAttrPRn,
		},

		ActionRuleProfileAttributes: rtctrlAttrPattr,
	}
}

func (rtctrlAttrP *ActionRuleProfile) ToMap() (map[string]string, error) {
	rtctrlAttrPMap, err := rtctrlAttrP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(rtctrlAttrPMap, "name", rtctrlAttrP.Name)

	A(rtctrlAttrPMap, "annotation", rtctrlAttrP.Annotation)

	A(rtctrlAttrPMap, "nameAlias", rtctrlAttrP.NameAlias)

	return rtctrlAttrPMap, err
}

func ActionRuleProfileFromContainerList(cont *container.Container, index int) *ActionRuleProfile {

	ActionRuleProfileCont := cont.S("imdata").Index(index).S(RtctrlattrpClassName, "attributes")
	return &ActionRuleProfile{
		BaseAttributes{
			DistinguishedName: G(ActionRuleProfileCont, "dn"),
			Description:       G(ActionRuleProfileCont, "descr"),
			Status:            G(ActionRuleProfileCont, "status"),
			ClassName:         RtctrlattrpClassName,
			Rn:                G(ActionRuleProfileCont, "rn"),
		},

		ActionRuleProfileAttributes{

			Name: G(ActionRuleProfileCont, "name"),

			Annotation: G(ActionRuleProfileCont, "annotation"),

			NameAlias: G(ActionRuleProfileCont, "nameAlias"),
		},
	}
}

func ActionRuleProfileFromContainer(cont *container.Container) *ActionRuleProfile {

	return ActionRuleProfileFromContainerList(cont, 0)
}

func ActionRuleProfileListFromContainer(cont *container.Container) []*ActionRuleProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ActionRuleProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = ActionRuleProfileFromContainerList(cont, i)
	}

	return arr
}
