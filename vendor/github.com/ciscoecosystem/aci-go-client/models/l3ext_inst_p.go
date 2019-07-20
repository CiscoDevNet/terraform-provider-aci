package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extinstpClassName = "l3extInstP"

type ExternalNetworkInstanceProfile struct {
	BaseAttributes
    ExternalNetworkInstanceProfileAttributes 
}
  
type ExternalNetworkInstanceProfileAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	ExceptionTag       string `json:",omitempty"`
	
    
	FloodOnEncap       string `json:",omitempty"`
	
    
	MatchT       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	PrefGrMemb       string `json:",omitempty"`
	
    
	Prio       string `json:",omitempty"`
	
    
	TargetDscp       string `json:",omitempty"`
	
    
}
   

func NewExternalNetworkInstanceProfile(l3extInstPRn, parentDn, description string, l3extInstPattr ExternalNetworkInstanceProfileAttributes) *ExternalNetworkInstanceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extInstPRn)  
	return &ExternalNetworkInstanceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extinstpClassName,
			Rn:                l3extInstPRn,
		},
        
		ExternalNetworkInstanceProfileAttributes: l3extInstPattr,
         
	}
}

func (l3extInstP *ExternalNetworkInstanceProfile) ToMap() (map[string]string, error) {
	l3extInstPMap, err := l3extInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(l3extInstPMap, "name",l3extInstP.Name)
	
	
    
	A(l3extInstPMap, "annotation",l3extInstP.Annotation)
	
    
	A(l3extInstPMap, "exceptionTag",l3extInstP.ExceptionTag)
	
    
	A(l3extInstPMap, "floodOnEncap",l3extInstP.FloodOnEncap)
	
    
	A(l3extInstPMap, "matchT",l3extInstP.MatchT)
	
    
	A(l3extInstPMap, "nameAlias",l3extInstP.NameAlias)
	
    
	A(l3extInstPMap, "prefGrMemb",l3extInstP.PrefGrMemb)
	
    
	A(l3extInstPMap, "prio",l3extInstP.Prio)
	
    
	A(l3extInstPMap, "targetDscp",l3extInstP.TargetDscp)
	
    
	

	return l3extInstPMap, err
}

func ExternalNetworkInstanceProfileFromContainerList(cont *container.Container, index int) *ExternalNetworkInstanceProfile {

	ExternalNetworkInstanceProfileCont := cont.S("imdata").Index(index).S(L3extinstpClassName, "attributes")
	return &ExternalNetworkInstanceProfile{
		BaseAttributes{
			DistinguishedName: G(ExternalNetworkInstanceProfileCont, "dn"),
			Description:       G(ExternalNetworkInstanceProfileCont, "descr"),
			Status:            G(ExternalNetworkInstanceProfileCont, "status"),
			ClassName:         L3extinstpClassName,
			Rn:                G(ExternalNetworkInstanceProfileCont, "rn"),
		},
        
		ExternalNetworkInstanceProfileAttributes{
		
		
			Name : G(ExternalNetworkInstanceProfileCont, "name"),
		
		
        
	        Annotation : G(ExternalNetworkInstanceProfileCont, "annotation"),
		
        
	        ExceptionTag : G(ExternalNetworkInstanceProfileCont, "exceptionTag"),
		
        
	        FloodOnEncap : G(ExternalNetworkInstanceProfileCont, "floodOnEncap"),
		
        
	        MatchT : G(ExternalNetworkInstanceProfileCont, "matchT"),
		
        
	        NameAlias : G(ExternalNetworkInstanceProfileCont, "nameAlias"),
		
        
	        PrefGrMemb : G(ExternalNetworkInstanceProfileCont, "prefGrMemb"),
		
        
	        Prio : G(ExternalNetworkInstanceProfileCont, "prio"),
		
        
	        TargetDscp : G(ExternalNetworkInstanceProfileCont, "targetDscp"),
		
        		
        },
        
	}
}

func ExternalNetworkInstanceProfileFromContainer(cont *container.Container) *ExternalNetworkInstanceProfile {

	return ExternalNetworkInstanceProfileFromContainerList(cont, 0)
}

func ExternalNetworkInstanceProfileListFromContainer(cont *container.Container) []*ExternalNetworkInstanceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ExternalNetworkInstanceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = ExternalNetworkInstanceProfileFromContainerList(cont, i)
	}

	return arr
}