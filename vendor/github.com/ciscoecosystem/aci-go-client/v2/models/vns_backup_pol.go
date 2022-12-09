package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnvnsBackupPol        = "uni/tn-%s/svcCont/backupPol-%s"
	RnvnsBackupPol        = "backupPol-%s"
	ParentDnvnsBackupPol  = "uni/tn-%s/svcCont"
	VnsbackuppolClassName = "vnsBackupPol" // PBR Backup Policy ClassName
	RnvnsSvcCont          = "svcCont"      // Service Container RN
	VnsSvcContClassName   = "vnsSvcCont"   // Service Container ClassName
)

type PBRBackupPolicy struct {
	BaseAttributes
	NameAliasAttribute
	PBRBackupPolicyAttributes
}

type PBRBackupPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewPBRBackupPolicy(vnsBackupPolRn, parentDn, description, nameAlias string, vnsBackupPolAttr PBRBackupPolicyAttributes) *PBRBackupPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsBackupPolRn)
	return &PBRBackupPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsbackuppolClassName,
			Rn:                vnsBackupPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		PBRBackupPolicyAttributes: vnsBackupPolAttr,
	}
}

func (vnsBackupPol *PBRBackupPolicy) ToMap() (map[string]string, error) {
	vnsBackupPolMap, err := vnsBackupPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsBackupPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsBackupPolMap, key, value)
	}

	A(vnsBackupPolMap, "annotation", vnsBackupPol.Annotation)
	A(vnsBackupPolMap, "name", vnsBackupPol.Name)
	return vnsBackupPolMap, err
}

func PBRBackupPolicyFromContainerList(cont *container.Container, index int) *PBRBackupPolicy {
	PBRBackupPolicyCont := cont.S("imdata").Index(index).S(VnsbackuppolClassName, "attributes")
	return &PBRBackupPolicy{
		BaseAttributes{
			DistinguishedName: G(PBRBackupPolicyCont, "dn"),
			Description:       G(PBRBackupPolicyCont, "descr"),
			Status:            G(PBRBackupPolicyCont, "status"),
			ClassName:         VnsbackuppolClassName,
			Rn:                G(PBRBackupPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(PBRBackupPolicyCont, "nameAlias"),
		},
		PBRBackupPolicyAttributes{
			Annotation: G(PBRBackupPolicyCont, "annotation"),
			Name:       G(PBRBackupPolicyCont, "name"),
		},
	}
}

func PBRBackupPolicyFromContainer(cont *container.Container) *PBRBackupPolicy {
	return PBRBackupPolicyFromContainerList(cont, 0)
}

func PBRBackupPolicyListFromContainer(cont *container.Container) []*PBRBackupPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PBRBackupPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = PBRBackupPolicyFromContainerList(cont, i)
	}

	return arr
}
