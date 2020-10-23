package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabrichifpolClassName = "fabricHIfPol"

type LinkLevelPolicy struct {
	BaseAttributes
	LinkLevelPolicyAttributes
}

type LinkLevelPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	AutoNeg string `json:",omitempty"`

	FecMode string `json:",omitempty"`

	LinkDebounce string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Speed string `json:",omitempty"`
}

func NewLinkLevelPolicy(fabricHIfPolRn, parentDn, description string, fabricHIfPolattr LinkLevelPolicyAttributes) *LinkLevelPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricHIfPolRn)
	return &LinkLevelPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabrichifpolClassName,
			Rn:                fabricHIfPolRn,
		},

		LinkLevelPolicyAttributes: fabricHIfPolattr,
	}
}

func (fabricHIfPol *LinkLevelPolicy) ToMap() (map[string]string, error) {
	fabricHIfPolMap, err := fabricHIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricHIfPolMap, "name", fabricHIfPol.Name)

	A(fabricHIfPolMap, "annotation", fabricHIfPol.Annotation)

	A(fabricHIfPolMap, "autoNeg", fabricHIfPol.AutoNeg)

	A(fabricHIfPolMap, "fecMode", fabricHIfPol.FecMode)

	A(fabricHIfPolMap, "linkDebounce", fabricHIfPol.LinkDebounce)

	A(fabricHIfPolMap, "nameAlias", fabricHIfPol.NameAlias)

	A(fabricHIfPolMap, "speed", fabricHIfPol.Speed)

	return fabricHIfPolMap, err
}

func LinkLevelPolicyFromContainerList(cont *container.Container, index int) *LinkLevelPolicy {

	LinkLevelPolicyCont := cont.S("imdata").Index(index).S(FabrichifpolClassName, "attributes")
	return &LinkLevelPolicy{
		BaseAttributes{
			DistinguishedName: G(LinkLevelPolicyCont, "dn"),
			Description:       G(LinkLevelPolicyCont, "descr"),
			Status:            G(LinkLevelPolicyCont, "status"),
			ClassName:         FabrichifpolClassName,
			Rn:                G(LinkLevelPolicyCont, "rn"),
		},

		LinkLevelPolicyAttributes{

			Name: G(LinkLevelPolicyCont, "name"),

			Annotation: G(LinkLevelPolicyCont, "annotation"),

			AutoNeg: G(LinkLevelPolicyCont, "autoNeg"),

			FecMode: G(LinkLevelPolicyCont, "fecMode"),

			LinkDebounce: G(LinkLevelPolicyCont, "linkDebounce"),

			NameAlias: G(LinkLevelPolicyCont, "nameAlias"),

			Speed: G(LinkLevelPolicyCont, "speed"),
		},
	}
}

func LinkLevelPolicyFromContainer(cont *container.Container) *LinkLevelPolicy {

	return LinkLevelPolicyFromContainerList(cont, 0)
}

func LinkLevelPolicyListFromContainer(cont *container.Container) []*LinkLevelPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LinkLevelPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = LinkLevelPolicyFromContainerList(cont, i)
	}

	return arr
}
