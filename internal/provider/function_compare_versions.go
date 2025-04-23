package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = CompareVersionsFunction{}

func NewCompareVersionsFunction() function.Function {
	return &CompareVersionsFunction{}
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
				Validators:  []function.StringParameterValidator{oneOf("==", "!=", ">", "<", ">=", "<=", "inside", "outside")},
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
	var result bool
	var err error

	if operator == "inside" || operator == "outside" {

		comparisons, err := GetVersionRangesListFromVersionRangesString(versionRanges)
		if err != nil {
			return false, err
		}

		for _, comparison := range comparisons {

			if comparison.Min == comparison.Max {

				// When the min and max are the same, the version is singular which requires single comparison
				result, err = CompareVersions(version, comparison.Min, "==")
				if err != nil {
					return false, err
				}

				// If the version is inside the range, the comparison must be true for result to be true
				// If the version is outside the range, the comparison must be false for result to be true
				if result {
					break
				}

			} else {

				// When the min and max are different, the version is a range
				// This requires two comparisons, one for each side of the range
				// Higher or equal than min and lower or equal than max
				resultMin, err := CompareVersions(version, comparison.Min, ">=")
				if err != nil {
					return false, err
				}

				resultMax := true
				// If max is "unlimited", we don't need to check the upper bound
				if comparison.Max != "unlimited" {
					resultMax, err = CompareVersions(version, comparison.Max, "<=")
					if err != nil {
						return false, err
					}
				}

				// If the version is inside the range, both comparison must be true
				// If the version is outside the range, one comparison must be false
				if resultMin && resultMax {
					result = true
					break
				}
			}
		}

		// If the operator is "outside", invert the result of the comparison that checks inside the range
		if operator == "outside" {
			result = !result
		}

	} else { // If the operator is "==", "!=", ">", "<", ">=", "<=", we assume no range is provided and use the CompareVersions function directly

		result, err = CompareVersions(version, versionRanges, operator)
		if err != nil {
			return false, err
		}

	}

	return result, nil
}

type VersionRange struct {
	Min string
	Max string
}

func GetVersionRangesListFromVersionRangesString(versionRanges string) ([]VersionRange, error) {

	var comparisons []VersionRange

	// Split the ranges string  into list of ranges "4.2(7f)-4.2(7w),5.2(1g)-5.2(1t)" -> ["4.2(7f)-4.2(7w)", "5.2(1g)-5.2(1t)"]
	for _, versionRange := range strings.Split(versionRanges, ",") {
		// Split the range into min and max versions "4.2(7f)-4.2(7w)" -> ["4.2(7f)", "4.2(7w)"]
		versions := strings.Split(versionRange, "-")
		if len(versions) == 2 {
			if versions[1] == "" {
				versions[1] = "unlimited"
			}
			comparisons = append(
				comparisons,
				VersionRange{
					Min: versions[0],
					Max: versions[1],
				},
			)
		} else if len(versions) == 1 {
			comparisons = append(comparisons, VersionRange{
				Min: versions[0],
				Max: versions[0],
			})
		} else {
			return nil, fmt.Errorf("Invalid Range Format: %s", versionRange)
		}
	}

	return comparisons, nil
}

type Version struct {
	Major int
	Minor int
	Patch int
	Tag   int
}

type VersionResult struct {
	Version *Version
	Error   string
}

func ParseVersion(rawVersion string) VersionResult {
	versionRegex := regexp.MustCompile(`(\d+)\.(\d+)\((\d+)([a-z])\)`)
	matches := versionRegex.FindStringSubmatch(rawVersion)
	if matches == nil {
		return VersionResult{Error: "unknown"}
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	tag := int(matches[4][0])

	return VersionResult{Version: &Version{Major: major, Minor: minor, Patch: patch, Tag: tag}}

}

func IsVersionGreater(v1, v2 Version) bool {
	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	} else if v1.Minor != v2.Minor {
		return v1.Minor > v2.Minor
	} else if v1.Patch != v2.Patch {
		return v1.Patch > v2.Patch
	}
	return v1.Tag > v2.Tag
}

func IsVersionEqual(v1, v2 Version) bool {
	return v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch == v2.Patch && v1.Tag == v2.Tag
}

func IsVersionNotEqual(v1, v2 Version) bool {
	return !IsVersionEqual(v1, v2)
}

func IsVersionLesser(v1, v2 Version) bool {
	if v1.Major != v2.Major {
		return v1.Major < v2.Major
	} else if v1.Minor != v2.Minor {
		return v1.Minor < v2.Minor
	} else if v1.Patch != v2.Patch {
		return v1.Patch < v2.Patch
	}
	return v1.Tag < v2.Tag
}

func IsVersionGreaterOrEqual(v1, v2 Version) bool {
	return IsVersionGreater(v1, v2) || IsVersionEqual(v1, v2)
}

func IsVersionLessOrEqual(v1, v2 Version) bool {
	return IsVersionLesser(v1, v2) || IsVersionEqual(v1, v2)
}

func CompareVersions(version1, version2, operator string) (bool, error) {
	versionResult1 := ParseVersion(version1)
	if versionResult1.Error == "unknown" {
		return false, errors.New(fmt.Sprintf("Invalid APIC version provided: %s", version1))
	}

	versionResult2 := ParseVersion(version2)
	if versionResult2.Error == "unknown" {
		return false, errors.New(fmt.Sprintf("Invalid APIC version provided: %s", version2))
	}

	v1 := *versionResult1.Version
	v2 := *versionResult2.Version

	var result bool
	switch operator {
	case "==":
		result = IsVersionEqual(v1, v2)
	case "!=":
		result = IsVersionNotEqual(v1, v2)
	case ">":
		result = IsVersionGreater(v1, v2)
	case "<":
		result = IsVersionLesser(v1, v2)
	case ">=":
		result = IsVersionGreaterOrEqual(v1, v2)
	case "<=":
		result = IsVersionLessOrEqual(v1, v2)
	default:
		return false, errors.New(fmt.Sprintf("Invalid operator: %s", operator))
	}

	return result, nil
}

var _ function.StringParameterValidator = oneOfValidator{}

type oneOfValidator struct {
	AllowedValues []string
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("value must be one of %v", v.AllowedValues)
}

func (v oneOfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("value must be one of `%v`", v.AllowedValues)
}

func (v oneOfValidator) ValidateParameterString(ctx context.Context, req function.StringParameterValidatorRequest, resp *function.StringParameterValidatorResponse) {

	if req.Value.IsUnknown() || req.Value.IsNull() {
		return
	}

	value := req.Value.ValueString()
	for _, allowedValue := range v.AllowedValues {
		if value == allowedValue {
			return
		}
	}

	resp.Error = function.NewArgumentFuncError(
		req.ArgumentPosition,
		fmt.Sprintf("Invalid Value: Value must be one of %v, got: %s.", v.AllowedValues, value),
	)
}

func oneOf(values ...string) oneOfValidator {
	return oneOfValidator{
		AllowedValues: values,
	}
}
