package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaPwdProfile        = "uni/userext/pwdprofile"
	RnaaaPwdProfile        = "pwdprofile"
	ParentDnaaaPwdProfile  = "uni/userext"
	AaapwdprofileClassName = "aaaPwdProfile"
)

type PasswordChangeExpirationPolicy struct {
	BaseAttributes
	NameAliasAttribute
	PasswordChangeExpirationPolicyAttributes
}

type PasswordChangeExpirationPolicyAttributes struct {
	Annotation           string `json:",omitempty"`
	ChangeCount          string `json:",omitempty"`
	ChangeDuringInterval string `json:",omitempty"`
	ChangeInterval       string `json:",omitempty"`
	ExpirationWarnTime   string `json:",omitempty"`
	HistoryCount         string `json:",omitempty"`
	Name                 string `json:",omitempty"`
	NoChangeInterval     string `json:",omitempty"`
}

func NewPasswordChangeExpirationPolicy(aaaPwdProfileRn, parentDn, description, nameAlias string, aaaPwdProfileAttr PasswordChangeExpirationPolicyAttributes) *PasswordChangeExpirationPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaPwdProfileRn)
	return &PasswordChangeExpirationPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaapwdprofileClassName,
			Rn:                aaaPwdProfileRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		PasswordChangeExpirationPolicyAttributes: aaaPwdProfileAttr,
	}
}

func (aaaPwdProfile *PasswordChangeExpirationPolicy) ToMap() (map[string]string, error) {
	aaaPwdProfileMap, err := aaaPwdProfile.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaPwdProfile.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaPwdProfileMap, key, value)
	}
	A(aaaPwdProfileMap, "annotation", aaaPwdProfile.Annotation)
	A(aaaPwdProfileMap, "changeCount", aaaPwdProfile.ChangeCount)
	A(aaaPwdProfileMap, "changeDuringInterval", aaaPwdProfile.ChangeDuringInterval)
	A(aaaPwdProfileMap, "changeInterval", aaaPwdProfile.ChangeInterval)
	A(aaaPwdProfileMap, "expirationWarnTime", aaaPwdProfile.ExpirationWarnTime)
	A(aaaPwdProfileMap, "historyCount", aaaPwdProfile.HistoryCount)
	A(aaaPwdProfileMap, "name", aaaPwdProfile.Name)
	A(aaaPwdProfileMap, "noChangeInterval", aaaPwdProfile.NoChangeInterval)
	return aaaPwdProfileMap, err
}

func PasswordChangeExpirationPolicyFromContainerList(cont *container.Container, index int) *PasswordChangeExpirationPolicy {
	PasswordChangeExpirationPolicyCont := cont.S("imdata").Index(index).S(AaapwdprofileClassName, "attributes")
	return &PasswordChangeExpirationPolicy{
		BaseAttributes{
			DistinguishedName: G(PasswordChangeExpirationPolicyCont, "dn"),
			Description:       G(PasswordChangeExpirationPolicyCont, "descr"),
			Status:            G(PasswordChangeExpirationPolicyCont, "status"),
			ClassName:         AaapwdprofileClassName,
			Rn:                G(PasswordChangeExpirationPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(PasswordChangeExpirationPolicyCont, "nameAlias"),
		},
		PasswordChangeExpirationPolicyAttributes{
			Annotation:           G(PasswordChangeExpirationPolicyCont, "annotation"),
			ChangeCount:          G(PasswordChangeExpirationPolicyCont, "changeCount"),
			ChangeDuringInterval: G(PasswordChangeExpirationPolicyCont, "changeDuringInterval"),
			ChangeInterval:       G(PasswordChangeExpirationPolicyCont, "changeInterval"),
			ExpirationWarnTime:   G(PasswordChangeExpirationPolicyCont, "expirationWarnTime"),
			HistoryCount:         G(PasswordChangeExpirationPolicyCont, "historyCount"),
			Name:                 G(PasswordChangeExpirationPolicyCont, "name"),
			NoChangeInterval:     G(PasswordChangeExpirationPolicyCont, "noChangeInterval"),
		},
	}
}

func PasswordChangeExpirationPolicyFromContainer(cont *container.Container) *PasswordChangeExpirationPolicy {
	return PasswordChangeExpirationPolicyFromContainerList(cont, 0)
}

func PasswordChangeExpirationPolicyListFromContainer(cont *container.Container) []*PasswordChangeExpirationPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PasswordChangeExpirationPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = PasswordChangeExpirationPolicyFromContainerList(cont, i)
	}
	return arr
}
