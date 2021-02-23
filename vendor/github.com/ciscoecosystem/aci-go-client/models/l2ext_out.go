package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L2extoutClassName = "l2extOut"

type L2Outside struct {
	BaseAttributes
	L2OutsideAttributes
}

type L2OutsideAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`
}

func NewL2Outside(l2extOutRn, parentDn, description string, l2extOutattr L2OutsideAttributes) *L2Outside {
	dn := fmt.Sprintf("%s/%s", parentDn, l2extOutRn)
	return &L2Outside{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L2extoutClassName,
			Rn:                l2extOutRn,
		},

		L2OutsideAttributes: l2extOutattr,
	}
}

func (l2extOut *L2Outside) ToMap() (map[string]string, error) {
	l2extOutMap, err := l2extOut.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l2extOutMap, "name", l2extOut.Name)

	A(l2extOutMap, "annotation", l2extOut.Annotation)

	A(l2extOutMap, "nameAlias", l2extOut.NameAlias)

	A(l2extOutMap, "targetDscp", l2extOut.TargetDscp)

	return l2extOutMap, err
}

func L2OutsideFromContainerList(cont *container.Container, index int) *L2Outside {

	L2OutsideCont := cont.S("imdata").Index(index).S(L2extoutClassName, "attributes")
	return &L2Outside{
		BaseAttributes{
			DistinguishedName: G(L2OutsideCont, "dn"),
			Description:       G(L2OutsideCont, "descr"),
			Status:            G(L2OutsideCont, "status"),
			ClassName:         L2extoutClassName,
			Rn:                G(L2OutsideCont, "rn"),
		},

		L2OutsideAttributes{

			Name: G(L2OutsideCont, "name"),

			Annotation: G(L2OutsideCont, "annotation"),

			NameAlias: G(L2OutsideCont, "nameAlias"),

			TargetDscp: G(L2OutsideCont, "targetDscp"),
		},
	}
}

func L2OutsideFromContainer(cont *container.Container) *L2Outside {

	return L2OutsideFromContainerList(cont, 0)
}

func L2OutsideListFromContainer(cont *container.Container) []*L2Outside {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L2Outside, length)

	for i := 0; i < length; i++ {

		arr[i] = L2OutsideFromContainerList(cont, i)
	}

	return arr
}
