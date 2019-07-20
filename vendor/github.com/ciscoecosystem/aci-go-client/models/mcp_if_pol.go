package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const McpifpolClassName = "mcpIfPol"

type MiscablingProtocolInterfacePolicy struct {
	BaseAttributes
    MiscablingProtocolInterfacePolicyAttributes 
}
  
type MiscablingProtocolInterfacePolicyAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	AdminSt       string `json:",omitempty"`
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewMiscablingProtocolInterfacePolicy(mcpIfPolRn, parentDn, description string, mcpIfPolattr MiscablingProtocolInterfacePolicyAttributes) *MiscablingProtocolInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, mcpIfPolRn)  
	return &MiscablingProtocolInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         McpifpolClassName,
			Rn:                mcpIfPolRn,
		},
        
		MiscablingProtocolInterfacePolicyAttributes: mcpIfPolattr,
         
	}
}

func (mcpIfPol *MiscablingProtocolInterfacePolicy) ToMap() (map[string]string, error) {
	mcpIfPolMap, err := mcpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(mcpIfPolMap, "name",mcpIfPol.Name)
	
	
    
	A(mcpIfPolMap, "adminSt",mcpIfPol.AdminSt)
	
    
	A(mcpIfPolMap, "annotation",mcpIfPol.Annotation)
	
    
	A(mcpIfPolMap, "nameAlias",mcpIfPol.NameAlias)
	
    
	

	return mcpIfPolMap, err
}

func MiscablingProtocolInterfacePolicyFromContainerList(cont *container.Container, index int) *MiscablingProtocolInterfacePolicy {

	MiscablingProtocolInterfacePolicyCont := cont.S("imdata").Index(index).S(McpifpolClassName, "attributes")
	return &MiscablingProtocolInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(MiscablingProtocolInterfacePolicyCont, "dn"),
			Description:       G(MiscablingProtocolInterfacePolicyCont, "descr"),
			Status:            G(MiscablingProtocolInterfacePolicyCont, "status"),
			ClassName:         McpifpolClassName,
			Rn:                G(MiscablingProtocolInterfacePolicyCont, "rn"),
		},
        
		MiscablingProtocolInterfacePolicyAttributes{
		
		
			Name : G(MiscablingProtocolInterfacePolicyCont, "name"),
		
		
        
	        AdminSt : G(MiscablingProtocolInterfacePolicyCont, "adminSt"),
		
        
	        Annotation : G(MiscablingProtocolInterfacePolicyCont, "annotation"),
		
        
	        NameAlias : G(MiscablingProtocolInterfacePolicyCont, "nameAlias"),
		
        		
        },
        
	}
}

func MiscablingProtocolInterfacePolicyFromContainer(cont *container.Container) *MiscablingProtocolInterfacePolicy {

	return MiscablingProtocolInterfacePolicyFromContainerList(cont, 0)
}

func MiscablingProtocolInterfacePolicyListFromContainer(cont *container.Container) []*MiscablingProtocolInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*MiscablingProtocolInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = MiscablingProtocolInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}