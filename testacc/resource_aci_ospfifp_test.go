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

func TestAccAciL3outOspfInterfaceProfile_Basic(t *testing.T) {
	var l3out_ospf_interface_profile_default models.OSPFInterfaceProfile
	var l3out_ospf_interface_profile_updated models.OSPFInterfaceProfile
	resourceName := "aci_l3out_ospf_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	randomValue := acctest.RandString(5)
	randomValueUpdated := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outOspfInterfaceProfileWithoutRequired(rName, rName, rName, rName, "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outOspfInterfaceProfileWithoutRequired(rName, rName, rName, rName, "auth_key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfInterfaceProfileExists(resourceName, &l3out_ospf_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_key", randomValue),
					resource.TestCheckResourceAttr(resourceName, "auth_key_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "none"),
					resource.TestCheckResourceAttr(resourceName, "relation_ospf_rs_if_pol", ""),
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccL3outOspfInterfaceProfileConfigWithOptionalValues(rName, rName, rName, rName, randomValueUpdated), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfInterfaceProfileExists(resourceName, &l3out_ospf_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_ospf_interface_profile"),
					resource.TestCheckResourceAttr(resourceName, "auth_key", randomValueUpdated),
					resource.TestCheckResourceAttr(resourceName, "auth_key_id", "2"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "md5"),
					resource.TestCheckResourceAttr(resourceName, "relation_ospf_rs_if_pol", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_key"},
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfigWithRequiredParams(rNameUpdated, randomValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfInterfaceProfileExists(resourceName, &l3out_ospf_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					testAccCheckAciL3outOspfInterfaceProfileIdNotEqual(&l3out_ospf_interface_profile_default, &l3out_ospf_interface_profile_updated),
				),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileConfigUpdateWithoutRequiredParameters(rName, "description", "test_coverage"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
			},
		},
	})
}

func TestAccAciL3outOspfInterfaceProfile_Update(t *testing.T) {
	var l3out_ospf_interface_profile_default models.OSPFInterfaceProfile
	var l3out_ospf_interface_profile_updated models.OSPFInterfaceProfile
	resourceName := "aci_l3out_ospf_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfInterfaceProfileExists(resourceName, &l3out_ospf_interface_profile_default),
				),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "auth_type", "simple"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfInterfaceProfileExists(resourceName, &l3out_ospf_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "simple"),
					testAccCheckAciL3outOspfInterfaceProfileIdEqual(&l3out_ospf_interface_profile_default, &l3out_ospf_interface_profile_updated),
				),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
			},
		},
	})
}

func TestAccAciL3outOspfInterfaceProfile_Negative(t *testing.T) {

	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileWithInValidParentDn(rName, rName, rName, rName),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "auth_key_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, "auth_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfig(rName, rName, rName, rName, randomValue),
			},
		},
	})
}

func testAccCheckAciL3outOspfInterfaceProfileExists(name string, l3out_ospf_interface_profile *models.OSPFInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Ospf Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Ospf Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_ospf_interface_profileFound := models.OSPFInterfaceProfileFromContainer(cont)
		if l3out_ospf_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Ospf Interface Profile %s not found", rs.Primary.ID)
		}
		*l3out_ospf_interface_profile = *l3out_ospf_interface_profileFound
		return nil
	}
}

func testAccCheckAciL3outOspfInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_ospf_interface_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_ospf_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_ospf_interface_profile := models.OSPFInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Ospf Interface Profile %s Still exists", l3out_ospf_interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outOspfInterfaceProfileIdEqual(m1, m2 *models.OSPFInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_ospf_interface_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outOspfInterfaceProfileIdNotEqual(m1, m2 *models.OSPFInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_ospf_interface_profile DNs are equal")
		}
		return nil
	}
}

func CreateL3outOspfInterfaceProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_interface_profile creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
		resource "aci_l3out_ospf_interface_profile" "test" {
		#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
			auth_key = "random"
			description = "created while acceptance testing"
	}
	`
	case "auth_key":
		rBlock += `
		resource "aci_l3out_ospf_interface_profile" "test" {
				logical_interface_profile_dn  = aci_logical_interface_profile.test.id
				#auth_key = "random"
				description = "created while acceptance testing"
			}`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
}

func CreateAccL3outOspfInterfaceProfileConfigWithRequiredParams(rName, randomValue string) string {
	fmt.Println("=== STEP  testing l3out_ospf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}
	`, rName, rName, rName, rName, randomValue)
	return resource
}

func CreateAccL3outOspfInterfaceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_ospf_interface_profile creation with : %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, rName, value, attribute, value)
	return resource
}

func CreateAccL3outOspfInterfaceProfileConfigUpdateWithoutRequiredParameters(rName, attribute, value string) string {
	fmt.Println("=== STEP  testing l3out_ospf_interface_profile updation without required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		auth_key = "random"
		%s = "%s"
	}
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccL3outOspfInterfaceProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, randomValue string) string {
	fmt.Println("=== STEP  testing l3out_ospf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, randomValue)
	return resource
}

func CreateAccL3outOspfInterfaceProfileWithInValidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_ospf_interface_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}invalid"	
		auth_key = "random"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
	return resource
}

func CreateAccL3outOspfInterfaceProfileConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, randomValue string) string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_interface_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_ospf_interface_profile"
		auth_key = "%s"
		auth_key_id = "2"
		auth_type = "md5"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, randomValue)

	return resource
}
