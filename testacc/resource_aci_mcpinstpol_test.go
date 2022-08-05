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

func TestAccAciMCPInstancePolicy_Basic(t *testing.T) {
	var mcp_instance_policy_default models.MiscablingProtocolInstancePolicy
	var mcp_instance_policy_updated models.MiscablingProtocolInstancePolicy
	resourceName := "aci_mcp_instance_policy.test"
	key := acctest.RandString(5)
	mcpInstancePolicy, err := aci.GetRemoteMiscablingProtocolInstancePolicy(sharedAciClient(), "uni/infra/mcpInstP-default")
	if err != nil {
		t.Errorf("reading initial config of MCP Instance Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateMCPInstancePolicyWithoutRequired(key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMCPInstancePolicyConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_default),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttrSet(resourceName, "admin_st"),
					resource.TestCheckResourceAttrSet(resourceName, "init_delay_time"),
					resource.TestCheckResourceAttrSet(resourceName, "loop_detect_mult"),
					resource.TestCheckResourceAttrSet(resourceName, "loop_protect_act"),
					resource.TestCheckResourceAttrSet(resourceName, "loop_protect_act"),
					resource.TestCheckResourceAttrSet(resourceName, "tx_freq"),
					resource.TestCheckResourceAttrSet(resourceName, "tx_freq_msec"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyConfigWithOptionalValues(key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_mcp_instance_policy"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "1"),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "0"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "100"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config:      CreateAccMCPInstancePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMCPInstancePolicyInitialConfig(key, mcpInstancePolicy),
			},
		},
	})
}

func TestAccAciMCPInstancePolicy_Update(t *testing.T) {
	var mcp_instance_policy_default models.MiscablingProtocolInstancePolicy
	var mcp_instance_policy_updated models.MiscablingProtocolInstancePolicy
	resourceName := "aci_mcp_instance_policy.test"
	key := acctest.RandString(5)
	mcpInstancePolicy, err := aci.GetRemoteMiscablingProtocolInstancePolicy(sharedAciClient(), "uni/infra/mcpInstP-default")
	if err != nil {
		t.Errorf("reading initial config of MCP Instance Policy")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMCPInstancePolicyConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_default),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{"pdu-per-vlan"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{"pdu-per-vlan", "stateful-ha"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "stateful-ha"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{"stateful-ha"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "stateful-ha"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{"stateful-ha", "pdu-per-vlan"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "stateful-ha"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "pdu-per-vlan"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", "[]"),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "admin_st", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "init_delay_time", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "1800"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "init_delay_time", "900"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "900"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_detect_mult", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "255"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_detect_mult", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "127"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_protect_act", "port-disable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act", "port-disable"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", "500"),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq", "299"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "299"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "150"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", "999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "999"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", "499"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMCPInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "499"),
					testAccCheckAciMCPInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyInitialConfig(key, mcpInstancePolicy),
			},
		},
	})
}

func TestAccAciMCPInstancePolicy_Negative(t *testing.T) {
	key := acctest.RandString(5)
	mcpInstancePolicy, err := aci.GetRemoteMiscablingProtocolInstancePolicy(sharedAciClient(), "uni/infra/mcpInstP-default")
	if err != nil {
		t.Errorf("reading initial config of MCP Instance Policy")
	}
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMCPInstancePolicyConfig(key),
			},

			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttrList(key, "ctrl", StringListtoString([]string{"pdu-per-vlan", "pdu-per-vlan"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "init_delay_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "init_delay_time", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "init_delay_time", "1801"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_detect_mult", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_detect_mult", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "loop_protect_act", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, "tx_freq_msec", "1000"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMCPInstancePolicyUpdatedAttr(key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMCPInstancePolicyInitialConfig(key, mcpInstancePolicy),
			},
		},
	})
}

func testAccCheckAciMCPInstancePolicyExists(name string, mcp_instance_policy *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("MCP Instance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MCP Instance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mcp_instance_policyFound := models.MiscablingProtocolInstancePolicyFromContainer(cont)
		if mcp_instance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("MCP Instance Policy %s not found", rs.Primary.ID)
		}
		*mcp_instance_policy = *mcp_instance_policyFound
		return nil
	}
}

func testAccCheckAciMCPInstancePolicyIdEqual(m1, m2 *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("mcp_instance_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMCPInstancePolicyIdNotEqual(m1, m2 *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("mcp_instance_policy DNs are equal")
		}
		return nil
	}
}

func CreateMCPInstancePolicyWithoutRequired(key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "key":
		rBlock += `
	resource "aci_mcp_instance_policy" "test" {
	
	#	key  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, key)
}

func CreateAccMCPInstancePolicyConfig(key string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy creation")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
	}
	`, key)
	return resource
}

func CreateAccMCPInstancePolicyConfigWithOptionalValues(key string) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mcp_instance_policy"
		admin_st = "enabled"
		ctrl = ["pdu-per-vlan"]
		init_delay_time = "0"
		loop_detect_mult = "1"
		loop_protect_act = "port-disable"
		tx_freq = "0"
		tx_freq_msec = "100"	
	}
	`, key)

	return resource
}

func CreateAccMCPInstancePolicyInitialConfig(key string, mcpInstancePolicy *models.MiscablingProtocolInstancePolicy) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy creation with Initial Config")
	resource := fmt.Sprintf(`
	resource "aci_mcp_instance_policy" "test" {
		key  = "%s"
		description = "%s"
		annotation = "%s"
		name_alias = "%s"
		admin_st = "%s"
		ctrl = %s
		init_delay_time = "%s"
		loop_detect_mult = "%s"
		loop_protect_act = "%s"
		tx_freq = "%s"
		tx_freq_msec = "%s"
	}
	`, key, mcpInstancePolicy.Description, mcpInstancePolicy.Annotation, mcpInstancePolicy.NameAlias, mcpInstancePolicy.AdminSt, StringListtoString(convertToStringArray(mcpInstancePolicy.Ctrl)), mcpInstancePolicy.InitDelayTime, mcpInstancePolicy.LoopDetectMult, mcpInstancePolicy.LoopProtectAct, mcpInstancePolicy.TxFreq, mcpInstancePolicy.TxFreqMsec)
	return resource
}

func CreateAccMCPInstancePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_mcp_instance_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mcp_instance_policy"
		admin_st = "enabled"
		ctrl = ["pdu-per-vlan"]
		init_delay_time = "1"
		loop_detect_mult = "2"
		loop_protect_act = "port-disable"
		tx_freq = "1"
		tx_freq_msec = "1"
	}
	`)

	return resource
}

func CreateAccMCPInstancePolicyUpdatedAttr(key, attribute, value string) string {
	fmt.Printf("=== STEP  testing mcp_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
		%s = "%s"
	}
	`, key, attribute, value)
	return resource
}

func CreateAccMCPInstancePolicyUpdatedAttrList(key, attribute, value string) string {
	fmt.Printf("=== STEP  testing mcp_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
		%s = %s
	}
	`, key, attribute, value)
	return resource
}
