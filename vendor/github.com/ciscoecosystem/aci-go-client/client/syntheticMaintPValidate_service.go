package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

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
