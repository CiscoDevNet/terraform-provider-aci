package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnSnmpUserP        = "user-%s"
	DnSnmpUserP        = "uni/fabric/snmppol-%s/user-%s"
	ParentDnSnmpUserP  = "uni/fabric/snmppol-%s"
	SnmpUserPClassName = "snmpUserP"
)

type SnmpUserProfile struct {
	BaseAttributes
	SnmpUserProfileAttributes
}

type SnmpUserProfileAttributes struct {
	Annotation string `json:",omitempty"`
	AuthKey    string `json:",omitempty"`
	AuthType   string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	PrivKey    string `json:",omitempty"`
	PrivType   string `json:",omitempty"`
}

func NewSnmpUserProfile(snmpUserPRn, parentDn, description string, snmpUserPAttr SnmpUserProfileAttributes) *SnmpUserProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, snmpUserPRn)
	return &SnmpUserProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         SnmpUserPClassName,
			Rn:                snmpUserPRn,
		},
		SnmpUserProfileAttributes: snmpUserPAttr,
	}
}

func (snmpUserP *SnmpUserProfile) ToMap() (map[string]string, error) {
	snmpUserPMap, err := snmpUserP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(snmpUserPMap, "annotation", snmpUserP.Annotation)
	A(snmpUserPMap, "authKey", snmpUserP.AuthKey)
	A(snmpUserPMap, "authType", snmpUserP.AuthType)
	A(snmpUserPMap, "name", snmpUserP.Name)
	A(snmpUserPMap, "nameAlias", snmpUserP.NameAlias)
	A(snmpUserPMap, "privKey", snmpUserP.PrivKey)
	A(snmpUserPMap, "privType", snmpUserP.PrivType)
	return snmpUserPMap, err
}

func SnmpUserProfileFromContainerList(cont *container.Container, index int) *SnmpUserProfile {
	SnmpUserProfileCont := cont.S("imdata").Index(index).S(SnmpUserPClassName, "attributes")
	return &SnmpUserProfile{
		BaseAttributes{
			DistinguishedName: G(SnmpUserProfileCont, "dn"),
			Description:       G(SnmpUserProfileCont, "descr"),
			Status:            G(SnmpUserProfileCont, "status"),
			ClassName:         SnmpUserPClassName,
			Rn:                G(SnmpUserProfileCont, "rn"),
		},
		SnmpUserProfileAttributes{
			Annotation: G(SnmpUserProfileCont, "annotation"),
			AuthKey:    G(SnmpUserProfileCont, "authKey"),
			AuthType:   G(SnmpUserProfileCont, "authType"),
			Name:       G(SnmpUserProfileCont, "name"),
			NameAlias:  G(SnmpUserProfileCont, "nameAlias"),
			PrivKey:    G(SnmpUserProfileCont, "privKey"),
			PrivType:   G(SnmpUserProfileCont, "privType"),
		},
	}
}

func SnmpUserProfileFromContainer(cont *container.Container) *SnmpUserProfile {
	return SnmpUserProfileFromContainerList(cont, 0)
}

func SnmpUserProfileListFromContainer(cont *container.Container) []*SnmpUserProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SnmpUserProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = SnmpUserProfileFromContainerList(cont, i)
	}

	return arr
}
