package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOutofbandStaticNode(tDn string, out_of_band_management_epg string, management_profile string, tenant string, description string, mgmtRsOoBStNodeattr models.OutofbandStaticNodeAttributes) (*models.OutofbandStaticNode, error) {
	rn := fmt.Sprintf("rsooBStNode-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s", tenant, management_profile, out_of_band_management_epg)
	mgmtRsOoBStNode := models.NewOutofbandStaticNode(rn, parentDn, description, mgmtRsOoBStNodeattr)
	err := sm.Save(mgmtRsOoBStNode)
	return mgmtRsOoBStNode, err
}

func (sm *ServiceManager) ReadOutofbandStaticNode(tDn string, out_of_band_management_epg string, management_profile string, tenant string) (*models.OutofbandStaticNode, error) {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s/rsooBStNode-[%s]", tenant, management_profile, out_of_band_management_epg, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtRsOoBStNode := models.OutofbandStaticNodeFromContainer(cont)
	return mgmtRsOoBStNode, nil
}

func (sm *ServiceManager) DeleteOutofbandStaticNode(tDn string, out_of_band_management_epg string, management_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s/rsooBStNode-[%s]", tenant, management_profile, out_of_band_management_epg, tDn)
	return sm.DeleteByDn(dn, models.MgmtrsoobstnodeClassName)
}

func (sm *ServiceManager) UpdateOutofbandStaticNode(tDn string, out_of_band_management_epg string, management_profile string, tenant string, description string, mgmtRsOoBStNodeattr models.OutofbandStaticNodeAttributes) (*models.OutofbandStaticNode, error) {
	rn := fmt.Sprintf("rsooBStNode-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s", tenant, management_profile, out_of_band_management_epg)
	mgmtRsOoBStNode := models.NewOutofbandStaticNode(rn, parentDn, description, mgmtRsOoBStNodeattr)

	mgmtRsOoBStNode.Status = "modified"
	err := sm.Save(mgmtRsOoBStNode)
	return mgmtRsOoBStNode, err

}

func (sm *ServiceManager) ListOutofbandStaticNode(out_of_band_management_epg string, management_profile string, tenant string) ([]*models.OutofbandStaticNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/mgmtp-%s/oob-%s/mgmtRsOoBStNode.json", baseurlStr, tenant, management_profile, out_of_band_management_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OutofbandStaticNodeListFromContainer(cont)

	return list, err
}
