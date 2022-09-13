package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DncloudCredentials        = "uni/tn-%s/credentials-%s"
	RncloudCredentials        = "credentials-%s"
	ParentDncloudCredentials  = "uni/tn-%s"
	CloudcredentialsClassName = "cloudCredentials"
)

type CloudCredentials struct {
	BaseAttributes
	NameAliasAttribute
	CloudCredentialsAttributes
}

type CloudCredentialsAttributes struct {
	Annotation    string `json:",omitempty"`
	ClientId      string `json:",omitempty"`
	Email         string `json:",omitempty"`
	HttpProxy     string `json:",omitempty"`
	Key           string `json:",omitempty"`
	KeyId         string `json:",omitempty"`
	Name          string `json:",omitempty"`
	RsaPrivateKey string `json:",omitempty"`
}

func NewCloudCredentials(cloudCredentialsRn, parentDn, nameAlias string, cloudCredentialsAttr CloudCredentialsAttributes) *CloudCredentials {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudCredentialsRn)
	return &CloudCredentials{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudcredentialsClassName,
			Rn:                cloudCredentialsRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		CloudCredentialsAttributes: cloudCredentialsAttr,
	}
}

func (cloudCredentials *CloudCredentials) ToMap() (map[string]string, error) {
	cloudCredentialsMap, err := cloudCredentials.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := cloudCredentials.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(cloudCredentialsMap, key, value)
	}

	A(cloudCredentialsMap, "clientId", cloudCredentials.ClientId)
	A(cloudCredentialsMap, "email", cloudCredentials.Email)
	A(cloudCredentialsMap, "httpProxy", cloudCredentials.HttpProxy)
	A(cloudCredentialsMap, "key", cloudCredentials.Key)
	A(cloudCredentialsMap, "keyId", cloudCredentials.KeyId)
	A(cloudCredentialsMap, "name", cloudCredentials.Name)
	A(cloudCredentialsMap, "rsaPrivateKey", cloudCredentials.RsaPrivateKey)
	return cloudCredentialsMap, err
}

func CloudCredentialsFromContainerList(cont *container.Container, index int) *CloudCredentials {
	CloudCredentialsCont := cont.S("imdata").Index(index).S(CloudcredentialsClassName, "attributes")
	return &CloudCredentials{
		BaseAttributes{
			DistinguishedName: G(CloudCredentialsCont, "dn"),
			Status:            G(CloudCredentialsCont, "status"),
			ClassName:         CloudcredentialsClassName,
			Rn:                G(CloudCredentialsCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(CloudCredentialsCont, "nameAlias"),
		},
		CloudCredentialsAttributes{
			ClientId:      G(CloudCredentialsCont, "clientId"),
			Email:         G(CloudCredentialsCont, "email"),
			HttpProxy:     G(CloudCredentialsCont, "httpProxy"),
			Key:           G(CloudCredentialsCont, "key"),
			KeyId:         G(CloudCredentialsCont, "keyId"),
			Name:          G(CloudCredentialsCont, "name"),
			RsaPrivateKey: G(CloudCredentialsCont, "rsaPrivateKey"),
		},
	}
}

func CloudCredentialsFromContainer(cont *container.Container) *CloudCredentials {
	return CloudCredentialsFromContainerList(cont, 0)
}

func CloudCredentialsListFromContainer(cont *container.Container) []*CloudCredentials {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudCredentials, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudCredentialsFromContainerList(cont, i)
	}

	return arr
}
