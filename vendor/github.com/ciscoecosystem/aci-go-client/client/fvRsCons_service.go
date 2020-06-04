package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateContractConsumer(tnVzBrCPName string, application_epg string, application_profile string, tenant string, fvRsConsattr models.ContractConsumerAttributes) (*models.ContractConsumer, error) {
	rn := fmt.Sprintf("rscons-%s", tnVzBrCPName)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsCons := models.NewContractConsumer(rn, parentDn, fvRsConsattr)
	err := sm.Save(fvRsCons)
	return fvRsCons, err
}

func (sm *ServiceManager) ReadContractConsumer(tnVzBrCPName string, application_epg string, application_profile string, tenant string) (*models.ContractConsumer, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rscons-%s", tenant, application_profile, application_epg, tnVzBrCPName)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsCons := models.ContractConsumerFromContainer(cont)
	return fvRsCons, nil
}

func (sm *ServiceManager) DeleteContractConsumer(tnVzBrCPName string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rscons-%s", tenant, application_profile, application_epg, tnVzBrCPName)
	return sm.DeleteByDn(dn, models.FvrsconsClassName)
}

func (sm *ServiceManager) UpdateContractConsumer(tnVzBrCPName string, application_epg string, application_profile string, tenant string, fvRsConsattr models.ContractConsumerAttributes) (*models.ContractConsumer, error) {
	rn := fmt.Sprintf("rscons-%s", tnVzBrCPName)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsCons := models.NewContractConsumer(rn, parentDn, fvRsConsattr)

	fvRsCons.Status = "modified"
	err := sm.Save(fvRsCons)
	return fvRsCons, err

}

func (sm *ServiceManager) ListContractConsumer(application_epg string, application_profile string, tenant string) ([]*models.ContractConsumer, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/fvRsCons.json", baseurlStr, tenant, application_profile, application_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ContractConsumerListFromContainer(cont)

	return list, err
}
