package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnpkiEp        = "uni/userext/pkiext"
	RnpkiEp        = "pkiext"
	ParentDnpkiEp  = "uni/userext"
	PkiepClassName = "pkiEp"
)

type PublicKeyManagement struct {
	BaseAttributes
	NameAliasAttribute
	PublicKeyManagementAttributes
}

type PublicKeyManagementAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewPublicKeyManagement(pkiEpRn, parentDn, description, nameAlias string, pkiEpAttr PublicKeyManagementAttributes) *PublicKeyManagement {
	dn := fmt.Sprintf("%s/%s", parentDn, pkiEpRn)
	return &PublicKeyManagement{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PkiepClassName,
			Rn:                pkiEpRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		PublicKeyManagementAttributes: pkiEpAttr,
	}
}

func (pkiEp *PublicKeyManagement) ToMap() (map[string]string, error) {
	pkiEpMap, err := pkiEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := pkiEp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(pkiEpMap, key, value)
	}
	A(pkiEpMap, "annotation", pkiEp.Annotation)
	A(pkiEpMap, "name", pkiEp.Name)
	return pkiEpMap, err
}

func PublicKeyManagementFromContainerList(cont *container.Container, index int) *PublicKeyManagement {
	PublicKeyManagementCont := cont.S("imdata").Index(index).S(PkiepClassName, "attributes")
	return &PublicKeyManagement{
		BaseAttributes{
			DistinguishedName: G(PublicKeyManagementCont, "dn"),
			Description:       G(PublicKeyManagementCont, "descr"),
			Status:            G(PublicKeyManagementCont, "status"),
			ClassName:         PkiepClassName,
			Rn:                G(PublicKeyManagementCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(PublicKeyManagementCont, "nameAlias"),
		},
		PublicKeyManagementAttributes{
			Annotation: G(PublicKeyManagementCont, "annotation"),
			Name:       G(PublicKeyManagementCont, "name"),
		},
	}
}

func PublicKeyManagementFromContainer(cont *container.Container) *PublicKeyManagement {
	return PublicKeyManagementFromContainerList(cont, 0)
}

func PublicKeyManagementListFromContainer(cont *container.Container) []*PublicKeyManagement {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PublicKeyManagement, length)
	for i := 0; i < length; i++ {
		arr[i] = PublicKeyManagementFromContainerList(cont, i)
	}
	return arr
}
