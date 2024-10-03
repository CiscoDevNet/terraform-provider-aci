package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestCompareVersionsFunction(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesFunction,
		Steps: []resource.TestStep{
			{
				Config: `
                output "test_equal" {
                    value = provider::aci::compare_versions("1.0(0a)", "==", "1.0(0a)")
                }
                output "test_greater" {
                    value = provider::aci::compare_versions("1.1(0a)", ">", "1.0(0a)")
                }
                output "test_lesser" {
                    value = provider::aci::compare_versions("1.0(0a)", "<", "1.1(0a)")
                }
                output "test_not_equal" {
                    value = provider::aci::compare_versions("1.0(0a)", "!=", "1.1(0a)")
                }
                output "test_greater_or_equal" {
                    value = provider::aci::compare_versions("1.1(0a)", ">=", "1.0(0a)")
                }
                output "test_less_or_equal" {
                    value = provider::aci::compare_versions("1.0(0a)", "<=", "1.1(0a)")
                }
                output "test_not_equal_false" {
                    value = provider::aci::compare_versions("1.0(0a)", "!=", "1.0(0a)")
                }
                output "test_greater_false" {
                    value = provider::aci::compare_versions("1.0(0a)", ">", "1.1(0a)")
                }
                output "test_lesser_false" {
                    value = provider::aci::compare_versions("1.1(0a)", "<", "1.0(0a)")
                }
                output "test_greater_or_equal_false" {
                    value = provider::aci::compare_versions("1.0(0a)", ">=", "1.1(0a)")
                }
                output "test_less_or_equal_false" {
                    value = provider::aci::compare_versions("1.1(0a)", "<=", "1.0(0a)")
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					CheckOutputBool("test_equal", true),
					CheckOutputBool("test_greater", true),
					CheckOutputBool("test_lesser", true),
					CheckOutputBool("test_not_equal", true),
					CheckOutputBool("test_greater_or_equal", true),
					CheckOutputBool("test_less_or_equal", true),
					CheckOutputBool("test_not_equal_false", false),
					CheckOutputBool("test_greater_false", false),
					CheckOutputBool("test_lesser_false", false),
					CheckOutputBool("test_greater_or_equal_false", false),
					CheckOutputBool("test_less_or_equal_false", false),
				),
			},
		},
	})
}
