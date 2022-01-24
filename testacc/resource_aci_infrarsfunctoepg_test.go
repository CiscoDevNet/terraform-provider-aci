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

func TestAccAciEPGsUsingFunction_Basic(t *testing.T) {
	var epgs_using_function_default models.EPGsUsingFunction
	var epgs_using_function_updated models.EPGsUsingFunction
	resourceName := "aci_epgs_using_function.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	tDn := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEPGsUsingFunctionWithoutRequired(rName, tDn, "access_generic_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEPGsUsingFunctionWithoutRequired(rName, tDn, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEPGsUsingFunctionWithoutRequired(rName, tDn, "encap"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(rName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_default),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", fmt.Sprintf("uni/infra/attentp-%s/gen-default", rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "unknown"),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfigWithOptionalValues(rName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", fmt.Sprintf("uni/infra/attentp-%s/gen-default", rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-2"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "vlan-2"),
					testAccCheckAciEPGsUsingFunctionIdEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEPGsUsingFunctionRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", fmt.Sprintf("uni/infra/attentp-%s/gen-default", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					testAccCheckAciEPGsUsingFunctionIdNotEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(rName, tDn),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", fmt.Sprintf("uni/infra/attentp-%s/gen-default", rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-1"),
					testAccCheckAciEPGsUsingFunctionIdNotEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(rName, tDn),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfigWithUpdatedEncap(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", fmt.Sprintf("uni/infra/attentp-%s/gen-default", rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-2"),
					testAccCheckAciEPGsUsingFunctionIdEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
		},
	})
}

func TestAccAciEPGsUsingFunction_Update(t *testing.T) {
	var epgs_using_function_default models.EPGsUsingFunction
	var epgs_using_function_updated models.EPGsUsingFunction
	resourceName := "aci_epgs_using_function.test"
	rName := makeTestVariable(acctest.RandString(5))
	tDn := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGsUsingFunctionConfig(rName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_default),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionUpdatedAttr(rName, tDn, "mode", "untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "untagged"),
					testAccCheckAciEPGsUsingFunctionIdEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(rName, tDn),
			},
		},
	})
}

func TestAccAciEPGsUsingFunction_Negative(t *testing.T) {
	tDn := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, tDn),
			},
			{
				Config:      CreateAccEPGsUsingFunctionWithInValidParentDn(infraAttEntityPName, tDn),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionWithInValidTDn(infraAttEntityPName, tDn),
				ExpectError: regexp.MustCompile(`Invalid target DN`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionWithInValidEncap(infraAttEntityPName, tDn, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, tDn, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, tDn, "instr_imedcy", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, tDn, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, tDn, "primary_encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, tDn),
			},
		},
	})
}

func TestAccAciEPGsUsingFunction_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	tDn := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGsUsingFunctionConfigMultiple(rName, tDn),
			},
		},
	})
}

func testAccCheckAciEPGsUsingFunctionExists(name string, epgs_using_function *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("EPGs Using Function %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EPGs Using Function dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		epgs_using_functionFound := models.EPGsUsingFunctionFromContainer(cont)
		if epgs_using_functionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("EPGs Using Function %s not found", rs.Primary.ID)
		}
		*epgs_using_function = *epgs_using_functionFound
		return nil
	}
}

func testAccCheckAciEPGsUsingFunctionDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing epgs_using_function destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_epgs_using_function" {
			cont, err := client.Get(rs.Primary.ID)
			epgs_using_function := models.EPGsUsingFunctionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("EPGs Using Function %s Still exists", epgs_using_function.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEPGsUsingFunctionIdEqual(m1, m2 *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("epgs_using_function DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEPGsUsingFunctionIdNotEqual(m1, m2 *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("epgs_using_function DNs are equal")
		}
		return nil
	}
}

func CreateEPGsUsingFunctionWithoutRequired(rName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation without ", attrName)
	rBlock := `
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}
	`
	switch attrName {
	case "access_generic_dn":
		rBlock += `
	resource "aci_epgs_using_function" "test" {
	#	access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`
	case "tdn":
		rBlock += `
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		#	tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`
	case "encap":
		rBlock += `
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
	#	encap = "vlan-1"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, rName)
}

func CreateAccEPGsUsingFunctionConfigWithRequiredParams(rName, rName1 string) string {
	fmt.Println("=== STEP  testing epgs_using_function creation with Updation required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`, rName, rName, rName, rName1)
	return resource
}

func CreateAccEPGsUsingFunctionConfigWithUpdatedEncap(rName, rName1 string) string {
	fmt.Println("=== STEP  testing epgs_using_function creation with Updation of Encap")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-2"
	}
	`, rName, rName, rName, rName1)
	return resource
}

func CreateAccEPGsUsingFunctionConfig(rName, tDn string) string {
	fmt.Println("=== STEP  testing epgs_using_function creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGsUsingFunctionConfigMultiple(rName, tDn string) string {
	fmt.Println("=== STEP  testing multiple epgs_using_function creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test1" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test2" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test3" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}

	resource "aci_epgs_using_function" "test1" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test1.id
		encap = "vlan-1"
	}

	resource "aci_epgs_using_function" "test2" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test2.id
		encap = "vlan-1"
	}

	resource "aci_epgs_using_function" "test3" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test3.id
		encap = "vlan-1"
	}
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	`, rName, rName, rName, rName+"1", rName+"2", rName+"3", rName)
	return resource
}

func CreateAccEPGsUsingFunctionWithInValidParentDn(rName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing epgs_using_function creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_tenant.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGsUsingFunctionWithInValidTDn(rName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing epgs_using_function creation with invalid TDn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_profile.test.id
		encap = "vlan-1"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGsUsingFunctionWithInValidEncap(rName, tDn, encap string) string {
	fmt.Println("=== STEP  Negative Case: testing epgs_using_function creation with invalid encap")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_profile.test.id
		encap = "%s"
	}
	`, rName, rName, rName, rName, encap)
	return resource
}
func CreateAccEPGsUsingFunctionConfigWithOptionalValues(rName, tDn string) string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-2"
		annotation = "orchestrator:terraform_testacc"
		instr_imedcy = "immediate"
		mode = "native"
		primary_encap = "vlan-2"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGsUsingFunctionRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_epgs_using_function" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_epgs_using_function"
		instr_imedcy = "immediate"
		mode = "native"
		primary_encap = "unknown"
	}
	`)
	return resource
}

func CreateAccEPGsUsingFunctionUpdatedAttr(rName, tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing epgs_using_function attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	}

	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
		name = "default"
	}

	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tdn  = aci_application_epg.test.id
		encap = "vlan-1"
		%s = "%s"
	}
	
	`, rName, rName, rName, rName, attribute, value)
	return resource
}
