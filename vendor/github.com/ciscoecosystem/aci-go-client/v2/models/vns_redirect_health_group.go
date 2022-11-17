package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnvnsRedirectHealthGroup        = "uni/tn-%s/svcCont/redirectHealthGroup-%s"
	RnvnsRedirectHealthGroup        = "redirectHealthGroup-%s"
	ParentDnvnsRedirectHealthGroup  = "uni/tn-%s/svcCont"
	VnsredirecthealthgroupClassName = "vnsRedirectHealthGroup"
)

type L4L7RedirectHealthGroup struct {
	BaseAttributes
	NameAliasAttribute
	L4L7RedirectHealthGroupAttributes
}

type L4L7RedirectHealthGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewL4L7RedirectHealthGroup(vnsRedirectHealthGroupRn, parentDn, description, nameAlias string, vnsRedirectHealthGroupAttr L4L7RedirectHealthGroupAttributes) *L4L7RedirectHealthGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsRedirectHealthGroupRn)
	return &L4L7RedirectHealthGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsredirecthealthgroupClassName,
			Rn:                vnsRedirectHealthGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		L4L7RedirectHealthGroupAttributes: vnsRedirectHealthGroupAttr,
	}
}

func (vnsRedirectHealthGroup *L4L7RedirectHealthGroup) ToMap() (map[string]string, error) {
	vnsRedirectHealthGroupMap, err := vnsRedirectHealthGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsRedirectHealthGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsRedirectHealthGroupMap, key, value)
	}

	A(vnsRedirectHealthGroupMap, "annotation", vnsRedirectHealthGroup.Annotation)
	A(vnsRedirectHealthGroupMap, "name", vnsRedirectHealthGroup.Name)
	return vnsRedirectHealthGroupMap, err
}

func L4L7RedirectHealthGroupFromContainerList(cont *container.Container, index int) *L4L7RedirectHealthGroup {
	L4L7RedirectHealthGroupCont := cont.S("imdata").Index(index).S(VnsredirecthealthgroupClassName, "attributes")
	return &L4L7RedirectHealthGroup{
		BaseAttributes{
			DistinguishedName: G(L4L7RedirectHealthGroupCont, "dn"),
			Description:       G(L4L7RedirectHealthGroupCont, "descr"),
			Status:            G(L4L7RedirectHealthGroupCont, "status"),
			ClassName:         VnsredirecthealthgroupClassName,
			Rn:                G(L4L7RedirectHealthGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(L4L7RedirectHealthGroupCont, "nameAlias"),
		},
		L4L7RedirectHealthGroupAttributes{
			Annotation: G(L4L7RedirectHealthGroupCont, "annotation"),
			Name:       G(L4L7RedirectHealthGroupCont, "name"),
		},
	}
}

func L4L7RedirectHealthGroupFromContainer(cont *container.Container) *L4L7RedirectHealthGroup {
	return L4L7RedirectHealthGroupFromContainerList(cont, 0)
}

func L4L7RedirectHealthGroupListFromContainer(cont *container.Container) []*L4L7RedirectHealthGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*L4L7RedirectHealthGroup, length)

	for i := 0; i < length; i++ {
		arr[i] = L4L7RedirectHealthGroupFromContainerList(cont, i)
	}

	return arr
}
