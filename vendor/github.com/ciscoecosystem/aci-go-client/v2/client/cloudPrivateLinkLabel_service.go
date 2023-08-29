package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudPrivateLinkLabel(name string, parentDn string, description string, cloudPrivateLinkLabelAttr models.CloudPrivateLinkLabelAttributes) (*models.CloudPrivateLinkLabel, error) {

	rn := fmt.Sprintf(models.RnCloudPrivateLinkLabel, name)
	cloudPrivateLinkLabel := models.NewCloudPrivateLinkLabel(rn, parentDn, description, cloudPrivateLinkLabelAttr)

	err := sm.Save(cloudPrivateLinkLabel)
	return cloudPrivateLinkLabel, err
}

func (sm *ServiceManager) ReadCloudPrivateLinkLabel(name string, parentDn string) (*models.CloudPrivateLinkLabel, error) {

	rn := fmt.Sprintf(models.RnCloudPrivateLinkLabel, name)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudPrivateLinkLabel := models.CloudPrivateLinkLabelFromContainer(cont)
	return cloudPrivateLinkLabel, nil
}

func (sm *ServiceManager) DeleteCloudPrivateLinkLabel(name string, parentDn string) error {

	rn := fmt.Sprintf(models.RnCloudPrivateLinkLabel, name)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.CloudPrivateLinkLabelClassName)
}

func (sm *ServiceManager) UpdateCloudPrivateLinkLabel(name string, parentDn string, description string, cloudPrivateLinkLabelAttr models.CloudPrivateLinkLabelAttributes) (*models.CloudPrivateLinkLabel, error) {

	rn := fmt.Sprintf(models.RnCloudPrivateLinkLabel, name)
	cloudPrivateLinkLabel := models.NewCloudPrivateLinkLabel(rn, parentDn, description, cloudPrivateLinkLabelAttr)

	cloudPrivateLinkLabel.Status = "modified"
	err := sm.Save(cloudPrivateLinkLabel)
	return cloudPrivateLinkLabel, err
}

func (sm *ServiceManager) ListCloudPrivateLinkLabel(parentDn string) ([]*models.CloudPrivateLinkLabel, error) {

	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.CloudPrivateLinkLabelClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudPrivateLinkLabelListFromContainer(cont)
	return list, err
}
