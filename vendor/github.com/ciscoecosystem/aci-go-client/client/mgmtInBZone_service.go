package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateInBManagedNodesZone(managed_node_connectivity_group string, description string, nameAlias string, mgmtInBZoneAttr models.InBManagedNodesZoneAttributes) (*models.InBManagedNodesZone, error) {
	rn := fmt.Sprintf(models.RnmgmtInBZone)
	parentDn := fmt.Sprintf(models.ParentDnmgmtInBZone, managed_node_connectivity_group)
	mgmtInBZone := models.NewInBManagedNodesZone(rn, parentDn, description, nameAlias, mgmtInBZoneAttr)
	err := sm.Save(mgmtInBZone)
	return mgmtInBZone, err
}

func (sm *ServiceManager) ReadInBManagedNodesZone(managed_node_connectivity_group string) (*models.InBManagedNodesZone, error) {
	dn := fmt.Sprintf(models.DnmgmtInBZone, managed_node_connectivity_group)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtInBZone := models.InBManagedNodesZoneFromContainer(cont)
	return mgmtInBZone, nil
}

func (sm *ServiceManager) DeleteInBManagedNodesZone(managed_node_connectivity_group string) error {
	dn := fmt.Sprintf(models.DnmgmtInBZone, managed_node_connectivity_group)
	return sm.DeleteByDn(dn, models.MgmtinbzoneClassName)
}

func (sm *ServiceManager) UpdateInBManagedNodesZone(managed_node_connectivity_group string, description string, nameAlias string, mgmtInBZoneAttr models.InBManagedNodesZoneAttributes) (*models.InBManagedNodesZone, error) {
	rn := fmt.Sprintf(models.RnmgmtInBZone)
	parentDn := fmt.Sprintf(models.ParentDnmgmtInBZone, managed_node_connectivity_group)
	mgmtInBZone := models.NewInBManagedNodesZone(rn, parentDn, description, nameAlias, mgmtInBZoneAttr)
	mgmtInBZone.Status = "modified"
	err := sm.Save(mgmtInBZone)
	return mgmtInBZone, err
}

func (sm *ServiceManager) ListInBManagedNodesZone(managed_node_connectivity_group string) ([]*models.InBManagedNodesZone, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/funcprof/grp-%s/mgmtInBZone.json", models.BaseurlStr, managed_node_connectivity_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.InBManagedNodesZoneListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationmgmtRsAddrInstFrommgmtInBZone(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsaddrInst", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "mgmtRsAddrInst", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationmgmtRsAddrInstFrommgmtInBZone(parentDn string) error {
	dn := fmt.Sprintf("%s/rsaddrInst", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsAddrInst")
}

func (sm *ServiceManager) ReadRelationmgmtRsAddrInstFrommgmtInBZone(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "mgmtRsAddrInst")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "mgmtRsAddrInst")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationmgmtRsInB(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsinB", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "mgmtRsInB", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationmgmtRsInB(parentDn string) error {
	dn := fmt.Sprintf("%s/rsinB", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsInB")
}

func (sm *ServiceManager) ReadRelationmgmtRsInB(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "mgmtRsInB")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "mgmtRsInB")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationmgmtRsInbEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsinbEpg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "mgmtRsInbEpg", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationmgmtRsInbEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsinbEpg", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsInbEpg")
}

func (sm *ServiceManager) ReadRelationmgmtRsInbEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "mgmtRsInbEpg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "mgmtRsInbEpg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
