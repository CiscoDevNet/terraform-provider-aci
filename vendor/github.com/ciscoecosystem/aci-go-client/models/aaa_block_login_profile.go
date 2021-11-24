package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaBlockLoginProfile        = "uni/userext/blockloginp"
	RnaaaBlockLoginProfile        = "blockloginp"
	ParentDnaaaBlockLoginProfile  = "uni/userext"
	AaablockloginprofileClassName = "aaaBlockLoginProfile"
)

type BlockUserLoginsPolicy struct {
	BaseAttributes
	NameAliasAttribute
	BlockUserLoginsPolicyAttributes
}

type BlockUserLoginsPolicyAttributes struct {
	Annotation              string `json:",omitempty"`
	BlockDuration           string `json:",omitempty"`
	EnableLoginBlock        string `json:",omitempty"`
	MaxFailedAttempts       string `json:",omitempty"`
	MaxFailedAttemptsWindow string `json:",omitempty"`
	Name                    string `json:",omitempty"`
}

func NewBlockUserLoginsPolicy(aaaBlockLoginProfileRn, parentDn, description, nameAlias string, aaaBlockLoginProfileAttr BlockUserLoginsPolicyAttributes) *BlockUserLoginsPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaBlockLoginProfileRn)
	return &BlockUserLoginsPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaablockloginprofileClassName,
			Rn:                aaaBlockLoginProfileRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		BlockUserLoginsPolicyAttributes: aaaBlockLoginProfileAttr,
	}
}

func (aaaBlockLoginProfile *BlockUserLoginsPolicy) ToMap() (map[string]string, error) {
	aaaBlockLoginProfileMap, err := aaaBlockLoginProfile.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaBlockLoginProfile.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaBlockLoginProfileMap, key, value)
	}
	A(aaaBlockLoginProfileMap, "annotation", aaaBlockLoginProfile.Annotation)
	A(aaaBlockLoginProfileMap, "blockDuration", aaaBlockLoginProfile.BlockDuration)
	A(aaaBlockLoginProfileMap, "enableLoginBlock", aaaBlockLoginProfile.EnableLoginBlock)
	A(aaaBlockLoginProfileMap, "maxFailedAttempts", aaaBlockLoginProfile.MaxFailedAttempts)
	A(aaaBlockLoginProfileMap, "maxFailedAttemptsWindow", aaaBlockLoginProfile.MaxFailedAttemptsWindow)
	A(aaaBlockLoginProfileMap, "name", aaaBlockLoginProfile.Name)
	return aaaBlockLoginProfileMap, err
}

func BlockUserLoginsPolicyFromContainerList(cont *container.Container, index int) *BlockUserLoginsPolicy {
	BlockUserLoginsPolicyCont := cont.S("imdata").Index(index).S(AaablockloginprofileClassName, "attributes")
	return &BlockUserLoginsPolicy{
		BaseAttributes{
			DistinguishedName: G(BlockUserLoginsPolicyCont, "dn"),
			Description:       G(BlockUserLoginsPolicyCont, "descr"),
			Status:            G(BlockUserLoginsPolicyCont, "status"),
			ClassName:         AaablockloginprofileClassName,
			Rn:                G(BlockUserLoginsPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(BlockUserLoginsPolicyCont, "nameAlias"),
		},
		BlockUserLoginsPolicyAttributes{
			Annotation:              G(BlockUserLoginsPolicyCont, "annotation"),
			BlockDuration:           G(BlockUserLoginsPolicyCont, "blockDuration"),
			EnableLoginBlock:        G(BlockUserLoginsPolicyCont, "enableLoginBlock"),
			MaxFailedAttempts:       G(BlockUserLoginsPolicyCont, "maxFailedAttempts"),
			MaxFailedAttemptsWindow: G(BlockUserLoginsPolicyCont, "maxFailedAttemptsWindow"),
			Name:                    G(BlockUserLoginsPolicyCont, "name"),
		},
	}
}

func BlockUserLoginsPolicyFromContainer(cont *container.Container) *BlockUserLoginsPolicy {
	return BlockUserLoginsPolicyFromContainerList(cont, 0)
}

func BlockUserLoginsPolicyListFromContainer(cont *container.Container) []*BlockUserLoginsPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*BlockUserLoginsPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = BlockUserLoginsPolicyFromContainerList(cont, i)
	}
	return arr
}
