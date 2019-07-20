package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VzfilterClassName = "vzFilter"

type Filter struct {
	BaseAttributes
    FilterAttributes 
}
  
type FilterAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewFilter(vzFilterRn, parentDn, description string, vzFilterattr FilterAttributes) *Filter {
	dn := fmt.Sprintf("%s/%s", parentDn, vzFilterRn)  
	return &Filter{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzfilterClassName,
			Rn:                vzFilterRn,
		},
        
		FilterAttributes: vzFilterattr,
         
	}
}

func (vzFilter *Filter) ToMap() (map[string]string, error) {
	vzFilterMap, err := vzFilter.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(vzFilterMap, "name",vzFilter.Name)
	
	
    
	A(vzFilterMap, "annotation",vzFilter.Annotation)
	
    
	A(vzFilterMap, "nameAlias",vzFilter.NameAlias)
	
    
	

	return vzFilterMap, err
}

func FilterFromContainerList(cont *container.Container, index int) *Filter {

	FilterCont := cont.S("imdata").Index(index).S(VzfilterClassName, "attributes")
	return &Filter{
		BaseAttributes{
			DistinguishedName: G(FilterCont, "dn"),
			Description:       G(FilterCont, "descr"),
			Status:            G(FilterCont, "status"),
			ClassName:         VzfilterClassName,
			Rn:                G(FilterCont, "rn"),
		},
        
		FilterAttributes{
		
		
			Name : G(FilterCont, "name"),
		
		
        
	        Annotation : G(FilterCont, "annotation"),
		
        
	        NameAlias : G(FilterCont, "nameAlias"),
		
        		
        },
        
	}
}

func FilterFromContainer(cont *container.Container) *Filter {

	return FilterFromContainerList(cont, 0)
}

func FilterListFromContainer(cont *container.Container) []*Filter {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Filter, length)

	for i := 0; i < length; i++ {

		arr[i] = FilterFromContainerList(cont, i)
	}

	return arr
}