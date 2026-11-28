package data

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/apicversion"
)

type VersionRange struct {
	// raw contains the original version range string from the meta file.
	raw string
	// The minimum version of the range.
	// The version is in the format "4.2(7f)".
	// When nil, there is no lower bound.
	min *apicversion.Version
	// The maximum version of the range.
	// The version is in the format "4.2(7w)".
	// When nil, there is no upper bound (unlimited).
	max *apicversion.Version
	// Derived availability can exclude a boundary without inventing a previous release.
	minExclusive bool
	maxExclusive bool
}

func (vr *VersionRange) Raw() string {
	// Raw returns the original version range string from the meta file.
	return vr.raw
}

func (vr *VersionRange) Min() *apicversion.Version {
	// Min returns the minimum version of the range, or nil if there is no lower bound.
	return vr.min
}

func (vr *VersionRange) Max() *apicversion.Version {
	// Max returns the maximum version of the range, or nil if there is no upper bound.
	return vr.max
}

func (vr *VersionRange) IsSingleVersion() bool {
	// IsSingleVersion returns true if min and max are equal (representing a single version, not a range).
	return vr.min != nil && vr.max != nil && *vr.min == *vr.max && !vr.minExclusive && !vr.maxExclusive
}

func (vr *VersionRange) String() string {
	// String returns the version range as a documentation-friendly string.
	// Examples:
	//   - Single version: "4.2(7f)"
	//   - Bounded range: "4.2(7f) to 4.2(7w)"
	//   - Unbounded upper: "4.2(7f) and later"
	//   - Unbounded lower: "up to 4.2(7w)"
	switch {
	case vr.IsSingleVersion():
		return vr.min.String()
	case vr.maxExclusive && vr.min == nil:
		return fmt.Sprintf("before %s", vr.max)
	case vr.maxExclusive && vr.minExclusive:
		return fmt.Sprintf("after %s and before %s", vr.min, vr.max)
	case vr.maxExclusive:
		return fmt.Sprintf("%s up to, but excluding, %s", vr.min, vr.max)
	case vr.minExclusive && vr.max == nil:
		return fmt.Sprintf("after %s", vr.min)
	case vr.minExclusive:
		return fmt.Sprintf("after %s through %s", vr.min, vr.max)
	case vr.min == nil:
		return fmt.Sprintf("up to %s", vr.max)
	case vr.max == nil:
		return fmt.Sprintf("%s and later", vr.min)
	default:
		return fmt.Sprintf("%s to %s", vr.min, vr.max)
	}
}

func NewVersionRange(metaVersionRange string) (*VersionRange, error) {
	// NewVersionRange parses a single version range string and returns a VersionRange struct.
	// Examples:
	//   - "1.0(1e)" -> single version (min == max)
	//   - "4.2(7f)-4.2(7w)" -> bounded range
	//   - "4.2(7f)-" -> unbounded range (max is nil)
	//   - "-4.2(7w)" -> unbounded range (min is nil)
	genLogger.Tracef("Parsing version range: %s", metaVersionRange)

	minVersion, maxVersion, err := parseVersionBounds(metaVersionRange)
	if err != nil {
		return nil, err
	}

	versionRange := &VersionRange{raw: metaVersionRange, min: minVersion, max: maxVersion}

	genLogger.Tracef("Successfully parsed version range '%s': %s", metaVersionRange, versionRange)
	return versionRange, nil
}

func parseVersionBounds(metaVersionRange string) (*apicversion.Version, *apicversion.Version, error) {
	// parseVersionBounds extracts min and max versions from a version range string.
	// Returns (minVersion, maxVersion, error).
	var minVersion, maxVersion *apicversion.Version

	if strings.Contains(metaVersionRange, "-") {
		// SplitN with n=2 ensures we only split on the first dash.
		// Extra dashes remain in the second part and will fail ParseVersion validation.
		metaVersionParts := strings.SplitN(metaVersionRange, "-", 2)
		minMetaVersion := metaVersionParts[0]
		maxMetaVersion := metaVersionParts[1]

		// Parse minimum version if present
		if minMetaVersion != "" {
			minVersionResult := apicversion.ParseVersion(minMetaVersion)
			if minVersionResult.Error != "" {
				return nil, nil, fmt.Errorf("invalid minimum version '%s': %s", minMetaVersion, minVersionResult.Error)
			}
			minVersion = minVersionResult.Version
		}

		// Parse maximum version if present
		if maxMetaVersion != "" {
			maxVersionResult := apicversion.ParseVersion(maxMetaVersion)
			if maxVersionResult.Error != "" {
				return nil, nil, fmt.Errorf("invalid maximum version '%s': %s", maxMetaVersion, maxVersionResult.Error)
			}
			maxVersion = maxVersionResult.Version
		}
	} else {
		versionResult := apicversion.ParseVersion(metaVersionRange)
		if versionResult.Error != "" {
			return nil, nil, fmt.Errorf("invalid version '%s': %s", metaVersionRange, versionResult.Error)
		}
		minVersion = versionResult.Version
		maxVersion = versionResult.Version
	}

	return minVersion, maxVersion, nil
}

type Versions struct {
	// raw contains the original version string from the meta file.
	raw string
	// ranges contains the list of version ranges where the class/property is supported.
	ranges []VersionRange
}

func (v *Versions) Raw() string {
	// Preserve source strings; derived inclusive ranges use the same metadata syntax.
	if v.raw != "" || v.NeedsBounds() {
		return v.raw
	}
	var ranges []string
	for _, interval := range v.ranges {
		value := interval.min.String()
		if !interval.IsSingleVersion() {
			value += "-"
			if interval.max != nil {
				value += interval.max.String()
			}
		}
		ranges = append(ranges, value)
	}
	return strings.Join(ranges, ",")
}

func (v *Versions) Ranges() []VersionRange {
	// Ranges returns the list of version ranges.
	return v.ranges
}

func (v *Versions) String() string {
	// String returns the version ranges as a formatted string for documentation.
	if len(v.ranges) == 0 {
		return "none"
	}
	var versionRanges []string
	for _, versionRange := range v.ranges {
		versionRanges = append(versionRanges, versionRange.String())
	}
	return strings.Join(versionRanges, ", ")
}

func documentationVersions(v *Versions) string {
	if v.raw == "" {
		return v.String()
	}
	formattedRanges := make([]string, 0, len(v.ranges))
	for _, rawRange := range strings.Split(v.Raw(), ",") {
		versionRange, err := NewVersionRange(strings.TrimSpace(rawRange))
		if err != nil {
			return v.String()
		}
		formattedRanges = append(formattedRanges, versionRange.String())
	}
	formatted := strings.Join(formattedRanges, ", ")
	if strings.HasSuffix(v.Raw(), "-") {
		return formatted + "."
	}
	return formatted
}

// Covers reports whether every range in other is fully supported by v.
func (v *Versions) Covers(other *Versions) bool {
	if v == nil || other == nil {
		return false
	}

	for _, otherRange := range other.ranges {
		if !slices.ContainsFunc(v.ranges, func(versionRange VersionRange) bool {
			return versionRangeCovers(versionRange, otherRange)
		}) {
			return false
		}
	}

	return true
}

func versionRangeCovers(versionRange, otherRange VersionRange) bool {
	if otherRange.min == nil {
		if versionRange.min != nil {
			return false
		}
	} else if versionRange.min != nil {
		comparison := sortVersions(versionRange.min, otherRange.min, true)
		if comparison > 0 || comparison == 0 && versionRange.minExclusive && !otherRange.minExclusive {
			return false
		}
	}

	if otherRange.max == nil {
		return versionRange.max == nil
	}
	if versionRange.max == nil {
		return true
	}
	comparison := sortVersions(versionRange.max, otherRange.max, false)
	return comparison > 0 || comparison == 0 && (!versionRange.maxExclusive || otherRange.maxExclusive)
}

func (v *Versions) Equal(other *Versions) bool {
	// Compare normalized bounds; identical source strings may have different derived endpoints.
	if v == nil || other == nil {
		return v == other
	}
	return slices.EqualFunc(v.ranges, other.ranges, func(a, b VersionRange) bool {
		return sortVersions(a.min, b.min, true) == 0 && sortVersions(a.max, b.max, false) == 0 && a.minExclusive == b.minExclusive && a.maxExclusive == b.maxExclusive
	})
}

func intersectVersionRanges(a, b VersionRange) (VersionRange, bool) {
	// Equal bounds inherit the stricter endpoint; an excluded point is an empty interval.
	result := a
	comparison := sortVersions(a.min, b.min, true)
	if comparison < 0 {
		result.min, result.minExclusive = b.min, b.minExclusive
	} else if comparison == 0 {
		result.minExclusive = a.minExclusive || b.minExclusive
	}
	comparison = sortVersions(a.max, b.max, false)
	if comparison > 0 {
		result.max, result.maxExclusive = b.max, b.maxExclusive
	} else if comparison == 0 {
		result.maxExclusive = a.maxExclusive || b.maxExclusive
	}
	if result.min != nil && result.max != nil {
		comparison = sortVersions(result.min, result.max, true)
		if comparison > 0 || comparison == 0 && (result.minExclusive || result.maxExclusive) {
			return VersionRange{}, false
		}
	}
	result.raw = ""
	return result, true
}

func (v *Versions) Intersect(other *Versions) *Versions {
	// Nil means no additional restriction; an empty non-nil range means no supported version.
	if v == nil {
		return other
	}
	if other == nil || other.Covers(v) {
		return v
	}
	if v.Covers(other) {
		return other
	}
	result := &Versions{}
	for _, a := range v.ranges {
		for _, b := range other.ranges {
			if intersection, ok := intersectVersionRanges(a, b); ok {
				result.ranges = append(result.ranges, intersection)
			}
		}
	}
	result.Sort()
	return result
}

func (v *Versions) Exclude(hidden *Versions) *Versions {
	// Split supported intervals around each hidden range, retaining gaps and boundary semantics.
	if hidden == nil {
		return v
	}
	result := &Versions{ranges: []VersionRange{{}}}
	if v != nil {
		result.ranges = slices.Clone(v.ranges)
	}
	for _, excluded := range hidden.ranges {
		var remaining []VersionRange
		for _, supported := range result.ranges {
			intersection, ok := intersectVersionRanges(supported, excluded)
			if !ok {
				remaining = append(remaining, supported)
				continue
			}
			if intersection.min != nil {
				before := VersionRange{max: intersection.min, maxExclusive: !intersection.minExclusive}
				if interval, ok := intersectVersionRanges(supported, before); ok {
					remaining = append(remaining, interval)
				}
			}
			if intersection.max != nil {
				after := VersionRange{min: intersection.max, minExclusive: !intersection.maxExclusive}
				if interval, ok := intersectVersionRanges(supported, after); ok {
					remaining = append(remaining, interval)
				}
			}
		}
		result.ranges = remaining
	}
	result.Sort()
	if result.Equal(v) {
		return v
	}
	return result
}

func (v *Versions) Go() string {
	// Emit normalized bounds rather than Raw(), which cannot carry exclusive cutoffs.
	var ranges []string
	for _, interval := range v.ranges {
		var fields []string
		if interval.min != nil {
			fields = append(fields, fmt.Sprintf("Min: %q", interval.min.String()))
		}
		if interval.max != nil {
			fields = append(fields, fmt.Sprintf("Max: %q", interval.max.String()))
		}
		if interval.minExclusive {
			fields = append(fields, "MinExclusive: true")
		}
		if interval.maxExclusive {
			fields = append(fields, "MaxExclusive: true")
		}
		ranges = append(ranges, "{"+strings.Join(fields, ", ")+"}")
	}
	return "[]apicversion.VersionRange{" + strings.Join(ranges, ", ") + "}"
}

func (v *Versions) NeedsBounds() bool {
	// Inclusive intervals retain string checks; excluded endpoints and empty availability need exact bounds.
	return len(v.ranges) == 0 || slices.ContainsFunc(v.ranges, func(interval VersionRange) bool {
		return interval.min == nil || interval.minExclusive || interval.maxExclusive
	})
}

func (v *Versions) TestExpression() string {
	// Keep ordinary range expressions stable and render derived bounds with their exact operators.
	if !v.NeedsBounds() {
		return fmt.Sprintf("provider::aci::compare_versions(var.apic_version, %q, %q)", "inside", v.Raw())
	}
	var intervals []string
	for _, interval := range v.ranges {
		var bounds []string
		if interval.min != nil {
			operator := ">="
			if interval.minExclusive {
				operator = ">"
			}
			bounds = append(bounds, fmt.Sprintf("provider::aci::compare_versions(var.apic_version, %q, %q)", operator, interval.min.String()))
		}
		if interval.max != nil {
			operator := "<="
			if interval.maxExclusive {
				operator = "<"
			}
			bounds = append(bounds, fmt.Sprintf("provider::aci::compare_versions(var.apic_version, %q, %q)", operator, interval.max.String()))
		}
		if len(bounds) == 0 {
			return "true"
		}
		intervals = append(intervals, "("+strings.Join(bounds, " && ")+")")
	}
	if len(intervals) == 0 {
		return "false"
	}
	return "(" + strings.Join(intervals, " || ") + ")"
}

func (v *Versions) Sort() {
	// Sort orders the version ranges from lowest to highest.
	// Ranges are compared first by their minimum version, then by their maximum version.
	// nil values are treated as unbounded: nil min comes first, nil max comes last.
	slices.SortFunc(v.ranges, sortVersionRanges)
}

func sortVersionRanges(a, b VersionRange) int {
	// sortVersionRanges compares two VersionRange structs for sorting.
	// Returns -1 if a < b, 0 if a == b, 1 if a > b.

	// Compare by minimum version first
	minCmp := sortVersions(a.min, b.min, true)
	if minCmp != 0 {
		return minCmp
	}

	// If min versions are equal, compare by max version
	maxCmp := sortVersions(a.max, b.max, false)
	if maxCmp == 0 && a.maxExclusive != b.maxExclusive {
		if a.maxExclusive {
			return -1
		}
		return 1
	}
	return maxCmp
}

func sortVersions(a, b *apicversion.Version, isMinBound bool) int {
	// sortVersions compares two version pointers for sorting.
	// isMinBound determines how nil values are ordered:
	//   - true (min bound): nil sorts first, representing "no lower bound" (earliest)
	//   - false (max bound): nil sorts last, representing "no upper bound" (latest)
	aIsNil := a == nil
	bIsNil := b == nil

	switch {
	case aIsNil && bIsNil:
		// Unlikely in practice, but needed for correctness if both bounds are unbounded.
		return 0
	case aIsNil && isMinBound:
		return -1
	case aIsNil:
		return 1
	case bIsNil && isMinBound:
		return 1
	case bIsNil:
		return -1
	case apicversion.IsVersionEqual(*a, *b):
		return 0
	case apicversion.IsVersionLesser(*a, *b):
		return -1
	default:
		return 1
	}
}

func NewVersions(metaVersions string) (*Versions, error) {
	// NewVersions parses a raw version string from meta files and returns a Versions struct.
	// The version string can contain multiple comma-separated ranges.
	// Each range can be:
	//   - A single version: "1.0(1e)" (min and max are the same)
	//   - A bounded range: "4.2(7f)-4.2(7w)" (min to max)
	//   - An unbounded range: "4.2(7f)-" (min to unlimited, max is nil)
	//
	// Example: "3.2(10e)-3.2(10g),3.2(7f)-" produces two ranges.
	genLogger.Tracef("Parsing version string: %s", metaVersions)

	ranges, err := parseVersionRanges(metaVersions)
	if err != nil {
		return nil, err
	}

	versions := &Versions{raw: metaVersions, ranges: ranges}
	versions.Sort()

	genLogger.Tracef("Successfully parsed version string '%s' into %d ranges.", metaVersions, len(versions.ranges))
	genLogger.Tracef("Constructed version string: %s", versions)
	return versions, nil
}

func parseVersionRanges(metaVersions string) ([]VersionRange, error) {
	// parseVersionRanges parses a comma-separated list of version ranges.
	// Example: "4.2(7f)-4.2(7w),5.2(1g)-" -> [VersionRange{4.2(7f), 4.2(7w)}, VersionRange{5.2(1g), nil}]
	var versionRanges []VersionRange

	for metaVersionRange := range strings.SplitSeq(metaVersions, ",") {
		versionRange, err := NewVersionRange(strings.TrimSpace(metaVersionRange))
		if err != nil {
			return nil, err
		}
		versionRanges = append(versionRanges, *versionRange)
	}

	return versionRanges, nil
}
