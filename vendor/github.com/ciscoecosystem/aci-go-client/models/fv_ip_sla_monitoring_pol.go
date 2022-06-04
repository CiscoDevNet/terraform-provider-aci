package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvIPSLAMonitoringPol        = "uni/tn-%s/ipslaMonitoringPol-%s"
	RnfvIPSLAMonitoringPol        = "ipslaMonitoringPol-%s"
	ParentDnfvIPSLAMonitoringPol  = "uni/tn-%s"
	FvipslamonitoringpolClassName = "fvIPSLAMonitoringPol"
)

type IPSLAMonitoringPolicy struct {
	BaseAttributes
	NameAliasAttribute
	IPSLAMonitoringPolicyAttributes
}

type IPSLAMonitoringPolicyAttributes struct {
	Annotation          string `json:",omitempty"`
	HttpUri             string `json:",omitempty"`
	HttpVersion         string `json:",omitempty"`
	Ipv4Tos             string `json:",omitempty"`
	Ipv6TrfClass        string `json:",omitempty"`
	Name                string `json:",omitempty"`
	ReqDataSize         string `json:",omitempty"`
	SlaDetectMultiplier string `json:",omitempty"`
	SlaFrequency        string `json:",omitempty"`
	SlaPort             string `json:",omitempty"`
	SlaType             string `json:",omitempty"`
	Threshold           string `json:",omitempty"`
	Timeout             string `json:",omitempty"`
}

func NewIPSLAMonitoringPolicy(fvIPSLAMonitoringPolRn, parentDn, description, nameAlias string, fvIPSLAMonitoringPolAttr IPSLAMonitoringPolicyAttributes) *IPSLAMonitoringPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, fvIPSLAMonitoringPolRn)
	return &IPSLAMonitoringPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvipslamonitoringpolClassName,
			Rn:                fvIPSLAMonitoringPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		IPSLAMonitoringPolicyAttributes: fvIPSLAMonitoringPolAttr,
	}
}

func (fvIPSLAMonitoringPol *IPSLAMonitoringPolicy) ToMap() (map[string]string, error) {
	fvIPSLAMonitoringPolMap, err := fvIPSLAMonitoringPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := fvIPSLAMonitoringPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(fvIPSLAMonitoringPolMap, key, value)
	}

	A(fvIPSLAMonitoringPolMap, "annotation", fvIPSLAMonitoringPol.Annotation)
	A(fvIPSLAMonitoringPolMap, "httpUri", fvIPSLAMonitoringPol.HttpUri)
	A(fvIPSLAMonitoringPolMap, "httpVersion", fvIPSLAMonitoringPol.HttpVersion)
	A(fvIPSLAMonitoringPolMap, "ipv4Tos", fvIPSLAMonitoringPol.Ipv4Tos)
	A(fvIPSLAMonitoringPolMap, "ipv6TrfClass", fvIPSLAMonitoringPol.Ipv6TrfClass)
	A(fvIPSLAMonitoringPolMap, "name", fvIPSLAMonitoringPol.Name)
	A(fvIPSLAMonitoringPolMap, "reqDataSize", fvIPSLAMonitoringPol.ReqDataSize)
	A(fvIPSLAMonitoringPolMap, "slaDetectMultiplier", fvIPSLAMonitoringPol.SlaDetectMultiplier)
	A(fvIPSLAMonitoringPolMap, "slaFrequency", fvIPSLAMonitoringPol.SlaFrequency)
	A(fvIPSLAMonitoringPolMap, "slaPort", fvIPSLAMonitoringPol.SlaPort)
	A(fvIPSLAMonitoringPolMap, "slaType", fvIPSLAMonitoringPol.SlaType)
	A(fvIPSLAMonitoringPolMap, "threshold", fvIPSLAMonitoringPol.Threshold)
	A(fvIPSLAMonitoringPolMap, "timeout", fvIPSLAMonitoringPol.Timeout)
	return fvIPSLAMonitoringPolMap, err
}

func IPSLAMonitoringPolicyFromContainerList(cont *container.Container, index int) *IPSLAMonitoringPolicy {
	IPSLAMonitoringPolicyCont := cont.S("imdata").Index(index).S(FvipslamonitoringpolClassName, "attributes")
	return &IPSLAMonitoringPolicy{
		BaseAttributes{
			DistinguishedName: G(IPSLAMonitoringPolicyCont, "dn"),
			Description:       G(IPSLAMonitoringPolicyCont, "descr"),
			Status:            G(IPSLAMonitoringPolicyCont, "status"),
			ClassName:         FvipslamonitoringpolClassName,
			Rn:                G(IPSLAMonitoringPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(IPSLAMonitoringPolicyCont, "nameAlias"),
		},
		IPSLAMonitoringPolicyAttributes{
			Annotation:          G(IPSLAMonitoringPolicyCont, "annotation"),
			HttpUri:             G(IPSLAMonitoringPolicyCont, "httpUri"),
			HttpVersion:         G(IPSLAMonitoringPolicyCont, "httpVersion"),
			Ipv4Tos:             G(IPSLAMonitoringPolicyCont, "ipv4Tos"),
			Ipv6TrfClass:        G(IPSLAMonitoringPolicyCont, "ipv6TrfClass"),
			Name:                G(IPSLAMonitoringPolicyCont, "name"),
			ReqDataSize:         G(IPSLAMonitoringPolicyCont, "reqDataSize"),
			SlaDetectMultiplier: G(IPSLAMonitoringPolicyCont, "slaDetectMultiplier"),
			SlaFrequency:        G(IPSLAMonitoringPolicyCont, "slaFrequency"),
			SlaPort:             G(IPSLAMonitoringPolicyCont, "slaPort"),
			SlaType:             G(IPSLAMonitoringPolicyCont, "slaType"),
			Threshold:           G(IPSLAMonitoringPolicyCont, "threshold"),
			Timeout:             G(IPSLAMonitoringPolicyCont, "timeout"),
		},
	}
}

func IPSLAMonitoringPolicyFromContainer(cont *container.Container) *IPSLAMonitoringPolicy {
	return IPSLAMonitoringPolicyFromContainerList(cont, 0)
}

func IPSLAMonitoringPolicyListFromContainer(cont *container.Container) []*IPSLAMonitoringPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IPSLAMonitoringPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = IPSLAMonitoringPolicyFromContainerList(cont, i)
	}

	return arr
}
