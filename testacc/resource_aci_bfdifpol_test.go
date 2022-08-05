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

func TestAccAciBfdInterfacePolicy_Basic(t *testing.T) {
	var bfd_interface_policy_default models.BFDInterfacePolicy
	var bfd_interface_policy_updated models.BFDInterfacePolicy
	resourceName := "aci_bfd_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBfdInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBfdInterfacePolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBfdInterfacePolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", ""),
					resource.TestCheckResourceAttr(resourceName, "detect_mult", "3"),
					resource.TestCheckResourceAttr(resourceName, "echo_admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "echo_rx_intvl", "50"),
					resource.TestCheckResourceAttr(resourceName, "min_rx_intvl", "50"),
					resource.TestCheckResourceAttr(resourceName, "min_tx_intvl", "50"),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bfd_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "opt-subif"),
					resource.TestCheckResourceAttr(resourceName, "detect_mult", "1"),
					resource.TestCheckResourceAttr(resourceName, "echo_admin_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "echo_rx_intvl", "999"),
					resource.TestCheckResourceAttr(resourceName, "min_rx_intvl", "999"),
					resource.TestCheckResourceAttr(resourceName, "min_tx_intvl", "999"),
					testAccCheckAciBfdInterfacePolicyIdEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBfdInterfacePolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBfdInterfacePolicyIdNotEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBfdInterfacePolicyIdNotEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
		},
	})
}

//detect_mult 50,
//echo_rx_intvl, min_rx_intvl, min_tx_intvl middle value
func TestAccAciBfdInterfacePolicy_Update(t *testing.T) {
	var bfd_interface_policy_default models.BFDInterfacePolicy
	var bfd_interface_policy_updated models.BFDInterfacePolicy
	resourceName := "aci_bfd_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBfdInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_default),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "detect_mult", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "detect_mult", "50"),
					testAccCheckAciBfdInterfacePolicyIdEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "echo_rx_intvl", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "echo_rx_intvl", "150"),
					testAccCheckAciBfdInterfacePolicyIdEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_rx_intvl", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "min_rx_intvl", "150"),
					testAccCheckAciBfdInterfacePolicyIdEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_tx_intvl", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdInterfacePolicyExists(resourceName, &bfd_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "min_tx_intvl", "150"),
					testAccCheckAciBfdInterfacePolicyIdEqual(&bfd_interface_policy_default, &bfd_interface_policy_updated),
				),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciBfdInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBfdInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccBfdInterfacePolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "detect_mult", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "detect_mult", "51"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "echo_admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "echo_rx_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "echo_rx_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "echo_rx_intvl", "1000"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_rx_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_rx_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_rx_intvl", "1000"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_tx_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_tx_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, "min_tx_intvl", "1000"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBfdInterfacePolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciBfdInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBfdInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBfdInterfacePolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciBfdInterfacePolicyExists(name string, bfd_interface_policy *models.BFDInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bfd Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bfd Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bfd_interface_policyFound := models.BFDInterfacePolicyFromContainer(cont)
		if bfd_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bfd Interface Policy %s not found", rs.Primary.ID)
		}
		*bfd_interface_policy = *bfd_interface_policyFound
		return nil
	}
}

func testAccCheckAciBfdInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bfd_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bfd_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			bfd_interface_policy := models.BFDInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bfd Interface Policy %s Still exists", bfd_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBfdInterfacePolicyIdEqual(m1, m2 *models.BFDInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bfd_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBfdInterfacePolicyIdNotEqual(m1, m2 *models.BFDInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bfd_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateBfdInterfacePolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bfd_interface_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bfd_interface_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBfdInterfacePolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy creation with updated required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bfd_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple bfd_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bfd_interface_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_aaa_domain" "test"{
		name = "%s"
	}
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_aaa_domain.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccBfdInterfacePolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bfd_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bfd_interface_policy"
		admin_st = "disabled"
		ctrl = "opt-subif"
		detect_mult = "1"
		echo_admin_st = "disabled"
		echo_rx_intvl = "999"
		min_rx_intvl = "999"
		min_tx_intvl = "999"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBfdInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bfd_interface_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_bfd_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bfd_interface_policy"
		admin_st = "disabled"
		ctrl = ["opt-subif"]
		detect_mult = "2"
		echo_admin_st = "disabled"
		echo_rx_intvl = "51"
		min_rx_intvl = "51"
		min_tx_intvl = "51"
	}
	`)

	return resource
}

func CreateAccBfdInterfacePolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bfd_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
