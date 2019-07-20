package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudepselectorClassName = "cloudEPSelector"

type CloudEndpointSelector struct {
	BaseAttributes
    CloudEndpointSelectorAttributes 
}
  
type CloudEndpointSelectorAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	MatchExpression       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewCloudEndpointSelector(cloudEPSelectorRn, parentDn, description string, cloudEPSelectorattr CloudEndpointSelectorAttributes) *CloudEndpointSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudEPSelectorRn)  
	return &CloudEndpointSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudepselectorClassName,
			Rn:                cloudEPSelectorRn,
		},
        
		CloudEndpointSelectorAttributes: cloudEPSelectorattr,
         
	}
}

func (cloudEPSelector *CloudEndpointSelector) ToMap() (map[string]string, error) {
	cloudEPSelectorMap, err := cloudEPSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(cloudEPSelectorMap, "name",cloudEPSelector.Name)
	
	
    
	A(cloudEPSelectorMap, "annotation",cloudEPSelector.Annotation)
	
    
	A(cloudEPSelectorMap, "matchExpression",cloudEPSelector.MatchExpression)
	
    
	A(cloudEPSelectorMap, "nameAlias",cloudEPSelector.NameAlias)
	
    
	

	return cloudEPSelectorMap, err
}

func CloudEndpointSelectorFromContainerList(cont *container.Container, index int) *CloudEndpointSelector {

	CloudEndpointSelectorCont := cont.S("imdata").Index(index).S(CloudepselectorClassName, "attributes")
	return &CloudEndpointSelector{
		BaseAttributes{
			DistinguishedName: G(CloudEndpointSelectorCont, "dn"),
			Description:       G(CloudEndpointSelectorCont, "descr"),
			Status:            G(CloudEndpointSelectorCont, "status"),
			ClassName:         CloudepselectorClassName,
			Rn:                G(CloudEndpointSelectorCont, "rn"),
		},
        
		CloudEndpointSelectorAttributes{
		
		
			Name : G(CloudEndpointSelectorCont, "name"),
		
		
        
	        Annotation : G(CloudEndpointSelectorCont, "annotation"),
		
        
	        MatchExpression : G(CloudEndpointSelectorCont, "matchExpression"),
		
        
	        NameAlias : G(CloudEndpointSelectorCont, "nameAlias"),
		
        		
        },
        
	}
}

func CloudEndpointSelectorFromContainer(cont *container.Container) *CloudEndpointSelector {

	return CloudEndpointSelectorFromContainerList(cont, 0)
}

func CloudEndpointSelectorListFromContainer(cont *container.Container) []*CloudEndpointSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudEndpointSelector, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudEndpointSelectorFromContainerList(cont, i)
	}

	return arr
}