package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const ConfigimportpClassName = "configImportP"

type ConfigurationImportPolicy struct {
	BaseAttributes
	ConfigurationImportPolicyAttributes
}

type ConfigurationImportPolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	FailOnDecryptErrors string `json:",omitempty"`

	FileName string `json:",omitempty"`

	ImportMode string `json:",omitempty"`

	ImportType string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Snapshot string `json:",omitempty"`
}

func NewConfigurationImportPolicy(configImportPRn, parentDn, description string, configImportPattr ConfigurationImportPolicyAttributes) *ConfigurationImportPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, configImportPRn)
	return &ConfigurationImportPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         ConfigimportpClassName,
			Rn:                configImportPRn,
		},

		ConfigurationImportPolicyAttributes: configImportPattr,
	}
}

func (configImportP *ConfigurationImportPolicy) ToMap() (map[string]string, error) {
	configImportPMap, err := configImportP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(configImportPMap, "name", configImportP.Name)

	A(configImportPMap, "adminSt", configImportP.AdminSt)

	A(configImportPMap, "annotation", configImportP.Annotation)

	A(configImportPMap, "failOnDecryptErrors", configImportP.FailOnDecryptErrors)

	A(configImportPMap, "fileName", configImportP.FileName)

	A(configImportPMap, "importMode", configImportP.ImportMode)

	A(configImportPMap, "importType", configImportP.ImportType)

	A(configImportPMap, "nameAlias", configImportP.NameAlias)

	A(configImportPMap, "snapshot", configImportP.Snapshot)

	return configImportPMap, err
}

func ConfigurationImportPolicyFromContainerList(cont *container.Container, index int) *ConfigurationImportPolicy {

	ConfigurationImportPolicyCont := cont.S("imdata").Index(index).S(ConfigimportpClassName, "attributes")
	return &ConfigurationImportPolicy{
		BaseAttributes{
			DistinguishedName: G(ConfigurationImportPolicyCont, "dn"),
			Description:       G(ConfigurationImportPolicyCont, "descr"),
			Status:            G(ConfigurationImportPolicyCont, "status"),
			ClassName:         ConfigimportpClassName,
			Rn:                G(ConfigurationImportPolicyCont, "rn"),
		},

		ConfigurationImportPolicyAttributes{

			Name: G(ConfigurationImportPolicyCont, "name"),

			AdminSt: G(ConfigurationImportPolicyCont, "adminSt"),

			Annotation: G(ConfigurationImportPolicyCont, "annotation"),

			FailOnDecryptErrors: G(ConfigurationImportPolicyCont, "failOnDecryptErrors"),

			FileName: G(ConfigurationImportPolicyCont, "fileName"),

			ImportMode: G(ConfigurationImportPolicyCont, "importMode"),

			ImportType: G(ConfigurationImportPolicyCont, "importType"),

			NameAlias: G(ConfigurationImportPolicyCont, "nameAlias"),

			Snapshot: G(ConfigurationImportPolicyCont, "snapshot"),
		},
	}
}

func ConfigurationImportPolicyFromContainer(cont *container.Container) *ConfigurationImportPolicy {

	return ConfigurationImportPolicyFromContainerList(cont, 0)
}

func ConfigurationImportPolicyListFromContainer(cont *container.Container) []*ConfigurationImportPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ConfigurationImportPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = ConfigurationImportPolicyFromContainerList(cont, i)
	}

	return arr
}
