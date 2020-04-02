package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const SpansrcgrpClassName = "spanSrcGrp"

type SPANSourceGroup struct {
	BaseAttributes
    SPANSourceGroupAttributes 
}
  
type SPANSourceGroupAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	AdminSt       string `json:",omitempty"`
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewSPANSourceGroup(spanSrcGrpRn, parentDn, description string, spanSrcGrpattr SPANSourceGroupAttributes) *SPANSourceGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, spanSrcGrpRn)  
	return &SPANSourceGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         SpansrcgrpClassName,
			Rn:                spanSrcGrpRn,
		},
        
		SPANSourceGroupAttributes: spanSrcGrpattr,
         
	}
}

func (spanSrcGrp *SPANSourceGroup) ToMap() (map[string]string, error) {
	spanSrcGrpMap, err := spanSrcGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(spanSrcGrpMap, "name",spanSrcGrp.Name)
	
	
    
	A(spanSrcGrpMap, "adminSt",spanSrcGrp.AdminSt)
	
    
	A(spanSrcGrpMap, "annotation",spanSrcGrp.Annotation)
	
    
	A(spanSrcGrpMap, "nameAlias",spanSrcGrp.NameAlias)
	
    
	

	return spanSrcGrpMap, err
}

func SPANSourceGroupFromContainerList(cont *container.Container, index int) *SPANSourceGroup {

	SPANSourceGroupCont := cont.S("imdata").Index(index).S(SpansrcgrpClassName, "attributes")
	return &SPANSourceGroup{
		BaseAttributes{
			DistinguishedName: G(SPANSourceGroupCont, "dn"),
			Description:       G(SPANSourceGroupCont, "descr"),
			Status:            G(SPANSourceGroupCont, "status"),
			ClassName:         SpansrcgrpClassName,
			Rn:                G(SPANSourceGroupCont, "rn"),
		},
        
		SPANSourceGroupAttributes{
		
		
			Name : G(SPANSourceGroupCont, "name"),
		
		
        
	        AdminSt : G(SPANSourceGroupCont, "adminSt"),
		
        
	        Annotation : G(SPANSourceGroupCont, "annotation"),
		
        
	        NameAlias : G(SPANSourceGroupCont, "nameAlias"),
		
        		
        },
        
	}
}

func SPANSourceGroupFromContainer(cont *container.Container) *SPANSourceGroup {

	return SPANSourceGroupFromContainerList(cont, 0)
}

func SPANSourceGroupListFromContainer(cont *container.Container) []*SPANSourceGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SPANSourceGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = SPANSourceGroupFromContainerList(cont, i)
	}

	return arr
}