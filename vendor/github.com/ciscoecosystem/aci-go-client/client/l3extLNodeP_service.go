package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	



	


)









func (sm *ServiceManager) CreateLogicalNodeProfile(name string ,l3_outside string ,tenant string , description string, l3extLNodePattr models.LogicalNodeProfileAttributes) (*models.LogicalNodeProfile, error) {	
	rn := fmt.Sprintf("lnodep-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant ,l3_outside )
	l3extLNodeP := models.NewLogicalNodeProfile(rn, parentDn, description, l3extLNodePattr)
	err := sm.Save(l3extLNodeP)
	return l3extLNodeP, err
}

func (sm *ServiceManager) ReadLogicalNodeProfile(name string ,l3_outside string ,tenant string ) (*models.LogicalNodeProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant ,l3_outside ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLNodeP := models.LogicalNodeProfileFromContainer(cont)
	return l3extLNodeP, nil
}

func (sm *ServiceManager) DeleteLogicalNodeProfile(name string ,l3_outside string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant ,l3_outside ,name )
	return sm.DeleteByDn(dn, models.L3extlnodepClassName)
}

func (sm *ServiceManager) UpdateLogicalNodeProfile(name string ,l3_outside string ,tenant string  ,description string, l3extLNodePattr models.LogicalNodeProfileAttributes) (*models.LogicalNodeProfile, error) {
	rn := fmt.Sprintf("lnodep-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant ,l3_outside )
	l3extLNodeP := models.NewLogicalNodeProfile(rn, parentDn, description, l3extLNodePattr)

    l3extLNodeP.Status = "modified"
	err := sm.Save(l3extLNodeP)
	return l3extLNodeP, err

}

func (sm *ServiceManager) ListLogicalNodeProfile(l3_outside string ,tenant string ) ([]*models.LogicalNodeProfile, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/l3extLNodeP.json", baseurlStr , tenant ,l3_outside )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.LogicalNodeProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsNodeL3OutAttFromLogicalNodeProfile( parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsnodeL3OutAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "l3extRsNodeL3OutAtt", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(parentDn , tDn string) error{
	dn := fmt.Sprintf("%s/rsnodeL3OutAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn , "l3extRsNodeL3OutAtt")
}

func (sm *ServiceManager) ReadRelationl3extRsNodeL3OutAttFromLogicalNodeProfile( parentDn string) (interface{},error) {
	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",baseurlStr,parentDn,"l3extRsNodeL3OutAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont,"l3extRsNodeL3OutAtt")
	
	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList{
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err
			





}

