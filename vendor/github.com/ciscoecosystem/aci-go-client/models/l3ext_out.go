package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extoutClassName = "l3extOut"

type L3Outside struct {
	BaseAttributes
    L3OutsideAttributes 
}
  
type L3OutsideAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	EnforceRtctrl       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	TargetDscp       string `json:",omitempty"`
	
    
}
   

func NewL3Outside(l3extOutRn, parentDn, description string, l3extOutattr L3OutsideAttributes) *L3Outside {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extOutRn)  
	return &L3Outside{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extoutClassName,
			Rn:                l3extOutRn,
		},
        
		L3OutsideAttributes: l3extOutattr,
         
	}
}

func (l3extOut *L3Outside) ToMap() (map[string]string, error) {
	l3extOutMap, err := l3extOut.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(l3extOutMap, "name",l3extOut.Name)
	
	
    
	A(l3extOutMap, "annotation",l3extOut.Annotation)
	
    
	A(l3extOutMap, "enforceRtctrl",l3extOut.EnforceRtctrl)
	
    
	A(l3extOutMap, "nameAlias",l3extOut.NameAlias)
	
    
	A(l3extOutMap, "targetDscp",l3extOut.TargetDscp)
	
    
	

	return l3extOutMap, err
}

func L3OutsideFromContainerList(cont *container.Container, index int) *L3Outside {

	L3OutsideCont := cont.S("imdata").Index(index).S(L3extoutClassName, "attributes")
	return &L3Outside{
		BaseAttributes{
			DistinguishedName: G(L3OutsideCont, "dn"),
			Description:       G(L3OutsideCont, "descr"),
			Status:            G(L3OutsideCont, "status"),
			ClassName:         L3extoutClassName,
			Rn:                G(L3OutsideCont, "rn"),
		},
        
		L3OutsideAttributes{
		
		
			Name : G(L3OutsideCont, "name"),
		
		
        
	        Annotation : G(L3OutsideCont, "annotation"),
		
        
	        EnforceRtctrl : G(L3OutsideCont, "enforceRtctrl"),
		
        
	        NameAlias : G(L3OutsideCont, "nameAlias"),
		
        
	        TargetDscp : G(L3OutsideCont, "targetDscp"),
		
        		
        },
        
	}
}

func L3OutsideFromContainer(cont *container.Container) *L3Outside {

	return L3OutsideFromContainerList(cont, 0)
}

func L3OutsideListFromContainer(cont *container.Container) []*L3Outside {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3Outside, length)

	for i := 0; i < length; i++ {

		arr[i] = L3OutsideFromContainerList(cont, i)
	}

	return arr
}