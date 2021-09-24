package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRemotePathofaFile_Basic(t *testing.T) {
	var remote_pathofa_file models.RemotePathofaFile
	description := "remote_pathofa_file created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRemotePathofaFileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists("aci_file_remote_path.test", &remote_pathofa_file),
					testAccCheckAciRemotePathofaFileAttributes(description, &remote_pathofa_file),
				),
			},
		},
	})
}

func TestAccAciRemotePathofaFile_Update(t *testing.T) {
	var remote_pathofa_file models.RemotePathofaFile
	description := "remote_pathofa_file created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRemotePathofaFileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists("aci_file_remote_path.test", &remote_pathofa_file),
					testAccCheckAciRemotePathofaFileAttributes(description, &remote_pathofa_file),
				),
			},
			{
				Config: testAccCheckAciRemotePathofaFileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists("aci_file_remote_path.test", &remote_pathofa_file),
					testAccCheckAciRemotePathofaFileAttributes(description, &remote_pathofa_file),
				),
			},
		},
	})
}

func testAccCheckAciRemotePathofaFileConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_file_remote_path" "test" {
		name 		= "test"
		description = "%s"
		annotation = "test_annotation"
		name_alias = "test_alias"
		auth_type = "usePassword"
		user_passwd = "password"
		user_name = "test_user"
		remote_port = "21"
		protocol = "sftp"
		host = "cisco.com"
	}
	`, description)
}

func testAccCheckAciRemotePathofaFileExists(name string, remote_pathofa_file *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Remote Path of a File %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Remote Path of a File dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		remote_pathofa_fileFound := models.RemotePathofaFileFromContainer(cont)
		if remote_pathofa_fileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Remote Path of a File %s not found", rs.Primary.ID)
		}
		*remote_pathofa_file = *remote_pathofa_fileFound
		return nil
	}
}

func testAccCheckAciRemotePathofaFileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_file_remote_path" {
			cont, err := client.Get(rs.Primary.ID)
			remote_pathofa_file := models.RemotePathofaFileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Remote Path of a File %s Still exists", remote_pathofa_file.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRemotePathofaFileAttributes(description string, remote_pathofa_file *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(remote_pathofa_file.DistinguishedName) {
			return fmt.Errorf("Bad file_remote_path %s", GetMOName(remote_pathofa_file.DistinguishedName))
		}

		if description != remote_pathofa_file.Description {
			return fmt.Errorf("Bad file_remote_path Description %s", remote_pathofa_file.Description)
		}

		if "test_annotation" != remote_pathofa_file.Annotation {
			return fmt.Errorf("Bad file_remote_path Annotation %s", remote_pathofa_file.Annotation)
		}

		if "test_alias" != remote_pathofa_file.NameAlias {
			return fmt.Errorf("Bad file_remote_path NameAlias %s", remote_pathofa_file.NameAlias)
		}

		if "usePassword" != remote_pathofa_file.AuthType {
			return fmt.Errorf("Bad file_remote_path AuthType %s", remote_pathofa_file.AuthType)
		}

		if "test_user" != remote_pathofa_file.UserName {
			return fmt.Errorf("Bad file_remote_path UserName %s", remote_pathofa_file.UserName)
		}

		if "21" != remote_pathofa_file.RemotePort {
			return fmt.Errorf("Bad file_remote_path RemotePort %s", remote_pathofa_file.RemotePort)
		}

		if "sftp" != remote_pathofa_file.Protocol {
			return fmt.Errorf("Bad file_remote_path Protocol %s", remote_pathofa_file.Protocol)
		}

		if "cisco.com" != remote_pathofa_file.Host {
			return fmt.Errorf("Bad file_remote_path Host %s", remote_pathofa_file.Host)
		}
		return nil
	}
}
