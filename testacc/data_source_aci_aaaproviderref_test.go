package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciProviderGroupMemberDataSource_Basic(t *testing.T) {
	resourceName := "aci_login_domain_provider.test"
	dataSourceName := "data.aci_login_domain_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateProviderGroupMemberDSWithoutRequired(aaaDuoProviderGroupName, rName, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateProviderGroupMemberDSWithoutRequired(aaaDuoProviderGroupName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccProviderGroupMemberConfigDataSource(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "order", resourceName, "order"),
				),
			},
			{
				Config:      CreateAccProviderGroupMemberDataSourceUpdate(aaaDuoProviderGroupName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccProviderGroupMemberDSWithInvalidName(aaaDuoProviderGroupName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccProviderGroupMemberDataSourceUpdatedResource(aaaDuoProviderGroupName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccProviderGroupMemberConfigDataSource(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateProviderGroupMemberDSWithoutRequired(aaaDuoProviderGroupName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	data "aci_login_domain_provider" "test" {
	#	parent_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
	#	name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaDuoProviderGroupName, rName)
}

func CreateAccProviderGroupMemberDSWithInvalidName(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "${aci_login_domain_provider.test.name}_invalid"
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccProviderGroupMemberDataSourceUpdate(aaaDuoProviderGroupName, rName, key, value string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		%s = "%s"
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName, key, value)
	return resource
}

func CreateAccProviderGroupMemberDataSourceUpdatedResource(aaaDuoProviderGroupName, rName, key, value string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_login_domain_provider" "test" {
		parent_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName, key, value)
	return resource
}
