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

func TestAccAciL3ExtSubnet_Basic(t *testing.T) {
	var l3_ext_subnet_default models.L3ExtSubnet
	var l3_ext_subnet_updated models.L3ExtSubnet
	resourceName := "aci_l3_ext_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.1.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	ipUpdated, _ := acctest.RandIpAddress("10.0.2.0/17")
	ipUpdated = fmt.Sprintf("%s/17", ipUpdated)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3ExtSubnetWithoutRequired(rName, rName, rName, ip, "external_network_instance_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3ExtSubnetWithoutRequired(rName, rName, rName, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3ExtSubnetConfig(rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_default),
					resource.TestCheckResourceAttr(resourceName, "external_network_instance_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "aggregate", ""),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_subnet_to_profile.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_subnet_to_rt_summ", ""),
				),
			},
			{
				Config: CreateAccL3ExtSubnetConfigWithOptionalValues(rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3_ext_subnet"),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "shared-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "external_network_instance_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_subnet_to_profile.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_subnet_to_rt_summ", ""),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3ExtSubnetRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3ExtSubnetWithInavalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extSubnet (.)+`),
			},
			{
				Config: CreateAccL3ExtSubnetConfigWithRequiredParams(rNameUpdated, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "external_network_instance_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					testAccCheckAciL3ExtSubnetIdNotEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetConfig(rName, rName, rName, ip),
			},
			{
				Config: CreateAccL3ExtSubnetConfigWithRequiredParams(rName, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "external_network_instance_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					testAccCheckAciL3ExtSubnetIdNotEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
		},
	})
}

func TestAccAciL3ExtSubnet_Update(t *testing.T) {
	var l3_ext_subnet_default models.L3ExtSubnet
	var l3_ext_subnet_updated models.L3ExtSubnet
	resourceName := "aci_l3_ext_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip := "0.0.0.0/0"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3ExtSubnetConfig(rName, rName, rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_default),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "aggregate", "import-rtctrl"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "import-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "aggregate", "export-rtctrl"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "export-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "export-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "import-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"export-rtctrl", "import-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"export-rtctrl", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-security", "shared-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-security", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "export-rtctrl", "import-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "import-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-security", "shared-security", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "export-rtctrl", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "import-security", "shared-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "import-security", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"export-rtctrl", "import-security", "shared-security"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-security"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"export-rtctrl", "import-security", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
			{
				Config: CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"import-rtctrl", "export-rtctrl", "import-security", "shared-security", "shared-rtctrl"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists(resourceName, &l3_ext_subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "import-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "export-rtctrl"),
					resource.TestCheckResourceAttr(resourceName, "scope.2", "import-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.3", "shared-security"),
					resource.TestCheckResourceAttr(resourceName, "scope.4", "shared-rtctrl"),
					testAccCheckAciL3ExtSubnetIdEqual(&l3_ext_subnet_default, &l3_ext_subnet_updated),
				),
			},
		},
	})
}

func TestAccAciL3ExtSubnet_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.3.0/18")
	ip = fmt.Sprintf("%s/18", ip)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3ExtSubnetConfig(rName, rName, rName, ip),
			},
			{
				Config:      CreateAccL3ExtSubnetWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extSubnet (.)+`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "aggregate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of(.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of(.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"export-rtctrl", "export-rtctrl"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttrList(rName, rName, rName, ip, "scope", StringListtoString([]string{"shared-security", "shared-rtctrl"})),
				ExpectError: regexp.MustCompile(`shared_security scope also needs import_security`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, "aggregate", "import-rtctrl"),
				ExpectError: regexp.MustCompile(`Invalid Configuration - Import/Export Subnet Aggregation Supported Only`),
			},
			{
				Config:      CreateAccL3ExtSubnetUpdatedAttr(rName, rName, rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3ExtSubnetConfig(rName, rName, rName, ip),
			},
		},
	})
}

func TestAccAciL3ExtSubnet_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.0.4.0/19")
	ip = fmt.Sprintf("%s/19", ip)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3ExtSubnetsConfig(rName, rName, rName, ip),
			},
		},
	})
}

func CreateAccL3ExtSubnetsConfig(fvTenantName, l3extOutName, l3extInstPName, ip string) string {
	fmt.Println("=== STEP  testing multiple l3_ext_subnet creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test1" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_external_network_instance_profile" "test2" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}

	resource "aci_external_network_instance_profile" "test3" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test1" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test1.id
		ip  = "%s"
	}

	resource "aci_l3_ext_subnet" "test2" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test2.id
		ip  = "%s"
	}

	resource "aci_l3_ext_subnet" "test3" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test3.id
		ip  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extInstPName+"1", l3extInstPName+"2", l3extInstPName+"3", ip, ip, ip)
	return resource
}

func CreateAccL3ExtSubnetWithInavalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing l3_ext_subnet creation with invalid IP")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = "${aci_external_network_instance_profile.test.id}"
		ip  = "%s0"
	}
	`, rName, rName, rName, ip)
	return resource
}

func testAccCheckAciL3ExtSubnetExists(name string, l3_ext_subnet *models.L3ExtSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Ext Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Ext Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_ext_subnetFound := models.L3ExtSubnetFromContainer(cont)
		if l3_ext_subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Ext Subnet %s not found", rs.Primary.ID)
		}
		*l3_ext_subnet = *l3_ext_subnetFound
		return nil
	}
}

func testAccCheckAciL3ExtSubnetDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3_ext_subnet destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3_ext_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			l3_ext_subnet := models.L3ExtSubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Ext Subnet %s Still exists", l3_ext_subnet.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3ExtSubnetIdEqual(m1, m2 *models.L3ExtSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3_ext_subnet DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3ExtSubnetIdNotEqual(m1, m2 *models.L3ExtSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3_ext_subnet DNs are equal")
		}
		return nil
	}
}

func CreateL3ExtSubnetWithoutRequired(fvTenantName, l3extOutName, l3extInstPName, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_ext_subnet creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	`
	switch attrName {
	case "external_network_instance_profile_dn":
		rBlock += `
	resource "aci_l3_ext_subnet" "test" {
	#	external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}
		`
	case "ip":
		rBlock += `
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
	#	ip  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extInstPName, ip)
}

func CreateAccL3ExtSubnetConfigWithRequiredParams(rName, ip string) string {
	fmt.Println("=== STEP  testing l3_ext_subnet creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccL3ExtSubnetConfig(fvTenantName, l3extOutName, l3extInstPName, ip string) string {
	fmt.Println("=== STEP  testing l3_ext_subnet creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip)
	return resource
}

func CreateAccL3ExtSubnetWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing l3_ext_subnet creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_tenant.test.id
		ip  = "%s"	
	}
	`, rName, ip)
	return resource
}

func CreateAccL3ExtSubnetConfigWithOptionalValues(fvTenantName, l3extOutName, l3extInstPName, ip string) string {
	fmt.Println("=== STEP  Basic: testing l3_ext_subnet creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = "${aci_external_network_instance_profile.test.id}"
		ip  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3_ext_subnet"
		scope = ["export-rtctrl"]
		aggregate = "shared-rtctrl"
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip)

	return resource
}

func CreateAccL3ExtSubnetRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3_ext_subnet updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3_ext_subnet" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_l3_ext_subnet"
		aggregate = "export-rtctrl"
		scope = ["export-rtctrl"]
	}
	`)

	return resource
}

func CreateAccL3ExtSubnetUpdatedAttr(fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_ext_subnet attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value)
	return resource
}

func CreateAccL3ExtSubnetUpdatedAttrList(fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_ext_subnet attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_external_network_instance_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3_ext_subnet" "test" {
		external_network_instance_profile_dn  = aci_external_network_instance_profile.test.id
		ip  = "%s"
		%s = %s
	}
	`, fvTenantName, l3extOutName, l3extInstPName, ip, attribute, value)
	return resource
}
