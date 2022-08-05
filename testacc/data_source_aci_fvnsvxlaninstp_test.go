package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVXLANPoolDataSource_Basic(t *testing.T) {
	resourceName := "aci_vxlan_pool.test"
	dataSourceName := "data.aci_vxlan_pool.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVXLANPoolDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVXLANPoolConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVXLANPoolDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVXLANPoolDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVXLANPoolDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVXLANPoolConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing vxlan_pool Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}

	data "aci_vxlan_pool" "test" {
	
		name  = aci_vxlan_pool.test.name
		depends_on = [ aci_vxlan_pool.test ]
	}
	`, rName)
	return resource
}

func CreateVXLANPoolDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vxlan_pool Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_vxlan_pool" "test" {
	
	#	name  = aci_vxlan_pool.test.name
		depends_on = [ aci_vxlan_pool.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccVXLANPoolDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing vxlan_pool Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}

	data "aci_vxlan_pool" "test" {
	
		name  = "${aci_vxlan_pool.test.name}_invalid"
		depends_on = [ aci_vxlan_pool.test ]
	}
	`, rName)
	return resource
}

func CreateAccVXLANPoolDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing vxlan_pool Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}

	data "aci_vxlan_pool" "test" {
	
		name  = aci_vxlan_pool.test.name
		%s = "%s"
		depends_on = [ aci_vxlan_pool.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccVXLANPoolDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing vxlan_pool Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_vxlan_pool" "test" {
	
		name  = aci_vxlan_pool.test.name
		depends_on = [ aci_vxlan_pool.test ]
	}
	`, rName, key, value)
	return resource
}
