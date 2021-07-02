package client

import (
	"github.com/ciscoecosystem/aci-go-client/models"
)

// This is a special operation which validations a target Switch (leaf/spine) config
// Action keyword is "Validate" instead of List/Get/Create because there is no creation done on the APIC
// In addition, the incoming MO is different then the outgoing MO

// Notes: PyQuery Based Request
// Requires a valid client APIC-Cookie to be present in the Header
// Ways to obtain the APIC-Cooki:
// (1) Use Username/Password for Client Authentication
// (2) Call client.Authenticate() prior to calling this function for AppUserName+Cert based clients
func (sm *ServiceManager) ValidateSyntheticSwitchMaintP(syntheticMaintPSwitchDetailsattr models.MaintPSwitchDetailsAttributes) ([]*models.SwitchMaintPValidate, error) {
	postUrlStr := "/mqapi2/deployment.query.json?mode=validateSwitchMaintP"
	rn := "switchDetails"
	parentDn := "synthuni"
	maintPSwitchDetails := models.NewMaintPSwitchDetails(rn, parentDn, "", syntheticMaintPSwitchDetailsattr)
	cont, err := sm.PostViaURL(postUrlStr, maintPSwitchDetails)
	list := models.SwitchMaintPValidateListFromContainer(cont)
	return list, err
}
