package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateBgpRouteSummarization(name string, tenant string, description string, bgpRtSummPolattr models.BgpRouteSummarizationAttributes) (*models.BgpRouteSummarization, error) {

	rn := fmt.Sprintf(models.RnBgpRtSummPol, name)

	parentDn := fmt.Sprintf(models.ParentDnBgpRtSummPol, tenant)
	bgpRtSummPol := models.NewBgpRouteSummarization(rn, parentDn, description, bgpRtSummPolattr)

	err := sm.Save(bgpRtSummPol)
	return bgpRtSummPol, err
}

func (sm *ServiceManager) ReadBgpRouteSummarization(name string, tenant string) (*models.BgpRouteSummarization, error) {

	rn := fmt.Sprintf(models.RnBgpRtSummPol, name)

	parentDn := fmt.Sprintf(models.ParentDnBgpRtSummPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	bgpRtSummPol := models.BgpRouteSummarizationFromContainer(cont)
	return bgpRtSummPol, nil
}

func (sm *ServiceManager) DeleteBgpRouteSummarization(name string, tenant string) error {

	rn := fmt.Sprintf(models.RnBgpRtSummPol, name)

	parentDn := fmt.Sprintf(models.ParentDnBgpRtSummPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.BgprtsummpolClassName)
}

func (sm *ServiceManager) UpdateBgpRouteSummarization(name string, tenant string, description string, bgpRtSummPolattr models.BgpRouteSummarizationAttributes) (*models.BgpRouteSummarization, error) {

	rn := fmt.Sprintf(models.RnBgpRtSummPol, name)

	parentDn := fmt.Sprintf(models.ParentDnBgpRtSummPol, tenant)
	bgpRtSummPol := models.NewBgpRouteSummarization(rn, parentDn, description, bgpRtSummPolattr)

	bgpRtSummPol.Status = "modified"
	err := sm.Save(bgpRtSummPol)
	return bgpRtSummPol, err
}

func (sm *ServiceManager) ListBgpRouteSummarization(tenant string) ([]*models.BgpRouteSummarization, error) {

	parentDn := fmt.Sprintf(models.ParentDnBgpRtSummPol, tenant)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.BgprtsummpolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BgpRouteSummarizationListFromContainer(cont)
	return list, err
}
