package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FirmwarefwgrpClassName = "firmwareFwGrp"

type FirmwareGroup struct {
	BaseAttributes
    FirmwareGroupAttributes 
}
  
type FirmwareGroupAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	FirmwareGroup_type       string `json:",omitempty"`
	
    
}
   

func NewFirmwareGroup(firmwareFwGrpRn, parentDn, description string, firmwareFwGrpattr FirmwareGroupAttributes) *FirmwareGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareFwGrpRn)  
	return &FirmwareGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FirmwarefwgrpClassName,
			Rn:                firmwareFwGrpRn,
		},
        
		FirmwareGroupAttributes: firmwareFwGrpattr,
         
	}
}

func (firmwareFwGrp *FirmwareGroup) ToMap() (map[string]string, error) {
	firmwareFwGrpMap, err := firmwareFwGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(firmwareFwGrpMap, "name",firmwareFwGrp.Name)
	
	
    
	A(firmwareFwGrpMap, "annotation",firmwareFwGrp.Annotation)
	
    
	A(firmwareFwGrpMap, "nameAlias",firmwareFwGrp.NameAlias)
	
    
	A(firmwareFwGrpMap, "type",firmwareFwGrp.FirmwareGroup_type)
	
    
	

	return firmwareFwGrpMap, err
}

func FirmwareGroupFromContainerList(cont *container.Container, index int) *FirmwareGroup {

	FirmwareGroupCont := cont.S("imdata").Index(index).S(FirmwarefwgrpClassName, "attributes")
	return &FirmwareGroup{
		BaseAttributes{
			DistinguishedName: G(FirmwareGroupCont, "dn"),
			Description:       G(FirmwareGroupCont, "descr"),
			Status:            G(FirmwareGroupCont, "status"),
			ClassName:         FirmwarefwgrpClassName,
			Rn:                G(FirmwareGroupCont, "rn"),
		},
        
		FirmwareGroupAttributes{
		
		
			Name : G(FirmwareGroupCont, "name"),
		
		
        
	        Annotation : G(FirmwareGroupCont, "annotation"),
		
        
	        NameAlias : G(FirmwareGroupCont, "nameAlias"),
		
        
	        FirmwareGroup_type : G(FirmwareGroupCont, "type"),
		
        		
        },
        
	}
}

func FirmwareGroupFromContainer(cont *container.Container) *FirmwareGroup {

	return FirmwareGroupFromContainerList(cont, 0)
}

func FirmwareGroupListFromContainer(cont *container.Container) []*FirmwareGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FirmwareGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = FirmwareGroupFromContainerList(cont, i)
	}

	return arr
}