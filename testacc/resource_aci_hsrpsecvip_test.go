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

func TestAccAciL3outHSRPSecondaryVIP_Basic(t *testing.T) {
	var l3out_hsrp_secondary_vip_default models.L3outHSRPSecondaryVIP
	var l3out_hsrp_secondary_vip_updated models.L3outHSRPSecondaryVIP
	resourceName := "aci_l3out_hsrp_secondary_vip.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.1.0.0/16")
	ipUpdated, _ := acctest.RandIpAddress("10.2.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHSRPSecondaryVIPWithoutRequired(rName, rName, rName, rName, rName, ip, "l3out_hsrp_interface_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outHSRPSecondaryVIPWithoutRequired(rName, rName, rName, rName, rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfig(rName, rName, rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_default),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_group_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", rName, rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "none"),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfigWithOptionalValues(rName, rName, rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_group_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", rName, rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_hsrp_secondary_vip"),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupMac-Conflicts-Other-Group"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPWithInavalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfigWithRequiredParams(rNameUpdated, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_group_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					testAccCheckAciL3outHSRPSecondaryVIPIdNotEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfig(rName, rName, rName, rName, rName, ip),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfigWithRequiredParams(rName, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_group_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", rName, rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					testAccCheckAciL3outHSRPSecondaryVIPIdNotEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
		},
	})
}

func TestAccAciL3outHSRPSecondaryVIP_Update(t *testing.T) {
	var l3out_hsrp_secondary_vip_default models.L3outHSRPSecondaryVIP
	var l3out_hsrp_secondary_vip_updated models.L3outHSRPSecondaryVIP
	resourceName := "aci_l3out_hsrp_secondary_vip.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.3.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfig(rName, rName, rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_default),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "GroupName-Conflicts-Other-Group"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupName-Conflicts-Other-Group"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "GroupVIP-Conflicts-Other-Group"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupVIP-Conflicts-Other-Group"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "Multiple-Version-On-Interface"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Multiple-Version-On-Interface"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "Secondary-vip-conflicts-if-ip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Secondary-vip-conflicts-if-ip"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "Secondary-vip-subnet-mismatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Secondary-vip-subnet-mismatch"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "group-vip-conflicts-if-ip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "group-vip-conflicts-if-ip"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", "group-vip-subnet-mismatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists(resourceName, &l3out_hsrp_secondary_vip_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "group-vip-subnet-mismatch"),
					testAccCheckAciL3outHSRPSecondaryVIPIdEqual(&l3out_hsrp_secondary_vip_default, &l3out_hsrp_secondary_vip_updated),
				),
			},
		},
	})
}

func TestAccAciL3outHSRPSecondaryVIP_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.4.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfig(rName, rName, rName, rName, rName, ip),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class hsrpSecVip (.)+`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, "config_issues", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outHSRPSecondaryVIPUpdatedAttr(rName, rName, rName, rName, rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfig(rName, rName, rName, rName, rName, ip),
			},
		},
	})
}

func TestAccAciL3outHSRPSecondaryVIP_MulipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip1, _ := acctest.RandIpAddress("10.5.0.0/16")
	ip2, _ := acctest.RandIpAddress("10.6.0.0/16")
	ip3, _ := acctest.RandIpAddress("10.7.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHSRPSecondaryVIPConfigs(rName, rName, rName, rName, rName, ip1, ip2, ip3),
			},
		},
	})
}

func CreateAccL3outHSRPSecondaryVIPWithInavalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_secondary_vip creation with invalid IP")
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s0"
	}
	`, rName, rName, rName, rName, rName, ip)
	return resource
}

func testAccCheckAciL3outHSRPSecondaryVIPExists(name string, l3out_hsrp_secondary_vip *models.L3outHSRPSecondaryVIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out HSRP Secondary VIP %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out HSRP Secondary VIP dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_secondary_vipFound := models.L3outHSRPSecondaryVIPFromContainer(cont)
		if l3out_hsrp_secondary_vipFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out HSRP Secondary VIP %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_secondary_vip = *l3out_hsrp_secondary_vipFound
		return nil
	}
}

func testAccCheckAciL3outHSRPSecondaryVIPDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_hsrp_secondary_vip" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_hsrp_secondary_vip := models.L3outHSRPSecondaryVIPFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out HSRP Secondary VIP %s Still exists", l3out_hsrp_secondary_vip.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outHSRPSecondaryVIPIdEqual(m1, m2 *models.L3outHSRPSecondaryVIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_secondary_vip DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outHSRPSecondaryVIPIdNotEqual(m1, m2 *models.L3outHSRPSecondaryVIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_secondary_vip DNs are equal")
		}
		return nil
	}
}

func CreateL3outHSRPSecondaryVIPWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_secondary_vip creation without ", attrName)
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	`
	switch attrName {
	case "l3out_hsrp_interface_group_dn":
		rBlock += `
	resource "aci_l3out_hsrp_secondary_vip" "test" {
	#	l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
		`
	case "ip":
		rBlock += `
	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
	#	ip  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)
}

func CreateAccL3outHSRPSecondaryVIPConfigWithRequiredParams(rName, ip string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_secondary_vip creation with parent resources name %s and ip %s\n", rName, ip)
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
	`, rName, rName, rName, rName, rName, ip)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip creation with required arguments only")
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
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}
	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPConfigs(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip1, ip2, ip3 string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_secondary_vip multiple creation")
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
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}
	resource "aci_l3out_hsrp_secondary_vip" "test1" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
	resource "aci_l3out_hsrp_secondary_vip" "test2" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
	resource "aci_l3out_hsrp_secondary_vip" "test3" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip1, ip2, ip3)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_hsrp_secondary_vip creation with invalid parent Dn")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_tenant.test.id
		ip  = "%s"	
	}
	`, rName, ip)
	return resource
}

func CreateAccL3outHSRPSecondaryVIPConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_secondary_vip creation with optional parameters")
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_secondary_vip"
		config_issues = "GroupMac-Conflicts-Other-Group"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip)

	return resource
}

func CreateAccL3outHSRPSecondaryVIPRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_secondary_vip updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3out_hsrp_secondary_vip" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_l3out_hsrp_secondary_vip"
		config_issues = "GroupMac-Conflicts-Other-Group"
	}
	`)

	return resource
}

func CreateAccL3outHSRPSecondaryVIPUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_secondary_vip attribute: %s=%s \n", attribute, value)
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

	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}

	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.test.id
		name = "%s"
		ip_obtain_mode = "learn"
	}

	resource "aci_l3out_hsrp_secondary_vip" "test" {
		l3out_hsrp_interface_group_dn  = aci_l3out_hsrp_interface_group.test.id
		ip  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, hsrpGroupPName, ip, attribute, value)
	return resource
}
