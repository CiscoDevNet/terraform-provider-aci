package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const DhcpoptionpolClassName = "dhcpOptionPol"

type DHCPOptionPolicy struct {
	BaseAttributes
	DHCPOptionPolicyAttributes
}

type DHCPOptionPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewDHCPOptionPolicy(dhcpOptionPolRn, parentDn, description string, dhcpOptionPolattr DHCPOptionPolicyAttributes) *DHCPOptionPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpOptionPolRn)
	return &DHCPOptionPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcpoptionpolClassName,
			Rn:                dhcpOptionPolRn,
		},

		DHCPOptionPolicyAttributes: dhcpOptionPolattr,
	}
}

func (dhcpOptionPol *DHCPOptionPolicy) ToMap() (map[string]string, error) {
	dhcpOptionPolMap, err := dhcpOptionPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpOptionPolMap, "name", dhcpOptionPol.Name)

	A(dhcpOptionPolMap, "annotation", dhcpOptionPol.Annotation)

	A(dhcpOptionPolMap, "nameAlias", dhcpOptionPol.NameAlias)

	return dhcpOptionPolMap, err
}

func DHCPOptionPolicyFromContainerList(cont *container.Container, index int) *DHCPOptionPolicy {

	DHCPOptionPolicyCont := cont.S("imdata").Index(index).S(DhcpoptionpolClassName, "attributes")
	return &DHCPOptionPolicy{
		BaseAttributes{
			DistinguishedName: G(DHCPOptionPolicyCont, "dn"),
			Description:       G(DHCPOptionPolicyCont, "descr"),
			Status:            G(DHCPOptionPolicyCont, "status"),
			ClassName:         DhcpoptionpolClassName,
			Rn:                G(DHCPOptionPolicyCont, "rn"),
		},

		DHCPOptionPolicyAttributes{

			Name: G(DHCPOptionPolicyCont, "name"),

			Annotation: G(DHCPOptionPolicyCont, "annotation"),

			NameAlias: G(DHCPOptionPolicyCont, "nameAlias"),
		},
	}
}

func DHCPOptionPolicyFromContainer(cont *container.Container) *DHCPOptionPolicy {

	return DHCPOptionPolicyFromContainerList(cont, 0)
}

func DHCPOptionPolicyListFromContainer(cont *container.Container) []*DHCPOptionPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*DHCPOptionPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = DHCPOptionPolicyFromContainerList(cont, i)
	}

	return arr
}
