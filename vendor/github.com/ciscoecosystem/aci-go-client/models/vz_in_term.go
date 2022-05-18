package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvzInTerm        = "uni/tn-%s/brc-%s/subj-%s/intmnl"
	RnvzInTerm        = "intmnl"
	ParentDnvzInTerm  = "uni/tn-%s/brc-%s/subj-%s"
	VzintermClassName = "vzInTerm"
)

type InTermSubject struct {
	BaseAttributes
	NameAliasAttribute
	InTermSubjectAttributes
}

type InTermSubjectAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Prio       string `json:",omitempty"`
	TargetDscp string `json:",omitempty"`
}

func NewInTermSubject(vzInTermRn, parentDn, description, nameAlias string, vzInTermAttr InTermSubjectAttributes) *InTermSubject {
	dn := fmt.Sprintf("%s/%s", parentDn, vzInTermRn)
	return &InTermSubject{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzintermClassName,
			Rn:                vzInTermRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		InTermSubjectAttributes: vzInTermAttr,
	}
}

func (vzInTerm *InTermSubject) ToMap() (map[string]string, error) {
	vzInTermMap, err := vzInTerm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vzInTerm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vzInTermMap, key, value)
	}

	A(vzInTermMap, "name", vzInTerm.Name)
	A(vzInTermMap, "prio", vzInTerm.Prio)
	A(vzInTermMap, "targetDscp", vzInTerm.TargetDscp)
	return vzInTermMap, err
}

func InTermSubjectFromContainerList(cont *container.Container, index int) *InTermSubject {
	InTermSubjectCont := cont.S("imdata").Index(index).S(VzintermClassName, "attributes")
	return &InTermSubject{
		BaseAttributes{
			DistinguishedName: G(InTermSubjectCont, "dn"),
			Description:       G(InTermSubjectCont, "descr"),
			Status:            G(InTermSubjectCont, "status"),
			ClassName:         VzintermClassName,
			Rn:                G(InTermSubjectCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(InTermSubjectCont, "nameAlias"),
		},
		InTermSubjectAttributes{
			Name:       G(InTermSubjectCont, "name"),
			Prio:       G(InTermSubjectCont, "prio"),
			TargetDscp: G(InTermSubjectCont, "targetDscp"),
		},
	}
}

func InTermSubjectFromContainer(cont *container.Container) *InTermSubject {
	return InTermSubjectFromContainerList(cont, 0)
}

func InTermSubjectListFromContainer(cont *container.Container) []*InTermSubject {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*InTermSubject, length)

	for i := 0; i < length; i++ {
		arr[i] = InTermSubjectFromContainerList(cont, i)
	}

	return arr
}
