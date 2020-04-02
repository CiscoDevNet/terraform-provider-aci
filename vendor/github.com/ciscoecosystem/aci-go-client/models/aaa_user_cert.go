package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const AaausercertClassName = "aaaUserCert"

type X509Certificate struct {
	BaseAttributes
    X509CertificateAttributes 
}
  
type X509CertificateAttributes struct {
	
	
	Name string `json:",omitempty"`
	
	
    
	Annotation       string `json:",omitempty"`
	
    
	Data       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
}
   

func NewX509Certificate(aaaUserCertRn, parentDn, description string, aaaUserCertattr X509CertificateAttributes) *X509Certificate {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaUserCertRn)  
	return &X509Certificate{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaausercertClassName,
			Rn:                aaaUserCertRn,
		},
        
		X509CertificateAttributes: aaaUserCertattr,
         
	}
}

func (aaaUserCert *X509Certificate) ToMap() (map[string]string, error) {
	aaaUserCertMap, err := aaaUserCert.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
	
	A(aaaUserCertMap, "name",aaaUserCert.Name)
	
	
    
	A(aaaUserCertMap, "annotation",aaaUserCert.Annotation)
	
    
	A(aaaUserCertMap, "data",aaaUserCert.Data)
	
    
	A(aaaUserCertMap, "nameAlias",aaaUserCert.NameAlias)
	
    
	

	return aaaUserCertMap, err
}

func X509CertificateFromContainerList(cont *container.Container, index int) *X509Certificate {

	X509CertificateCont := cont.S("imdata").Index(index).S(AaausercertClassName, "attributes")
	return &X509Certificate{
		BaseAttributes{
			DistinguishedName: G(X509CertificateCont, "dn"),
			Description:       G(X509CertificateCont, "descr"),
			Status:            G(X509CertificateCont, "status"),
			ClassName:         AaausercertClassName,
			Rn:                G(X509CertificateCont, "rn"),
		},
        
		X509CertificateAttributes{
		
		
			Name : G(X509CertificateCont, "name"),
		
		
        
	        Annotation : G(X509CertificateCont, "annotation"),
		
        
	        Data : G(X509CertificateCont, "data"),
		
        
	        NameAlias : G(X509CertificateCont, "nameAlias"),
		
        		
        },
        
	}
}

func X509CertificateFromContainer(cont *container.Container) *X509Certificate {

	return X509CertificateFromContainerList(cont, 0)
}

func X509CertificateListFromContainer(cont *container.Container) []*X509Certificate {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*X509Certificate, length)

	for i := 0; i < length; i++ {

		arr[i] = X509CertificateFromContainerList(cont, i)
	}

	return arr
}