package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BfdifpClassName = "bfdIfP"

type BFDInterfaceProfile struct {
	BaseAttributes
	BFDInterfaceProfileAttributes
}

type BFDInterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`

	Key string `json:",omitempty"`

	KeyId string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	InterfaceProfileType string `json:",omitempty"`

	Userdom string `json:",omitempty"`
}

func NewBFDInterfaceProfile(bfdIfPRn, parentDn, description string, bfdIfPattr BFDInterfaceProfileAttributes) *BFDInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bfdIfPRn)
	return &BFDInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BfdifpClassName,
			Rn:                bfdIfPRn,
		},

		BFDInterfaceProfileAttributes: bfdIfPattr,
	}
}

func (bfdIfP *BFDInterfaceProfile) ToMap() (map[string]string, error) {
	bfdIfPMap, err := bfdIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bfdIfPMap, "annotation", bfdIfP.Annotation)

	A(bfdIfPMap, "key", bfdIfP.Key)

	A(bfdIfPMap, "keyId", bfdIfP.KeyId)

	A(bfdIfPMap, "nameAlias", bfdIfP.NameAlias)

	A(bfdIfPMap, "type", bfdIfP.InterfaceProfileType)

	A(bfdIfPMap, "userdom", bfdIfP.Userdom)

	return bfdIfPMap, err
}

func BFDInterfaceProfileFromContainerList(cont *container.Container, index int) *BFDInterfaceProfile {

	InterfaceProfileCont := cont.S("imdata").Index(index).S(BfdifpClassName, "attributes")
	return &BFDInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         BfdifpClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},

		BFDInterfaceProfileAttributes{

			Annotation: G(InterfaceProfileCont, "annotation"),

			Key: G(InterfaceProfileCont, "key"),

			KeyId: G(InterfaceProfileCont, "keyId"),

			NameAlias: G(InterfaceProfileCont, "nameAlias"),

			InterfaceProfileType: G(InterfaceProfileCont, "type"),

			Userdom: G(InterfaceProfileCont, "userdom"),
		},
	}
}

func BFDInterfaceProfileFromContainer(cont *container.Container) *BFDInterfaceProfile {

	return BFDInterfaceProfileFromContainerList(cont, 0)
}

func BFDInterfaceProfileListFromContainer(cont *container.Container) []*BFDInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BFDInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = BFDInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
