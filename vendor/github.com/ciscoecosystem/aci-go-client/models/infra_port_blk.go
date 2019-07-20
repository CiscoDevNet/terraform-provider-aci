package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraportblkClassName = "infraPortBlk"

type AccessPortBlock struct {
	BaseAttributes
    AccessPortBlockAttributes 
}
  
type AccessPortBlockAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	FromCard       string `json:",omitempty"`
	
    
	FromPort       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	ToCard       string `json:",omitempty"`
	
    
	ToPort       string `json:",omitempty"`
	
    
}
   

func NewAccessPortBlock(infraPortBlkRn, parentDn, description string, infraPortBlkattr AccessPortBlockAttributes) *AccessPortBlock {
	dn := fmt.Sprintf("%s/%s", parentDn, infraPortBlkRn)  
	return &AccessPortBlock{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraportblkClassName,
			Rn:                infraPortBlkRn,
		},
        
		AccessPortBlockAttributes: infraPortBlkattr,
         
	}
}

func (infraPortBlk *AccessPortBlock) ToMap() (map[string]string, error) {
	infraPortBlkMap, err := infraPortBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(infraPortBlkMap, "name",infraPortBlk.Name)
	
	
    
	A(infraPortBlkMap, "annotation",infraPortBlk.Annotation)
	
    
	A(infraPortBlkMap, "fromCard",infraPortBlk.FromCard)
	
    
	A(infraPortBlkMap, "fromPort",infraPortBlk.FromPort)
	
    
	A(infraPortBlkMap, "nameAlias",infraPortBlk.NameAlias)
	
    
	A(infraPortBlkMap, "toCard",infraPortBlk.ToCard)
	
    
	A(infraPortBlkMap, "toPort",infraPortBlk.ToPort)
	
    
	

	return infraPortBlkMap, err
}

func AccessPortBlockFromContainerList(cont *container.Container, index int) *AccessPortBlock {

	AccessPortBlockCont := cont.S("imdata").Index(index).S(InfraportblkClassName, "attributes")
	return &AccessPortBlock{
		BaseAttributes{
			DistinguishedName: G(AccessPortBlockCont, "dn"),
			Description:       G(AccessPortBlockCont, "descr"),
			Status:            G(AccessPortBlockCont, "status"),
			ClassName:         InfraportblkClassName,
			Rn:                G(AccessPortBlockCont, "rn"),
		},
        
		AccessPortBlockAttributes{
		
		
			Name : G(AccessPortBlockCont, "name"),
		
		
        
	        Annotation : G(AccessPortBlockCont, "annotation"),
		
        
	        FromCard : G(AccessPortBlockCont, "fromCard"),
		
        
	        FromPort : G(AccessPortBlockCont, "fromPort"),
		
        
	        NameAlias : G(AccessPortBlockCont, "nameAlias"),
		
        
	        ToCard : G(AccessPortBlockCont, "toCard"),
		
        
	        ToPort : G(AccessPortBlockCont, "toPort"),
		
        		
        },
        
	}
}

func AccessPortBlockFromContainer(cont *container.Container) *AccessPortBlock {

	return AccessPortBlockFromContainerList(cont, 0)
}

func AccessPortBlockListFromContainer(cont *container.Container) []*AccessPortBlock {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*AccessPortBlock, length)

	for i := 0; i < length; i++ {

		arr[i] = AccessPortBlockFromContainerList(cont, i)
	}

	return arr
}