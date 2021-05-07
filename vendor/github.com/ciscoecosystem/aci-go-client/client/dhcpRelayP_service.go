package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateDHCPRelayPolicy(name string, tenant string, description string, dhcpRelayPattr models.DHCPRelayPolicyAttributes) (*models.DHCPRelayPolicy, error) {
	rn := fmt.Sprintf("relayp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	dhcpRelayP := models.NewDHCPRelayPolicy(rn, parentDn, description, dhcpRelayPattr)
	err := sm.Save(dhcpRelayP)
	return dhcpRelayP, err
}

func (sm *ServiceManager) ReadDHCPRelayPolicy(name string, tenant string) (*models.DHCPRelayPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/relayp-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpRelayP := models.DHCPRelayPolicyFromContainer(cont)
	return dhcpRelayP, nil
}

func (sm *ServiceManager) DeleteDHCPRelayPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/relayp-%s", tenant, name)
	return sm.DeleteByDn(dn, models.DhcprelaypClassName)
}

func (sm *ServiceManager) UpdateDHCPRelayPolicy(name string, tenant string, description string, dhcpRelayPattr models.DHCPRelayPolicyAttributes) (*models.DHCPRelayPolicy, error) {
	rn := fmt.Sprintf("relayp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	dhcpRelayP := models.NewDHCPRelayPolicy(rn, parentDn, description, dhcpRelayPattr)

	dhcpRelayP.Status = "modified"
	err := sm.Save(dhcpRelayP)
	return dhcpRelayP, err

}

func (sm *ServiceManager) ListDHCPRelayPolicy(tenant string) ([]*models.DHCPRelayPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/dhcpRelayP.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.DHCPRelayPolicyListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationdhcpRsProvFromDHCPRelayPolicy(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsprov-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s", "annotation": "orchestrator:terraform", "tDn": "%s"				
			}
		}
	}`, "dhcpRsProv", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationdhcpRsProvFromDHCPRelayPolicy(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsprov-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "dhcpRsProv")
}

func (sm *ServiceManager) ReadRelationdhcpRsProvFromDHCPRelayPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "dhcpRsProv")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "dhcpRsProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
