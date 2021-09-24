package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOOBManagedNodesZone(managed_node_connectivity_group string, description string, nameAlias string, mgmtOoBZoneAttr models.OOBManagedNodesZoneAttributes) (*models.OOBManagedNodesZone, error) {
	rn := fmt.Sprintf(models.RnmgmtOoBZone)
	parentDn := fmt.Sprintf(models.ParentDnmgmtOoBZone, managed_node_connectivity_group)
	mgmtOoBZone := models.NewOOBManagedNodesZone(rn, parentDn, description, nameAlias, mgmtOoBZoneAttr)
	err := sm.Save(mgmtOoBZone)
	return mgmtOoBZone, err
}

func (sm *ServiceManager) ReadOOBManagedNodesZone(managed_node_connectivity_group string) (*models.OOBManagedNodesZone, error) {
	dn := fmt.Sprintf(models.DnmgmtOoBZone, managed_node_connectivity_group)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtOoBZone := models.OOBManagedNodesZoneFromContainer(cont)
	return mgmtOoBZone, nil
}

func (sm *ServiceManager) DeleteOOBManagedNodesZone(managed_node_connectivity_group string) error {
	dn := fmt.Sprintf(models.DnmgmtOoBZone, managed_node_connectivity_group)
	return sm.DeleteByDn(dn, models.MgmtoobzoneClassName)
}

func (sm *ServiceManager) UpdateOOBManagedNodesZone(managed_node_connectivity_group string, description string, nameAlias string, mgmtOoBZoneAttr models.OOBManagedNodesZoneAttributes) (*models.OOBManagedNodesZone, error) {
	rn := fmt.Sprintf(models.RnmgmtOoBZone)
	parentDn := fmt.Sprintf(models.ParentDnmgmtOoBZone, managed_node_connectivity_group)
	mgmtOoBZone := models.NewOOBManagedNodesZone(rn, parentDn, description, nameAlias, mgmtOoBZoneAttr)
	mgmtOoBZone.Status = "modified"
	err := sm.Save(mgmtOoBZone)
	return mgmtOoBZone, err
}

func (sm *ServiceManager) ListOOBManagedNodesZone(managed_node_connectivity_group string) ([]*models.OOBManagedNodesZone, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/funcprof/grp-%s/mgmtOoBZone.json", models.BaseurlStr, managed_node_connectivity_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.OOBManagedNodesZoneListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationmgmtRsAddrInst(parentDn, annotation, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationmgmtRsAddrInst(parentDn string) error {
	dn := fmt.Sprintf("%s/rsaddrInst", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsAddrInst")
}

func (sm *ServiceManager) ReadRelationmgmtRsAddrInst(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) CreateRelationmgmtRsOoB(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsooB", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "mgmtRsOoB", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationmgmtRsOoB(parentDn string) error {
	dn := fmt.Sprintf("%s/rsooB", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsOoB")
}

func (sm *ServiceManager) ReadRelationmgmtRsOoB(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "mgmtRsOoB")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "mgmtRsOoB")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationmgmtRsOobEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsoobEpg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "mgmtRsOobEpg", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationmgmtRsOobEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsoobEpg", parentDn)
	return sm.DeleteByDn(dn, "mgmtRsOobEpg")
}

func (sm *ServiceManager) ReadRelationmgmtRsOobEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "mgmtRsOobEpg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "mgmtRsOobEpg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
