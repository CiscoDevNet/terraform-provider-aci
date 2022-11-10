package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudAccount        = "uni/tn-%s/act-[%s]-vendor-%s"
	RncloudAccount        = "act-[%s]-vendor-%s"
	ParentDncloudAccount  = "uni/tn-%s"
	CloudaccountClassName = "cloudAccount"
)

type CloudAccount struct {
	BaseAttributes
	NameAliasAttribute
	CloudAccountAttributes
}

type CloudAccountAttributes struct {
	Annotation string `json:",omitempty"`
	AccessType string `json:",omitempty"`
	Account_id string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Vendor     string `json:",omitempty"`
}

func NewCloudAccount(cloudAccountRn, parentDn, nameAlias string, cloudAccountAttr CloudAccountAttributes) *CloudAccount {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudAccountRn)
	return &CloudAccount{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudaccountClassName,
			Rn:                cloudAccountRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		CloudAccountAttributes: cloudAccountAttr,
	}
}

func (cloudAccount *CloudAccount) ToMap() (map[string]string, error) {
	cloudAccountMap, err := cloudAccount.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := cloudAccount.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(cloudAccountMap, key, value)
	}

	A(cloudAccountMap, "accessType", cloudAccount.AccessType)
	A(cloudAccountMap, "id", cloudAccount.Account_id)
	A(cloudAccountMap, "name", cloudAccount.Name)
	A(cloudAccountMap, "vendor", cloudAccount.Vendor)
	return cloudAccountMap, err
}

func CloudAccountFromContainerList(cont *container.Container, index int) *CloudAccount {
	AccountCont := cont.S("imdata").Index(index).S(CloudaccountClassName, "attributes")
	return &CloudAccount{
		BaseAttributes{
			DistinguishedName: G(AccountCont, "dn"),
			Status:            G(AccountCont, "status"),
			ClassName:         CloudaccountClassName,
			Rn:                G(AccountCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AccountCont, "nameAlias"),
		},
		CloudAccountAttributes{
			AccessType: G(AccountCont, "accessType"),
			Account_id: G(AccountCont, "id"),
			Name:       G(AccountCont, "name"),
			Vendor:     G(AccountCont, "vendor"),
		},
	}
}

func CloudAccountFromContainer(cont *container.Container) *CloudAccount {
	return CloudAccountFromContainerList(cont, 0)
}

func CloudAccountListFromContainer(cont *container.Container) []*CloudAccount {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudAccount, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudAccountFromContainerList(cont, i)
	}

	return arr
}
