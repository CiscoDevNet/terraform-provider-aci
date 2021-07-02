package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const LacplagpolClassName = "lacpLagPol"

type LACPPolicy struct {
	BaseAttributes
	LACPPolicyAttributes
}

type LACPPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	MaxLinks string `json:",omitempty"`

	MinLinks string `json:",omitempty"`

	Mode string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLACPPolicy(lacpLagPolRn, parentDn, description string, lacpLagPolattr LACPPolicyAttributes) *LACPPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, lacpLagPolRn)
	return &LACPPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LacplagpolClassName,
			Rn:                lacpLagPolRn,
		},

		LACPPolicyAttributes: lacpLagPolattr,
	}
}

func (lacpLagPol *LACPPolicy) ToMap() (map[string]string, error) {
	lacpLagPolMap, err := lacpLagPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(lacpLagPolMap, "name", lacpLagPol.Name)

	A(lacpLagPolMap, "annotation", lacpLagPol.Annotation)

	A(lacpLagPolMap, "ctrl", lacpLagPol.Ctrl)

	A(lacpLagPolMap, "maxLinks", lacpLagPol.MaxLinks)

	A(lacpLagPolMap, "minLinks", lacpLagPol.MinLinks)

	A(lacpLagPolMap, "mode", lacpLagPol.Mode)

	A(lacpLagPolMap, "nameAlias", lacpLagPol.NameAlias)

	return lacpLagPolMap, err
}

func LACPPolicyFromContainerList(cont *container.Container, index int) *LACPPolicy {

	LACPPolicyCont := cont.S("imdata").Index(index).S(LacplagpolClassName, "attributes")
	return &LACPPolicy{
		BaseAttributes{
			DistinguishedName: G(LACPPolicyCont, "dn"),
			Description:       G(LACPPolicyCont, "descr"),
			Status:            G(LACPPolicyCont, "status"),
			ClassName:         LacplagpolClassName,
			Rn:                G(LACPPolicyCont, "rn"),
		},

		LACPPolicyAttributes{

			Name: G(LACPPolicyCont, "name"),

			Annotation: G(LACPPolicyCont, "annotation"),

			Ctrl: G(LACPPolicyCont, "ctrl"),

			MaxLinks: G(LACPPolicyCont, "maxLinks"),

			MinLinks: G(LACPPolicyCont, "minLinks"),

			Mode: G(LACPPolicyCont, "mode"),

			NameAlias: G(LACPPolicyCont, "nameAlias"),
		},
	}
}

func LACPPolicyFromContainer(cont *container.Container) *LACPPolicy {

	return LACPPolicyFromContainerList(cont, 0)
}

func LACPPolicyListFromContainer(cont *container.Container) []*LACPPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LACPPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = LACPPolicyFromContainerList(cont, i)
	}

	return arr
}
