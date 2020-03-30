package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const ConfigexportpClassName = "configExportP"

type ConfigurationExportPolicy struct {
	BaseAttributes
	ConfigurationExportPolicyAttributes
}

type ConfigurationExportPolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Format string `json:",omitempty"`

	IncludeSecureFields string `json:",omitempty"`

	MaxSnapshotCount string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Snapshot string `json:",omitempty"`

	TargetDn string `json:",omitempty"`
}

func NewConfigurationExportPolicy(configExportPRn, parentDn, description string, configExportPattr ConfigurationExportPolicyAttributes) *ConfigurationExportPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, configExportPRn)
	return &ConfigurationExportPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         ConfigexportpClassName,
			Rn:                configExportPRn,
		},

		ConfigurationExportPolicyAttributes: configExportPattr,
	}
}

func (configExportP *ConfigurationExportPolicy) ToMap() (map[string]string, error) {
	configExportPMap, err := configExportP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(configExportPMap, "name", configExportP.Name)

	A(configExportPMap, "adminSt", configExportP.AdminSt)

	A(configExportPMap, "annotation", configExportP.Annotation)

	A(configExportPMap, "format", configExportP.Format)

	A(configExportPMap, "includeSecureFields", configExportP.IncludeSecureFields)

	A(configExportPMap, "maxSnapshotCount", configExportP.MaxSnapshotCount)

	A(configExportPMap, "nameAlias", configExportP.NameAlias)

	A(configExportPMap, "snapshot", configExportP.Snapshot)

	A(configExportPMap, "targetDn", configExportP.TargetDn)

	return configExportPMap, err
}

func ConfigurationExportPolicyFromContainerList(cont *container.Container, index int) *ConfigurationExportPolicy {

	ConfigurationExportPolicyCont := cont.S("imdata").Index(index).S(ConfigexportpClassName, "attributes")
	return &ConfigurationExportPolicy{
		BaseAttributes{
			DistinguishedName: G(ConfigurationExportPolicyCont, "dn"),
			Description:       G(ConfigurationExportPolicyCont, "descr"),
			Status:            G(ConfigurationExportPolicyCont, "status"),
			ClassName:         ConfigexportpClassName,
			Rn:                G(ConfigurationExportPolicyCont, "rn"),
		},

		ConfigurationExportPolicyAttributes{

			Name: G(ConfigurationExportPolicyCont, "name"),

			AdminSt: G(ConfigurationExportPolicyCont, "adminSt"),

			Annotation: G(ConfigurationExportPolicyCont, "annotation"),

			Format: G(ConfigurationExportPolicyCont, "format"),

			IncludeSecureFields: G(ConfigurationExportPolicyCont, "includeSecureFields"),

			MaxSnapshotCount: G(ConfigurationExportPolicyCont, "maxSnapshotCount"),

			NameAlias: G(ConfigurationExportPolicyCont, "nameAlias"),

			Snapshot: G(ConfigurationExportPolicyCont, "snapshot"),

			TargetDn: G(ConfigurationExportPolicyCont, "targetDn"),
		},
	}
}

func ConfigurationExportPolicyFromContainer(cont *container.Container) *ConfigurationExportPolicy {

	return ConfigurationExportPolicyFromContainerList(cont, 0)
}

func ConfigurationExportPolicyListFromContainer(cont *container.Container) []*ConfigurationExportPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ConfigurationExportPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = ConfigurationExportPolicyFromContainerList(cont, i)
	}

	return arr
}
