package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsrsnodetoldevClassName = "vnsRsNodeToLDev"

type RelationfromaAbsNodetoanLDev struct {
	BaseAttributes
    RelationfromaAbsNodetoanLDevAttributes 
}
  
type RelationfromaAbsNodetoanLDevAttributes struct {
	
    
	Annotation       string `json:",omitempty"`
	
    
	TDn       string `json:",omitempty"`
	
    
}
   

func NewRelationfromaAbsNodetoanLDev(vnsRsNodeToLDevRn, parentDn, description string, vnsRsNodeToLDevattr RelationfromaAbsNodetoanLDevAttributes) *RelationfromaAbsNodetoanLDev {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsRsNodeToLDevRn)  
	return &RelationfromaAbsNodetoanLDev{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsrsnodetoldevClassName,
			Rn:                vnsRsNodeToLDevRn,
		},
        
		RelationfromaAbsNodetoanLDevAttributes: vnsRsNodeToLDevattr,
         
	}
}

func (vnsRsNodeToLDev *RelationfromaAbsNodetoanLDev) ToMap() (map[string]string, error) {
	vnsRsNodeToLDevMap, err := vnsRsNodeToLDev.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
    
	A(vnsRsNodeToLDevMap, "annotation",vnsRsNodeToLDev.Annotation)
	
    
	A(vnsRsNodeToLDevMap, "tDn",vnsRsNodeToLDev.TDn)
	
    
	

	return vnsRsNodeToLDevMap, err
}

func RelationfromaAbsNodetoanLDevFromContainerList(cont *container.Container, index int) *RelationfromaAbsNodetoanLDev {

	RelationfromaAbsNodetoanLDevCont := cont.S("imdata").Index(index).S(VnsrsnodetoldevClassName, "attributes")
	return &RelationfromaAbsNodetoanLDev{
		BaseAttributes{
			DistinguishedName: G(RelationfromaAbsNodetoanLDevCont, "dn"),
			Description:       G(RelationfromaAbsNodetoanLDevCont, "descr"),
			Status:            G(RelationfromaAbsNodetoanLDevCont, "status"),
			ClassName:         VnsrsnodetoldevClassName,
			Rn:                G(RelationfromaAbsNodetoanLDevCont, "rn"),
		},
        
		RelationfromaAbsNodetoanLDevAttributes{
		
        
	        Annotation : G(RelationfromaAbsNodetoanLDevCont, "annotation"),
		
        
	        TDn : G(RelationfromaAbsNodetoanLDevCont, "tDn"),
		
        		
        },
        
	}
}

func RelationfromaAbsNodetoanLDevFromContainer(cont *container.Container) *RelationfromaAbsNodetoanLDev {

	return RelationfromaAbsNodetoanLDevFromContainerList(cont, 0)
}

func RelationfromaAbsNodetoanLDevListFromContainer(cont *container.Container) []*RelationfromaAbsNodetoanLDev {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*RelationfromaAbsNodetoanLDev, length)

	for i := 0; i < length; i++ {

		arr[i] = RelationfromaAbsNodetoanLDevFromContainerList(cont, i)
	}

	return arr
}