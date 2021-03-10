package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const DhcplblClassName = "dhcpLbl"

type BDDHCPLabel struct {
	BaseAttributes
	BDDHCPLabelAttributes
}

type BDDHCPLabelAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Owner string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewBDDHCPLabel(dhcpLblRn, parentDn, description string, dhcpLblattr BDDHCPLabelAttributes) *BDDHCPLabel {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpLblRn)
	return &BDDHCPLabel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcplblClassName,
			Rn:                dhcpLblRn,
		},

		BDDHCPLabelAttributes: dhcpLblattr,
	}
}

func (dhcpLbl *BDDHCPLabel) ToMap() (map[string]string, error) {
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

func BDDHCPLabelFromContainerList(cont *container.Container, index int) *BDDHCPLabel {

	BDDHCPLabelCont := cont.S("imdata").Index(index).S(DhcplblClassName, "attributes")
	return &BDDHCPLabel{
		BaseAttributes{
			DistinguishedName: G(BDDHCPLabelCont, "dn"),
			Description:       G(BDDHCPLabelCont, "descr"),
			Status:            G(BDDHCPLabelCont, "status"),
			ClassName:         DhcplblClassName,
			Rn:                G(BDDHCPLabelCont, "rn"),
		},

		BDDHCPLabelAttributes{

			Name: G(BDDHCPLabelCont, "name"),

			Annotation: G(BDDHCPLabelCont, "annotation"),

			NameAlias: G(BDDHCPLabelCont, "nameAlias"),

			Owner: G(BDDHCPLabelCont, "owner"),

			Tag: G(BDDHCPLabelCont, "tag"),
		},
	}
}

func BDDHCPLabelFromContainer(cont *container.Container) *BDDHCPLabel {

	return BDDHCPLabelFromContainerList(cont, 0)
}

func BDDHCPLabelListFromContainer(cont *container.Container) []*BDDHCPLabel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BDDHCPLabel, length)

	for i := 0; i < length; i++ {

		arr[i] = BDDHCPLabelFromContainerList(cont, i)
	}

	return arr
}
