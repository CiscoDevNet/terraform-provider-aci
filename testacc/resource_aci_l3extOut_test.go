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

func TestAccAciL3Outside_Basic(t *testing.T) {
	var l3outside_default models.L3Outside
	var l3outside_updated models.L3Outside
	resourceName := "aci_l3_outside.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccL3OutsideWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3OutsideWithoutTenantDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3OutsideConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.0", "export"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dampening_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_interleak_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_dom_att", ""),
				),
			},
			{
				Config: CreateAccL3OutsideConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_l3out"),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_out"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.0", "export"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.1", "import"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dampening_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_interleak_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_dom_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_out_to_bd_public_subnet_holder.#", "0"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccl3outsideConfigWithAnotherName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					testAccCheckAciL3OutsideIdNotEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideConfig(rName),
			},
			{
				Config: CreateAccl3outsideConfigWithAnotherTenantDn(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					testAccCheckAciL3OutsideIdNotEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config:      CreateAccl3outsideConfigUpdateWithoutTenantdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccl3outsideConfigUpdateWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccl3outsideConfigUpdateWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of out-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccL3OutsideConfig(rName),
			},
		},
	})
}

func TestAccAciL3Outside_Update(t *testing.T) {
	var l3outside_default models.L3Outside
	var l3outside_updated models.L3Outside
	resourceName := "aci_l3_outside.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3OutsideConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_default),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF32"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF33"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS4"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF42"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF43"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS5"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS6"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", "CS7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS7"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideUpdatedAttrList(rName, "enforce_rtctrl", StringListtoString([]string{"import", "export"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.0", "import"),
					resource.TestCheckResourceAttr(resourceName, "enforce_rtctrl.1", "export"),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
		},
	})
}
func TestAccAciL3Outside_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3OutsideConfig(rName),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttrList(rName, "enforce_rtctrl", StringListtoString([]string{"import"})),
				ExpectError: regexp.MustCompile(`Invalid Configuration Unenforced Route Control is not supported for Export direction.`),
			},
			{
				Config:      CreateAccl3outsideConfigWithInvalidTenantdn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extOut (.)+`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttrList(rName, "enforce_rtctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected enforce_rtctrl.0 to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttrList(rName, "enforce_rtctrl", StringListtoString([]string{"export", "export"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+ got (.)+`),
			},
			{
				Config:      CreateAccL3OutsideUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3OutsideConfig(rName),
			},
		},
	})
}

func TestAccAciL3Outside_reltionalParameters(t *testing.T) {
	var l3outside_default models.L3Outside
	var l3outside_updated models.L3Outside
	resourceName := "aci_l3_outside.test"
	rName := makeTestVariable(acctest.RandString(5))
	rsRelName1 := acctest.RandString(5)
	rsRelName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3OutsideConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_default),
				),
			},
			{
				Config: CreateAccL3OutsidUpdatedL3OutsideIntial(rName, rsRelName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_out_to_bd_public_subnet_holder.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_dom_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_interleak_pol", fmt.Sprintf("uni/tn-%s/prof-%s", rName, rsRelName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_ectx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rsRelName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dampening_pol.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_l3ext_rs_dampening_pol.*", map[string]string{
						"af":                     "ipv4-ucast",
						"tn_rtctrl_profile_name": fmt.Sprintf("uni/tn-%s/prof-%s", rName, rsRelName1),
					}),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsidUpdatedL3OutsideFinal(rName, rsRelName1, rsRelName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckl3OutsideExists(resourceName, &l3outside_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_out_to_bd_public_subnet_holder.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_dom_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_interleak_pol", fmt.Sprintf("uni/tn-%s/prof-%s", rName, rsRelName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_ectx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rsRelName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dampening_pol.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_l3ext_rs_dampening_pol.*", map[string]string{
						"af":                     "ipv4-ucast",
						"tn_rtctrl_profile_name": fmt.Sprintf("uni/tn-%s/prof-%s", rName, rsRelName1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_l3ext_rs_dampening_pol.*", map[string]string{
						"af":                     "ipv6-ucast",
						"tn_rtctrl_profile_name": fmt.Sprintf("uni/tn-%s/prof-%s", rName, rsRelName2),
					}),
					testAccCheckAciL3OutsideIdEqual(&l3outside_default, &l3outside_updated),
				),
			},
			{
				Config: CreateAccL3OutsideConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_out_to_bd_public_subnet_holder.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_dom_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_interleak_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_ectx", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dampening_pol.#", "0"),
				),
			},
		},
	})
}

func TestAccAciL3Outside_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3OutsideConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckl3OutsideExists(name string, l3Outside *models.L3Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("l3Outside %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No l3Outside dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3OutsideFound := models.L3OutsideFromContainer(cont)
		if l3OutsideFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("l3Outside %s not found", rs.Primary.ID)
		}
		*l3Outside = *l3OutsideFound
		return nil
	}
}

func testAccCheckAciL3OutsideDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing L3Outside destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3_outside" {
			cont, err := client.Get(rs.Primary.ID)
			aci := models.L3OutsideFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3Outside %s Still exists", aci.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3OutsideIdEqual(l3outside1, l3outside2 *models.L3Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if l3outside1.DistinguishedName != l3outside2.DistinguishedName {
			return fmt.Errorf("L3Outside DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3OutsideIdNotEqual(l3outside1, l3outside2 *models.L3Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if l3outside1.DistinguishedName == l3outside2.DistinguishedName {
			return fmt.Errorf("L3Outside DNs are equal")
		}
		return nil
	}
}

func CreateAccL3OutsideWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3Outside creation without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
	  }
	`, rName)
	return resource
}

func CreateAccL3OutsideWithoutTenantDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3Outside creation without giving Tenant dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
	  }
	`, rName, rName)
	return resource
}

func CreateAccl3outsideConfigUpdateWithoutTenantdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3outside update without giving Tenant Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name    = "%s"
	  }
	`, rName, rName)
	return resource
}

func CreateAccl3outsideConfigUpdateWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3outside update without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn  = aci_tenant.test.id
	  }
	`, rName)
	return resource
}

func CreateAccl3outsideConfigUpdateWithInvalidName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing L3outside creation with parent resource name %s and name %s\n", parentName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name= "%s"
	  }
	`, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithAnotherName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing l3outside creation with different l3outside name %s \n", rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	`, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithAnotherTenantDn(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing l3outside creation with different parent %s \n", parentName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	`, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithInvalidTenantdn(rName string) string {
	fmt.Printf("=== STEP  Basic: testing l3outside creation with invalid tenant dn \n")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	}
	resource"aci_application_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name       = "%s"
	}
	resource "aci_l3_outside" "test" {
		name = "%s"
		tenant_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccL3OutsideConfig(rName string) string {
	fmt.Println("=== STEP testing L3Outside creation with required parameters only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn      = aci_tenant.test.id
	    name           = "%s"
	    description    = "from terraform"	
	    annotation     = "tag_l3out"
		enforce_rtctrl = ["export","import"]
		name_alias     = "alias_out"
		target_dscp    = "CS0"
	  }
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccL3OutsideUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
		%s = %s
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccL3OutsideConfigMultiple(rName string) string {
	fmt.Println("=== STEP  creating multiple l3Outside")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
			  tenant_dn      = aci_tenant.test.id
			  name           = "%s"
	  }
	  resource "aci_l3_outside" "test1" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	resource "aci_l3_outside" "test2" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccL3OutsidUpdatedL3OutsideIntial(rName, rsRelName string) string {
	fmt.Println("=== STEP  Relation Parameters: testing l3outside creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_bgp_route_control_profile" "test" {
		parent_dn                  = aci_tenant.test.id
		name                       = "%s"
	  }
	  resource "aci_vrf" "test" {
		tenant_dn              = aci_tenant.test.id
		name                   = "%s"
	  }
	  resource "aci_l3_outside" "test" {
		 tenant_dn      = aci_tenant.test.id
	     name           = "%s"
		  relation_l3ext_rs_dampening_pol {
		  af =  "ipv4-ucast"
		  tn_rtctrl_profile_name = aci_bgp_route_control_profile.test.id
		  }
		  relation_l3ext_rs_ectx = aci_vrf.test.id
		  relation_l3ext_rs_interleak_pol = aci_bgp_route_control_profile.test.id
	  }
	`, rName, rsRelName, rsRelName, rName)

	return resource
}
func CreateAccL3OutsidUpdatedL3OutsideFinal(rName, rsRelName1, rsRelName2 string) string {
	fmt.Println("=== STEP  Relation Parameters: testing l3outside creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_bgp_route_control_profile" "test" {
		parent_dn                  = aci_tenant.test.id
		name                       = "%s"
	  }
	  resource "aci_bgp_route_control_profile" "test1" {
		parent_dn                  = aci_tenant.test.id
		name                       = "%s"
	  }
	  resource "aci_vrf" "test" {
		tenant_dn              = aci_tenant.test.id
		name                   = "%s"
	  }
	  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        name           = "%s"
        relation_l3ext_rs_dampening_pol {
		af =  "ipv4-ucast"
		tn_rtctrl_profile_name = aci_bgp_route_control_profile.test.id
	    }
        relation_l3ext_rs_dampening_pol {
		af =  "ipv6-ucast"
		tn_rtctrl_profile_name = aci_bgp_route_control_profile.test1.id
	    }
	relation_l3ext_rs_ectx = aci_vrf.test.id
	relation_l3ext_rs_interleak_pol = aci_bgp_route_control_profile.test1.id
}
	`, rName, rsRelName1, rsRelName2, rsRelName2, rName)
	return resource
}
