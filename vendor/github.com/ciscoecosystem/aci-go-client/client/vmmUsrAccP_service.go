package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVMMCredential(name string, vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmUsrAccPAttr models.VMMCredentialAttributes) (*models.VMMCredential, error) {
	rn := fmt.Sprintf(models.RnvmmUsrAccP, name)
	parentDn := fmt.Sprintf(models.ParentDnvmmUsrAccP, provider_profile_vendor, vmm_domain)
	vmmUsrAccP := models.NewVMMCredential(rn, parentDn, description, nameAlias, vmmUsrAccPAttr)
	err := sm.Save(vmmUsrAccP)
	return vmmUsrAccP, err
}

func (sm *ServiceManager) ReadVMMCredential(name string, vmm_domain string, provider_profile_vendor string) (*models.VMMCredential, error) {
	dn := fmt.Sprintf(models.DnvmmUsrAccP, provider_profile_vendor, vmm_domain, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmUsrAccP := models.VMMCredentialFromContainer(cont)
	return vmmUsrAccP, nil
}

func (sm *ServiceManager) DeleteVMMCredential(name string, vmm_domain string, provider_profile_vendor string) error {
	dn := fmt.Sprintf(models.DnvmmUsrAccP, provider_profile_vendor, vmm_domain, name)
	return sm.DeleteByDn(dn, models.VmmusraccpClassName)
}

func (sm *ServiceManager) UpdateVMMCredential(name string, vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmUsrAccPAttr models.VMMCredentialAttributes) (*models.VMMCredential, error) {
	rn := fmt.Sprintf(models.RnvmmUsrAccP, name)
	parentDn := fmt.Sprintf(models.ParentDnvmmUsrAccP, provider_profile_vendor, vmm_domain)
	vmmUsrAccP := models.NewVMMCredential(rn, parentDn, description, nameAlias, vmmUsrAccPAttr)
	vmmUsrAccP.Status = "modified"
	err := sm.Save(vmmUsrAccP)
	return vmmUsrAccP, err
}

func (sm *ServiceManager) ListVMMCredential(vmm_domain string, provider_profile_vendor string) ([]*models.VMMCredential, error) {
	dnUrl := fmt.Sprintf("%s/uni/vmmp-%s/dom-%s/vmmUsrAccP.json", models.BaseurlStr, provider_profile_vendor, vmm_domain)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.VMMCredentialListFromContainer(cont)
	return list, err
}
