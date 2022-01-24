package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFvCEpDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_client_end_point.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFvCEpConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "fvcep_objects.0.application_profile_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fvcep_objects.0.epg_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fvcep_objects.0.mac"),
					resource.TestCheckResourceAttr(dataSourceName, "fvcep_objects.0.tenant_name", "infra"),
				),
			},
			{
				Config:      CreateAccFvCEpDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFvCEpConfigDataSource(),
			},
		},
	})
}

func CreateAccFvCEpDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing client_end_point Data Source with random parameters")
	resource := fmt.Sprintf(`
	data "aci_client_end_point" "test" {
		%s = "%s"
	}
	`, key, value)
	return resource
}

func CreateAccFvCEpConfigDataSource() string {
	fmt.Println("=== STEP  testing client_end_point Data Source with required arguments only")
	resource := fmt.Sprintln(`
	data "aci_client_end_point" "test" {}
	`)
	return resource
}
