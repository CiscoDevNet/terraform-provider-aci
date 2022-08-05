package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLocalUserDataSource_Basic(t *testing.T) {
	resourceName := "aci_local_user.test"
	dataSourceName := "data.aci_local_user.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	password := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	passwordOther := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLocalUserDSWithoutRequired(rName, password, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLocalUserConfigDataSource(rName, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_status", resourceName, "account_status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cert_attribute", resourceName, "cert_attribute"),
					resource.TestCheckResourceAttrPair(dataSourceName, "clear_pwd_history", resourceName, "clear_pwd_history"),
					resource.TestCheckResourceAttrPair(dataSourceName, "email", resourceName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceName, "expiration", resourceName, "expiration"),
					resource.TestCheckResourceAttrPair(dataSourceName, "expires", resourceName, "expires"),
					resource.TestCheckResourceAttrPair(dataSourceName, "first_name", resourceName, "first_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "last_name", resourceName, "last_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "otpenable", resourceName, "otpenable"),
					resource.TestCheckResourceAttrPair(dataSourceName, "otpkey", resourceName, "otpkey"),
					resource.TestCheckResourceAttrPair(dataSourceName, "phone", resourceName, "phone"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pwd_life_time", resourceName, "pwd_life_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pwd_update_required", resourceName, "pwd_update_required"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rbac_string", resourceName, "rbac_string"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unix_user_id", resourceName, "unix_user_id"),
				),
			},
			{
				Config:      CreateAccLocalUserDataSourceUpdate(rName, password, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLocalUserDSWithInvalidName(rName, password),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLocalUserDataSourceUpdatedResource(rName, passwordOther, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLocalUserConfigDataSource(rName, pwd string) string {
	fmt.Println("=== STEP  testing local_user Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}

	data "aci_local_user" "test" {
	
		name  = aci_local_user.test.name
		depends_on = [ aci_local_user.test ]
	}
	`, rName, pwd)
	return resource
}

func CreateLocalUserDSWithoutRequired(rName, pwd, attrName string) string {
	fmt.Println("=== STEP  Basic: testing local_user Data Source without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_local_user" "test" {
	
	#	name  = aci_local_user.test.name
		depends_on = [ aci_local_user.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, pwd)
}

func CreateAccLocalUserDSWithInvalidName(rName, pwd string) string {
	fmt.Println("=== STEP  testing local_user Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}

	data "aci_local_user" "test" {
	
		name  = "${aci_local_user.test.name}_invalid"
		depends_on = [ aci_local_user.test ]
	}
	`, rName, pwd)
	return resource
}

func CreateAccLocalUserDataSourceUpdate(rName, pwd, key, value string) string {
	fmt.Println("=== STEP  testing local_user Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}

	data "aci_local_user" "test" {
	
		name  = aci_local_user.test.name
		%s = "%s"
		depends_on = [ aci_local_user.test ]
	}
	`, rName, pwd, key, value)
	return resource
}

func CreateAccLocalUserDataSourceUpdatedResource(rName, pwd, key, value string) string {
	fmt.Println("=== STEP  testing local_user Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
		%s = "%s"
	}

	data "aci_local_user" "test" {
	
		name  = aci_local_user.test.name
		depends_on = [ aci_local_user.test ]
	}
	`, rName, pwd, key, value)
	return resource
}
