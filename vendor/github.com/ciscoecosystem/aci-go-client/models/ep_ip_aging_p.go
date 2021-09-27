package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnepIpAgingP        = "uni/infra/ipAgingP-%s"
	RnepIpAgingP        = "ipAgingP-%s"
	ParentDnepIpAgingP  = "uni/infra"
	EpipagingpClassName = "epIpAgingP"
)

type IPAgingPolicy struct {
	BaseAttributes
	NameAliasAttribute
	IPAgingPolicyAttributes
}

type IPAgingPolicyAttributes struct {
	AdminSt    string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewIPAgingPolicy(epIpAgingPRn, parentDn, description, nameAlias string, epIpAgingPAttr IPAgingPolicyAttributes) *IPAgingPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, epIpAgingPRn)
	return &IPAgingPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         EpipagingpClassName,
			Rn:                epIpAgingPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		IPAgingPolicyAttributes: epIpAgingPAttr,
	}
}

func (epIpAgingP *IPAgingPolicy) ToMap() (map[string]string, error) {
	epIpAgingPMap, err := epIpAgingP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := epIpAgingP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(epIpAgingPMap, key, value)
	}
	A(epIpAgingPMap, "adminSt", epIpAgingP.AdminSt)
	A(epIpAgingPMap, "annotation", epIpAgingP.Annotation)
	A(epIpAgingPMap, "name", epIpAgingP.Name)
	return epIpAgingPMap, err
}

func IPAgingPolicyFromContainerList(cont *container.Container, index int) *IPAgingPolicy {
	IPAgingPolicyCont := cont.S("imdata").Index(index).S(EpipagingpClassName, "attributes")
	return &IPAgingPolicy{
		BaseAttributes{
			DistinguishedName: G(IPAgingPolicyCont, "dn"),
			Description:       G(IPAgingPolicyCont, "descr"),
			Status:            G(IPAgingPolicyCont, "status"),
			ClassName:         EpipagingpClassName,
			Rn:                G(IPAgingPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(IPAgingPolicyCont, "nameAlias"),
		},
		IPAgingPolicyAttributes{
			AdminSt:    G(IPAgingPolicyCont, "adminSt"),
			Annotation: G(IPAgingPolicyCont, "annotation"),
			Name:       G(IPAgingPolicyCont, "name"),
		},
	}
}

func IPAgingPolicyFromContainer(cont *container.Container) *IPAgingPolicy {
	return IPAgingPolicyFromContainerList(cont, 0)
}

func IPAgingPolicyListFromContainer(cont *container.Container) []*IPAgingPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IPAgingPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = IPAgingPolicyFromContainerList(cont, i)
	}
	return arr
}
