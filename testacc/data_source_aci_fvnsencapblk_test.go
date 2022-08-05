package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRangesDataSource_Basic(t *testing.T) {
	resourceName := "aci_ranges.test"
	dataSourceName := "data.aci_ranges.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	from := strconv.Itoa(acctest.RandIntRange(1, 10))
	to := strconv.Itoa(acctest.RandIntRange(11, 20))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRangesDSWithoutRequired(rName, from, to, "vlan_pool_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRangesDSWithoutRequired(rName, from, to, "from"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRangesDSWithoutRequired(rName, from, to, "to"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRangesConfigDataSource(rName, from, to),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "from", resourceName, "from"),
					resource.TestCheckResourceAttrPair(dataSourceName, "to", resourceName, "to"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "alloc_mode", resourceName, "alloc_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vlan_pool_dn", resourceName, "vlan_pool_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "role", resourceName, "role"),
				),
			},
			{
				Config:      CreateAccRangesDataSourceUpdate(rName, from, to, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccRangesDSWithInvalidName(rName, from, to),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},
			{
				Config: CreateAccRangesDataSourceUpdatedResource(rName, from, to, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccRangesConfigDataSource(rName, from, to string) string {
	fmt.Println("=== STEP  testing ranges Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}	
	
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}

	data "aci_ranges" "test" {
		vlan_pool_dn = aci_ranges.test.vlan_pool_dn
		from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
	`, rName, from, to)
	return resource
}

func CreateRangesDSWithoutRequired(rName, from, to, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ranges Data Source without ", attrName)
	rBlock := `

	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}

	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}
	`
	switch attrName {
	case "from":
		rBlock += `
	data "aci_ranges" "test" {
		vlan_pool_dn = aci_ranges.test.vlan_pool_dn
	#	from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
		`
	case "to":
		rBlock += `
	data "aci_ranges" "test" {
		vlan_pool_dn = aci_ranges.test.vlan_pool_dn
		from  = aci_ranges.test.from
	#	to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
	`
	case "vlan_pool_dn":
		rBlock += `
	data "aci_ranges" "test" {
	#	vlan_pool_dn = aci_ranges.test.vlan_pool_dn
		from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, rName, from, to)
}

func CreateAccRangesDSWithInvalidName(rName, from, to string) string {
	fmt.Println("=== STEP  testing ranges Data Source with Invalid Name")
	resource := fmt.Sprintf(`

	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}

	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}

	data "aci_ranges" "test" {
		vlan_pool_dn = "${aci_ranges.test.vlan_pool_dn}_invalid"
		from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
	`, rName, from, to)
	return resource
}

func CreateAccRangesDataSourceUpdate(rName, from, to, key, value string) string {
	fmt.Println("=== STEP  testing ranges Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}

	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}

	data "aci_ranges" "test" {
		vlan_pool_dn = aci_ranges.test.vlan_pool_dn
		from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		%s = "%s"
		depends_on = [ aci_ranges.test ]
	}
	`, rName, from, to, key, value)
	return resource
}

func CreateAccRangesDataSourceUpdatedResource(rName, from, to, key, value string) string {
	fmt.Println("=== STEP  testing ranges Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}

	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
		%s = "%s"
	}

	data "aci_ranges" "test" {
		vlan_pool_dn = aci_ranges.test.vlan_pool_dn
		from  = aci_ranges.test.from
		to  = aci_ranges.test.to
		depends_on = [ aci_ranges.test ]
	}
	`, rName, from, to, key, value)
	return resource
}
