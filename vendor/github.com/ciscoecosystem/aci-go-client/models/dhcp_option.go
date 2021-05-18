package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const DhcpoptionClassName = "dhcpOption"

type DHCPOption struct {
	BaseAttributes
	DHCPOptionAttributes
}

type DHCPOptionAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Data string `json:",omitempty"`

	DHCPOption_id string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewDHCPOption(dhcpOptionRn, parentDn string, dhcpOptionattr DHCPOptionAttributes) *DHCPOption {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpOptionRn)
	return &DHCPOption{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         DhcpoptionClassName,
			Rn:                dhcpOptionRn,
		},

		DHCPOptionAttributes: dhcpOptionattr,
	}
}

func (dhcpOption *DHCPOption) ToMap() (map[string]string, error) {
	dhcpOptionMap, err := dhcpOption.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpOptionMap, "name", dhcpOption.Name)

	A(dhcpOptionMap, "annotation", dhcpOption.Annotation)

	A(dhcpOptionMap, "data", dhcpOption.Data)

	A(dhcpOptionMap, "id", dhcpOption.DHCPOption_id)

	A(dhcpOptionMap, "nameAlias", dhcpOption.NameAlias)

	return dhcpOptionMap, err
}

func DHCPOptionFromContainerList(cont *container.Container, index int) *DHCPOption {

	DHCPOptionCont := cont.S("imdata").Index(index).S(DhcpoptionClassName, "attributes")
	return &DHCPOption{
		BaseAttributes{
			DistinguishedName: G(DHCPOptionCont, "dn"),
			Status:            G(DHCPOptionCont, "status"),
			ClassName:         DhcpoptionClassName,
			Rn:                G(DHCPOptionCont, "rn"),
		},

		DHCPOptionAttributes{

			Name: G(DHCPOptionCont, "name"),

			Annotation: G(DHCPOptionCont, "annotation"),

			Data: G(DHCPOptionCont, "data"),

			DHCPOption_id: G(DHCPOptionCont, "id"),

			NameAlias: G(DHCPOptionCont, "nameAlias"),
		},
	}
}

func DHCPOptionFromContainer(cont *container.Container) *DHCPOption {

	return DHCPOptionFromContainerList(cont, 0)
}

func DHCPOptionListFromContainer(cont *container.Container) []*DHCPOption {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*DHCPOption, length)

	for i := 0; i < length; i++ {

		arr[i] = DHCPOptionFromContainerList(cont, i)
	}

	return arr
}
