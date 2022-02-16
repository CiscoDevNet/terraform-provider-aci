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

func TestAccAciVXLANPool_Basic(t *testing.T) {
	var vxlan_pool_default models.VXLANPool
	var vxlan_pool_updated models.VXLANPool
	resourceName := "aci_vxlan_pool.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVXLANPoolWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVXLANPoolConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists(resourceName, &vxlan_pool_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVXLANPoolConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists(resourceName, &vxlan_pool_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vxlan_pool"),

					testAccCheckAciVXLANPoolIdEqual(&vxlan_pool_default, &vxlan_pool_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVXLANPoolConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVXLANPoolRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccVXLANPoolConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists(resourceName, &vxlan_pool_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVXLANPoolIdNotEqual(&vxlan_pool_default, &vxlan_pool_updated),
				),
			},
		},
	})
}

func TestAccAciVXLANPool_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVXLANPoolConfig(rName),
			},

			{
				Config:      CreateAccVXLANPoolUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVXLANPoolUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVXLANPoolUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVXLANPoolUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVXLANPoolConfig(rName),
			},
		},
	})
}

func TestAccAciVXLANPool_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVXLANPoolConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciVXLANPoolExists(name string, vxlan_pool *models.VXLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VXLAN Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VXLAN Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vxlan_poolFound := models.VXLANPoolFromContainer(cont)
		if vxlan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VXLAN Pool %s not found", rs.Primary.ID)
		}
		*vxlan_pool = *vxlan_poolFound
		return nil
	}
}

func testAccCheckAciVXLANPoolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vxlan_pool destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vxlan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vxlan_pool := models.VXLANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VXLAN Pool %s Still exists", vxlan_pool.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVXLANPoolIdEqual(m1, m2 *models.VXLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vxlan_pool DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVXLANPoolIdNotEqual(m1, m2 *models.VXLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vxlan_pool DNs are equal")
		}
		return nil
	}
}

func CreateVXLANPoolWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vxlan_pool creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vxlan_pool" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccVXLANPoolConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing vxlan_pool creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccVXLANPoolConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing vxlan_pool creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVXLANPoolConfig(rName string) string {
	fmt.Println("=== STEP  testing vxlan_pool creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVXLANPoolConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple vxlan_pool creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccVXLANPoolConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing vxlan_pool creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vxlan_pool"
		
	}
	`, rName)

	return resource
}

func CreateAccVXLANPoolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vxlan_pool updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vxlan_pool" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vxlan_pool"
		
	}
	`)

	return resource
}

func CreateAccVXLANPoolUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vxlan_pool attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccVXLANPoolUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vxlan_pool attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vxlan_pool" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
