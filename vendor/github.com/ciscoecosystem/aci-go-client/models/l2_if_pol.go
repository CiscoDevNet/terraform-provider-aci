package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L2ifpolClassName = "l2IfPol"

type L2InterfacePolicy struct {
	BaseAttributes
    L2InterfacePolicyAttributes 
}
  
type L2InterfacePolicyAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	Qinq       string `json:",omitempty"`
	
    
	Vepa       string `json:",omitempty"`
	
    
	VlanScope       string `json:",omitempty"`
	
    
}
   

func NewL2InterfacePolicy(l2IfPolRn, parentDn, description string, l2IfPolattr L2InterfacePolicyAttributes) *L2InterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l2IfPolRn)  
	return &L2InterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L2ifpolClassName,
			Rn:                l2IfPolRn,
		},
        
		L2InterfacePolicyAttributes: l2IfPolattr,
         
	}
}

func (l2IfPol *L2InterfacePolicy) ToMap() (map[string]string, error) {
	l2IfPolMap, err := l2IfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(l2IfPolMap, "name",l2IfPol.Name)
	
	
    
	A(l2IfPolMap, "annotation",l2IfPol.Annotation)
	
    
	A(l2IfPolMap, "nameAlias",l2IfPol.NameAlias)
	
    
	A(l2IfPolMap, "qinq",l2IfPol.Qinq)
	
    
	A(l2IfPolMap, "vepa",l2IfPol.Vepa)
	
    
	A(l2IfPolMap, "vlanScope",l2IfPol.VlanScope)
	
    
	

	return l2IfPolMap, err
}

func L2InterfacePolicyFromContainerList(cont *container.Container, index int) *L2InterfacePolicy {

	L2InterfacePolicyCont := cont.S("imdata").Index(index).S(L2ifpolClassName, "attributes")
	return &L2InterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(L2InterfacePolicyCont, "dn"),
			Description:       G(L2InterfacePolicyCont, "descr"),
			Status:            G(L2InterfacePolicyCont, "status"),
			ClassName:         L2ifpolClassName,
			Rn:                G(L2InterfacePolicyCont, "rn"),
		},
        
		L2InterfacePolicyAttributes{
		
		
			Name : G(L2InterfacePolicyCont, "name"),
		
		
        
	        Annotation : G(L2InterfacePolicyCont, "annotation"),
		
        
	        NameAlias : G(L2InterfacePolicyCont, "nameAlias"),
		
        
	        Qinq : G(L2InterfacePolicyCont, "qinq"),
		
        
	        Vepa : G(L2InterfacePolicyCont, "vepa"),
		
        
	        VlanScope : G(L2InterfacePolicyCont, "vlanScope"),
		
        		
        },
        
	}
}

func L2InterfacePolicyFromContainer(cont *container.Container) *L2InterfacePolicy {

	return L2InterfacePolicyFromContainerList(cont, 0)
}

func L2InterfacePolicyListFromContainer(cont *container.Container) []*L2InterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L2InterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = L2InterfacePolicyFromContainerList(cont, i)
	}

	return arr
}