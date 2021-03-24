package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const HsrpifpolClassName = "hsrpIfPol"

type HSRPInterfacePolicy struct {
	BaseAttributes
	HSRPInterfacePolicyAttributes
}

type HSRPInterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	Delay string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	ReloadDelay string `json:",omitempty"`
}

func NewHSRPInterfacePolicy(hsrpIfPolRn, parentDn, description string, hsrpIfPolattr HSRPInterfacePolicyAttributes) *HSRPInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, hsrpIfPolRn)
	return &HSRPInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         HsrpifpolClassName,
			Rn:                hsrpIfPolRn,
		},

		HSRPInterfacePolicyAttributes: hsrpIfPolattr,
	}
}

func (hsrpIfPol *HSRPInterfacePolicy) ToMap() (map[string]string, error) {
	hsrpIfPolMap, err := hsrpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(hsrpIfPolMap, "name", hsrpIfPol.Name)

	A(hsrpIfPolMap, "annotation", hsrpIfPol.Annotation)

	A(hsrpIfPolMap, "ctrl", hsrpIfPol.Ctrl)

	A(hsrpIfPolMap, "delay", hsrpIfPol.Delay)

	A(hsrpIfPolMap, "nameAlias", hsrpIfPol.NameAlias)

	A(hsrpIfPolMap, "reloadDelay", hsrpIfPol.ReloadDelay)

	return hsrpIfPolMap, err
}

func HSRPInterfacePolicyFromContainerList(cont *container.Container, index int) *HSRPInterfacePolicy {

	HSRPInterfacePolicyCont := cont.S("imdata").Index(index).S(HsrpifpolClassName, "attributes")
	return &HSRPInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(HSRPInterfacePolicyCont, "dn"),
			Description:       G(HSRPInterfacePolicyCont, "descr"),
			Status:            G(HSRPInterfacePolicyCont, "status"),
			ClassName:         HsrpifpolClassName,
			Rn:                G(HSRPInterfacePolicyCont, "rn"),
		},

		HSRPInterfacePolicyAttributes{

			Name: G(HSRPInterfacePolicyCont, "name"),

			Annotation: G(HSRPInterfacePolicyCont, "annotation"),

			Ctrl: G(HSRPInterfacePolicyCont, "ctrl"),

			Delay: G(HSRPInterfacePolicyCont, "delay"),

			NameAlias: G(HSRPInterfacePolicyCont, "nameAlias"),

			ReloadDelay: G(HSRPInterfacePolicyCont, "reloadDelay"),
		},
	}
}

func HSRPInterfacePolicyFromContainer(cont *container.Container) *HSRPInterfacePolicy {

	return HSRPInterfacePolicyFromContainerList(cont, 0)
}

func HSRPInterfacePolicyListFromContainer(cont *container.Container) []*HSRPInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*HSRPInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = HSRPInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}
