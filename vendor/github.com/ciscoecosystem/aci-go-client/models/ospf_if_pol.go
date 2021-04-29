package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const OspfifpolClassName = "ospfIfPol"

type OSPFInterfacePolicy struct {
	BaseAttributes
	OSPFInterfacePolicyAttributes
}

type OSPFInterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Cost string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	DeadIntvl string `json:",omitempty"`

	HelloIntvl string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	NwT string `json:",omitempty"`

	PfxSuppress string `json:",omitempty"`

	Prio string `json:",omitempty"`

	RexmitIntvl string `json:",omitempty"`

	XmitDelay string `json:",omitempty"`
}

func NewOSPFInterfacePolicy(ospfIfPolRn, parentDn, description string, ospfIfPolattr OSPFInterfacePolicyAttributes) *OSPFInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, ospfIfPolRn)
	return &OSPFInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         OspfifpolClassName,
			Rn:                ospfIfPolRn,
		},

		OSPFInterfacePolicyAttributes: ospfIfPolattr,
	}
}

func (ospfIfPol *OSPFInterfacePolicy) ToMap() (map[string]string, error) {
	ospfIfPolMap, err := ospfIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ospfIfPolMap, "name", ospfIfPol.Name)

	A(ospfIfPolMap, "annotation", ospfIfPol.Annotation)

	A(ospfIfPolMap, "cost", ospfIfPol.Cost)

	A(ospfIfPolMap, "ctrl", ospfIfPol.Ctrl)

	A(ospfIfPolMap, "deadIntvl", ospfIfPol.DeadIntvl)

	A(ospfIfPolMap, "helloIntvl", ospfIfPol.HelloIntvl)

	A(ospfIfPolMap, "nameAlias", ospfIfPol.NameAlias)

	A(ospfIfPolMap, "nwT", ospfIfPol.NwT)

	A(ospfIfPolMap, "pfxSuppress", ospfIfPol.PfxSuppress)

	A(ospfIfPolMap, "prio", ospfIfPol.Prio)

	A(ospfIfPolMap, "rexmitIntvl", ospfIfPol.RexmitIntvl)

	A(ospfIfPolMap, "xmitDelay", ospfIfPol.XmitDelay)

	return ospfIfPolMap, err
}

func OSPFInterfacePolicyFromContainerList(cont *container.Container, index int) *OSPFInterfacePolicy {

	OSPFInterfacePolicyCont := cont.S("imdata").Index(index).S(OspfifpolClassName, "attributes")
	return &OSPFInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(OSPFInterfacePolicyCont, "dn"),
			Description:       G(OSPFInterfacePolicyCont, "descr"),
			Status:            G(OSPFInterfacePolicyCont, "status"),
			ClassName:         OspfifpolClassName,
			Rn:                G(OSPFInterfacePolicyCont, "rn"),
		},

		OSPFInterfacePolicyAttributes{

			Name: G(OSPFInterfacePolicyCont, "name"),

			Annotation: G(OSPFInterfacePolicyCont, "annotation"),

			Cost: G(OSPFInterfacePolicyCont, "cost"),

			Ctrl: G(OSPFInterfacePolicyCont, "ctrl"),

			DeadIntvl: G(OSPFInterfacePolicyCont, "deadIntvl"),

			HelloIntvl: G(OSPFInterfacePolicyCont, "helloIntvl"),

			NameAlias: G(OSPFInterfacePolicyCont, "nameAlias"),

			NwT: G(OSPFInterfacePolicyCont, "nwT"),

			PfxSuppress: G(OSPFInterfacePolicyCont, "pfxSuppress"),

			Prio: G(OSPFInterfacePolicyCont, "prio"),

			RexmitIntvl: G(OSPFInterfacePolicyCont, "rexmitIntvl"),

			XmitDelay: G(OSPFInterfacePolicyCont, "xmitDelay"),
		},
	}
}

func OSPFInterfacePolicyFromContainer(cont *container.Container) *OSPFInterfacePolicy {

	return OSPFInterfacePolicyFromContainerList(cont, 0)
}

func OSPFInterfacePolicyListFromContainer(cont *container.Container) []*OSPFInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OSPFInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = OSPFInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}
