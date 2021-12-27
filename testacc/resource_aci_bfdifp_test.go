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

func TestAccAciL3outBfdInterfaceProfile_Basic(t *testing.T) {
	var l3out_bfd_interface_profile_default models.BFDInterfaceProfile
	var l3out_bfd_interface_profile_updated models.BFDInterfaceProfile
	resourceName := "aci_l3out_bfd_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	pNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBfdInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBfdInterfaceProfileWithoutRequired(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBfdInterfaceProfileExists(resourceName, &l3out_bfd_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation","orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description",""),
					resource.TestCheckResourceAttr(resourceName, "name_alias",""),
					resource.TestCheckResourceAttr(resourceName, "key_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "interface_profile_type", "none"),
					
				),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfigWithOptionalValues(rName), 
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBfdInterfaceProfileExists(resourceName, &l3out_bfd_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_bfd_interface_profile"),
					resource.TestCheckResourceAttr(resourceName, "key", "1234"),
					resource.TestCheckResourceAttr(resourceName, "key_id", "255"),
					resource.TestCheckResourceAttr(resourceName, "interface_profile_type", "sha1"),
					testAccCheckAciL3outBfdInterfaceProfileIdEqual(&l3out_bfd_interface_profile_default,&l3out_bfd_interface_profile_updated),
				),
			},  
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfigWithRequiredParams(rName, pNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBfdInterfaceProfileExists(resourceName, &l3out_bfd_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, pNameUpdated)),
					testAccCheckAciL3outBfdInterfaceProfileIdNotEqual(&l3out_bfd_interface_profile_default, &l3out_bfd_interface_profile_updated),
				),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileUpdateWithoutRequiredParams(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciL3outBfdInterfaceProfile_Update(t *testing.T) {
	var l3out_bfd_interface_profile_default models.BFDInterfaceProfile
	var l3out_bfd_interface_profile_updated models.BFDInterfaceProfile
	resourceName := "aci_l3out_bfd_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBfdInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBfdInterfaceProfileExists(resourceName, &l3out_bfd_interface_profile_default),
				),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName, "key_id", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBfdInterfaceProfileExists(resourceName, &l3out_bfd_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "key_id", "100"),
					testAccCheckAciL3outBfdInterfaceProfileIdEqual(&l3out_bfd_interface_profile_default,&l3out_bfd_interface_profile_updated),
				),
			},
		},
	})
}

func TestAccAciL3outBfdInterfaceProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBfdInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(` unknown property value (.)+, name dn, class bfdIfP (.)+`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "key", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "key_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name keyId, class bfdIfP (.)+`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "key_id", "0"),
				ExpectError: regexp.MustCompile(`Property keyId of (.)+ is out of range`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "key_id", "-10"),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name keyId, class bfdIfP (.)+`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "key_id", "256"),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name keyId, class bfdIfP (.)+`),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , "interface_profile_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+to be one of(.)+, got(.)+`),
			},
			
			{
				Config:      CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName , randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciL3outBfdInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBfdInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBfdInterfaceProfileMultiple(rName),
			},
		},
	})
}

func CreateAccL3outBfdInterfaceProfileMultiple(rName  string) string {
	fmt.Println("=== STEP Testing Multiple l3out_bfd_interface_profile creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_logical_interface_profile" "test1" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_bfd_interface_profile" "test1" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test1.id
	}
	`, rName, rName, rName, rName, rName+"1")
	return resource
}

func testAccCheckAciL3outBfdInterfaceProfileExists(name string, l3out_bfd_interface_profile *models.BFDInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Bfd Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Bfd Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_bfd_interface_profileFound := models.BFDInterfaceProfileFromContainer(cont)
		if l3out_bfd_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Bfd Interface Profile %s not found", rs.Primary.ID)
		}
		*l3out_bfd_interface_profile = *l3out_bfd_interface_profileFound
		return nil
	}
}

func testAccCheckAciL3outBfdInterfaceProfileDestroy(s *terraform.State) error {	
	fmt.Println("=== STEP Testing l3out_bfd_interface_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		 if rs.Type == "aci_l3out_bfd_interface_profile" {
			cont,err := client.Get(rs.Primary.ID)
			l3out_bfd_interface_profile := models.BFDInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Bfd Interface Profile %s Still exists",l3out_bfd_interface_profile.DistinguishedName)
			}
		}else{
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outBfdInterfaceProfileIdEqual(m1, m2 *models.BFDInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_bfd_interface_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outBfdInterfaceProfileIdNotEqual(m1, m2 *models.BFDInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_bfd_interface_profile DNs are equal")
		}
		return nil
	}
}

func CreateL3outBfdInterfaceProfileWithoutRequired() string {
	fmt.Println("=== STEP  Basic: Testing l3out_bfd_interface_profile creation without required attribute")
	resource := fmt.Sprintln(`
	resource "aci_l3out_bfd_interface_profile" "test" {
		}
	`,)
	return resource
}

func CreateAccL3outBfdInterfaceProfileConfigWithRequiredParams(rName, pNameupdated string) string {
	fmt.Println("=== STEP Testing l3out_bfd_interface_profile creation with different logical_interface_profile")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
	`, rName, rName, rName, pNameupdated )
	return resource
}


func CreateAccL3outBfdInterfaceProfileConfig(rName  string) string {
	fmt.Println("=== STEP Testing l3out_bfd_interface_profile creation with required arguement only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outBfdInterfaceProfileWithInValidParentDn(rName  string) string {
	fmt.Println("=== STEP  Negative Case: Testing l3out_bfd_interface_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_node_profile.test.id	
	}
	`, rName, rName, rName, )
	return resource
}


func CreateAccL3outBfdInterfaceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: Testing l3out_bfd_interface_profile creation with optional parameters")
	resource := fmt.Sprintf(`
		resource "aci_tenant" "test" {
			name 		= "%s"
		}
		
		resource "aci_l3_outside" "test" {
			name 		= "%s"
			tenant_dn = aci_tenant.test.id
		}
		
		resource "aci_logical_node_profile" "test" {
			name 		= "%s"
			l3_outside_dn = aci_l3_outside.test.id
		}
		
		resource "aci_logical_interface_profile" "test" {
			name 		= "%s"
			logical_node_profile_dn = aci_logical_node_profile.test.id
		}
		resource "aci_l3out_bfd_interface_profile" "test" {
			logical_interface_profile_dn  = aci_logical_interface_profile.test.id
			description = "created while acceptance testing"
			annotation = "orchestrator:terraform_testacc"
			name_alias = "test_l3out_bfd_interface_profile"
			key = "1234"
			key_id = "255"
			interface_profile_type = "sha1"
		}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outBfdInterfaceProfileUpdateWithoutRequiredParams() string {
	fmt.Println("=== STEP  Basic: Testing l3out_bfd_interface_profile updation without required parameter")
	resource := fmt.Sprintln(`
		resource "aci_l3out_bfd_interface_profile" "test" {
			description = "created while acceptance testing1"
			annotation = "orchestrator:terraform_testacc1"
			name_alias = "test_l3out_bfd_interface_profile"
			key = "1234"
			key_id = "255"
			interface_profile_type = "sha1"
		}
	`)
	return resource
}

func CreateAccL3outBfdInterfaceProfileUpdatedAttr(rName, attribute,value string) string {
	fmt.Printf("=== STEP Testing l3out_bfd_interface_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
	}
	`, rName, rName, rName, rName ,attribute,value)
	return resource
}