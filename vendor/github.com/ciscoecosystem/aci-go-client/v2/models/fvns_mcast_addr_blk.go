package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnfvnsMcastAddrBlk        = "uni/infra/maddrns-%s/fromaddr-[%s]-toaddr-[%s]"
	RnfvnsMcastAddrBlk        = "fromaddr-[%s]-toaddr-[%s]"
	ParentDnfvnsMcastAddrBlk  = "uni/infra/maddrns-%s"
	FvnsmcastaddrblkClassName = "fvnsMcastAddrBlk"
)

type MulticastAddressBlock struct {
	BaseAttributes
	MulticastAddressBlockAttributes
}

type MulticastAddressBlockAttributes struct {
	Annotation string `json:",omitempty"`
	From       string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	To         string `json:",omitempty"`
}

func NewMulticastAddressBlock(fvnsMcastAddrBlkRn, parentDn, description string, fvnsMcastAddrBlkAttr MulticastAddressBlockAttributes) *MulticastAddressBlock {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsMcastAddrBlkRn)
	return &MulticastAddressBlock{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsmcastaddrblkClassName,
			Rn:                fvnsMcastAddrBlkRn,
		},
		MulticastAddressBlockAttributes: fvnsMcastAddrBlkAttr,
	}
}

func (fvnsMcastAddrBlk *MulticastAddressBlock) ToMap() (map[string]string, error) {
	fvnsMcastAddrBlkMap, err := fvnsMcastAddrBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsMcastAddrBlkMap, "annotation", fvnsMcastAddrBlk.Annotation)
	A(fvnsMcastAddrBlkMap, "from", fvnsMcastAddrBlk.From)
	A(fvnsMcastAddrBlkMap, "name", fvnsMcastAddrBlk.Name)
	A(fvnsMcastAddrBlkMap, "nameAlias", fvnsMcastAddrBlk.NameAlias)
	A(fvnsMcastAddrBlkMap, "to", fvnsMcastAddrBlk.To)
	return fvnsMcastAddrBlkMap, err
}

func MulticastAddressBlockFromContainerList(cont *container.Container, index int) *MulticastAddressBlock {
	MulticastAddressBlockCont := cont.S("imdata").Index(index).S(FvnsmcastaddrblkClassName, "attributes")
	return &MulticastAddressBlock{
		BaseAttributes{
			DistinguishedName: G(MulticastAddressBlockCont, "dn"),
			Description:       G(MulticastAddressBlockCont, "descr"),
			Status:            G(MulticastAddressBlockCont, "status"),
			ClassName:         FvnsmcastaddrblkClassName,
			Rn:                G(MulticastAddressBlockCont, "rn"),
		},
		MulticastAddressBlockAttributes{
			Annotation: G(MulticastAddressBlockCont, "annotation"),
			From:       G(MulticastAddressBlockCont, "from"),
			Name:       G(MulticastAddressBlockCont, "name"),
			NameAlias:  G(MulticastAddressBlockCont, "nameAlias"),
			To:         G(MulticastAddressBlockCont, "to"),
		},
	}
}

func MulticastAddressBlockFromContainer(cont *container.Container) *MulticastAddressBlock {
	return MulticastAddressBlockFromContainerList(cont, 0)
}

func MulticastAddressBlockListFromContainer(cont *container.Container) []*MulticastAddressBlock {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MulticastAddressBlock, length)

	for i := 0; i < length; i++ {
		arr[i] = MulticastAddressBlockFromContainerList(cont, i)
	}

	return arr
}
