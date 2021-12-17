package acctest

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

const tdn1 = "topology/pod-1/paths-101/pathep-[eth1/12]"
const tdn2 = "topology/pod-1/paths-101/pathep-[eth1/6]"
const multdn1 = "topology/pod-1/paths-101/pathep-[eth1/21]"
const multdn2 = "topology/pod-1/paths-101/pathep-[eth1/30]"
const multdn3 = "topology/pod-1/paths-101/pathep-[eth1/1]"

func TestAccAciStaticPath_Basic(t *testing.T) {
	var static_path_default models.StaticPath
	var static_path_updated models.StaticPath
	resourceName := "aci_epg_to_static_path.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateStaticPathWithoutApplicationEpg(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateStaticPathWithoutTdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateStaticPathWithoutEncap(rName),
				ExpectError: regexp.MustCompile(`Validation failed: Encap not specifiedRn=(.)+`),
			},
			{
				Config: CreateAccStaticPathConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "tdn", tdn1),
				),
			},
			{
				Config: CreateAccStaticPathConfigWithOptionalParameters(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "annotation"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "vlan-5"),
					resource.TestCheckResourceAttr(resourceName, "tdn", tdn1),
					testAccCheckAciStaticPathIdEqual(&static_path_default, &static_path_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccStaticPathRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccStaticPathConfigWithEpgAndTdn(rOtherName, tdn1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_updated),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rOtherName, rOtherName, rOtherName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", tdn1),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					testAccCheckAciStaticPathIdNotEqual(&static_path_default, &static_path_updated),
				),
			},
			{
				Config: CreateAccStaticPathConfig(rName),
			},
			{
				Config: CreateAccStaticPathConfigWithEpgAndTdn(rName, tdn2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_updated),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", tdn2),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					testAccCheckAciStaticPathIdNotEqual(&static_path_default, &static_path_updated),
				),
			},
		},
	})
}

func TestAccAciStaticPath_Update(t *testing.T) {
	var static_path_default models.StaticPath
	var static_path_updated models.StaticPath
	resourceName := "aci_epg_to_static_path.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccStaticPathConfigForUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_default),
				),
			},
			{
				Config: CreateAccStaticPathUpdated(rName, "mode", "untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists(resourceName, &static_path_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "untagged"),
					testAccCheckAciStaticPathIdEqual(&static_path_default, &static_path_updated),
				),
			},
		},
	})
}

func TestAccAciStaticPath_NegativCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccStaticPathConfigNegative(rName),
			},
			{
				Config:      CreateAccStaticPathWithInvalidEpgDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvRsPathAtt (.)+,`),
			},
			{
				Config:      CreateAccStaticPathWithInvalidTDn(rName, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvRsPathAtt (.)+,`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegativeEncap(rName, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name encap, class fvRsPathAtt (.)+`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, "instr_imedcy", randomValue),
				ExpectError: regexp.MustCompile(`expected instr_imedcy to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected mode to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, "primary_encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name primaryEncap, class fvRsPathAtt (.)+`),
			},
			{
				Config:      CreateAccStaticPathUpdatedNegative(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccStaticPathConfigNegative(rName),
			},
		},
	})
}

func TestAccAciStaticPath_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccStaticPathsConfig(rName),
			},
		},
	})
}

func CreateAccStaticPathsConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path multiple creation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test1" {
		application_epg_dn  = aci_application_epg.test.id
		tdn  = "%s"
		encap = "vlan-101"
	}

	resource "aci_epg_to_static_path" "test2" {
		application_epg_dn  = aci_application_epg.test.id
		tdn  = "%s"
		encap = "vlan-102"
	}

	resource "aci_epg_to_static_path" "test3" {
		application_epg_dn  = aci_application_epg.test.id
		tdn  = "%s"
		encap = "vlan-103"
	}
	`, rName, rName, rName, multdn1, multdn2, multdn3)
	return resource
}

func CreateAccStaticPathConfigWithEpgAndTdn(rName, tdn string) string {
	fmt.Printf("=== STEP  Basic: testing epg_to_static_path creation with resource name %s and tdn %s\n", rName, tdn)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn  = aci_application_epg.test.id
		tdn  = "%s"
		encap = "vlan-1"
	}
	`, rName, rName, rName, tdn)
	return resource
}

func CreateAccStaticPathRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path updation without required parameters")
	resource := fmt.Sprintln(`

	resource "aci_epg_to_static_path" "test" {
		annotation = "tag"
		description = "description"
		instr_imedcy = "immediate"
		mode = "native"
		primary_encap = "vlan-500"
	  }
	`)
	return resource
}

func CreateAccStaticPathConfigWithOptionalParameters(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation with all optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn  = aci_application_epg.test.id
		tdn  = "%s"
		encap = "vlan-1"
		annotation = "annotation"
		description = "description"
		instr_imedcy = "immediate"
		mode = "native"
		primary_encap = "vlan-5"
	  }
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateAccStaticPathWithInvalidTDn(rName, rVal string) string {
	fmt.Println("=== STEP  Negative cases: testing epg_to_static_path creation with invalid tdn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test"{
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-27"
	}
	`, rName, rName, rName, rVal)
	return resource
}

func CreateAccStaticPathWithInvalidEpgDn(rName string) string {
	fmt.Println("=== STEP  Negative cases: testing epg_to_static_path creation with invalid application_epg_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_tenant.test.id
		tdn = "%s"
		encap = "vlan-27"
	}
	`, rName, tdn1)
	return resource
}

func CreateAccStaticPathUpdatedNegativeEncap(rName, encap string) string {
	fmt.Printf("=== STEP  Negative cases: testing epg_to_static_path creation with encap = %s\n", encap)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "%s"
	}
	`, rName, rName, rName, tdn1, encap)
	return resource
}

func CreateAccStaticPathUpdatedNegative(rName, key, value string) string {
	fmt.Printf("=== STEP  Negative cases: testing epg_to_static_path creation with %s = %s\n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-27"
		%s = "%s"
	}
	`, rName, rName, rName, tdn1, key, value)
	return resource
}

func CreateAccStaticPathUpdated(rName, key, value string) string {
	fmt.Printf("=== STEP  Update: testing epg_to_static_path creation with %s = %s\n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-15"
		%s = "%s"
	}
	`, rName, rName, rName, tdn1, key, value)
	return resource
}

func CreateAccStaticPathConfigForUpdate(rName string) string {
	fmt.Println("=== STEP  Update: testing epg_to_static_path creation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-15"
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateAccStaticPathConfigNegative(rName string) string {
	fmt.Println("=== STEP  Negative cases: testing epg_to_static_path creation with required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-27"
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateAccStaticPathConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation with required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
		encap = "vlan-1"
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateStaticPathWithoutEncap(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation without encap")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "example" {
		application_epg_dn = aci_application_epg.test.id
		tdn = "%s"
	}
	`, rName, rName, rName, tdn1)
	return resource
}

func CreateStaticPathWithoutTdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation without tDn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epg_to_static_path" "example" {
		application_epg_dn = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`, rName, rName, rName)
	return resource
}

func CreateStaticPathWithoutApplicationEpg() string {
	fmt.Println("=== STEP  Basic: testing epg_to_static_path creation without application_epg_dn")
	resource := fmt.Sprintf(`
	resource "aci_epg_to_static_path" "test" {
		tdn = "%s"
		encap = "vlan-1"
	}
	`, tdn1)
	return resource
}

func testAccCheckAciStaticPathExists(name string, static_path *models.StaticPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Static Path %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Path dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		static_pathFound := models.StaticPathFromContainer(cont)
		if static_pathFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Static Path %s not found", rs.Primary.ID)
		}
		*static_path = *static_pathFound
		return nil
	}
}

func testAccCheckAciStaticPathDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing epg_to_static_path destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_epg_to_static_path" {
			cont, err := client.Get(rs.Primary.ID)
			static_path := models.StaticPathFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Static Path %s Still exists", static_path.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciStaticPathIdEqual(m1, m2 *models.StaticPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("static_path DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciStaticPathIdNotEqual(m1, m2 *models.StaticPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("static_path DNs are equal")
		}
		return nil
	}
}
