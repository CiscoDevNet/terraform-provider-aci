package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlMatchRtDest        = "uni/tn-%s/subj-%s/dest-[%s]"
	RnrtctrlMatchRtDest        = "dest-[%s]"
	ParentDnrtctrlMatchRtDest  = "uni/tn-%s/subj-%s"
	RtctrlmatchrtdestClassName = "rtctrlMatchRtDest"
)

type MatchRouteDestinationRule struct {
	BaseAttributes
	NameAliasAttribute
	MatchRouteDestinationRuleAttributes
}

type MatchRouteDestinationRuleAttributes struct {
	Aggregate  string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	FromPfxLen string `json:",omitempty"`
	Ip         string `json:",omitempty"`
	Name       string `json:",omitempty"`
	ToPfxLen   string `json:",omitempty"`
}

func NewMatchRouteDestinationRule(rtctrlMatchRtDestRn, parentDn, description, nameAlias string, rtctrlMatchRtDestAttr MatchRouteDestinationRuleAttributes) *MatchRouteDestinationRule {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlMatchRtDestRn)
	return &MatchRouteDestinationRule{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlmatchrtdestClassName,
			Rn:                rtctrlMatchRtDestRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MatchRouteDestinationRuleAttributes: rtctrlMatchRtDestAttr,
	}
}

func (rtctrlMatchRtDest *MatchRouteDestinationRule) ToMap() (map[string]string, error) {
	rtctrlMatchRtDestMap, err := rtctrlMatchRtDest.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := rtctrlMatchRtDest.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(rtctrlMatchRtDestMap, key, value)
	}
	A(rtctrlMatchRtDestMap, "aggregate", rtctrlMatchRtDest.Aggregate)
	A(rtctrlMatchRtDestMap, "annotation", rtctrlMatchRtDest.Annotation)
	A(rtctrlMatchRtDestMap, "fromPfxLen", rtctrlMatchRtDest.FromPfxLen)
	A(rtctrlMatchRtDestMap, "ip", rtctrlMatchRtDest.Ip)
	A(rtctrlMatchRtDestMap, "name", rtctrlMatchRtDest.Name)
	A(rtctrlMatchRtDestMap, "toPfxLen", rtctrlMatchRtDest.ToPfxLen)
	return rtctrlMatchRtDestMap, err
}

func MatchRouteDestinationRuleFromContainerList(cont *container.Container, index int) *MatchRouteDestinationRule {
	MatchRouteDestinationRuleCont := cont.S("imdata").Index(index).S(RtctrlmatchrtdestClassName, "attributes")
	return &MatchRouteDestinationRule{
		BaseAttributes{
			DistinguishedName: G(MatchRouteDestinationRuleCont, "dn"),
			Description:       G(MatchRouteDestinationRuleCont, "descr"),
			Status:            G(MatchRouteDestinationRuleCont, "status"),
			ClassName:         RtctrlmatchrtdestClassName,
			Rn:                G(MatchRouteDestinationRuleCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MatchRouteDestinationRuleCont, "nameAlias"),
		},
		MatchRouteDestinationRuleAttributes{
			Aggregate:  G(MatchRouteDestinationRuleCont, "aggregate"),
			Annotation: G(MatchRouteDestinationRuleCont, "annotation"),
			FromPfxLen: G(MatchRouteDestinationRuleCont, "fromPfxLen"),
			Ip:         G(MatchRouteDestinationRuleCont, "ip"),
			Name:       G(MatchRouteDestinationRuleCont, "name"),
			ToPfxLen:   G(MatchRouteDestinationRuleCont, "toPfxLen"),
		},
	}
}

func MatchRouteDestinationRuleFromContainer(cont *container.Container) *MatchRouteDestinationRule {
	return MatchRouteDestinationRuleFromContainerList(cont, 0)
}

func MatchRouteDestinationRuleListFromContainer(cont *container.Container) []*MatchRouteDestinationRule {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MatchRouteDestinationRule, length)
	for i := 0; i < length; i++ {
		arr[i] = MatchRouteDestinationRuleFromContainerList(cont, i)
	}
	return arr
}
