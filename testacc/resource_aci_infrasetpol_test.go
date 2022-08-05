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

func TestAccAciFabricWideSettingsPolicy_Basic(t *testing.T) {
	var fabric_wide_settings_policy_default models.FabricWideSettingsPolicy
	var fabric_wide_settings_policy_updated models.FabricWideSettingsPolicy
	resourceName := "aci_fabric_wide_settings.test"
	fabricWideSettings, err := aci.GetRemoteFabricWideSettingsPolicy(sharedAciClient(), "uni/infra/settings")
	if err != nil {
		t.Errorf("reading initial config of fabricWideSettings")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_default),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_wide_settings_policy"),
					resource.TestCheckResourceAttr(resourceName, "disable_ep_dampening", "yes"),
					resource.TestCheckResourceAttr(resourceName, "enable_mo_streaming", "yes"),
					resource.TestCheckResourceAttr(resourceName, "enforce_subnet_check", "yes"),
					resource.TestCheckResourceAttr(resourceName, "opflexp_authenticate_clients", "yes"),
					resource.TestCheckResourceAttr(resourceName, "opflexp_use_ssl", "yes"),
					resource.TestCheckResourceAttr(resourceName, "restrict_infra_vlan_traffic", "yes"),
					resource.TestCheckResourceAttr(resourceName, "unicast_xr_ep_learn_disable", "yes"),
					resource.TestCheckResourceAttr(resourceName, "validate_overlapping_vlans", "yes"),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_remote_leaf_direct"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: RestoreFabricWideSetting(fabricWideSettings),
			},
		},
	})
}

func TestAccAciFabricWideSettingsPolicy_Update(t *testing.T) {
	var fabric_wide_settings_policy_default models.FabricWideSettingsPolicy
	var fabric_wide_settings_policy_updated models.FabricWideSettingsPolicy
	resourceName := "aci_fabric_wide_settings.test"
	fabricWideSettings, err := aci.GetRemoteFabricWideSettingsPolicy(sharedAciClient(), "uni/infra/settings")
	if err != nil {
		t.Errorf("reading initial config of fabricWideSettings")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_default),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("disable_ep_dampening", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "disable_ep_dampening", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("enable_mo_streaming", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "enable_mo_streaming", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("enforce_subnet_check", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "enforce_subnet_check", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_authenticate_clients", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "opflexp_authenticate_clients", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_use_ssl", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "opflexp_use_ssl", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("restrict_infra_vlan_traffic", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "restrict_infra_vlan_traffic", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("unicast_xr_ep_learn_disable", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "unicast_xr_ep_learn_disable", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyUpdatedAttr("validate_overlapping_vlans", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "validate_overlapping_vlans", "no"),
					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_policy_default, &fabric_wide_settings_policy_updated),
				),
			},
			{
				Config: RestoreFabricWideSetting(fabricWideSettings),
			},
		},
	})
}

func TestAccAciFabricWideSettingsPolicy_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	fabricWideSettings, err := aci.GetRemoteFabricWideSettingsPolicy(sharedAciClient(), "uni/infra/settings")
	if err != nil {
		t.Errorf("reading initial config of fabricWideSettings")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("disable_ep_dampening", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enable_mo_streaming", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enable_remote_leaf_direct", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enforce_subnet_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("name", acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_authenticate_clients", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_use_ssl", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("restrict_infra_vlan_traffic", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("unicast_xr_ep_learn_disable", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("validate_overlapping_vlans", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: RestoreFabricWideSetting(fabricWideSettings),
			},
		},
	})
}

func RestoreFabricWideSetting(m *models.FabricWideSettingsPolicy) string {
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		disable_ep_dampening = "%s"
		enable_mo_streaming = "%s"
		enforce_subnet_check = "%s"
		opflexp_authenticate_clients = "%s"
		opflexp_use_ssl = "%s"
		restrict_infra_vlan_traffic = "%s"
		unicast_xr_ep_learn_disable = "%s"
		validate_overlapping_vlans = "%s"
	}
	`, m.Description, m.Annotation, m.NameAlias, m.DisableEpDampening, m.EnableMoStreaming, m.EnforceSubnetCheck, m.OpflexpAuthenticateClients, m.OpflexpUseSsl, m.RestrictInfraVLANTraffic, m.UnicastXrEpLearnDisable, m.ValidateOverlappingVlans)

	return resource
}

func testAccCheckAciFabricWideSettingsPolicyExists(name string, fabric_wide_settings_policy *models.FabricWideSettingsPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Wide Settings Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Wide Settings Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_wide_settings_policyFound := models.FabricWideSettingsPolicyFromContainer(cont)
		if fabric_wide_settings_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Wide Settings Policy %s not found", rs.Primary.ID)
		}
		*fabric_wide_settings_policy = *fabric_wide_settings_policyFound
		return nil
	}
}

func testAccCheckAciFabricWideSettingsPolicyIdEqual(m1, m2 *models.FabricWideSettingsPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_wide_settings_policy DNs are not equal")
		}
		return nil
	}
}

func CreateAccFabricWideSettingsPolicyConfig() string {
	fmt.Println("=== STEP  testing fabric_wide_settings_policy creation")
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {}
	`)
	return resource
}

func CreateAccFabricWideSettingsPolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing fabric_wide_settings_policy creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_wide_settings_policy"
		disable_ep_dampening = "yes"
		enable_mo_streaming = "yes"
		enforce_subnet_check = "yes"
		opflexp_authenticate_clients = "yes"
		opflexp_use_ssl = "yes"
		restrict_infra_vlan_traffic = "yes"
		unicast_xr_ep_learn_disable = "yes"
		validate_overlapping_vlans = "yes"
		name = "test"
	}
	`)

	return resource
}

func CreateAccFabricWideSettingsPolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_wide_settings_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {

		%s = "%s"
	}
	`, attribute, value)
	return resource
}
