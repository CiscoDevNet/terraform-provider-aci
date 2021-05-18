package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBgpBestPathPolicy(name string, tenant string, description string, bgpBestPathCtrlPolattr models.BgpBestPathPolicyAttributes) (*models.BgpBestPathPolicy, error) {
	rn := fmt.Sprintf("bestpath-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(rn, parentDn, description, bgpBestPathCtrlPolattr)
	err := sm.Save(bgpBestPathCtrlPol)
	return bgpBestPathCtrlPol, err
}

func (sm *ServiceManager) ReadBgpBestPathPolicy(name string, tenant string) (*models.BgpBestPathPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/bestpath-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpBestPathCtrlPol := models.BgpBestPathPolicyFromContainer(cont)
	return bgpBestPathCtrlPol, nil
}

func (sm *ServiceManager) DeleteBgpBestPathPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bestpath-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BgpbestpathctrlpolClassName)
}

func (sm *ServiceManager) UpdateBgpBestPathPolicy(name string, tenant string, description string, bgpBestPathCtrlPolattr models.BgpBestPathPolicyAttributes) (*models.BgpBestPathPolicy, error) {
	rn := fmt.Sprintf("bestpath-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(rn, parentDn, description, bgpBestPathCtrlPolattr)

	bgpBestPathCtrlPol.Status = "modified"
	err := sm.Save(bgpBestPathCtrlPol)
	return bgpBestPathCtrlPol, err

}

func (sm *ServiceManager) ListBgpBestPathPolicy(tenant string) ([]*models.BgpBestPathPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bgpBestPathCtrlPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BgpBestPathPolicyListFromContainer(cont)

	return list, err
}
