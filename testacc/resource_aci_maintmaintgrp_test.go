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

func TestAccAciPODMaintenanceGroup_Basic(t *testing.T) {
	var pod_maintenance_group_default models.PODMaintenanceGroup
	var pod_maintenance_group_updated models.PODMaintenanceGroup
	resourceName := "aci_pod_maintenance_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePODMaintenanceGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPODMaintenanceGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "switch"),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_type", "range"),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_pod_maintenance_group"),

					resource.TestCheckResourceAttr(resourceName, "fwtype", "catalog"),

					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_type", "ALL"),

					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccPODMaintenanceGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccPODMaintenanceGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccPODMaintenanceGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciPODMaintenanceGroupIdNotEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
		},
	})
}

func TestAccAciPODMaintenanceGroup_Update(t *testing.T) {
	var pod_maintenance_group_default models.PODMaintenanceGroup
	var pod_maintenance_group_updated models.PODMaintenanceGroup
	resourceName := "aci_pod_maintenance_group.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPODMaintenanceGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_default),
				),
			},

			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", "config"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "config"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", "controller"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "controller"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", "plugin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "plugin"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", "pluginPackage"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "pluginPackage"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", "vpod"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "fwtype", "vpod"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "pod_maintenance_group_type", "ALL_IN_POD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_type", "ALL_IN_POD"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupUpdatedAttr(rName, "pod_maintenance_group_type", "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists(resourceName, &pod_maintenance_group_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_type", "range"),
					testAccCheckAciPODMaintenanceGroupIdEqual(&pod_maintenance_group_default, &pod_maintenance_group_updated),
				),
			},
			{
				Config: CreateAccPODMaintenanceGroupConfig(rName),
			},
		},
	})
}

func TestAccAciPODMaintenanceGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPODMaintenanceGroupConfig(rName),
			},

			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, "fwtype", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, "pod_maintenance_group_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPODMaintenanceGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPODMaintenanceGroupConfig(rName),
			},
		},
	})
}

func TestAccAciPODMaintenanceGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPODMaintenanceGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciPODMaintenanceGroupExists(name string, pod_maintenance_group *models.PODMaintenanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Pod Maintenance Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pod Maintenance Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		pod_maintenance_groupFound := models.PODMaintenanceGroupFromContainer(cont)
		if pod_maintenance_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Pod Maintenance Group %s not found", rs.Primary.ID)
		}
		*pod_maintenance_group = *pod_maintenance_groupFound
		return nil
	}
}

func testAccCheckAciPODMaintenanceGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing pod_maintenance_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_pod_maintenance_group" {
			cont, err := client.Get(rs.Primary.ID)
			pod_maintenance_group := models.PODMaintenanceGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Pod Maintenance Group %s Still exists", pod_maintenance_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPODMaintenanceGroupIdEqual(m1, m2 *models.PODMaintenanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("pod_maintenance_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciPODMaintenanceGroupIdNotEqual(m1, m2 *models.PODMaintenanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("pod_maintenance_group DNs are equal")
		}
		return nil
	}
}

func CreatePODMaintenanceGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing pod_maintenance_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_pod_maintenance_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPODMaintenanceGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccPODMaintenanceGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPODMaintenanceGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPODMaintenanceGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple pod_maintenance_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccPODMaintenanceGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing pod_maintenance_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_pod_maintenance_group"
		fwtype = "catalog"
		pod_maintenance_group_type = "ALL"
		
	}
	`, rName)

	return resource
}

func CreateAccPODMaintenanceGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing pod_maintenance_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_pod_maintenance_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_pod_maintenance_group"
		fwtype = "catalog"
		pod_maintenance_group_type = "ALL"
		
	}
	`)

	return resource
}

func CreateAccPODMaintenanceGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing pod_maintenance_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
