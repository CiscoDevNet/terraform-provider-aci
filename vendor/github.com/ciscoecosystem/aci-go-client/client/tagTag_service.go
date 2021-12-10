package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTag(key string, parentDn string, tagTagAttr models.TagAttributes) (*models.Tag, error) {
	rn := fmt.Sprintf(models.RnTagTag, key)
	tagTag := models.NewTag(rn, parentDn, tagTagAttr)
	err := sm.Save(tagTag)
	return tagTag, err
}

func (sm *ServiceManager) ReadTag(key string, parentDn string) (*models.Tag, error) {
	dn := fmt.Sprintf(models.DnTagTag, parentDn, key)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tagTag := models.TagFromContainer(cont)
	return tagTag, nil
}

func (sm *ServiceManager) DeleteTag(key string, parentDn string) error {
	dn := fmt.Sprintf(models.DnTagTag, parentDn, key)
	return sm.DeleteByDn(dn, models.TagTagClassName)
}

func (sm *ServiceManager) UpdateTag(key string, parentDn string, tagTagAttr models.TagAttributes) (*models.Tag, error) {
	rn := fmt.Sprintf(models.RnTagTag, key)
	tagTag := models.NewTag(rn, parentDn, tagTagAttr)
	tagTag.Status = "modified"
	err := sm.Save(tagTag)
	return tagTag, err
}

func (sm *ServiceManager) ListTag() ([]*models.Tag, error) {
	dnUrl := fmt.Sprintf("%s/tagTag.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TagListFromContainer(cont)
	return list, err
}
