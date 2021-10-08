package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnisisDomPol        = "uni/fabric/isisDomP-%s"
	RnisisDomPol        = "isisDomP-%s"
	ParentDnisisDomPol  = "uni/fabric"
	IsisdompolClassName = "isisDomPol"
)

type ISISDomainPolicy struct {
	BaseAttributes
	NameAliasAttribute
	ISISDomainPolicyAttributes
}

type ISISDomainPolicyAttributes struct {
	Annotation      string `json:",omitempty"`
	Mtu             string `json:",omitempty"`
	Name            string `json:",omitempty"`
	RedistribMetric string `json:",omitempty"`
}

func NewISISDomainPolicy(isisDomPolRn, parentDn, description, nameAlias string, isisDomPolAttr ISISDomainPolicyAttributes) *ISISDomainPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, isisDomPolRn)
	return &ISISDomainPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IsisdompolClassName,
			Rn:                isisDomPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ISISDomainPolicyAttributes: isisDomPolAttr,
	}
}

func (isisDomPol *ISISDomainPolicy) ToMap() (map[string]string, error) {
	isisDomPolMap, err := isisDomPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := isisDomPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(isisDomPolMap, key, value)
	}
	A(isisDomPolMap, "annotation", isisDomPol.Annotation)
	A(isisDomPolMap, "mtu", isisDomPol.Mtu)
	A(isisDomPolMap, "name", isisDomPol.Name)
	A(isisDomPolMap, "redistribMetric", isisDomPol.RedistribMetric)
	return isisDomPolMap, err
}

func ISISDomainPolicyFromContainerList(cont *container.Container, index int) *ISISDomainPolicy {
	ISISDomainPolicyCont := cont.S("imdata").Index(index).S(IsisdompolClassName, "attributes")
	return &ISISDomainPolicy{
		BaseAttributes{
			DistinguishedName: G(ISISDomainPolicyCont, "dn"),
			Description:       G(ISISDomainPolicyCont, "descr"),
			Status:            G(ISISDomainPolicyCont, "status"),
			ClassName:         IsisdompolClassName,
			Rn:                G(ISISDomainPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ISISDomainPolicyCont, "nameAlias"),
		},
		ISISDomainPolicyAttributes{
			Annotation:      G(ISISDomainPolicyCont, "annotation"),
			Mtu:             G(ISISDomainPolicyCont, "mtu"),
			Name:            G(ISISDomainPolicyCont, "name"),
			RedistribMetric: G(ISISDomainPolicyCont, "redistribMetric"),
		},
	}
}

func ISISDomainPolicyFromContainer(cont *container.Container) *ISISDomainPolicy {
	return ISISDomainPolicyFromContainerList(cont, 0)
}

func ISISDomainPolicyListFromContainer(cont *container.Container) []*ISISDomainPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ISISDomainPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = ISISDomainPolicyFromContainerList(cont, i)
	}
	return arr
}
