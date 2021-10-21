package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvrsconsClassName = "fvRsCons"

type ContractConsumer struct {
	BaseAttributes
	ContractConsumerAttributes
}

type ContractConsumerAttributes struct {
	TnVzBrCPName string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Prio string `json:",omitempty"`

	TDn string `json:",omitempty"`
}

func NewContractConsumer(fvRsConsRn, parentDn string, fvRsConsattr ContractConsumerAttributes) *ContractConsumer {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsConsRn)
	return &ContractConsumer{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvrsconsClassName,
			Rn:                fvRsConsRn,
		},

		ContractConsumerAttributes: fvRsConsattr,
	}
}

func (fvRsCons *ContractConsumer) ToMap() (map[string]string, error) {
	fvRsConsMap, err := fvRsCons.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvRsConsMap, "tnVzBrCPName", fvRsCons.TnVzBrCPName)

	A(fvRsConsMap, "annotation", fvRsCons.Annotation)

	A(fvRsConsMap, "prio", fvRsCons.Prio)

	A(fvRsConsMap, "tDn", fvRsCons.TDn)

	return fvRsConsMap, err
}

func ContractConsumerFromContainerList(cont *container.Container, index int) *ContractConsumer {

	ContractConsumerCont := cont.S("imdata").Index(index).S(FvrsconsClassName, "attributes")
	return &ContractConsumer{
		BaseAttributes{
			DistinguishedName: G(ContractConsumerCont, "dn"),
			Status:            G(ContractConsumerCont, "status"),
			ClassName:         FvrsconsClassName,
			Rn:                G(ContractConsumerCont, "rn"),
		},

		ContractConsumerAttributes{

			TnVzBrCPName: G(ContractConsumerCont, "tnVzBrCPName"),

			Annotation: G(ContractConsumerCont, "annotation"),

			Prio: G(ContractConsumerCont, "prio"),

			TDn: G(ContractConsumerCont, "tDn"),
		},
	}
}

func ContractConsumerFromContainer(cont *container.Container) *ContractConsumer {

	return ContractConsumerFromContainerList(cont, 0)
}

func ContractConsumerListFromContainer(cont *container.Container) []*ContractConsumer {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ContractConsumer, length)

	for i := 0; i < length; i++ {

		arr[i] = ContractConsumerFromContainerList(cont, i)
	}

	return arr
}
