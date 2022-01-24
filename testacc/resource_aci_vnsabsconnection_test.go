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

func TestAccAciConnection_Basic(t *testing.T) {
	var connection_default models.Connection
	var connection_updated models.Connection
	resourceName := "aci_connection.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateConnectionWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "l4_l7_service_graph_template_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateConnectionWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConnectionConfig(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists(resourceName, &connection_default),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", fvTenantName, vnsAbsGraphName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "adj_type", "L2"),
					// resource.TestCheckResourceAttr(resourceName, "conn_dir", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "conn_type", "external"),
					resource.TestCheckResourceAttr(resourceName, "direct_connect", "no"),
					resource.TestCheckResourceAttr(resourceName, "unicast_route", "yes"),
				),
			},
			{
				Config: CreateAccConnectionConfigWithOptionalValues(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists(resourceName, &connection_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", fvTenantName, vnsAbsGraphName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_connection"),

					resource.TestCheckResourceAttr(resourceName, "adj_type", "L2"),

					resource.TestCheckResourceAttr(resourceName, "conn_dir", "consumer"),

					resource.TestCheckResourceAttr(resourceName, "conn_type", "internal"),

					resource.TestCheckResourceAttr(resourceName, "direct_connect", "yes"),

					resource.TestCheckResourceAttr(resourceName, "unicast_route", "yes"),

					testAccCheckAciConnectionIdEqual(&connection_default, &connection_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"l4_l7_service_graph_template_dn"},
			},
			{
				Config: CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "adj_type", "L3"),
			},
			{
				Config: CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "unicast_route", "yes"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "adj_type", "L3"),
					resource.TestCheckResourceAttr(resourceName, "unicast_route", "yes"),
				),
			},
			{
				Config:      CreateAccConnectionConfigUpdatedName(fvTenantName, vnsAbsGraphName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccConnectionRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConnectionConfigWithRequiredParams(rName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists(resourceName, &connection_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciConnectionIdNotEqual(&connection_default, &connection_updated),
				),
			},
			{
				Config: CreateAccConnectionConfig(fvTenantName, vnsAbsGraphName, rName),
			},
			{
				Config: CreateAccConnectionConfigWithRequiredParams(rName, rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConnectionExists(resourceName, &connection_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciConnectionIdNotEqual(&connection_default, &connection_updated),
				),
			},
		},
	})
}

func TestAccAciConnection_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConnectionConfig(fvTenantName, vnsAbsGraphName, rName),
			},
			{
				Config:      CreateAccConnectionWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "adj_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "conn_dir", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "conn_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "direct_connect", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "unicast_route", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccConnectionConfig(fvTenantName, vnsAbsGraphName, rName),
			},
		},
	})
}

func TestAccAciConnection_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConnectionConfigMultiple(fvTenantName, vnsAbsGraphName, rName),
			},
		},
	})
}

func testAccCheckAciConnectionExists(name string, connection *models.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Connection %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Connection dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		connectionFound := models.ConnectionFromContainer(cont)
		if connectionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Connection %s not found", rs.Primary.ID)
		}
		*connection = *connectionFound
		return nil
	}
}

func testAccCheckAciConnectionDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing connection destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_connection" {
			cont, err := client.Get(rs.Primary.ID)
			connection := models.ConnectionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Connection %s Still exists", connection.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConnectionIdEqual(m1, m2 *models.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("connection DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciConnectionIdNotEqual(m1, m2 *models.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("connection DNs are equal")
		}
		return nil
	}
}

func CreateConnectionWithoutRequired(fvTenantName, vnsAbsGraphName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing connection creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "l4_l7_service_graph_template_dn":
		rBlock += `
	resource "aci_connection" "test" {
	#	l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsAbsGraphName, rName)
}

func CreateAccConnectionConfigWithRequiredParams(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing connection creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}
func CreateAccConnectionConfigUpdatedName(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing connection creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccConnectionConfig(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing connection creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccConnectionConfigMultiple(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing multiple connection creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccConnectionWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing connection creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccConnectionConfigWithOptionalValues(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  Basic: testing connection creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = "${aci_l4_l7_service_graph_template.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_connection"
		adj_type = "L2"
		conn_dir = "consumer"
		conn_type = "internal"
		direct_connect = "yes"
		unicast_route = "yes"
		
	}
	`, fvTenantName, vnsAbsGraphName, rName)

	return resource
}

func CreateAccConnectionRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing connection updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_connection" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_connection"
		adj_type = "L3"
		conn_dir = "consumer"
		conn_type = "internal"
		direct_connect = "yes"
		unicast_route = "no"
		
	}
	`)

	return resource
}

func CreateAccConnectionUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing connection attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName, attribute, value)
	return resource
}
