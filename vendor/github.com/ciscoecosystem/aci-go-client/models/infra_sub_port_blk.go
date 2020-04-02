package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrasubportblkClassName = "infraSubPortBlk"

type AccessSubPortBlock struct {
	BaseAttributes
	AccessSubPortBlockAttributes
}

type AccessSubPortBlockAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	FromCard string `json:",omitempty"`

	FromPort string `json:",omitempty"`

	FromSubPort string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	ToCard string `json:",omitempty"`

	ToPort string `json:",omitempty"`

	ToSubPort string `json:",omitempty"`
}

func NewAccessSubPortBlock(infraSubPortBlkRn, parentDn, description string, infraSubPortBlkattr AccessSubPortBlockAttributes) *AccessSubPortBlock {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSubPortBlkRn)
	return &AccessSubPortBlock{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrasubportblkClassName,
			Rn:                infraSubPortBlkRn,
		},

		AccessSubPortBlockAttributes: infraSubPortBlkattr,
	}
}

func (infraSubPortBlk *AccessSubPortBlock) ToMap() (map[string]string, error) {
	infraSubPortBlkMap, err := infraSubPortBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraSubPortBlkMap, "name", infraSubPortBlk.Name)

	A(infraSubPortBlkMap, "annotation", infraSubPortBlk.Annotation)

	A(infraSubPortBlkMap, "fromCard", infraSubPortBlk.FromCard)

	A(infraSubPortBlkMap, "fromPort", infraSubPortBlk.FromPort)

	A(infraSubPortBlkMap, "fromSubPort", infraSubPortBlk.FromSubPort)

	A(infraSubPortBlkMap, "nameAlias", infraSubPortBlk.NameAlias)

	A(infraSubPortBlkMap, "toCard", infraSubPortBlk.ToCard)

	A(infraSubPortBlkMap, "toPort", infraSubPortBlk.ToPort)

	A(infraSubPortBlkMap, "toSubPort", infraSubPortBlk.ToSubPort)

	return infraSubPortBlkMap, err
}

func AccessSubPortBlockFromContainerList(cont *container.Container, index int) *AccessSubPortBlock {

	AccessSubPortBlockCont := cont.S("imdata").Index(index).S(InfrasubportblkClassName, "attributes")
	return &AccessSubPortBlock{
		BaseAttributes{
			DistinguishedName: G(AccessSubPortBlockCont, "dn"),
			Description:       G(AccessSubPortBlockCont, "descr"),
			Status:            G(AccessSubPortBlockCont, "status"),
			ClassName:         InfrasubportblkClassName,
			Rn:                G(AccessSubPortBlockCont, "rn"),
		},

		AccessSubPortBlockAttributes{

			Name: G(AccessSubPortBlockCont, "name"),

			Annotation: G(AccessSubPortBlockCont, "annotation"),

			FromCard: G(AccessSubPortBlockCont, "fromCard"),

			FromPort: G(AccessSubPortBlockCont, "fromPort"),

			FromSubPort: G(AccessSubPortBlockCont, "fromSubPort"),

			NameAlias: G(AccessSubPortBlockCont, "nameAlias"),

			ToCard: G(AccessSubPortBlockCont, "toCard"),

			ToPort: G(AccessSubPortBlockCont, "toPort"),

			ToSubPort: G(AccessSubPortBlockCont, "toSubPort"),
		},
	}
}

func AccessSubPortBlockFromContainer(cont *container.Container) *AccessSubPortBlock {

	return AccessSubPortBlockFromContainerList(cont, 0)
}

func AccessSubPortBlockListFromContainer(cont *container.Container) []*AccessSubPortBlock {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AccessSubPortBlock, length)

	for i := 0; i < length; i++ {

		arr[i] = AccessSubPortBlockFromContainerList(cont, i)
	}

	return arr
}
