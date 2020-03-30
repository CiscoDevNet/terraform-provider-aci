package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VztabooClassName = "vzTaboo"

type TabooContract struct {
	BaseAttributes
	TabooContractAttributes
}

type TabooContractAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewTabooContract(vzTabooRn, parentDn, description string, vzTabooattr TabooContractAttributes) *TabooContract {
	dn := fmt.Sprintf("%s/%s", parentDn, vzTabooRn)
	return &TabooContract{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VztabooClassName,
			Rn:                vzTabooRn,
		},

		TabooContractAttributes: vzTabooattr,
	}
}

func (vzTaboo *TabooContract) ToMap() (map[string]string, error) {
	vzTabooMap, err := vzTaboo.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vzTabooMap, "name", vzTaboo.Name)

	A(vzTabooMap, "annotation", vzTaboo.Annotation)

	A(vzTabooMap, "nameAlias", vzTaboo.NameAlias)

	return vzTabooMap, err
}

func TabooContractFromContainerList(cont *container.Container, index int) *TabooContract {

	TabooContractCont := cont.S("imdata").Index(index).S(VztabooClassName, "attributes")
	return &TabooContract{
		BaseAttributes{
			DistinguishedName: G(TabooContractCont, "dn"),
			Description:       G(TabooContractCont, "descr"),
			Status:            G(TabooContractCont, "status"),
			ClassName:         VztabooClassName,
			Rn:                G(TabooContractCont, "rn"),
		},

		TabooContractAttributes{

			Name: G(TabooContractCont, "name"),

			Annotation: G(TabooContractCont, "annotation"),

			NameAlias: G(TabooContractCont, "nameAlias"),
		},
	}
}

func TabooContractFromContainer(cont *container.Container) *TabooContract {

	return TabooContractFromContainerList(cont, 0)
}

func TabooContractListFromContainer(cont *container.Container) []*TabooContract {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*TabooContract, length)

	for i := 0; i < length; i++ {

		arr[i] = TabooContractFromContainerList(cont, i)
	}

	return arr
}
