package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const DhcplblClassName = "dhcpLbl"

type DHCPRelayLabel struct {
	BaseAttributes
	DHCPRelayLabelAttributes
}

type DHCPRelayLabelAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Owner string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewDHCPRelayLabel(dhcpLblRn, parentDn, description string, dhcpLblattr DHCPRelayLabelAttributes) *DHCPRelayLabel {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpLblRn)
	return &DHCPRelayLabel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcplblClassName,
			Rn:                dhcpLblRn,
		},

		DHCPRelayLabelAttributes: dhcpLblattr,
	}
}

func (dhcpLbl *DHCPRelayLabel) ToMap() (map[string]string, error) {
	dhcpLblMap, err := dhcpLbl.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpLblMap, "name", dhcpLbl.Name)

	A(dhcpLblMap, "annotation", dhcpLbl.Annotation)

	A(dhcpLblMap, "nameAlias", dhcpLbl.NameAlias)

	A(dhcpLblMap, "owner", dhcpLbl.Owner)

	A(dhcpLblMap, "tag", dhcpLbl.Tag)

	return dhcpLblMap, err
}

func DHCPRelayLabelFromContainerList(cont *container.Container, index int) *DHCPRelayLabel {

	DHCPRelayLabelCont := cont.S("imdata").Index(index).S(DhcplblClassName, "attributes")
	return &DHCPRelayLabel{
		BaseAttributes{
			DistinguishedName: G(DHCPRelayLabelCont, "dn"),
			Description:       G(DHCPRelayLabelCont, "descr"),
			Status:            G(DHCPRelayLabelCont, "status"),
			ClassName:         DhcplblClassName,
			Rn:                G(DHCPRelayLabelCont, "rn"),
		},

		DHCPRelayLabelAttributes{

			Name: G(DHCPRelayLabelCont, "name"),

			Annotation: G(DHCPRelayLabelCont, "annotation"),

			NameAlias: G(DHCPRelayLabelCont, "nameAlias"),

			Owner: G(DHCPRelayLabelCont, "owner"),

			Tag: G(DHCPRelayLabelCont, "tag"),
		},
	}
}

func DHCPRelayLabelFromContainer(cont *container.Container) *DHCPRelayLabel {

	return DHCPRelayLabelFromContainerList(cont, 0)
}

func DHCPRelayLabelListFromContainer(cont *container.Container) []*DHCPRelayLabel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*DHCPRelayLabel, length)

	for i := 0; i < length; i++ {

		arr[i] = DHCPRelayLabelFromContainerList(cont, i)
	}

	return arr
}
