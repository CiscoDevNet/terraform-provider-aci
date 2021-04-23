package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvmmUsrAccP        = "uni/vmmp-%s/dom-%s/usracc-%s"
	RnvmmUsrAccP        = "usracc-%s"
	ParentDnvmmUsrAccP  = "uni/vmmp-%s/dom-%s"
	VmmusraccpClassName = "vmmUsrAccP"
)

type VMMCredential struct {
	BaseAttributes
	NameAliasAttribute
	VMMCredentialAttributes
}

type VMMCredentialAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Pwd        string `json:",omitempty"`
	Usr        string `json:",omitempty"`
}

func NewVMMCredential(vmmUsrAccPRn, parentDn, description, nameAlias string, vmmUsrAccPAttr VMMCredentialAttributes) *VMMCredential {
	dn := fmt.Sprintf("%s/%s", parentDn, vmmUsrAccPRn)
	return &VMMCredential{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VmmusraccpClassName,
			Rn:                vmmUsrAccPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		VMMCredentialAttributes: vmmUsrAccPAttr,
	}
}

func (vmmUsrAccP *VMMCredential) ToMap() (map[string]string, error) {
	vmmUsrAccPMap, err := vmmUsrAccP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := vmmUsrAccP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(vmmUsrAccPMap, key, value)
	}
	A(vmmUsrAccPMap, "annotation", vmmUsrAccP.Annotation)
	A(vmmUsrAccPMap, "name", vmmUsrAccP.Name)
	A(vmmUsrAccPMap, "pwd", vmmUsrAccP.Pwd)
	A(vmmUsrAccPMap, "usr", vmmUsrAccP.Usr)
	return vmmUsrAccPMap, err
}

func VMMCredentialFromContainerList(cont *container.Container, index int) *VMMCredential {
	VMMCredentialCont := cont.S("imdata").Index(index).S(VmmusraccpClassName, "attributes")
	return &VMMCredential{
		BaseAttributes{
			DistinguishedName: G(VMMCredentialCont, "dn"),
			Description:       G(VMMCredentialCont, "descr"),
			Status:            G(VMMCredentialCont, "status"),
			ClassName:         VmmusraccpClassName,
			Rn:                G(VMMCredentialCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(VMMCredentialCont, "nameAlias"),
		},
		VMMCredentialAttributes{
			Annotation: G(VMMCredentialCont, "annotation"),
			Name:       G(VMMCredentialCont, "name"),
			Pwd:        G(VMMCredentialCont, "pwd"),
			Usr:        G(VMMCredentialCont, "usr"),
		},
	}
}

func VMMCredentialFromContainer(cont *container.Container) *VMMCredential {
	return VMMCredentialFromContainerList(cont, 0)
}

func VMMCredentialListFromContainer(cont *container.Container) []*VMMCredential {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*VMMCredential, length)
	for i := 0; i < length; i++ {
		arr[i] = VMMCredentialFromContainerList(cont, i)
	}
	return arr
}
