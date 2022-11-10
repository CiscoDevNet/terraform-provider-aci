package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetASPathASN        = "uni/tn-%s/attr-%s/saspath-%s/asn-%s"
	RnrtctrlSetASPathASN        = "asn-%s"
	ParentDnrtctrlSetASPathASN  = "uni/tn-%s/attr-%s/saspath-%s"
	RtctrlsetaspathasnClassName = "rtctrlSetASPathASN"
)

type ASNumber struct {
	BaseAttributes
	NameAliasAttribute
	ASNumberAttributes
}

type ASNumberAttributes struct {
	Annotation string `json:",omitempty"`
	Asn        string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Order      string `json:",omitempty"`
}

func NewASNumber(rtctrlSetASPathASNRn, parentDn, description, nameAlias string, rtctrlSetASPathASNAttr ASNumberAttributes) *ASNumber {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetASPathASNRn)
	return &ASNumber{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetaspathasnClassName,
			Rn:                rtctrlSetASPathASNRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ASNumberAttributes: rtctrlSetASPathASNAttr,
	}
}

func (rtctrlSetASPathASN *ASNumber) ToMap() (map[string]string, error) {
	rtctrlSetASPathASNMap, err := rtctrlSetASPathASN.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetASPathASN.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetASPathASNMap, key, value)
	}

	A(rtctrlSetASPathASNMap, "annotation", rtctrlSetASPathASN.Annotation)
	A(rtctrlSetASPathASNMap, "asn", rtctrlSetASPathASN.Asn)
	A(rtctrlSetASPathASNMap, "name", rtctrlSetASPathASN.Name)
	A(rtctrlSetASPathASNMap, "order", rtctrlSetASPathASN.Order)
	return rtctrlSetASPathASNMap, err
}

func ASNumberFromContainerList(cont *container.Container, index int) *ASNumber {
	ASNumberCont := cont.S("imdata").Index(index).S(RtctrlsetaspathasnClassName, "attributes")
	return &ASNumber{
		BaseAttributes{
			DistinguishedName: G(ASNumberCont, "dn"),
			Description:       G(ASNumberCont, "descr"),
			Status:            G(ASNumberCont, "status"),
			ClassName:         RtctrlsetaspathasnClassName,
			Rn:                G(ASNumberCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ASNumberCont, "nameAlias"),
		},
		ASNumberAttributes{
			Annotation: G(ASNumberCont, "annotation"),
			Asn:        G(ASNumberCont, "asn"),
			Name:       G(ASNumberCont, "name"),
			Order:      G(ASNumberCont, "order"),
		},
	}
}

func ASNumberFromContainer(cont *container.Container) *ASNumber {
	return ASNumberFromContainerList(cont, 0)
}

func ASNumberListFromContainer(cont *container.Container) []*ASNumber {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ASNumber, length)

	for i := 0; i < length; i++ {
		arr[i] = ASNumberFromContainerList(cont, i)
	}

	return arr
}
