package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAnnotation(key string, parentDn string, tagAnnotationAttr models.AnnotationAttributes) (*models.Annotation, error) {
	rn := fmt.Sprintf(models.RnTagAnnotation, key)
	tagAnnotation := models.NewAnnotation(rn, parentDn, tagAnnotationAttr)
	err := sm.Save(tagAnnotation)
	return tagAnnotation, err
}

func (sm *ServiceManager) ReadAnnotation(key string, parentDn string) (*models.Annotation, error) {
	dn := fmt.Sprintf(models.DnTagAnnotation, parentDn, key)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tagAnnotation := models.AnnotationFromContainer(cont)
	return tagAnnotation, nil
}

func (sm *ServiceManager) DeleteAnnotation(key string, parentDn string) error {
	dn := fmt.Sprintf(models.DnTagAnnotation, parentDn, key)
	return sm.DeleteByDn(dn, models.TagAnnotationClassName)
}

func (sm *ServiceManager) UpdateAnnotation(key string, parentDn string, tagAnnotationAttr models.AnnotationAttributes) (*models.Annotation, error) {
	rn := fmt.Sprintf(models.RnTagAnnotation, key)
	tagAnnotation := models.NewAnnotation(rn, parentDn, tagAnnotationAttr)
	tagAnnotation.Status = "modified"
	err := sm.Save(tagAnnotation)
	return tagAnnotation, err
}

func (sm *ServiceManager) ListAnnotation() ([]*models.Annotation, error) {
	dnUrl := fmt.Sprintf("%s/tagAnnotation.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AnnotationListFromContainer(cont)
	return list, err
}
