package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraaccbndlgrpClassName = "infraAccBndlGrp"

type PCVPCInterfacePolicyGroup struct {
	BaseAttributes
    PCVPCInterfacePolicyGroupAttributes 
}
  
type PCVPCInterfacePolicyGroupAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	LagT       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewPCVPCInterfacePolicyGroup(infraAccBndlGrpRn, parentDn, description string, infraAccBndlGrpattr PCVPCInterfacePolicyGroupAttributes) *PCVPCInterfacePolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAccBndlGrpRn)  
	return &PCVPCInterfacePolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraaccbndlgrpClassName,
			Rn:                infraAccBndlGrpRn,
		},
        
		PCVPCInterfacePolicyGroupAttributes: infraAccBndlGrpattr,
         
	}
}

func (infraAccBndlGrp *PCVPCInterfacePolicyGroup) ToMap() (map[string]string, error) {
	infraAccBndlGrpMap, err := infraAccBndlGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(infraAccBndlGrpMap, "name",infraAccBndlGrp.Name)
	
	
    
	A(infraAccBndlGrpMap, "annotation",infraAccBndlGrp.Annotation)
	
    
	A(infraAccBndlGrpMap, "lagT",infraAccBndlGrp.LagT)
	
    
	A(infraAccBndlGrpMap, "nameAlias",infraAccBndlGrp.NameAlias)
	
    
	

	return infraAccBndlGrpMap, err
}

func PCVPCInterfacePolicyGroupFromContainerList(cont *container.Container, index int) *PCVPCInterfacePolicyGroup {

	PCVPCInterfacePolicyGroupCont := cont.S("imdata").Index(index).S(InfraaccbndlgrpClassName, "attributes")
	return &PCVPCInterfacePolicyGroup{
		BaseAttributes{
			DistinguishedName: G(PCVPCInterfacePolicyGroupCont, "dn"),
			Description:       G(PCVPCInterfacePolicyGroupCont, "descr"),
			Status:            G(PCVPCInterfacePolicyGroupCont, "status"),
			ClassName:         InfraaccbndlgrpClassName,
			Rn:                G(PCVPCInterfacePolicyGroupCont, "rn"),
		},
        
		PCVPCInterfacePolicyGroupAttributes{
		
		
			Name : G(PCVPCInterfacePolicyGroupCont, "name"),
		
		
        
	        Annotation : G(PCVPCInterfacePolicyGroupCont, "annotation"),
		
        
	        LagT : G(PCVPCInterfacePolicyGroupCont, "lagT"),
		
        
	        NameAlias : G(PCVPCInterfacePolicyGroupCont, "nameAlias"),
		
        		
        },
        
	}
}

func PCVPCInterfacePolicyGroupFromContainer(cont *container.Container) *PCVPCInterfacePolicyGroup {

	return PCVPCInterfacePolicyGroupFromContainerList(cont, 0)
}

func PCVPCInterfacePolicyGroupListFromContainer(cont *container.Container) []*PCVPCInterfacePolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*PCVPCInterfacePolicyGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = PCVPCInterfacePolicyGroupFromContainerList(cont, i)
	}

	return arr
}