package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateLACPMemberPolicy(name string, description string, lacpIfPolAttr models.LACPMemberPolicyAttributes) (*models.LACPMemberPolicy, error) {
	rn := fmt.Sprintf(models.RnlacpIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnlacpIfPol)
	lacpIfPol := models.NewLACPMemberPolicy(rn, parentDn, description, lacpIfPolAttr)
	err := sm.Save(lacpIfPol)
	return lacpIfPol, err
}

func (sm *ServiceManager) ReadLACPMemberPolicy(name string) (*models.LACPMemberPolicy, error) {
	dn := fmt.Sprintf(models.DnlacpIfPol, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	lacpIfPol := models.LACPMemberPolicyFromContainer(cont)
	return lacpIfPol, nil
}

func (sm *ServiceManager) DeleteLACPMemberPolicy(name string) error {
	dn := fmt.Sprintf(models.DnlacpIfPol, name)
	return sm.DeleteByDn(dn, models.LacpifpolClassName)
}

func (sm *ServiceManager) UpdateLACPMemberPolicy(name string, description string, lacpIfPolAttr models.LACPMemberPolicyAttributes) (*models.LACPMemberPolicy, error) {
	rn := fmt.Sprintf(models.RnlacpIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnlacpIfPol)
	lacpIfPol := models.NewLACPMemberPolicy(rn, parentDn, description, lacpIfPolAttr)
	lacpIfPol.Status = "modified"
	err := sm.Save(lacpIfPol)
	return lacpIfPol, err
}

func (sm *ServiceManager) ListLACPMemberPolicy() ([]*models.LACPMemberPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/lacpIfPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LACPMemberPolicyListFromContainer(cont)
	return list, err
}
