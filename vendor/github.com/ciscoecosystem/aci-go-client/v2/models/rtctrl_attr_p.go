package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlAttrP        = "uni/tn-%s/attr-%s"
	RnrtctrlAttrP        = "attr-%s"
	ParentDnrtctrlAttrP  = "uni/tn-%s"
	RtctrlattrpClassName = "rtctrlAttrP"
)

type ActionRuleProfile struct {
	BaseAttributes
	NameAliasAttribute
	ActionRuleProfileAttributes
}

type ActionRuleProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewActionRuleProfile(rtctrlAttrPRn, parentDn, description, nameAlias string, rtctrlAttrPAttr ActionRuleProfileAttributes) *ActionRuleProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlAttrPRn)
	return &ActionRuleProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlattrpClassName,
			Rn:                rtctrlAttrPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ActionRuleProfileAttributes: rtctrlAttrPAttr,
	}
}

func (rtctrlAttrP *ActionRuleProfile) ToMap() (map[string]string, error) {
	rtctrlAttrPMap, err := rtctrlAttrP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlAttrP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlAttrPMap, key, value)
	}

	A(rtctrlAttrPMap, "annotation", rtctrlAttrP.Annotation)
	A(rtctrlAttrPMap, "name", rtctrlAttrP.Name)
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
		NameAliasAttribute{
			NameAlias: G(ActionRuleProfileCont, "nameAlias"),
		},
		ActionRuleProfileAttributes{

			Name: G(ActionRuleProfileCont, "name"),

			Annotation: G(ActionRuleProfileCont, "annotation"),
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
