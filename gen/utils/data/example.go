package data

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ExampleVariantEnum identifies which public example is being constructed.
type ExampleVariantEnum int

const (
	// MinimumExample contains only the configuration needed to create the resource.
	MinimumExample ExampleVariantEnum = iota
	// FullExample contains all configurable attributes.
	FullExample
)

// ExampleValue is a literal, expression, reference, or collection that can render itself as HCL.
type ExampleValue interface {
	HCL() string
}

// ExampleLiteral is a Terraform string literal, including a valid empty string.
type ExampleLiteral string

// HCL returns the quoted and escaped Terraform string literal.
func (v ExampleLiteral) HCL() string {
	return string(hclwrite.TokensForValue(cty.StringVal(string(v))).Bytes())
}

// ExampleExpression is an explicitly authored Terraform expression.
type ExampleExpression string

// HCL returns the expression unchanged so Terraform can evaluate it.
func (v ExampleExpression) HCL() string {
	return string(v)
}

// ExampleReference identifies an attribute on the standard example resource or data source.
// ResourceName excludes the aci_ prefix.
type ExampleReference struct {
	Type          ReferenceTypeEnum
	ResourceName  string
	AttributeName string
}

// HCL returns the Terraform attribute path represented by the reference.
func (r ExampleReference) HCL() string {
	prefix := ""
	if r.Type == DataSourceReference {
		prefix = "data."
	}
	return fmt.Sprintf("%saci_%s.example.%s", prefix, r.ResourceName, r.AttributeName)
}

// ExampleCollection is an ordered collection of example values.
type ExampleCollection []ExampleValue

// HCL returns the collection as a Terraform collection expression.
func (v ExampleCollection) HCL() string {
	values := make([]string, 0, len(v))
	for _, value := range v {
		values = append(values, value.HCL())
	}
	return "[" + strings.Join(values, ", ") + "]"
}

// ExampleAttribute is one ordered attribute in an example object.
type ExampleAttribute struct {
	Name  string
	Value ExampleValue
}

// ExampleObject is the recursive body shared by resources and nested child objects.
type ExampleObject struct {
	Attributes []ExampleAttribute
	Children   []ExampleChild
}

// ExampleChild describes one nested object attribute and its cardinality.
type ExampleChild struct {
	Name     string
	IsSingle bool
	Body     ExampleObject
}

// Example is one complete public example for a representative parent placement.
type Example struct {
	Platform PlatformTypeEnum
	Body     ExampleObject
}

// DataSourceHCL renders the complete, formatter-normalized data-source block
// shared by the standalone example and generated documentation.
func (e Example) DataSourceHCL(resourceName string) string {
	var contents strings.Builder
	if e.Platform == Cloud {
		contents.WriteString("# This example is only applicable to Cisco Cloud Network Controller.\n")
	}
	fmt.Fprintf(&contents, "data %q %q {\n", "aci_"+resourceName, "example")
	writeExampleObjectHCL(&contents, e.Body)
	contents.WriteString("}\n")
	return strings.TrimSuffix(string(hclwrite.Format([]byte(contents.String()))), "\n")
}

func writeExampleObjectHCL(contents *strings.Builder, object ExampleObject) {
	for _, attribute := range object.Attributes {
		fmt.Fprintf(contents, "%s = %s\n", attribute.Name, attribute.Value.HCL())
	}
	for _, child := range object.Children {
		if child.IsSingle {
			fmt.Fprintf(contents, "%s = {\n", child.Name)
			writeExampleObjectHCL(contents, child.Body)
			contents.WriteString("}\n")
			continue
		}
		fmt.Fprintf(contents, "%s = [{\n", child.Name)
		writeExampleObjectHCL(contents, child.Body)
		contents.WriteString("}]\n")
	}
}

// ResourceExamples contains both resource artifacts and files copied beside them.
type ResourceExamples struct {
	Minimum Example
	Full    Example
	Files   []string
}

// HasDistinctFullExample reports whether the all-attributes example adds
// anything to the minimum example.
func (e ResourceExamples) HasDistinctFullExample() bool {
	return !reflect.DeepEqual(e.Minimum, e.Full)
}

func newExampleReference(referenceType ReferenceTypeEnum, resourceName, attributeName string) ExampleReference {
	return ExampleReference{
		Type:          referenceType,
		ResourceName:  resourceName,
		AttributeName: attributeName,
	}
}
