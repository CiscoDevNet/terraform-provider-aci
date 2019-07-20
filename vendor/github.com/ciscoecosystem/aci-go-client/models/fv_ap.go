package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvapClassName = "fvAp"

type ApplicationProfile struct {
	BaseAttributes
    ApplicationProfileAttributes 
}
  
type ApplicationProfileAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	Prio       string `json:",omitempty"`
	
    
}
   

func NewApplicationProfile(fvApRn, parentDn, description string, fvApattr ApplicationProfileAttributes) *ApplicationProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, fvApRn)  
	return &ApplicationProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvapClassName,
			Rn:                fvApRn,
		},
        
		ApplicationProfileAttributes: fvApattr,
         
	}
}

func (fvAp *ApplicationProfile) ToMap() (map[string]string, error) {
	fvApMap, err := fvAp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(fvApMap, "name",fvAp.Name)
	
	
    
	A(fvApMap, "annotation",fvAp.Annotation)
	
    
	A(fvApMap, "nameAlias",fvAp.NameAlias)
	
    
	A(fvApMap, "prio",fvAp.Prio)
	
    
	

	return fvApMap, err
}

func ApplicationProfileFromContainerList(cont *container.Container, index int) *ApplicationProfile {

	ApplicationProfileCont := cont.S("imdata").Index(index).S(FvapClassName, "attributes")
	return &ApplicationProfile{
		BaseAttributes{
			DistinguishedName: G(ApplicationProfileCont, "dn"),
			Description:       G(ApplicationProfileCont, "descr"),
			Status:            G(ApplicationProfileCont, "status"),
			ClassName:         FvapClassName,
			Rn:                G(ApplicationProfileCont, "rn"),
		},
        
		ApplicationProfileAttributes{
		
		
			Name : G(ApplicationProfileCont, "name"),
		
		
        
	        Annotation : G(ApplicationProfileCont, "annotation"),
		
        
	        NameAlias : G(ApplicationProfileCont, "nameAlias"),
		
        
	        Prio : G(ApplicationProfileCont, "prio"),
		
        		
        },
        
	}
}

func ApplicationProfileFromContainer(cont *container.Container) *ApplicationProfile {

	return ApplicationProfileFromContainerList(cont, 0)
}

func ApplicationProfileListFromContainer(cont *container.Container) []*ApplicationProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ApplicationProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = ApplicationProfileFromContainerList(cont, i)
	}

	return arr
}