package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaUserEp        = "uni/userext"
	RnaaaUserEp        = "userext"
	ParentDnaaaUserEp  = "uni"
	AaauserepClassName = "aaaUserEp"
)

type UserManagement struct {
	BaseAttributes
	NameAliasAttribute
	UserManagementAttributes
}

type UserManagementAttributes struct {
	Annotation       string `json:",omitempty"`
	Name             string `json:",omitempty"`
	PwdStrengthCheck string `json:",omitempty"`
}

func NewUserManagement(aaaUserEpRn, parentDn, description, nameAlias string, aaaUserEpAttr UserManagementAttributes) *UserManagement {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaUserEpRn)
	return &UserManagement{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaauserepClassName,
			Rn:                aaaUserEpRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		UserManagementAttributes: aaaUserEpAttr,
	}
}

func (aaaUserEp *UserManagement) ToMap() (map[string]string, error) {
	aaaUserEpMap, err := aaaUserEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaUserEp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaUserEpMap, key, value)
	}
	A(aaaUserEpMap, "annotation", aaaUserEp.Annotation)
	A(aaaUserEpMap, "name", aaaUserEp.Name)
	A(aaaUserEpMap, "pwdStrengthCheck", aaaUserEp.PwdStrengthCheck)
	return aaaUserEpMap, err
}

func UserManagementFromContainerList(cont *container.Container, index int) *UserManagement {
	UserManagementCont := cont.S("imdata").Index(index).S(AaauserepClassName, "attributes")
	return &UserManagement{
		BaseAttributes{
			DistinguishedName: G(UserManagementCont, "dn"),
			Description:       G(UserManagementCont, "descr"),
			Status:            G(UserManagementCont, "status"),
			ClassName:         AaauserepClassName,
			Rn:                G(UserManagementCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(UserManagementCont, "nameAlias"),
		},
		UserManagementAttributes{
			Annotation:       G(UserManagementCont, "annotation"),
			Name:             G(UserManagementCont, "name"),
			PwdStrengthCheck: G(UserManagementCont, "pwdStrengthCheck"),
		},
	}
}

func UserManagementFromContainer(cont *container.Container) *UserManagement {
	return UserManagementFromContainerList(cont, 0)
}

func UserManagementListFromContainer(cont *container.Container) []*UserManagement {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*UserManagement, length)
	for i := 0; i < length; i++ {
		arr[i] = UserManagementFromContainerList(cont, i)
	}
	return arr
}
