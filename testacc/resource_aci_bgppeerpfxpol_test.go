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

func TestAccAciBGPPeerPrefix_Basic(t *testing.T) {
	var bgp_peer_prefix_default models.BGPPeerPrefixPolicy
	var bgp_peer_prefix_updated models.BGPPeerPrefixPolicy
	resourceName := "aci_bgp_peer_prefix.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPPeerPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPPeerPrefixWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPPeerPrefixWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "action", "reject"),
					resource.TestCheckResourceAttr(resourceName, "max_pfx", "20000"),
					resource.TestCheckResourceAttr(resourceName, "restart_time", "infinite"),
					resource.TestCheckResourceAttr(resourceName, "thresh", "75"),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bgp_peer_prefix"),
					resource.TestCheckResourceAttr(resourceName, "action", "log"),
					resource.TestCheckResourceAttr(resourceName, "max_pfx", "2"),
					resource.TestCheckResourceAttr(resourceName, "restart_time", "1"),
					resource.TestCheckResourceAttr(resourceName, "thresh", "1"),
					testAccCheckAciBGPPeerPrefixIdEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBGPPeerPrefixConfig(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdateWithoutRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPPeerPrefixConfigWithUpdatedRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBGPPeerPrefixIdNotEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccBGPPeerPrefixConfigWithUpdatedRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBGPPeerPrefixIdNotEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
		},
	})
}

func TestAccAciBGPPeerPrefix_Update(t *testing.T) {
	var bgp_peer_prefix_default models.BGPPeerPrefixPolicy
	var bgp_peer_prefix_updated models.BGPPeerPrefixPolicy
	resourceName := "aci_bgp_peer_prefix.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPPeerPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_default),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "action", "restart"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "action", "restart"),
					testAccCheckAciBGPPeerPrefixIdEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "action", "shut"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "action", "shut"),
					testAccCheckAciBGPPeerPrefixIdEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "max_pfx", "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "max_pfx", "1"),
					testAccCheckAciBGPPeerPrefixIdEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "thresh", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixExists(resourceName, &bgp_peer_prefix_updated),
					resource.TestCheckResourceAttr(resourceName, "thresh", "100"),
					testAccCheckAciBGPPeerPrefixIdEqual(&bgp_peer_prefix_default, &bgp_peer_prefix_updated),
				),
			},
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciBGPPeerPrefix_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPPeerPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccBGPPeerPrefixWithInValidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "action", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "max_pfx", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "max_pfx", "3000002"),
				ExpectError: regexp.MustCompile(`Property maxPfx of (.)* is out of range`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "restart_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "restart_time", "0"),
				ExpectError: regexp.MustCompile(`Property restartTime of (.)* is out of range`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "thresh", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, "thresh", "200"),
				ExpectError: regexp.MustCompile(`Property thresh of (.)* is out of range`),
			},
			{
				Config:      CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBGPPeerPrefixConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciBGPPeerPrefix_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPPeerPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPPeerPrefixesConfig(rName, rName),
			},
		},
	})
}

func testAccCheckAciBGPPeerPrefixExists(name string, bgp_peer_prefix *models.BGPPeerPrefixPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP Peer Prefix %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP Peer Prefix dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_peer_prefixFound := models.BGPPeerPrefixPolicyFromContainer(cont)
		if bgp_peer_prefixFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP Peer Prefix %s not found", rs.Primary.ID)
		}
		*bgp_peer_prefix = *bgp_peer_prefixFound
		return nil
	}
}

func testAccCheckAciBGPPeerPrefixDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bgp_peer_prefix destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_peer_prefix" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_peer_prefix := models.BGPPeerPrefixPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP Peer Prefix %s Still exists", bgp_peer_prefix.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBGPPeerPrefixIdEqual(m1, m2 *models.BGPPeerPrefixPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_peer_prefix DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBGPPeerPrefixIdNotEqual(m1, m2 *models.BGPPeerPrefixPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_peer_prefix DNs are equal")
		}
		return nil
	}
}

func CreateBGPPeerPrefixWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_peer_prefix creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bgp_peer_prefix" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPPeerPrefixConfigWithUpdatedRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing bgp_peer_prefix creation with updated required arguments with Tenant name %s and BGP Peer Prefix %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPPeerPrefixConfig(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing bgp_peer_prefix creation with required arguments only with name %s\n", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPPeerPrefixesConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing Multiple bgp_peer_prefixes creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_peer_prefix" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_peer_prefix" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_peer_prefix" "test3" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccBGPPeerPrefixWithInValidParentDn(prName, rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bgp_peer_prefix creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_aaa_domain" "test"{
		name = "%s"
	}
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_aaa_domain.test.id
		name  = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccBGPPeerPrefixConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_peer_prefix creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_peer_prefix"
		action = "log"
		max_pfx = "2"
		restart_time = "1"
		thresh = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBGPPeerPrefixUpdateWithoutRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bgp_peer_prefix updation without required parameters")
	resource := fmt.Sprint(`
	resource "aci_bgp_peer_prefix" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_peer_prefix"
		action = "log"
		max_pfx = "2"
		restart_time = "2"
		thresh = "2"
	}
	`)

	return resource
}

func CreateAccBGPPeerPrefixUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_peer_prefix attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
