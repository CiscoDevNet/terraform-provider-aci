package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvrspathattClassName = "fvRsPathAtt"

type StaticPath struct {
	BaseAttributes
    StaticPathAttributes 
}
  
type StaticPathAttributes struct {
	
	
	TDn string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	Encap       string `json:",omitempty"`
	
    
	InstrImedcy       string `json:",omitempty"`
	
    
	Mode       string `json:",omitempty"`
	
    
	PrimaryEncap       string `json:",omitempty"`
	
    	
    
}
   

func NewStaticPath(fvRsPathAttRn, parentDn, description string, fvRsPathAttattr StaticPathAttributes) *StaticPath {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsPathAttRn)  
	return &StaticPath{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvrspathattClassName,
			Rn:                fvRsPathAttRn,
		},
        
		StaticPathAttributes: fvRsPathAttattr,
         
	}
}

func (fvRsPathAtt *StaticPath) ToMap() (map[string]string, error) {
	fvRsPathAttMap, err := fvRsPathAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(fvRsPathAttMap, "tDn",fvRsPathAtt.TDn)
	
	
    
	A(fvRsPathAttMap, "annotation",fvRsPathAtt.Annotation)
	
    
	A(fvRsPathAttMap, "encap",fvRsPathAtt.Encap)
	
    
	A(fvRsPathAttMap, "instrImedcy",fvRsPathAtt.InstrImedcy)
	
    
	A(fvRsPathAttMap, "mode",fvRsPathAtt.Mode)
	
    
	A(fvRsPathAttMap, "primaryEncap",fvRsPathAtt.PrimaryEncap)
	
    
	
    
	

	return fvRsPathAttMap, err
}

func StaticPathFromContainerList(cont *container.Container, index int) *StaticPath {

	StaticPathCont := cont.S("imdata").Index(index).S(FvrspathattClassName, "attributes")
	return &StaticPath{
		BaseAttributes{
			DistinguishedName: G(StaticPathCont, "dn"),
			Description:       G(StaticPathCont, "descr"),
			Status:            G(StaticPathCont, "status"),
			ClassName:         FvrspathattClassName,
			Rn:                G(StaticPathCont, "rn"),
		},
        
		StaticPathAttributes{
		
		
			TDn : G(StaticPathCont, "tDn"),
		
		
        
	        Annotation : G(StaticPathCont, "annotation"),
		
        
	        Encap : G(StaticPathCont, "encap"),
		
        
	        InstrImedcy : G(StaticPathCont, "instrImedcy"),
		
        
	        Mode : G(StaticPathCont, "mode"),
		
        
	        PrimaryEncap : G(StaticPathCont, "primaryEncap"),
		
        
		
        		
        },
        
	}
}

func StaticPathFromContainer(cont *container.Container) *StaticPath {

	return StaticPathFromContainerList(cont, 0)
}

func StaticPathListFromContainer(cont *container.Container) []*StaticPath {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*StaticPath, length)

	for i := 0; i < length; i++ {

		arr[i] = StaticPathFromContainerList(cont, i)
	}

	return arr
}