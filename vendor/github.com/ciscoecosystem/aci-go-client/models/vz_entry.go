package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VzentryClassName = "vzEntry"

type FilterEntry struct {
	BaseAttributes
    FilterEntryAttributes 
}
  
type FilterEntryAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	ApplyToFrag       string `json:",omitempty"`
	
    
	ArpOpc       string `json:",omitempty"`
	
    
	DFromPort       string `json:",omitempty"`
	
    
	DToPort       string `json:",omitempty"`
	
    
	EtherT       string `json:",omitempty"`
	
    
	Icmpv4T       string `json:",omitempty"`
	
    
	Icmpv6T       string `json:",omitempty"`
	
    
	MatchDscp       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	Prot       string `json:",omitempty"`
	
    
	SFromPort       string `json:",omitempty"`
	
    
	SToPort       string `json:",omitempty"`
	
    
	Stateful       string `json:",omitempty"`
	
    
	TcpRules       string `json:",omitempty"`
	
    
}
   

func NewFilterEntry(vzEntryRn, parentDn, description string, vzEntryattr FilterEntryAttributes) *FilterEntry {
	dn := fmt.Sprintf("%s/%s", parentDn, vzEntryRn)  
	return &FilterEntry{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VzentryClassName,
			Rn:                vzEntryRn,
		},
        
		FilterEntryAttributes: vzEntryattr,
         
	}
}

func (vzEntry *FilterEntry) ToMap() (map[string]string, error) {
	vzEntryMap, err := vzEntry.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(vzEntryMap, "name",vzEntry.Name)
	
	
    
	A(vzEntryMap, "annotation",vzEntry.Annotation)
	
    
	A(vzEntryMap, "applyToFrag",vzEntry.ApplyToFrag)
	
    
	A(vzEntryMap, "arpOpc",vzEntry.ArpOpc)
	
    
	A(vzEntryMap, "dFromPort",vzEntry.DFromPort)
	
    
	A(vzEntryMap, "dToPort",vzEntry.DToPort)
	
    
	A(vzEntryMap, "etherT",vzEntry.EtherT)
	
    
	A(vzEntryMap, "icmpv4T",vzEntry.Icmpv4T)
	
    
	A(vzEntryMap, "icmpv6T",vzEntry.Icmpv6T)
	
    
	A(vzEntryMap, "matchDscp",vzEntry.MatchDscp)
	
    
	A(vzEntryMap, "nameAlias",vzEntry.NameAlias)
	
    
	A(vzEntryMap, "prot",vzEntry.Prot)
	
    
	A(vzEntryMap, "sFromPort",vzEntry.SFromPort)
	
    
	A(vzEntryMap, "sToPort",vzEntry.SToPort)
	
    
	A(vzEntryMap, "stateful",vzEntry.Stateful)
	
    
	A(vzEntryMap, "tcpRules",vzEntry.TcpRules)
	
    
	

	return vzEntryMap, err
}

func FilterEntryFromContainerList(cont *container.Container, index int) *FilterEntry {

	FilterEntryCont := cont.S("imdata").Index(index).S(VzentryClassName, "attributes")
	return &FilterEntry{
		BaseAttributes{
			DistinguishedName: G(FilterEntryCont, "dn"),
			Description:       G(FilterEntryCont, "descr"),
			Status:            G(FilterEntryCont, "status"),
			ClassName:         VzentryClassName,
			Rn:                G(FilterEntryCont, "rn"),
		},
        
		FilterEntryAttributes{
		
		
			Name : G(FilterEntryCont, "name"),
		
		
        
	        Annotation : G(FilterEntryCont, "annotation"),
		
        
	        ApplyToFrag : G(FilterEntryCont, "applyToFrag"),
		
        
	        ArpOpc : G(FilterEntryCont, "arpOpc"),
		
        
	        DFromPort : G(FilterEntryCont, "dFromPort"),
		
        
	        DToPort : G(FilterEntryCont, "dToPort"),
		
        
	        EtherT : G(FilterEntryCont, "etherT"),
		
        
	        Icmpv4T : G(FilterEntryCont, "icmpv4T"),
		
        
	        Icmpv6T : G(FilterEntryCont, "icmpv6T"),
		
        
	        MatchDscp : G(FilterEntryCont, "matchDscp"),
		
        
	        NameAlias : G(FilterEntryCont, "nameAlias"),
		
        
	        Prot : G(FilterEntryCont, "prot"),
		
        
	        SFromPort : G(FilterEntryCont, "sFromPort"),
		
        
	        SToPort : G(FilterEntryCont, "sToPort"),
		
        
	        Stateful : G(FilterEntryCont, "stateful"),
		
        
	        TcpRules : G(FilterEntryCont, "tcpRules"),
		
        		
        },
        
	}
}

func FilterEntryFromContainer(cont *container.Container) *FilterEntry {

	return FilterEntryFromContainerList(cont, 0)
}

func FilterEntryListFromContainer(cont *container.Container) []*FilterEntry {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FilterEntry, length)

	for i := 0; i < length; i++ {

		arr[i] = FilterEntryFromContainerList(cont, i)
	}

	return arr
}