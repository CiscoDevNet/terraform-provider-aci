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

func TestAccAciL3outPathAttachmentSecondaryIp_Basic(t *testing.T) {
	var l3out_path_attachment_secondary_ip_default models.L3outPathAttachmentSecondaryIp
	var l3out_path_attachment_secondary_ip_updated models.L3outPathAttachmentSecondaryIp
	resourceName := "aci_l3out_path_attachment_secondary_ip.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("10.0.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.1.0.0/16")
	addrUpdated = fmt.Sprintf("%s/16", addrUpdated)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3outPathAttachmentSecondaryIpWithoutRequired(rName, pathEp1, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outPathAttachmentSecondaryIpWithoutRequired(rName, pathEp1, addr, "l3out_path_attachment_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfig(rName, pathEp1, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists(resourceName, &l3out_path_attachment_secondary_ip_default),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "l3out_path_attachment_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", rName, rName, rName, rName, pathEp1)),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfigWithOptionalValues(rName, pathEp1, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists(resourceName, &l3out_path_attachment_secondary_ip_updated),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_path_attachment_secondary_ip"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "l3out_path_attachment_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", rName, rName, rName, rName, pathEp1)),
					testAccCheckAciL3outPathAttachmentSecondaryIpIdEqual(&l3out_path_attachment_secondary_ip_default, &l3out_path_attachment_secondary_ip_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpWithInvalidIP(rName, pathEp1),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfigWithRequiredParams(rName, pathEp1, addrUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists(resourceName, &l3out_path_attachment_secondary_ip_updated),
					resource.TestCheckResourceAttr(resourceName, "addr", addrUpdated),
					resource.TestCheckResourceAttr(resourceName, "l3out_path_attachment_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", rName, rName, rName, rName, pathEp1)),
					testAccCheckAciL3outPathAttachmentSecondaryIpIdNotEqual(&l3out_path_attachment_secondary_ip_default, &l3out_path_attachment_secondary_ip_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfig(rName, pathEp1, addr),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfigWithRequiredParams(rNameUpdated, pathEp1, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists(resourceName, &l3out_path_attachment_secondary_ip_updated),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "l3out_path_attachment_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated, pathEp1)),
					testAccCheckAciL3outPathAttachmentSecondaryIpIdNotEqual(&l3out_path_attachment_secondary_ip_default, &l3out_path_attachment_secondary_ip_updated),
				),
			},
		},
	})
}

func TestAccAciL3outPathAttachmentSecondaryIp_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	addr, _ := acctest.RandIpAddress("10.2.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfig(rName, pathEp1, addr),
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpConfigWithInvalidParentDn(rName, addr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(rName, pathEp1, addr, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(rName, pathEp1, addr, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(rName, pathEp1, addr, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(rName, pathEp1, addr, "ipv6_dad", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(rName, pathEp1, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfig(rName, pathEp1, addr),
			},
		},
	})
}

func TestAccAciL3outPathAttachmentSecondaryIp_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfigMultiple(rName, pathEp1),
			},
		},
	})
}

func testAccCheckAciL3outPathAttachmentSecondaryIpExists(name string, l3out_path_attachment_secondary_ip *models.L3outPathAttachmentSecondaryIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Path Attachment Secondary Ip %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Path Attachment Secondary Ip dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_path_attachment_secondary_ipFound := models.L3outPathAttachmentSecondaryIpFromContainer(cont)
		if l3out_path_attachment_secondary_ipFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Path Attachment Secondary Ip %s not found", rs.Primary.ID)
		}
		*l3out_path_attachment_secondary_ip = *l3out_path_attachment_secondary_ipFound
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentSecondaryIpDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_path_attachment_secondary_ip" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_path_attachment_secondary_ip := models.L3outPathAttachmentSecondaryIpFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Path Attachment Secondary Ip %s Still exists", l3out_path_attachment_secondary_ip.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outPathAttachmentSecondaryIpIdEqual(m1, m2 *models.L3outPathAttachmentSecondaryIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_path_attachment_secondary_ip DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentSecondaryIpIdNotEqual(m1, m2 *models.L3outPathAttachmentSecondaryIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_path_attachment_secondary_ip DNs are equal")
		}
		return nil
	}
}

func CreateL3outPathAttachmentSecondaryIpWithoutRequired(rName, tdn, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment_secondary_ip creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}
	`
	switch attrName {
	case "addr":
		rBlock += `
	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
	#	addr  = "%s"
	}
		`
	case "l3out_path_attachment_dn":
		rBlock += `
	resource "aci_l3out_path_attachment_secondary_ip" "test" {
	#	l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, rName, tdn, addr)
}

func CreateAccL3outPathAttachmentSecondaryIpConfigWithRequiredParams(name, tdn, addr string) string {
	fmt.Printf("=== STEP  testing l3out_path_attachment_secondary_ip creation with parent resource name %s, tdn %s and addr %s\n", name, tdn, addr)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}
	`, name, name, name, name, tdn, addr)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpWithInvalidIP(name, tdn string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip creation with invalid ip")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		addr  = "%s"
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
	}
	`, name, name, name, name, tdn, name)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpConfigWithInvalidParentDn(name, addr string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip creation with invalid parent dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_tenant.test.id
		addr  = "%s"
	}
	`, name, addr)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpConfig(name, tdn, addr string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}
	`, name, name, name, name, tdn, addr)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpConfigMultiple(name, tdn string) string {
	fmt.Println("=== STEP  testing multiple l3out_path_attachment_secondary_ip creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "10.3.0.${count.index+1}/16"
		count = 5
	}
	`, name, name, name, name, tdn)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpConfigWithOptionalValues(name, tdn, addr string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment_secondary_ip creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_path_attachment_secondary_ip"
		ipv6_dad = "disabled"
		
	}
	`, name, name, name, name, tdn, addr)

	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment_secondary_ip updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_path_attachment_secondary_ip"
		ipv6_dad = "disabled"
		
	}
	`)

	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpUpdatedAttr(name, tdn, addr, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_path_attachment_secondary_ip attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
		%s = "%s"
	}
	`, name, name, name, name, tdn, addr, attribute, value)
	return resource
}
