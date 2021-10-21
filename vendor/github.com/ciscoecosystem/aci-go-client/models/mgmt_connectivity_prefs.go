package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnmgmtConnectivityPrefs        = "uni/fabric/connectivityPrefs"
	RnmgmtConnectivityPrefs        = "connectivityPrefs"
	ParentDnmgmtConnectivityPrefs  = "uni/fabric"
	MgmtconnectivityprefsClassName = "mgmtConnectivityPrefs"
)

type Mgmtconnectivitypreference struct {
	BaseAttributes
	NameAliasAttribute
	MgmtconnectivitypreferenceAttributes
}

type MgmtconnectivitypreferenceAttributes struct {
	Annotation    string `json:",omitempty"`
	InterfacePref string `json:",omitempty"`
	Name          string `json:",omitempty"`
}

func NewMgmtconnectivitypreference(mgmtConnectivityPrefsRn, parentDn, description, nameAlias string, mgmtConnectivityPrefsAttr MgmtconnectivitypreferenceAttributes) *Mgmtconnectivitypreference {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtConnectivityPrefsRn)
	return &Mgmtconnectivitypreference{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtconnectivityprefsClassName,
			Rn:                mgmtConnectivityPrefsRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MgmtconnectivitypreferenceAttributes: mgmtConnectivityPrefsAttr,
	}
}

func (mgmtConnectivityPrefs *Mgmtconnectivitypreference) ToMap() (map[string]string, error) {
	mgmtConnectivityPrefsMap, err := mgmtConnectivityPrefs.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := mgmtConnectivityPrefs.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(mgmtConnectivityPrefsMap, key, value)
	}
	A(mgmtConnectivityPrefsMap, "annotation", mgmtConnectivityPrefs.Annotation)
	A(mgmtConnectivityPrefsMap, "interfacePref", mgmtConnectivityPrefs.InterfacePref)
	A(mgmtConnectivityPrefsMap, "name", mgmtConnectivityPrefs.Name)
	return mgmtConnectivityPrefsMap, err
}

func MgmtconnectivitypreferenceFromContainerList(cont *container.Container, index int) *Mgmtconnectivitypreference {
	MgmtconnectivitypreferenceCont := cont.S("imdata").Index(index).S(MgmtconnectivityprefsClassName, "attributes")
	return &Mgmtconnectivitypreference{
		BaseAttributes{
			DistinguishedName: G(MgmtconnectivitypreferenceCont, "dn"),
			Description:       G(MgmtconnectivitypreferenceCont, "descr"),
			Status:            G(MgmtconnectivitypreferenceCont, "status"),
			ClassName:         MgmtconnectivityprefsClassName,
			Rn:                G(MgmtconnectivitypreferenceCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MgmtconnectivitypreferenceCont, "nameAlias"),
		},
		MgmtconnectivitypreferenceAttributes{
			Annotation:    G(MgmtconnectivitypreferenceCont, "annotation"),
			InterfacePref: G(MgmtconnectivitypreferenceCont, "interfacePref"),
			Name:          G(MgmtconnectivitypreferenceCont, "name"),
		},
	}
}

func MgmtconnectivitypreferenceFromContainer(cont *container.Container) *Mgmtconnectivitypreference {
	return MgmtconnectivitypreferenceFromContainerList(cont, 0)
}

func MgmtconnectivitypreferenceListFromContainer(cont *container.Container) []*Mgmtconnectivitypreference {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*Mgmtconnectivitypreference, length)
	for i := 0; i < length; i++ {
		arr[i] = MgmtconnectivitypreferenceFromContainerList(cont, i)
	}
	return arr
}
