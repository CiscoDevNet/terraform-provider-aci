package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateBFDMultihopNodePolicy(name string, tenant string, description string, bfdMhNodePolAttr models.BFDMultihopNodePolicyAttributes) (*models.BFDMultihopNodePolicy, error) {

	rn := fmt.Sprintf(models.RnBfdMhNodePol, name)

	parentDn := fmt.Sprintf(models.ParentDnBfdMhNodePol, tenant)
	bfdMhNodePol := models.NewBFDMultihopNodePolicy(rn, parentDn, description, bfdMhNodePolAttr)

	err := sm.Save(bfdMhNodePol)
	return bfdMhNodePol, err
}

func (sm *ServiceManager) ReadBFDMultihopNodePolicy(name string, tenant string) (*models.BFDMultihopNodePolicy, error) {

	rn := fmt.Sprintf(models.RnBfdMhNodePol, name)

	parentDn := fmt.Sprintf(models.ParentDnBfdMhNodePol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	bfdMhNodePol := models.BFDMultihopNodePolicyFromContainer(cont)
	return bfdMhNodePol, nil
}

func (sm *ServiceManager) DeleteBFDMultihopNodePolicy(name string, tenant string) error {

	rn := fmt.Sprintf(models.RnBfdMhNodePol, name)

	parentDn := fmt.Sprintf(models.ParentDnBfdMhNodePol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.BfdMhNodePolClassName)
}

func (sm *ServiceManager) UpdateBFDMultihopNodePolicy(name string, tenant string, description string, bfdMhNodePolAttr models.BFDMultihopNodePolicyAttributes) (*models.BFDMultihopNodePolicy, error) {

	rn := fmt.Sprintf(models.RnBfdMhNodePol, name)

	parentDn := fmt.Sprintf(models.ParentDnBfdMhNodePol, tenant)
	bfdMhNodePol := models.NewBFDMultihopNodePolicy(rn, parentDn, description, bfdMhNodePolAttr)

	bfdMhNodePol.Status = "modified"
	err := sm.Save(bfdMhNodePol)
	return bfdMhNodePol, err
}

func (sm *ServiceManager) ListBFDMultihopNodePolicy(tenant string) ([]*models.BFDMultihopNodePolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnBfdMhNodePol, tenant)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.BfdMhNodePolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BFDMultihopNodePolicyListFromContainer(cont)
	return list, err
}
