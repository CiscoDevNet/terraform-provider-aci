package models

import (
	"encoding/json"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const naPayload = `
	"nameAlias": "%s",
`

type NameAliasAttribute struct {
	NameAlias string `json:"nameAlias"`
}

func (na *NameAliasAttribute) ToJson() (string, error) {
	data, err := json.Marshal(na)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (na *NameAliasAttribute) ToMap() (map[string]string, error) {
	jsonData, err := na.ToJson()
	if err != nil {
		return nil, err
	}
	cont, err := container.ParseJSON([]byte(jsonData))
	if err != nil {
		return nil, err
	}
	return toStringMap(cont.Data()), nil
}
