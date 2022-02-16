package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSPANSourcedestinationGroupMatchLabelDataSource_Basic(t *testing.T) {
	resourceName := "aci_span_sourcedestination_group_match_label.test"
	dataSourceName := "data.aci_span_sourcedestination_group_match_label.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	
	fvTenantName := makeTestVariable(acctest.RandString(5))
	spanSrcGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		ProviderFactories:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSPANSourcedestinationGroupMatchLabelDSWithoutRequired(fvTenantName, spanSrcGrpName, rName,"span_source_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSPANSourcedestinationGroupMatchLabelDSWithoutRequired(fvTenantName, spanSrcGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelConfigDataSource(fvTenantName, spanSrcGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "span_source_group_dn", resourceName, "span_source_group_dn",),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
					
				),
			},
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelDataSourceUpdate(fvTenantName, spanSrcGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			
			{
				Config:      CreateAccSPANSourcedestinationGroupMatchLabelDSWithInvalidParentDn(fvTenantName, spanSrcGrpName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			
			{
				Config: CreateAccSPANSourcedestinationGroupMatchLabelDataSourceUpdatedResource(fvTenantName, spanSrcGrpName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}


func CreateAccSPANSourcedestinationGroupMatchLabelConfigDataSource(fvTenantName, spanSrcGrpName, rName string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}

	data "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = aci_span_sourcedestination_group_match_label.test.name
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
	`, fvTenantName, spanSrcGrpName, rName)
	return resource
}

func CreateSPANSourcedestinationGroupMatchLabelDSWithoutRequired(fvTenantName, spanSrcGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing span_sourcedestination_group_match_label Data Source without ",attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "span_source_group_dn":
		rBlock += `
	data "aci_span_sourcedestination_group_match_label" "test" {
	#	span_source_group_dn  = aci_span_source_group.test.id
		name  = aci_span_sourcedestination_group_match_label.test.name
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
	#	name  = aci_span_sourcedestination_group_match_label.test.name
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock,fvTenantName, spanSrcGrpName, rName)
}

func CreateAccSPANSourcedestinationGroupMatchLabelDSWithInvalidParentDn(fvTenantName, spanSrcGrpName, rName string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}

	data "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "${aci_span_sourcedestination_group_match_label.test.name}_invalid"
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
	`, fvTenantName, spanSrcGrpName, rName)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelDataSourceUpdate(fvTenantName, spanSrcGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
	}

	data "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = aci_span_sourcedestination_group_match_label.test.name
		%s = "%s"
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
	`, fvTenantName, spanSrcGrpName, rName,key,value)
	return resource
}

func CreateAccSPANSourcedestinationGroupMatchLabelDataSourceUpdatedResource(fvTenantName, spanSrcGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing span_sourcedestination_group_match_label Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_span_sourcedestination_group_match_label" "test" {
		span_source_group_dn  = aci_span_source_group.test.id
		name  = aci_span_sourcedestination_group_match_label.test.name
		depends_on = [ aci_span_sourcedestination_group_match_label.test ]
	}
	`, fvTenantName, spanSrcGrpName, rName,key,value)
	return resource
}