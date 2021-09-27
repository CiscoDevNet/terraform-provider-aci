package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateIPAgingPolicy(name string, description string, nameAlias string, epIpAgingPAttr models.IPAgingPolicyAttributes) (*models.IPAgingPolicy, error) {
	rn := fmt.Sprintf(models.RnepIpAgingP, name)
	parentDn := fmt.Sprintf(models.ParentDnepIpAgingP)
	epIpAgingP := models.NewIPAgingPolicy(rn, parentDn, description, nameAlias, epIpAgingPAttr)
	err := sm.Save(epIpAgingP)
	return epIpAgingP, err
}

func (sm *ServiceManager) ReadIPAgingPolicy(name string) (*models.IPAgingPolicy, error) {
	dn := fmt.Sprintf(models.DnepIpAgingP, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	epIpAgingP := models.IPAgingPolicyFromContainer(cont)
	return epIpAgingP, nil
}

func (sm *ServiceManager) DeleteIPAgingPolicy(name string) error {
	dn := fmt.Sprintf(models.DnepIpAgingP, name)
	return sm.DeleteByDn(dn, models.EpipagingpClassName)
}

func (sm *ServiceManager) UpdateIPAgingPolicy(name string, description string, nameAlias string, epIpAgingPAttr models.IPAgingPolicyAttributes) (*models.IPAgingPolicy, error) {
	rn := fmt.Sprintf(models.RnepIpAgingP, name)
	parentDn := fmt.Sprintf(models.ParentDnepIpAgingP)
	epIpAgingP := models.NewIPAgingPolicy(rn, parentDn, description, nameAlias, epIpAgingPAttr)
	epIpAgingP.Status = "modified"
	err := sm.Save(epIpAgingP)
	return epIpAgingP, err
}

func (sm *ServiceManager) ListIPAgingPolicy() ([]*models.IPAgingPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/epIpAgingP.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.IPAgingPolicyListFromContainer(cont)
	return list, err
}
