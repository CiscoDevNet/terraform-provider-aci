package models

// classVersions stores the raw APIC support range registered by every
// generated class model. Generated resource and data-source adapters use it to
// exclude child classes that are unavailable on the connected APIC version.
var classVersions = map[string]string{}

func registerClassVersion(className, versions string) {
	classVersions[className] = versions
}

// ClassVersions returns the raw APIC version range registered for a class.
func ClassVersions(className string) (string, bool) {
	versions, ok := classVersions[className]
	return versions, ok
}
