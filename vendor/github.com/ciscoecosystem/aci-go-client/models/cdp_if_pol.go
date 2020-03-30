package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CdpifpolClassName = "cdpIfPol"

type CDPInterfacePolicy struct {
	BaseAttributes
	CDPInterfacePolicyAttributes
}

type CDPInterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewCDPInterfacePolicy(cdpIfPolRn, parentDn, description string, cdpIfPolattr CDPInterfacePolicyAttributes) *CDPInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, cdpIfPolRn)
	return &CDPInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CdpifpolClassName,
			Rn:                cdpIfPolRn,
		},

		CDPInterfacePolicyAttributes: cdpIfPolattr,
	}
}

func (cdpIfPol *CDPInterfacePolicy) ToMap() (map[string]string, error) {
	cdpIfPolMap, err := cdpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cdpIfPolMap, "name", cdpIfPol.Name)

	A(cdpIfPolMap, "adminSt", cdpIfPol.AdminSt)

	A(cdpIfPolMap, "annotation", cdpIfPol.Annotation)

	A(cdpIfPolMap, "nameAlias", cdpIfPol.NameAlias)

	return cdpIfPolMap, err
}

func CDPInterfacePolicyFromContainerList(cont *container.Container, index int) *CDPInterfacePolicy {

	CDPInterfacePolicyCont := cont.S("imdata").Index(index).S(CdpifpolClassName, "attributes")
	return &CDPInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(CDPInterfacePolicyCont, "dn"),
			Description:       G(CDPInterfacePolicyCont, "descr"),
			Status:            G(CDPInterfacePolicyCont, "status"),
			ClassName:         CdpifpolClassName,
			Rn:                G(CDPInterfacePolicyCont, "rn"),
		},

		CDPInterfacePolicyAttributes{

			Name: G(CDPInterfacePolicyCont, "name"),

			AdminSt: G(CDPInterfacePolicyCont, "adminSt"),

			Annotation: G(CDPInterfacePolicyCont, "annotation"),

			NameAlias: G(CDPInterfacePolicyCont, "nameAlias"),
		},
	}
}

func CDPInterfacePolicyFromContainer(cont *container.Container) *CDPInterfacePolicy {

	return CDPInterfacePolicyFromContainerList(cont, 0)
}

func CDPInterfacePolicyListFromContainer(cont *container.Container) []*CDPInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CDPInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = CDPInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}
