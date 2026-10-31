package apicversion

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CompareVersionsRange compares an APIC version against either a single
// version or a comma-separated set of version ranges.
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
				result, err = CompareVersions(version, comparison.Min, "==")
				if err != nil {
					return false, err
				}
				if result {
					break
				}
			} else {
				resultMin, err := CompareVersions(version, comparison.Min, ">=")
				if err != nil {
					return false, err
				}

				resultMax := true
				if comparison.Max != "unlimited" {
					resultMax, err = CompareVersions(version, comparison.Max, "<=")
					if err != nil {
						return false, err
					}
				}

				if resultMin && resultMax {
					result = true
					break
				}
			}
		}

		if operator == "outside" {
			result = !result
		}
	} else {
		result, err = CompareVersions(version, versionRanges, operator)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}

// VersionRange stores the raw lower and upper bounds used by APIC metadata.
type VersionRange struct {
	Min string
	Max string
}

// GetVersionRangesListFromVersionRangesString parses APIC's comma-separated
// version range format while preserving its unlimited upper-bound marker.
func GetVersionRangesListFromVersionRangesString(versionRanges string) ([]VersionRange, error) {
	var comparisons []VersionRange

	for _, versionRange := range strings.Split(versionRanges, ",") {
		versions := strings.Split(versionRange, "-")
		if len(versions) == 2 {
			if versions[1] == "" {
				versions[1] = "unlimited"
			}
			comparisons = append(comparisons, VersionRange{Min: versions[0], Max: versions[1]})
		} else if len(versions) == 1 {
			comparisons = append(comparisons, VersionRange{Min: versions[0], Max: versions[0]})
		} else {
			return nil, fmt.Errorf("invalid range format: %s", versionRange)
		}
	}

	return comparisons, nil
}

// Version is the comparable representation of an APIC release string.
type Version struct {
	Major int
	Minor int
	Patch int
	Tag   int
}

func (v Version) String() string {
	if v.Tag != 0 {
		return fmt.Sprintf("%d.%d(%d%c)", v.Major, v.Minor, v.Patch, v.Tag)
	}
	return fmt.Sprintf("%d.%d(%d)", v.Major, v.Minor, v.Patch)
}

// VersionResult preserves the existing parser contract used by generator and
// provider compatibility callers.
type VersionResult struct {
	Version *Version
	Error   string
}

// ParseVersion parses the numeric and optional letter components of an APIC
// release string.
func ParseVersion(rawVersion string) VersionResult {
	versionRegex := regexp.MustCompile(`(\d+)\.(\d+)\((\d+)([a-z]?)\)`)
	matches := versionRegex.FindStringSubmatch(rawVersion)
	if matches == nil {
		return VersionResult{Error: "unknown"}
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	tag := 0
	if matches[4] != "" {
		tag = int(matches[4][0])
	}

	return VersionResult{Version: &Version{Major: major, Minor: minor, Patch: patch, Tag: tag}}
}

// IsVersionGreater reports whether v1 is later than v2.
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

// IsVersionEqual reports whether both APIC version components are identical.
func IsVersionEqual(v1, v2 Version) bool {
	return v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch == v2.Patch && v1.Tag == v2.Tag
}

// IsVersionNotEqual reports whether any APIC version component differs.
func IsVersionNotEqual(v1, v2 Version) bool {
	return !IsVersionEqual(v1, v2)
}

// IsVersionLesser reports whether v1 is earlier than v2.
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

// IsVersionGreaterOrEqual reports whether v1 is equal to or later than v2.
func IsVersionGreaterOrEqual(v1, v2 Version) bool {
	return IsVersionGreater(v1, v2) || IsVersionEqual(v1, v2)
}

// IsVersionLessOrEqual reports whether v1 is equal to or earlier than v2.
func IsVersionLessOrEqual(v1, v2 Version) bool {
	return IsVersionLesser(v1, v2) || IsVersionEqual(v1, v2)
}

// CompareVersions compares two singular APIC versions with a standard
// comparison operator.
func CompareVersions(version1, version2, operator string) (bool, error) {
	versionResult1 := ParseVersion(version1)
	if versionResult1.Error == "unknown" {
		// Preserve the existing provider-function diagnostic text.
		return false, fmt.Errorf("Invalid APIC version provided: %s", version1) //nolint:staticcheck
	}

	versionResult2 := ParseVersion(version2)
	if versionResult2.Error == "unknown" {
		// Preserve the existing provider-function diagnostic text.
		return false, fmt.Errorf("Invalid APIC version provided: %s", version2) //nolint:staticcheck
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
		// Preserve the existing provider-function diagnostic text.
		return false, fmt.Errorf("Invalid operator: %s", operator) //nolint:staticcheck
	}

	return result, nil
}
