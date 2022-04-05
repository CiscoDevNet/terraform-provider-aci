package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnrtctrlMatchCommFactor        = "uni/tn-%s/subj-%s/commtrm-%s/commfct-%s"
	RnrtctrlMatchCommFactor        = "commfct-%s"
	ParentDnrtctrlMatchCommFactor  = "uni/tn-%s/subj-%s/commtrm-%s"
	RtctrlmatchcommfactorClassName = "rtctrlMatchCommFactor"
)

type MatchCommunityFactor struct {
	BaseAttributes
	NameAliasAttribute
	MatchCommunityFactorAttributes
}

type MatchCommunityFactorAttributes struct {
	Annotation string `json:",omitempty"`
	Community  string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Scope      string `json:",omitempty"`
}

func NewMatchCommunityFactor(rtctrlMatchCommFactorRn, parentDn, description, nameAlias string, rtctrlMatchCommFactorAttr MatchCommunityFactorAttributes) *MatchCommunityFactor {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlMatchCommFactorRn)
	return &MatchCommunityFactor{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlmatchcommfactorClassName,
			Rn:                rtctrlMatchCommFactorRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MatchCommunityFactorAttributes: rtctrlMatchCommFactorAttr,
	}
}

func (rtctrlMatchCommFactor *MatchCommunityFactor) ToMap() (map[string]string, error) {
	rtctrlMatchCommFactorMap, err := rtctrlMatchCommFactor.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlMatchCommFactor.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlMatchCommFactorMap, key, value)
	}

	A(rtctrlMatchCommFactorMap, "annotation", rtctrlMatchCommFactor.Annotation)
	A(rtctrlMatchCommFactorMap, "community", rtctrlMatchCommFactor.Community)
	A(rtctrlMatchCommFactorMap, "name", rtctrlMatchCommFactor.Name)
	A(rtctrlMatchCommFactorMap, "scope", rtctrlMatchCommFactor.Scope)
	return rtctrlMatchCommFactorMap, err
}

func MatchCommunityFactorFromContainerList(cont *container.Container, index int) *MatchCommunityFactor {
	MatchCommunityFactorCont := cont.S("imdata").Index(index).S(RtctrlmatchcommfactorClassName, "attributes")
	return &MatchCommunityFactor{
		BaseAttributes{
			DistinguishedName: G(MatchCommunityFactorCont, "dn"),
			Description:       G(MatchCommunityFactorCont, "descr"),
			Status:            G(MatchCommunityFactorCont, "status"),
			ClassName:         RtctrlmatchcommfactorClassName,
			Rn:                G(MatchCommunityFactorCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MatchCommunityFactorCont, "nameAlias"),
		},
		MatchCommunityFactorAttributes{
			Annotation: G(MatchCommunityFactorCont, "annotation"),
			Community:  G(MatchCommunityFactorCont, "community"),
			Name:       G(MatchCommunityFactorCont, "name"),
			Scope:      G(MatchCommunityFactorCont, "scope"),
		},
	}
}

func MatchCommunityFactorFromContainer(cont *container.Container) *MatchCommunityFactor {
	return MatchCommunityFactorFromContainerList(cont, 0)
}

func MatchCommunityFactorListFromContainer(cont *container.Container) []*MatchCommunityFactor {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MatchCommunityFactor, length)

	for i := 0; i < length; i++ {
		arr[i] = MatchCommunityFactorFromContainerList(cont, i)
	}

	return arr
}
