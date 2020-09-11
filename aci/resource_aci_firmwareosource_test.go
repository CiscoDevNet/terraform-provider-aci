package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciFirmwareDownloadTask_Basic(t *testing.T) {
	var firmware_download_task models.FirmwareDownloadTask
	description := "firmware_download_task created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwareDownloadTaskConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists("aci_firmware_download_task.foofirmware_download_task", &firmware_download_task),
					testAccCheckAciFirmwareDownloadTaskAttributes(description, &firmware_download_task),
				),
			},
			{
				ResourceName:      "aci_firmware_download_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFirmwareDownloadTask_update(t *testing.T) {
	var firmware_download_task models.FirmwareDownloadTask
	description := "firmware_download_task created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwareDownloadTaskConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists("aci_firmware_download_task.foofirmware_download_task", &firmware_download_task),
					testAccCheckAciFirmwareDownloadTaskAttributes(description, &firmware_download_task),
				),
			},
			{
				Config: testAccCheckAciFirmwareDownloadTaskConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists("aci_firmware_download_task.foofirmware_download_task", &firmware_download_task),
					testAccCheckAciFirmwareDownloadTaskAttributes(description, &firmware_download_task),
				),
			},
		},
	})
}

func testAccCheckAciFirmwareDownloadTaskConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_firmware_download_task" "foofirmware_download_task" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  auth_pass  = "password"
		  auth_type  = "usePassword"
		  dnld_task_flip  = "example"
		  identity_private_key_contents  = "example"
		  identity_private_key_passphrase  = "example"
		  identity_public_key_contents  = "example"
		  load_catalog_if_exists_and_newer  = "no"
		  name_alias  = "example"
		  password  = "example"
		  polling_interval  = "example"
		  proto  = "http"
		  url  = "example"
		  user  = "example"
		}
	`, description)
}

func testAccCheckAciFirmwareDownloadTaskExists(name string, firmware_download_task *models.FirmwareDownloadTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Download Task %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Download Task dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_download_taskFound := models.FirmwareDownloadTaskFromContainer(cont)
		if firmware_download_taskFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Download Task %s not found", rs.Primary.ID)
		}
		*firmware_download_task = *firmware_download_taskFound
		return nil
	}
}

func testAccCheckAciFirmwareDownloadTaskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_firmware_download_task" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_download_task := models.FirmwareDownloadTaskFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Download Task %s Still exists", firmware_download_task.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFirmwareDownloadTaskAttributes(description string, firmware_download_task *models.FirmwareDownloadTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != firmware_download_task.Description {
			return fmt.Errorf("Bad firmware_download_task Description %s", firmware_download_task.Description)
		}

		if "example" != firmware_download_task.Name {
			return fmt.Errorf("Bad firmware_download_task name %s", firmware_download_task.Name)
		}

		if "example" != firmware_download_task.Annotation {
			return fmt.Errorf("Bad firmware_download_task annotation %s", firmware_download_task.Annotation)
		}

		if "password" != firmware_download_task.AuthPass {
			return fmt.Errorf("Bad firmware_download_task auth_pass %s", firmware_download_task.AuthPass)
		}

		if "usePassword" != firmware_download_task.AuthType {
			return fmt.Errorf("Bad firmware_download_task auth_type %s", firmware_download_task.AuthType)
		}

		if "example" != firmware_download_task.DnldTaskFlip {
			return fmt.Errorf("Bad firmware_download_task dnld_task_flip %s", firmware_download_task.DnldTaskFlip)
		}

		if "example" != firmware_download_task.IdentityPrivateKeyContents {
			return fmt.Errorf("Bad firmware_download_task identity_private_key_contents %s", firmware_download_task.IdentityPrivateKeyContents)
		}

		if "example" != firmware_download_task.IdentityPrivateKeyPassphrase {
			return fmt.Errorf("Bad firmware_download_task identity_private_key_passphrase %s", firmware_download_task.IdentityPrivateKeyPassphrase)
		}

		if "example" != firmware_download_task.IdentityPublicKeyContents {
			return fmt.Errorf("Bad firmware_download_task identity_public_key_contents %s", firmware_download_task.IdentityPublicKeyContents)
		}

		if "no" != firmware_download_task.LoadCatalogIfExistsAndNewer {
			return fmt.Errorf("Bad firmware_download_task load_catalog_if_exists_and_newer %s", firmware_download_task.LoadCatalogIfExistsAndNewer)
		}

		if "example" != firmware_download_task.NameAlias {
			return fmt.Errorf("Bad firmware_download_task name_alias %s", firmware_download_task.NameAlias)
		}

		if "example" != firmware_download_task.Password {
			return fmt.Errorf("Bad firmware_download_task password %s", firmware_download_task.Password)
		}

		if "example" != firmware_download_task.PollingInterval {
			return fmt.Errorf("Bad firmware_download_task polling_interval %s", firmware_download_task.PollingInterval)
		}

		if "http" != firmware_download_task.Proto {
			return fmt.Errorf("Bad firmware_download_task proto %s", firmware_download_task.Proto)
		}

		if "example" != firmware_download_task.Url {
			return fmt.Errorf("Bad firmware_download_task url %s", firmware_download_task.Url)
		}

		if "example" != firmware_download_task.User {
			return fmt.Errorf("Bad firmware_download_task user %s", firmware_download_task.User)
		}

		return nil
	}
}
