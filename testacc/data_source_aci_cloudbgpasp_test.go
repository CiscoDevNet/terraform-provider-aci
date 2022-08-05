package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAutonomousSystemProfileDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_autonomous_system_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAutonomousSystemProfileConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "asn"),
				),
			},
			{
				Config:      CreateAccAutonomousSystemProfileDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAutonomousSystemProfileConfigDataSource(),
			},
		},
	})
}

func CreateAccAutonomousSystemProfileConfigDataSource() string {
	fmt.Println("=== STEP  testing autonomous_system_profile Data Source with required arguments only")
	resource := fmt.Sprintln(`
	data "aci_autonomous_system_profile" "test" {}
	`)
	return resource
}

func CreateAccAutonomousSystemProfileDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing autonomous_system_profile Data Source with random attribute")
	resource := fmt.Sprintf(`

	data "aci_autonomous_system_profile" "test" {
		%s = "%s"
	}
	`, key, value)
	return resource
}
