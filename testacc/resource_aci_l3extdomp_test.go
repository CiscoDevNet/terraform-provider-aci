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

func TestAccAciL3DomainProfile_Basic(t *testing.T) {
	var l3_domain_profile_default models.L3DomainProfile
	var l3_domain_profile_updated models.L3DomainProfile
	resourceName := "aci_l3_domain_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3DomainProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3DomainProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
				),
			},
			{
				Config: CreateAccL3DomainProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3_domain_profile"),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_extnw_rs_out.#", "0"),
					testAccCheckAciL3DomainProfileIdEqual(&l3_domain_profile_default, &l3_domain_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3DomainProfileConfigWithRequiredParams(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL3DomainProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3DomainProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL3DomainProfileIdNotEqual(&l3_domain_profile_default, &l3_domain_profile_updated),
				),
			},
		},
	})
}

func TestAccAciL3DomainProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3DomainProfileConfig(rName),
			},
			{
				Config:      CreateAccL3DomainProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3DomainProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3DomainProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3DomainProfileConfig(rName),
			},
		},
	})
}

func TestAccAciL3DomainProfile_RelationParameters(t *testing.T) {
	var l3_domain_profile_default models.L3DomainProfile
	var l3_domain_profile_updated models.L3DomainProfile
	rName := makeTestVariable(acctest.RandString(5))
	relResName1 := makeTestVariable(acctest.RandString(5))
	relResName2 := makeTestVariable(acctest.RandString(5))
	resourceName := "aci_l3_domain_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3DomainProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
				),
			},
			{
				Config: CreateAccL3DomainProfileRelParams(rName, relResName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", fmt.Sprintf("uni/infra/vlanns-[%s]-static", relResName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_extnw_rs_out.#", "0"),
					testAccCheckAciL3DomainProfileIdEqual(&l3_domain_profile_default, &l3_domain_profile_updated),
				),
			},
			{
				Config: CreateAccL3DomainProfileRelParams(rName, relResName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", fmt.Sprintf("uni/infra/vlanns-[%s]-static", relResName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_extnw_rs_out.#", "0"),
					testAccCheckAciL3DomainProfileIdEqual(&l3_domain_profile_default, &l3_domain_profile_updated),
				),
			},
			{
				Config: CreateAccL3DomainProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists(resourceName, &l3_domain_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_dom_vxlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vip_addr_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_infra_rs_vlan_ns_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_extnw_rs_out.#", "0"),
					testAccCheckAciL3DomainProfileIdEqual(&l3_domain_profile_default, &l3_domain_profile_updated),
				),
			},
		},
	})
}

func TestAccAciL3DomainProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3DomainProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciL3DomainProfileExists(name string, l3_domain_profile *models.L3DomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Domain Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Domain Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_domain_profileFound := models.L3DomainProfileFromContainer(cont)
		if l3_domain_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Domain Profile %s not found", rs.Primary.ID)
		}
		*l3_domain_profile = *l3_domain_profileFound
		return nil
	}
}

func testAccCheckAciL3DomainProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3_domain_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3_domain_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3_domain_profile := models.L3DomainProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Domain Profile %s Still exists", l3_domain_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3DomainProfileIdEqual(m1, m2 *models.L3DomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3_domain_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3DomainProfileIdNotEqual(m1, m2 *models.L3DomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3_domain_profile DNs are equal")
		}
		return nil
	}
}

func CreateL3DomainProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_domain_profile creation without ", attrName)
	rBlock := `

	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_l3_domain_profile" "test" {

	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL3DomainProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l3_domain_profile creation with name", rName)
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL3DomainProfileRelParams(rName, relName string) string {
	fmt.Println("=== STEP  testing l3_domain_profile creation relation resource name", relName)
	resource := fmt.Sprintf(`
	resource "aci_vlan_pool" "test" {
		name = "%s"
		alloc_mode = "static"
	}

	resource "aci_l3_domain_profile" "test" {
		name  = "%s"
		relation_infra_rs_vlan_ns = aci_vlan_pool.test.id
	}
	`, relName, rName)
	return resource
}

func CreateAccL3DomainProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple l3_domain_profile creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccL3DomainProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing l3_domain_profile creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL3DomainProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3_domain_profile creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3_domain_profile"

	}
	`, rName)

	return resource
}

func CreateAccL3DomainProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3_domain_profile update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3_domain_profile" "test" {
		annotation = "tag"
		name_alias = "test_l3_domain_profile"

	}
	`)

	return resource
}

func CreateAccL3DomainProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_domain_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
