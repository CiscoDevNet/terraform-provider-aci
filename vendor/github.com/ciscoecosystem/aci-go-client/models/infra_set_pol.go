package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DninfraSetPol        = "uni/infra/settings"
	RninfraSetPol        = "settings"
	ParentDninfraSetPol  = "uni/infra"
	InfrasetpolClassName = "infraSetPol"
)

type FabricWideSettingsPolicy struct {
	BaseAttributes
	NameAliasAttribute
	FabricWideSettingsPolicyAttributes
}

type FabricWideSettingsPolicyAttributes struct {
	Annotation                 string `json:",omitempty"`
	DisableEpDampening         string `json:",omitempty"`
	EnableMoStreaming          string `json:",omitempty"`
	EnableRemoteLeafDirect     string `json:",omitempty"`
	EnforceSubnetCheck         string `json:",omitempty"`
	Name                       string `json:",omitempty"`
	OpflexpAuthenticateClients string `json:",omitempty"`
	OpflexpUseSsl              string `json:",omitempty"`
	RestrictInfraVLANTraffic   string `json:",omitempty"`
	UnicastXrEpLearnDisable    string `json:",omitempty"`
	ValidateOverlappingVlans   string `json:",omitempty"`
}

func NewFabricWideSettingsPolicy(infraSetPolRn, parentDn, description, nameAlias string, infraSetPolAttr FabricWideSettingsPolicyAttributes) *FabricWideSettingsPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSetPolRn)
	return &FabricWideSettingsPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrasetpolClassName,
			Rn:                infraSetPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		FabricWideSettingsPolicyAttributes: infraSetPolAttr,
	}
}

func (infraSetPol *FabricWideSettingsPolicy) ToMap() (map[string]string, error) {
	infraSetPolMap, err := infraSetPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := infraSetPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(infraSetPolMap, key, value)
	}
	A(infraSetPolMap, "annotation", infraSetPol.Annotation)
	A(infraSetPolMap, "disableEpDampening", infraSetPol.DisableEpDampening)
	A(infraSetPolMap, "enableMoStreaming", infraSetPol.EnableMoStreaming)
	A(infraSetPolMap, "enableRemoteLeafDirect", infraSetPol.EnableRemoteLeafDirect)
	A(infraSetPolMap, "enforceSubnetCheck", infraSetPol.EnforceSubnetCheck)
	A(infraSetPolMap, "name", infraSetPol.Name)
	A(infraSetPolMap, "opflexpAuthenticateClients", infraSetPol.OpflexpAuthenticateClients)
	A(infraSetPolMap, "opflexpUseSsl", infraSetPol.OpflexpUseSsl)
	A(infraSetPolMap, "restrictInfraVLANTraffic", infraSetPol.RestrictInfraVLANTraffic)
	A(infraSetPolMap, "unicastXrEpLearnDisable", infraSetPol.UnicastXrEpLearnDisable)
	A(infraSetPolMap, "validateOverlappingVlans", infraSetPol.ValidateOverlappingVlans)
	return infraSetPolMap, err
}

func FabricWideSettingsPolicyFromContainerList(cont *container.Container, index int) *FabricWideSettingsPolicy {
	FabricWideSettingsPolicyCont := cont.S("imdata").Index(index).S(InfrasetpolClassName, "attributes")
	return &FabricWideSettingsPolicy{
		BaseAttributes{
			DistinguishedName: G(FabricWideSettingsPolicyCont, "dn"),
			Description:       G(FabricWideSettingsPolicyCont, "descr"),
			Status:            G(FabricWideSettingsPolicyCont, "status"),
			ClassName:         InfrasetpolClassName,
			Rn:                G(FabricWideSettingsPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(FabricWideSettingsPolicyCont, "nameAlias"),
		},
		FabricWideSettingsPolicyAttributes{
			Annotation:                 G(FabricWideSettingsPolicyCont, "annotation"),
			DisableEpDampening:         G(FabricWideSettingsPolicyCont, "disableEpDampening"),
			EnableMoStreaming:          G(FabricWideSettingsPolicyCont, "enableMoStreaming"),
			EnableRemoteLeafDirect:     G(FabricWideSettingsPolicyCont, "enableRemoteLeafDirect"),
			EnforceSubnetCheck:         G(FabricWideSettingsPolicyCont, "enforceSubnetCheck"),
			Name:                       G(FabricWideSettingsPolicyCont, "name"),
			OpflexpAuthenticateClients: G(FabricWideSettingsPolicyCont, "opflexpAuthenticateClients"),
			OpflexpUseSsl:              G(FabricWideSettingsPolicyCont, "opflexpUseSsl"),
			RestrictInfraVLANTraffic:   G(FabricWideSettingsPolicyCont, "restrictInfraVLANTraffic"),
			UnicastXrEpLearnDisable:    G(FabricWideSettingsPolicyCont, "unicastXrEpLearnDisable"),
			ValidateOverlappingVlans:   G(FabricWideSettingsPolicyCont, "validateOverlappingVlans"),
		},
	}
}

func FabricWideSettingsPolicyFromContainer(cont *container.Container) *FabricWideSettingsPolicy {
	return FabricWideSettingsPolicyFromContainerList(cont, 0)
}

func FabricWideSettingsPolicyListFromContainer(cont *container.Container) []*FabricWideSettingsPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*FabricWideSettingsPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = FabricWideSettingsPolicyFromContainerList(cont, i)
	}
	return arr
}
