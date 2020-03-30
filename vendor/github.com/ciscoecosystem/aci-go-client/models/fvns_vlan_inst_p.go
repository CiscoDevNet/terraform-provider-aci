package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvnsvlaninstpClassName = "fvnsVlanInstP"

type VLANPool struct {
	BaseAttributes
	VLANPoolAttributes
}

type VLANPoolAttributes struct {
	Name string `json:",omitempty"`

	AllocMode string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewVLANPool(fvnsVlanInstPRn, parentDn, description string, fvnsVlanInstPattr VLANPoolAttributes) *VLANPool {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsVlanInstPRn)
	return &VLANPool{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsvlaninstpClassName,
			//Rn:                fvnsVlanInstPRn,
		},

		VLANPoolAttributes: fvnsVlanInstPattr,
	}
}

func (fvnsVlanInstP *VLANPool) ToMap() (map[string]string, error) {
	fvnsVlanInstPMap, err := fvnsVlanInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsVlanInstPMap, "name", fvnsVlanInstP.Name)

	A(fvnsVlanInstPMap, "allocMode", fvnsVlanInstP.AllocMode)

	A(fvnsVlanInstPMap, "annotation", fvnsVlanInstP.Annotation)

	A(fvnsVlanInstPMap, "nameAlias", fvnsVlanInstP.NameAlias)

	return fvnsVlanInstPMap, err
}

func VLANPoolFromContainerList(cont *container.Container, index int) *VLANPool {

	VLANPoolCont := cont.S("imdata").Index(index).S(FvnsvlaninstpClassName, "attributes")
	return &VLANPool{
		BaseAttributes{
			DistinguishedName: G(VLANPoolCont, "dn"),
			Description:       G(VLANPoolCont, "descr"),
			Status:            G(VLANPoolCont, "status"),
			ClassName:         FvnsvlaninstpClassName,
			//Rn:                G(VLANPoolCont, "rn"),
		},

		VLANPoolAttributes{

			Name: G(VLANPoolCont, "name"),

			AllocMode: G(VLANPoolCont, "allocMode"),

			Annotation: G(VLANPoolCont, "annotation"),

			NameAlias: G(VLANPoolCont, "nameAlias"),
		},
	}
}

func VLANPoolFromContainer(cont *container.Container) *VLANPool {

	return VLANPoolFromContainerList(cont, 0)
}

func VLANPoolListFromContainer(cont *container.Container) []*VLANPool {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VLANPool, length)

	for i := 0; i < length; i++ {

		arr[i] = VLANPoolFromContainerList(cont, i)
	}

	return arr
}
