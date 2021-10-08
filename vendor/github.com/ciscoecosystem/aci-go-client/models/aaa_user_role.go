package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaUserRole        = "uni/userext/user-%s/userdomain-%s/role-%s"
	RnaaaUserRole        = "role-%s"
	ParentDnaaaUserRole  = "uni/userext/user-%s/userdomain-%s"
	AaauserroleClassName = "aaaUserRole"
)

type UserRole struct {
	BaseAttributes
	NameAliasAttribute
	UserRoleAttributes
}

type UserRoleAttributes struct {
	Name       string `json:",omitempty"`
	PrivType   string `json:",omitempty"`
	Annotation string `json:,"omitempty"`
}

func NewUserRole(aaaUserRoleRn, parentDn, description, nameAlias string, aaaUserRoleAttr UserRoleAttributes) *UserRole {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaUserRoleRn)
	return &UserRole{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaauserroleClassName,
			Rn:                aaaUserRoleRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		UserRoleAttributes: aaaUserRoleAttr,
	}
}

func (aaaUserRole *UserRole) ToMap() (map[string]string, error) {
	aaaUserRoleMap, err := aaaUserRole.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaUserRole.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaUserRoleMap, key, value)
	}
	A(aaaUserRoleMap, "name", aaaUserRole.Name)
	A(aaaUserRoleMap, "privType", aaaUserRole.PrivType)
	A(aaaUserRoleMap, "annotation", aaaUserRole.Annotation)
	return aaaUserRoleMap, err
}

func UserRoleFromContainerList(cont *container.Container, index int) *UserRole {
	UserRoleCont := cont.S("imdata").Index(index).S(AaauserroleClassName, "attributes")
	return &UserRole{
		BaseAttributes{
			DistinguishedName: G(UserRoleCont, "dn"),
			Description:       G(UserRoleCont, "descr"),
			Status:            G(UserRoleCont, "status"),
			ClassName:         AaauserroleClassName,
			Rn:                G(UserRoleCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(UserRoleCont, "nameAlias"),
		},
		UserRoleAttributes{
			Name:       G(UserRoleCont, "name"),
			PrivType:   G(UserRoleCont, "privType"),
			Annotation: G(UserRoleCont, "annotation"),
		},
	}
}

func UserRoleFromContainer(cont *container.Container) *UserRole {
	return UserRoleFromContainerList(cont, 0)
}

func UserRoleListFromContainer(cont *container.Container) []*UserRole {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*UserRole, length)
	for i := 0; i < length; i++ {
		arr[i] = UserRoleFromContainerList(cont, i)
	}
	return arr
}
