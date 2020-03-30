package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const AaauserClassName = "aaaUser"

type LocalUser struct {
	BaseAttributes
	LocalUserAttributes
}

type LocalUserAttributes struct {
	Name string `json:",omitempty"`

	AccountStatus string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	CertAttribute string `json:",omitempty"`

	ClearPwdHistory string `json:",omitempty"`

	Email string `json:",omitempty"`

	Expiration string `json:",omitempty"`

	Expires string `json:",omitempty"`

	FirstName string `json:",omitempty"`

	LastName string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Otpenable string `json:",omitempty"`

	Otpkey string `json:",omitempty"`

	Phone string `json:",omitempty"`

	Pwd string `json:",omitempty"`

	PwdLifeTime string `json:",omitempty"`

	PwdUpdateRequired string `json:",omitempty"`

	RbacString string `json:",omitempty"`

	UnixUserId string `json:",omitempty"`
}

func NewLocalUser(aaaUserRn, parentDn, description string, aaaUserattr LocalUserAttributes) *LocalUser {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaUserRn)
	return &LocalUser{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaauserClassName,
			Rn:                aaaUserRn,
		},

		LocalUserAttributes: aaaUserattr,
	}
}

func (aaaUser *LocalUser) ToMap() (map[string]string, error) {
	aaaUserMap, err := aaaUser.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(aaaUserMap, "name", aaaUser.Name)

	A(aaaUserMap, "accountStatus", aaaUser.AccountStatus)

	A(aaaUserMap, "annotation", aaaUser.Annotation)

	A(aaaUserMap, "certAttribute", aaaUser.CertAttribute)

	A(aaaUserMap, "clearPwdHistory", aaaUser.ClearPwdHistory)

	A(aaaUserMap, "email", aaaUser.Email)

	A(aaaUserMap, "expiration", aaaUser.Expiration)

	A(aaaUserMap, "expires", aaaUser.Expires)

	A(aaaUserMap, "firstName", aaaUser.FirstName)

	A(aaaUserMap, "lastName", aaaUser.LastName)

	A(aaaUserMap, "nameAlias", aaaUser.NameAlias)

	A(aaaUserMap, "otpenable", aaaUser.Otpenable)

	A(aaaUserMap, "otpkey", aaaUser.Otpkey)

	A(aaaUserMap, "phone", aaaUser.Phone)

	A(aaaUserMap, "pwd", aaaUser.Pwd)

	A(aaaUserMap, "pwdLifeTime", aaaUser.PwdLifeTime)

	A(aaaUserMap, "pwdUpdateRequired", aaaUser.PwdUpdateRequired)

	A(aaaUserMap, "rbacString", aaaUser.RbacString)

	A(aaaUserMap, "unixUserId", aaaUser.UnixUserId)

	return aaaUserMap, err
}

func LocalUserFromContainerList(cont *container.Container, index int) *LocalUser {

	LocalUserCont := cont.S("imdata").Index(index).S(AaauserClassName, "attributes")
	return &LocalUser{
		BaseAttributes{
			DistinguishedName: G(LocalUserCont, "dn"),
			Description:       G(LocalUserCont, "descr"),
			Status:            G(LocalUserCont, "status"),
			ClassName:         AaauserClassName,
			Rn:                G(LocalUserCont, "rn"),
		},

		LocalUserAttributes{

			Name: G(LocalUserCont, "name"),

			AccountStatus: G(LocalUserCont, "accountStatus"),

			Annotation: G(LocalUserCont, "annotation"),

			CertAttribute: G(LocalUserCont, "certAttribute"),

			ClearPwdHistory: G(LocalUserCont, "clearPwdHistory"),

			Email: G(LocalUserCont, "email"),

			Expiration: G(LocalUserCont, "expiration"),

			Expires: G(LocalUserCont, "expires"),

			FirstName: G(LocalUserCont, "firstName"),

			LastName: G(LocalUserCont, "lastName"),

			NameAlias: G(LocalUserCont, "nameAlias"),

			Otpenable: G(LocalUserCont, "otpenable"),

			Otpkey: G(LocalUserCont, "otpkey"),

			Phone: G(LocalUserCont, "phone"),

			Pwd: G(LocalUserCont, "pwd"),

			PwdLifeTime: G(LocalUserCont, "pwdLifeTime"),

			PwdUpdateRequired: G(LocalUserCont, "pwdUpdateRequired"),

			RbacString: G(LocalUserCont, "rbacString"),

			UnixUserId: G(LocalUserCont, "unixUserId"),
		},
	}
}

func LocalUserFromContainer(cont *container.Container) *LocalUser {

	return LocalUserFromContainerList(cont, 0)
}

func LocalUserListFromContainer(cont *container.Container) []*LocalUser {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LocalUser, length)

	for i := 0; i < length; i++ {

		arr[i] = LocalUserFromContainerList(cont, i)
	}

	return arr
}
