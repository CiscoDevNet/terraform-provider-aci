package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvrsprovClassName = "fvRsProv"

type ContractProvider struct {
	BaseAttributes
	ContractProviderAttributes
}

type ContractProviderAttributes struct {
	TnVzBrCPName string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	MatchT string `json:",omitempty"`

	Prio string `json:",omitempty"`

	TDn string `json:",omitempty"`
}

func NewContractProvider(fvRsProvRn, parentDn string, fvRsProvattr ContractProviderAttributes) *ContractProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsProvRn)
	return &ContractProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvrsprovClassName,
			Rn:                fvRsProvRn,
		},

		ContractProviderAttributes: fvRsProvattr,
	}
}

func (fvRsProv *ContractProvider) ToMap() (map[string]string, error) {
	fvRsProvMap, err := fvRsProv.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvRsProvMap, "tnVzBrCPName", fvRsProv.TnVzBrCPName)

	A(fvRsProvMap, "annotation", fvRsProv.Annotation)

	A(fvRsProvMap, "matchT", fvRsProv.MatchT)

	A(fvRsProvMap, "prio", fvRsProv.Prio)

	A(fvRsProvMap, "tDn", fvRsProv.TDn)

	return fvRsProvMap, err
}

func ContractProviderFromContainerList(cont *container.Container, index int) *ContractProvider {

	ContractProviderCont := cont.S("imdata").Index(index).S(FvrsprovClassName, "attributes")
	return &ContractProvider{
		BaseAttributes{
			DistinguishedName: G(ContractProviderCont, "dn"),
			Status:            G(ContractProviderCont, "status"),
			ClassName:         FvrsprovClassName,
			Rn:                G(ContractProviderCont, "rn"),
		},

		ContractProviderAttributes{

			TnVzBrCPName: G(ContractProviderCont, "tnVzBrCPName"),

			Annotation: G(ContractProviderCont, "annotation"),

			MatchT: G(ContractProviderCont, "matchT"),

			Prio: G(ContractProviderCont, "prio"),

			TDn: G(ContractProviderCont, "tDn"),
		},
	}
}

func ContractProviderFromContainer(cont *container.Container) *ContractProvider {

	return ContractProviderFromContainerList(cont, 0)
}

func ContractProviderListFromContainer(cont *container.Container) []*ContractProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ContractProvider, length)

	for i := 0; i < length; i++ {

		arr[i] = ContractProviderFromContainerList(cont, i)
	}

	return arr
}
