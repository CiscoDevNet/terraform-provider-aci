package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciISISDomainPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_isis_domain_policy.test"
	dataSourceName := "data.aci_isis_domain_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	isisDomPol, err := aci.GetRemoteISISDomainPolicy(sharedAciClient(), "uni/fabric/isisDomP-default")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}
	isisLvlComp, err := aci.GetRemoteISISLevel(sharedAciClient(), "uni/fabric/isisDomP-default/lvl-l1")
	if err != nil {
		t.Errorf("reading initial config of isisLvlComp")
	}
	fmt.Println(*isisLvlComp)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "isis_level_name", resourceName, "isis_level_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mtu", resourceName, "mtu"),
					resource.TestCheckResourceAttrPair(dataSourceName, "redistrib_metric", resourceName, "redistrib_metric"),
					resource.TestCheckResourceAttrPair(dataSourceName, "isis_level_type", resourceName, "isis_level_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsp_fast_flood", resourceName, "lsp_fast_flood"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsp_gen_init_intvl", resourceName, "lsp_gen_init_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsp_gen_max_intvl", resourceName, "lsp_gen_max_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lsp_gen_sec_intvl", resourceName, "lsp_gen_sec_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_comp_init_intvl", resourceName, "spf_comp_init_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_comp_max_intvl", resourceName, "spf_comp_max_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spf_comp_sec_intvl", resourceName, "spf_comp_sec_intvl"),
				),
			},
			{
				Config:      CreateAccISISDomainPolicyDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccISISDomainPolicyDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreISISDomainPolicyToInitConfig(isisDomPol, isisLvlComp),
			},
		},
	})
}

func CreateAccISISDomainPolicyConfigDataSource() string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	}

	data "aci_isis_domain_policy" "test" {
	}
	`)
	return resource
}

func CreateAccISISDomainPolicyDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	}

	data "aci_isis_domain_policy" "test" {
		%s = "%s"
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, key, value)
	return resource
}

func CreateAccISISDomainPolicyDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
		%s = "%s"
	}

	data "aci_isis_domain_policy" "test" {
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, key, value)
	return resource
}
