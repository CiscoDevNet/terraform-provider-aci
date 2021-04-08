package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const HsrpgrouppolClassName = "hsrpGroupPol"

type HSRPGroupPolicy struct {
	BaseAttributes
	HSRPGroupPolicyAttributes
}

type HSRPGroupPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	HelloIntvl string `json:",omitempty"`

	HoldIntvl string `json:",omitempty"`

	Key string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PreemptDelayMin string `json:",omitempty"`

	PreemptDelayReload string `json:",omitempty"`

	PreemptDelaySync string `json:",omitempty"`

	Prio string `json:",omitempty"`

	SecureAuthKey string `json:",omitempty"`

	Timeout string `json:",omitempty"`

	HSRPGroupPolicy_type string `json:",omitempty"`
}

func NewHSRPGroupPolicy(hsrpGroupPolRn, parentDn, description string, hsrpGroupPolattr HSRPGroupPolicyAttributes) *HSRPGroupPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, hsrpGroupPolRn)
	return &HSRPGroupPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         HsrpgrouppolClassName,
			Rn:                hsrpGroupPolRn,
		},

		HSRPGroupPolicyAttributes: hsrpGroupPolattr,
	}
}

func (hsrpGroupPol *HSRPGroupPolicy) ToMap() (map[string]string, error) {
	hsrpGroupPolMap, err := hsrpGroupPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(hsrpGroupPolMap, "name", hsrpGroupPol.Name)

	A(hsrpGroupPolMap, "annotation", hsrpGroupPol.Annotation)

	A(hsrpGroupPolMap, "ctrl", hsrpGroupPol.Ctrl)

	A(hsrpGroupPolMap, "helloIntvl", hsrpGroupPol.HelloIntvl)

	A(hsrpGroupPolMap, "holdIntvl", hsrpGroupPol.HoldIntvl)

	A(hsrpGroupPolMap, "key", hsrpGroupPol.Key)

	A(hsrpGroupPolMap, "nameAlias", hsrpGroupPol.NameAlias)

	A(hsrpGroupPolMap, "preemptDelayMin", hsrpGroupPol.PreemptDelayMin)

	A(hsrpGroupPolMap, "preemptDelayReload", hsrpGroupPol.PreemptDelayReload)

	A(hsrpGroupPolMap, "preemptDelaySync", hsrpGroupPol.PreemptDelaySync)

	A(hsrpGroupPolMap, "prio", hsrpGroupPol.Prio)

	A(hsrpGroupPolMap, "secureAuthKey", hsrpGroupPol.SecureAuthKey)

	A(hsrpGroupPolMap, "timeout", hsrpGroupPol.Timeout)

	A(hsrpGroupPolMap, "type", hsrpGroupPol.HSRPGroupPolicy_type)

	return hsrpGroupPolMap, err
}

func HSRPGroupPolicyFromContainerList(cont *container.Container, index int) *HSRPGroupPolicy {

	HSRPGroupPolicyCont := cont.S("imdata").Index(index).S(HsrpgrouppolClassName, "attributes")
	return &HSRPGroupPolicy{
		BaseAttributes{
			DistinguishedName: G(HSRPGroupPolicyCont, "dn"),
			Description:       G(HSRPGroupPolicyCont, "descr"),
			Status:            G(HSRPGroupPolicyCont, "status"),
			ClassName:         HsrpgrouppolClassName,
			Rn:                G(HSRPGroupPolicyCont, "rn"),
		},

		HSRPGroupPolicyAttributes{

			Name: G(HSRPGroupPolicyCont, "name"),

			Annotation: G(HSRPGroupPolicyCont, "annotation"),

			Ctrl: G(HSRPGroupPolicyCont, "ctrl"),

			HelloIntvl: G(HSRPGroupPolicyCont, "helloIntvl"),

			HoldIntvl: G(HSRPGroupPolicyCont, "holdIntvl"),

			Key: G(HSRPGroupPolicyCont, "key"),

			NameAlias: G(HSRPGroupPolicyCont, "nameAlias"),

			PreemptDelayMin: G(HSRPGroupPolicyCont, "preemptDelayMin"),

			PreemptDelayReload: G(HSRPGroupPolicyCont, "preemptDelayReload"),

			PreemptDelaySync: G(HSRPGroupPolicyCont, "preemptDelaySync"),

			Prio: G(HSRPGroupPolicyCont, "prio"),

			SecureAuthKey: G(HSRPGroupPolicyCont, "secureAuthKey"),

			Timeout: G(HSRPGroupPolicyCont, "timeout"),

			HSRPGroupPolicy_type: G(HSRPGroupPolicyCont, "type"),
		},
	}
}

func HSRPGroupPolicyFromContainer(cont *container.Container) *HSRPGroupPolicy {

	return HSRPGroupPolicyFromContainerList(cont, 0)
}

func HSRPGroupPolicyListFromContainer(cont *container.Container) []*HSRPGroupPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*HSRPGroupPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = HSRPGroupPolicyFromContainerList(cont, i)
	}

	return arr
}
