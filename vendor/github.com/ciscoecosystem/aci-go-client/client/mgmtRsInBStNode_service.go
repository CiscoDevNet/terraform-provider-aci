package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateInbandStaticNode(tDn string, in_band_management_epg string, management_profile string, tenant string, description string, mgmtRsInBStNodeattr models.InbandStaticNodeAttributes) (*models.InbandStaticNode, error) {
	rn := fmt.Sprintf("rsinBStNode-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s", tenant, management_profile, in_band_management_epg)
	mgmtRsInBStNode := models.NewInbandStaticNode(rn, parentDn, description, mgmtRsInBStNodeattr)
	err := sm.Save(mgmtRsInBStNode)
	return mgmtRsInBStNode, err
}

func (sm *ServiceManager) ReadInbandStaticNode(tDn string, in_band_management_epg string, management_profile string, tenant string) (*models.InbandStaticNode, error) {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s/rsinBStNode-[%s]", tenant, management_profile, in_band_management_epg, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtRsInBStNode := models.InbandStaticNodeFromContainer(cont)
	return mgmtRsInBStNode, nil
}

func (sm *ServiceManager) DeleteInbandStaticNode(tDn string, in_band_management_epg string, management_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s/rsinBStNode-[%s]", tenant, management_profile, in_band_management_epg, tDn)
	return sm.DeleteByDn(dn, models.MgmtrsinbstnodeClassName)
}

func (sm *ServiceManager) UpdateInbandStaticNode(tDn string, in_band_management_epg string, management_profile string, tenant string, description string, mgmtRsInBStNodeattr models.InbandStaticNodeAttributes) (*models.InbandStaticNode, error) {
	rn := fmt.Sprintf("rsinBStNode-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s", tenant, management_profile, in_band_management_epg)
	mgmtRsInBStNode := models.NewInbandStaticNode(rn, parentDn, description, mgmtRsInBStNodeattr)

	mgmtRsInBStNode.Status = "modified"
	err := sm.Save(mgmtRsInBStNode)
	return mgmtRsInBStNode, err

}

func (sm *ServiceManager) ListInbandStaticNode(in_band_management_epg string, management_profile string, tenant string) ([]*models.InbandStaticNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/mgmtp-%s/inb-%s/mgmtRsInBStNode.json", baseurlStr, tenant, management_profile, in_band_management_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.InbandStaticNodeListFromContainer(cont)

	return list, err
}
