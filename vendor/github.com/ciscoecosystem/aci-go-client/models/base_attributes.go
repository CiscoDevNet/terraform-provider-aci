package models

import (
	"encoding/json"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const baPayload = `
	"dn": "%s",
	"name": "%s",
	"descr": "%s",
`

type BaseAttributes struct {
	DistinguishedName string `json:"dn"`
	Status            string `json:"status"`
	Description       string `json:"descr"`
	ClassName         string `json:"-"`
	Rn                string `json:"rn"`
}

func (ba *BaseAttributes) ToJson() (string, error) {
	data, err := json.Marshal(ba)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (ba *BaseAttributes) ToMap() (map[string]string, error) {

	jsonData, err := ba.ToJson()
	if err != nil {
		return nil, err
	}
	cont, err := container.ParseJSON([]byte(jsonData))

	if err != nil {
		return nil, err
	}
	//

	cont.Set(ba.ClassName, "classname")

	if err != nil {
		return nil, err
	}
	//
	return toStringMap(cont.Data()), nil
}
