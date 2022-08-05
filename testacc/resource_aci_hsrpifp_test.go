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

func TestAccAciL3outHSRPInterfaceProfile_Basic(t *testing.T) {
	var l3out_hsrp_interface_profile_default models.L3outHSRPInterfaceProfile
	var l3out_hsrp_interface_profile_updated models.L3outHSRPInterfaceProfile
	resourceName := "aci_l3out_hsrp_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHSRPInterfaceProfileWithoutRequired(rName, rName, rName, rName, "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "version", "v1"),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", ""),
				),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigWithOptionalValues(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_hsrp_interface_profile"),
					resource.TestCheckResourceAttr(resourceName, "version", "v2"),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", ""),
					testAccCheckAciL3outHSRPInterfaceProfileIdEqual(&l3out_hsrp_interface_profile_default, &l3out_hsrp_interface_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					testAccCheckAciL3outHSRPInterfaceProfileIdNotEqual(&l3out_hsrp_interface_profile_default, &l3out_hsrp_interface_profile_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
			},
		},
	})
}

func TestAccAciL3outHSRPInterfaceProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class hsrpIfP (.)+`),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileUpdatedAttr(rName, rName, rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileUpdatedAttr(rName, rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileUpdatedAttr(rName, rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPInterfaceProfileUpdatedAttr(rName, rName, rName, rName, "version", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outHSRPInterfaceProfileUpdatedAttr(rName, rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
			},
		},
	})
}

func TestAccAciL3outHSRPInterfaceProfile_RelationParameters(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	relName1 := makeTestVariable(acctest.RandString(5))
	relName2 := makeTestVariable(acctest.RandString(5))
	var l3out_hsrp_interface_profile_default models.L3outHSRPInterfaceProfile
	var l3out_hsrp_interface_profile_updated models.L3outHSRPInterfaceProfile
	resourceName := "aci_l3out_hsrp_interface_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", ""),
				),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigRelParameter(rName, relName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", fmt.Sprintf("uni/tn-%s/hsrpIfPol-%s", rName, relName1)),
					testAccCheckAciL3outHSRPInterfaceProfileIdEqual(&l3out_hsrp_interface_profile_default, &l3out_hsrp_interface_profile_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigRelParameter(rName, relName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", fmt.Sprintf("uni/tn-%s/hsrpIfPol-%s", rName, relName2)),
					testAccCheckAciL3outHSRPInterfaceProfileIdEqual(&l3out_hsrp_interface_profile_default, &l3out_hsrp_interface_profile_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfig(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists(resourceName, &l3out_hsrp_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_if_pol", ""),
				),
			},
		},
	})
}

func TestAccAciL3outHSRPInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPInterfaceProfileConfigs(rName),
			},
		},
	})
}

func testAccCheckAciL3outHSRPInterfaceProfileExists(name string, l3out_hsrp_interface_profile *models.L3outHSRPInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out HSRP Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out HSRP Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_interface_profileFound := models.L3outHSRPInterfaceProfileFromContainer(cont)
		if l3out_hsrp_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out HSRP Interface Profile %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_interface_profile = *l3out_hsrp_interface_profileFound
		return nil
	}
}

func testAccCheckAciL3outHSRPInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_hsrp_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_hsrp_interface_profile := models.L3outHSRPInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out HSRP Interface Profile %s Still exists", l3out_hsrp_interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outHSRPInterfaceProfileIdEqual(m1, m2 *models.L3outHSRPInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outHSRPInterfaceProfileIdNotEqual(m1, m2 *models.L3outHSRPInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_profile DNs are equal")
		}
		return nil
	}
}

func CreateL3outHSRPInterfaceProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_profile creation without ", attrName)
	rBlock := `
	
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
	
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	resource "aci_l3out_hsrp_interface_profile" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
}

func CreateAccL3outHSRPInterfaceProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile creation with updated resource name", rName)
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileConfigs(rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile multiple creation with required arguments")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_logical_interface_profile" "test1" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_logical_interface_profile" "test2" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_logical_interface_profile" "test3" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test1" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test1.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test2" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test2.id
	}

	resource "aci_l3out_hsrp_interface_profile" "test3" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test3.id
	}
	`, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccL3outHSRPInterfaceProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_profile creation with required arguments only")
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileConfigRelParameter(rName, relName string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_profile creation with resource name %s and relation resource name %s\n", rName, relName)
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

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	  }

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		relation_hsrp_rs_if_pol = aci_hsrp_interface_policy.test.id
	}
	`, rName, rName, rName, rName, relName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_hsrp_interface_profile creation with invalid logical_interface_profile_dn")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_tenant.test.id	
	}
	`, rName)
	return resource
}

func CreateAccL3outHSRPInterfaceProfileConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_profile creation with optional parameters")
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_interface_profile"
		version = "v2"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName)

	return resource
}

func CreateAccL3outHSRPInterfaceProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_profile updation without required parameter")
	resource := fmt.Sprintln(`
	resource "aci_l3out_hsrp_interface_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_interface_profile"
		version = "v2"
	}
	`)

	return resource
}

func CreateAccL3outHSRPInterfaceProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_profile attribute: %s=%s \n", attribute, value)
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attribute, value)
	return resource
}
