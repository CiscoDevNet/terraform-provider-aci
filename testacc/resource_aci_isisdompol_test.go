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
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciISISDomainPolicy_Basic(t *testing.T) {
	var isis_domain_policy_default models.ISISDomainPolicy
	var isis_domain_policy_updated models.ISISDomainPolicy
	resourceName := "aci_isis_domain_policy.test"

	isisDomPol, err := aci.GetRemoteISISDomainPolicy(sharedAciClient(), "uni/fabric/isisDomP-default")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}
	isisLvlComp, err := aci.GetRemoteISISLevel(sharedAciClient(), "uni/fabric/isisDomP-default/lvl-l1")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_default),
					// all default values varies in server hence skipping default testcase
				),
			},
			{
				Config: CreateAccISISDomainPolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_isis_domain_policy"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "257"),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "2"),
					resource.TestCheckResourceAttr(resourceName, "lsp_fast_flood", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_init_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_max_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_sec_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_init_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_max_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_sec_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "isis_level_type", "l1"),

					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restoreISISDomainPolicyToInitConfig(isisDomPol, isisLvlComp),
			},
		},
	})
}

func TestAccAciISISDomainPolicy_Update(t *testing.T) {
	var isis_domain_policy_default models.ISISDomainPolicy
	var isis_domain_policy_updated models.ISISDomainPolicy
	resourceName := "aci_isis_domain_policy.test"

	isisDomPol, err := aci.GetRemoteISISDomainPolicy(sharedAciClient(), "uni/fabric/isisDomP-default")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}
	isisLvlComp, err := aci.GetRemoteISISLevel(sharedAciClient(), "uni/fabric/isisDomP-default/lvl-l1")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_default),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_fast_flood", "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_fast_flood", "enabled"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("mtu", "4352"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "4352"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("mtu", "2048"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "2048"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("redistrib_metric", "63"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "63"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("redistrib_metric", "31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "31"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_init_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_init_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_init_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_init_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_max_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_max_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_max_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_max_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_sec_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_sec_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_sec_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "lsp_gen_sec_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_init_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_init_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_init_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_init_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_max_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_max_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_max_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_max_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_sec_intvl", "120000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_sec_intvl", "120000"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr("spf_comp_sec_intvl", "59975"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_comp_sec_intvl", "59975"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyConfig(),
			},
			{
				Config: restoreISISDomainPolicyToInitConfig(isisDomPol, isisLvlComp),
			},
		},
	})
}

func TestAccAciISISDomainPolicy_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	isisDomPol, err := aci.GetRemoteISISDomainPolicy(sharedAciClient(), "uni/fabric/isisDomP-default")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}
	isisLvlComp, err := aci.GetRemoteISISLevel(sharedAciClient(), "uni/fabric/isisDomP-default/lvl-l1")
	if err != nil {
		t.Errorf("reading initial config of isisDomPol")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfig(),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("mtu", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("mtu", "255"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("mtu", "4353"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("redistrib_metric", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("redistrib_metric", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("redistrib_metric", "64"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_fast_flood", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_init_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_init_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_init_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_max_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_max_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_max_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_sec_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_sec_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("lsp_gen_sec_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_init_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_init_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_init_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_max_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_max_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_max_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_sec_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_sec_intvl", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("spf_comp_sec_intvl", "120001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr("isis_level_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccISISDomainPolicyConfig(),
			},
			{
				Config: restoreISISDomainPolicyToInitConfig(isisDomPol, isisLvlComp),
			},
		},
	})
}

func testAccCheckAciISISDomainPolicyExists(name string, isis_domain_policy *models.ISISDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("ISIS Domain Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ISIS Domain Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		isis_domain_policyFound := models.ISISDomainPolicyFromContainer(cont)
		if isis_domain_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("ISIS Domain Policy %s not found", rs.Primary.ID)
		}
		*isis_domain_policy = *isis_domain_policyFound
		return nil
	}
}

func testAccCheckAciISISDomainPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing isis_domain_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_isis_domain_policy" {
			cont, err := client.Get(rs.Primary.ID)
			isis_domain_policy := models.ISISDomainPolicyFromContainer(cont)
			if err != nil {
				return fmt.Errorf("ISIS Domain Policy %s Still exists", isis_domain_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciISISDomainPolicyIdEqual(m1, m2 *models.ISISDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("isis_domain_policy DNs are not equal")
		}
		return nil
	}
}

func CreateAccISISDomainPolicyConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with updated naming arguments")
	resource := `
	
	resource "aci_isis_domain_policy" "test" {
	
		
	}
	`
	return resource
}
func CreateAccISISDomainPolicyConfigUpdatedName() string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with invalid name = ")
	resource := `
	
	resource "aci_isis_domain_policy" "test" {
	
		
	}
	`
	return resource
}

func CreateAccISISDomainPolicyConfig() string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with required arguments only")
	resource := `
	
	resource "aci_isis_domain_policy" "test" {
	
		
	}
	`
	return resource
}

func CreateAccISISDomainPolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing isis_domain_policy creation with optional parameters")
	resource := `
	
	resource "aci_isis_domain_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_isis_domain_policy"
		mtu = "257"
		redistrib_metric = "2"
		lsp_fast_flood = "disabled"
		lsp_gen_init_intvl = "60"
		lsp_gen_max_intvl = "60"
		lsp_gen_sec_intvl = "60"
		spf_comp_init_intvl = "60"
		spf_comp_max_intvl = "60"
		spf_comp_sec_intvl = "60"
		isis_level_type = "l1"
	}
	`

	return resource
}

func CreateAccISISDomainPolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing isis_domain_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func restoreISISDomainPolicyToInitConfig(isisDomPol *models.ISISDomainPolicy, isisLvlComp *models.ISISLevel) string {
	fmt.Println("=== STEP  Basic: restoring isis_domain_policy to original configuration")

	resource := fmt.Sprintf(`
	resource "aci_isis_domain_policy" "test" {
		annotation = "%s"
		mtu = "%s"
		redistrib_metric = "%s"
		description = "%s"
		name_alias = "%s"
		lsp_fast_flood = "%s"
		lsp_gen_init_intvl = "%s"
		lsp_gen_max_intvl = "%s"
		lsp_gen_sec_intvl = "%s"
		spf_comp_init_intvl = "%s"
		spf_comp_max_intvl = "%s"
		spf_comp_sec_intvl = "%s"
		isis_level_name = "%s"
		isis_level_type = "%s"
		
	}
	`, isisDomPol.Annotation,
		isisDomPol.Mtu,
		isisDomPol.RedistribMetric,
		isisDomPol.Description,
		isisDomPol.NameAlias,
		isisLvlComp.LspFastFlood,
		isisLvlComp.LspGenInitIntvl,
		isisLvlComp.LspGenMaxIntvl,
		isisLvlComp.LspGenSecIntvl,
		isisLvlComp.SpfCompInitIntvl,
		isisLvlComp.SpfCompMaxIntvl,
		isisLvlComp.SpfCompSecIntvl,
		isisLvlComp.Name,
		isisLvlComp.ISISLevel_type,
	)

	return resource
}
