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

func TestAccAciL3outHsrpInterfaceGroup_Basic(t *testing.T) {
	var l3out_hsrp_interface_group_default models.HSRPGroupProfile
	var l3out_hsrp_interface_group_updated models.HSRPGroupProfile
	resourceName := "aci_l3out_hsrp_interface_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHsrpInterfaceGroupWithoutRequired(rName, rName, rName, rName, rName, "l3out_hsrp_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outHsrpInterfaceGroupWithoutRequired(rName, rName, rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName, rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "none"),
					resource.TestCheckResourceAttr(resourceName, "group_af", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "ip", "10.20.30.40"),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "admin"),
					resource.TestCheckResourceAttr(resourceName, "mac", "00:00:00:00:00:00"),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", ""),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithOptionalValues(rName, rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_hsrp_interface_group"),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "group_af", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "125"),
					resource.TestCheckResourceAttr(resourceName, "group_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "ip", "10.20.30.50"),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "admin"),
					resource.TestCheckResourceAttr(resourceName, "mac", "02:10:45:00:00:56"),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", ""),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName, rName, rName, rName, rName),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "l3out_hsrp_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_Update(t *testing.T) {
	var l3out_hsrp_interface_group_default models.HSRPGroupProfile
	var l3out_hsrp_interface_group_updated models.HSRPGroupProfile
	resourceName := "aci_l3out_hsrp_interface_group.test"
	rName1 := makeTestVariable(acctest.RandString(5))
	rName2 := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithIPObtainModeLearn(rName1, rName1, rName1, rName1, rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "learn"),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "GroupName-Conflicts-Other-Group"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupName-Conflicts-Other-Group"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "GroupVIP-Conflicts-Other-Group"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "GroupVIP-Conflicts-Other-Group"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "Multiple-Version-On-Interface"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Multiple-Version-On-Interface"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "Secondary-vip-conflicts-if-ip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Secondary-vip-conflicts-if-ip"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "Secondary-vip-subnet-mismatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "Secondary-vip-subnet-mismatch"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "group-vip-conflicts-if-ip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "group-vip-conflicts-if-ip"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", "group-vip-subnet-mismatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "group-vip-subnet-mismatch"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_id", "2025"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "group_id", "2025"),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithIPObtainModeAuto(rName2, rName2, rName2, rName2, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "auto"),
					resource.TestCheckResourceAttr(resourceName, "group_af", "ipv6"),
				),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_Negative(t *testing.T) {
	rName1 := makeTestVariable(acctest.RandString(5))
	rName2 := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName1, rName1, rName1, rName1, rName1),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupWithInValidParentDn(rName1),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class hsrpGroupP (.)+`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "config_issues", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of(.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_af", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of(.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "ip_obtain_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of(.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_name", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`property groupName of (.)+ failed validation`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "ip", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_id", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_id", "256"),
				ExpectError: regexp.MustCompile(`Invalid Configuration HSRP V1 group id range is 0-255`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, "group_id", "4096"),
				ExpectError: regexp.MustCompile(`Property (.)+ of (.)+ is out of range`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(rName1, rName1, rName1, rName1, rName1, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+ is not expected here.`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupForAdmin(rName2, rName2, rName2, rName1, rName1),
				ExpectError: regexp.MustCompile(`Invalid Configuration VIP Address cannot be empty with Admin Mode`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupForAuto(rName2, rName2, rName2, rName1, rName1),
				ExpectError: regexp.MustCompile(`Invalid Configuration VIP configuration should be NULL, if learn/auto configuration is enabled.`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupForLearn(rName2, rName2, rName2, rName2, rName2),
				ExpectError: regexp.MustCompile(`Invalid Configuration VIP configuration should be NULL, if learn/auto configuration is enabled.`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupForAdminAF6(rName2, rName2, rName2, rName2, rName2),
				ExpectError: regexp.MustCompile(`Invalid Configuration HSRP V1 group is V4 only group`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName1, rName1, rName1, rName1, rName1),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_RelationParameter(t *testing.T) {
	relName1 := makeTestVariable(acctest.RandString(5))
	relName2 := makeTestVariable(acctest.RandString(5))
	var l3out_hsrp_interface_group_default models.HSRPGroupProfile
	var l3out_hsrp_interface_group_updated models.HSRPGroupProfile
	resourceName := "aci_l3out_hsrp_interface_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName, rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", ""),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigRel(rName, rName, rName, rName, rName, relName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", fmt.Sprintf("uni/tn-%s/hsrpGroupPol-%s", rName, relName1)),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigRel(rName, rName, rName, rName, rName, relName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", fmt.Sprintf("uni/tn-%s/hsrpGroupPol-%s", rName, relName2)),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(rName, rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "relation_hsrp_rs_group_pol", ""),
					testAccCheckAciL3outHsrpInterfaceGroupIdEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigs(rName),
			},
		},
	})
}

func testAccCheckAciL3outHsrpInterfaceGroupExists(name string, l3out_hsrp_interface_group *models.HSRPGroupProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Hsrp Interface Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Hsrp Interface Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_interface_groupFound := models.HSRPGroupProfileFromContainer(cont)
		if l3out_hsrp_interface_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Hsrp Interface Group %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_interface_group = *l3out_hsrp_interface_groupFound
		return nil
	}
}

func testAccCheckAciL3outHsrpInterfaceGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_hsrp_interface_group" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_hsrp_interface_group := models.HSRPGroupProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Hsrp Interface Group %s Still exists", l3out_hsrp_interface_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outHsrpInterfaceGroupIdEqual(m1, m2 *models.HSRPGroupProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(m1, m2 *models.HSRPGroupProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_group DNs are equal")
		}
		return nil
	}
}

func CreateL3outHsrpInterfaceGroupWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation without ", attrName)
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	`
	switch attrName {
	case "l3out_hsrp_interface_profile_dn":
		rBlock += `
	resource "aci_l3out_hsrp_interface_group" "test" {
	#	l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
}

func CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_group creation with parent resources name %s and resource name %s \n", prName, rName)
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip = "10.20.30.50"
	}
	`, prName, prName, prName, prName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfigs(rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group muliple creation")
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test1" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip = "10.20.30.50"
	}

	resource "aci_l3out_hsrp_interface_group" "test2" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip = "10.20.30.50"
	}

	resource "aci_l3out_hsrp_interface_group" "test3" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip = "10.20.30.50"
	}
	`, rName, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfigRel(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, relName string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_group creation with relation resource name %s\n", relName)
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
		relation_hsrp_rs_group_pol = aci_hsrp_group_policy.test.id
	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	  }
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, relName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group creation with required arguments and ip")
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
    	ip = "10.20.30.40"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_hsrp_interface_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation with optional parameters")
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_interface_group"
		config_issues = "GroupMac-Conflicts-Other-Group"
		group_af = "ipv4"
		group_id = "125"
		group_name = "test"
		ip = "10.20.30.50"
		ip_obtain_mode = "admin"
		mac = "02:10:45:00:00:56"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfigWithIPObtainModeAuto(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation with  ip_obtain_mode = auto and group_af = ipv6")
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "auto"
		group_af = "ipv6"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupConfigWithIPObtainModeLearn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation with  ip_obtain_mode = learn and required arguments")
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		version = "v2"
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "learn"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3out_hsrp_interface_group" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_l3out_hsrp_interface_group"
		config_issues = "GroupMac-Conflicts-Other-Group"
		group_af = "ipv6"
		group_id = "1"
		group_name = "test"
		ip_obtain_mode = "auto"
	}
	`)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupForLearn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group when ip_obtain_mode is learn and ip provided")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "learn"
		ip = "1.2.3.4"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupForAuto(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group when ip_obtain_mode is admin and ip provided")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "auto"
		ip = "1.2.3.4"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupForAdminAF6(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group when ip_obtain_mode is admin and group_af is ipv6")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "admin"
		group_af = "ipv6"
		ip = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupForAdmin(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group when ip_obtain_mode is admin and no ip provided")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		ip_obtain_mode = "admin"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_group attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
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
    	logical_node_profile_dn = aci_logical_node_profile.test.id
    	name = "%s"
  	}
	
	resource "aci_l3out_hsrp_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn  = aci_l3out_hsrp_interface_profile.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, attribute, value)
	return resource
}
