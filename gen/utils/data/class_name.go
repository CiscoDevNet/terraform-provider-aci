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

func NewClassName(className string) (*ClassName, error) {
	// Creates a ClassName struct with all representations of the class name.
	// Handles both formats:
	//   - Full class name: "fvTenant"
	//   - Meta-style with colon: "fv:Tenant"
	// Returns an error if the class name cannot be split into package and short parts.
	genLogger.Trace(fmt.Sprintf("Creating ClassName from '%s'.", className))

	className, err := sanitizeClassName(className)
	if err != nil {
		return nil, err
	}

	packageName, short, err := splitClassNameToPackageNameAndShortName(className)
	if err != nil {
		return nil, err
	}

	genLogger.Debug(fmt.Sprintf("ClassName created: full='%s', packageName='%s', short='%s'.", className, packageName, short))

	return &ClassName{
		full:        className,
		packageName: packageName,
		short:       short,
		capitalized: cases.Title(language.Und, cases.NoLower).String(className),
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

func sanitizeClassName(className string) (string, error) {
	// Removes the ':' separator from class names if present.
	// Handles both formats:
	//   - Meta-style with colon: "fv:Tenant" -> "fvTenant"
	//   - Already sanitized: "fvTenant" -> "fvTenant" (no change)
	// Returns an error if multiple colons are detected.
	if strings.Contains(className, ":") {
		if strings.Count(className, ":") > 1 {
			return "", fmt.Errorf("invalid class name '%s': multiple colons detected", className)
		}
		className = strings.Replace(className, ":", "", 1)

	}
	return className, nil
}

func splitClassNameToPackageNameAndShortName(className string) (string, string, error) {
	// Splits the class name into the package and short name.
	// Handles both formats:
	//   - Full class name: "fvTenant" -> ("fv", "Tenant")
	//   - Meta-style with colon: "fv:Tenant" -> ("fv", "Tenant")
	// The package and short names are used for the meta file download, documentation links and lookup in the raw data.
	// Returns an error if the class name cannot be split (no uppercase letter found after package prefix).
	genLogger.Trace(fmt.Sprintf("Splitting class name '%s' for name space separation.", className))

	className, err := sanitizeClassName(className)
	if err != nil {
		return "", "", err
	}

	index := strings.IndexFunc(className, unicode.IsUpper)
	if index <= 0 {
		genLogger.Error(fmt.Sprintf("Failed to split class name '%s' for name space separation.", className))
		return "", "", fmt.Errorf("failed to split class name '%s' for name space separation", className)
	}

	packageName := className[:index]
	shortName := className[index:]

	genLogger.Debug(fmt.Sprintf("Class name '%s' got split into package name '%s' and short name '%s'.", className, packageName, shortName))

	return packageName, shortName, nil
}
