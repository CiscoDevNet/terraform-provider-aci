package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLoginDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_login_domain.test"
	dataSourceName := "data.aci_login_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLoginDomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoginDomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_group", resourceName, "provider_group"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm", resourceName, "realm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_sub_type", resourceName, "realm_sub_type"),
				),
			},
			{
				Config:      CreateAccLoginDomainDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLoginDomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLoginDomainDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLoginDomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Testing login_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}

	data "aci_login_domain" "test" {
	
		name  = aci_login_domain.test.name
		depends_on = [ aci_login_domain.test ]
	}
	`, rName)
	return resource
}

func CreateLoginDomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Testing login_domain Data Source without ", attrName)
	rBlock := `
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_login_domain" "test" {
	
	#	name  = aci_login_domain.test.name
		depends_on = [ aci_login_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLoginDomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Testing login_domain Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}

	data "aci_login_domain" "test" {
	
		name  = "${aci_login_domain.test.name}_invalid"
		depends_on = [ aci_login_domain.test ]
	}
	`, rName)
	return resource
}

func CreateAccLoginDomainDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  Testing login_domain Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}

	data "aci_login_domain" "test" {
	
		name  = aci_login_domain.test.name
		%s = "%s"
		depends_on = [ aci_login_domain.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLoginDomainDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  Testing login_domain Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_login_domain" "test" {
	
		name  = aci_login_domain.test.name
		depends_on = [ aci_login_domain.test ]
	}
	`, rName, key, value)
	return resource
}
