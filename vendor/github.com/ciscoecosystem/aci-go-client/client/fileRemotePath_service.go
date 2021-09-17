package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRemotePathofaFile(name string, description string, nameAlias string, fileRemotePathAttr models.RemotePathofaFileAttributes) (*models.RemotePathofaFile, error) {
	rn := fmt.Sprintf(models.RnfileRemotePath, name)
	parentDn := fmt.Sprintf(models.ParentDnfileRemotePath)
	fileRemotePath := models.NewRemotePathofaFile(rn, parentDn, description, nameAlias, fileRemotePathAttr)
	err := sm.Save(fileRemotePath)
	return fileRemotePath, err
}

func (sm *ServiceManager) ReadRemotePathofaFile(name string) (*models.RemotePathofaFile, error) {
	dn := fmt.Sprintf(models.DnfileRemotePath, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fileRemotePath := models.RemotePathofaFileFromContainer(cont)
	return fileRemotePath, nil
}

func (sm *ServiceManager) DeleteRemotePathofaFile(name string) error {
	dn := fmt.Sprintf(models.DnfileRemotePath, name)
	return sm.DeleteByDn(dn, models.FileremotepathClassName)
}

func (sm *ServiceManager) UpdateRemotePathofaFile(name string, description string, nameAlias string, fileRemotePathAttr models.RemotePathofaFileAttributes) (*models.RemotePathofaFile, error) {
	rn := fmt.Sprintf(models.RnfileRemotePath, name)
	parentDn := fmt.Sprintf(models.ParentDnfileRemotePath)
	fileRemotePath := models.NewRemotePathofaFile(rn, parentDn, description, nameAlias, fileRemotePathAttr)
	fileRemotePath.Status = "modified"
	err := sm.Save(fileRemotePath)
	return fileRemotePath, err
}

func (sm *ServiceManager) ListRemotePathofaFile() ([]*models.RemotePathofaFile, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/fileRemotePath.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RemotePathofaFileListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationfileRsARemoteHostToEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsARemoteHostToEpg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "fileRsARemoteHostToEpg", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationfileRsARemoteHostToEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsARemoteHostToEpg", parentDn)
	return sm.DeleteByDn(dn, "fileRsARemoteHostToEpg")
}

func (sm *ServiceManager) ReadRelationfileRsARemoteHostToEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fileRsARemoteHostToEpg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fileRsARemoteHostToEpg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationfileRsARemoteHostToEpp(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsARemoteHostToEpp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "fileRsARemoteHostToEpp", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationfileRsARemoteHostToEpp(parentDn string) error {
	dn := fmt.Sprintf("%s/rsARemoteHostToEpp", parentDn)
	return sm.DeleteByDn(dn, "fileRsARemoteHostToEpp")
}

func (sm *ServiceManager) ReadRelationfileRsARemoteHostToEpp(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fileRsARemoteHostToEpp")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fileRsARemoteHostToEpp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
