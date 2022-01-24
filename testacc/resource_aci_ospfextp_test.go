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

func TestAccAciL3outOspfExternalPolicy_Basic(t *testing.T) {
	var l3out_ospf_external_policy_default models.L3outOspfExternalPolicy
	var l3out_ospf_external_policy_updated models.L3outOspfExternalPolicy
	resourceName := "aci_l3out_ospf_external_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outOspfExternalPolicyWithoutRequired(rName, rName, "l3_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_default),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "area_cost", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "redistribute"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.1", "summary"),
					resource.TestCheckResourceAttr(resourceName, "area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "area_type", "nssa"),
					resource.TestCheckResourceAttr(resourceName, "multipod_internal", "no"),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_ospf_external_policy"),
					resource.TestCheckResourceAttr(resourceName, "area_cost", "0"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "suppress-fa"),
					resource.TestCheckResourceAttr(resourceName, "area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "area_type", "stub"),
					resource.TestCheckResourceAttr(resourceName, "multipod_internal", "no"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rNameUpdated, rNameUpdated)),
					testAccCheckAciL3outOspfExternalPolicyIdNotEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciL3outOspfExternalPolicy_Update(t *testing.T) {
	var l3out_ospf_external_policy_default models.L3outOspfExternalPolicy
	var l3out_ospf_external_policy_updated models.L3outOspfExternalPolicy
	var l3out_ospf_external_policy_default_infra models.L3outOspfExternalPolicy
	var l3out_ospf_external_policy_updated_infra models.L3outOspfExternalPolicy
	resourceName := "aci_l3out_ospf_external_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outOspfExternalPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_default),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"redistribute"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "redistribute"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"summary"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "summary"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"redistribute", "summary", "suppress-fa"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "redistribute"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.1", "summary"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.2", "suppress-fa"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"summary", "suppress-fa"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "summary"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.1", "suppress-fa"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"unspecified"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "unspecified"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"suppress-fa", "summary", "redistribute"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "suppress-fa"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.1", "summary"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.2", "redistribute"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"suppress-fa"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "area_ctrl.0", "suppress-fa"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_type", "regular"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_type", "regular"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_id", "0.0.1.1"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_id", "0.1.1.1"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "1.1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_id", "1.1.1.1"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_id", "backbone"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_cost", "16777215"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "area_cost", "16777215"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default, &l3out_ospf_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfigInfra(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_default_infra),
				),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyUpdatedAttrInfra(rName, "multipod_internal", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists(resourceName, &l3out_ospf_external_policy_updated_infra),
					resource.TestCheckResourceAttr(resourceName, "multipod_internal", "yes"),
					testAccCheckAciL3outOspfExternalPolicyIdEqual(&l3out_ospf_external_policy_default_infra, &l3out_ospf_external_policy_default_infra),
				),
			},
		},
	})
}

func TestAccAciL3outOspfExternalPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outOspfExternalPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_cost", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_cost", "16777216"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"redistribute", "redistribute"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttrList(rName, rName, "area_ctrl", StringListtoString([]string{"redistribute", "unspecified"})),
				ExpectError: regexp.MustCompile(`should't be used along with other values`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "256"),
				ExpectError: regexp.MustCompile(`Invalid value for area_id`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "256.1"),
				ExpectError: regexp.MustCompile(`Invalid value for area_id`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "256.1.1"),
				ExpectError: regexp.MustCompile(`Invalid value for area_id`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "256.1.1.1"),
				ExpectError: regexp.MustCompile(`Invalid value for area_id`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "1.1.1.1.1"),
				ExpectError: regexp.MustCompile(`Invalid value for area_id`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_id", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "area_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "multipod_internal", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, "multipod_internal", "yes"),
				ExpectError: regexp.MustCompile(`Invalid Configuration OSPF Multipod Internal can be set only under infra tenant.`),
			},
			{
				Config:      CreateAccL3outOspfExternalPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outOspfExternalPolicyConfig(rName, rName),
			},
		},
	})
}

func testAccCheckAciL3outOspfExternalPolicyExists(name string, l3out_ospf_external_policy *models.L3outOspfExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Ospf External Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Ospf External Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_ospf_external_policyFound := models.L3outOspfExternalPolicyFromContainer(cont)
		if l3out_ospf_external_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Ospf External Policy %s not found", rs.Primary.ID)
		}
		*l3out_ospf_external_policy = *l3out_ospf_external_policyFound
		return nil
	}
}

func testAccCheckAciL3outOspfExternalPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_ospf_external_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_ospf_external_policy := models.L3outOspfExternalPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Ospf External Policy %s Still exists", l3out_ospf_external_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outOspfExternalPolicyIdEqual(m1, m2 *models.L3outOspfExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_ospf_external_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outOspfExternalPolicyIdNotEqual(m1, m2 *models.L3outOspfExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_ospf_external_policy DNs are equal")
		}
		return nil
	}
}

func CreateL3outOspfExternalPolicyWithoutRequired(fvTenantName, l3extOutName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_external_policy creation without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	`
	switch attrName {
	case "l3_outside_dn":
		rBlock += `
	resource "aci_l3out_ospf_external_policy" "test" {
	#	l3_outside_dn  = aci_l3_outside.test.id

	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName)
}

func CreateAccL3outOspfExternalPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy creation with updated l3_outside name")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`, rName, rName)
	return resource
}

func CreateAccL3outOspfExternalPolicyConfigInfra(rName string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy creation with infra tenant")
	resource := fmt.Sprintf(`

	data "aci_tenant" "test" {
		name 		= "infra"
	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = data.aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`, rName)
	return resource
}

func CreateAccL3outOspfExternalPolicyConfig(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_ospf_external_policy creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outOspfExternalPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_ospf_external_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccL3outOspfExternalPolicyConfigWithOptionalValues(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_external_policy creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = "${aci_l3_outside.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_ospf_external_policy"
		area_cost = "0"
		area_ctrl = ["suppress-fa"]
		area_id = "1"
		area_type = "stub"
		multipod_internal = "no"
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outOspfExternalPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_ospf_external_policy update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3out_ospf_external_policy" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_l3out_ospf_external_policy"
		area_cost = "1"
		area_ctrl = ["suppress-fa"]
		area_id = "backbone"
		area_type = "regular"
		multipod_internal = "yes"
	}
	`)

	return resource
}

func CreateAccL3outOspfExternalPolicyUpdatedAttrInfra(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_ospf_external_policy attribute: %s = %s for infra tenant\n", attribute, value)
	resource := fmt.Sprintf(`

	data "aci_tenant" "test" {
		name 		= "infra"
	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = data.aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccL3outOspfExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_ospf_external_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, attribute, value)
	return resource
}

func CreateAccL3outOspfExternalPolicyUpdatedAttrList(fvTenantName, l3extOutName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_ospf_external_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l3out_ospf_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = %s
	}
	`, fvTenantName, l3extOutName, attribute, value)
	return resource
}
