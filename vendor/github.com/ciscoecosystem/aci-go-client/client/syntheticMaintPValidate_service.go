package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)
// Notes: PyQuery Based Request
// Requires a valid client APIC-Cookie to be present in the Header
// Ways to obtain the APIC-Cooki:
// (1) Use Username/Password for Client Authentication
// (2) Call client.Authenticate() prior to calling this function for AppUserName+Cert based clients

func (sm *ServiceManager) ListSyntheticMaintPValidateCtrlrMaintP(targetVersion string) ([]*models.MaintPValidate, error) {
	baseurlStr := "/mqapi2/deployment.query.json?mode=validateCtrlrMaintP"
	dnUrl := fmt.Sprintf("%s&targetVersion=%s", baseurlStr, targetVersion)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MaintPValidateListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) ListSyntheticMaintPValidateCtrlrMaintPReport(report string) ([]*models.MaintPValidate, error) {
	baseurlStr := "/mqapi2/deployment.query.json?mode=validateCtrlrMaintP"
	dnUrl := fmt.Sprintf("%s&report=", baseurlStr, report)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MaintPValidateListFromContainer(cont)
	return list, err
}
