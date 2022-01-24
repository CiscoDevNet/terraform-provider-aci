package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFexBundleGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_fex_bundle_group.test"
	dataSourceName := "data.aci_fex_bundle_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	infraFexPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFexBundleGroupDSWithoutRequired(infraFexPName, rName, "fex_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFexBundleGroupDSWithoutRequired(infraFexPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFexBundleGroupConfigDataSource(infraFexPName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "fex_profile_dn", resourceName, "fex_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccFexBundleGroupDataSourceUpdate(infraFexPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFexBundleGroupDSWithInvalidParentDn(infraFexPName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccFexBundleGroupDataSourceUpdatedResource(infraFexPName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFexBundleGroupConfigDataSource(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing fex_bundle_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}

	data "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = aci_fex_bundle_group.test.name
		depends_on = [ aci_fex_bundle_group.test ]
	}
	`, infraFexPName, rName)
	return resource
}

func CreateFexBundleGroupDSWithoutRequired(infraFexPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fex_bundle_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "fex_profile_dn":
		rBlock += `
	data "aci_fex_bundle_group" "test" {
	#	fex_profile_dn  = aci_fex_profile.test.id
		name  = aci_fex_bundle_group.test.name
		depends_on = [ aci_fex_bundle_group.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
	#	name  = aci_fex_bundle_group.test.name
		depends_on = [ aci_fex_bundle_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraFexPName, rName)
}

func CreateAccFexBundleGroupDSWithInvalidParentDn(infraFexPName, rName string) string {
	fmt.Println("=== STEP  testing fex_bundle_group Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}

	data "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "${aci_fex_bundle_group.test.name}_invalid"
		depends_on = [ aci_fex_bundle_group.test ]
	}
	`, infraFexPName, rName)
	return resource
}

func CreateAccFexBundleGroupDataSourceUpdate(infraFexPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing fex_bundle_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
	}

	data "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = aci_fex_bundle_group.test.name
		%s = "%s"
		depends_on = [ aci_fex_bundle_group.test ]
	}
	`, infraFexPName, rName, key, value)
	return resource
}

func CreateAccFexBundleGroupDataSourceUpdatedResource(infraFexPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing fex_bundle_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_fex_bundle_group" "test" {
		fex_profile_dn  = aci_fex_profile.test.id
		name  = aci_fex_bundle_group.test.name
		depends_on = [ aci_fex_bundle_group.test ]
	}
	`, infraFexPName, rName, key, value)
	return resource
}
