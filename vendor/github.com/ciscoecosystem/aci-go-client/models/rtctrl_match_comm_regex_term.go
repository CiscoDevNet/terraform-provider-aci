package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlMatchCommRegexTerm        = "uni/tn-%s/subj-%s/commrxtrm-%s"
	RnrtctrlMatchCommRegexTerm        = "commrxtrm-%s"
	ParentDnrtctrlMatchCommRegexTerm  = "uni/tn-%s/subj-%s"
	RtctrlmatchcommregextermClassName = "rtctrlMatchCommRegexTerm"
)

type MatchRuleBasedonCommunityRegularExpression struct {
	BaseAttributes
	NameAliasAttribute
	MatchRuleBasedonCommunityRegularExpressionAttributes
}

type MatchRuleBasedonCommunityRegularExpressionAttributes struct {
	Annotation string `json:",omitempty"`
	CommType   string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Regex      string `json:",omitempty"`
}

func NewMatchRuleBasedonCommunityRegularExpression(rtctrlMatchCommRegexTermRn, parentDn, description, nameAlias string, rtctrlMatchCommRegexTermAttr MatchRuleBasedonCommunityRegularExpressionAttributes) *MatchRuleBasedonCommunityRegularExpression {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlMatchCommRegexTermRn)
	return &MatchRuleBasedonCommunityRegularExpression{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlmatchcommregextermClassName,
			Rn:                rtctrlMatchCommRegexTermRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MatchRuleBasedonCommunityRegularExpressionAttributes: rtctrlMatchCommRegexTermAttr,
	}
}

func (rtctrlMatchCommRegexTerm *MatchRuleBasedonCommunityRegularExpression) ToMap() (map[string]string, error) {
	rtctrlMatchCommRegexTermMap, err := rtctrlMatchCommRegexTerm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlMatchCommRegexTerm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlMatchCommRegexTermMap, key, value)
	}

	A(rtctrlMatchCommRegexTermMap, "annotation", rtctrlMatchCommRegexTerm.Annotation)
	A(rtctrlMatchCommRegexTermMap, "commType", rtctrlMatchCommRegexTerm.CommType)
	A(rtctrlMatchCommRegexTermMap, "name", rtctrlMatchCommRegexTerm.Name)
	A(rtctrlMatchCommRegexTermMap, "regex", rtctrlMatchCommRegexTerm.Regex)
	return rtctrlMatchCommRegexTermMap, err
}

func MatchRuleBasedonCommunityRegularExpressionFromContainerList(cont *container.Container, index int) *MatchRuleBasedonCommunityRegularExpression {
	MatchRuleBasedonCommunityRegularExpressionCont := cont.S("imdata").Index(index).S(RtctrlmatchcommregextermClassName, "attributes")
	return &MatchRuleBasedonCommunityRegularExpression{
		BaseAttributes{
			DistinguishedName: G(MatchRuleBasedonCommunityRegularExpressionCont, "dn"),
			Description:       G(MatchRuleBasedonCommunityRegularExpressionCont, "descr"),
			Status:            G(MatchRuleBasedonCommunityRegularExpressionCont, "status"),
			ClassName:         RtctrlmatchcommregextermClassName,
			Rn:                G(MatchRuleBasedonCommunityRegularExpressionCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MatchRuleBasedonCommunityRegularExpressionCont, "nameAlias"),
		},
		MatchRuleBasedonCommunityRegularExpressionAttributes{
			Annotation: G(MatchRuleBasedonCommunityRegularExpressionCont, "annotation"),
			CommType:   G(MatchRuleBasedonCommunityRegularExpressionCont, "commType"),
			Name:       G(MatchRuleBasedonCommunityRegularExpressionCont, "name"),
			Regex:      G(MatchRuleBasedonCommunityRegularExpressionCont, "regex"),
		},
	}
}

func MatchRuleBasedonCommunityRegularExpressionFromContainer(cont *container.Container) *MatchRuleBasedonCommunityRegularExpression {
	return MatchRuleBasedonCommunityRegularExpressionFromContainerList(cont, 0)
}

func MatchRuleBasedonCommunityRegularExpressionListFromContainer(cont *container.Container) []*MatchRuleBasedonCommunityRegularExpression {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MatchRuleBasedonCommunityRegularExpression, length)

	for i := 0; i < length; i++ {
		arr[i] = MatchRuleBasedonCommunityRegularExpressionFromContainerList(cont, i)
	}

	return arr
}
