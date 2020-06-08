package client

import (
	"github.com/ciscoecosystem/aci-go-client/models"
)

// This is a special operation which validations a target Switch (leaf/spine) config
// Action keyword is "Validate" instead of List/Get/Create because there is no creation done on the APIC
// In addition, the incoming MO is different then the outgoing MO
func (sm *ServiceManager) ValidateSyntheticSwitchMaintP(syntheticMaintPSwitchDetailsattr models.MaintPSwitchDetailsAttributes) ([]*models.SwitchMaintPValidate, error) {
	postUrlStr := "/mqapi2/deployment.query.json?mode=validateSwitchMaintP"
    rn := "switchDetails"
    parentDn := "synthuni"
    maintPSwitchDetails := models.NewMaintPSwitchDetails(rn, parentDn, "", syntheticMaintPSwitchDetailsattr)
    cont, err := sm.PostViaURL(postUrlStr, maintPSwitchDetails)
	list := models.SwitchMaintPValidateListFromContainer(cont)
	return list, err
}
