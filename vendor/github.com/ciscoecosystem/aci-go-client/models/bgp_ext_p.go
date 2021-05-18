package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpextpClassName = "bgpExtP"

type L3outBgpExternalPolicy struct {
	BaseAttributes
	L3outBgpExternalPolicyAttributes
}

type L3outBgpExternalPolicyAttributes struct {
	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outBgpExternalPolicy(bgpExtPRn, parentDn, description string, bgpExtPattr L3outBgpExternalPolicyAttributes) *L3outBgpExternalPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpExtPRn)
	return &L3outBgpExternalPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgpextpClassName,
			Rn:                bgpExtPRn,
		},

		L3outBgpExternalPolicyAttributes: bgpExtPattr,
	}
}

func (bgpExtP *L3outBgpExternalPolicy) ToMap() (map[string]string, error) {
	bgpExtPMap, err := bgpExtP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpExtPMap, "annotation", bgpExtP.Annotation)

	A(bgpExtPMap, "nameAlias", bgpExtP.NameAlias)

	return bgpExtPMap, err
}

func L3outBgpExternalPolicyFromContainerList(cont *container.Container, index int) *L3outBgpExternalPolicy {

	L3outBgpExternalPolicyCont := cont.S("imdata").Index(index).S(BgpextpClassName, "attributes")
	return &L3outBgpExternalPolicy{
		BaseAttributes{
			DistinguishedName: G(L3outBgpExternalPolicyCont, "dn"),
			Description:       G(L3outBgpExternalPolicyCont, "descr"),
			Status:            G(L3outBgpExternalPolicyCont, "status"),
			ClassName:         BgpextpClassName,
			Rn:                G(L3outBgpExternalPolicyCont, "rn"),
		},

		L3outBgpExternalPolicyAttributes{

			Annotation: G(L3outBgpExternalPolicyCont, "annotation"),

			NameAlias: G(L3outBgpExternalPolicyCont, "nameAlias"),
		},
	}
}

func L3outBgpExternalPolicyFromContainer(cont *container.Container) *L3outBgpExternalPolicy {

	return L3outBgpExternalPolicyFromContainerList(cont, 0)
}

func L3outBgpExternalPolicyListFromContainer(cont *container.Container) []*L3outBgpExternalPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outBgpExternalPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outBgpExternalPolicyFromContainerList(cont, i)
	}

	return arr
}
