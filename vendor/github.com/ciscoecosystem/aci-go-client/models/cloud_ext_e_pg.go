package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudextepgClassName = "cloudExtEPg"

type CloudExternalEPg struct {
	BaseAttributes
    CloudExternalEPgAttributes 
}
  
type CloudExternalEPgAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	ExceptionTag       string `json:",omitempty"`
	
    
	FloodOnEncap       string `json:",omitempty"`
	
    
	MatchT       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	PrefGrMemb       string `json:",omitempty"`
	
    
	Prio       string `json:",omitempty"`
	
    
	RouteReachability       string `json:",omitempty"`
	
    
}
   

func NewCloudExternalEPg(cloudExtEPgRn, parentDn, description string, cloudExtEPgattr CloudExternalEPgAttributes) *CloudExternalEPg {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudExtEPgRn)  
	return &CloudExternalEPg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudextepgClassName,
			Rn:                cloudExtEPgRn,
		},
        
		CloudExternalEPgAttributes: cloudExtEPgattr,
         
	}
}

func (cloudExtEPg *CloudExternalEPg) ToMap() (map[string]string, error) {
	cloudExtEPgMap, err := cloudExtEPg.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(cloudExtEPgMap, "name",cloudExtEPg.Name)
	
	
    
	A(cloudExtEPgMap, "annotation",cloudExtEPg.Annotation)
	
    
	A(cloudExtEPgMap, "exceptionTag",cloudExtEPg.ExceptionTag)
	
    
	A(cloudExtEPgMap, "floodOnEncap",cloudExtEPg.FloodOnEncap)
	
    
	A(cloudExtEPgMap, "matchT",cloudExtEPg.MatchT)
	
    
	A(cloudExtEPgMap, "nameAlias",cloudExtEPg.NameAlias)
	
    
	A(cloudExtEPgMap, "prefGrMemb",cloudExtEPg.PrefGrMemb)
	
    
	A(cloudExtEPgMap, "prio",cloudExtEPg.Prio)
	
    
	A(cloudExtEPgMap, "routeReachability",cloudExtEPg.RouteReachability)
	
    
	

	return cloudExtEPgMap, err
}

func CloudExternalEPgFromContainerList(cont *container.Container, index int) *CloudExternalEPg {

	CloudExternalEPgCont := cont.S("imdata").Index(index).S(CloudextepgClassName, "attributes")
	return &CloudExternalEPg{
		BaseAttributes{
			DistinguishedName: G(CloudExternalEPgCont, "dn"),
			Description:       G(CloudExternalEPgCont, "descr"),
			Status:            G(CloudExternalEPgCont, "status"),
			ClassName:         CloudextepgClassName,
			Rn:                G(CloudExternalEPgCont, "rn"),
		},
        
		CloudExternalEPgAttributes{
		
		
			Name : G(CloudExternalEPgCont, "name"),
		
		
        
	        Annotation : G(CloudExternalEPgCont, "annotation"),
		
        
	        ExceptionTag : G(CloudExternalEPgCont, "exceptionTag"),
		
        
	        FloodOnEncap : G(CloudExternalEPgCont, "floodOnEncap"),
		
        
	        MatchT : G(CloudExternalEPgCont, "matchT"),
		
        
	        NameAlias : G(CloudExternalEPgCont, "nameAlias"),
		
        
	        PrefGrMemb : G(CloudExternalEPgCont, "prefGrMemb"),
		
        
	        Prio : G(CloudExternalEPgCont, "prio"),
		
        
	        RouteReachability : G(CloudExternalEPgCont, "routeReachability"),
		
        		
        },
        
	}
}

func CloudExternalEPgFromContainer(cont *container.Container) *CloudExternalEPg {

	return CloudExternalEPgFromContainerList(cont, 0)
}

func CloudExternalEPgListFromContainer(cont *container.Container) []*CloudExternalEPg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudExternalEPg, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudExternalEPgFromContainerList(cont, i)
	}

	return arr
}