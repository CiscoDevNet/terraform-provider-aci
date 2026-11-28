package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/apicversion"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = CompareVersionsFunction{}

func NewCompareVersionsFunction() function.Function {
	return &CompareVersionsFunction{}
}

func init() {
	registerFunction("compare_versions", NewCompareVersionsFunction)
}

type CompareVersionsFunction struct{}

func (f CompareVersionsFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "compare_versions"
}

func (f CompareVersionsFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Compare two version strings",
		Description: "Given two APIC version strings and a comparison operator, returns a boolean result of the comparison.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "version1",
				Description: "First version string",
			},
			function.StringParameter{
				Name:        "operator",
				Description: "Comparison operator",
				Validators: []function.StringParameterValidator{
					stringvalidator.OneOf("==", "!=", ">", "<", ">=", "<=", "inside", "outside"),
				},
			},
			function.StringParameter{
				Name:        "version2",
				Description: "Second version string",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f CompareVersionsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var version1, version2, operator string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &version1, &operator, &version2))
	if resp.Error != nil {
		return
	}

	result, err := CompareVersionsRange(version1, version2, operator)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

func CompareVersionsRange(version, versionRanges, operator string) (bool, error) {
	return apicversion.CompareVersionsRange(version, versionRanges, operator)
}
