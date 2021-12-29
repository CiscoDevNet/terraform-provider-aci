package testacc

import (
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
