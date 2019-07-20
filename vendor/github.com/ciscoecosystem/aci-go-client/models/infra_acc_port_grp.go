package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraaccportgrpClassName = "infraAccPortGrp"

type LeafAccessPortPolicyGroup struct {
	BaseAttributes
    LeafAccessPortPolicyGroupAttributes 
}
  
type LeafAccessPortPolicyGroupAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewLeafAccessPortPolicyGroup(infraAccPortGrpRn, parentDn, description string, infraAccPortGrpattr LeafAccessPortPolicyGroupAttributes) *LeafAccessPortPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAccPortGrpRn)  
	return &LeafAccessPortPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraaccportgrpClassName,
			Rn:                infraAccPortGrpRn,
		},
        
		LeafAccessPortPolicyGroupAttributes: infraAccPortGrpattr,
         
	}
}

func (infraAccPortGrp *LeafAccessPortPolicyGroup) ToMap() (map[string]string, error) {
	infraAccPortGrpMap, err := infraAccPortGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(infraAccPortGrpMap, "name",infraAccPortGrp.Name)
	
	
    
	A(infraAccPortGrpMap, "annotation",infraAccPortGrp.Annotation)
	
    
	A(infraAccPortGrpMap, "nameAlias",infraAccPortGrp.NameAlias)
	
    
	

	return infraAccPortGrpMap, err
}

func LeafAccessPortPolicyGroupFromContainerList(cont *container.Container, index int) *LeafAccessPortPolicyGroup {

	LeafAccessPortPolicyGroupCont := cont.S("imdata").Index(index).S(InfraaccportgrpClassName, "attributes")
	return &LeafAccessPortPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(LeafAccessPortPolicyGroupCont, "dn"),
			Description:       G(LeafAccessPortPolicyGroupCont, "descr"),
			Status:            G(LeafAccessPortPolicyGroupCont, "status"),
			ClassName:         InfraaccportgrpClassName,
			Rn:                G(LeafAccessPortPolicyGroupCont, "rn"),
		},
        
		LeafAccessPortPolicyGroupAttributes{
		
		
			Name : G(LeafAccessPortPolicyGroupCont, "name"),
		
		
        
	        Annotation : G(LeafAccessPortPolicyGroupCont, "annotation"),
		
        
	        NameAlias : G(LeafAccessPortPolicyGroupCont, "nameAlias"),
		
        		
        },
        
	}
}

func LeafAccessPortPolicyGroupFromContainer(cont *container.Container) *LeafAccessPortPolicyGroup {

	return LeafAccessPortPolicyGroupFromContainerList(cont, 0)
}

func LeafAccessPortPolicyGroupListFromContainer(cont *container.Container) []*LeafAccessPortPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LeafAccessPortPolicyGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = LeafAccessPortPolicyGroupFromContainerList(cont, i)
	}

	return arr
}