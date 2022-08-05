package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/terraform-providers/terraform-provider-aci/aci"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVPCExplicitProtectionGroup_Basic(t *testing.T) {
	var vpc_explicit_protection_group_default models.VPCExplicitProtectionGroup
	var vpc_explicit_protection_group_updated models.VPCExplicitProtectionGroup
	resourceName := "aci_vpc_explicit_protection_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	switch1 := "503"
	switch2 := "504"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVPCExplicitProtectionGroupWithoutRequired(switch1, switch2, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVPCExplicitProtectionGroupWithoutRequired(switch1, switch2, rName, "switch1"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVPCExplicitProtectionGroupWithoutRequired(switch1, switch2, rName, "switch2"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(switch1, switch2, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_default),
					resource.TestCheckResourceAttr(resourceName, "switch1", switch1),
					resource.TestCheckResourceAttr(resourceName, "switch2", switch2),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpc_explicit_protection_group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpc_domain_policy", ""),
				),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigWithOptionalValues(switch1, switch2, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),
					resource.TestCheckResourceAttr(resourceName, "switch1", switch1),
					resource.TestCheckResourceAttr(resourceName, "switch2", switch2),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "vpc_explicit_protection_group_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_domain_policy", "test"),
					testAccCheckAciVPCExplicitProtectionGroupIdEqual(&vpc_explicit_protection_group_default, &vpc_explicit_protection_group_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"switch1", "switch2"},
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(switch2, switch1, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "switch1", switch2),
					resource.TestCheckResourceAttr(resourceName, "switch2", switch1),
					testAccCheckAciVPCExplicitProtectionGroupIdEqual(&vpc_explicit_protection_group_default, &vpc_explicit_protection_group_updated),
				),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(switch1, switch2, rName),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigUpdatedName(switch1, switch2, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(rName, switch2, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(switch1, rName, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams("0", switch2, rName),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams("16001", switch2, rName),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(switch1, "0", rName),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(switch1, "16001", rName),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(switch1, switch2, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVPCExplicitProtectionGroupIdNotEqual(&vpc_explicit_protection_group_default, &vpc_explicit_protection_group_updated),
				),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_Update(t *testing.T) {
	var vpc_explicit_protection_group_default models.VPCExplicitProtectionGroup
	var vpc_explicit_protection_group_updated models.VPCExplicitProtectionGroup
	resourceName := "aci_vpc_explicit_protection_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	switch1 := "505"
	switch2 := "506"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(switch1, switch2, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_default),
				),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_explicit_protection_group_id", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),
					resource.TestCheckResourceAttr(resourceName, "vpc_explicit_protection_group_id", "1000"),
				),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_explicit_protection_group_id", "500"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),
					resource.TestCheckResourceAttr(resourceName, "vpc_explicit_protection_group_id", "500"),
				),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	switch1 := "507"
	switch2 := "508"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(switch1, switch2, rName),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_explicit_protection_group_id", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_explicit_protection_group_id", "1001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_explicit_protection_group_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, "vpc_domain_policy", acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(switch1, switch2, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(switch1, switch2, rName),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciVPCExplicitProtectionGroupExists(name string, vpc_explicit_protection_group *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Explicit Protection Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		vpc_explicit_protection_groupFound, err := aci.GetRemoteVPCExplicitProtectionGroup(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if vpc_explicit_protection_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", rs.Primary.ID)
		}
		*vpc_explicit_protection_group = *vpc_explicit_protection_groupFound
		return nil
	}
}

func testAccCheckAciVPCExplicitProtectionGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vpc_explicit_protection_group" {
			cont, err := client.Get(rs.Primary.ID)
			vpc_explicit_protection_group := models.VPCExplicitProtectionGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VPC Explicit Protection Group %s Still exists", vpc_explicit_protection_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVPCExplicitProtectionGroupIdEqual(m1, m2 *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vpc_explicit_protection_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVPCExplicitProtectionGroupIdNotEqual(m1, m2 *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vpc_explicit_protection_group DNs are equal")
		}
		return nil
	}
}

func CreateVPCExplicitProtectionGroupWithoutRequired(sw1, sw2, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group creation without ", attrName)
	rBlock := `

	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vpc_explicit_protection_group" "test" {
	#	name  = "%s"
		switch1 = "%s"
		switch2 = "%s"
	}
		`
	case "switch1":
		rBlock += `
	resource "aci_vpc_explicit_protection_group" "test" {
		name  = "%s"
	#	switch1 = "%s"
		switch2 = "%s"
	}
		`
	case "switch2":
		rBlock += `
	resource "aci_vpc_explicit_protection_group" "test" {
		name  = "%s"
		switch1 = "%s"
	#	switch2 = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, sw1, sw2)
}

func CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(sw1, sw2, rName string) string {
	fmt.Printf("=== STEP  testing vpc_explicit_protection_group creation with switch1 %s,switch2 %s and resource name %s\n", sw1, sw2, rName)
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}
	`, sw1, sw2, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfigUpdatedName(sw1, sw2, rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}
	`, sw1, sw2, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfig(sw1, sw2, rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}
	`, sw1, sw2, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple vpc_explicit_protection_group creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "51${count.index*2}"
		switch2 = "51${count.index*2+1}"
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfigWithOptionalValues(sw1, sw2, rName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		vpc_domain_policy = "test"
		vpc_explicit_protection_group_id = "1"
	}
	`, sw1, sw2, rName)

	return resource
}

func CreateAccVPCExplicitProtectionGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vpc_explicit_protection_group" "test" {
		annotation = "orchestrator:terraform_testacc"
		vpc_domain_policy = "test"
		vpc_explicit_protection_group_id = "1"
	}
	`)

	return resource
}

func CreateAccVPCExplicitProtectionGroupUpdatedAttr(sw1, sw2, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vpc_explicit_protection_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
		%s = "%s"
	}
	`, sw1, sw2, rName, attribute, value)
	return resource
}
