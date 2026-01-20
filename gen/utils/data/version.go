package data

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
)

type VersionRange struct {
	// raw contains the original version range string from the meta file.
	raw string
	// The minimum version of the range.
	// The version is in the format "4.2(7f)".
	// When nil, there is no lower bound.
	min *provider.Version
	// The maximum version of the range.
	// The version is in the format "4.2(7w)".
	// When nil, there is no upper bound (unlimited).
	max *provider.Version
}

func (vr *VersionRange) Raw() string {
	// Raw returns the original version range string from the meta file.
	return vr.raw
}

func (vr *VersionRange) Min() *provider.Version {
	// Min returns the minimum version of the range, or nil if there is no lower bound.
	return vr.min
}

func (vr *VersionRange) Max() *provider.Version {
	// Max returns the maximum version of the range, or nil if there is no upper bound.
	return vr.max
}

func (vr *VersionRange) IsSingleVersion() bool {
	// IsSingleVersion returns true if min and max are equal (representing a single version, not a range).
	return vr.min != nil && vr.max != nil && *vr.min == *vr.max
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
	case vr.min == nil:
		return fmt.Sprintf("up to %s", vr.max)
	case vr.max == nil:
		return fmt.Sprintf("%s and later", vr.min)
	default:
		return fmt.Sprintf("%s to %s", vr.min, vr.max)
	}
}

type Versions struct {
	// raw contains the original version string from the meta file.
	raw string
	// ranges contains the list of version ranges where the class/property is supported.
	ranges []VersionRange
}

func (v *Versions) Raw() string {
	// Raw returns the original version string from the meta file.
	return v.raw
}

func (v *Versions) Ranges() []VersionRange {
	// Ranges returns the list of version ranges.
	return v.ranges
}

func (v *Versions) String() string {
	// String returns the version ranges as a formatted string for documentation.
	var versionRanges []string
	for _, versionRange := range v.ranges {
		versionRanges = append(versionRanges, versionRange.String())
	}
	return strings.Join(versionRanges, ", ")
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
	return sortVersions(a.max, b.max, false)
}

func sortVersions(a, b *provider.Version, isMinBound bool) int {
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
	case provider.IsVersionEqual(*a, *b):
		return 0
	case provider.IsVersionLesser(*a, *b):
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
	genLogger.Trace(fmt.Sprintf("Parsing version string: %s", metaVersions))

	ranges, err := parseVersionRanges(metaVersions)
	if err != nil {
		return nil, err
	}

	versions := &Versions{raw: metaVersions, ranges: ranges}
	versions.Sort()

	genLogger.Trace(fmt.Sprintf("Successfully parsed version string '%s' into %d ranges.", metaVersions, len(versions.ranges)))
	genLogger.Trace(fmt.Sprintf("Constructed version string: %s", versions))
	return versions, nil
}

func parseVersionRanges(metaVersions string) ([]VersionRange, error) {
	// parseVersionRanges parses a comma-separated list of version ranges.
	// Example: "4.2(7f)-4.2(7w),5.2(1g)-" -> [VersionRange{4.2(7f), 4.2(7w)}, VersionRange{5.2(1g), nil}]
	var versionRanges []VersionRange

	for metaVersionRange := range strings.SplitSeq(metaVersions, ",") {
		versionRange, err := parseVersionRange(metaVersionRange)
		if err != nil {
			return nil, err
		}
		versionRanges = append(versionRanges, *versionRange)
	}

	return versionRanges, nil
}

func parseVersionRange(metaVersionRange string) (*VersionRange, error) {
	// parseVersionRange parses a single version range string.
	// Examples:
	//   - "1.0(1e)" -> single version (min == max)
	//   - "4.2(7f)-4.2(7w)" -> bounded range
	//   - "4.2(7f)-" -> unbounded range (max is nil)
	//   - "-4.2(7w)" -> unbounded range (min is nil)
	var minVersion, maxVersion *provider.Version

	// Check if it's a range
	if strings.Contains(metaVersionRange, "-") {
		// SplitN with n=2 ensures we only split on the first dash.
		// Extra dashes remain in the second part and will fail ParseVersion validation.
		metaVersionParts := strings.SplitN(metaVersionRange, "-", 2)
		minMetaVersion := metaVersionParts[0]
		maxMetaVersion := metaVersionParts[1]

		// Parse minimum version if present
		if minMetaVersion != "" {
			minVersionResult := provider.ParseVersion(minMetaVersion)
			if minVersionResult.Error != "" {
				return nil, fmt.Errorf("invalid minimum version '%s': %s", minMetaVersion, minVersionResult.Error)
			}
			minVersion = minVersionResult.Version
		}

		// Parse maximum version if present
		if maxMetaVersion != "" {
			maxVersionResult := provider.ParseVersion(maxMetaVersion)
			if maxVersionResult.Error != "" {
				return nil, fmt.Errorf("invalid maximum version '%s': %s", maxMetaVersion, maxVersionResult.Error)
			}
			maxVersion = maxVersionResult.Version
		}
	} else {
		versionResult := provider.ParseVersion(metaVersionRange)
		if versionResult.Error != "" {
			return nil, fmt.Errorf("invalid version '%s': %s", metaVersionRange, versionResult.Error)
		}
		minVersion = versionResult.Version
		maxVersion = versionResult.Version
	}

	return &VersionRange{raw: metaVersionRange, min: minVersion, max: maxVersion}, nil
}
