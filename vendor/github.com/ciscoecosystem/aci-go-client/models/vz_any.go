package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VzanyClassName = "vzAny"

type Any struct {
	BaseAttributes
    AnyAttributes 
}
  
type AnyAttributes struct {
	
    
	Annotation       string `json:",omitempty"`
	
    
	MatchT       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	PrefGrMemb       string `json:",omitempty"`
	
    
}
   

func NewAny(vzAnyRn, parentDn, description string, vzAnyattr AnyAttributes) *Any {
	dn := fmt.Sprintf("%s/%s", parentDn, vzAnyRn)  
	return &Any{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzanyClassName,
			Rn:                vzAnyRn,
		},
        
		AnyAttributes: vzAnyattr,
         
	}
}

func (vzAny *Any) ToMap() (map[string]string, error) {
	vzAnyMap, err := vzAny.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
    
	A(vzAnyMap, "annotation",vzAny.Annotation)
	
    
	A(vzAnyMap, "matchT",vzAny.MatchT)
	
    
	A(vzAnyMap, "nameAlias",vzAny.NameAlias)
	
    
	A(vzAnyMap, "prefGrMemb",vzAny.PrefGrMemb)
	
    
	

	return vzAnyMap, err
}

func AnyFromContainerList(cont *container.Container, index int) *Any {

	AnyCont := cont.S("imdata").Index(index).S(VzanyClassName, "attributes")
	return &Any{
		BaseAttributes{
			DistinguishedName: G(AnyCont, "dn"),
			Description:       G(AnyCont, "descr"),
			Status:            G(AnyCont, "status"),
			ClassName:         VzanyClassName,
			Rn:                G(AnyCont, "rn"),
		},
        
		AnyAttributes{
		
        
	        Annotation : G(AnyCont, "annotation"),
		
        
	        MatchT : G(AnyCont, "matchT"),
		
        
	        NameAlias : G(AnyCont, "nameAlias"),
		
        
	        PrefGrMemb : G(AnyCont, "prefGrMemb"),
		
        		
        },
        
	}
}

func AnyFromContainer(cont *container.Container) *Any {

	return AnyFromContainerList(cont, 0)
}

func AnyListFromContainer(cont *container.Container) []*Any {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Any, length)

	for i := 0; i < length; i++ {

		arr[i] = AnyFromContainerList(cont, i)
	}

	return arr
}