package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRanges_Basic(t *testing.T) {
	var ranges_default models.Ranges
	var ranges_updated models.Ranges
	resourceName := "aci_ranges.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	from := strconv.Itoa(acctest.RandIntRange(1, 5))
	fromUpdated := strconv.Itoa(acctest.RandIntRange(5, 10))
	to := strconv.Itoa(acctest.RandIntRange(11, 15))
	toUpdated := strconv.Itoa(acctest.RandIntRange(15, 20))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateRangesWithoutRequired(rName, from, to, "vlan_pool_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRangesWithoutRequired(rName, from, to, "from"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRangesWithoutRequired(rName, from, to, "to"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRangesConfig(rName, from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_default),
					resource.TestCheckResourceAttr(resourceName, "vlan_pool_dn", fmt.Sprintf("uni/infra/vlanns-[%s]-dynamic", rName)),
					resource.TestCheckResourceAttr(resourceName, "from", fmt.Sprintf("vlan-%s", from)),
					resource.TestCheckResourceAttr(resourceName, "to", fmt.Sprintf("vlan-%s", to)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "role", "external"),
				),
			},
			{
				Config: CreateAccRangesConfigWithOptionalValues(rName, from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "vlan_pool_dn", fmt.Sprintf("uni/infra/vlanns-[%s]-dynamic", rName)),
					resource.TestCheckResourceAttr(resourceName, "from", fmt.Sprintf("vlan-%s", from)),
					resource.TestCheckResourceAttr(resourceName, "to", fmt.Sprintf("vlan-%s", to)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ranges"),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", "dynamic"),
					resource.TestCheckResourceAttr(resourceName, "role", "internal"),
					testAccCheckAciRangesIdEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRangesRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRangesConfigWithRequiredParams(rNameUpdated, from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "vlan_pool_dn", fmt.Sprintf("uni/infra/vlanns-[%s]-dynamic", rNameUpdated)),
					testAccCheckAciRangesIdNotEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				Config: CreateAccRangesConfig(rName, from, to),
			},
			{
				Config: CreateAccRangesConfigWithRequiredParams(rName, fromUpdated, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "from", fmt.Sprintf("vlan-%s", fromUpdated)),
					testAccCheckAciRangesIdNotEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				Config: CreateAccRangesConfig(rName, from, to),
			},
			{
				Config: CreateAccRangesConfigWithRequiredParams(rName, from, toUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "to", fmt.Sprintf("vlan-%s", toUpdated)),
					testAccCheckAciRangesIdNotEqual(&ranges_default, &ranges_updated),
				),
			},
		},
	})
}

func TestAccAciRanges_Update(t *testing.T) {
	var ranges_default models.Ranges
	var ranges_updated models.Ranges
	resourceName := "aci_ranges.test"
	rName := makeTestVariable(acctest.RandString(5))
	from := strconv.Itoa(acctest.RandIntRange(1, 10))
	to := strconv.Itoa(acctest.RandIntRange(11, 20))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRangesConfig(rName, from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_default),
				),
			},
			{
				Config: CreateAccRangesConfig(rName, "4095", "4095"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "from", "vlan-4095"),
					resource.TestCheckResourceAttr(resourceName, "to", "vlan-4095"),
					testAccCheckAciRangesIdNotEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				Config: CreateAccRangesConfig(rName, from, "4095"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "to", "vlan-4095"),
					testAccCheckAciRangesIdNotEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				Config: CreateAccRangesUpdatedAttr(rName, from, to, "alloc_mode", "static"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists(resourceName, &ranges_updated),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", "static"),
					testAccCheckAciRangesIdEqual(&ranges_default, &ranges_updated),
				),
			},
			{
				Config: CreateAccRangesConfig(rName, from, to),
			},
		},
	})
}

func TestAccAciRanges_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	from := strconv.Itoa(acctest.RandIntRange(1, 10))
	to := strconv.Itoa(acctest.RandIntRange(11, 20))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRangesConfig(rName, from, to),
			},
			{
				Config:      CreateAccRangesWithInvalidVlanPoolDn(rName, from, to),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRangesConfig(rName, randomValue, to),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRangesConfig(rName, from, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRangesConfig(rName, "0", to),
				ExpectError: regexp.MustCompile(`Invalid encapsulation type. Only VLAN encapsulation type is allowed`),
			},
			{
				Config:      CreateAccRangesConfig(rName, to, from),
				ExpectError: regexp.MustCompile(`Range (.)* is invalid. From value cannot be larger than To value`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, "alloc_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, "role", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccRangesUpdatedAttr(rName, from, to, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRangesConfig(rName, from, to),
			},
		},
	})
}

func TestAccAciRanges_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRangesConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciRangesExists(name string, ranges *models.Ranges) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Ranges %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ranges dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		rangesFound := models.RangesFromContainer(cont)
		if rangesFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Ranges %s not found", rs.Primary.ID)
		}
		*ranges = *rangesFound
		return nil
	}
}

func testAccCheckAciRangesDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ranges destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ranges" {
			cont, err := client.Get(rs.Primary.ID)
			ranges := models.RangesFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Ranges %s Still exists", ranges.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRangesIdEqual(m1, m2 *models.Ranges) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ranges DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRangesIdNotEqual(m1, m2 *models.Ranges) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ranges DNs are equal")
		}
		return nil
	}
}

func CreateRangesWithoutRequired(rName, from, to, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ranges creation without ", attrName)
	rBlock := `
	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}	
	`
	switch attrName {
	case "vlan_pool_dn":
		rBlock += `
	resource "aci_ranges" "test" {
	#    vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}
	`
	case "from":
		rBlock += `
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
	#	from  = "vlan-%s"
		to  = "vlan-%s"
	}
	`
	case "to":
		rBlock += `
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
	#	to  = "vlan-%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName, from, to)
}

func CreateAccRangesConfigWithRequiredParams(rName, from, to string) string {
	fmt.Println("=== STEP  testing ranges creation with updated Required arguments")
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
	`, rName, from, to)
	return resource
}

func CreateAccRangesConfig(rName, from, to string) string {
	fmt.Println("=== STEP  testing ranges creation with required arguments only")
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
	`, rName, from, to)
	return resource
}

func CreateAccRangesWithInvalidVlanPoolDn(rName, from, to string) string {
	fmt.Println("=== STEP  testing ranges creation with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name = "%s"
	}	
	
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_tenant.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
	}
	`, rName, from, to)
	return resource
}

func CreateAccRangesConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple ranges creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}	
	
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-1"
		to  = "vlan-2"
	}
	
	resource "aci_ranges" "test1" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-3"
		to  = "vlan-4"
	}

	resource "aci_ranges" "test2" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-5"
		to  = "vlan-6"
	}

	resource "aci_ranges" "test3" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-7"
		to  = "vlan-8"
	}
	`, rName)
	return resource
}

func CreateAccRangesConfigWithOptionalValues(rName, from, to string) string {
	fmt.Println("=== STEP  Basic: testing ranges creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "dynamic"
	}	
	
	resource "aci_ranges" "test" {
		vlan_pool_dn = aci_vlan_pool.test.id
		from  = "vlan-%s"
		to  = "vlan-%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ranges"
		alloc_mode = "dynamic"
		role = "internal"
	}
	`, rName, from, to)

	return resource
}

func CreateAccRangesRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ranges updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ranges" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ranges"
		alloc_mode = "dynamic"
		role = "internal"		
	}
	`)
	return resource
}

func CreateAccRangesUpdatedAttr(rName, from, to, attribute, value string) string {
	fmt.Printf("=== STEP  testing ranges attribute: %s = %s \n", attribute, value)
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
	`, rName, from, to, attribute, value)
	return resource
}
