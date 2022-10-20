package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnleakInternalPrefix        = "uni/tn-%s/ctx-%s/leakroutes/leakintprefix-[%s]"
	RnleakInternalPrefix        = "leakintprefix-[%s]"
	ParentDnleakInternalPrefix  = "uni/tn-%s/ctx-%s/leakroutes"
	LeakInternalPrefixClassName = "leakInternalPrefix"
)

type LeakInternalPrefix struct {
	BaseAttributes
	NameAliasAttribute
	LeakInternalPrefixAttributes
}

type LeakInternalPrefixAttributes struct {
	Annotation      string `json:",omitempty"`
	Ip              string `json:",omitempty"`
	Name            string `json:",omitempty"`
	LessThanOrEqual string `json:",omitempty"`
}

func NewLeakInternalPrefix(leakInternalPrefixRn, parentDn, description, nameAlias string, leakInternalPrefixAttr LeakInternalPrefixAttributes) *LeakInternalPrefix {
	dn := fmt.Sprintf("%s/%s", parentDn, leakInternalPrefixRn)
	return &LeakInternalPrefix{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LeakInternalPrefixClassName,
			Rn:                leakInternalPrefixRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LeakInternalPrefixAttributes: leakInternalPrefixAttr,
	}
}

func (leakInternalPrefix *LeakInternalPrefix) ToMap() (map[string]string, error) {
	leakInternalPrefixMap, err := leakInternalPrefix.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := leakInternalPrefix.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(leakInternalPrefixMap, key, value)
	}

	A(leakInternalPrefixMap, "annotation", leakInternalPrefix.Annotation)
	A(leakInternalPrefixMap, "ip", leakInternalPrefix.Ip)
	A(leakInternalPrefixMap, "name", leakInternalPrefix.Name)
	A(leakInternalPrefixMap, "le", leakInternalPrefix.LessThanOrEqual)
	return leakInternalPrefixMap, err
}

func LeakInternalPrefixFromContainerList(cont *container.Container, index int) *LeakInternalPrefix {
	LeakInternalPrefixCont := cont.S("imdata").Index(index).S(LeakInternalPrefixClassName, "attributes")
	return &LeakInternalPrefix{
		BaseAttributes{
			DistinguishedName: G(LeakInternalPrefixCont, "dn"),
			Description:       G(LeakInternalPrefixCont, "descr"),
			Status:            G(LeakInternalPrefixCont, "status"),
			ClassName:         LeakInternalPrefixClassName,
			Rn:                G(LeakInternalPrefixCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LeakInternalPrefixCont, "nameAlias"),
		},
		LeakInternalPrefixAttributes{
			Annotation:      G(LeakInternalPrefixCont, "annotation"),
			Ip:              G(LeakInternalPrefixCont, "ip"),
			Name:            G(LeakInternalPrefixCont, "name"),
			LessThanOrEqual: G(LeakInternalPrefixCont, "le"),
		},
	}
}

func LeakInternalPrefixFromContainer(cont *container.Container) *LeakInternalPrefix {
	return LeakInternalPrefixFromContainerList(cont, 0)
}

func LeakInternalPrefixListFromContainer(cont *container.Container) []*LeakInternalPrefix {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LeakInternalPrefix, length)

	for i := 0; i < length; i++ {
		arr[i] = LeakInternalPrefixFromContainerList(cont, i)
	}

	return arr
}
