package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MonepgpolClassName = "monEPGPol"

type MonitoringPolicy struct {
	BaseAttributes
	MonitoringPolicyAttributes
}

type MonitoringPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewMonitoringPolicy(monEPGPolRn, parentDn, description string, monEPGPolattr MonitoringPolicyAttributes) *MonitoringPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, monEPGPolRn)
	return &MonitoringPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MonepgpolClassName,
			Rn:                monEPGPolRn,
		},

		MonitoringPolicyAttributes: monEPGPolattr,
	}
}

func (monEPGPol *MonitoringPolicy) ToMap() (map[string]string, error) {
	monEPGPolMap, err := monEPGPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(monEPGPolMap, "name", monEPGPol.Name)

	A(monEPGPolMap, "annotation", monEPGPol.Annotation)

	A(monEPGPolMap, "nameAlias", monEPGPol.NameAlias)

	return monEPGPolMap, err
}

func MonitoringPolicyFromContainerList(cont *container.Container, index int) *MonitoringPolicy {

	MonitoringPolicyCont := cont.S("imdata").Index(index).S(MonepgpolClassName, "attributes")
	return &MonitoringPolicy{
		BaseAttributes{
			DistinguishedName: G(MonitoringPolicyCont, "dn"),
			Description:       G(MonitoringPolicyCont, "descr"),
			Status:            G(MonitoringPolicyCont, "status"),
			ClassName:         MonepgpolClassName,
			Rn:                G(MonitoringPolicyCont, "rn"),
		},

		MonitoringPolicyAttributes{

			Name: G(MonitoringPolicyCont, "name"),

			Annotation: G(MonitoringPolicyCont, "annotation"),

			NameAlias: G(MonitoringPolicyCont, "nameAlias"),
		},
	}
}

func MonitoringPolicyFromContainer(cont *container.Container) *MonitoringPolicy {

	return MonitoringPolicyFromContainerList(cont, 0)
}

func MonitoringPolicyListFromContainer(cont *container.Container) []*MonitoringPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MonitoringPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = MonitoringPolicyFromContainerList(cont, i)
	}

	return arr
}
