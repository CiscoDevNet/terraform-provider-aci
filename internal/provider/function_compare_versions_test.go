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
                }
                output "test_inside_single_below_false" {
                    value = provider::aci::compare_versions("4.1(7g)", "inside", "4.2(7f)")
                }
                output "test_inside_single_true" {
                    value = provider::aci::compare_versions("4.2(7f)", "inside", "4.2(7f)")
                }
                output "test_inside_single_above_false" {
                    value = provider::aci::compare_versions("5.2(7f)", "inside", "4.2(7f)")
                }
                output "test_inside_single_above_with_-_true" {
                    value = provider::aci::compare_versions("5.2(7f)", "inside", "4.2(7f)-")
                }
                output "test_inside_single_ower_boundary_with_-_true" {
                    value = provider::aci::compare_versions("4.2(7f)", "inside", "4.2(7f)-")
                }
                output "test_inside_single_below_with_-_false" {
                    value = provider::aci::compare_versions("4.1(7f)", "inside", "4.2(7f)-")
                }
                output "test_inside_single_range_below_false" {
                    value = provider::aci::compare_versions("4.1(7g)", "inside", "4.2(7f)-4.2(7w)")
                }
                output "test_inside_single_range_lower_boundary_true" {
                    value = provider::aci::compare_versions("4.2(7f)", "inside", "4.2(7f)-4.2(7w)")
                }
                output "test_inside_single_range_true" {
                    value = provider::aci::compare_versions("4.2(7g)", "inside", "4.2(7f)-4.2(7w)")
                }
                output "test_inside_single_range_upper_boundary_true" {
                    value = provider::aci::compare_versions("4.2(7w)", "inside", "4.2(7f)-4.2(7w)")
                }
                output "test_inside_single_range_above_false" {
                    value = provider::aci::compare_versions("4.2(7z)", "inside", "4.2(7f)-4.2(7w)")
                }
                output "test_inside_double_range_first_below_false" {
                    value = provider::aci::compare_versions("4.1(7g)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_first_lower_boundary_true" {
                    value = provider::aci::compare_versions("4.2(7f)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_first_true" {
                    value = provider::aci::compare_versions("4.2(7g)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_first_upper_boundary_true" {
                    value = provider::aci::compare_versions("4.2(7w)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_first_above_false" {
                    value = provider::aci::compare_versions("4.2(7z)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_second_below_false" {
                    value = provider::aci::compare_versions("5.2(1a)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_second_lower_boundary_true" {
                    value = provider::aci::compare_versions("5.2(1g)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_second_true" {
                    value = provider::aci::compare_versions("5.2(1h)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_second_upper_boundary_true" {
                    value = provider::aci::compare_versions("5.2(1k)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_double_range_second_above_false" {
                    value = provider::aci::compare_versions("5.2(1z)", "inside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_single_true" {
                    value = provider::aci::compare_versions("4.2(7f)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_single_false" {
                    value = provider::aci::compare_versions("4.1(7f)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_false" {
                    value = provider::aci::compare_versions("5.2(1a)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_lower_boundary_true" {
                    value = provider::aci::compare_versions("5.2(1g)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_true" {
                    value = provider::aci::compare_versions("5.2(1h)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_upper_boundary_true" {
                    value = provider::aci::compare_versions("5.2(1k)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_inside_single_first_range_second_above_false" {
                    value = provider::aci::compare_versions("5.2(1z)", "inside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_range_below_true" {
                    value = provider::aci::compare_versions("4.1(7g)", "outside", "4.2(7f)-4.2(7w)")
                }
                output "test_outside_single_range_lower_boundary_false" {
                    value = provider::aci::compare_versions("4.2(7f)", "outside", "4.2(7f)-4.2(7w)")
                }
                output "test_outside_single_range_false" {
                    value = provider::aci::compare_versions("4.2(7g)", "outside", "4.2(7f)-4.2(7w)")
                }
                output "test_outside_single_range_upper_boundary_false" {
                    value = provider::aci::compare_versions("4.2(7w)", "outside", "4.2(7f)-4.2(7w)")
                }
                output "test_outside_single_range_above_true" {
                    value = provider::aci::compare_versions("4.2(7z)", "outside", "4.2(7f)-4.2(7w)")
                }
                output "test_outside_double_range_first_below_true" {
                    value = provider::aci::compare_versions("4.1(7g)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_first_lower_boundary_false" {
                    value = provider::aci::compare_versions("4.2(7f)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_first_false" {
                    value = provider::aci::compare_versions("4.2(7g)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_first_upper_boundary_false" {
                    value = provider::aci::compare_versions("4.2(7w)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_first_above_true" {
                    value = provider::aci::compare_versions("4.2(7z)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_second_below_true" {
                    value = provider::aci::compare_versions("5.2(1a)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_second_lower_boundary_false" {
                    value = provider::aci::compare_versions("5.2(1g)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_second_false" {
                    value = provider::aci::compare_versions("5.2(1h)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_second_upper_boundary_false" {
                    value = provider::aci::compare_versions("5.2(1k)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_double_range_second_above_true" {
                    value = provider::aci::compare_versions("5.2(1z)", "outside", "4.2(7f)-4.2(7w),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_single_false" {
                    value = provider::aci::compare_versions("4.2(7f)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_single_true" {
                    value = provider::aci::compare_versions("4.1(7f)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_true" {
                    value = provider::aci::compare_versions("5.2(1a)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_lower_boundary_false" {
                    value = provider::aci::compare_versions("5.2(1g)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_false" {
                    value = provider::aci::compare_versions("5.2(1h)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_upper_boundary_false" {
                    value = provider::aci::compare_versions("5.2(1k)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                output "test_outside_single_first_range_second_above_true" {
                    value = provider::aci::compare_versions("5.2(1z)", "outside", "4.2(7f),5.2(1g)-5.2(1k)")
                }
                `,
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
					CheckOutputBool("test_inside_single_below_false", false),
					CheckOutputBool("test_inside_single_true", true),
					CheckOutputBool("test_inside_single_above_false", false),
					CheckOutputBool("test_inside_single_above_with_-_true", true),
					CheckOutputBool("test_inside_single_ower_boundary_with_-_true", true),
					CheckOutputBool("test_inside_single_below_with_-_false", false),
					CheckOutputBool("test_inside_single_range_below_false", false),
					CheckOutputBool("test_inside_single_range_lower_boundary_true", true),
					CheckOutputBool("test_inside_single_range_true", true),
					CheckOutputBool("test_inside_single_range_upper_boundary_true", true),
					CheckOutputBool("test_inside_single_range_above_false", false),
					CheckOutputBool("test_inside_double_range_first_below_false", false),
					CheckOutputBool("test_inside_double_range_first_lower_boundary_true", true),
					CheckOutputBool("test_inside_double_range_first_true", true),
					CheckOutputBool("test_inside_double_range_first_upper_boundary_true", true),
					CheckOutputBool("test_inside_double_range_first_above_false", false),
					CheckOutputBool("test_inside_double_range_second_below_false", false),
					CheckOutputBool("test_inside_double_range_second_lower_boundary_true", true),
					CheckOutputBool("test_inside_double_range_second_true", true),
					CheckOutputBool("test_inside_double_range_second_upper_boundary_true", true),
					CheckOutputBool("test_inside_double_range_second_above_false", false),
					CheckOutputBool("test_inside_single_first_range_second_single_true", true),
					CheckOutputBool("test_inside_single_first_range_second_single_false", false),
					CheckOutputBool("test_inside_single_first_range_second_false", false),
					CheckOutputBool("test_inside_single_first_range_second_lower_boundary_true", true),
					CheckOutputBool("test_inside_single_first_range_second_true", true),
					CheckOutputBool("test_inside_single_first_range_second_upper_boundary_true", true),
					CheckOutputBool("test_inside_single_first_range_second_above_false", false),
					CheckOutputBool("test_outside_single_range_below_true", true),
					CheckOutputBool("test_outside_single_range_lower_boundary_false", false),
					CheckOutputBool("test_outside_single_range_false", false),
					CheckOutputBool("test_outside_single_range_upper_boundary_false", false),
					CheckOutputBool("test_outside_single_range_above_true", true),
					CheckOutputBool("test_outside_double_range_first_below_true", true),
					CheckOutputBool("test_outside_double_range_first_lower_boundary_false", false),
					CheckOutputBool("test_outside_double_range_first_false", false),
					CheckOutputBool("test_outside_double_range_first_upper_boundary_false", false),
					CheckOutputBool("test_outside_double_range_first_above_true", true),
					CheckOutputBool("test_outside_double_range_second_below_true", true),
					CheckOutputBool("test_outside_double_range_second_lower_boundary_false", false),
					CheckOutputBool("test_outside_double_range_second_false", false),
					CheckOutputBool("test_outside_double_range_second_upper_boundary_false", false),
					CheckOutputBool("test_outside_double_range_second_above_true", true),
					CheckOutputBool("test_outside_single_first_range_second_single_false", false),
					CheckOutputBool("test_outside_single_first_range_second_single_true", true),
					CheckOutputBool("test_outside_single_first_range_second_true", true),
					CheckOutputBool("test_outside_single_first_range_second_lower_boundary_false", false),
					CheckOutputBool("test_outside_single_first_range_second_false", false),
					CheckOutputBool("test_outside_single_first_range_second_upper_boundary_false", false),
					CheckOutputBool("test_outside_single_first_range_second_above_true", true),
				),
			},
		},
	})
}
