package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudTemplateRegion(parentDn string, cloudtemplateRegionDetailAttr models.CloudTemplateRegionAttributes) (*models.CloudTemplateRegion, error) {

	cloudtemplateRegionDetail := models.NewCloudTemplateRegion(models.RnCloudtemplateRegionDetail, parentDn, cloudtemplateRegionDetailAttr)

	err := sm.Save(cloudtemplateRegionDetail)
	return cloudtemplateRegionDetail, err
}

func (sm *ServiceManager) ReadCloudTemplateRegion(parentDn string) (*models.CloudTemplateRegion, error) {

	dn := fmt.Sprintf("%s/%s", parentDn, models.RnCloudtemplateRegionDetail)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateRegionDetail := models.CloudTemplateRegionFromContainer(cont)
	return cloudtemplateRegionDetail, nil
}

func (sm *ServiceManager) DeleteCloudTemplateRegion(parentDn string) error {

	dn := fmt.Sprintf("%s/%s", parentDn, models.RnCloudtemplateRegionDetail)

	return sm.DeleteByDn(dn, models.CloudtemplateRegionDetailClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateRegion(parentDn string, cloudtemplateRegionDetailAttr models.CloudTemplateRegionAttributes) (*models.CloudTemplateRegion, error) {

	cloudtemplateRegionDetail := models.NewCloudTemplateRegion(models.RnCloudtemplateRegionDetail, parentDn, cloudtemplateRegionDetailAttr)

	cloudtemplateRegionDetail.Status = "modified"
	err := sm.Save(cloudtemplateRegionDetail)
	return cloudtemplateRegionDetail, err
}

func (sm *ServiceManager) ListCloudTemplateRegion(parentDn string) ([]*models.CloudTemplateRegion, error) {

	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.CloudtemplateRegionDetailClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateRegionListFromContainer(cont)
	return list, err
}
