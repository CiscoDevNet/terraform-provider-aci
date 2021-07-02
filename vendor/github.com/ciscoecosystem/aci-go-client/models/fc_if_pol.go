package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FcifpolClassName = "fcIfPol"

type InterfaceFCPolicy struct {
	BaseAttributes
	InterfaceFCPolicyAttributes
}

type InterfaceFCPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Automaxspeed string `json:",omitempty"`

	FillPattern string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PortMode string `json:",omitempty"`

	RxBBCredit string `json:",omitempty"`

	Speed string `json:",omitempty"`

	TrunkMode string `json:",omitempty"`
}

func NewInterfaceFCPolicy(fcIfPolRn, parentDn, description string, fcIfPolattr InterfaceFCPolicyAttributes) *InterfaceFCPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, fcIfPolRn)
	return &InterfaceFCPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FcifpolClassName,
			Rn:                fcIfPolRn,
		},

		InterfaceFCPolicyAttributes: fcIfPolattr,
	}
}

func (fcIfPol *InterfaceFCPolicy) ToMap() (map[string]string, error) {
	fcIfPolMap, err := fcIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fcIfPolMap, "name", fcIfPol.Name)

	A(fcIfPolMap, "annotation", fcIfPol.Annotation)

	A(fcIfPolMap, "automaxspeed", fcIfPol.Automaxspeed)

	A(fcIfPolMap, "fillPattern", fcIfPol.FillPattern)

	A(fcIfPolMap, "nameAlias", fcIfPol.NameAlias)

	A(fcIfPolMap, "portMode", fcIfPol.PortMode)

	A(fcIfPolMap, "rxBBCredit", fcIfPol.RxBBCredit)

	A(fcIfPolMap, "speed", fcIfPol.Speed)

	A(fcIfPolMap, "trunkMode", fcIfPol.TrunkMode)

	return fcIfPolMap, err
}

func InterfaceFCPolicyFromContainerList(cont *container.Container, index int) *InterfaceFCPolicy {

	InterfaceFCPolicyCont := cont.S("imdata").Index(index).S(FcifpolClassName, "attributes")
	return &InterfaceFCPolicy{
		BaseAttributes{
			DistinguishedName: G(InterfaceFCPolicyCont, "dn"),
			Description:       G(InterfaceFCPolicyCont, "descr"),
			Status:            G(InterfaceFCPolicyCont, "status"),
			ClassName:         FcifpolClassName,
			Rn:                G(InterfaceFCPolicyCont, "rn"),
		},

		InterfaceFCPolicyAttributes{

			Name: G(InterfaceFCPolicyCont, "name"),

			Annotation: G(InterfaceFCPolicyCont, "annotation"),

			Automaxspeed: G(InterfaceFCPolicyCont, "automaxspeed"),

			FillPattern: G(InterfaceFCPolicyCont, "fillPattern"),

			NameAlias: G(InterfaceFCPolicyCont, "nameAlias"),

			PortMode: G(InterfaceFCPolicyCont, "portMode"),

			RxBBCredit: G(InterfaceFCPolicyCont, "rxBBCredit"),

			Speed: G(InterfaceFCPolicyCont, "speed"),

			TrunkMode: G(InterfaceFCPolicyCont, "trunkMode"),
		},
	}
}

func InterfaceFCPolicyFromContainer(cont *container.Container) *InterfaceFCPolicy {

	return InterfaceFCPolicyFromContainerList(cont, 0)
}

func InterfaceFCPolicyListFromContainer(cont *container.Container) []*InterfaceFCPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*InterfaceFCPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = InterfaceFCPolicyFromContainerList(cont, i)
	}

	return arr
}
