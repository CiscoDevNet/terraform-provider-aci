package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VzbrcpClassName = "vzBrCP"

type Contract struct {
	BaseAttributes
	ContractAttributes
}

type ContractAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Prio string `json:",omitempty"`

	Scope string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`
}

func NewContract(vzBrCPRn, parentDn, description string, vzBrCPattr ContractAttributes) *Contract {
	dn := fmt.Sprintf("%s/%s", parentDn, vzBrCPRn)
	return &Contract{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzbrcpClassName,
			Rn:                vzBrCPRn,
		},

		ContractAttributes: vzBrCPattr,
	}
}

func (vzBrCP *Contract) ToMap() (map[string]string, error) {
	vzBrCPMap, err := vzBrCP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vzBrCPMap, "name", vzBrCP.Name)

	A(vzBrCPMap, "annotation", vzBrCP.Annotation)

	A(vzBrCPMap, "nameAlias", vzBrCP.NameAlias)

	A(vzBrCPMap, "prio", vzBrCP.Prio)

	A(vzBrCPMap, "scope", vzBrCP.Scope)

	A(vzBrCPMap, "targetDscp", vzBrCP.TargetDscp)

	return vzBrCPMap, err
}

func ContractFromContainerList(cont *container.Container, index int) *Contract {

	ContractCont := cont.S("imdata").Index(index).S(VzbrcpClassName, "attributes")
	return &Contract{
		BaseAttributes{
			DistinguishedName: G(ContractCont, "dn"),
			Description:       G(ContractCont, "descr"),
			Status:            G(ContractCont, "status"),
			ClassName:         VzbrcpClassName,
			Rn:                G(ContractCont, "rn"),
		},

		ContractAttributes{

			Name: G(ContractCont, "name"),

			Annotation: G(ContractCont, "annotation"),

			NameAlias: G(ContractCont, "nameAlias"),

			Prio: G(ContractCont, "prio"),

			Scope: G(ContractCont, "scope"),

			TargetDscp: G(ContractCont, "targetDscp"),
		},
	}
}

func ContractFromContainer(cont *container.Container) *Contract {

	return ContractFromContainerList(cont, 0)
}

func ContractListFromContainer(cont *container.Container) []*Contract {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Contract, length)

	for i := 0; i < length; i++ {

		arr[i] = ContractFromContainerList(cont, i)
	}

	return arr
}
