package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnbfdMhIfP        = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s/bfdMhIfP"
	RnbfdMhIfP        = "bfdMhIfP"
	ParentDnbfdMhIfP  = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s"
	BfdmhifpClassName = "bfdMhIfP"
)

type AciBfdMultihopInterfaceProfile struct {
	BaseAttributes
	NameAliasAttribute
	AciBfdMultihopInterfaceProfileAttributes
}

type AciBfdMultihopInterfaceProfileAttributes struct {
	Annotation            string `json:",omitempty"`
	Key                   string `json:",omitempty"`
	KeyId                 string `json:",omitempty"`
	Name                  string `json:",omitempty"`
	InterfaceProfile_type string `json:",omitempty"`
}

func NewAciBfdMultihopInterfaceProfile(bfdMhIfPRn, parentDn, description, nameAlias string, bfdMhIfPAttr AciBfdMultihopInterfaceProfileAttributes) *AciBfdMultihopInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bfdMhIfPRn)
	return &AciBfdMultihopInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BfdmhifpClassName,
			Rn:                bfdMhIfPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AciBfdMultihopInterfaceProfileAttributes: bfdMhIfPAttr,
	}
}

func (bfdMhIfP *AciBfdMultihopInterfaceProfile) ToMap() (map[string]string, error) {
	bfdMhIfPMap, err := bfdMhIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := bfdMhIfP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(bfdMhIfPMap, key, value)
	}

	A(bfdMhIfPMap, "annotation", bfdMhIfP.Annotation)
	A(bfdMhIfPMap, "key", bfdMhIfP.Key)
	A(bfdMhIfPMap, "keyId", bfdMhIfP.KeyId)
	A(bfdMhIfPMap, "name", bfdMhIfP.Name)
	A(bfdMhIfPMap, "InterfaceProfile_type", bfdMhIfP.InterfaceProfile_type)
	return bfdMhIfPMap, err
}

func AciBfdMultihopInterfaceProfileFromContainerList(cont *container.Container, index int) *AciBfdMultihopInterfaceProfile {
	AciBfdMultihopInterfaceProfileCont := cont.S("imdata").Index(index).S(BfdmhifpClassName, "attributes")
	return &AciBfdMultihopInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(AciBfdMultihopInterfaceProfileCont, "dn"),
			Description:       G(AciBfdMultihopInterfaceProfileCont, "descr"),
			Status:            G(AciBfdMultihopInterfaceProfileCont, "status"),
			ClassName:         BfdmhifpClassName,
			Rn:                G(AciBfdMultihopInterfaceProfileCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AciBfdMultihopInterfaceProfileCont, "nameAlias"),
		},
		AciBfdMultihopInterfaceProfileAttributes{
			Annotation:            G(AciBfdMultihopInterfaceProfileCont, "annotation"),
			Key:                   G(AciBfdMultihopInterfaceProfileCont, "key"),
			KeyId:                 G(AciBfdMultihopInterfaceProfileCont, "keyId"),
			Name:                  G(AciBfdMultihopInterfaceProfileCont, "name"),
			InterfaceProfile_type: G(AciBfdMultihopInterfaceProfileCont, "InterfaceProfile_type"),
		},
	}
}

func AciBfdMultihopInterfaceProfileFromContainer(cont *container.Container) *AciBfdMultihopInterfaceProfile {
	return AciBfdMultihopInterfaceProfileFromContainerList(cont, 0)
}

func AciBfdMultihopInterfaceProfileListFromContainer(cont *container.Container) []*AciBfdMultihopInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AciBfdMultihopInterfaceProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = AciBfdMultihopInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
