package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateContractInterfaceRelationship(tnVzCPIfName string, application_epg string, application_profile string, tenant string, fvRsConsIfAttr models.ContractInterfaceRelationshipAttributes) (*models.ContractInterfaceRelationship, error) {
	rn := fmt.Sprintf(models.RnfvRsConsIf, tnVzCPIfName)
	parentDn := fmt.Sprintf(models.ParentDnfvRsConsIf, tenant, application_profile, application_epg)
	fvRsConsIf := models.NewContractInterfaceRelationship(rn, parentDn, fvRsConsIfAttr)
	err := sm.Save(fvRsConsIf)
	return fvRsConsIf, err
}

func (sm *ServiceManager) ReadContractInterfaceRelationship(tnVzCPIfName string, application_epg string, application_profile string, tenant string) (*models.ContractInterfaceRelationship, error) {
	dn := fmt.Sprintf(models.DnfvRsConsIf, tenant, application_profile, application_epg, tnVzCPIfName)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsConsIf := models.ContractInterfaceRelationshipFromContainer(cont)
	return fvRsConsIf, nil
}

func (sm *ServiceManager) DeleteContractInterfaceRelationship(tnVzCPIfName string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvRsConsIf, tenant, application_profile, application_epg, tnVzCPIfName)
	return sm.DeleteByDn(dn, models.FvrsconsifClassName)
}

func (sm *ServiceManager) UpdateContractInterfaceRelationship(tnVzCPIfName string, application_epg string, application_profile string, tenant string, fvRsConsIfAttr models.ContractInterfaceRelationshipAttributes) (*models.ContractInterfaceRelationship, error) {
	rn := fmt.Sprintf(models.RnfvRsConsIf, tnVzCPIfName)
	parentDn := fmt.Sprintf(models.ParentDnfvRsConsIf, tenant, application_profile, application_epg)
	fvRsConsIf := models.NewContractInterfaceRelationship(rn, parentDn, fvRsConsIfAttr)
	fvRsConsIf.Status = "modified"
	err := sm.Save(fvRsConsIf)
	return fvRsConsIf, err
}

func (sm *ServiceManager) ListContractInterfaceRelationship(application_epg string, application_profile string, tenant string) ([]*models.ContractInterfaceRelationship, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/fvRsConsIf.json", models.BaseurlStr, tenant, application_profile, application_epg)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ContractInterfaceRelationshipListFromContainer(cont)
	return list, err
}
