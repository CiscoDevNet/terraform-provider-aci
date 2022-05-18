package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvRsConsIf        = "uni/tn-%s/ap-%s/epg-%s/rsconsIf-%s"
	RnfvRsConsIf        = "rsconsIf-%s"
	ParentDnfvRsConsIf  = "uni/tn-%s/ap-%s/epg-%s"
	FvrsconsifClassName = "fvRsConsIf"
)

type ContractInterfaceRelationship struct {
	BaseAttributes
	ContractInterfaceRelationshipAttributes
}

type ContractInterfaceRelationshipAttributes struct {
	Annotation   string `json:",omitempty"`
	Prio         string `json:",omitempty"`
	TnVzCPIfName string `json:",omitempty"`
	tDn          string `json:",omitempty"`
}

func NewContractInterfaceRelationship(fvRsConsIfRn, parentDn string, fvRsConsIfAttr ContractInterfaceRelationshipAttributes) *ContractInterfaceRelationship {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsConsIfRn)
	return &ContractInterfaceRelationship{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvrsconsifClassName,
			Rn:                fvRsConsIfRn,
		},
		ContractInterfaceRelationshipAttributes: fvRsConsIfAttr,
	}
}

func (fvRsConsIf *ContractInterfaceRelationship) ToMap() (map[string]string, error) {
	fvRsConsIfMap, err := fvRsConsIf.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvRsConsIfMap, "prio", fvRsConsIf.Prio)
	A(fvRsConsIfMap, "tnVzCPIfName", fvRsConsIf.TnVzCPIfName)
	A(fvRsConsIfMap, "tDn", fvRsConsIf.tDn)
	A(fvRsConsIfMap, "annotation", fvRsConsIf.Annotation)
	return fvRsConsIfMap, err
}

func ContractInterfaceRelationshipFromContainerList(cont *container.Container, index int) *ContractInterfaceRelationship {
	ContractInterfaceRelationshipCont := cont.S("imdata").Index(index).S(FvrsconsifClassName, "attributes")
	return &ContractInterfaceRelationship{
		BaseAttributes{
			DistinguishedName: G(ContractInterfaceRelationshipCont, "dn"),
			Status:            G(ContractInterfaceRelationshipCont, "status"),
			ClassName:         FvrsconsifClassName,
			Rn:                G(ContractInterfaceRelationshipCont, "rn"),
		},
		ContractInterfaceRelationshipAttributes{
			Prio:         G(ContractInterfaceRelationshipCont, "prio"),
			TnVzCPIfName: G(ContractInterfaceRelationshipCont, "tnVzCPIfName"),
			tDn:          G(ContractInterfaceRelationshipCont, "tDn"),
			Annotation:   G(ContractInterfaceRelationshipCont, "annotation"),
		},
	}
}

func ContractInterfaceRelationshipFromContainer(cont *container.Container) *ContractInterfaceRelationship {
	return ContractInterfaceRelationshipFromContainerList(cont, 0)
}

func ContractInterfaceRelationshipListFromContainer(cont *container.Container) []*ContractInterfaceRelationship {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ContractInterfaceRelationship, length)

	for i := 0; i < length; i++ {
		arr[i] = ContractInterfaceRelationshipFromContainerList(cont, i)
	}

	return arr
}
