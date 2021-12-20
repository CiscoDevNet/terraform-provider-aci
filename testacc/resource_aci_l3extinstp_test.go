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

func TestAccAciExternalNetworkInstanceProfile_Basic(t *testing.T) {
	var external_network_instance_profile_default models.ExternalNetworkInstanceProfile
	var external_network_instance_profile_update models.ExternalNetworkInstanceProfile
	resourceName := "aci_external_network_instance_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccExternalNetworkInstanceProfileWithoutL3Out(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_inst_p_to_nat_mapping_epg", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_inst_p_to_profile.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_inst_p_to_dom_p", ""),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "annotation", "annotation"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "0"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "name_alias"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_inst_p_to_nat_mapping_epg", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_inst_p_to_profile.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l3_inst_p_to_dom_p", ""),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of instP-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rName, rName)),
					testAccCheckAciExternalNetworkInstanceProfileIdNotEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfig(rName),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfigWithParentAndName(rOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", rOther, rOther)),
					testAccCheckAciExternalNetworkInstanceProfileIdNotEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
		},
	})
}

func TestAccAciExternalNetworkInstanceProfile_Update(t *testing.T) {
	var external_network_instance_profile_default models.ExternalNetworkInstanceProfile
	var external_network_instance_profile_update models.ExternalNetworkInstanceProfile
	resourceName := "aci_external_network_instance_profile.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccExternalNetworkInstanceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_default),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "exception_tag", "512"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "512"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF32"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF33"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS4"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF42"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF43"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS5"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS6"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", "CS7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists(resourceName, &external_network_instance_profile_update),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS7"),
					testAccCheckAciExternalNetworkInstanceProfileIdEqual(&external_network_instance_profile_default, &external_network_instance_profile_update),
				),
			},
		},
	})
}

func TestAccAciExternalNetworkInstanceProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccExternalNetworkInstanceProfileConfig(rName),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileConfigWithInvalidL3Out(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extInstP (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "exception_tag", "-1"),
				ExpectError: regexp.MustCompile(`property is out of range`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "exception_tag", "513"),
				ExpectError: regexp.MustCompile(`property is out of range`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "flood_on_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected flood_on_encap to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "pref_gr_memb", randomValue),
				ExpectError: regexp.MustCompile(`expected pref_gr_memb to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccExternalNetworkInstanceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciExternalNetworkInstanceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccExternalNetworkInstanceProfileConfigs(rName),
			},
		},
	})
}

func CreateAccExternalNetworkInstanceProfileConfigWithInvalidL3Out(rName string) string {
	fmt.Println("=== STEP  testing external_network_instance_profile updation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
		l3_outside_dn  = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing external_network_instance_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_external_network_instance_profile" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccExternalNetworkInstanceProfileConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing external_network_instance_profile creation with l3_outside name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_external_network_instance_profile" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
	}
	`, prName, prName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile updation without required fields")
	resource := fmt.Sprintln(`
	resource "aci_external_network_instance_profile" "test" {
		annotation = "tag"
		description = "description"
		exception_tag = "0"
		flood_on_encap = "enabled"
		match_t = "All"
		name_alias = "name_alias"
		pref_gr_memb = "include"
		prio = "level1"
		target_dscp = "CS0"
	}
	`)
	return resource
}

func CreateAccExternalNetworkInstanceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  testing external_network_instance_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_external_network_instance_profile" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		name = "%s"
		annotation = "annotation"
		description = "description"
		exception_tag = "0"
		flood_on_encap = "enabled"
		match_t = "All"
		name_alias = "name_alias"
		pref_gr_memb = "include"
		prio = "level1"
		target_dscp = "CS0"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileConfigs(rName string) string {
	fmt.Println("=== STEP  creating multiple external_network_instance_profile")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_external_network_instance_profile" "test1" {
        l3_outside_dn  = aci_l3_outside.test.id
        name = "%s"
    }

	resource "aci_external_network_instance_profile" "test2" {
        l3_outside_dn  = aci_l3_outside.test.id
        name = "%s"
    }

	resource "aci_external_network_instance_profile" "test3" {
        l3_outside_dn  = aci_l3_outside.test.id
        name = "%s"
    }
	`, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccExternalNetworkInstanceProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing external_network_instance_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
        name = "%s"
    }
	`, rName, rName, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileWithoutL3Out(rName string) string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile creation without creating l3_outside")
	resource := fmt.Sprintf(`
	resource "aci_external_network_instance_profile" "test" {
        name = "%s"
    }
	`, rName)
	return resource
}

func CreateAccExternalNetworkInstanceProfileWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing external_network_instance_profile creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_external_network_instance_profile" "test" {
        l3_outside_dn  = aci_l3_outside.test.id
    }
	`, rName, rName)
	return resource
}

func testAccCheckAciExternalNetworkInstanceProfileExists(name string, external_network_instance_profile *models.ExternalNetworkInstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("External Network Instance Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No External Network Instance Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		external_network_instance_profileFound := models.ExternalNetworkInstanceProfileFromContainer(cont)
		if external_network_instance_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("External Network Instance Profile %s not found", rs.Primary.ID)
		}
		*external_network_instance_profile = *external_network_instance_profileFound
		return nil
	}
}

func testAccCheckAciExternalNetworkInstanceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing external_network_instance_profile destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_external_network_instance_profile" {
			cont, err := client.Get(rs.Primary.ID)
			external_network_instance_profile := models.ExternalNetworkInstanceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("External Network Instance Profile %s Still exists", external_network_instance_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciExternalNetworkInstanceProfileIdEqual(enip1, enip2 *models.ExternalNetworkInstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if enip1.DistinguishedName != enip2.DistinguishedName {
			return fmt.Errorf("ExternalNetworkInstanceProfile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciExternalNetworkInstanceProfileIdNotEqual(enip1, enip2 *models.ExternalNetworkInstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if enip1.DistinguishedName == enip2.DistinguishedName {
			return fmt.Errorf("ExternalNetworkInstanceProfile DNs are equal")
		}
		return nil
	}
}
