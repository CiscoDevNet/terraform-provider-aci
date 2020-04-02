package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvnsvxlaninstpClassName = "fvnsVxlanInstP"

type VXLANPool struct {
	BaseAttributes
	VXLANPoolAttributes
}

type VXLANPoolAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewVXLANPool(fvnsVxlanInstPRn, parentDn, description string, fvnsVxlanInstPattr VXLANPoolAttributes) *VXLANPool {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsVxlanInstPRn)
	return &VXLANPool{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsvxlaninstpClassName,
			Rn:                fvnsVxlanInstPRn,
		},

		VXLANPoolAttributes: fvnsVxlanInstPattr,
	}
}

func (fvnsVxlanInstP *VXLANPool) ToMap() (map[string]string, error) {
	fvnsVxlanInstPMap, err := fvnsVxlanInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsVxlanInstPMap, "name", fvnsVxlanInstP.Name)

	A(fvnsVxlanInstPMap, "annotation", fvnsVxlanInstP.Annotation)

	A(fvnsVxlanInstPMap, "nameAlias", fvnsVxlanInstP.NameAlias)

	return fvnsVxlanInstPMap, err
}

func VXLANPoolFromContainerList(cont *container.Container, index int) *VXLANPool {

	VXLANPoolCont := cont.S("imdata").Index(index).S(FvnsvxlaninstpClassName, "attributes")
	return &VXLANPool{
		BaseAttributes{
			DistinguishedName: G(VXLANPoolCont, "dn"),
			Description:       G(VXLANPoolCont, "descr"),
			Status:            G(VXLANPoolCont, "status"),
			ClassName:         FvnsvxlaninstpClassName,
			Rn:                G(VXLANPoolCont, "rn"),
		},

		VXLANPoolAttributes{

			Name: G(VXLANPoolCont, "name"),

			Annotation: G(VXLANPoolCont, "annotation"),

			NameAlias: G(VXLANPoolCont, "nameAlias"),
		},
	}
}

func VXLANPoolFromContainer(cont *container.Container) *VXLANPool {

	return VXLANPoolFromContainerList(cont, 0)
}

func VXLANPoolListFromContainer(cont *container.Container) []*VXLANPool {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VXLANPool, length)

	for i := 0; i < length; i++ {

		arr[i] = VXLANPoolFromContainerList(cont, i)
	}

	return arr
}
