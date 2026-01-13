package data

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ClassName struct {
	full        string
	packageName string
	short       string
	capitalized string
	metaStyle   string
}

func NewClassName(name string) (*ClassName, error) {
	// Creates a ClassName from either a full class name (e.g., "fvTenant")
	// or a meta-style class name with colon separator (e.g., "fv:Tenant").
	// Returns an error if the class name cannot be split into package and short parts.
	genLogger.Trace(fmt.Sprintf("Creating ClassName from '%s'.", name))

	// If the name contains a colon, sanitize it first (meta-style format).
	if strings.Contains(name, ":") {
		var err error
		name, err = sanitizeClassName(name)
		if err != nil {
			return nil, err
		}
	}

	packageName, short, err := splitClassNameToPackageNameAndShortName(name)
	if err != nil {
		return nil, err
	}

	genLogger.Debug(fmt.Sprintf("ClassName created: full='%s', packageName='%s', short='%s'.", name, packageName, short))

	return &ClassName{
		full:        name,
		packageName: packageName,
		short:       short,
		capitalized: cases.Title(language.Und, cases.NoLower).String(name),
		metaStyle:   fmt.Sprintf("%s:%s", packageName, short),
	}, nil
}

func (cn *ClassName) Capitalized() string {
	// Returns the class name with the first letter capitalized (e.g., "FvTenant").
	// Used for function names in generated code.
	return cn.capitalized
}

func (cn *ClassName) Package() string {
	// Returns the package prefix (e.g., "fv").
	return cn.packageName
}

func (cn *ClassName) Short() string {
	// Returns the short name without the package prefix (e.g., "Tenant").
	return cn.short
}

func (cn *ClassName) MetaStyle() string {
	// Returns the class name in meta file format with colon separator (e.g., "fv:Tenant").
	return cn.metaStyle
}

func (cn *ClassName) String() string {
	// Returns the full class name (alias for Full).
	return cn.full
}

func sanitizeClassName(classNameWithColon string) (string, error) {
	// Removes the ':' separator from class names retrieved from meta.
	// Example: "fv:Tenant" -> "fvTenant"
	// Returns an error if multiple colons are detected.
	if strings.Count(classNameWithColon, ":") > 1 {
		return "", fmt.Errorf("invalid class name '%s': multiple colons detected", classNameWithColon)
	}
	return strings.Replace(classNameWithColon, ":", "", 1), nil
}

func splitClassNameToPackageNameAndShortName(className string) (string, string, error) {
	// Splits the class name into the package and short name.
	// The package and short names are used for the meta file download, documentation links and lookup in the raw data.
	var shortName, packageName string
	genLogger.Trace(fmt.Sprintf("Splitting class name '%s' for name space separation.", className))
	for index, character := range className {
		if unicode.IsUpper(character) {
			shortName = className[index:]
			packageName = className[:index]
			break
		}
	}

	genLogger.Debug(fmt.Sprintf("Class name '%s' got split into package name '%s' and short name '%s'.", className, packageName, shortName))

	if packageName == "" || shortName == "" {
		genLogger.Error(fmt.Sprintf("Failed to split class name '%s' for name space separation.", className))
		return "", "", fmt.Errorf("failed to split class name '%s' for name space separation", className)
	}

	return packageName, shortName, nil
}
