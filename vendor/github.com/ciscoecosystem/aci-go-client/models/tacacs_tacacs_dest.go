package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DntacacsTacacsDest        = "uni/fabric/tacacsgroup-%s/tacacsdest-%s-port-%s"
	RntacacsTacacsDest        = "tacacsdest-%s-port-%s"
	ParentDntacacsTacacsDest  = "uni/fabric/tacacsgroup-%s"
	TacacstacacsdestClassName = "tacacsTacacsDest"
)

type TACACSDestination struct {
	BaseAttributes
	NameAliasAttribute
	TACACSDestinationAttributes
}

type TACACSDestinationAttributes struct {
	Annotation   string `json:",omitempty"`
	AuthProtocol string `json:",omitempty"`
	Host         string `json:",omitempty"`
	Key          string `json:",omitempty"`
	Name         string `json:",omitempty"`
	Port         string `json:",omitempty"`
}

func NewTACACSDestination(tacacsTacacsDestRn, parentDn, description, nameAlias string, tacacsTacacsDestAttr TACACSDestinationAttributes) *TACACSDestination {
	dn := fmt.Sprintf("%s/%s", parentDn, tacacsTacacsDestRn)
	return &TACACSDestination{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         TacacstacacsdestClassName,
			Rn:                tacacsTacacsDestRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TACACSDestinationAttributes: tacacsTacacsDestAttr,
	}
}

func (tacacsTacacsDest *TACACSDestination) ToMap() (map[string]string, error) {
	tacacsTacacsDestMap, err := tacacsTacacsDest.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := tacacsTacacsDest.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(tacacsTacacsDestMap, key, value)
	}
	A(tacacsTacacsDestMap, "annotation", tacacsTacacsDest.Annotation)
	A(tacacsTacacsDestMap, "authProtocol", tacacsTacacsDest.AuthProtocol)
	A(tacacsTacacsDestMap, "host", tacacsTacacsDest.Host)
	A(tacacsTacacsDestMap, "key", tacacsTacacsDest.Key)
	A(tacacsTacacsDestMap, "name", tacacsTacacsDest.Name)
	A(tacacsTacacsDestMap, "port", tacacsTacacsDest.Port)
	return tacacsTacacsDestMap, err
}

func TACACSDestinationFromContainerList(cont *container.Container, index int) *TACACSDestination {
	TACACSDestinationCont := cont.S("imdata").Index(index).S(TacacstacacsdestClassName, "attributes")
	return &TACACSDestination{
		BaseAttributes{
			DistinguishedName: G(TACACSDestinationCont, "dn"),
			Description:       G(TACACSDestinationCont, "descr"),
			Status:            G(TACACSDestinationCont, "status"),
			ClassName:         TacacstacacsdestClassName,
			Rn:                G(TACACSDestinationCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TACACSDestinationCont, "nameAlias"),
		},
		TACACSDestinationAttributes{
			Annotation:   G(TACACSDestinationCont, "annotation"),
			AuthProtocol: G(TACACSDestinationCont, "authProtocol"),
			Host:         G(TACACSDestinationCont, "host"),
			Key:          G(TACACSDestinationCont, "key"),
			Name:         G(TACACSDestinationCont, "name"),
			Port:         G(TACACSDestinationCont, "port"),
		},
	}
}

func TACACSDestinationFromContainer(cont *container.Container) *TACACSDestination {
	return TACACSDestinationFromContainerList(cont, 0)
}

func TACACSDestinationListFromContainer(cont *container.Container) []*TACACSDestination {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TACACSDestination, length)
	for i := 0; i < length; i++ {
		arr[i] = TACACSDestinationFromContainerList(cont, i)
	}
	return arr
}
