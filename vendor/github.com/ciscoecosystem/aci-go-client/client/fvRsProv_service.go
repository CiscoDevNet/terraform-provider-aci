package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateContractProvider(tnVzBrCPName string, application_epg string, application_profile string, tenant string, fvRsProvattr models.ContractProviderAttributes) (*models.ContractProvider, error) {
	rn := fmt.Sprintf("rsprov-%s", tnVzBrCPName)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsProv := models.NewContractProvider(rn, parentDn, fvRsProvattr)
	err := sm.Save(fvRsProv)
	return fvRsProv, err
}

func (sm *ServiceManager) ReadContractProvider(tnVzBrCPName string, application_epg string, application_profile string, tenant string) (*models.ContractProvider, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsprov-%s", tenant, application_profile, application_epg, tnVzBrCPName)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsProv := models.ContractProviderFromContainer(cont)
	return fvRsProv, nil
}

func (sm *ServiceManager) DeleteContractProvider(tnVzBrCPName string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsprov-%s", tenant, application_profile, application_epg, tnVzBrCPName)
	return sm.DeleteByDn(dn, models.FvrsprovClassName)
}

func (sm *ServiceManager) UpdateContractProvider(tnVzBrCPName string, application_epg string, application_profile string, tenant string, fvRsProvattr models.ContractProviderAttributes) (*models.ContractProvider, error) {
	rn := fmt.Sprintf("rsprov-%s", tnVzBrCPName)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsProv := models.NewContractProvider(rn, parentDn, fvRsProvattr)

	fvRsProv.Status = "modified"
	err := sm.Save(fvRsProv)
	return fvRsProv, err

}

func (sm *ServiceManager) ListContractProvider(application_epg string, application_profile string, tenant string) ([]*models.ContractProvider, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/fvRsProv.json", baseurlStr, tenant, application_profile, application_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ContractProviderListFromContainer(cont)

	return list, err
}
