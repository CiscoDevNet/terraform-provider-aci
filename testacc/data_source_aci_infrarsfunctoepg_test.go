package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEPGsUsingFunctionDataSource_Basic(t *testing.T) {
	resourceName := "aci_epgs_using_function.test"
	dataSourceName := "data.aci_epgs_using_function.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	tDn := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEPGsUsingFunctionDSWithoutRequired(infraAttEntityPName, tDn, "access_generic_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEPGsUsingFunctionDSWithoutRequired(infraAttEntityPName, tDn, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			// {
			// 	Config:      CreateEPGsUsingFunctionDSWithoutRequired(infraAttEntityPName, tDn, "encap"),
			// 	ExpectError: regexp.MustCompile(`Missing required argument`),
			// },
			{
				Config: CreateAccEPGsUsingFunctionConfigDataSource(infraAttEntityPName, tDn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "access_generic_dn", resourceName, "access_generic_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tdn", resourceName, "tdn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap", resourceName, "encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instr_imedcy", resourceName, "instr_imedcy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "primary_encap", resourceName, "primary_encap"),
				),
			},
			{
				Config:      CreateAccEPGsUsingFunctionDataSourceUpdate(infraAttEntityPName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionDSWithInvalidParentDn(infraAttEntityPName, tDn),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEPGsUsingFunctionDataSourceUpdatedResource(infraAttEntityPName, tDn, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEPGsUsingFunctionConfigDataSource(infraAttEntityPName, tDn string) string {
	fmt.Println("=== STEP  testing epgs_using_function Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}

	data "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_epgs_using_function.test.access_generic_dn
		tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
	`, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName)
	return resource
}

func CreateEPGsUsingFunctionDSWithoutRequired(infraAttEntityPName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`
	switch attrName {
	case "access_generic_dn":
		rBlock += `
	data "aci_epgs_using_function" "test" {
	#	access_generic_dn  = aci_epgs_using_function.test.access_generic_dn
		tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
		`
	case "tdn":
		rBlock += `
	data "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_epgs_using_function.test.access_generic_dn
	#	tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
	`
	case "encap":
		rBlock += `
	data "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_epgs_using_function.test.access_generic_dn
		tdn  = aci_epgs_using_function.test.tdn
	#	encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName)
}

func CreateAccEPGsUsingFunctionDSWithInvalidParentDn(infraAttEntityPName, tDn string) string {
	fmt.Println("=== STEP  testing epgs_using_function Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}

	data "aci_epgs_using_function" "test" {
		access_generic_dn  = "${aci_epgs_using_function.test.access_generic_dn}_invalid"
		tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
	`, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName)
	return resource
}

func CreateAccEPGsUsingFunctionDataSourceUpdate(infraAttEntityPName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing epgs_using_function Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}

	data "aci_epgs_using_function" "test" {
		access_generic_dn  = "${aci_epgs_using_function.test.access_generic_dn}_invalid"
		tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		%s = "%s"
		depends_on = [ aci_epgs_using_function.test ]
	}
	`, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, key, value)
	return resource
}

func CreateAccEPGsUsingFunctionDataSourceUpdatedResource(infraAttEntityPName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing epgs_using_function Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
		%s = "%s"
	}

	data "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_epgs_using_function.test.access_generic_dn
		tdn  = aci_epgs_using_function.test.tdn
		encap = aci_epgs_using_function.test.encap
		depends_on = [ aci_epgs_using_function.test ]
	}
	`, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, infraAttEntityPName, key, value)
	return resource
}
