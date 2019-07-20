package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraprovaccClassName = "infraProvAcc"

type VlanEncapsulationforVxlanTraffic struct {
	BaseAttributes
    VlanEncapsulationforVxlanTrafficAttributes 
}
  
type VlanEncapsulationforVxlanTrafficAttributes struct {
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewVlanEncapsulationforVxlanTraffic(infraProvAccRn, parentDn, description string, infraProvAccattr VlanEncapsulationforVxlanTrafficAttributes) *VlanEncapsulationforVxlanTraffic {
	dn := fmt.Sprintf("%s/%s", parentDn, infraProvAccRn)  
	return &VlanEncapsulationforVxlanTraffic{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraprovaccClassName,
			Rn:                infraProvAccRn,
		},
        
		VlanEncapsulationforVxlanTrafficAttributes: infraProvAccattr,
         
	}
}

func (infraProvAcc *VlanEncapsulationforVxlanTraffic) ToMap() (map[string]string, error) {
	infraProvAccMap, err := infraProvAcc.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
    
	A(infraProvAccMap, "annotation",infraProvAcc.Annotation)
	
    
	A(infraProvAccMap, "nameAlias",infraProvAcc.NameAlias)
	
    
	

	return infraProvAccMap, err
}

func VlanEncapsulationforVxlanTrafficFromContainerList(cont *container.Container, index int) *VlanEncapsulationforVxlanTraffic {

	VlanEncapsulationforVxlanTrafficCont := cont.S("imdata").Index(index).S(InfraprovaccClassName, "attributes")
	return &VlanEncapsulationforVxlanTraffic{
		BaseAttributes{
			DistinguishedName: G(VlanEncapsulationforVxlanTrafficCont, "dn"),
			Description:       G(VlanEncapsulationforVxlanTrafficCont, "descr"),
			Status:            G(VlanEncapsulationforVxlanTrafficCont, "status"),
			ClassName:         InfraprovaccClassName,
			Rn:                G(VlanEncapsulationforVxlanTrafficCont, "rn"),
		},
        
		VlanEncapsulationforVxlanTrafficAttributes{
		
        
	        Annotation : G(VlanEncapsulationforVxlanTrafficCont, "annotation"),
		
        
	        NameAlias : G(VlanEncapsulationforVxlanTrafficCont, "nameAlias"),
		
        		
        },
        
	}
}

func VlanEncapsulationforVxlanTrafficFromContainer(cont *container.Container) *VlanEncapsulationforVxlanTraffic {

	return VlanEncapsulationforVxlanTrafficFromContainerList(cont, 0)
}

func VlanEncapsulationforVxlanTrafficListFromContainer(cont *container.Container) []*VlanEncapsulationforVxlanTraffic {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VlanEncapsulationforVxlanTraffic, length)

	for i := 0; i < length; i++ {

		arr[i] = VlanEncapsulationforVxlanTrafficFromContainerList(cont, i)
	}

	return arr
}