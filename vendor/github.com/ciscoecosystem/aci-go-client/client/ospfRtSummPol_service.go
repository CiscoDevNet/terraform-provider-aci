package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOspfRouteSummarization(name string, tenant string, description string, ospfRtSummPolattr models.OspfRouteSummarizationAttributes) (*models.OspfRouteSummarization, error) {
	rn := fmt.Sprintf("ospfrtsumm-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	ospfRtSummPol := models.NewOspfRouteSummarization(rn, parentDn, description, ospfRtSummPolattr)
	err := sm.Save(ospfRtSummPol)
	return ospfRtSummPol, err
}

func (sm *ServiceManager) ReadOspfRouteSummarization(name string, tenant string) (*models.OspfRouteSummarization, error) {
	dn := fmt.Sprintf("uni/tn-%s/ospfrtsumm-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfRtSummPol := models.OspfRouteSummarizationFromContainer(cont)
	return ospfRtSummPol, nil
}

func (sm *ServiceManager) DeleteOspfRouteSummarization(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ospfrtsumm-%s", tenant, name)
	return sm.DeleteByDn(dn, models.OspfrtsummpolClassName)
}

func (sm *ServiceManager) UpdateOspfRouteSummarization(name string, tenant string, description string, ospfRtSummPolattr models.OspfRouteSummarizationAttributes) (*models.OspfRouteSummarization, error) {
	rn := fmt.Sprintf("ospfrtsumm-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	ospfRtSummPol := models.NewOspfRouteSummarization(rn, parentDn, description, ospfRtSummPolattr)

	ospfRtSummPol.Status = "modified"
	err := sm.Save(ospfRtSummPol)
	return ospfRtSummPol, err

}

func (sm *ServiceManager) ListOspfRouteSummarization(tenant string) ([]*models.OspfRouteSummarization, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ospfRtSummPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OspfRouteSummarizationListFromContainer(cont)

	return list, err
}
