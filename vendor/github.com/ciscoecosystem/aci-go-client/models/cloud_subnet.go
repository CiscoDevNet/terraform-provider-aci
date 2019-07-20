package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudsubnetClassName = "cloudSubnet"

type CloudSubnet struct {
	BaseAttributes
    CloudSubnetAttributes 
}
  
type CloudSubnetAttributes struct {
	
	
	Ip string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
    
	NameAlias       string `json:",omitempty"`
	
    
	Scope       string `json:",omitempty"`
	
    
	Usage       string `json:",omitempty"`
	
    
}
   

func NewCloudSubnet(cloudSubnetRn, parentDn, description string, cloudSubnetattr CloudSubnetAttributes) *CloudSubnet {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudSubnetRn)  
	return &CloudSubnet{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudsubnetClassName,
			Rn:                cloudSubnetRn,
		},
        
		CloudSubnetAttributes: cloudSubnetattr,
         
	}
}

func (cloudSubnet *CloudSubnet) ToMap() (map[string]string, error) {
	cloudSubnetMap, err := cloudSubnet.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(cloudSubnetMap, "ip",cloudSubnet.Ip)
	
	
    
	A(cloudSubnetMap, "annotation",cloudSubnet.Annotation)
	
    
    
	A(cloudSubnetMap, "nameAlias",cloudSubnet.NameAlias)
	
    
	A(cloudSubnetMap, "scope",cloudSubnet.Scope)
	
    
	A(cloudSubnetMap, "usage",cloudSubnet.Usage)
	
    
	

	return cloudSubnetMap, err
}

func CloudSubnetFromContainerList(cont *container.Container, index int) *CloudSubnet {

	CloudSubnetCont := cont.S("imdata").Index(index).S(CloudsubnetClassName, "attributes")
	return &CloudSubnet{
		BaseAttributes{
			DistinguishedName: G(CloudSubnetCont, "dn"),
			Description:       G(CloudSubnetCont, "descr"),
			Status:            G(CloudSubnetCont, "status"),
			ClassName:         CloudsubnetClassName,
			Rn:                G(CloudSubnetCont, "rn"),
		},
        
		CloudSubnetAttributes{
		
		
			Ip : G(CloudSubnetCont, "ip"),
		
		
        
	        Annotation : G(CloudSubnetCont, "annotation"),
		
        
        
	        NameAlias : G(CloudSubnetCont, "nameAlias"),
		
        
	        Scope : G(CloudSubnetCont, "scope"),
		
        
	        Usage : G(CloudSubnetCont, "usage"),
		
        		
        },
        
	}
}

func CloudSubnetFromContainer(cont *container.Container) *CloudSubnet {

	return CloudSubnetFromContainerList(cont, 0)
}

func CloudSubnetListFromContainer(cont *container.Container) []*CloudSubnet {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudSubnet, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudSubnetFromContainerList(cont, i)
	}

	return arr
}