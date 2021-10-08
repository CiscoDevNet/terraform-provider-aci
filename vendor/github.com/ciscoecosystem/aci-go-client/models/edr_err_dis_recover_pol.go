package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnedrErrDisRecoverPol        = "uni/infra/edrErrDisRecoverPol-%s"
	RnedrErrDisRecoverPol        = "edrErrDisRecoverPol-%s"
	ParentDnedrErrDisRecoverPol  = "uni/infra"
	EdrerrdisrecoverpolClassName = "edrErrDisRecoverPol"
)

type ErrorDisabledRecoveryPolicy struct {
	BaseAttributes
	NameAliasAttribute
	ErrorDisabledRecoveryPolicyAttributes
}

type ErrorDisabledRecoveryPolicyAttributes struct {
	Annotation       string `json:",omitempty"`
	ErrDisRecovIntvl string `json:",omitempty"`
	Name             string `json:",omitempty"`
}

func NewErrorDisabledRecoveryPolicy(edrErrDisRecoverPolRn, parentDn, description, nameAlias string, edrErrDisRecoverPolAttr ErrorDisabledRecoveryPolicyAttributes) *ErrorDisabledRecoveryPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, edrErrDisRecoverPolRn)
	return &ErrorDisabledRecoveryPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         EdrerrdisrecoverpolClassName,
			Rn:                edrErrDisRecoverPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ErrorDisabledRecoveryPolicyAttributes: edrErrDisRecoverPolAttr,
	}
}

func (edrErrDisRecoverPol *ErrorDisabledRecoveryPolicy) ToMap() (map[string]string, error) {
	edrErrDisRecoverPolMap, err := edrErrDisRecoverPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := edrErrDisRecoverPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(edrErrDisRecoverPolMap, key, value)
	}
	A(edrErrDisRecoverPolMap, "annotation", edrErrDisRecoverPol.Annotation)
	A(edrErrDisRecoverPolMap, "errDisRecovIntvl", edrErrDisRecoverPol.ErrDisRecovIntvl)
	A(edrErrDisRecoverPolMap, "name", edrErrDisRecoverPol.Name)
	return edrErrDisRecoverPolMap, err
}

func ErrorDisabledRecoveryPolicyFromContainerList(cont *container.Container, index int) *ErrorDisabledRecoveryPolicy {
	ErrorDisabledRecoveryPolicyCont := cont.S("imdata").Index(index).S(EdrerrdisrecoverpolClassName, "attributes")
	return &ErrorDisabledRecoveryPolicy{
		BaseAttributes{
			DistinguishedName: G(ErrorDisabledRecoveryPolicyCont, "dn"),
			Description:       G(ErrorDisabledRecoveryPolicyCont, "descr"),
			Status:            G(ErrorDisabledRecoveryPolicyCont, "status"),
			ClassName:         EdrerrdisrecoverpolClassName,
			Rn:                G(ErrorDisabledRecoveryPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ErrorDisabledRecoveryPolicyCont, "nameAlias"),
		},
		ErrorDisabledRecoveryPolicyAttributes{
			Annotation:       G(ErrorDisabledRecoveryPolicyCont, "annotation"),
			ErrDisRecovIntvl: G(ErrorDisabledRecoveryPolicyCont, "errDisRecovIntvl"),
			Name:             G(ErrorDisabledRecoveryPolicyCont, "name"),
		},
	}
}

func ErrorDisabledRecoveryPolicyFromContainer(cont *container.Container) *ErrorDisabledRecoveryPolicy {
	return ErrorDisabledRecoveryPolicyFromContainerList(cont, 0)
}

func ErrorDisabledRecoveryPolicyListFromContainer(cont *container.Container) []*ErrorDisabledRecoveryPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ErrorDisabledRecoveryPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = ErrorDisabledRecoveryPolicyFromContainerList(cont, i)
	}
	return arr
}
