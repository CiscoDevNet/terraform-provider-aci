package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L2portsecuritypolClassName = "l2PortSecurityPol"

type PortSecurityPolicy struct {
	BaseAttributes
	PortSecurityPolicyAttributes
}

type PortSecurityPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Maximum string `json:",omitempty"`

	Mode string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Timeout string `json:",omitempty"`

	Violation string `json:",omitempty"`
}

func NewPortSecurityPolicy(l2PortSecurityPolRn, parentDn, description string, l2PortSecurityPolattr PortSecurityPolicyAttributes) *PortSecurityPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l2PortSecurityPolRn)
	return &PortSecurityPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L2portsecuritypolClassName,
			Rn:                l2PortSecurityPolRn,
		},

		PortSecurityPolicyAttributes: l2PortSecurityPolattr,
	}
}

func (l2PortSecurityPol *PortSecurityPolicy) ToMap() (map[string]string, error) {
	l2PortSecurityPolMap, err := l2PortSecurityPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l2PortSecurityPolMap, "name", l2PortSecurityPol.Name)

	A(l2PortSecurityPolMap, "annotation", l2PortSecurityPol.Annotation)

	A(l2PortSecurityPolMap, "maximum", l2PortSecurityPol.Maximum)

	A(l2PortSecurityPolMap, "mode", l2PortSecurityPol.Mode)

	A(l2PortSecurityPolMap, "nameAlias", l2PortSecurityPol.NameAlias)

	A(l2PortSecurityPolMap, "timeout", l2PortSecurityPol.Timeout)

	A(l2PortSecurityPolMap, "violation", l2PortSecurityPol.Violation)

	return l2PortSecurityPolMap, err
}

func PortSecurityPolicyFromContainerList(cont *container.Container, index int) *PortSecurityPolicy {

	PortSecurityPolicyCont := cont.S("imdata").Index(index).S(L2portsecuritypolClassName, "attributes")
	return &PortSecurityPolicy{
		BaseAttributes{
			DistinguishedName: G(PortSecurityPolicyCont, "dn"),
			Description:       G(PortSecurityPolicyCont, "descr"),
			Status:            G(PortSecurityPolicyCont, "status"),
			ClassName:         L2portsecuritypolClassName,
			Rn:                G(PortSecurityPolicyCont, "rn"),
		},

		PortSecurityPolicyAttributes{

			Name: G(PortSecurityPolicyCont, "name"),

			Annotation: G(PortSecurityPolicyCont, "annotation"),

			Maximum: G(PortSecurityPolicyCont, "maximum"),

			Mode: G(PortSecurityPolicyCont, "mode"),

			NameAlias: G(PortSecurityPolicyCont, "nameAlias"),

			Timeout: G(PortSecurityPolicyCont, "timeout"),

			Violation: G(PortSecurityPolicyCont, "violation"),
		},
	}
}

func PortSecurityPolicyFromContainer(cont *container.Container) *PortSecurityPolicy {

	return PortSecurityPolicyFromContainerList(cont, 0)
}

func PortSecurityPolicyListFromContainer(cont *container.Container) []*PortSecurityPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*PortSecurityPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = PortSecurityPolicyFromContainerList(cont, i)
	}

	return arr
}
