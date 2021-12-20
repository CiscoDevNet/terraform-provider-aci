package acctest

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

func TestAccAciEndpointSecurityGroup_Basic(t *testing.T) {
	var esg_default models.EndpointSecurityGroup
	var esg_updated models.EndpointSecurityGroup
	resourceName := "aci_endpoint_security_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccEndpointSecurityGroupWithoutAP(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "unenforced"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_name_alias"),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "enforced"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "0"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config:      CreateAccEndpointSecurityGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of esg-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					testAccCheckAciEndpointSecurityGroupIdNotEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfig(rName),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", prOther, prOther)),
					testAccCheckAciEndpointSecurityGroupIdNotEqual(&esg_default, &esg_updated),
				),
			},
		},
	})
}

func TestAccAciEndpointSecurityGroup_Update(t *testing.T) {
	var esg_default models.EndpointSecurityGroup
	var esg_updated models.EndpointSecurityGroup
	resourceName := "aci_endpoint_security_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointSecurityGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_default),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists(resourceName, &esg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciEndpointSecurityGroupIdEqual(&esg_default, &esg_updated),
				),
			},
		},
	})
}

func TestAccAciEndpointSecurityGroup_NegativeCases(t *testing.T) {
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointSecurityGroupConfig(rName),
			},
			{
				Config:      CreateAccEndpointSecurityGroupWithInValidAPDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvESg (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected flood_on_encap to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "pc_enf_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected pc_enf_pref to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected pref_gr_memb to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointSecurityGroupConfig(rName),
			},
		},
	})
}

func TestAccAciEndpointSecurityGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointSecurityGroupsConfig(rName),
			},
		},
	})
}

func CreateAccEndpointSecurityGroupsConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group multiple creation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test1" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_endpoint_security_group" "test2" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_endpoint_security_group" "test3" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func testAccCheckAciEndpointSecurityGroupIdNotEqual(esg1, esg2 *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if esg1.DistinguishedName == esg2.DistinguishedName {
			return fmt.Errorf("Endpoint Security Group DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupIdEqual(esg1, esg2 *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if esg1.DistinguishedName != esg2.DistinguishedName {
			return fmt.Errorf("Endpoint Security Group DNs are not equal")
		}
		return nil
	}
}

func CreateAccEndpointSecurityGroupWithInValidAPDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing endpoint_security_group creation with invalid application_profile_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	  }
	  
	  resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_tenant.test.id
	  }
	`, rName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing endpoint_security_group creation with parente resource name %s and name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	  
	  resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
		annotation = "test_annotation"
		description = "test_description"
		flood_on_encap = "enabled"
		match_t = "All"
		name_alias = "test_name_alias"
		pc_enf_pref = "enforced"
		pref_gr_memb = "include"
		prio = "level1"
	  }
	`, prName, prName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group updation without required fields")
	resource := fmt.Sprintln(`
	resource "aci_endpoint_security_group" "test" {
		annotation = "tag"
		description = "test_description"
		flood_on_encap = "enabled"
		match_t = "All"
		name_alias = "test_name_alias"
		pc_enf_pref = "enforced"
		pref_gr_memb = "include"
		prio = "level1"
	  }
	`)
	return resource
}

func CreateAccEndpointSecurityGroupConfigWithUpdatedAttr(rName, key, value string) string {
	fmt.Printf("=== STEP  testing endpoint_security_group creation with %s = %s\n", key, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	  
	  resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
		%s = "%s"
	  }
	`, rName, rName, rName, key, value)
	return resource
}

func CreateAccEndpointSecurityGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  testing endpoint_security_group creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	  }
	  
	  resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
		annotation = "test_annotation"
		description = "test_description"
		flood_on_encap = "enabled"
		match_t = "All"
		name_alias = "test_name_alias"
		pc_enf_pref = "enforced"
		pref_gr_memb = "include"
		prio = "level1"
	  }
	`, rName, rName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing endpoint_security_group creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	 }
	  
	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	  
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccEndpointSecurityGroupWithoutAP(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group creation without creating application profile")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_security_group" "test" {
		name = "%s"
	  }
	`, rName)
	return resource
}

func CreateAccEndpointSecurityGroupWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_security_group creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	  
	 resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	 }
	  
	 resource "aci_endpoint_security_group" "test" {
		application_profile_dn = aci_application_profile.test.id
	 }
	`, rName, rName)
	return resource
}

func testAccCheckAciEndpointSecurityGroupExists(name string, endpoint_security_group *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Security Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Security Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_security_groupFound := models.EndpointSecurityGroupFromContainer(cont)
		if endpoint_security_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Security Group %s not found", rs.Primary.ID)
		}
		*endpoint_security_group = *endpoint_security_groupFound
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_security_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_security_group" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_security_group := models.EndpointSecurityGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Security Group %s Still exists", endpoint_security_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}
