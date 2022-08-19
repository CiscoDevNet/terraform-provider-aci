package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvnsMcastAddrInstP        = "uni/infra/maddrns-%s"
	RnfvnsMcastAddrInstP        = "maddrns-%s"
	ParentDnfvnsMcastAddrInstP  = "uni/infra"
	FvnsmcastaddrinstpClassName = "fvnsMcastAddrInstP"
)

type MulticastAddressPool struct {
	BaseAttributes
	MulticastAddressPoolAttributes
}

type MulticastAddressPoolAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewMulticastAddressPool(fvnsMcastAddrInstPRn, parentDn, description string, fvnsMcastAddrInstPAttr MulticastAddressPoolAttributes) *MulticastAddressPool {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsMcastAddrInstPRn)
	return &MulticastAddressPool{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsmcastaddrinstpClassName,
			Rn:                fvnsMcastAddrInstPRn,
		},
		MulticastAddressPoolAttributes: fvnsMcastAddrInstPAttr,
	}
}

func (fvnsMcastAddrInstP *MulticastAddressPool) ToMap() (map[string]string, error) {
	fvnsMcastAddrInstPMap, err := fvnsMcastAddrInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsMcastAddrInstPMap, "annotation", fvnsMcastAddrInstP.Annotation)
	A(fvnsMcastAddrInstPMap, "name", fvnsMcastAddrInstP.Name)
	A(fvnsMcastAddrInstPMap, "nameAlias", fvnsMcastAddrInstP.NameAlias)
	return fvnsMcastAddrInstPMap, err
}

func MulticastAddressPoolFromContainerList(cont *container.Container, index int) *MulticastAddressPool {
	MulticastAddressPoolCont := cont.S("imdata").Index(index).S(FvnsmcastaddrinstpClassName, "attributes")
	return &MulticastAddressPool{
		BaseAttributes{
			DistinguishedName: G(MulticastAddressPoolCont, "dn"),
			Description:       G(MulticastAddressPoolCont, "descr"),
			Status:            G(MulticastAddressPoolCont, "status"),
			ClassName:         FvnsmcastaddrinstpClassName,
			Rn:                G(MulticastAddressPoolCont, "rn"),
		},
		MulticastAddressPoolAttributes{
			Annotation: G(MulticastAddressPoolCont, "annotation"),
			Name:       G(MulticastAddressPoolCont, "name"),
			NameAlias:  G(MulticastAddressPoolCont, "nameAlias"),
		},
	}
}

func MulticastAddressPoolFromContainer(cont *container.Container) *MulticastAddressPool {
	return MulticastAddressPoolFromContainerList(cont, 0)
}

func MulticastAddressPoolListFromContainer(cont *container.Container) []*MulticastAddressPool {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MulticastAddressPool, length)

	for i := 0; i < length; i++ {
		arr[i] = MulticastAddressPoolFromContainerList(cont, i)
	}

	return arr
}
