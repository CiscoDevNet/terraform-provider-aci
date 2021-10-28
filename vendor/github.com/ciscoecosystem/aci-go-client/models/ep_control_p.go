package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnepControlP        = "uni/infra/epCtrlP-%s"
	RnepControlP        = "epCtrlP-%s"
	ParentDnepControlP  = "uni/infra"
	EpcontrolpClassName = "epControlP"
)

type EndpointControlPolicy struct {
	BaseAttributes
	NameAliasAttribute
	EndpointControlPolicyAttributes
}

type EndpointControlPolicyAttributes struct {
	AdminSt            string `json:",omitempty"`
	Annotation         string `json:",omitempty"`
	HoldIntvl          string `json:",omitempty"`
	Name               string `json:",omitempty"`
	RogueEpDetectIntvl string `json:",omitempty"`
	RogueEpDetectMult  string `json:",omitempty"`
}

func NewEndpointControlPolicy(epControlPRn, parentDn, description, nameAlias string, epControlPAttr EndpointControlPolicyAttributes) *EndpointControlPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, epControlPRn)
	return &EndpointControlPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         EpcontrolpClassName,
			Rn:                epControlPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EndpointControlPolicyAttributes: epControlPAttr,
	}
}

func (epControlP *EndpointControlPolicy) ToMap() (map[string]string, error) {
	epControlPMap, err := epControlP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := epControlP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(epControlPMap, key, value)
	}
	A(epControlPMap, "adminSt", epControlP.AdminSt)
	A(epControlPMap, "annotation", epControlP.Annotation)
	A(epControlPMap, "holdIntvl", epControlP.HoldIntvl)
	A(epControlPMap, "name", epControlP.Name)
	A(epControlPMap, "rogueEpDetectIntvl", epControlP.RogueEpDetectIntvl)
	A(epControlPMap, "rogueEpDetectMult", epControlP.RogueEpDetectMult)
	return epControlPMap, err
}

func EndpointControlPolicyFromContainerList(cont *container.Container, index int) *EndpointControlPolicy {
	EndpointControlPolicyCont := cont.S("imdata").Index(index).S(EpcontrolpClassName, "attributes")
	return &EndpointControlPolicy{
		BaseAttributes{
			DistinguishedName: G(EndpointControlPolicyCont, "dn"),
			Description:       G(EndpointControlPolicyCont, "descr"),
			Status:            G(EndpointControlPolicyCont, "status"),
			ClassName:         EpcontrolpClassName,
			Rn:                G(EndpointControlPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EndpointControlPolicyCont, "nameAlias"),
		},
		EndpointControlPolicyAttributes{
			AdminSt:            G(EndpointControlPolicyCont, "adminSt"),
			Annotation:         G(EndpointControlPolicyCont, "annotation"),
			HoldIntvl:          G(EndpointControlPolicyCont, "holdIntvl"),
			Name:               G(EndpointControlPolicyCont, "name"),
			RogueEpDetectIntvl: G(EndpointControlPolicyCont, "rogueEpDetectIntvl"),
			RogueEpDetectMult:  G(EndpointControlPolicyCont, "rogueEpDetectMult"),
		},
	}
}

func EndpointControlPolicyFromContainer(cont *container.Container) *EndpointControlPolicy {
	return EndpointControlPolicyFromContainerList(cont, 0)
}

func EndpointControlPolicyListFromContainer(cont *container.Container) []*EndpointControlPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EndpointControlPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = EndpointControlPolicyFromContainerList(cont, i)
	}
	return arr
}
