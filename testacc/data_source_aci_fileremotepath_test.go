package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRemotePathofaFileDataSource_Basic(t *testing.T) {
	resourceName := "aci_file_remote_path.test"
	dataSourceName := "data.aci_file_remote_path.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	host := "cisco.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFileRemotePathDSWithoutRequired(rName, host, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFileRemotePathConfigDataSource(rName, host),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_type", resourceName, "auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host", resourceName, "host"),
					resource.TestCheckResourceAttrPair(dataSourceName, "protocol", resourceName, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_path", resourceName, "remote_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_port", resourceName, "remote_port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_name", resourceName, "user_name"),
				),
			},
			{
				Config:      CreateAccFileRemotePathDataSourceUpdate(rName, host, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFileRemotePathDSWithInvalidName(rName, host),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFileRemotePathDataSourceUpdatedResource(rName, host, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFileRemotePathConfigDataSource(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"

	}

	data "aci_file_remote_path" "test" {
		name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, host)
	return resource
}

func CreateFileRemotePathDSWithoutRequired(rName, host, attrName string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path Data Source without ", attrName)
	rBlock := `
	
	resource "aci_file_remote_path" "test" {
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"

	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_file_remote_path" "test" {
	#	name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, host)
}

func CreateAccFileRemotePathDSWithInvalidName(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"

	}

	data "aci_file_remote_path" "test" {
	
		name  = "${aci_file_remote_path.test.name}_invalid"
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, host)
	return resource
}

func CreateAccFileRemotePathDataSourceUpdate(rName, host, key, value string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"

	}

	data "aci_file_remote_path" "test" {
	
		name  = aci_file_remote_path.test.name
		%s = "%s"
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, host, key, value)
	return resource
}

func CreateAccFileRemotePathDataSourceUpdatedResource(rName, host, key, value string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"

		%s = "%s"
	}

	data "aci_file_remote_path" "test" {
	
		name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, host, key, value)
	return resource
}
