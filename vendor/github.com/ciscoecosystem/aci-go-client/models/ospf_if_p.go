package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const OspfifpClassName = "ospfIfP"

type OSPFInterfaceProfile struct {
	BaseAttributes
	OSPFInterfaceProfileAttributes
}

type OSPFInterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`

	AuthKey string `json:",omitempty"`

	AuthKeyId string `json:",omitempty"`

	AuthType string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewOSPFInterfaceProfile(ospfIfPRn, parentDn, description string, ospfIfPattr OSPFInterfaceProfileAttributes) *OSPFInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, ospfIfPRn)
	return &OSPFInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         OspfifpClassName,
			Rn:                ospfIfPRn,
		},

		OSPFInterfaceProfileAttributes: ospfIfPattr,
	}
}

func (ospfIfP *OSPFInterfaceProfile) ToMap() (map[string]string, error) {
	ospfIfPMap, err := ospfIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ospfIfPMap, "annotation", ospfIfP.Annotation)

	A(ospfIfPMap, "authKey", ospfIfP.AuthKey)

	A(ospfIfPMap, "authKeyId", ospfIfP.AuthKeyId)

	A(ospfIfPMap, "authType", ospfIfP.AuthType)

	A(ospfIfPMap, "nameAlias", ospfIfP.NameAlias)

	return ospfIfPMap, err
}

func OSPFInterfaceProfileFromContainerList(cont *container.Container, index int) *OSPFInterfaceProfile {

	InterfaceProfileCont := cont.S("imdata").Index(index).S(OspfifpClassName, "attributes")
	return &OSPFInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         OspfifpClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},

		OSPFInterfaceProfileAttributes{

			Annotation: G(InterfaceProfileCont, "annotation"),

			AuthKey: G(InterfaceProfileCont, "authKey"),

			AuthKeyId: G(InterfaceProfileCont, "authKeyId"),

			AuthType: G(InterfaceProfileCont, "authType"),

			NameAlias: G(InterfaceProfileCont, "nameAlias"),
		},
	}
}

func OSPFInterfaceProfileFromContainer(cont *container.Container) *OSPFInterfaceProfile {

	return OSPFInterfaceProfileFromContainerList(cont, 0)
}

func OSPFInterfaceProfileListFromContainer(cont *container.Container) []*OSPFInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OSPFInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = OSPFInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}
