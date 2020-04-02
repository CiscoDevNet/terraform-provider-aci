package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FirmwarefwpClassName = "firmwareFwP"

type FirmwarePolicy struct {
	BaseAttributes
	FirmwarePolicyAttributes
}

type FirmwarePolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	EffectiveOnReboot string `json:",omitempty"`

	IgnoreCompat string `json:",omitempty"`

	InternalLabel string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Version string `json:",omitempty"`

	VersionCheckOverride string `json:",omitempty"`
}

func NewFirmwarePolicy(firmwareFwPRn, parentDn, description string, firmwareFwPattr FirmwarePolicyAttributes) *FirmwarePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareFwPRn)
	return &FirmwarePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FirmwarefwpClassName,
			Rn:                firmwareFwPRn,
		},

		FirmwarePolicyAttributes: firmwareFwPattr,
	}
}

func (firmwareFwP *FirmwarePolicy) ToMap() (map[string]string, error) {
	firmwareFwPMap, err := firmwareFwP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareFwPMap, "name", firmwareFwP.Name)

	A(firmwareFwPMap, "annotation", firmwareFwP.Annotation)

	A(firmwareFwPMap, "effectiveOnReboot", firmwareFwP.EffectiveOnReboot)

	A(firmwareFwPMap, "ignoreCompat", firmwareFwP.IgnoreCompat)

	A(firmwareFwPMap, "internalLabel", firmwareFwP.InternalLabel)

	A(firmwareFwPMap, "nameAlias", firmwareFwP.NameAlias)

	A(firmwareFwPMap, "version", firmwareFwP.Version)

	A(firmwareFwPMap, "versionCheckOverride", firmwareFwP.VersionCheckOverride)

	return firmwareFwPMap, err
}

func FirmwarePolicyFromContainerList(cont *container.Container, index int) *FirmwarePolicy {

	FirmwarePolicyCont := cont.S("imdata").Index(index).S(FirmwarefwpClassName, "attributes")
	return &FirmwarePolicy{
		BaseAttributes{
			DistinguishedName: G(FirmwarePolicyCont, "dn"),
			Description:       G(FirmwarePolicyCont, "descr"),
			Status:            G(FirmwarePolicyCont, "status"),
			ClassName:         FirmwarefwpClassName,
			Rn:                G(FirmwarePolicyCont, "rn"),
		},

		FirmwarePolicyAttributes{

			Name: G(FirmwarePolicyCont, "name"),

			Annotation: G(FirmwarePolicyCont, "annotation"),

			EffectiveOnReboot: G(FirmwarePolicyCont, "effectiveOnReboot"),

			IgnoreCompat: G(FirmwarePolicyCont, "ignoreCompat"),

			InternalLabel: G(FirmwarePolicyCont, "internalLabel"),

			NameAlias: G(FirmwarePolicyCont, "nameAlias"),

			Version: G(FirmwarePolicyCont, "version"),

			VersionCheckOverride: G(FirmwarePolicyCont, "versionCheckOverride"),
		},
	}
}

func FirmwarePolicyFromContainer(cont *container.Container) *FirmwarePolicy {

	return FirmwarePolicyFromContainerList(cont, 0)
}

func FirmwarePolicyListFromContainer(cont *container.Container) []*FirmwarePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FirmwarePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = FirmwarePolicyFromContainerList(cont, i)
	}

	return arr
}
