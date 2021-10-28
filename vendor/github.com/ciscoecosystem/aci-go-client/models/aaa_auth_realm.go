package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaAuthRealm        = "uni/userext/authrealm"
	RnaaaAuthRealm        = "authrealm"
	ParentDnaaaAuthRealm  = "uni/userext"
	AaaauthrealmClassName = "aaaAuthRealm"
)

type AAAAuthentication struct {
	BaseAttributes
	NameAliasAttribute
	AAAAuthenticationAttributes
}

type AAAAuthenticationAttributes struct {
	Annotation    string `json:",omitempty"`
	DefRolePolicy string `json:",omitempty"`
	Name          string `json:",omitempty"`
}

func NewAAAAuthentication(aaaAuthRealmRn, parentDn, description, nameAlias string, aaaAuthRealmAttr AAAAuthenticationAttributes) *AAAAuthentication {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaAuthRealmRn)
	return &AAAAuthentication{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaauthrealmClassName,
			Rn:                aaaAuthRealmRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AAAAuthenticationAttributes: aaaAuthRealmAttr,
	}
}

func (aaaAuthRealm *AAAAuthentication) ToMap() (map[string]string, error) {
	aaaAuthRealmMap, err := aaaAuthRealm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaAuthRealm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaAuthRealmMap, key, value)
	}
	A(aaaAuthRealmMap, "annotation", aaaAuthRealm.Annotation)
	A(aaaAuthRealmMap, "defRolePolicy", aaaAuthRealm.DefRolePolicy)
	A(aaaAuthRealmMap, "name", aaaAuthRealm.Name)
	return aaaAuthRealmMap, err
}

func AAAAuthenticationFromContainerList(cont *container.Container, index int) *AAAAuthentication {
	AAAAuthenticationCont := cont.S("imdata").Index(index).S(AaaauthrealmClassName, "attributes")
	return &AAAAuthentication{
		BaseAttributes{
			DistinguishedName: G(AAAAuthenticationCont, "dn"),
			Description:       G(AAAAuthenticationCont, "descr"),
			Status:            G(AAAAuthenticationCont, "status"),
			ClassName:         AaaauthrealmClassName,
			Rn:                G(AAAAuthenticationCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AAAAuthenticationCont, "nameAlias"),
		},
		AAAAuthenticationAttributes{
			Annotation:    G(AAAAuthenticationCont, "annotation"),
			DefRolePolicy: G(AAAAuthenticationCont, "defRolePolicy"),
			Name:          G(AAAAuthenticationCont, "name"),
		},
	}
}

func AAAAuthenticationFromContainer(cont *container.Container) *AAAAuthentication {
	return AAAAuthenticationFromContainerList(cont, 0)
}

func AAAAuthenticationListFromContainer(cont *container.Container) []*AAAAuthentication {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AAAAuthentication, length)
	for i := 0; i < length; i++ {
		arr[i] = AAAAuthenticationFromContainerList(cont, i)
	}
	return arr
}
