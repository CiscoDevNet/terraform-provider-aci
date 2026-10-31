package models

import "github.com/CiscoDevNet/terraform-provider-aci/v2/internal/apicversion"

// classVersions stores the raw APIC support range registered by every
// generated class model. Generated resource and data-source adapters use it to
// exclude subtree classes that are unavailable on the connected APIC version.
var classVersions = map[string]string{}

func registerClassVersion(className, versions string) {
	classVersions[className] = versions
}

// ClassVersions returns the raw APIC version range registered for a class.
func ClassVersions(className string) (string, bool) {
	versions, ok := classVersions[className]
	return versions, ok
}

func supportedSubtreeClasses(apicVersion string, subtreeClasses []string) []string {
	supportedClasses := make([]string, 0, len(subtreeClasses))
	seen := make(map[string]bool, len(subtreeClasses))
	for _, subtreeClass := range subtreeClasses {
		if seen[subtreeClass] {
			continue
		}

		versions, found := classVersions[subtreeClass]
		if !found {
			continue
		}

		supported, err := apicversion.CompareVersionsRange(apicVersion, versions, "inside")
		// Omit classes with unusable version metadata, matching the legacy request
		// behavior: requesting an unsupported subtree class can make APIC reject the read.
		if err != nil || !supported {
			continue
		}

		seen[subtreeClass] = true
		supportedClasses = append(supportedClasses, subtreeClass)
	}
	return supportedClasses
}
