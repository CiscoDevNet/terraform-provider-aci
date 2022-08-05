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

func TestAccAciQOSInstancePolicy_Basic(t *testing.T) {
	var qos_instance_policy_default models.QOSInstancePolicy
	var qos_instance_policy_updated models.QOSInstancePolicy
	resourceName := "aci_qos_instance_policy.test"
	QOSInstancePolicy, err := aci.GetRemoteQOSInstancePolicy(sharedAciClient(), "uni/infra/qosinst-default")
	if err != nil {
		t.Errorf("reading initial config of QOSInstancePolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "etrap_age_timer"),
					resource.TestCheckResourceAttrSet(resourceName, "etrap_bw_thresh"),
					resource.TestCheckResourceAttrSet(resourceName, "etrap_byte_ct"),
					resource.TestCheckResourceAttrSet(resourceName, "etrap_st"),
					resource.TestCheckResourceAttrSet(resourceName, "ctrl"),
					resource.TestCheckResourceAttrSet(resourceName, "fabric_flush_interval"),
					resource.TestCheckResourceAttrSet(resourceName, "fabric_flush_st"),
					resource.TestCheckResourceAttrSet(resourceName, "uburst_spine_queues"),
					resource.TestCheckResourceAttrSet(resourceName, "uburst_tor_queues"),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_qos_instance_policy"),
					resource.TestCheckResourceAttr(resourceName, "etrap_st", "yes"),
					resource.TestCheckResourceAttr(resourceName, "etrap_age_timer", "0"),
					resource.TestCheckResourceAttr(resourceName, "etrap_bw_thresh", "0"),
					resource.TestCheckResourceAttr(resourceName, "etrap_byte_ct", "0"),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "100"),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_st", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "none"),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "0"),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "0"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restoreQOSInstancePolicy(QOSInstancePolicy),
			},
		},
	})
}

func TestAccAciQOSInstancePolicy_Update(t *testing.T) {
	var qos_instance_policy_default models.QOSInstancePolicy
	var qos_instance_policy_updated models.QOSInstancePolicy
	resourceName := "aci_qos_instance_policy.test"
	QOSInstancePolicy, err := aci.GetRemoteQOSInstancePolicy(sharedAciClient(), "uni/infra/qosinst-default")
	if err != nil {
		t.Errorf("reading initial config of QOSInstancePolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_default),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("etrap_age_timer", "500"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "etrap_age_timer", "500"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("etrap_bw_thresh", "500"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "etrap_bw_thresh", "500"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("etrap_byte_ct", "500"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "etrap_byte_ct", "500"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("etrap_st", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "etrap_st", "no"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_interval", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "1000"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_interval", "550"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "550"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_st", "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_st", "no"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("uburst_spine_queues", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "100"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("uburst_spine_queues", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "50"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("uburst_tor_queues", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "100"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("uburst_tor_queues", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "50"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr("ctrl", "none"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "none"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: restoreQOSInstancePolicy(QOSInstancePolicy),
			},
		},
	})
}

func TestAccAciQOSInstancePolicy_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	QOSInstancePolicy, err := aci.GetRemoteQOSInstancePolicy(sharedAciClient(), "uni/infra/qosinst-default")
	if err != nil {
		t.Errorf("reading initial config of QOSInstancePolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfig(),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_age_timer", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_age_timer", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_bw_thresh", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_bw_thresh", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_byte_ct", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_byte_ct", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("etrap_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_interval", "99"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_interval", "1001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("fabric_flush_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_spine_queues", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_spine_queues", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_spine_queues", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_tor_queues", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_tor_queues", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr("uburst_tor_queues", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restoreQOSInstancePolicy(QOSInstancePolicy),
			},
		},
	})
}

func testAccCheckAciQOSInstancePolicyExists(name string, qos_instance_policy *models.QOSInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("QOS Instance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No QOS Instance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		qos_instance_policyFound := models.QOSInstancePolicyFromContainer(cont)
		if qos_instance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("QOS Instance Policy %s not found", rs.Primary.ID)
		}
		*qos_instance_policy = *qos_instance_policyFound
		return nil
	}
}

func testAccCheckAciQOSInstancePolicyIdEqual(m1, m2 *models.QOSInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("qos_instance_policy DNs are not equal")
		}
		return nil
	}
}

func restoreQOSInstancePolicy(QOSInstancePolicy *models.QOSInstancePolicy) string {
	var resource string
	if QOSInstancePolicy.Ctrl == "dot1p-preserve" {
		resource = fmt.Sprintf(`
		resource "aci_qos_instance_policy" "test" {
				name_alias            = "%s"
				description           = "%s"
				etrap_age_timer       = "%s" 
				etrap_bw_thresh       = "%s"
				etrap_byte_ct         = "%s"
				etrap_st              = "%s"
				fabric_flush_interval = "%s"
				fabric_flush_st       = "%s"
				annotation            = "%s"
				ctrl                  = "dot1p-preserve"
				uburst_spine_queues   = "%s"
				uburst_tor_queues     = "%s"
			  }
		`, QOSInstancePolicy.NameAlias, QOSInstancePolicy.Description, QOSInstancePolicy.EtrapAgeTimer, QOSInstancePolicy.EtrapBwThresh, QOSInstancePolicy.EtrapByteCt, QOSInstancePolicy.EtrapSt, QOSInstancePolicy.FabricFlushInterval, QOSInstancePolicy.FabricFlushSt, QOSInstancePolicy.Annotation, QOSInstancePolicy.UburstSpineQueues, QOSInstancePolicy.UburstTorQueues)
	} else {
		resource = fmt.Sprintf(`
		resource "aci_qos_instance_policy" "test" {
				name_alias            = "%s"
				description           = "%s"
				etrap_age_timer       = "%s" 
				etrap_bw_thresh       = "%s"
				etrap_byte_ct         = "%s"
				etrap_st              = "%s"
				fabric_flush_interval = "%s"
				fabric_flush_st       = "%s"
				annotation            = "%s"
				ctrl                  = "none"
				uburst_spine_queues   = "%s"
				uburst_tor_queues     = "%s"
			  }
		`, QOSInstancePolicy.NameAlias, QOSInstancePolicy.Description, QOSInstancePolicy.EtrapAgeTimer, QOSInstancePolicy.EtrapBwThresh, QOSInstancePolicy.EtrapByteCt, QOSInstancePolicy.EtrapSt, QOSInstancePolicy.FabricFlushInterval, QOSInstancePolicy.FabricFlushSt, QOSInstancePolicy.Annotation, QOSInstancePolicy.UburstSpineQueues, QOSInstancePolicy.UburstTorQueues)
	}
	return resource
}

func CreateAccQOSInstancePolicyConfig() string {
	fmt.Println("=== STEP  testing qos_instance_policy creation")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	}
	`)
	return resource
}

func CreateAccQOSInstancePolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing qos_instance_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_qos_instance_policy"
		etrap_st = "yes"
		etrap_age_timer = "0"
		etrap_bw_thresh = "0"
		etrap_byte_ct = "0"
		fabric_flush_interval = "100"
		fabric_flush_st = "yes"
		ctrl = "none"
		uburst_spine_queues = "0"
		uburst_tor_queues = "0"
	}
	`)

	return resource
}

func CreateAccQOSInstancePolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing qos_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
