package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvRsCloudAccount        = "uni/tn-%s/rsCloudAccount"
	RnfvRsCloudAccount        = "rsCloudAccount"
	ParentDnfvRsCloudAccount  = "uni/tn-%s"
	FvrscloudaccountClassName = "fvRsCloudAccount"
)

type TenantToCloudAccountAssociation struct {
	BaseAttributes
	NameAliasAttribute
	TenantToCloudAccountAssociationAttributes
}

type TenantToCloudAccountAssociationAttributes struct {
	Annotation string `json:",omitempty"`
	TDn        string `json:",omitempty"`
}

func NewTenantToCloudAccountAssociation(fvRsCloudAccountRn, parentDn, nameAlias string, fvRsCloudAccountAttr TenantToCloudAccountAssociationAttributes) *TenantToCloudAccountAssociation {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsCloudAccountRn)
	return &TenantToCloudAccountAssociation{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvrscloudaccountClassName,
			Rn:                fvRsCloudAccountRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TenantToCloudAccountAssociationAttributes: fvRsCloudAccountAttr,
	}
}

func (fvRsCloudAccount *TenantToCloudAccountAssociation) ToMap() (map[string]string, error) {
	fvRsCloudAccountMap, err := fvRsCloudAccount.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := fvRsCloudAccount.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(fvRsCloudAccountMap, key, value)
	}

	A(fvRsCloudAccountMap, "annotation", fvRsCloudAccount.Annotation)
	A(fvRsCloudAccountMap, "tDn", fvRsCloudAccount.TDn)
	return fvRsCloudAccountMap, err
}

func TenantToCloudAccountAssociationFromContainerList(cont *container.Container, index int) *TenantToCloudAccountAssociation {
	TenantToCloudAccountAssociationCont := cont.S("imdata").Index(index).S(FvrscloudaccountClassName, "attributes")
	return &TenantToCloudAccountAssociation{
		BaseAttributes{
			DistinguishedName: G(TenantToCloudAccountAssociationCont, "dn"),
			Status:            G(TenantToCloudAccountAssociationCont, "status"),
			ClassName:         FvrscloudaccountClassName,
			Rn:                G(TenantToCloudAccountAssociationCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TenantToCloudAccountAssociationCont, "nameAlias"),
		},
		TenantToCloudAccountAssociationAttributes{
			Annotation: G(TenantToCloudAccountAssociationCont, "annotation"),
			TDn:        G(TenantToCloudAccountAssociationCont, "tDn"),
		},
	}
}

func TenantToCloudAccountAssociationFromContainer(cont *container.Container) *TenantToCloudAccountAssociation {
	return TenantToCloudAccountAssociationFromContainerList(cont, 0)
}

func TenantToCloudAccountAssociationListFromContainer(cont *container.Container) []*TenantToCloudAccountAssociation {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TenantToCloudAccountAssociation, length)

	for i := 0; i < length; i++ {
		arr[i] = TenantToCloudAccountAssociationFromContainerList(cont, i)
	}

	return arr
}
