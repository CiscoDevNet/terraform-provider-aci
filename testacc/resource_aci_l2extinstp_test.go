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

func TestAccAciL2outExternalEpg_Basic(t *testing.T) {
	var l2_out_extepg_default models.L2outExternalEpg
	var l2_out_extepg_updated models.L2outExternalEpg
	resourceName := "aci_l2out_extepg.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	l2extOutName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL2outExternalEpgWithoutRequired(fvTenantName, l2extOutName, rName, "l2_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL2outExternalEpgWithoutRequired(fvTenantName, l2extOutName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_default),
					resource.TestCheckResourceAttr(resourceName, "l2_outside_dn", fmt.Sprintf("uni/tn-%s/l2out-%s", fvTenantName, l2extOutName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", ""),
				),
			},
			{
				Config: CreateAccL2outExternalEpgConfigWithOptionalValues(fvTenantName, l2extOutName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "l2_outside_dn", fmt.Sprintf("uni/tn-%s/l2out-%s", fvTenantName, l2extOutName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l2_out_extepg"),

					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),

					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),

					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),

					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "example"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL2outExternalEpgConfigUpdatedName(fvTenantName, l2extOutName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL2outExternalEpgRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2outExternalEpgConfigWithRequiredParams(rName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "l2_outside_dn", fmt.Sprintf("uni/tn-%s/l2out-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciL2outExternalEpgIdNotEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
			},
			{
				Config: CreateAccL2outExternalEpgConfigWithRequiredParams(rName, rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "l2_outside_dn", fmt.Sprintf("uni/tn-%s/l2out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL2outExternalEpgIdNotEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
		},
	})
}

func TestAccAciL2outExternalEpg_Update(t *testing.T) {
	var l2_out_extepg_default models.L2outExternalEpg
	var l2_out_extepg_updated models.L2outExternalEpg
	resourceName := "aci_l2out_extepg.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	l2extOutName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_default),
				),
			},

			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists(resourceName, &l2_out_extepg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciL2outExternalEpgIdEqual(&l2_out_extepg_default, &l2_out_extepg_updated),
				),
			},
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
			},
		},
	})
}

func TestAccAciL2outExternalEpg_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	l2extOutName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
			},
			{
				Config:      CreateAccL2outExternalEpgWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName),
			},
		},
	})
}

func testAccCheckAciL2outExternalEpgExists(name string, l2_out_extepg *models.L2outExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2Out Extepg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2Out Extepg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_out_extepgFound := models.L2outExternalEpgFromContainer(cont)
		if l2_out_extepgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2Out Extepg %s not found", rs.Primary.ID)
		}
		*l2_out_extepg = *l2_out_extepgFound
		return nil
	}
}

func testAccCheckAciL2outExternalEpgDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l2_out_extepg destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l2out_extepg" {
			cont, err := client.Get(rs.Primary.ID)
			l2_out_extepg := models.L2outExternalEpgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2Out Extepg %s Still exists", l2_out_extepg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL2outExternalEpgIdEqual(m1, m2 *models.L2outExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l2_out_extepg DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL2outExternalEpgIdNotEqual(m1, m2 *models.L2outExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l2_out_extepg DNs are equal")
		}
		return nil
	}
}

func CreateL2outExternalEpgWithoutRequired(fvTenantName, l2extOutName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_out_extepg creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "l2_outside_dn":
		rBlock += `
	resource "aci_l2out_extepg" "test" {
	#	l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l2extOutName, rName)
}

func CreateAccL2outExternalEpgConfigWithRequiredParams(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  testing l2_out_extepg creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}
	`, fvTenantName, l2extOutName, rName)
	return resource
}
func CreateAccL2outExternalEpgConfigUpdatedName(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  testing l2_out_extepg creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}
	`, fvTenantName, l2extOutName, rName)
	return resource
}

func CreateAccL2outExternalEpgConfig(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  testing l2_out_extepg creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
	}
	`, fvTenantName, l2extOutName, rName)
	return resource
}

func CreateAccL2outExternalEpgWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l2_out_extepg creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccL2outExternalEpgConfigWithOptionalValues(fvTenantName, l2extOutName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l2_out_extepg creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = "${aci_l2_outside.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_out_extepg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		target_dscp = "AF11"
		exception_tag = "example"
	}
	`, fvTenantName, l2extOutName, rName)

	return resource
}

func CreateAccL2outExternalEpgRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l2_out_extepg updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l2out_extepg" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_out_extepg"
		flood_on_encap = "enabled"
		match_t = "All"
		pref_gr_memb = "include"
		prio = "level1"
		target_dscp = "1"
		
	}
	`)

	return resource
}

func CreateAccL2outExternalEpgUpdatedAttr(fvTenantName, l2extOutName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l2_out_extepg attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l2out_extepg" "test" {
		l2_outside_dn  = aci_l2_outside.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l2extOutName, rName, attribute, value)
	return resource
}
