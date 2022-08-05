package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVLANPool_Basic(t *testing.T) {
	var vlan_pool_default models.VLANPool
	var vlan_pool_updated models.VLANPool
	resourceName := "aci_vlan_pool.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"
	allocModeUpdated := "static"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVLANPoolWithoutRequired(rName, allocMode, "alloc_mode"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVLANPoolWithoutRequired(rName, allocMode, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVLANPoolConfig(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists(resourceName, &vlan_pool_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocMode),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVLANPoolConfigWithOptionalValues(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists(resourceName, &vlan_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocMode),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vlan_pool"),
					testAccCheckAciVLANPoolIdEqual(&vlan_pool_default, &vlan_pool_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVLANPoolConfigUpdatedName(acctest.RandString(65), allocMode),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config: CreateAccVLANPoolConfig(rName, allocMode),
			},
			{
				Config:      CreateAccVLANPoolConfig(rName, acctest.RandString(5)),
				ExpectError: regexp.MustCompile(`expected (.)* to be one of (.)*`),
			},
			{
				Config:      CreateAccVLANPoolRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVLANPoolConfigWithRequiredParams(rNameUpdated, allocMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists(resourceName, &vlan_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVLANPoolIdNotEqual(&vlan_pool_default, &vlan_pool_updated),
				),
			},
			{
				Config: CreateAccVLANPoolConfig(rName, allocMode),
			},
			{
				Config: CreateAccVLANPoolConfigWithRequiredParams(rName, allocModeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists(resourceName, &vlan_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocModeUpdated),
					testAccCheckAciVLANPoolIdNotEqual(&vlan_pool_default, &vlan_pool_updated),
				),
			},
		},
	})
}

func TestAccAciVLANPool_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVLANPoolConfig(rName, allocMode),
			},
			{
				Config:      CreateAccVLANPoolUpdatedAttr(rName, allocMode, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVLANPoolUpdatedAttr(rName, allocMode, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVLANPoolUpdatedAttr(rName, allocMode, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVLANPoolUpdatedAttr(rName, allocMode, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVLANPoolConfig(rName, allocMode),
			},
		},
	})
}

func TestAccAciVLANPool_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVLANPoolConfigMultiple(rName, allocMode),
			},
		},
	})
}

func testAccCheckAciVLANPoolExists(name string, vlan_pool *models.VLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VLAN Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VLAN Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vlan_poolFound := models.VLANPoolFromContainer(cont)
		if vlan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VLAN Pool %s not found", rs.Primary.ID)
		}
		*vlan_pool = *vlan_poolFound
		return nil
	}
}

func testAccCheckAciVLANPoolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vlan_pool destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vlan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vlan_pool := models.VLANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VLAN Pool %s Still exists", vlan_pool.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVLANPoolIdEqual(m1, m2 *models.VLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vlan_pool DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVLANPoolIdNotEqual(m1, m2 *models.VLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vlan_pool DNs are equal")
		}
		return nil
	}
}

func CreateVLANPoolWithoutRequired(rName, allocMode, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vlan_pool creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vlan_pool" "test" {
	
	#	name  = "%s"
		alloc_mode  = "%s"
	}
		`
	case "alloc_mode":
		rBlock += `
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
	#	alloc_mode  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, allocMode)
}

func CreateAccVLANPoolConfigWithRequiredParams(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vlan_pool creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}
func CreateAccVLANPoolConfigUpdatedName(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vlan_pool creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVLANPoolConfig(rName, allocMode string) string {
	fmt.Printf("=== STEP  testing vlan_pool creation with required arguments name %s and alloc_mode %s\n", rName, allocMode)
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVLANPoolConfigMultiple(rName, allocMode string) string {
	fmt.Println("=== STEP  testing multiple vlan_pool creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s_${count.index}"
		alloc_mode  = "dynamic"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccVLANPoolConfigWithOptionalValues(rName, allocMode string) string {
	fmt.Println("=== STEP  Basic: testing vlan_pool creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vlan_pool"
		
	}
	`, rName, allocMode)

	return resource
}

func CreateAccVLANPoolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vlan_pool updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vlan_pool" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vlan_pool"
		
	}
	`)

	return resource
}

func CreateAccVLANPoolUpdatedAttr(rName, allocMode, attribute, value string) string {
	fmt.Printf("=== STEP  testing vlan_pool attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vlan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		%s = "%s"
	}
	`, rName, allocMode, attribute, value)
	return resource
}
