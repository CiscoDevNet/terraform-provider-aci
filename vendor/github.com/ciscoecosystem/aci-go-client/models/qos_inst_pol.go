package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnqosInstPol        = "uni/infra/qosinst-%s"
	RnqosInstPol        = "qosinst-%s"
	ParentDnqosInstPol  = "uni/infra"
	QosinstpolClassName = "qosInstPol"
)

type QOSInstancePolicy struct {
	BaseAttributes
	NameAliasAttribute
	QOSInstancePolicyAttributes
}

type QOSInstancePolicyAttributes struct {
	EtrapAgeTimer       string `json:",omitempty"`
	EtrapBwThresh       string `json:",omitempty"`
	EtrapByteCt         string `json:",omitempty"`
	EtrapSt             string `json:",omitempty"`
	FabricFlushInterval string `json:",omitempty"`
	FabricFlushSt       string `json:",omitempty"`
	Annotation          string `json:",omitempty"`
	Ctrl                string `json:",omitempty"`
	Name                string `json:",omitempty"`
	UburstSpineQueues   string `json:",omitempty"`
	UburstTorQueues     string `json:",omitempty"`
}

func NewQOSInstancePolicy(qosInstPolRn, parentDn, description, nameAlias string, qosInstPolAttr QOSInstancePolicyAttributes) *QOSInstancePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, qosInstPolRn)
	return &QOSInstancePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         QosinstpolClassName,
			Rn:                qosInstPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		QOSInstancePolicyAttributes: qosInstPolAttr,
	}
}

func (qosInstPol *QOSInstancePolicy) ToMap() (map[string]string, error) {
	qosInstPolMap, err := qosInstPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := qosInstPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(qosInstPolMap, key, value)
	}
	A(qosInstPolMap, "EtrapAgeTimer", qosInstPol.EtrapAgeTimer)
	A(qosInstPolMap, "EtrapBwThresh", qosInstPol.EtrapBwThresh)
	A(qosInstPolMap, "EtrapByteCt", qosInstPol.EtrapByteCt)
	A(qosInstPolMap, "EtrapSt", qosInstPol.EtrapSt)
	A(qosInstPolMap, "FabricFlushInterval", qosInstPol.FabricFlushInterval)
	A(qosInstPolMap, "FabricFlushSt", qosInstPol.FabricFlushSt)
	A(qosInstPolMap, "annotation", qosInstPol.Annotation)
	A(qosInstPolMap, "ctrl", qosInstPol.Ctrl)
	A(qosInstPolMap, "name", qosInstPol.Name)
	A(qosInstPolMap, "uburstSpineQueues", qosInstPol.UburstSpineQueues)
	A(qosInstPolMap, "uburstTorQueues", qosInstPol.UburstTorQueues)
	return qosInstPolMap, err
}

func QOSInstancePolicyFromContainerList(cont *container.Container, index int) *QOSInstancePolicy {
	QOSInstancePolicyCont := cont.S("imdata").Index(index).S(QosinstpolClassName, "attributes")
	return &QOSInstancePolicy{
		BaseAttributes{
			DistinguishedName: G(QOSInstancePolicyCont, "dn"),
			Description:       G(QOSInstancePolicyCont, "descr"),
			Status:            G(QOSInstancePolicyCont, "status"),
			ClassName:         QosinstpolClassName,
			Rn:                G(QOSInstancePolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(QOSInstancePolicyCont, "nameAlias"),
		},
		QOSInstancePolicyAttributes{
			EtrapAgeTimer:       G(QOSInstancePolicyCont, "EtrapAgeTimer"),
			EtrapBwThresh:       G(QOSInstancePolicyCont, "EtrapBwThresh"),
			EtrapByteCt:         G(QOSInstancePolicyCont, "EtrapByteCt"),
			EtrapSt:             G(QOSInstancePolicyCont, "EtrapSt"),
			FabricFlushInterval: G(QOSInstancePolicyCont, "FabricFlushInterval"),
			FabricFlushSt:       G(QOSInstancePolicyCont, "FabricFlushSt"),
			Annotation:          G(QOSInstancePolicyCont, "annotation"),
			Ctrl:                G(QOSInstancePolicyCont, "ctrl"),
			Name:                G(QOSInstancePolicyCont, "name"),
			UburstSpineQueues:   G(QOSInstancePolicyCont, "uburstSpineQueues"),
			UburstTorQueues:     G(QOSInstancePolicyCont, "uburstTorQueues"),
		},
	}
}

func QOSInstancePolicyFromContainer(cont *container.Container) *QOSInstancePolicy {
	return QOSInstancePolicyFromContainerList(cont, 0)
}

func QOSInstancePolicyListFromContainer(cont *container.Container) []*QOSInstancePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*QOSInstancePolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = QOSInstancePolicyFromContainerList(cont, i)
	}
	return arr
}
