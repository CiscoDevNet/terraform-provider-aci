package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnsnmpCommunityP        = "%s/community-%s"
	RnsnmpCommunityP        = "community-%s"
	SnmpcommunitypClassName = "snmpCommunityP"
)

type SNMPCommunity struct {
	BaseAttributes
	NameAliasAttribute
	SNMPCommunityAttributes
}

type SNMPCommunityAttributes struct {
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
}

func NewSNMPCommunity(snmpCommunityPRn, parentDn, description, nameAlias string, snmpCommunityPAttr SNMPCommunityAttributes) *SNMPCommunity {
	dn := fmt.Sprintf("%s/%s", parentDn, snmpCommunityPRn)
	return &SNMPCommunity{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         SnmpcommunitypClassName,
			Rn:                snmpCommunityPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SNMPCommunityAttributes: snmpCommunityPAttr,
	}
}

func (snmpCommunityP *SNMPCommunity) ToMap() (map[string]string, error) {
	snmpCommunityPMap, err := snmpCommunityP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := snmpCommunityP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(snmpCommunityPMap, key, value)
	}

	A(snmpCommunityPMap, "annotation", snmpCommunityP.Annotation)
	A(snmpCommunityPMap, "name", snmpCommunityP.Name)
	return snmpCommunityPMap, err
}

func SNMPCommunityFromContainerList(cont *container.Container, index int) *SNMPCommunity {
	SNMPCommunityCont := cont.S("imdata").Index(index).S(SnmpcommunitypClassName, "attributes")
	return &SNMPCommunity{
		BaseAttributes{
			DistinguishedName: G(SNMPCommunityCont, "dn"),
			Description:       G(SNMPCommunityCont, "descr"),
			Status:            G(SNMPCommunityCont, "status"),
			ClassName:         SnmpcommunitypClassName,
			Rn:                G(SNMPCommunityCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SNMPCommunityCont, "nameAlias"),
		},
		SNMPCommunityAttributes{
			Annotation: G(SNMPCommunityCont, "annotation"),
			Name:       G(SNMPCommunityCont, "name"),
		},
	}
}

func SNMPCommunityFromContainer(cont *container.Container) *SNMPCommunity {
	return SNMPCommunityFromContainerList(cont, 0)
}

func SNMPCommunityListFromContainer(cont *container.Container) []*SNMPCommunity {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SNMPCommunity, length)

	for i := 0; i < length; i++ {
		arr[i] = SNMPCommunityFromContainerList(cont, i)
	}

	return arr
}
