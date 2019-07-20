package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvsubnetClassName = "fvSubnet"

type Subnet struct {
	BaseAttributes
    SubnetAttributes 
}
  
type SubnetAttributes struct {
	
	
	Ip string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	Ctrl       string `json:",omitempty"`
	
    
    
	NameAlias       string `json:",omitempty"`
	
    
	Preferred       string `json:",omitempty"`
	
    
	Scope       string `json:",omitempty"`
	
    
	Virtual       string `json:",omitempty"`
	
    
}
   

func NewSubnet(fvSubnetRn, parentDn, description string, fvSubnetattr SubnetAttributes) *Subnet {
	dn := fmt.Sprintf("%s/%s", parentDn, fvSubnetRn)  
	return &Subnet{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvsubnetClassName,
			Rn:                fvSubnetRn,
		},
        
		SubnetAttributes: fvSubnetattr,
         
	}
}

func (fvSubnet *Subnet) ToMap() (map[string]string, error) {
	fvSubnetMap, err := fvSubnet.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(fvSubnetMap, "ip",fvSubnet.Ip)
	
	
    
	A(fvSubnetMap, "annotation",fvSubnet.Annotation)
	
    
	A(fvSubnetMap, "ctrl",fvSubnet.Ctrl)
	
    
    
	A(fvSubnetMap, "nameAlias",fvSubnet.NameAlias)
	
    
	A(fvSubnetMap, "preferred",fvSubnet.Preferred)
	
    
	A(fvSubnetMap, "scope",fvSubnet.Scope)
	
    
	A(fvSubnetMap, "virtual",fvSubnet.Virtual)
	
    
	

	return fvSubnetMap, err
}

func SubnetFromContainerList(cont *container.Container, index int) *Subnet {

	SubnetCont := cont.S("imdata").Index(index).S(FvsubnetClassName, "attributes")
	return &Subnet{
		BaseAttributes{
			DistinguishedName: G(SubnetCont, "dn"),
			Description:       G(SubnetCont, "descr"),
			Status:            G(SubnetCont, "status"),
			ClassName:         FvsubnetClassName,
			Rn:                G(SubnetCont, "rn"),
		},
        
		SubnetAttributes{
		
		
			Ip : G(SubnetCont, "ip"),
		
		
        
	        Annotation : G(SubnetCont, "annotation"),
		
        
	        Ctrl : G(SubnetCont, "ctrl"),
		
        
        
	        NameAlias : G(SubnetCont, "nameAlias"),
		
        
	        Preferred : G(SubnetCont, "preferred"),
		
        
	        Scope : G(SubnetCont, "scope"),
		
        
	        Virtual : G(SubnetCont, "virtual"),
		
        		
        },
        
	}
}

func SubnetFromContainer(cont *container.Container) *Subnet {

	return SubnetFromContainerList(cont, 0)
}

func SubnetListFromContainer(cont *container.Container) []*Subnet {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Subnet, length)

	for i := 0; i < length; i++ {

		arr[i] = SubnetFromContainerList(cont, i)
	}

	return arr
}