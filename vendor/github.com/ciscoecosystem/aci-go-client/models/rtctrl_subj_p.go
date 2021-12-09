package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlSubjP        = "uni/tn-%s/subj-%s"
	RnrtctrlSubjP        = "subj-%s"
	ParentDnrtctrlSubjP  = "uni/tn-%s"
	RtctrlsubjpClassName = "rtctrlSubjP"
)

type MatchRule struct {
	BaseAttributes
	NameAliasAttribute
	MatchRuleAttributes
}

type MatchRuleAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewMatchRule(rtctrlSubjPRn, parentDn, description, nameAlias string, rtctrlSubjPAttr MatchRuleAttributes) *MatchRule {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSubjPRn)
	return &MatchRule{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsubjpClassName,
			Rn:                rtctrlSubjPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MatchRuleAttributes: rtctrlSubjPAttr,
	}
}

func (rtctrlSubjP *MatchRule) ToMap() (map[string]string, error) {
	rtctrlSubjPMap, err := rtctrlSubjP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := rtctrlSubjP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(rtctrlSubjPMap, key, value)
	}
	A(rtctrlSubjPMap, "annotation", rtctrlSubjP.Annotation)
	A(rtctrlSubjPMap, "name", rtctrlSubjP.Name)
	return rtctrlSubjPMap, err
}

func MatchRuleFromContainerList(cont *container.Container, index int) *MatchRule {
	MatchRuleCont := cont.S("imdata").Index(index).S(RtctrlsubjpClassName, "attributes")
	return &MatchRule{
		BaseAttributes{
			DistinguishedName: G(MatchRuleCont, "dn"),
			Description:       G(MatchRuleCont, "descr"),
			Status:            G(MatchRuleCont, "status"),
			ClassName:         RtctrlsubjpClassName,
			Rn:                G(MatchRuleCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MatchRuleCont, "nameAlias"),
		},
		MatchRuleAttributes{
			Annotation: G(MatchRuleCont, "annotation"),
			Name:       G(MatchRuleCont, "name"),
		},
	}
}

func MatchRuleFromContainer(cont *container.Container) *MatchRule {
	return MatchRuleFromContainerList(cont, 0)
}

func MatchRuleListFromContainer(cont *container.Container) []*MatchRule {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MatchRule, length)
	for i := 0; i < length; i++ {
		arr[i] = MatchRuleFromContainerList(cont, i)
	}
	return arr
}
