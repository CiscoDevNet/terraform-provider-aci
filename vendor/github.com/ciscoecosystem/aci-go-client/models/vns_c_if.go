package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvnsCIf        = "%s/cIf-[%s]"
	RnvnsCIf        = "cIf-[%s]"
	VnscifClassName = "vnsCIf"
)

type ConcreteInterface struct {
	BaseAttributes
	NameAliasAttribute
	ConcreteInterfaceAttributes
}

type ConcreteInterfaceAttributes struct {
	Encap      string `json:",omitempty"`
	Name       string `json:",omitempty"`
	VnicName   string `json:",omitempty"`
	Annotation string `json:",omitempty"`
}

func NewConcreteInterface(vnsCIfRn, parentDn, nameAlias string, vnsCIfAttr ConcreteInterfaceAttributes) *ConcreteInterface {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsCIfRn)
	return &ConcreteInterface{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VnscifClassName,
			Rn:                vnsCIfRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ConcreteInterfaceAttributes: vnsCIfAttr,
	}
}

func (vnsCIf *ConcreteInterface) ToMap() (map[string]string, error) {
	vnsCIfMap, err := vnsCIf.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsCIf.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsCIfMap, key, value)
	}

	A(vnsCIfMap, "encap", vnsCIf.Encap)
	A(vnsCIfMap, "name", vnsCIf.Name)
	A(vnsCIfMap, "vnicName", vnsCIf.VnicName)
	A(vnsCIfMap, "annotation", vnsCIf.Annotation)
	return vnsCIfMap, err
}

func ConcreteInterfaceFromContainerList(cont *container.Container, index int) *ConcreteInterface {
	ConcreteInterfaceCont := cont.S("imdata").Index(index).S(VnscifClassName, "attributes")
	return &ConcreteInterface{
		BaseAttributes{
			DistinguishedName: G(ConcreteInterfaceCont, "dn"),
			Status:            G(ConcreteInterfaceCont, "status"),
			ClassName:         VnscifClassName,
			Rn:                G(ConcreteInterfaceCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ConcreteInterfaceCont, "nameAlias"),
		},
		ConcreteInterfaceAttributes{
			Encap:      G(ConcreteInterfaceCont, "encap"),
			Name:       G(ConcreteInterfaceCont, "name"),
			VnicName:   G(ConcreteInterfaceCont, "vnicName"),
			Annotation: G(ConcreteInterfaceCont, "annotation"),
		},
	}
}

func ConcreteInterfaceFromContainer(cont *container.Container) *ConcreteInterface {
	return ConcreteInterfaceFromContainerList(cont, 0)
}

func ConcreteInterfaceListFromContainer(cont *container.Container) []*ConcreteInterface {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ConcreteInterface, length)

	for i := 0; i < length; i++ {
		arr[i] = ConcreteInterfaceFromContainerList(cont, i)
	}

	return arr
}
