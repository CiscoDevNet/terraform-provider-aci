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

func TestAccAciOSPFInterfacePolicy_Basic(t *testing.T) {
	var ospf_interface_policy_default models.OSPFInterfacePolicy
	var ospf_interface_policy_updated models.OSPFInterfacePolicy
	resourceName := "aci_ospf_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOSPFInterfacePolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOSPFInterfacePolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "cost", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "40"),
					resource.TestCheckResourceAttr(resourceName, "hello_intvl", "10"),
					resource.TestCheckResourceAttr(resourceName, "nw_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "pfx_suppress", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "prio", "1"),
					resource.TestCheckResourceAttr(resourceName, "rexmit_intvl", "5"),
					resource.TestCheckResourceAttr(resourceName, "xmit_delay", "1"),
				),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ospf_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "cost", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "advert-subnet"),
					resource.TestCheckResourceAttr(resourceName, "dead_intvl", "2"),
					resource.TestCheckResourceAttr(resourceName, "hello_intvl", "2"),

					resource.TestCheckResourceAttr(resourceName, "nw_t", "bcast"),

					resource.TestCheckResourceAttr(resourceName, "pfx_suppress", "disable"),
					resource.TestCheckResourceAttr(resourceName, "prio", "1"),
					resource.TestCheckResourceAttr(resourceName, "rexmit_intvl", "2"),
					resource.TestCheckResourceAttr(resourceName, "xmit_delay", "2"),

					testAccCheckAciOSPFInterfacePolicyIdEqual(&ospf_interface_policy_default, &ospf_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccOSPFInterfacePolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccOSPFInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciOSPFInterfacePolicyIdNotEqual(&ospf_interface_policy_default, &ospf_interface_policy_updated),
				),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciOSPFInterfacePolicyIdNotEqual(&ospf_interface_policy_default, &ospf_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciOSPFInterfacePolicy_Update(t *testing.T) {
	var ospf_interface_policy_default models.OSPFInterfacePolicy
	var ospf_interface_policy_updated models.OSPFInterfacePolicy
	resourceName := "aci_ospf_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_default),
				),
			},

			{
				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"advert-subnet"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "advert-subnet"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"advert-subnet", "bfd"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "advert-subnet"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "bfd"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"bfd"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bfd"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"bfd", "mtu-ignore"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "mtu-ignore"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"mtu-ignore"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "mtu-ignore"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"advert-subnet", "bfd", "mtu-ignore", "passive"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "advert-subnet"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "mtu-ignore"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.3", "passive"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"mtu-ignore", "passive"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "mtu-ignore"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "passive"),
				),
			},
			{

				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"passive"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "passive"),
				),
			},
			{
				Config: CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"passive", "mtu-ignore", "bfd", "advert-subnet"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists(resourceName, &ospf_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "passive"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "mtu-ignore"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.3", "advert-subnet"),
				),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciOSPFInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "cost", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "cost", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "cost", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"advert-subnet", "advert-subnet"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "dead_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "dead_intvl", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "dead_intvl", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "hello_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "hello_intvl", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "hello_intvl", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "nw_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "pfx_suppress", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "prio", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "prio", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "rexmit_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "rexmit_intvl", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "rexmit_intvl", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "xmit_delay", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "xmit_delay", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, "xmit_delay", "451"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciOSPFInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFInterfacePolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciOSPFInterfacePolicyExists(name string, ospf_interface_policy *models.OSPFInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("OSPF Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OSPF Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_interface_policyFound := models.OSPFInterfacePolicyFromContainer(cont)
		if ospf_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("OSPF Interface Policy %s not found", rs.Primary.ID)
		}
		*ospf_interface_policy = *ospf_interface_policyFound
		return nil
	}
}

func testAccCheckAciOSPFInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ospf_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ospf_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_interface_policy := models.OSPFInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("OSPF Interface Policy %s Still exists", ospf_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOSPFInterfacePolicyIdEqual(m1, m2 *models.OSPFInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ospf_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciOSPFInterfacePolicyIdNotEqual(m1, m2 *models.OSPFInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ospf_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateOSPFInterfacePolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_interface_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_ospf_interface_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOSPFInterfacePolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccOSPFInterfacePolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFInterfacePolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFInterfacePolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple ospf_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFInterfacePolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing ospf_interface_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccOSPFInterfacePolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ospf_interface_policy"
		cost = "1"
		ctrl = ["advert-subnet"]
		dead_intvl = "2"
		hello_intvl = "2"
		nw_t = "bcast"
		pfx_suppress = "disable"
		prio = "1"
		rexmit_intvl = "2"
		xmit_delay = "2"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccOSPFInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ospf_interface_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ospf_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ospf_interface_policy"
		cost = "1"
		ctrl = ["advert-subnet"]
		dead_intvl = "2"
		hello_intvl = "2"
		nw_t = "bcast"
		pfx_suppress = "disable"
		prio = "1"
		rexmit_intvl = "2"
		xmit_delay = "2"
		
	}
	`)

	return resource
}

func CreateAccOSPFInterfacePolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ospf_interface_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}

func CreateAccOSPFInterfacePolicyUpdatedAttrList(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ospf_interface_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = %s
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
