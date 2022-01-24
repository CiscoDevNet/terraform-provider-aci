package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVLANPoolDataSource_Basic(t *testing.T) {
	resourceName := "aci_vlan_pool.test"
	dataSourceName := "data.aci_vlan_pool.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVLANPoolDSWithoutRequired(rName, allocMode, "alloc_mode"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateVLANPoolDSWithoutRequired(rName, allocMode, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVLANPoolConfigDataSource(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "alloc_mode", resourceName, "alloc_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVLANPoolDataSourceUpdate(rName, allocMode, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVLANPoolDSWithInvalidName(rName, allocMode),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVLANPoolDataSourceUpdatedResource(rName, allocMode, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVLANPoolDSWithInvalidName(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vlan_pool Data Source with Invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vlan_pool" "test" {
		name  = "${aci_vlan_pool.test.name}_invalid"
		alloc_mode  = aci_vlan_pool.test.alloc_mode
		depends_on = [ aci_vlan_pool.test ]
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVLANPoolConfigDataSource(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vlan_pool Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vlan_pool" "test" {
	
		name  = aci_vlan_pool.test.name
		alloc_mode  = aci_vlan_pool.test.alloc_mode
		depends_on = [ aci_vlan_pool.test ]
	}
	`, rName, allocMode)
	return resource
}

func CreateVLANPoolDSWithoutRequired(rName, allocMode, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vlan_pool Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_vlan_pool" "test" {
	
	#	name  = "%s"
		alloc_mode  = "%s"
		depends_on = [ aci_vlan_pool.test ]
	}
		`
	case "alloc_mode":
		rBlock += `
	data "aci_vlan_pool" "test" {
	
		name  = "%s"
	#	alloc_mode  = "%s"
		depends_on = [ aci_vlan_pool.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, allocMode)
}

func CreateAccVLANPoolDataSourceUpdate(rName, allocMode, key, value string) string {
	fmt.Println("=== STEP  testing vlan_pool Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vlan_pool" "test" {
	
		name  = aci_vlan_pool.test.name
		alloc_mode  = aci_vlan_pool.test.alloc_mode
		%s = "%s"
		depends_on = [ aci_vlan_pool.test ]
	}
	`, rName, allocMode, key, value)
	return resource
}

func CreateAccVLANPoolDataSourceUpdatedResource(rName, allocMode, key, value string) string {
	fmt.Println("=== STEP  testing vlan_pool Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		%s = "%s"
	}

	data "aci_vlan_pool" "test" {
	
		name  = aci_vlan_pool.test.name
		alloc_mode  = aci_vlan_pool.test.alloc_mode
		depends_on = [ aci_vlan_pool.test ]
	}
	`, rName, allocMode, key, value)
	return resource
}
