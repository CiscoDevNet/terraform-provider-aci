package data

import (
	"fmt"
	"math"
	"slices"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
)

type Property struct {
	// The name of the property in the resource and datasource schemas.
	AttributeName string
	// Indicates if the property is computed in the resource schemas.
	Computed bool
	// Indicates if the property is deprecated in the resource and datasource schemas.
	// Deprecated properties include a warning the resource and datasource schemas.
	// Driven by the meta `isDeprecated` flag with an optional definition override (logical OR).
	Deprecated bool
	// The deprecated APIC versions for the property.
	// Used to indicate versions where the property is deprecated but still functional.
	// Driven by the meta `deprecatedSince` value with an optional definition override.
	DeprecatedVersions *Versions
	// Documentation specific information for the property.
	Documentation PropertyDocumentation
	// Hidden indicates that the property is no longer accepted by the APIC API.
	// Driven by the meta `isHidden` flag with an optional definition override (logical OR).
	Hidden bool
	// The hidden APIC versions for the property.
	// Driven by the meta `hiddenSince` value with an optional definition override.
	HiddenVersions *Versions
	// Migration specific information for the property.
	// This is a map that contains the migration value details of the attribute for a specific schema version.
	MigrationValues map[int]MigrationValue
	// Indicates if a property is optional in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Optional bool
	// The name of the property in the APIC class.
	PropertyName string
	// Indicates if the property is read-only in the resource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	ReadOnly bool
	// Indicates if the property is required in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Required bool
	// Indicates if the property is sensitive in the resource and datasource schemas.
	// Exposed as a separate bool because it directly maps to a Terraform schema construct, which makes templating easier.
	Sensitive bool
	// The supported APIC versions for the property.
	// Each version range is separated by a comma, ex "4.2(7f)-4.2(7w),5.2(1g)-".
	// The first version is the minimum version and the second version is the maximum version.
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
	SupportedVersions *Versions
	// Test specific information for the property.
	// This is used to generate the test cases and examples for the property.
	TestValues []TestValue
	// Validation specific information for the property.
	// In the meta file for the class this is a regex statement that is used to validate the property.
	Validators []Validator
	// Specifies the valid values for the property when only certain values are allowed as input.
	// Sourced from the meta `validValues` array, with definition-driven AddValidValues and RemoveValidValues overrides.
	ValidValues ValidValues
	// The ValueTypeEnum type is used to indicate the type of the property in the resource and datasource schemas.
	ValueType ValueTypeEnum
	// The global meta definition containing global overrides. Unexported because it is only used internally by setter methods.
	globalDefinition GlobalMetaDefinition
	// The meta file details for the property. Unexported because it is only used internally by setter methods.
	metaDetails map[string]any
	// The property definition overrides from the class definition file. Unexported because it is only used internally by setter methods.
	propertyDefinition PropertyDefinition
}

type MigrationValue struct {
	// The name of the property in the legacy resource schema.
	AttributeName string
	// Indicates if a property is computed in the legacy resource schema.
	Computed bool
	// Indicates if a property is optional in the legacy resource schema.
	Optional bool
	// Indicates if a property is required in the legacy resource schema.
	Required bool
	// The type of the legacy attribute.
	Type ValueTypeEnum
}

type RegexStatement struct {
	// The regex string.
	Regex string
	// The type of the regex statement.
	Type RegexStatementTypeEnum
}

type RegexStatementTypeEnum int

// The enumeration options of the RegexStatement Type.
const (
	// Include indicates that the value must match the regex statement.
	Include RegexStatementTypeEnum = iota + 1
)

type TestValue struct {
	// The changed value of the property to be used in the test when a property is allowed to be changed without destruction of the resource.
	Changed []string
	// The initial value of the property to be used in the test.
	// This is set to the default value of the property in APIC when it is not a required value.
	Initial []string
}

type Validator struct {
	// If the property has a range of values, these are the minimum and maximum values of the range.
	Max int64
	Min int64
	// If the property has one or more regex statements it requires to match.
	RegexList []RegexStatement
}

// ValidValues is a map keyed by the wire value (e.g. "1") with details about that valid value.
// The map is sourced from the meta `validValues` array, with optional add/remove overrides from the property definition.
type ValidValues map[string]ValidValue

// ValidValue describes a single allowed value for a property.
// Currently it only carries the localName (the human-readable identifier, e.g. "level3") but is
// designed to be extended later with additional fields without changing the shape of the
// surrounding map. For example, when valid values become version-specific, a Versions field
// can be added here so that templates can render only the values supported by a given APIC
// version. Other anticipated extensions include Label, Comment, and PlatformFlavors.
type ValidValue struct {
	// The localName of the valid value (e.g. "level3").
	LocalName string
}

// LocalNamesList returns the localNames of all valid values, sorted alphabetically.
func (vv ValidValues) LocalNamesList() []string {
	localNames := make([]string, 0, len(vv))
	for _, entry := range vv {
		localNames = append(localNames, entry.LocalName)
	}
	slices.Sort(localNames)
	return localNames
}

// ValuesList returns the wire values of all valid values, sorted alphabetically.
func (vv ValidValues) ValuesList() []string {
	values := make([]string, 0, len(vv))
	for value := range vv {
		values = append(values, value)
	}
	slices.Sort(values)
	return values
}

// ValueLocalNameMap returns a flattened map of wire value to localName, useful for templates.
func (vv ValidValues) ValueLocalNameMap() map[string]string {
	out := make(map[string]string, len(vv))
	for value, entry := range vv {
		out[value] = entry.LocalName
	}
	return out
}

// ValueTypeEnum identifies the data shape of a property and is the single dispatch point
// for type-specific schema and template behavior. New named custom types should be added
// here as constants and registered in parseValueType so the same definition `value_type`
// override key can select them.
type ValueTypeEnum int

// The enumeration options of the ValueType.
const (
	// String indicates that the property is a plain string value.
	String ValueTypeEnum = iota + 1
	// Set indicates that the property is a set value (driven by meta uitype "bitmask").
	Set
	// IpAddress indicates that the property is an IP address (IPv4 or IPv6); driven by
	// meta `validateAsIPv4OrIPv6`. Renders with the IP-address custom type for parsing,
	// validation, and semantic-equality (e.g. zero-padding normalization).
	IpAddress
	// SemanticEquality indicates that the property has both ValidValues and Validators,
	// meaning the wire form (e.g. "22") and the human form (e.g. "ssh") must compare equal.
	// Templates render this with the semantic-equality custom type.
	SemanticEquality
)

// parseValueType maps the snake_case form used in property definition `value_type` overrides
// to the typed enum. Returns (0, false) for unknown values. Extend this when adding new
// ValueTypeEnum entries so the override key vocabulary stays in sync with the enum.
func parseValueType(rawType string) (ValueTypeEnum, bool) {
	switch rawType {
	case "string":
		return String, true
	case "set":
		return Set, true
	case "ip_address":
		return IpAddress, true
	case "semantic_equality":
		return SemanticEquality, true
	default:
		return 0, false
	}
}

// knownStringUiTypes are the meta `uitype` values that legitimately render as a string-typed
// schema attribute. Listed explicitly so that genuinely new uitype values surface via WARN.
var knownStringUiTypes = map[string]struct{}{
	"":         {}, // missing key
	"string":   {},
	"enum":     {},
	"number":   {},
	"boolean":  {},
	"auto":     {},
	"password": {}, // sensitive string; sensitivity is handled separately via meta `secure`.
}

func NewProperty(name string, details map[string]any, definition PropertyDefinition, globalDefinition GlobalMetaDefinition) (*Property, error) {
	genLogger.Tracef("Creating new property struct for property: %s.", name)

	property := &Property{
		PropertyName:       name,
		globalDefinition:   globalDefinition,
		metaDetails:        details,
		propertyDefinition: definition,
	}

	err := property.setPropertyData()
	if err != nil {
		return nil, err
	}

	genLogger.Tracef("Successfully created new property struct for property: %s.", name)

	return property, nil
}

func (p *Property) setPropertyData() error {
	genLogger.Debugf("Setting property data for property '%s'.", p.PropertyName)

	p.setAttributeName()

	p.setDeprecated()

	err := p.setDeprecatedVersions()
	if err != nil {
		return err
	}

	p.setHidden()

	err = p.setHiddenVersions()
	if err != nil {
		return err
	}

	// TODO: add function to set MigrationValues
	p.setMigrationValues()

	p.setRequired()

	// setOptional is called after setRequired because it depends on p.Required.
	p.setOptional()

	p.setReadOnly()

	// setComputed is called after setRequired because it depends on p.Required.
	p.setComputed()

	p.setSensitive()

	// TODO: add function to set TestValues
	p.setTestValues()

	err = p.setValidators()
	if err != nil {
		return err
	}

	err = p.setValidValues()
	if err != nil {
		return err
	}

	err = p.setValueType()
	if err != nil {
		return err
	}

	// setDocumentation is called after setValidValues and setValueType because
	// setDefaultValues depends on p.ValidValues and p.ValueType.
	err = p.setDocumentation()
	if err != nil {
		return err
	}

	err = p.setSupportedVersions()
	if err != nil {
		return err
	}

	genLogger.Debugf("Successfully set property data for property '%s'.", p.PropertyName)
	return nil
}

func (p *Property) setAttributeName() {
	// Determine the attribute name of the property.
	// Priority: per-class definition override > global attribute name override > default snake_case derivation.
	genLogger.Debugf("Setting AttributeName for property '%s'.", p.PropertyName)

	if p.propertyDefinition.AttributeName != "" {
		p.AttributeName = p.propertyDefinition.AttributeName
	} else if override, ok := p.globalDefinition.AttributeNameOverrides[p.PropertyName]; ok {
		p.AttributeName = override
	} else {
		p.AttributeName = utils.Underscore(p.PropertyName)
	}

	genLogger.Debugf("Successfully set AttributeName '%s' for property '%s'.", p.AttributeName, p.PropertyName)
}

func (p *Property) setComputed() {
	// Determine if the property is computed.
	// By default all properties are computed except required properties,
	// because optional properties can have server-side defaults and read-only properties are always computed.
	genLogger.Debugf("Setting Computed for property '%s'.", p.PropertyName)

	p.Computed = !p.Required

	genLogger.Debugf("Successfully set Computed '%t' for property '%s'.", p.Computed, p.PropertyName)
}

func (p *Property) setDeprecated() {
	// Determine if the property is deprecated.
	genLogger.Debugf("Setting Deprecated for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Deprecated {
		p.Deprecated = true
	} else if metaDeprecated, ok := p.metaDetails["isDeprecated"].(bool); ok {
		p.Deprecated = metaDeprecated
	}

	genLogger.Debugf("Successfully set Deprecated '%t' for property '%s'.", p.Deprecated, p.PropertyName)
}

func (p *Property) setDeprecatedVersions() error {
	// Determine the deprecated APIC versions for the property from the definition file (override) or meta `deprecatedSince`.
	genLogger.Debugf("Setting DeprecatedVersions for property '%s'.", p.PropertyName)

	deprecatedVersions := p.propertyDefinition.DeprecatedVersions
	if deprecatedVersions == "" {
		deprecatedVersions, _ = p.metaDetails["deprecatedSince"].(string)
	}
	if deprecatedVersions == "" {
		genLogger.Debugf("No DeprecatedVersions specified for property '%s'.", p.PropertyName)
		return nil
	}

	parsedVersions, err := NewVersions(deprecatedVersions)
	if err != nil {
		return fmt.Errorf("failed to parse deprecated versions for property '%s': %w", p.PropertyName, err)
	}
	p.DeprecatedVersions = parsedVersions

	genLogger.Debugf("Successfully set DeprecatedVersions for property '%s'. Versions: '%s'", p.PropertyName, p.DeprecatedVersions)
	return nil
}

func (p *Property) setHidden() {
	// Determine if the property is hidden by the APIC API.
	genLogger.Debugf("Setting Hidden for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Hidden {
		p.Hidden = true
	} else if metaHidden, ok := p.metaDetails["isHidden"].(bool); ok {
		p.Hidden = metaHidden
	}

	genLogger.Debugf("Successfully set Hidden '%t' for property '%s'.", p.Hidden, p.PropertyName)
}

func (p *Property) setHiddenVersions() error {
	// Determine the hidden APIC versions for the property from the definition file (override) or meta `hiddenSince`.
	genLogger.Debugf("Setting HiddenVersions for property '%s'.", p.PropertyName)

	hiddenVersions := p.propertyDefinition.HiddenVersions
	if hiddenVersions == "" {
		hiddenVersions, _ = p.metaDetails["hiddenSince"].(string)
	}
	if hiddenVersions == "" {
		genLogger.Debugf("No HiddenVersions specified for property '%s'.", p.PropertyName)
		return nil
	}

	parsedVersions, err := NewVersions(hiddenVersions)
	if err != nil {
		return fmt.Errorf("failed to parse hidden versions for property '%s': %w", p.PropertyName, err)
	}
	p.HiddenVersions = parsedVersions

	genLogger.Debugf("Successfully set HiddenVersions for property '%s'. Versions: '%s'", p.PropertyName, p.HiddenVersions)
	return nil
}

func (p *Property) setMigrationValues() {
	// Determine the migration values for the property.
	genLogger.Debugf("Setting MigrationValues for property '%s'.", p.PropertyName)
	genLogger.Debugf("Successfully set MigrationValues for property '%s'.", p.PropertyName)
}

func (p *Property) setOptional() {
	// Determine if the property is optional.
	// A property is optional when the definition restriction is "optional",
	// or when the meta file indicates isConfigurable is true and the property is not required and the restriction is not "read_only".
	genLogger.Debugf("Setting Optional for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Restriction == "optional" {
		p.Optional = true
	} else if p.propertyDefinition.Restriction != "read_only" && p.metaDetails["isConfigurable"] == true && !p.Required {
		p.Optional = true
	}

	genLogger.Debugf("Successfully set Optional '%t' for property '%s'.", p.Optional, p.PropertyName)
}

func (p *Property) setReadOnly() {
	// Determine if the property is read-only.
	// A property is read-only only when the definition restriction is "read_only".
	// This is used to include isConfigurable=false properties as read-only attributes in the schema.
	genLogger.Debugf("Setting ReadOnly for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Restriction == "read_only" {
		p.ReadOnly = true
	}

	genLogger.Debugf("Successfully set ReadOnly '%t' for property '%s'.", p.ReadOnly, p.PropertyName)
}

func (p *Property) setRequired() {
	// Determine if the property is required.
	// A property is required when the definition restriction is "required",
	// or when the meta file indicates isConfigurable and isNaming are both true.
	genLogger.Debugf("Setting Required for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Restriction == "required" {
		p.Required = true
	} else if p.metaDetails["isConfigurable"] == true && p.metaDetails["isNaming"] == true {
		p.Required = true
	}

	genLogger.Debugf("Successfully set Required '%t' for property '%s'.", p.Required, p.PropertyName)
}

func (p *Property) setSensitive() {
	// Determine if the property is sensitive.
	// A property is sensitive when the definition override is true,
	// or when the meta file indicates secure is true.
	genLogger.Debugf("Setting Sensitive for property '%s'.", p.PropertyName)

	if p.propertyDefinition.Sensitive {
		p.Sensitive = true
	} else if p.metaDetails["secure"] == true {
		p.Sensitive = true
	}

	genLogger.Debugf("Successfully set Sensitive '%t' for property '%s'.", p.Sensitive, p.PropertyName)
}

func (p *Property) setTestValues() {
	// Determine the test values for the property.
	genLogger.Debugf("Setting TestValues for property '%s'.", p.PropertyName)
	genLogger.Debugf("Successfully set TestValues for property '%s'.", p.PropertyName)
}

func (p *Property) setValidators() error {
	// Determine the validators for the property.
	// Driven by the meta `validators` array; the property definition replaces the meta entirely when non-empty.
	genLogger.Debugf("Setting Validators for property '%s'.", p.PropertyName)

	var validators []Validator
	var err error

	if len(p.propertyDefinition.Validators) > 0 {
		validators, err = validatorsFromDefinition(p.propertyDefinition.Validators)
	} else {
		validators, err = parseValidatorsFromMeta(p.metaDetails["validators"])
	}
	if err != nil {
		return fmt.Errorf("failed to parse validators for property '%s': %w", p.PropertyName, err)
	}

	p.Validators = validators

	genLogger.Debugf("Successfully set Validators for property '%s'. Count: %d", p.PropertyName, len(p.Validators))
	return nil
}

// parseValidatorsFromMeta converts the raw `validators` value from the meta JSON (already decoded into any) into a typed slice.
// Returns nil when the value is absent. Returns an error on shape mismatch or unknown regex statement type.
func parseValidatorsFromMeta(rawValidators any) ([]Validator, error) {
	if rawValidators == nil {
		return nil, nil
	}

	validatorList, ok := rawValidators.([]any)
	if !ok {
		return nil, fmt.Errorf("expected validators to be a list, got %T", rawValidators)
	}
	if len(validatorList) == 0 {
		return nil, nil
	}

	validators := make([]Validator, 0, len(validatorList))
	for index, rawValidator := range validatorList {
		validatorMap, ok := rawValidator.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected validator entry %d to be a map, got %T", index, rawValidator)
		}

		minValue, err := readOptionalInt64(validatorMap, "min")
		if err != nil {
			return nil, fmt.Errorf("validator entry %d: %w", index, err)
		}
		maxValue, err := readOptionalInt64(validatorMap, "max")
		if err != nil {
			return nil, fmt.Errorf("validator entry %d: %w", index, err)
		}

		validator := Validator{Min: minValue, Max: maxValue}

		if rawRegexes, hasRegexes := validatorMap["regexs"]; hasRegexes && rawRegexes != nil {
			regexList, ok := rawRegexes.([]any)
			if !ok {
				return nil, fmt.Errorf("validator entry %d: expected regexs to be a list, got %T", index, rawRegexes)
			}
			regexStatements := make([]RegexStatement, 0, len(regexList))
			for regexIndex, rawRegexEntry := range regexList {
				regexMap, ok := rawRegexEntry.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("validator entry %d regex %d: expected map, got %T", index, regexIndex, rawRegexEntry)
				}
				regexValue, _ := regexMap["regex"].(string)
				typeValue, _ := regexMap["type"].(string)
				parsedType, err := parseRegexStatementType(typeValue)
				if err != nil {
					return nil, fmt.Errorf("validator entry %d regex %d: %w", index, regexIndex, err)
				}
				regexStatements = append(regexStatements, RegexStatement{Regex: regexValue, Type: parsedType})
			}
			validator.RegexList = regexStatements
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

// validatorsFromDefinition converts ValidatorDefinition entries from a property definition into typed Validator structs.
func validatorsFromDefinition(definitionValidators []ValidatorDefinition) ([]Validator, error) {
	validators := make([]Validator, 0, len(definitionValidators))
	for index, definitionValidator := range definitionValidators {
		validator := Validator{Min: definitionValidator.Min, Max: definitionValidator.Max}

		if len(definitionValidator.RegexList) > 0 {
			regexStatements := make([]RegexStatement, 0, len(definitionValidator.RegexList))
			for regexIndex, definitionRegex := range definitionValidator.RegexList {
				parsedType, err := parseRegexStatementType(definitionRegex.Type)
				if err != nil {
					return nil, fmt.Errorf("validator entry %d regex %d: %w", index, regexIndex, err)
				}
				regexStatements = append(regexStatements, RegexStatement{Regex: definitionRegex.Regex, Type: parsedType})
			}
			validator.RegexList = regexStatements
		}

		validators = append(validators, validator)
	}
	return validators, nil
}

// parseRegexStatementType maps the raw string form to the typed enum.
func parseRegexStatementType(rawType string) (RegexStatementTypeEnum, error) {
	switch rawType {
	case "include":
		return Include, nil
	default:
		return 0, fmt.Errorf("unknown regex statement type %q", rawType)
	}
}

// readOptionalInt64 returns the value at key as int64. The zero value is returned when the key is absent or nil.
// Returns an error when the value is not a JSON number or when the number is not an integer within int64 range.
// (encoding/json decodes JSON numbers into float64 when the target is any, so we type-assert float64 first.)
func readOptionalInt64(source map[string]any, key string) (int64, error) {
	rawValue, present := source[key]
	if !present || rawValue == nil {
		return 0, nil
	}
	numericValue, ok := rawValue.(float64)
	if !ok {
		return 0, fmt.Errorf("expected %s to be a number, got %T", key, rawValue)
	}
	if numericValue != math.Trunc(numericValue) || numericValue < math.MinInt64 || numericValue > math.MaxInt64 {
		return 0, fmt.Errorf("expected %s to be an integer, got %v", key, numericValue)
	}
	return int64(numericValue), nil
}

func (p *Property) setValidValues() error {
	// Determine the valid values for the property.
	// Driven by the meta `validValues` array; the property definition can add or remove entries by localName.
	// Entries with localName "defaultValue" are skipped as they only carry the default value information.
	genLogger.Debugf("Setting ValidValues for property '%s'.", p.PropertyName)

	validValues := ValidValues{}

	removeSet := make(map[string]struct{}, len(p.propertyDefinition.RemoveValidValues))
	for _, name := range p.propertyDefinition.RemoveValidValues {
		removeSet[name] = struct{}{}
	}

	// Warn when a name appears in both AddValidValues and RemoveValidValues: the Add wins but the
	// definition is contradictory and likely a mistake.
	for _, localName := range p.propertyDefinition.AddValidValues {
		if _, contradicts := removeSet[localName]; contradicts {
			genLogger.Warnf("AddValidValues %q also listed in RemoveValidValues for property %q; the Add takes precedence.", localName, p.PropertyName)
		}
	}

	if rawValidValues := p.metaDetails["validValues"]; rawValidValues != nil {
		validValueList, ok := rawValidValues.([]any)
		if !ok {
			return fmt.Errorf("failed to parse validValues for property '%s': expected validValues to be a list, got %T", p.PropertyName, rawValidValues)
		}
		for index, rawEntry := range validValueList {
			entry, ok := rawEntry.(map[string]any)
			if !ok {
				return fmt.Errorf("failed to parse validValues for property '%s': expected validValues entry %d to be a map, got %T", p.PropertyName, index, rawEntry)
			}
			localName, localOk := entry["localName"].(string)
			value, valueOk := entry["value"].(string)
			if !localOk || !valueOk {
				return fmt.Errorf("failed to parse validValues for property '%s': validValues entry %d is missing or has non-string localName/value", p.PropertyName, index)
			}
			if localName == "defaultValue" {
				continue
			}
			if _, drop := removeSet[localName]; drop {
				// Delete on hit so any names still in removeSet after the loop are stale.
				delete(removeSet, localName)
				continue
			}
			if existing, exists := validValues[value]; exists {
				// Multiple localNames can legitimately alias the same wire value in the meta
				// (e.g. "aqua" and "cyan" both map to "0x00FFFF"). Keep the first-seen entry
				// under the wire-value key and fall back to inserting the alias under its
				// localName key so it is still addressable (mirroring how AddValidValues
				// entries are keyed). If the localName key is also already taken (e.g. the
				// same localName appears more than once in the meta, or a prior alias already
				// claimed it), skip the alias entirely with a warning. Use RemoveValidValues
				// in the property definition if a different localName should win the
				// wire-value slot.
				if _, localTaken := validValues[localName]; localTaken {
					genLogger.Warnf("Duplicate validValues value %q for property %q: keeping localName %q, skipping alias %q (localName key already in use).", value, p.PropertyName, existing.LocalName, localName)
					continue
				}
				genLogger.Warnf("Duplicate validValues value %q for property %q: keeping localName %q under value key, registering alias %q under its localName key.", value, p.PropertyName, existing.LocalName, localName)
				validValues[localName] = ValidValue{LocalName: localName}
				continue
			}
			validValues[value] = ValidValue{LocalName: localName}
		}
	}

	for name := range removeSet {
		genLogger.Warnf("RemoveValidValues %q not found in meta for property %q.", name, p.PropertyName)
	}

	for _, localName := range p.propertyDefinition.AddValidValues {
		if _, exists := validValues[localName]; exists {
			genLogger.Warnf("AddValidValues %q already present in meta for property %q.", localName, p.PropertyName)
		}
		validValues[localName] = ValidValue{LocalName: localName}
	}

	p.ValidValues = validValues

	genLogger.Debugf("Successfully set ValidValues for property '%s'. Count: %d", p.PropertyName, len(p.ValidValues))
	return nil
}

func (p *Property) setValueType() error {
	// Determine the value type of the property.
	// Precedence (first match wins):
	//   1. Definition `value_type` override (errors on unknown).
	//   2. Meta `uitype == "bitmask"` -> Set.
	//   3. Meta `validateAsIPv4OrIPv6 == true` -> IpAddress.
	//   4. Has both ValidValues and Validators -> SemanticEquality.
	//      (Reads fields populated earlier in setPropertyData by setValidators and
	//      setValidValues; pipeline order matters here.)
	//   5. Meta `uitype` in the known string-rendered set (or missing) -> String silently.
	//   6. Anything else -> String + WARN so new meta vocabulary surfaces.
	genLogger.Debugf("Setting ValueType for property '%s'.", p.PropertyName)

	if override := p.propertyDefinition.ValueType; override != "" {
		parsed, ok := parseValueType(override)
		if !ok {
			return fmt.Errorf("failed to parse value_type for property '%s': unknown value_type %q", p.PropertyName, override)
		}
		p.ValueType = parsed
		genLogger.Debugf("Successfully set ValueType '%v' for property '%s' from definition override.", p.ValueType, p.PropertyName)
		return nil
	}

	uiType, _ := p.metaDetails["uitype"].(string)

	switch {
	case uiType == "bitmask":
		p.ValueType = Set
	case p.metaDetails["validateAsIPv4OrIPv6"] == true:
		p.ValueType = IpAddress
	case len(p.ValidValues) > 0 && len(p.Validators) > 0:
		p.ValueType = SemanticEquality
	default:
		if _, known := knownStringUiTypes[uiType]; !known {
			genLogger.Warnf("Unmapped meta uiType %q for property %q; defaulting to String.", uiType, p.PropertyName)
		}
		p.ValueType = String
	}

	genLogger.Debugf("Successfully set ValueType '%v' for property '%s'.", p.ValueType, p.PropertyName)
	return nil
}

func (p *Property) setSupportedVersions() error {
	// Determine the supported APIC versions for the property.
	genLogger.Debugf("Setting SupportedVersions for property '%s'.", p.PropertyName)

	// Initialize with versions from PropertyDefinition, if not defined set the versions from meta file.
	metaVersions := p.propertyDefinition.SupportedVersions
	if metaVersions == "" {
		metaVersions, _ = p.metaDetails["versions"].(string)
	}

	if metaVersions == "" {
		genLogger.Debugf("No SupportedVersions specified for property '%s'.", p.PropertyName)
		return nil
	}

	versions, err := NewVersions(metaVersions)
	if err != nil {
		return fmt.Errorf("failed to parse versions for property '%s': %w", p.PropertyName, err)
	}
	p.SupportedVersions = versions

	genLogger.Debugf("Successfully set SupportedVersions for property '%s'. Versions: '%s'", p.PropertyName, p.SupportedVersions)
	return nil
}
