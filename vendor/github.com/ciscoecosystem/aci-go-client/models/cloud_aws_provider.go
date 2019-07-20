package models


import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudawsproviderClassName = "cloudAwsProvider"

type CloudAWSProvider struct {
	BaseAttributes
    CloudAWSProviderAttributes 
}
  
type CloudAWSProviderAttributes struct {
	
    
	AccessKeyId       string `json:",omitempty"`
	
    
	AccountId       string `json:",omitempty"`
	
    
	Annotation       string `json:",omitempty"`
	
    
	Email       string `json:",omitempty"`
	
    
	HttpProxy       string `json:",omitempty"`
	
    
	IsAccountInOrg       string `json:",omitempty"`
	
    
	IsTrusted       string `json:",omitempty"`
	
    
	NameAlias       string `json:",omitempty"`
	
    
	ProviderId       string `json:",omitempty"`
	
    
	Region       string `json:",omitempty"`
	
    
	SecretAccessKey       string `json:",omitempty"`
	
    
}
   

func NewCloudAWSProvider(cloudAwsProviderRn, parentDn, description string, cloudAwsProviderattr CloudAWSProviderAttributes) *CloudAWSProvider {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudAwsProviderRn)  
	return &CloudAWSProvider{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudawsproviderClassName,
			Rn:                cloudAwsProviderRn,
		},
        
		CloudAWSProviderAttributes: cloudAwsProviderattr,
         
	}
}

func (cloudAwsProvider *CloudAWSProvider) ToMap() (map[string]string, error) {
	cloudAwsProviderMap, err := cloudAwsProvider.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	
    
	A(cloudAwsProviderMap, "accessKeyId",cloudAwsProvider.AccessKeyId)
	
    
	A(cloudAwsProviderMap, "accountId",cloudAwsProvider.AccountId)
	
    
	A(cloudAwsProviderMap, "annotation",cloudAwsProvider.Annotation)
	
    
	A(cloudAwsProviderMap, "email",cloudAwsProvider.Email)
	
    
	A(cloudAwsProviderMap, "httpProxy",cloudAwsProvider.HttpProxy)
	
    
	A(cloudAwsProviderMap, "isAccountInOrg",cloudAwsProvider.IsAccountInOrg)
	
    
	A(cloudAwsProviderMap, "isTrusted",cloudAwsProvider.IsTrusted)
	
    
	A(cloudAwsProviderMap, "nameAlias",cloudAwsProvider.NameAlias)
	
    
	A(cloudAwsProviderMap, "providerId",cloudAwsProvider.ProviderId)
	
    
	A(cloudAwsProviderMap, "region",cloudAwsProvider.Region)
	
    
	A(cloudAwsProviderMap, "secretAccessKey",cloudAwsProvider.SecretAccessKey)
	
    
	

	return cloudAwsProviderMap, err
}

func CloudAWSProviderFromContainerList(cont *container.Container, index int) *CloudAWSProvider {

	CloudAWSProviderCont := cont.S("imdata").Index(index).S(CloudawsproviderClassName, "attributes")
	return &CloudAWSProvider{
		BaseAttributes{
			DistinguishedName: G(CloudAWSProviderCont, "dn"),
			Description:       G(CloudAWSProviderCont, "descr"),
			Status:            G(CloudAWSProviderCont, "status"),
			ClassName:         CloudawsproviderClassName,
			Rn:                G(CloudAWSProviderCont, "rn"),
		},
        
		CloudAWSProviderAttributes{
		
        
	        AccessKeyId : G(CloudAWSProviderCont, "accessKeyId"),
		
        
	        AccountId : G(CloudAWSProviderCont, "accountId"),
		
        
	        Annotation : G(CloudAWSProviderCont, "annotation"),
		
        
	        Email : G(CloudAWSProviderCont, "email"),
		
        
	        HttpProxy : G(CloudAWSProviderCont, "httpProxy"),
		
        
	        IsAccountInOrg : G(CloudAWSProviderCont, "isAccountInOrg"),
		
        
	        IsTrusted : G(CloudAWSProviderCont, "isTrusted"),
		
        
	        NameAlias : G(CloudAWSProviderCont, "nameAlias"),
		
        
	        ProviderId : G(CloudAWSProviderCont, "providerId"),
		
        
	        Region : G(CloudAWSProviderCont, "region"),
		
        
	        SecretAccessKey : G(CloudAWSProviderCont, "secretAccessKey"),
		
        		
        },
        
	}
}

func CloudAWSProviderFromContainer(cont *container.Container) *CloudAWSProvider {

	return CloudAWSProviderFromContainerList(cont, 0)
}

func CloudAWSProviderListFromContainer(cont *container.Container) []*CloudAWSProvider {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudAWSProvider, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudAWSProviderFromContainerList(cont, i)
	}

	return arr
}