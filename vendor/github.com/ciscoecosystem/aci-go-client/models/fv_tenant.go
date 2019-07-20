package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvtenantClassName = "fvTenant"

type Tenant struct {
	BaseAttributes
    TenantAttributes 
}
  
type TenantAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewTenant(fvTenantRn, parentDn, description string, fvTenantattr TenantAttributes) *Tenant {
	dn := fmt.Sprintf("%s/%s", parentDn, fvTenantRn)  
	return &Tenant{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvtenantClassName,
			Rn:                fvTenantRn,
		},
        
		TenantAttributes: fvTenantattr,
         
	}
}

func (fvTenant *Tenant) ToMap() (map[string]string, error) {
	fvTenantMap, err := fvTenant.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(fvTenantMap, "name",fvTenant.Name)
	
	
    
	A(fvTenantMap, "annotation",fvTenant.Annotation)
	
    
	A(fvTenantMap, "nameAlias",fvTenant.NameAlias)
	
    
	

	return fvTenantMap, err
}

func TenantFromContainerList(cont *container.Container, index int) *Tenant {

	TenantCont := cont.S("imdata").Index(index).S(FvtenantClassName, "attributes")
	return &Tenant{
		BaseAttributes{
			DistinguishedName: G(TenantCont, "dn"),
			Description:       G(TenantCont, "descr"),
			Status:            G(TenantCont, "status"),
			ClassName:         FvtenantClassName,
			Rn:                G(TenantCont, "rn"),
		},
        
		TenantAttributes{
		
		
			Name : G(TenantCont, "name"),
		
		
        
	        Annotation : G(TenantCont, "annotation"),
		
        
	        NameAlias : G(TenantCont, "nameAlias"),
		
        		
        },
        
	}
}

func TenantFromContainer(cont *container.Container) *Tenant {

	return TenantFromContainerList(cont, 0)
}

func TenantListFromContainer(cont *container.Container) []*Tenant {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Tenant, length)

	for i := 0; i < length; i++ {

		arr[i] = TenantFromContainerList(cont, i)
	}

	return arr
}