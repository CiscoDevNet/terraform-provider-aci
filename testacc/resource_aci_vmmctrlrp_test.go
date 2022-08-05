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

func TestAccAciVMMController_Basic(t *testing.T) {
	var vmm_controller_default models.VMMController
	var vmm_controller_updated models.VMMController
	resourceName := "aci_vmm_controller.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.0.0/16")
	ipUpdated, _ := acctest.RandIpAddress("10.1.0.0/16")
	rootContName := makeTestVariable(acctest.RandString(5))
	rootContNameUpdated := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMControllerWithoutRequired(vmmDomPName, rName, ip, rootContName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMControllerWithoutRequired(vmmDomPName, rName, ip, rootContName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMControllerWithoutRequired(vmmDomPName, rName, ip, rootContName, "host_or_ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMControllerWithoutRequired(vmmDomPName, rName, ip, rootContName, "root_cont_name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_default),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "host_or_ip", ip),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "root_cont_name", rootContName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "dvs_version", "unmanaged"),
					resource.TestCheckResourceAttr(resourceName, "mode", "default"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_err_msg", ""),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.0", ""),
					resource.TestCheckResourceAttr(resourceName, "n1kv_stats_mode", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "port", "0"),
					resource.TestCheckResourceAttr(resourceName, "scope", "vm"),
					resource.TestCheckResourceAttr(resourceName, "seq_num", "0"),
					resource.TestCheckResourceAttr(resourceName, "stats_mode", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "vxlan_depl_pref", "vxlan"),
				),
			},
			{
				Config: CreateAccVMMControllerConfigWithOptionalValues(vmmDomPName, rNameUpdated, ip, rootContName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "host_or_ip", ip),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "root_cont_name", rootContName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vmm_controller"),
					resource.TestCheckResourceAttr(resourceName, "dvs_version", "5.1"),
					resource.TestCheckResourceAttr(resourceName, "inventory_trig_st", "autoTriggered"),
					resource.TestCheckResourceAttr(resourceName, "mode", "cf"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_err_msg", "Error"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.0", "aaacert-invalid"),
					resource.TestCheckResourceAttr(resourceName, "n1kv_stats_mode", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "port", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope", "MicrosoftSCVMM"),
					resource.TestCheckResourceAttr(resourceName, "stats_mode", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "vxlan_depl_pref", "nsx"),
					testAccCheckAciVMMControllerIdNotEqual(&vmm_controller_default, &vmm_controller_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"inventory_trig_st"},
			},
			{
				Config:      CreateAccVMMControllerConfigUpdatedName(vmmDomPName, acctest.RandString(65), ip, rootContName),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVMMControllerRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMControllerConfigWithRequiredParams(rNameUpdated, rName, ip, rootContName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciVMMControllerIdNotEqual(&vmm_controller_default, &vmm_controller_updated),
				),
			},
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
			},
			{
				Config: CreateAccVMMControllerConfigWithRequiredParams(vmmDomPName, rNameUpdated, ipUpdated, rootContNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "host_or_ip", ipUpdated),
					resource.TestCheckResourceAttr(resourceName, "root_cont_name", rootContNameUpdated),
					testAccCheckAciVMMControllerIdNotEqual(&vmm_controller_default, &vmm_controller_updated),
				),
			},
		},
	})
}

func TestAccAciVMMController_Update(t *testing.T) {
	var vmm_controller_default models.VMMController
	var vmm_controller_updated models.VMMController
	resourceName := "aci_vmm_controller.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.1.0.0/16")
	rootContName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_default),
				),
			},
			{
				Config: CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "mode", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "unknown"),
					testAccCheckAciVMMControllerIdEqual(&vmm_controller_default, &vmm_controller_updated),
				),
			},
			{
				Config: CreateAccVMMControllerUpdatedAttrList(vmmDomPName, rName, ip, rootContName, "msft_config_issues", StringListtoString([]string{"aaacert-invalid", "duplicate-mac-in-inventory", "duplicate-rootContName", "invalid-object-in-inventory", "invalid-rootContName", "inventory-failed", "missing-hostGroup-in-cloud", "missing-rootContName", "zero-mac-in-inventory"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists(resourceName, &vmm_controller_updated),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.#", "9"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.0", "aaacert-invalid"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.1", "duplicate-mac-in-inventory"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.2", "duplicate-rootContName"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.3", "invalid-object-in-inventory"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.4", "invalid-rootContName"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.5", "inventory-failed"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.6", "missing-hostGroup-in-cloud"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.7", "missing-rootContName"),
					resource.TestCheckResourceAttr(resourceName, "msft_config_issues.8", "zero-mac-in-inventory"),
					testAccCheckAciVMMControllerIdEqual(&vmm_controller_default, &vmm_controller_updated),
				),
			},
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
			},
		},
	})
}

func TestAccAciVMMController_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.2.0.0/16")
	rootContName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
			},
			{
				Config:      CreateAccVMMControllerWithInValidParentDn(vmmDomPName, rName, ip, rootContName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMControllerConfig(vmmDomPName, rName, rName, rootContName),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "dvs_version", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "inventory_trig_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttrList(vmmDomPName, rName, ip, rootContName, "msft_config_issues", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttrList(vmmDomPName, rName, ip, rootContName, "msft_config_issues", StringListtoString([]string{"aaacert-invalid", "aaacert-invalid"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "n1kv_stats_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "scope", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "seq_num", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "stats_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, "vxlan_depl_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName),
			},
		},
	})
}

func TestAccAciVMMController_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	rootContName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMControllerConfigMultiple(vmmDomPName, rName, rootContName),
			},
		},
	})
}

func testAccCheckAciVMMControllerExists(name string, vmm_controller *models.VMMController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Controller %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Controller dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_controllerFound := models.VMMControllerFromContainer(cont)
		if vmm_controllerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Controller %s not found", rs.Primary.ID)
		}
		*vmm_controller = *vmm_controllerFound
		return nil
	}
}

func testAccCheckAciVMMControllerDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vmm_controller destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vmm_controller" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_controller := models.VMMControllerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Controller %s Still exists", vmm_controller.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVMMControllerIdEqual(m1, m2 *models.VMMController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vmm_controller DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVMMControllerIdNotEqual(m1, m2 *models.VMMController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vmm_controller DNs are equal")
		}
		return nil
	}
}

func CreateVMMControllerWithoutRequired(vmmDomPName, rName, hostOrIp, rootContName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_controller creation without ", attrName)
	rBlock := `
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
	resource "aci_vmm_controller" "test" {
	#	vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	#	name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
		`
	case "host_or_ip":
		rBlock += `
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	#	host_or_ip = "%s"
		root_cont_name = "%s"
	}
		`
	case "root_cont_name":
		rBlock += `
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
	#	root_cont_name = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn, rName, hostOrIp, rootContName)
}

func CreateAccVMMControllerConfigWithRequiredParams(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)
	return resource
}
func CreateAccVMMControllerConfigUpdatedRequiredArguments(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller creation with invalid argument")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`, providerProfileDn, vmmDomPName, rName, ip, rootContName)
	return resource
}

func CreateAccVMMControllerConfig(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)
	return resource
}

func CreateAccVMMControllerConfigUpdatedName(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)
	return resource
}
func CreateAccVMMControllerConfigMultiple(vmmDomPName, rName, rootContName string) string {
	fmt.Println("=== STEP  testing multiple vmm_controller creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		host_or_ip = "10.4.0.${count.index}"
		root_cont_name = "%s"
		name  = "%s_${count.index}"
		count = 5
	}
	`, vmmDomPName, providerProfileDn, rootContName, rName)
	return resource
}

func CreateAccVMMControllerWithInValidParentDn(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  Negative Case: testing vmm_controller creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_tenant.test.id
		name  = "%s"	
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`, vmmDomPName, rName, ip, rootContName)
	return resource
}

func CreateAccVMMControllerConfigWithOptionalValues(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_controller creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = "${aci_vmm_domain.test.id}"
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vmm_controller"
		dvs_version = "5.1"
		inventory_trig_st = "autoTriggered"
		mode = "cf"
		msft_config_err_msg = "Error"
		msft_config_issues = ["aaacert-invalid"]
		n1kv_stats_mode = "disabled"
		port = "1"
		scope = "MicrosoftSCVMM"
		stats_mode = "enabled"
		vxlan_depl_pref = "nsx"
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)

	return resource
}

func CreateAccVMMControllerRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vmm_controller updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vmm_controller" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vmm_controller"
		dvs_version = "5.1"
		inventory_trig_st = "autoTriggered"
		mode = "cf"
		msft_config_err_msg = ""
		msft_config_issues = ["aaacert-invalid"]
		n1kv_stats_mode = "disabled"
		port = "1"
		scope = "MicrosoftSCVMM"
		stats_mode = "enabled"
		vxlan_depl_pref = "nsx"
		
	}
	`)

	return resource
}

func CreateAccVMMControllerUpdatedAttr(vmmDomPName, rName, ip, rootContName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vmm_controller attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
		%s = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName, attribute, value)
	return resource
}

func CreateAccVMMControllerUpdatedAttrList(vmmDomPName, rName, ip, rootContName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vmm_controller attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
		%s = %s
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName, attribute, value)
	return resource
}
