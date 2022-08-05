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

func TestAccAciLocalUser_Basic(t *testing.T) {
	var local_user_default models.LocalUser
	var local_user_updated models.LocalUser
	resourceName := "aci_local_user.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	password := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password1 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password2 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLocalUserWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLocalUserConfig(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "account_status", "active"),
					resource.TestCheckResourceAttr(resourceName, "cert_attribute", ""),
					resource.TestCheckResourceAttr(resourceName, "email", ""),
					resource.TestCheckResourceAttr(resourceName, "expiration", "never"),
					resource.TestCheckResourceAttr(resourceName, "expires", "no"),
					resource.TestCheckResourceAttr(resourceName, "first_name", ""),
					resource.TestCheckResourceAttr(resourceName, "last_name", ""),
					resource.TestCheckResourceAttr(resourceName, "otpenable", "no"),
					resource.TestCheckResourceAttr(resourceName, "otpkey", "DISABLEDDISABLED"),
					resource.TestCheckResourceAttr(resourceName, "phone", ""),
					resource.TestCheckResourceAttr(resourceName, "pwd", password),
					resource.TestCheckResourceAttr(resourceName, "pwd_life_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "pwd_update_required", "no"),
					resource.TestCheckResourceAttr(resourceName, "rbac_string", ""),
					resource.TestCheckResourceAttrSet(resourceName, "unix_user_id"),
				),
			},
			{
				Config: CreateAccLocalUserConfigWithOptionalValues(rName, password1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_local_user"),
					resource.TestCheckResourceAttr(resourceName, "account_status", "inactive"),
					resource.TestCheckResourceAttr(resourceName, "cert_attribute", "test_cert_attribute"),
					resource.TestCheckResourceAttr(resourceName, "clear_pwd_history", "yes"),
					resource.TestCheckResourceAttr(resourceName, "email", "test@email.com"),
					resource.TestCheckResourceAttr(resourceName, "expiration", "2030-12-12T00:00:00.000+00:00"),
					resource.TestCheckResourceAttr(resourceName, "expires", "yes"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "test_first_name"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "test_last_name"),
					resource.TestCheckResourceAttr(resourceName, "otpenable", "yes"),
					resource.TestCheckResourceAttr(resourceName, "otpkey", "667PNG5QJW2VAVQA"),
					resource.TestCheckResourceAttr(resourceName, "phone", "1234567890"),
					resource.TestCheckResourceAttr(resourceName, "pwd", password1),
					resource.TestCheckResourceAttr(resourceName, "pwd_life_time", "3650"),
					resource.TestCheckResourceAttr(resourceName, "pwd_update_required", "yes"),
					resource.TestCheckResourceAttr(resourceName, "rbac_string", "test_rbac_string"),
					resource.TestCheckResourceAttrSet(resourceName, "unix_user_id"),
					testAccCheckAciLocalUserIdEqual(&local_user_default, &local_user_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pwd", "clear_pwd_history", "otpkey"},
			},
			{
				Config: CreateAccLocalUserUpdatedAttr(rName, password, "pwd_life_time", "1825"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_updated),
					resource.TestCheckResourceAttr(resourceName, "pwd_life_time", "1825"),
					testAccCheckAciLocalUserIdEqual(&local_user_default, &local_user_updated),
				),
			},
			{
				Config:      CreateAccLocalUserConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLocalUserRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLocalUserConfigWithRequiredParams(rNameUpdated, password2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLocalUserIdNotEqual(&local_user_default, &local_user_updated),
				),
			},
		},
	})
}

func TestAccAciLocalUser_Update(t *testing.T) {
	var local_user_default models.LocalUser
	var local_user_updated models.LocalUser
	resourceName := "aci_local_user.test"
	rName := makeTestVariable(acctest.RandString(5))
	password := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	passwodUpdated := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLocalUserConfig(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_default),
				),
			},
			{
				Config: CreateAccLocalUserUpdatedAttr(rName, passwodUpdated, "clear_pwd_history", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists(resourceName, &local_user_updated),
					resource.TestCheckResourceAttr(resourceName, "clear_pwd_history", "no"),
					testAccCheckAciLocalUserIdEqual(&local_user_default, &local_user_updated),
				),
			},
			{
				Config: CreateAccLocalUserConfig(rName, passwodUpdated),
			},
		},
	})
}

func TestAccAciLocalUser_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	password1 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password2 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password3 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password4 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	password5 := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	passwordWithoutSpecialChars := fmt.Sprintf("%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(2, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLocalUserConfig(rName, password1),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password1, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password1, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password1, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password1, "account_status", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password2, "cert_attribute", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password2, "clear_pwd_history", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password2, "email", randomValue),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password2, "expiration", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password2, "expires", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password3, "first_name", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password3, "last_name", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password3, "otpenable", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password3, "otpkey", randomValue),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password3, "phone", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttrPassword(rName, randomValue),
				ExpectError: regexp.MustCompile(`password must be minimum 8 characters`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttrPassword(rName, acctest.RandString(8)),
				ExpectError: regexp.MustCompile(`password strength check`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttrPassword(rName, passwordWithoutSpecialChars),
				ExpectError: regexp.MustCompile(`password strength check`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password4, "pwd_life_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password4, "pwd_life_time", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password4, "pwd_life_time", "3651"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password4, "pwd_update_required", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLocalUserUpdatedAttr(rName, password4, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLocalUserConfig(rName, password5),
			},
		},
	})
}

func TestAccAciLocalUser_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	password := fmt.Sprintf("%s%s%s", acctest.RandStringFromCharSet(6, "abcdefghjiklmnopqrstuvwxyz"), acctest.RandStringFromCharSet(1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), acctest.RandStringFromCharSet(1, "@#$%^&*()-_!`~+="))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLocalUserConfigMultiple(rName, password),
			},
		},
	})
}

func testAccCheckAciLocalUserExists(name string, local_user *models.LocalUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LocalUser %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LocalUser dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		local_userFound := models.LocalUserFromContainer(cont)
		if local_userFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LocalUser %s not found", rs.Primary.ID)
		}
		*local_user = *local_userFound
		return nil
	}
}

func testAccCheckAciLocalUserDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing local_user destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_local_user" {
			cont, err := client.Get(rs.Primary.ID)
			local_user := models.LocalUserFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LocalUser %s Still exists", local_user.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLocalUserIdEqual(m1, m2 *models.LocalUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("local_user DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLocalUserIdNotEqual(m1, m2 *models.LocalUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("local_user DNs are equal")
		}
		return nil
	}
}

func CreateLocalUserWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing local_user creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_local_user" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLocalUserConfigWithRequiredParams(rName, pwd string) string {
	fmt.Println("=== STEP  testing local_user creation with updated name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}
	`, rName, pwd)
	return resource
}
func CreateAccLocalUserConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing local_user creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLocalUserConfig(rName, pwd string) string {
	fmt.Println("=== STEP  testing local_user creation with required arguments and valid pwd only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}
	`, rName, pwd)
	return resource
}

func CreateAccLocalUserConfigMultiple(rName, pwd string) string {
	fmt.Println("=== STEP  testing multiple local_user creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s_${count.index}"
		pwd = "%s${count.index}"
		count = 5
	}
	`, rName, pwd)
	return resource
}

func CreateAccLocalUserConfigWithOptionalValues(rName, pwd string) string {
	fmt.Println("=== STEP  Basic: testing local_user creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_local_user"
		account_status = "inactive"
		cert_attribute = "test_cert_attribute"
		clear_pwd_history = "yes"
		email = "test@email.com"
		expiration = "2030-12-12T00:00:00.000+00:00"
		expires = "yes"
		first_name = "test_first_name"
		last_name = "test_last_name"
		otpenable = "yes"
		otpkey = "667PNG5QJW2VAVQA"
		phone = "1234567890"
		pwd = "%s"
		pwd_life_time = "3650"
		pwd_update_required = "yes"
		rbac_string = "test_rbac_string"
	}
	`, rName, pwd)

	return resource
}

func CreateAccLocalUserRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing local_user updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_local_user"
		account_status = "blocked"
		cert_attribute = ""
		clear_pwd_history = "yes"
		email = ""
		expiration = ""
		expires = "yes"
		first_name = ""
		last_name = ""
		otpenable = "yes"
		otpkey = ""
		phone = ""
		pwd = ""
		pwd_life_time = "1"
		pwd_update_required = "yes"
		rbac_string = ""
		restricted_rbac_user = "yes"
		
	}
	`)

	return resource
}

func CreateAccLocalUserUpdatedAttrPassword(rName, password string) string {
	fmt.Printf("=== STEP  testing local_user attribute: pwd = %s \n", password)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
	}
	`, rName, password)
	return resource
}

func CreateAccLocalUserUpdatedAttr(rName, password, attribute, value string) string {
	fmt.Printf("=== STEP  testing local_user attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		pwd = "%s"
		%s = "%s"
	}
	`, rName, password, attribute, value)
	return resource
}

func CreateAccLocalUserUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing local_user attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
