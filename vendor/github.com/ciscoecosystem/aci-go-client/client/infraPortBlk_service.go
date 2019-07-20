package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ciscoecosystem/aci-go-client/container"



	


)









func (sm *ServiceManager) CreateAccessPortBlock(name string ,access_port_selector_type string , access_port_selector string ,leaf_interface_profile string , description string, infraPortBlkattr models.AccessPortBlockAttributes) (*models.AccessPortBlock, error) {	
	rn := fmt.Sprintf("portblk-%s",name)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile ,access_port_selector , access_port_selector_type )
	infraPortBlk := models.NewAccessPortBlock(rn, parentDn, description, infraPortBlkattr)
	err := sm.Save(infraPortBlk)
	return infraPortBlk, err
}

func (sm *ServiceManager) ReadAccessPortBlock(name string ,access_port_selector_type string , access_port_selector string ,leaf_interface_profile string ) (*models.AccessPortBlock, error) {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/portblk-%s", leaf_interface_profile ,access_port_selector ,access_port_selector_type ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraPortBlk := models.AccessPortBlockFromContainer(cont)
	return infraPortBlk, nil
}

func (sm *ServiceManager) DeleteAccessPortBlock(name string ,access_port_selector_type string , access_port_selector string ,leaf_interface_profile string ) error {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/portblk-%s", leaf_interface_profile ,access_port_selector ,access_port_selector_type ,name )
	return sm.DeleteByDn(dn, models.InfraportblkClassName)
}

func (sm *ServiceManager) UpdateAccessPortBlock(name string ,access_port_selector_type string , access_port_selector string ,leaf_interface_profile string  ,description string, infraPortBlkattr models.AccessPortBlockAttributes) (*models.AccessPortBlock, error) {
	rn := fmt.Sprintf("portblk-%s",name)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile ,access_port_selector , access_port_selector_type )
	infraPortBlk := models.NewAccessPortBlock(rn, parentDn, description, infraPortBlkattr)

    infraPortBlk.Status = "modified"
	err := sm.Save(infraPortBlk)
	return infraPortBlk, err

}

func (sm *ServiceManager) ListAccessPortBlock(access_port_selector_type string , access_port_selector string ,leaf_interface_profile string ) ([]*models.AccessPortBlock, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/infra/accportprof-%s/hports-%s-typ-%s/infraPortBlk.json", baseurlStr , leaf_interface_profile ,access_port_selector , access_port_selector_type )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessPortBlockListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsAccBndlSubgrpFromAccessPortBlock( parentDn, tnInfraAccBndlSubgrpName string) error {
	dn := fmt.Sprintf("%s/rsaccBndlSubgrp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
			}
		}
	}`, "infraRsAccBndlSubgrp", dn,tnInfraAccBndlSubgrpName))

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

func (sm *ServiceManager) DeleteRelationinfraRsAccBndlSubgrpFromAccessPortBlock(parentDn string) error{
	dn := fmt.Sprintf("%s/rsaccBndlSubgrp", parentDn)
	return sm.DeleteByDn(dn , "infraRsAccBndlSubgrp")
}

func (sm *ServiceManager) ReadRelationinfraRsAccBndlSubgrpFromAccessPortBlock( parentDn string) (interface{},error) {
	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",baseurlStr,parentDn,"infraRsAccBndlSubgrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont,"infraRsAccBndlSubgrp")
	
	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil,err
	}
		





}

