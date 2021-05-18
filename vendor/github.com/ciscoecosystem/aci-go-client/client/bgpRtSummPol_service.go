package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBgpRouteSummarization(name string, tenant string, description string, bgpRtSummPolattr models.BgpRouteSummarizationAttributes) (*models.BgpRouteSummarization, error) {
	rn := fmt.Sprintf("bgprtsum-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpRtSummPol := models.NewBgpRouteSummarization(rn, parentDn, description, bgpRtSummPolattr)
	err := sm.Save(bgpRtSummPol)
	return bgpRtSummPol, err
}

func (sm *ServiceManager) ReadBgpRouteSummarization(name string, tenant string) (*models.BgpRouteSummarization, error) {
	dn := fmt.Sprintf("uni/tn-%s/bgprtsum-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpRtSummPol := models.BgpRouteSummarizationFromContainer(cont)
	return bgpRtSummPol, nil
}

func (sm *ServiceManager) DeleteBgpRouteSummarization(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bgprtsum-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BgprtsummpolClassName)
}

func (sm *ServiceManager) UpdateBgpRouteSummarization(name string, tenant string, description string, bgpRtSummPolattr models.BgpRouteSummarizationAttributes) (*models.BgpRouteSummarization, error) {
	rn := fmt.Sprintf("bgprtsum-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpRtSummPol := models.NewBgpRouteSummarization(rn, parentDn, description, bgpRtSummPolattr)

	bgpRtSummPol.Status = "modified"
	err := sm.Save(bgpRtSummPol)
	return bgpRtSummPol, err

}

func (sm *ServiceManager) ListBgpRouteSummarization(tenant string) ([]*models.BgpRouteSummarization, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bgpRtSummPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BgpRouteSummarizationListFromContainer(cont)

	return list, err
}
