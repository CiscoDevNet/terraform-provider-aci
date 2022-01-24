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

func TestAccAciVMMDomain_Basic(t *testing.T) {
	var vmm_domain_default models.VMMDomain
	var vmm_domain_updated models.VMMDomain
	resourceName := "aci_vmm_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMDomainWithoutRequired(rName, "provider_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMDomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_default),
					resource.TestCheckResourceAttr(resourceName, "provider_profile_dn", vmmProvProfileDn),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "read-write"),
					resource.TestCheckResourceAttr(resourceName, "arp_learning", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ave_time_out", "30"),
					resource.TestCheckResourceAttr(resourceName, "config_infra_pg", "no"),
					resource.TestCheckResourceAttr(resourceName, "ctrl_knob", "epDpVerify"),
					resource.TestCheckResourceAttr(resourceName, "delimiter", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_ave", "no"),
					resource.TestCheckResourceAttr(resourceName, "enable_tag", "no"),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "enf_pref", "hw"),
					resource.TestCheckResourceAttr(resourceName, "ep_inventory_type", "on-link"),
					resource.TestCheckResourceAttr(resourceName, "ep_ret_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "hv_avail_monitor", "no"),
					resource.TestCheckResourceAttr(resourceName, "mcast_addr", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "mode", "default"),
					resource.TestCheckResourceAttr(resourceName, "pref_encap_mode", "unspecified"),
				),
			},
			{
				Config: CreateAccVMMDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "provider_profile_dn", "uni/vmmp-VMware"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vmm_domain"),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "read-write"),
					resource.TestCheckResourceAttr(resourceName, "arp_learning", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ave_time_out", "10"),
					resource.TestCheckResourceAttr(resourceName, "config_infra_pg", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ctrl_knob", "none"),
					resource.TestCheckResourceAttr(resourceName, "delimiter", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_ave", "yes"),
					resource.TestCheckResourceAttr(resourceName, "enable_tag", "yes"),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "vlan"),
					resource.TestCheckResourceAttr(resourceName, "enf_pref", "sw"),
					resource.TestCheckResourceAttr(resourceName, "ep_inventory_type", "none"),
					resource.TestCheckResourceAttr(resourceName, "ep_ret_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "hv_avail_monitor", "yes"),
					resource.TestCheckResourceAttr(resourceName, "mcast_addr", "224.0.0.22"),
					resource.TestCheckResourceAttr(resourceName, "mode", "default"),
					resource.TestCheckResourceAttr(resourceName, "pref_encap_mode", "vlan"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVMMDomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMDomainConfigWithRequiredParams(vmmProvProfileDn, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "provider_profile_dn", vmmProvProfileDn),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVMMDomainIdNotEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainConfig(rName),
			},
			{
				Config: CreateAccVMMDomainConfigWithRequiredParams(vmmProvProfileDnOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "provider_profile_dn", vmmProvProfileDnOther),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciVMMDomainIdNotEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
		},
	})
}

func TestAccAciVMMDomain_Update(t *testing.T) {
	var vmm_domain_default models.VMMDomain
	var vmm_domain_updated models.VMMDomain
	resourceName := "aci_vmm_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameOther := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_default),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "pref_encap_mode", "vxlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "pref_encap_mode", "vxlan"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "enf_pref", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "enf_pref", "unknown"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "encap_mode", "vxlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "vxlan"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "ave_time_out", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "ave_time_out", "300"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "ave_time_out", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "ave_time_out", "150"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "ep_ret_time", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "ep_ret_time", "300"),
					testAccCheckAciVMMDomainIdEqual(&vmm_domain_default, &vmm_domain_updated),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rNameOther, "access_mode", "read-only"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "read-only"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "delimiter", "_"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "delimiter", "_"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rNameOther, "mode", "n1kv"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "n1kv"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "mode", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "unknown"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rNameOther, "mode", "ovs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "ovs"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "mode", "k8s"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "k8s"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rNameOther, "mode", "rhev"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "rhev"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "mode", "cf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "cf"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rNameOther, "mode", "openshift"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "openshift"),
				),
			},
			{
				Config: CreateAccVMMDomainUpdatedAttr(rName, "mode", "cf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists(resourceName, &vmm_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "cf"),
				),
			},
			{
				Config: CreateAccVMMDomainConfig(rName),
			},
		},
	})
}

func TestAccAciVMMDomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMDomainConfig(rName),
			},
			{
				Config:      CreateAccVMMDomainWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "access_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "arp_learning", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ave_time_out", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ave_time_out", "9"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ave_time_out", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "config_infra_pg", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ctrl_knob", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "delimiter", randomValue),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "enable_ave", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "enable_tag", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "encap_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "enf_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ep_inventory_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ep_ret_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ep_ret_time", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "ep_ret_time", "601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "hv_avail_monitor", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "mcast_addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, "pref_encap_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccVMMDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVMMDomainConfig(rName),
			},
		},
	})
}

func TestAccAciVMMDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMDomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciVMMDomainExists(name string, vmm_domain *models.VMMDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_domainFound := models.VMMDomainFromContainer(cont)
		if vmm_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Domain %s not found", rs.Primary.ID)
		}
		*vmm_domain = *vmm_domainFound
		return nil
	}
}

func testAccCheckAciVMMDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vmm_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vmm_domain" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_domain := models.VMMDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Domain %s Still exists", vmm_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVMMDomainIdEqual(m1, m2 *models.VMMDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vmm_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVMMDomainIdNotEqual(m1, m2 *models.VMMDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vmm_domain DNs are equal")
		}
		return nil
	}
}

func CreateVMMDomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_domain creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "provider_profile_dn":
		rBlock += `
	resource "aci_vmm_domain" "test" {
	#	provider_profile_dn  = "%s"
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, vmmProvProfileDn)
}

func CreateAccVMMDomainConfigWithRequiredParams(provProfileDn, rName string) string {
	fmt.Printf("=== STEP  testing vmm_domain creation with provider_profile_dn %s and name %s\n", provProfileDn, rName)
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}
	`, provProfileDn, rName)
	return resource
}

func CreateAccVMMDomainConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing vmm_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}
	`, vmmProvProfileDn, rName)
	return resource
}

func CreateAccVMMDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing vmm_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}
	`, vmmProvProfileDn, rName)
	return resource
}

func CreateAccVMMDomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple vmm_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s_${count.index}"
		count = 5
	}
	`, vmmProvProfileDn, rName)
	return resource
}

func CreateAccVMMDomainWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vmm_domain creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccVMMDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  =  "%s"
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		ave_time_out = "10"
		config_infra_pg = "yes"
		ctrl_knob = "none"
		enable_ave = "yes"
		enable_tag = "yes"
		encap_mode = "vlan"
		enf_pref = "sw"
		ep_inventory_type = "none"
		ep_ret_time = "600"
		hv_avail_monitor = "yes"
		mcast_addr = "224.0.0.22"
		pref_encap_mode = "vlan"
		name_alias = "test_vmm_domain"
	}
	`, vmmProvProfileDn, rName)

	return resource
}

func CreateAccVMMDomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vmm_domain updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_vmm_domain" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vmm_domain"
		access_mode = "read-only"
		arp_learning = "enabled"
		ave_time_out = "11"
		config_infra_pg = "yes"
		ctrl_knob = "none"
		delimiter = ""
		enable_ave = "yes"
		enable_tag = "yes"
		enable_vm_folder = "yes"
		encap_mode = "ivxlan"
		enf_pref = "sw"
		ep_inventory_type = "none"
		ep_ret_time = "1"
		hv_avail_monitor = "yes"
		mcast_addr = ""
		mode = "cf"
		pref_encap_mode = "vlan"
		
	}
	`)

	return resource
}

func CreateAccVMMDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vmm_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
		%s = "%s"
	}
	`, vmmProvProfileDn, rName, attribute, value)
	return resource
}
