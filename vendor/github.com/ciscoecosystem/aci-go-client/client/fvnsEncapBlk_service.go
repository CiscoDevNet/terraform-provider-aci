package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRanges(to string, _from string, vlan_pool_allocMode string, vlan_pool string, description string, fvnsEncapBlkattr models.RangesAttributes) (*models.Ranges, error) {
	rn := fmt.Sprintf("from-[%s]-to-[%s]", _from, to)
	parentDn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s", vlan_pool, vlan_pool_allocMode)
	fvnsEncapBlk := models.NewRanges(rn, parentDn, description, fvnsEncapBlkattr)
	err := sm.Save(fvnsEncapBlk)
	return fvnsEncapBlk, err
}

func (sm *ServiceManager) ReadRanges(to string, _from string, vlan_pool_allocMode string, vlan_pool string) (*models.Ranges, error) {
	dn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s/from-[%s]-to-[%s]", vlan_pool, vlan_pool_allocMode, _from, to)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsEncapBlk := models.RangesFromContainer(cont)
	return fvnsEncapBlk, nil
}

func (sm *ServiceManager) DeleteRanges(to string, _from string, vlan_pool_allocMode string, vlan_pool string) error {
	dn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s/from-[%s]-to-[%s]", vlan_pool, vlan_pool_allocMode, _from, to)
	return sm.DeleteByDn(dn, models.FvnsencapblkClassName)
}

func (sm *ServiceManager) UpdateRanges(to string, _from string, vlan_pool_allocMode string, vlan_pool string, description string, fvnsEncapBlkattr models.RangesAttributes) (*models.Ranges, error) {
	rn := fmt.Sprintf("from-[%s]-to-[%s]", _from, to)
	parentDn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s", vlan_pool, vlan_pool_allocMode)
	fvnsEncapBlk := models.NewRanges(rn, parentDn, description, fvnsEncapBlkattr)

	fvnsEncapBlk.Status = "modified"
	err := sm.Save(fvnsEncapBlk)
	return fvnsEncapBlk, err

}

func (sm *ServiceManager) ListRanges(vlan_pool_allocMode string, vlan_pool string) ([]*models.Ranges, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/vlanns-[%s]-%s/fvnsEncapBlk.json", baseurlStr, vlan_pool, vlan_pool_allocMode)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.RangesListFromContainer(cont)

	return list, err
}
