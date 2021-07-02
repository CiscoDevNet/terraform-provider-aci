package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateConnection(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsConnectionattr models.ConnectionAttributes) (*models.Connection, error) {
	rn := fmt.Sprintf("AbsConnection-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsConnection := models.NewConnection(rn, parentDn, description, vnsAbsConnectionattr)
	err := sm.Save(vnsAbsConnection)
	return vnsAbsConnection, err
}

func (sm *ServiceManager) ReadConnection(name string, l4_l7_service_graph_template string, tenant string) (*models.Connection, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsConnection-%s", tenant, l4_l7_service_graph_template, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsConnection := models.ConnectionFromContainer(cont)
	return vnsAbsConnection, nil
}

func (sm *ServiceManager) DeleteConnection(name string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsConnection-%s", tenant, l4_l7_service_graph_template, name)
	return sm.DeleteByDn(dn, models.VnsabsconnectionClassName)
}

func (sm *ServiceManager) UpdateConnection(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsConnectionattr models.ConnectionAttributes) (*models.Connection, error) {
	rn := fmt.Sprintf("AbsConnection-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsConnection := models.NewConnection(rn, parentDn, description, vnsAbsConnectionattr)

	vnsAbsConnection.Status = "modified"
	err := sm.Save(vnsAbsConnection)
	return vnsAbsConnection, err

}

func (sm *ServiceManager) ListConnection(l4_l7_service_graph_template string, tenant string) ([]*models.Connection, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/vnsAbsConnection.json", baseurlStr, tenant, l4_l7_service_graph_template)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConnectionListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsAbsCopyConnectionFromConnection(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsabsCopyConnection-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s", "annotation":"orchestrator:terraform"		
			}
		}
	}`, "vnsRsAbsCopyConnection", dn))

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

	return CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) DeleteRelationvnsRsAbsCopyConnectionFromConnection(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsabsCopyConnection-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "vnsRsAbsCopyConnection")
}

func (sm *ServiceManager) ReadRelationvnsRsAbsCopyConnectionFromConnection(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsAbsCopyConnection")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsAbsCopyConnection")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationvnsRsAbsConnectionConnsFromConnection(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsabsConnectionConns-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s", "annotation":"orchestrator:terraform"			
			}
		}
	}`, "vnsRsAbsConnectionConns", dn))

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

	return CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) DeleteRelationvnsRsAbsConnectionConnsFromConnection(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsabsConnectionConns-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "vnsRsAbsConnectionConns")
}

func (sm *ServiceManager) ReadRelationvnsRsAbsConnectionConnsFromConnection(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsAbsConnectionConns")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsAbsConnectionConns")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
