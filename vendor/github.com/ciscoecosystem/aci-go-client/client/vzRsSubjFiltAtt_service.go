package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSubjectFilter(tnVzFilterName string, contract_subject string, contract string, tenant string, vzRsSubjFiltAttAttr models.SubjectFilterAttributes) (*models.SubjectFilter, error) {
	rn := fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName)
	parentDn := fmt.Sprintf(models.ParentDnvzRsSubjFiltAtt, tenant, contract, contract_subject)
	vzRsSubjFiltAtt := models.NewSubjectFilter(rn, parentDn, vzRsSubjFiltAttAttr)
	err := sm.Save(vzRsSubjFiltAtt)
	return vzRsSubjFiltAtt, err
}

func (sm *ServiceManager) ReadSubjectFilter(tnVzFilterName string, contract_subject string, contract string, tenant string) (*models.SubjectFilter, error) {
	dn := fmt.Sprintf(models.DnvzRsSubjFiltAtt, tenant, contract, contract_subject, tnVzFilterName)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzRsSubjFiltAtt := models.SubjectFilterFromContainer(cont)
	return vzRsSubjFiltAtt, nil
}

func (sm *ServiceManager) DeleteSubjectFilter(tnVzFilterName string, contract_subject string, contract string, tenant string) error {
	dn := fmt.Sprintf(models.DnvzRsSubjFiltAtt, tenant, contract, contract_subject, tnVzFilterName)
	return sm.DeleteByDn(dn, models.VzrssubjfiltattClassName)
}

func (sm *ServiceManager) UpdateSubjectFilter(tnVzFilterName string, contract_subject string, contract string, tenant string, vzRsSubjFiltAttAttr models.SubjectFilterAttributes) (*models.SubjectFilter, error) {
	rn := fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName)
	parentDn := fmt.Sprintf(models.ParentDnvzRsSubjFiltAtt, tenant, contract, contract_subject)
	vzRsSubjFiltAtt := models.NewSubjectFilter(rn, parentDn, vzRsSubjFiltAttAttr)
	vzRsSubjFiltAtt.Status = "modified"
	err := sm.Save(vzRsSubjFiltAtt)
	return vzRsSubjFiltAtt, err
}

func (sm *ServiceManager) ListSubjectFilter(contract_subject string, contract string, tenant string) ([]*models.SubjectFilter, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/brc-%s/subj-%s/vzRsSubjFiltAtt.json", models.BaseurlStr, tenant, contract, contract_subject)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SubjectFilterListFromContainer(cont)
	return list, err
}
