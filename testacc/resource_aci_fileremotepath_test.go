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

func TestAccAciRemotePathofaFile_Basic(t *testing.T) {
	var file_remote_path_default models.RemotePathofaFile
	var file_remote_path_updated models.RemotePathofaFile
	resourceName := "aci_file_remote_path.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	host := "cisco.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRemotePathofaFileWithoutRequired(rName, host, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRemotePathofaFileWithoutRequired(rName, host, "host"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRemotePathofaFileConfig(rName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "usePassword"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "sftp"),
					resource.TestCheckResourceAttr(resourceName, "remote_path", ""),
					resource.TestCheckResourceAttr(resourceName, "remote_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_name", ""),
					resource.TestCheckResourceAttr(resourceName, "user_passwd", "cisco"),
				),
			},
			{
				Config: CreateAccRemotePathofaFileConfigWithOptionalValues(rName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_file_remote_path"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "usePassword"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "scp"),
					resource.TestCheckResourceAttr(resourceName, "remote_port", "18"),

					resource.TestCheckResourceAttr(resourceName, "remote_path", "/example"),

					resource.TestCheckResourceAttr(resourceName, "user_name", "example"),

					resource.TestCheckResourceAttr(resourceName, "user_passwd", "example"),

					testAccCheckAciRemotePathofaFileIdEqual(&file_remote_path_default, &file_remote_path_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"user_passwd",
				},
			},
			{
				Config:      CreateAccRemotePathofaFileConfigUpdatedName(acctest.RandString(65), host),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRemotePathofaFileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccRemotePathofaFileConfigWithRequiredParams(rNameUpdated, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRemotePathofaFileIdNotEqual(&file_remote_path_default, &file_remote_path_updated),
				),
			},
		},
	})
}

func TestAccAciRemotePathofaFile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	host := "cisco.com"

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccRemotePathofaFileWhenUserPasswordIsnotGiven(rName, host),
				ExpectError: regexp.MustCompile(`user_passwd must be set when auth_type is usePassword`),
			},
			{
				Config: CreateAccRemotePathofaFileConfig(rName, host),
			},
			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "auth_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "remote_path", randomValue),
				ExpectError: regexp.MustCompile(`The first character of remote_path should be '/'`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "remote_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "remote_port", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, "remote_port", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, host, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRemotePathofaFileConfig(rName, host),
			},
		},
	})
}

func TestAccAciRemotePathofaFile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	host := "cisco.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccfileRemotePathConfigMultiple(rName, host),
			},
		},
	})
}

func testAccCheckAciRemotePathofaFileExists(name string, file_remote_path *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("File Remote Path %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No File Remote Path dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		file_remote_pathFound := models.RemotePathofaFileFromContainer(cont)
		if file_remote_pathFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("File Remote Path %s not found", rs.Primary.ID)
		}
		*file_remote_path = *file_remote_pathFound
		return nil
	}
}

func testAccCheckAciRemotePathofaFileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing file_remote_path destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_file_remote_path" {
			cont, err := client.Get(rs.Primary.ID)
			file_remote_path := models.RemotePathofaFileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("File Remote Path %s Still exists", file_remote_path.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRemotePathofaFileIdEqual(m1, m2 *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("file_remote_path DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRemotePathofaFileIdNotEqual(m1, m2 *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("file_remote_path DNs are equal")
		}
		return nil
	}
}

func CreateRemotePathofaFileWithoutRequired(rName, host, attrName string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_file_remote_path" "test" {
	
	#	name  = "%s"
		host = "%s"
	}
		`
	case "host":
		rBlock += `
	resource "aci_file_remote_path" "test" {
		name  = "%s"
	#	host = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, host)
}

func CreateAccRemotePathofaFileConfigWithRequiredParams(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"
	}
	`, rName, host)
	return resource
}

func CreateAccRemotePathofaFileConfigUpdatedName(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with invalid name = ", rName, host)
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"	
		user_passwd = "cisco"
	}
	`, rName, host)
	return resource
}

func CreateAccRemotePathofaFileConfig(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"
	}
	`, rName, host)
	return resource
}

func CreateAccfileRemotePathConfigMultiple(rName, host string) string {
	fmt.Println("=== STEP  testing multiple file_remote_path creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s_${count.index}"
		host = "%s"
		user_passwd = "cisco"
		count = 5
	}
	`, rName, host)
	return resource
}

func CreateAccRemotePathofaFileConfigWithOptionalValues(rName, host string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path creation with optional parameters")

	resource := fmt.Sprintf(`
	resource "aci_file_remote_path" "test" {
		name  = "%s"
		host = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_file_remote_path"
		auth_type = "usePassword"
		// identity_private_key_contents = ""
		// identity_private_key_passphrase = ""
		// identity_public_key_contents = ""
		protocol = "scp"
		remote_path = "/example"
		user_name = "example"
		user_passwd = "example"
		remote_port = "18"
	}
	`, rName, host)

	return resource
}

func CreateAccRemotePathofaFileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing file_remote_path updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_file_remote_path" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_file_remote_path"
		auth_type = "usePassword"
		// identity_private_key_contents = ""
		// identity_private_key_passphrase = ""
		// identity_public_key_contents = ""
		protocol = "scp"
		remote_path = "/example"
		user_name = "example"
		user_passwd = "example"
		
	}
	`)

	return resource
}

func CreateAccRemotePathofaFileWhenUserPasswordIsnotGiven(rName, host string) string {
	fmt.Println("=== STEP  testing file_remote_path when auth_type is usePassword and user_passwd is not given")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
		name  = "%s"
		host = "%s"
	}
	`, rName, host)
	return resource
}

func CreateAccRemotePathofaFileUpdatedAttr(rName, host, attribute, value string) string {
	fmt.Printf("=== STEP  testing file_remote_path attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		user_passwd = "cisco"
		%s = "%s"
	}
	`, rName, host, attribute, value)
	return resource
}

func CreateAccRemotePathofaFileUpdatedPasswd(rName, host, attribute, value string) string {
	fmt.Printf("=== STEP  testing file_remote_path attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		host = "%s"
		%s = "%s"
	}
	`, rName, host, attribute, value)
	return resource
}
