package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VzsubjClassName = "vzSubj"

type ContractSubject struct {
	BaseAttributes
    ContractSubjectAttributes 
}
  
type ContractSubjectAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	ConsMatchT       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	Prio       string `json:",omitempty"`
	
    
	ProvMatchT       string `json:",omitempty"`
	
    
	RevFltPorts       string `json:",omitempty"`
	
    
	TargetDscp       string `json:",omitempty"`
	
    
}
   

func NewContractSubject(vzSubjRn, parentDn, description string, vzSubjattr ContractSubjectAttributes) *ContractSubject {
	dn := fmt.Sprintf("%s/%s", parentDn, vzSubjRn)  
	return &ContractSubject{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzsubjClassName,
			Rn:                vzSubjRn,
		},
        
		ContractSubjectAttributes: vzSubjattr,
         
	}
}

func (vzSubj *ContractSubject) ToMap() (map[string]string, error) {
	vzSubjMap, err := vzSubj.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(vzSubjMap, "name",vzSubj.Name)
	
	
    
	A(vzSubjMap, "annotation",vzSubj.Annotation)
	
    
	A(vzSubjMap, "consMatchT",vzSubj.ConsMatchT)
	
    
	A(vzSubjMap, "nameAlias",vzSubj.NameAlias)
	
    
	A(vzSubjMap, "prio",vzSubj.Prio)
	
    
	A(vzSubjMap, "provMatchT",vzSubj.ProvMatchT)
	
    
	A(vzSubjMap, "revFltPorts",vzSubj.RevFltPorts)
	
    
	A(vzSubjMap, "targetDscp",vzSubj.TargetDscp)
	
    
	

	return vzSubjMap, err
}

func ContractSubjectFromContainerList(cont *container.Container, index int) *ContractSubject {

	ContractSubjectCont := cont.S("imdata").Index(index).S(VzsubjClassName, "attributes")
	return &ContractSubject{
		BaseAttributes{
			DistinguishedName: G(ContractSubjectCont, "dn"),
			Description:       G(ContractSubjectCont, "descr"),
			Status:            G(ContractSubjectCont, "status"),
			ClassName:         VzsubjClassName,
			Rn:                G(ContractSubjectCont, "rn"),
		},
        
		ContractSubjectAttributes{
		
		
			Name : G(ContractSubjectCont, "name"),
		
		
        
	        Annotation : G(ContractSubjectCont, "annotation"),
		
        
	        ConsMatchT : G(ContractSubjectCont, "consMatchT"),
		
        
	        NameAlias : G(ContractSubjectCont, "nameAlias"),
		
        
	        Prio : G(ContractSubjectCont, "prio"),
		
        
	        ProvMatchT : G(ContractSubjectCont, "provMatchT"),
		
        
	        RevFltPorts : G(ContractSubjectCont, "revFltPorts"),
		
        
	        TargetDscp : G(ContractSubjectCont, "targetDscp"),
		
        		
        },
        
	}
}

func ContractSubjectFromContainer(cont *container.Container) *ContractSubject {

	return ContractSubjectFromContainerList(cont, 0)
}

func ContractSubjectListFromContainer(cont *container.Container) []*ContractSubject {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ContractSubject, length)

	for i := 0; i < length; i++ {

		arr[i] = ContractSubjectFromContainerList(cont, i)
	}

	return arr
}