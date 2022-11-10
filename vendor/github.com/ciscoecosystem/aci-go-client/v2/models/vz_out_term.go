package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnvzOutTerm        = "uni/tn-%s/brc-%s/subj-%s/outtmnl"
	RnvzOutTerm        = "outtmnl"
	ParentDnvzOutTerm  = "uni/tn-%s/brc-%s/subj-%s"
	VzouttermClassName = "vzOutTerm"
)

type OutTermSubject struct {
	BaseAttributes
	NameAliasAttribute
	OutTermSubjectAttributes
}

type OutTermSubjectAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Prio       string `json:",omitempty"`
	TargetDscp string `json:",omitempty"`
}

func NewOutTermSubject(vzOutTermRn, parentDn, description, nameAlias string, vzOutTermAttr OutTermSubjectAttributes) *OutTermSubject {
	dn := fmt.Sprintf("%s/%s", parentDn, vzOutTermRn)
	return &OutTermSubject{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzouttermClassName,
			Rn:                vzOutTermRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		OutTermSubjectAttributes: vzOutTermAttr,
	}
}

func (vzOutTerm *OutTermSubject) ToMap() (map[string]string, error) {
	vzOutTermMap, err := vzOutTerm.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vzOutTerm.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vzOutTermMap, key, value)
	}

	A(vzOutTermMap, "name", vzOutTerm.Name)
	A(vzOutTermMap, "prio", vzOutTerm.Prio)
	A(vzOutTermMap, "targetDscp", vzOutTerm.TargetDscp)
	return vzOutTermMap, err
}

func OutTermSubjectFromContainerList(cont *container.Container, index int) *OutTermSubject {
	OutTermSubjectCont := cont.S("imdata").Index(index).S(VzouttermClassName, "attributes")
	return &OutTermSubject{
		BaseAttributes{
			DistinguishedName: G(OutTermSubjectCont, "dn"),
			Description:       G(OutTermSubjectCont, "descr"),
			Status:            G(OutTermSubjectCont, "status"),
			ClassName:         VzouttermClassName,
			Rn:                G(OutTermSubjectCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(OutTermSubjectCont, "nameAlias"),
		},
		OutTermSubjectAttributes{
			Name:       G(OutTermSubjectCont, "name"),
			Prio:       G(OutTermSubjectCont, "prio"),
			TargetDscp: G(OutTermSubjectCont, "targetDscp"),
		},
	}
}

func OutTermSubjectFromContainer(cont *container.Container) *OutTermSubject {
	return OutTermSubjectFromContainerList(cont, 0)
}

func OutTermSubjectListFromContainer(cont *container.Container) []*OutTermSubject {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*OutTermSubject, length)

	for i := 0; i < length; i++ {
		arr[i] = OutTermSubjectFromContainerList(cont, i)
	}

	return arr
}
