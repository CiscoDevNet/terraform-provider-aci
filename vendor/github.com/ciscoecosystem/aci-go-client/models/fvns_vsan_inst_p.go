package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvnsvsaninstpClassName = "fvnsVsanInstP"

type VSANPool struct {
	BaseAttributes
	VSANPoolAttributes
}

type VSANPoolAttributes struct {
	Name string `json:",omitempty"`

	AllocMode string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewVSANPool(fvnsVsanInstPRn, parentDn, description string, fvnsVsanInstPattr VSANPoolAttributes) *VSANPool {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsVsanInstPRn)
	return &VSANPool{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsvsaninstpClassName,
			//Rn:                fvnsVsanInstPRn,
		},

		VSANPoolAttributes: fvnsVsanInstPattr,
	}
}

func (fvnsVsanInstP *VSANPool) ToMap() (map[string]string, error) {
	fvnsVsanInstPMap, err := fvnsVsanInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsVsanInstPMap, "name", fvnsVsanInstP.Name)

	A(fvnsVsanInstPMap, "allocMode", fvnsVsanInstP.AllocMode)

	A(fvnsVsanInstPMap, "annotation", fvnsVsanInstP.Annotation)

	A(fvnsVsanInstPMap, "nameAlias", fvnsVsanInstP.NameAlias)

	return fvnsVsanInstPMap, err
}

func VSANPoolFromContainerList(cont *container.Container, index int) *VSANPool {

	VSANPoolCont := cont.S("imdata").Index(index).S(FvnsvsaninstpClassName, "attributes")
	return &VSANPool{
		BaseAttributes{
			DistinguishedName: G(VSANPoolCont, "dn"),
			Description:       G(VSANPoolCont, "descr"),
			Status:            G(VSANPoolCont, "status"),
			ClassName:         FvnsvsaninstpClassName,
			//Rn:                G(VSANPoolCont, "rn"),
		},

		VSANPoolAttributes{

			Name: G(VSANPoolCont, "name"),

			AllocMode: G(VSANPoolCont, "allocMode"),

			Annotation: G(VSANPoolCont, "annotation"),

			NameAlias: G(VSANPoolCont, "nameAlias"),
		},
	}
}

func VSANPoolFromContainer(cont *container.Container) *VSANPool {

	return VSANPoolFromContainerList(cont, 0)
}

func VSANPoolListFromContainer(cont *container.Container) []*VSANPool {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VSANPool, length)

	for i := 0; i < length; i++ {

		arr[i] = VSANPoolFromContainerList(cont, i)
	}

	return arr
}
