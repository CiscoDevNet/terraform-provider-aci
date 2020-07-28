package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnssvcredirectpolClassName = "vnsSvcRedirectPol"

type ServiceRedirectPolicy struct {
	BaseAttributes
	ServiceRedirectPolicyAttributes
}

type ServiceRedirectPolicyAttributes struct {
	Name string `json:",omitempty"`

	AnycastEnabled string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	DestType string `json:",omitempty"`

	HashingAlgorithm string `json:",omitempty"`

	MaxThresholdPercent string `json:",omitempty"`

	MinThresholdPercent string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	ProgramLocalPodOnly string `json:",omitempty"`

	ResilientHashEnabled string `json:",omitempty"`

	ThresholdDownAction string `json:",omitempty"`

	ThresholdEnable string `json:",omitempty"`
}

func NewServiceRedirectPolicy(vnsSvcRedirectPolRn, parentDn, description string, vnsSvcRedirectPolattr ServiceRedirectPolicyAttributes) *ServiceRedirectPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsSvcRedirectPolRn)
	return &ServiceRedirectPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnssvcredirectpolClassName,
			Rn:                vnsSvcRedirectPolRn,
		},

		ServiceRedirectPolicyAttributes: vnsSvcRedirectPolattr,
	}
}

func (vnsSvcRedirectPol *ServiceRedirectPolicy) ToMap() (map[string]string, error) {
	vnsSvcRedirectPolMap, err := vnsSvcRedirectPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsSvcRedirectPolMap, "name", vnsSvcRedirectPol.Name)

	A(vnsSvcRedirectPolMap, "AnycastEnabled", vnsSvcRedirectPol.AnycastEnabled)

	A(vnsSvcRedirectPolMap, "annotation", vnsSvcRedirectPol.Annotation)

	A(vnsSvcRedirectPolMap, "destType", vnsSvcRedirectPol.DestType)

	A(vnsSvcRedirectPolMap, "hashingAlgorithm", vnsSvcRedirectPol.HashingAlgorithm)

	A(vnsSvcRedirectPolMap, "maxThresholdPercent", vnsSvcRedirectPol.MaxThresholdPercent)

	A(vnsSvcRedirectPolMap, "minThresholdPercent", vnsSvcRedirectPol.MinThresholdPercent)

	A(vnsSvcRedirectPolMap, "nameAlias", vnsSvcRedirectPol.NameAlias)

	A(vnsSvcRedirectPolMap, "programLocalPodOnly", vnsSvcRedirectPol.ProgramLocalPodOnly)

	A(vnsSvcRedirectPolMap, "resilientHashEnabled", vnsSvcRedirectPol.ResilientHashEnabled)

	A(vnsSvcRedirectPolMap, "thresholdDownAction", vnsSvcRedirectPol.ThresholdDownAction)

	A(vnsSvcRedirectPolMap, "thresholdEnable", vnsSvcRedirectPol.ThresholdEnable)

	return vnsSvcRedirectPolMap, err
}

func ServiceRedirectPolicyFromContainerList(cont *container.Container, index int) *ServiceRedirectPolicy {

	ServiceRedirectPolicyCont := cont.S("imdata").Index(index).S(VnssvcredirectpolClassName, "attributes")
	return &ServiceRedirectPolicy{
		BaseAttributes{
			DistinguishedName: G(ServiceRedirectPolicyCont, "dn"),
			Description:       G(ServiceRedirectPolicyCont, "descr"),
			Status:            G(ServiceRedirectPolicyCont, "status"),
			ClassName:         VnssvcredirectpolClassName,
			Rn:                G(ServiceRedirectPolicyCont, "rn"),
		},

		ServiceRedirectPolicyAttributes{

			Name: G(ServiceRedirectPolicyCont, "name"),

			AnycastEnabled: G(ServiceRedirectPolicyCont, "AnycastEnabled"),

			Annotation: G(ServiceRedirectPolicyCont, "annotation"),

			DestType: G(ServiceRedirectPolicyCont, "destType"),

			HashingAlgorithm: G(ServiceRedirectPolicyCont, "hashingAlgorithm"),

			MaxThresholdPercent: G(ServiceRedirectPolicyCont, "maxThresholdPercent"),

			MinThresholdPercent: G(ServiceRedirectPolicyCont, "minThresholdPercent"),

			NameAlias: G(ServiceRedirectPolicyCont, "nameAlias"),

			ProgramLocalPodOnly: G(ServiceRedirectPolicyCont, "programLocalPodOnly"),

			ResilientHashEnabled: G(ServiceRedirectPolicyCont, "resilientHashEnabled"),

			ThresholdDownAction: G(ServiceRedirectPolicyCont, "thresholdDownAction"),

			ThresholdEnable: G(ServiceRedirectPolicyCont, "thresholdEnable"),
		},
	}
}

func ServiceRedirectPolicyFromContainer(cont *container.Container) *ServiceRedirectPolicy {

	return ServiceRedirectPolicyFromContainerList(cont, 0)
}

func ServiceRedirectPolicyListFromContainer(cont *container.Container) []*ServiceRedirectPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ServiceRedirectPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = ServiceRedirectPolicyFromContainerList(cont, i)
	}

	return arr
}
