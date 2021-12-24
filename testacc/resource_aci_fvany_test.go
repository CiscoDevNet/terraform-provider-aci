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

func TestAccAciAny_Basic(t *testing.T) {
	var any_default models.Any
	var any_updated models.Any
	resourceName := "aci_any.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccAnyWithoutVRFdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "vrf_dn", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rName)),
				),
			},
			{
				Config: CreateAccAnyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "vzAny Description"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_any"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_any"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "vrf_dn", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rName)),
					testAccCheckAciAnyIdEqual(&any_default, &any_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccAnyConfigWithAnotherVRFdn(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_updated),
					testAccCheckAciAnyIdNotEqual(&any_default, &any_updated),
				),
			},
			{
				Config:      CreateAccAnyConfigUpdateWithoutVRFdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnyConfig(rName),
			},
		},
	})
}
func TestAccAciAny_Update(t *testing.T) {
	var any_default models.Any
	var any_updated models.Any
	resourceName := "aci_any.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_default),
				),
			},
			{
				Config: CreateAccAnyUpdatedAttr(rName, "match_t", "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					testAccCheckAciAnyIdEqual(&any_default, &any_updated),
				),
			},
			{
				Config: CreateAccAnyUpdatedAttr(rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciAnyIdEqual(&any_default, &any_updated),
				),
			},
		},
	})
}
func TestAccAciAny_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnyConfig(rName),
			},
			{
				Config:      CreateAccAnyConfigWithInvaliVRFdn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzAny (.)+`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected match_t to be one of (.)+ got (.)+`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected pref_gr_memb to be one of (.)+ got (.)+`),
			},
			{
				Config:      CreateAccAnyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAnyConfig(rName),
			},
		},
	})
}

func TestAccAciAny_reltionalParameters(t *testing.T) {
	var any_default models.Any
	var any_rel1 models.Any
	var any_rel2 models.Any
	resourceName := "aci_any.test"
	rName := makeTestVariable(acctest.RandString(5))
	rsRelName1 := acctest.RandString(5)
	rsRelName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_default),
				),
			},
			{
				Config: CreateAccAnyUpdatedAnyIntial(rName, rsRelName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons_if.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, rsRelName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_prov.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName1)),
					testAccCheckAciAnyIdEqual(&any_default, &any_rel1),
				),
			},
			{
				Config: CreateAccAnyUpdatedAnyFinal(rName, rsRelName1, rsRelName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists(resourceName, &any_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons_if.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, rsRelName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, rsRelName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_prov.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_any_to_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rsRelName2)),
					testAccCheckAciAnyIdEqual(&any_default, &any_rel1),
				),
			},
			{
				Config: CreateAccAnyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_any_to_prov.#", "0"),
				),
			},
		},
	})
}

func TestAccAciAny_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciAnyExists(name string, any *models.Any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Any %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Any dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		anyFound := models.AnyFromContainer(cont)
		if anyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Any %s not found", rs.Primary.ID)
		}
		*any = *anyFound
		return nil
	}
}

func testAccCheckAciAnyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing any destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_any" {
			cont, err := client.Get(rs.Primary.ID)
			aci := models.AnyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Any %s Still exists", aci.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAnyIdEqual(any1, any2 *models.Any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if any1.DistinguishedName != any2.DistinguishedName {
			return fmt.Errorf("Any DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAnyIdNotEqual(any1, any2 *models.Any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if any1.DistinguishedName == any2.DistinguishedName {
			return fmt.Errorf("Any DNs are equal")
		}
		return nil
	}
}

func CreateAccAnyWithoutVRFdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing any creation without vrf_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = %s
	}
	resource "aci_vrf" "test" {
		tenant_dn              = aci_tenant.test.id
		name                   = "%s"
	}
	resource "aci_any" "test" {
	}
	`, rName, rName)
	return resource
}
func CreateAccAnyConfigUpdateWithoutVRFdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing any update without required arguement")
	resource := fmt.Sprintln(`
	resource "aci_any" "test" {
		description  = "vzAny Description1"
		annotation   = "tag_any"
		match_t      = "AtmostOne"
		name_alias   = "alias_any"
		pref_gr_memb = "enabled"
	}
	`, rName, rName)
	return resource
}
func CreateAccAnyConfigWithAnotherVRFdn(rName, rOther string) string {
	fmt.Printf("=== STEP  Basic: testing any creation with different vrf name %s \n", rOther)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	`, rName, rOther)
	return resource
}
func CreateAccAnyConfigWithInvaliVRFdn(rName string) string {
	fmt.Printf("=== STEP  Basic: testing any creation with invalid vrf_dn \n")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_tenant.test.id
	}
	`, rName, rName)
	return resource
}
func CreateAccAnyConfig(rName string) string {
	fmt.Println("=== STEP testing any creation with required arguement")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	`, rName, rName)
	return resource
}

func CreateAccAnyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing any creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		description  = "vzAny Description"
		annotation   = "tag_any"
		match_t      = "AtmostOne"
		name_alias   = "alias_any"
		pref_gr_memb = "enabled"
	}
	`, rName, rName)
	return resource
}

func CreateAccAnyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccAnyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  creating multiple anys")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_vrf" "test1" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}	
	resource "aci_vrf" "test2" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	resource "aci_any" "test1" {
		vrf_dn = aci_vrf.test1.id
	}
	resource "aci_any" "test2" {
		vrf_dn = aci_vrf.test2.id
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccAnyUpdatedAnyIntial(rName, rsRelName string) string {
	fmt.Println("=== STEP  Relation Parameters: testing any creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_contract" "test" {
        tenant_dn   =  aci_tenant.test.id
        name = "%s"
    }
	resource "aci_imported_contract" "example" {
		tenant_dn   = aci_tenant.test.id
		name        = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		relation_vz_rs_any_to_cons = [aci_contract.test.id]
		relation_vz_rs_any_to_cons_if = [aci_imported_contract.example.id]
		relation_vz_rs_any_to_prov = [aci_contract.test.id]
	}
	`, rName, rsRelName, rsRelName, rName)

	return resource
}
func CreateAccAnyUpdatedAnyFinal(rName, rsRelName1, rsRelName2 string) string {
	fmt.Println("=== STEP  Relation Parameters: testing any creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_contract" "test1" {
        tenant_dn   =  aci_tenant.test.id
        name = "%s"
    }
	resource "aci_contract" "test2" {
        tenant_dn   =  aci_tenant.test.id
        name = "%s"
    }
	resource "aci_imported_contract" "example1" {
		tenant_dn   = aci_tenant.test.id
		name        = "%s"
	}
	resource "aci_imported_contract" "example2" {
		tenant_dn   = aci_tenant.test.id
		name        = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		relation_vz_rs_any_to_cons = [aci_contract.test1.id, aci_contract.test2.id]
		relation_vz_rs_any_to_cons_if = [aci_imported_contract.example1.id,aci_imported_contract.example2.id]
		relation_vz_rs_any_to_prov = [aci_contract.test1.id, aci_contract.test2.id]
	}
	`, rName, rsRelName1, rsRelName2, rsRelName1, rsRelName2, rName)
	return resource
}
