package testacc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("aci_tenant",
		&resource.Sweeper{
			Name: "aci_tenant",
			F:    aciTenantSweeper,
		})
	resource.AddTestSweepers("aci_fabric_if_pol",
		&resource.Sweeper{
			Name: "aci_fabric_if_pol",
			F:    aciFabricIfPolSweeper,
		})
	resource.AddTestSweepers("aci_interface_fc_policy",
		&resource.Sweeper{
			Name: "aci_interface_fc_policy",
			F:    aciInterfaceFCPolicySweeper,
		})
	resource.AddTestSweepers("aci_leaf_interface_profile",
		&resource.Sweeper{
			Name: "aci_leaf_interface_profile",
			F:    aciLeafInterfaceProfileSweeper,
		})
	resource.AddTestSweepers("aci_leaf_breakout_port_group",
		&resource.Sweeper{
			Name: "aci_leaf_breakout_port_group",
			F:    aciLeafBreakoutPortGroupSweeper,
		})
	resource.AddTestSweepers("aci_lacp_policy",
		&resource.Sweeper{
			Name: "aci_lacp_policy",
			F:    aciLeafBreakoutPortGroupSweeper,
		})
	resource.AddTestSweepers("aci_radius_provider",
		&resource.Sweeper{
			Name: "aci_radius_provider",
			F:    aciRadiusProviderSweeper,
		})
}

func aciLeafBreakoutPortGroupSweeper(_ string) error {
	className := "infraBrkoutPortGrp"
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL(fmt.Sprintf("/api/node/class/%s.json", className))
	instances := models.LeafBreakoutPortGroupListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, className)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciRadiusProviderSweeper(_ string) error {
	className := "aaaRadiusProvider"
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL(fmt.Sprintf("/api/node/class/%s.json", className))
	instances := models.RADIUSProviderListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, className)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciLACPPolicySweeper(_ string) error {
	className := "lacpLagPol"
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL(fmt.Sprintf("/api/node/class/%s.json", className))
	instances := models.LACPPolicyListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, className)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciTenantSweeper(_ string) error {
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL("/api/node/class/fvTenant.json")
	instances := models.TenantListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, "fvTenant")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciFabricIfPolSweeper(_ string) error {
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL("/api/node/class/fabricHIfPol.json")
	instances := models.LinkLevelPolicyListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, "fabricHIfPol")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciInterfaceFCPolicySweeper(_ string) error {
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL("/api/node/class/fcIfPol.json")
	instances := models.InterfaceFCPolicyListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, "fcIfPol")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func aciLeafInterfaceProfileSweeper(_ string) error {
	aciClient := sharedAciClient()
	cont, _ := aciClient.GetViaURL("/api/node/class/infraAccPortP.json")
	instances := models.LeafInterfaceProfileListFromContainer(cont)
	for _, instance := range instances {
		dn := instance.DistinguishedName
		if strings.HasPrefix(GetMOName(dn), "acctest") {
			err := aciClient.DeleteByDn(dn, "infraAccPortP")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
