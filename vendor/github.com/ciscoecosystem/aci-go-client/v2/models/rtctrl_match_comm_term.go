package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlMatchCommTerm        = "uni/tn-%s/subj-%s/commtrm-%s"
	RnrtctrlMatchCommTerm        = "commtrm-%s"
	ParentDnrtctrlMatchCommTerm  = "uni/tn-%s/subj-%s"
	RtctrlmatchcommtermClassName = "rtctrlMatchCommTerm"
)

type MatchCommunityTerm struct {
	BaseAttributes
	NameAliasAttribute
	MatchCommunityTermAttributes
}

type MatchCommunityTermAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewMatchCommunityTerm(rtctrlMatchCommTermRn, parentDn, description, nameAlias string, rtctrlMatchCommTermAttr MatchCommunityTermAttributes) *MatchCommunityTerm {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlMatchCommTermRn)
	return &MatchCommunityTerm{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlmatchcommtermClassName,
			Rn:                rtctrlMatchCommTermRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MatchCommunityTermAttributes: rtctrlMatchCommTermAttr,
	}
}

func (rtctrlMatchCommTerm *MatchCommunityTerm) ToMap() (map[string]string, error) {
	rtctrlMatchCommTermMap, err := rtctrlMatchCommTerm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlMatchCommTerm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlMatchCommTermMap, key, value)
	}

	A(rtctrlMatchCommTermMap, "annotation", rtctrlMatchCommTerm.Annotation)
	A(rtctrlMatchCommTermMap, "name", rtctrlMatchCommTerm.Name)
	return rtctrlMatchCommTermMap, err
}

func MatchCommunityTermFromContainerList(cont *container.Container, index int) *MatchCommunityTerm {
	MatchCommunityTermCont := cont.S("imdata").Index(index).S(RtctrlmatchcommtermClassName, "attributes")
	return &MatchCommunityTerm{
		BaseAttributes{
			DistinguishedName: G(MatchCommunityTermCont, "dn"),
			Description:       G(MatchCommunityTermCont, "descr"),
			Status:            G(MatchCommunityTermCont, "status"),
			ClassName:         RtctrlmatchcommtermClassName,
			Rn:                G(MatchCommunityTermCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MatchCommunityTermCont, "nameAlias"),
		},
		MatchCommunityTermAttributes{
			Annotation: G(MatchCommunityTermCont, "annotation"),
			Name:       G(MatchCommunityTermCont, "name"),
		},
	}
}

func MatchCommunityTermFromContainer(cont *container.Container) *MatchCommunityTerm {
	return MatchCommunityTermFromContainerList(cont, 0)
}

func MatchCommunityTermListFromContainer(cont *container.Container) []*MatchCommunityTerm {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MatchCommunityTerm, length)

	for i := 0; i < length; i++ {
		arr[i] = MatchCommunityTermFromContainerList(cont, i)
	}

	return arr
}
