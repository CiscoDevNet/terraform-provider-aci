// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceL3extRsRedistributePolWithL3extOut(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "apic", "4.2(1i)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testConfigL3extRsRedistributePolDataSourceDependencyWithL3extOut,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_l3out_redistribute_policy.test", "route_control_profile_name", "test_tn_rtctrl_profile_name"),
					resource.TestCheckResourceAttr("data.aci_l3out_redistribute_policy.test", "source", "direct"),
					resource.TestCheckResourceAttr("data.aci_l3out_redistribute_policy.test", "annotation", "orchestrator:terraform"),
				),
			},
			{
				Config:      testConfigL3extRsRedistributePolNotExistingL3extOut,
				ExpectError: regexp.MustCompile("Failed to read aci_l3out_redistribute_policy data source"),
			},
		},
	})
}

const testConfigL3extRsRedistributePolDataSourceDependencyWithL3extOut = testConfigL3extRsRedistributePolMinDependencyWithL3extOut + `
data "aci_l3out_redistribute_policy" "test" {
  parent_dn = aci_l3_outside.test.id
  route_control_profile_name = "test_tn_rtctrl_profile_name"
  source = "direct"
  depends_on = [aci_l3out_redistribute_policy.test]
}
`

const testConfigL3extRsRedistributePolNotExistingL3extOut = testConfigL3extOutMinDependencyWithFvTenant + `
data "aci_l3out_redistribute_policy" "test_non_existing" {
  parent_dn = aci_l3_outside.test.id
  route_control_profile_name = "non_existing_tn_rtctrl_profile_name"
  source = "static"
}
`
