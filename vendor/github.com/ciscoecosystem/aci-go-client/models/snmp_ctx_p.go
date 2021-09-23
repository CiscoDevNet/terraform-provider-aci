package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnsnmpCtxP        = "uni/tn-%s/ctx-%s/snmpctx"
	RnsnmpCtxP        = "snmpctx"
	ParentDnsnmpCtxP  = "uni/tn-%s/ctx-%s"
	SnmpctxpClassName = "snmpCtxP"
)

type SNMPContextProfile struct {
	BaseAttributes
	NameAliasAttribute
	SNMPContextProfileAttributes
}

type SNMPContextProfileAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewSNMPContextProfile(snmpCtxPRn, parentDn, nameAlias string, snmpCtxPAttr SNMPContextProfileAttributes) *SNMPContextProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, snmpCtxPRn)
	return &SNMPContextProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         SnmpctxpClassName,
			Rn:                snmpCtxPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SNMPContextProfileAttributes: snmpCtxPAttr,
	}
}

func (snmpCtxP *SNMPContextProfile) ToMap() (map[string]string, error) {
	snmpCtxPMap, err := snmpCtxP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := snmpCtxP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(snmpCtxPMap, key, value)
	}
	A(snmpCtxPMap, "annotation", snmpCtxP.Annotation)
	A(snmpCtxPMap, "name", snmpCtxP.Name)
	return snmpCtxPMap, err
}

func SNMPContextProfileFromContainerList(cont *container.Container, index int) *SNMPContextProfile {
	SNMPContextProfileCont := cont.S("imdata").Index(index).S(SnmpctxpClassName, "attributes")
	return &SNMPContextProfile{
		BaseAttributes{
			DistinguishedName: G(SNMPContextProfileCont, "dn"),
			Status:            G(SNMPContextProfileCont, "status"),
			ClassName:         SnmpctxpClassName,
			Rn:                G(SNMPContextProfileCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SNMPContextProfileCont, "nameAlias"),
		},
		SNMPContextProfileAttributes{
			Annotation: G(SNMPContextProfileCont, "annotation"),
			Name:       G(SNMPContextProfileCont, "name"),
		},
	}
}

func SNMPContextProfileFromContainer(cont *container.Container) *SNMPContextProfile {
	return SNMPContextProfileFromContainerList(cont, 0)
}

func SNMPContextProfileListFromContainer(cont *container.Container) []*SNMPContextProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SNMPContextProfile, length)
	for i := 0; i < length; i++ {
		arr[i] = SNMPContextProfileFromContainerList(cont, i)
	}
	return arr
}
