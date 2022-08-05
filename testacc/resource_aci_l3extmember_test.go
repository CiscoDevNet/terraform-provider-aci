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

func TestAccAciL3outVPCMember_Basic(t *testing.T) {
	var l3out_vpc_member_default models.L3outVPCMember
	var l3out_vpc_member_updated models.L3outVPCMember
	resourceName := "aci_l3out_vpc_member.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outVPCMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outVPCMemberWithoutRequired(rName, "side"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outVPCMemberWithoutRequired(rName, "side"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outVPCMemberConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists(resourceName, &l3out_vpc_member_default),
					resource.TestCheckResourceAttr(resourceName, "leaf_port_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[topology/pod-1/paths-101/pathep-[eth1/1]]", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "side", "A"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "addr", ""),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "::"),
				),
			},
			{
				Config: CreateAccL3outVPCMemberConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists(resourceName, &l3out_vpc_member_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_port_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[topology/pod-1/paths-101/pathep-[eth1/1]]", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "side", "A"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_vpc_member"),
					resource.TestCheckResourceAttr(resourceName, "addr", "10.0.0.5/16"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", ""),
					testAccCheckAciL3outVPCMemberIdEqual(&l3out_vpc_member_default, &l3out_vpc_member_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outVPCMemberRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outVPCMemberConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists(resourceName, &l3out_vpc_member_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_port_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[topology/pod-1/paths-101/pathep-[eth1/1]]", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "side", "B"),
					testAccCheckAciL3outVPCMemberIdNotEqual(&l3out_vpc_member_default, &l3out_vpc_member_updated),
				),
			},
		},
	})
}

func TestAccAciL3outVPCMember_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outVPCMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outVPCMemberConfig(rName),
			},
			{
				Config:      CreateAccL3outVPCMemberUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outVPCMemberUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outVPCMemberUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outVPCMemberUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outVPCMemberConfig(rName),
			},
		},
	})
}

func testAccCheckAciL3outVPCMemberExists(name string, l3out_vpc_member *models.L3outVPCMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out VPC Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out VPC Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_vpc_memberFound := models.L3outVPCMemberFromContainer(cont)
		if l3out_vpc_memberFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out VPC Member %s not found", rs.Primary.ID)
		}
		*l3out_vpc_member = *l3out_vpc_memberFound
		return nil
	}
}

func testAccCheckAciL3outVPCMemberDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_vpc_member destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_vpc_member" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_vpc_member := models.L3outVPCMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out VPC Member %s Still exists", l3out_vpc_member.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outVPCMemberIdEqual(m1, m2 *models.L3outVPCMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_vpc_member DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outVPCMemberIdNotEqual(m1, m2 *models.L3outVPCMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_vpc_member DNs are equal")
		}
		return nil
	}
}

func CreateL3outVPCMemberWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_vpc_member creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}

	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}

	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	`
	switch attrName {
	case "side":
		rBlock += `
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
	#	side  = "A"
	}
		`
	case "leaf_port_dn":
		rBlock += `
	resource "aci_l3out_vpc_member" "test" {
	#	leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, rName)
}

func CreateAccL3outVPCMemberConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member creation with updated naming arguments")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "B"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outVPCMemberConfig(rName string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outVPCMemberConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_vpc_member creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_vpc_member"
		addr = "10.0.0.5/16"
		ipv6_dad = "disabled"
		ll_addr = ""
	}
	`, rName, rName, rName, rName)

	return resource
}

func CreateAccL3outVPCMemberRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_vpc_member updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_vpc_member" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_vpc_member"
		addr = ""
		ipv6_dad = "disabled"
		ll_addr = ""
	}
	`)

	return resource
}

func CreateAccL3outVPCMemberUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_vpc_member attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
		%s = "%s"
	}
	`, rName, rName, rName, rName, attribute, value)
	return resource
}
