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

func TestAccAciL3outPathAttachment_Basic(t *testing.T) {
	var l3out_path_attachment_default models.L3outPathAttachment
	var l3out_path_attachment_updated models.L3outPathAttachment
	resourceName := "aci_l3out_path_attachment.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outPathAttachmentWithoutRequired(rName, rName, rName, rName, pathEp1, "ext-svi", "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outPathAttachmentWithoutRequired(rName, rName, rName, rName, pathEp1, "ext-svi", "target_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outPathAttachmentWithoutRequired(rName, rName, rName, rName, pathEp1, "ext-svi", "if_inst_t"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp1, "ext-svi"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "addr", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "autostate", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "encap_scope", "local"),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "::"),
					resource.TestCheckResourceAttr(resourceName, "mac", "00:22:BD:F8:19:FF"),
					resource.TestCheckResourceAttr(resourceName, "mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithOptionalValues(rName, rName, rName, rName, pathEp1, "ext-svi"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "addr", "10.20.30.40/16"),
					resource.TestCheckResourceAttr(resourceName, "autostate", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "encap_scope", "ctx"),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "fe80::1"),
					resource.TestCheckResourceAttr(resourceName, "mac", "00:22:BD:F8:19:F0"),
					resource.TestCheckResourceAttr(resourceName, "mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "576"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccL3outPathAttachmentRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, rName, "ext-svi"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, pathEp1, rName),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, pathEp1, "l3-port"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "l3-port"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, pathEp1, "sub-interface"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "sub-interface"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, pathEp1, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "unspecified"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp1, "ext-svi"),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithRequiredParams(rNameUpdated, pathEp1, "ext-svi"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp1),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					testAccCheckAciL3outPathAttachmentIdNotEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp1, "ext-svi"),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigWithRequiredParams(rNameUpdated, pathEp2, "ext-svi"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "target_dn", pathEp2),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					testAccCheckAciL3outPathAttachmentIdNotEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
		},
	})
}

func TestAccAciL3outPathAttachment_Update(t *testing.T) {
	var l3out_path_attachment_default models.L3outPathAttachment
	var l3out_path_attachment_updated models.L3outPathAttachment
	resourceName := "aci_l3out_path_attachment.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp3, "ext-svi"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_default),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "mode", "untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "untagged"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "mtu", "9216"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "9216"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "mtu", "4896"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "4896"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "CS0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp3, "ext-svi", "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists(resourceName, &l3out_path_attachment_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciL3outPathAttachmentIdEqual(&l3out_path_attachment_default, &l3out_path_attachment_updated),
				),
			},
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp3, "ext-svi"),
			},
		},
	})
}

func TestAccAciL3outPathAttachment_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp4, "ext-svi"),
			},
			{
				Config:      CreateAccL3outPathAttachmentWithInValidParentDn(rName, pathEp4),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "autostate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "encap_scope", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "ipv6_dad", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "ll_addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "mac", acctest.RandStringFromCharSet(5, "ghijklmnopqrstuvwxyz")),
				ExpectError: regexp.MustCompile(`invalid MAC format`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "mtu", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "mtu", "575"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "mtu", "9217"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outPathAttachmentUpdatedAttr(rName, rName, rName, rName, pathEp4, "ext-svi", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outPathAttachmentConfig(rName, rName, rName, rName, pathEp4, "ext-svi"),
			},
		},
	})
}

func testAccCheckAciL3outPathAttachmentExists(name string, l3out_path_attachment *models.L3outPathAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Path Attachment %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Path Attachment dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_path_attachmentFound := models.L3outPathAttachmentFromContainer(cont)
		if l3out_path_attachmentFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Path Attachment %s not found", rs.Primary.ID)
		}
		*l3out_path_attachment = *l3out_path_attachmentFound
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_path_attachment destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_path_attachment" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_path_attachment := models.L3outPathAttachmentFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Path Attachment %s Still exists", l3out_path_attachment.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outPathAttachmentIdEqual(m1, m2 *models.L3outPathAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_path_attachment DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentIdNotEqual(m1, m2 *models.L3outPathAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_path_attachment DNs are equal")
		}
		return nil
	}
}

func CreateL3outPathAttachmentWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment creation without ", attrName)
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
	
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	resource "aci_l3out_path_attachment" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "%s"
	}
		`
	case "target_dn":
		rBlock += `
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	#	target_dn  = "%s"
		if_inst_t = "%s"
	}
		`
	case "if_inst_t":
		rBlock += `
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
	#	if_inst_t = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType)
}

func CreateAccL3outPathAttachmentConfigWithUpdatedIfInstT(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, ifInstT string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment creation with required arguments only")
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
		if_inst_t = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, ifInstT)
	return resource
}
func CreateAccL3outPathAttachmentConfigWithRequiredParams(rName, tdn, instType string) string {
	fmt.Printf("=== STEP  testing l3out_path_attachment creation with parent resource name %s, target_dn %s and if_inst_t %s\n", rName, tdn, instType)
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
		if_inst_t = "%s"
	}
	`, rName, rName, rName, rName, tdn, instType)
	return resource
}

func CreateAccL3outPathAttachmentConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment creation with required arguments only")
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
		if_inst_t = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType)
	return resource
}

func CreateAccL3outPathAttachmentWithInValidParentDn(rName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_path_attachment creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_tenant.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"	
	}
	`, rName, tDn)
	return resource
}

func CreateAccL3outPathAttachmentConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment creation with optional parameters")
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
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}"
		target_dn  = "%s"
		if_inst_t = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		addr = "10.20.30.40/16"
		autostate = "enabled"
		encap_scope = "ctx"
		ipv6_dad = "disabled"
		ll_addr = "fe80::1"
		mac = "00:22:BD:F8:19:F0"
		mode = "native"
		mtu = "576"
		target_dscp = "AF11"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType)

	return resource
}

func CreateAccL3outPathAttachmentRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_path_attachment" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		addr = ""
		autostate = "enabled"
		encap = ""
		encap_scope = "ctx"
		ipv6_dad = "disabled"
		ll_addr = ""
		mac = ""
		mode = "native"
		mtu = "577"
		target_dscp = "1"
	}
	`)

	return resource
}

func CreateAccL3outPathAttachmentUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_path_attachment attribute: %s = %s \n", attribute, value)
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
		if_inst_t = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, instType, attribute, value)
	return resource
}
