package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricexplicitgepClassName = "fabricExplicitGEp"

type VPCExplicitProtectionGroup struct {
	BaseAttributes
	VPCExplicitProtectionGroupAttributes
}

type VPCExplicitProtectionGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	VPCExplicitProtectionGroup_id string `json:",omitempty"`
	Switch1                       string `json:",omitempty"`
	Switch2                       string `json:",omitempty"`
	VpcDomainPolicy               string `json:",omitempty"`
}

func NewVPCExplicitProtectionGroup(fabricExplicitGEpRn, parentDn string, fabricExplicitGEpattr VPCExplicitProtectionGroupAttributes) *VPCExplicitProtectionGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricExplicitGEpRn)
	return &VPCExplicitProtectionGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FabricexplicitgepClassName,
			Rn:                fabricExplicitGEpRn,
		},

		VPCExplicitProtectionGroupAttributes: fabricExplicitGEpattr,
	}
}

func (fabricExplicitGEp *VPCExplicitProtectionGroup) ToMap() (map[string]string, error) {
	fabricExplicitGEpMap, err := fabricExplicitGEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricExplicitGEpMap, "name", fabricExplicitGEp.Name)

	A(fabricExplicitGEpMap, "annotation", fabricExplicitGEp.Annotation)

	A(fabricExplicitGEpMap, "id", fabricExplicitGEp.VPCExplicitProtectionGroup_id)
	A(fabricExplicitGEpMap, "switch1", fabricExplicitGEp.Switch1)
	A(fabricExplicitGEpMap, "switch2", fabricExplicitGEp.Switch2)
	A(fabricExplicitGEpMap, "vpc_domain_policy", fabricExplicitGEp.VpcDomainPolicy)

	return fabricExplicitGEpMap, err
}

func VPCExplicitProtectionGroupFromContainerList(cont *container.Container, index int) *VPCExplicitProtectionGroup {

	VPCExplicitProtectionGroupCont := cont.S("imdata").Index(index).S(FabricexplicitgepClassName, "attributes")
	ChildContList, err := cont.S("imdata").Index(index).S(FabricexplicitgepClassName, "children").Children()
	if err != nil {
		return nil
	}
	Switch1 := ""
	Switch2 := ""
	VpcDomainPolicy := ""
	for _, childCont := range ChildContList {
		if childCont.Exists("fabricNodePEp") && Switch2 == "" {
			Switch2 = G(childCont.S("fabricNodePEp", "attributes"), "id")

		} else if childCont.Exists("fabricNodePEp") && Switch1 == "" {
			Switch1 = G(childCont.S("fabricNodePEp", "attributes"), "id")
		} else if childCont.Exists("fabricRsVpcInstPol") {
			VpcDomainPolicy = G(childCont.S("fabricRsVpcInstPol", "attributes"), "tnVpcInstPolName")
		}

	}
	return &VPCExplicitProtectionGroup{
		BaseAttributes{
			DistinguishedName: G(VPCExplicitProtectionGroupCont, "dn"),
			Status:            G(VPCExplicitProtectionGroupCont, "status"),
			ClassName:         FabricexplicitgepClassName,
			Rn:                G(VPCExplicitProtectionGroupCont, "rn"),
		},

		VPCExplicitProtectionGroupAttributes{

			Name: G(VPCExplicitProtectionGroupCont, "name"),

			Annotation: G(VPCExplicitProtectionGroupCont, "annotation"),

			VPCExplicitProtectionGroup_id: G(VPCExplicitProtectionGroupCont, "id"),
			Switch1:                       Switch1,
			Switch2:                       Switch2,
			VpcDomainPolicy:               VpcDomainPolicy,
		},
	}
}

func VPCExplicitProtectionGroupFromContainer(cont *container.Container) *VPCExplicitProtectionGroup {

	return VPCExplicitProtectionGroupFromContainerList(cont, 0)
}

func VPCExplicitProtectionGroupListFromContainer(cont *container.Container) []*VPCExplicitProtectionGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VPCExplicitProtectionGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = VPCExplicitProtectionGroupFromContainerList(cont, i)
	}

	return arr
}
