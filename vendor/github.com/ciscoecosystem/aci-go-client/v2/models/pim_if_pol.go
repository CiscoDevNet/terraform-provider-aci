package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimIfPol        = "pimifpol-%s"
	DnPimIfPol        = "uni/tn-%s/pimifpol-%s"
	ParentDnPimIfPol  = "uni/tn-%s"
	PimIfPolClassName = "pimIfPol"
)

type PIMInterfacePolicy struct {
	BaseAttributes
	PIMInterfacePolicyAttributes
}

type PIMInterfacePolicyAttributes struct {
	Annotation    string `json:",omitempty"`
	AuthKey       string `json:",omitempty"`
	AuthT         string `json:",omitempty"`
	Ctrl          string `json:",omitempty"`
	DrDelay       string `json:",omitempty"`
	DrPrio        string `json:",omitempty"`
	HelloItvl     string `json:",omitempty"`
	JpInterval    string `json:",omitempty"`
	Name          string `json:",omitempty"`
	NameAlias     string `json:",omitempty"`
	SecureAuthKey string `json:",omitempty"`
}

func NewPIMInterfacePolicy(pimIfPolRn, parentDn, description string, pimIfPolAttr PIMInterfacePolicyAttributes) *PIMInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, pimIfPolRn)
	return &PIMInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimIfPolClassName,
			Rn:                pimIfPolRn,
		},
		PIMInterfacePolicyAttributes: pimIfPolAttr,
	}
}

func (pimIfPol *PIMInterfacePolicy) ToMap() (map[string]string, error) {
	pimIfPolMap, err := pimIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimIfPolMap, "annotation", pimIfPol.Annotation)
	A(pimIfPolMap, "authKey", pimIfPol.AuthKey)
	A(pimIfPolMap, "authT", pimIfPol.AuthT)
	A(pimIfPolMap, "ctrl", pimIfPol.Ctrl)
	A(pimIfPolMap, "drDelay", pimIfPol.DrDelay)
	A(pimIfPolMap, "drPrio", pimIfPol.DrPrio)
	A(pimIfPolMap, "helloItvl", pimIfPol.HelloItvl)
	A(pimIfPolMap, "jpInterval", pimIfPol.JpInterval)
	A(pimIfPolMap, "name", pimIfPol.Name)
	A(pimIfPolMap, "nameAlias", pimIfPol.NameAlias)
	A(pimIfPolMap, "secureAuthKey", pimIfPol.SecureAuthKey)
	return pimIfPolMap, err
}

func PIMInterfacePolicyFromContainerList(cont *container.Container, index int) *PIMInterfacePolicy {
	PIMInterfacePolicyCont := cont.S("imdata").Index(index).S(PimIfPolClassName, "attributes")
	return &PIMInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(PIMInterfacePolicyCont, "dn"),
			Description:       G(PIMInterfacePolicyCont, "descr"),
			Status:            G(PIMInterfacePolicyCont, "status"),
			ClassName:         PimIfPolClassName,
			Rn:                G(PIMInterfacePolicyCont, "rn"),
		},
		PIMInterfacePolicyAttributes{
			Annotation:    G(PIMInterfacePolicyCont, "annotation"),
			AuthKey:       G(PIMInterfacePolicyCont, "authKey"),
			AuthT:         G(PIMInterfacePolicyCont, "authT"),
			Ctrl:          G(PIMInterfacePolicyCont, "ctrl"),
			DrDelay:       G(PIMInterfacePolicyCont, "drDelay"),
			DrPrio:        G(PIMInterfacePolicyCont, "drPrio"),
			HelloItvl:     G(PIMInterfacePolicyCont, "helloItvl"),
			JpInterval:    G(PIMInterfacePolicyCont, "jpInterval"),
			Name:          G(PIMInterfacePolicyCont, "name"),
			NameAlias:     G(PIMInterfacePolicyCont, "nameAlias"),
			SecureAuthKey: G(PIMInterfacePolicyCont, "secureAuthKey"),
		},
	}
}

func PIMInterfacePolicyFromContainer(cont *container.Container) *PIMInterfacePolicy {
	return PIMInterfacePolicyFromContainerList(cont, 0)
}

func PIMInterfacePolicyListFromContainer(cont *container.Container) []*PIMInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMInterfacePolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}
