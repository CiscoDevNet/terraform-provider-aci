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

func TestAccAciVMMCredential_Basic(t *testing.T) {
	var vmm_credential_default models.VMMCredential
	var vmm_credential_updated models.VMMCredential
	resourceName := "aci_vmm_credential.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMCredentialWithoutRequired(vmmDomPName, rName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMCredentialWithoutRequired(vmmDomPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMCredentialConfig(vmmDomPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMCredentialExists(resourceName, &vmm_credential_default),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					// resource.TestCheckResourceAttr(resourceName, "pwd", ""),
					resource.TestCheckResourceAttr(resourceName, "usr", ""),
				),
			},
			{
				Config: CreateAccVMMCredentialConfigWithOptionalValues(vmmDomPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMCredentialExists(resourceName, &vmm_credential_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vmm_credential"),
					resource.TestCheckResourceAttr(resourceName, "pwd", "ABCD"),
					resource.TestCheckResourceAttr(resourceName, "usr", "ABCD"),
					testAccCheckAciVMMCredentialIdEqual(&vmm_credential_default, &vmm_credential_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pwd"},
			},
			{
				Config:      CreateAccVMMCredentialConfigUpdatedName(vmmDomPName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVMMCredentialRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMCredentialConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMCredentialExists(resourceName, &vmm_credential_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciVMMCredentialIdNotEqual(&vmm_credential_default, &vmm_credential_updated),
				),
			},
			{
				Config: CreateAccVMMCredentialConfig(vmmDomPName, rName),
			},
			{
				Config: CreateAccVMMCredentialConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMCredentialExists(resourceName, &vmm_credential_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVMMCredentialIdNotEqual(&vmm_credential_default, &vmm_credential_updated),
				),
			},
		},
	})
}

func TestAccAciVMMCredential_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMCredentialConfig(vmmDomPName, rName),
			},
			{
				Config:      CreateAccVMMCredentialWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, "usr", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVMMCredentialConfig(vmmDomPName, rName),
			},
		},
	})
}

func TestAccAciVMMCredential_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVMMCredentialConfigMultiple(vmmDomPName, rName),
			},
		},
	})
}

func testAccCheckAciVMMCredentialExists(name string, vmm_credential *models.VMMCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Credential %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Credential dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_credentialFound := models.VMMCredentialFromContainer(cont)
		if vmm_credentialFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Credential %s not found", rs.Primary.ID)
		}
		*vmm_credential = *vmm_credentialFound
		return nil
	}
}

func testAccCheckAciVMMCredentialDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vmm_credential destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vmm_credential" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_credential := models.VMMCredentialFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Credential %s Still exists", vmm_credential.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVMMCredentialIdEqual(m1, m2 *models.VMMCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vmm_credential DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVMMCredentialIdNotEqual(m1, m2 *models.VMMCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vmm_credential DNs are equal")
		}
		return nil
	}
}

func CreateVMMCredentialWithoutRequired(vmmDomPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_credential creation without ", attrName)
	rBlock := `
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
	resource "aci_vmm_credential" "test" {
	#	vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn, rName)
}

func CreateAccVMMCredentialConfigWithRequiredParams(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing vmm_credential creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}

	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}
func CreateAccVMMCredentialConfigUpdatedName(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing vmm_credential creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}

func CreateAccVMMCredentialConfig(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing vmm_credential creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}

func CreateAccVMMCredentialConfigMultiple(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing multiple vmm_credential creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}

func CreateAccVMMCredentialWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vmm_credential creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccVMMCredentialConfigWithOptionalValues(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_credential creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = "${aci_vmm_domain.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vmm_credential"
		pwd = "ABCD"
		usr = "ABCD"
		
	}
	`, vmmDomPName, providerProfileDn, rName)

	return resource
}

func CreateAccVMMCredentialRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vmm_credential updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vmm_credential" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vmm_credential"
		pwd = ""
		usr = ""
	}
	`)

	return resource
}

func CreateAccVMMCredentialUpdatedAttr(vmmDomPName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vmm_credential attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		%s = "%s"
	}
	`, vmmDomPName, providerProfileDn, rName, attribute, value)
	return resource
}
