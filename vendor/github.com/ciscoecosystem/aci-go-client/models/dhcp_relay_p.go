package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const DhcprelaypClassName = "dhcpRelayP"

type DHCPRelayPolicy struct {
	BaseAttributes
	DHCPRelayPolicyAttributes
}

type DHCPRelayPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Mode string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Owner string `json:",omitempty"`
}

func NewDHCPRelayPolicy(dhcpRelayPRn, parentDn, description string, dhcpRelayPattr DHCPRelayPolicyAttributes) *DHCPRelayPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpRelayPRn)
	return &DHCPRelayPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcprelaypClassName,
			Rn:                dhcpRelayPRn,
		},

		DHCPRelayPolicyAttributes: dhcpRelayPattr,
	}
}

func (dhcpRelayP *DHCPRelayPolicy) ToMap() (map[string]string, error) {
	dhcpRelayPMap, err := dhcpRelayP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpRelayPMap, "name", dhcpRelayP.Name)

	A(dhcpRelayPMap, "annotation", dhcpRelayP.Annotation)

	A(dhcpRelayPMap, "mode", dhcpRelayP.Mode)

	A(dhcpRelayPMap, "nameAlias", dhcpRelayP.NameAlias)

	A(dhcpRelayPMap, "owner", dhcpRelayP.Owner)

	return dhcpRelayPMap, err
}

func DHCPRelayPolicyFromContainerList(cont *container.Container, index int) *DHCPRelayPolicy {

	DHCPRelayPolicyCont := cont.S("imdata").Index(index).S(DhcprelaypClassName, "attributes")
	return &DHCPRelayPolicy{
		BaseAttributes{
			DistinguishedName: G(DHCPRelayPolicyCont, "dn"),
			Description:       G(DHCPRelayPolicyCont, "descr"),
			Status:            G(DHCPRelayPolicyCont, "status"),
			ClassName:         DhcprelaypClassName,
			Rn:                G(DHCPRelayPolicyCont, "rn"),
		},

		DHCPRelayPolicyAttributes{

			Name: G(DHCPRelayPolicyCont, "name"),

			Annotation: G(DHCPRelayPolicyCont, "annotation"),

			Mode: G(DHCPRelayPolicyCont, "mode"),

			NameAlias: G(DHCPRelayPolicyCont, "nameAlias"),

			Owner: G(DHCPRelayPolicyCont, "owner"),
		},
	}
}

func DHCPRelayPolicyFromContainer(cont *container.Container) *DHCPRelayPolicy {

	return DHCPRelayPolicyFromContainerList(cont, 0)
}

func DHCPRelayPolicyListFromContainer(cont *container.Container) []*DHCPRelayPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*DHCPRelayPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = DHCPRelayPolicyFromContainerList(cont, i)
	}

	return arr
}
