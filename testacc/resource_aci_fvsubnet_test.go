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

func TestAccAciSubnet_Basic(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_updated models.Subnet
	resourceName := "aci_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	ipother, _ := acctest.RandIpAddress("10.21.0.0/16")
	ipother = fmt.Sprintf("%s/16", ipother)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSubnetWithoutParentDn(ip),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSubnetWithoutIP(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),  // ctrl is list type attribute with default value "nd", comparing length of ctrl with 1
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"), // comparing 0'th index element of ctrl with "nd"
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "preferred", "no"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),       // scope is list type attribute with default value "private", comparing length of scope with 1
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"), // comparing 0'th index element of scope with "nd"
					resource.TestCheckResourceAttr(resourceName, "virtual", "no"),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
				),
			},
			{
				Config: CreateAccSubnetConfigWithOptionalValues(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_subnet"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "querier"),
					resource.TestCheckResourceAttr(resourceName, "description", "subnet"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_subnet"),
					resource.TestCheckResourceAttr(resourceName, "preferred", "yes"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					resource.TestCheckResourceAttr(resourceName, "virtual", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSubnetRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccSubnetWithInvalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvSubnet (.)+`),
			},
			{
				Config: CreateAccSubnetConfigWithParentDnAndName(rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rOtherName, rOtherName)),
					testAccCheckAciSubnetIdNotEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
			{
				Config: CreateAccSubnetConfigWithParentDnAndName(rName, ipother),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ip", ipother),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
					testAccCheckAciSubnetIdNotEqual(&subnet_default, &subnet_updated),
				),
			},
		},
	})
}

func TestAccAciSubnet_Update(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_updated models.Subnet
	resourceName := "aci_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"unspecified"})), // StringListtoString will convert array of string to string with quoated elements, so that it can be passed into configuration
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "unspecified"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"no-default-gateway"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "no-default-gateway"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "no-default-gateway"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "no-default-gateway"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"no-default-gateway", "querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "no-default-gateway"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "no-default-gateway", "querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "no-default-gateway"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"querier", "no-default-gateway", "nd"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "querier"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "no-default-gateway"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "nd"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd"})),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public", "shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"shared", "public"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "public"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
		},
	})
}

func TestAccAciSubnet_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
			{
				Config:      CreateAccSubnetWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`is not valid bridge_domain_dn`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "virtual", randomValue),
				ExpectError: regexp.MustCompile(`expected virtual to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "preferred", randomValue),
				ExpectError: regexp.MustCompile(`expected preferred to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected scope.0 to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public", "public"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public"})),
				ExpectError: regexp.MustCompile(`Invalid Configuration : Subnet scope cannot be both private and public`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public", "shared"})),
				ExpectError: regexp.MustCompile(`Invalid Configuration : Subnet scope cannot be both private and public`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected ctrl.0 to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "nd"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"unspecified", "nd"})),
				ExpectError: regexp.MustCompile(`(.)+ should't be used along with other values`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
		},
	})
}

func TestAccAciSubnet_reltionalParameters(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_rel1 models.Subnet
	var subnet_rel2 models.Subnet
	resourceName := "aci_subnet.test"
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	relRes1 := makeTestVariable(acctest.RandString(5))
	relRes2 := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
				),
			},
			{
				Config: CreateAccSubnetUpdatedbdSubnetIntial(rName, ip, relRes1, relRes1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, relRes1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, relRes1)),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_rel1),
				),
			},
			{
				Config: CreateAccSubnetUpdatedbdSubnetFinal(rName, ip, relRes2, relRes1, relRes2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, relRes2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, relRes1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, relRes2)),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_rel2),
				),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "0"),
				),
			},
		},
	})
}

func TestAccAciSubnet_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip1 := fmt.Sprintf("%s/16", ip)
	ip2 := fmt.Sprintf("%s/17", ip)
	ip3 := fmt.Sprintf("%s/18", ip)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetsConfig(rName, ip1, ip2, ip3),
			},
		},
	})
}

func testAccCheckAciSubnetExists(name string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound
		return nil
	}
}

func testAccCheckAciSubnetDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing subnet destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetIdEqual(sn1, sn2 *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sn1.DistinguishedName != sn2.DistinguishedName {
			return fmt.Errorf("Subnet DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSubnetIdNotEqual(sn1, sn2 *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sn1.DistinguishedName == sn2.DistinguishedName {
			return fmt.Errorf("Subnet DNs are equal")
		}
		return nil
	}
}

func CreateSubnetWithoutParentDn(ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation without creating parent resource")
	resource := fmt.Sprintf(`
	resource "aci_subnet" "test" {
		ip = "%s"
	}
	`, ip)
	return resource
}

func CreateSubnetWithoutIP(rName string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation without giving ip")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccSubnetConfig(rName, ip string) string {
	fmt.Println("=== STEP  testing subnet creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetsConfig(rName, ip1, ip2, ip3 string) string {
	fmt.Println("=== STEP  creating multiple subnet")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test1"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	resource "aci_subnet" "test2"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}

	resource "aci_subnet" "test3"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}
	`, rName, rName, ip1, ip2, ip3)
	return resource
}

func CreateAccSubnetWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing subnet creation with invalid parent_dn")
	resource := fmt.Sprintf(`
	resource "aci_fc_domain" "test" {
		name = "%s"
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_fc_domain.test.id
		ip = "%s"
	}
	`, rName, ip)
	return resource
}

func CreateAccSubnetConfigWithOptionalValues(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		description = "subnet"
        annotation = "tag_subnet"
        ctrl = ["nd", "querier"]
        name_alias = "alias_subnet"
        preferred = "yes"
        scope = ["private", "shared"]
        virtual = "yes"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing subnet update without optional parameters")
	resource := fmt.Sprintln(`

	resource "aci_subnet" "test" {
		description = "subnet"
        annotation = "tag_subnet"
        ctrl = ["nd", "querier"]
        name_alias = "alias_subnet"
        preferred = "yes"
        scope = ["private", "shared"]
        virtual = "yes"
	}
	`)
	return resource
}

func CreateAccSubnetUpdatedbdSubnetIntial(rName, ip, bdSubnetToProfileName, bdSubnetToOutName string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test" {
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		relation_fv_rs_bd_subnet_to_profile = aci_bgp_route_control_profile.test.id
		relation_fv_rs_bd_subnet_to_out = [aci_l3_outside.test1.id]
	}
	`, rName, bdSubnetToProfileName, bdSubnetToOutName, rName, ip)
	return resource
}

func CreateAccSubnetUpdatedbdSubnetFinal(rName, ip, bdSubnetToProfileName, bdSubnetToOutName1, bdSubnetToOutName2 string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test" {
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test2" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		relation_fv_rs_bd_subnet_to_profile = aci_bgp_route_control_profile.test.id
		relation_fv_rs_bd_subnet_to_out = [aci_l3_outside.test1.id,aci_l3_outside.test2.id]
	}
	`, rName, bdSubnetToProfileName, bdSubnetToOutName1, bdSubnetToOutName2, rName, ip)
	return resource
}

func CreateAccSubnetConfigWithParentDnAndName(rName, ip string) string {
	fmt.Printf("=== STEP  Basic: testing subnet creation with parent resource name %s and ip %s\n", rName, ip)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetWithInvalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with invalid IP")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s0"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetUpdatedAttr(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing subnet attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		%s = "%s"
	}
	`, rName, rName, ip, attribute, value)
	return resource
}

func CreateAccSubnetUpdatedAttrList(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		%s = %s
	}
	`, rName, rName, ip, attribute, value)
	return resource
}
