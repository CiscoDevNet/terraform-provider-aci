package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSDestinationDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_accounting_destination.test"
	dataSourceName := "data.aci_tacacs_accounting_destination.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	host := fmt.Sprintf("%s.com", acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTACACSDestinationDSWithoutRequired(tacacsGroupName, host, "tacacs_accounting_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTACACSDestinationDSWithoutRequired(tacacsGroupName, host, "host"),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},
			{
				Config: CreateAccTACACSDestinationConfigDataSource(tacacsGroupName, host),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tacacs_accounting_dn", resourceName, "tacacs_accounting_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host", resourceName, "host"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_protocol", resourceName, "auth_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccTACACSDestinationDataSourceUpdate(tacacsGroupName, host, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSDestinationDSWithInvalidParentDn(tacacsGroupName, host),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccTACACSDestinationDataSourceUpdatedResource(tacacsGroupName, host, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSDestinationConfigDataSource(tacacsGroupName, host string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = aci_tacacs_accounting_destination.test.host
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host)
	return resource
}

func CreateTACACSDestinationDSWithoutRequired(tacacsGroupName, host, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
	}
	`
	switch attrName {
	case "tacacs_accounting_dn":
		rBlock += `
	data "aci_tacacs_accounting_destination" "test" {
	#	tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = aci_tacacs_accounting_destination.test.host	
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
		`
	case "host":
		rBlock += `
	data "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
	#	host  = aci_tacacs_accounting_destination.test.host
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, tacacsGroupName, host)
}

func CreateAccTACACSDestinationDSWithInvalidParentDn(tacacsGroupName, host string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = "${aci_tacacs_accounting.test.id}_invalid"
		host  = "${aci_tacacs_accounting_destination.test.host}"
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host)
	return resource
}

func CreateAccTACACSDestinationDataSourceUpdate(tacacsGroupName, host, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = aci_tacacs_accounting_destination.test.host
		%s = "%s"
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, key, value)
	return resource
}

func CreateAccTACACSDestinationDataSourceUpdatedResource(tacacsGroupName, host, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tacacs_accounting.test.id
		host  = aci_tacacs_accounting_destination.test.host
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, key, value)
	return resource
}
