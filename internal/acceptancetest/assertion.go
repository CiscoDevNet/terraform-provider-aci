package acceptancetest

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AttributeValue is a literal expectation or a reference to another attribute in state.
type AttributeValue struct {
	ExpectedValue string
	PairAddress   string
	PairAttribute string
	PairFormat    string
}

func CheckSemanticStringAttribute(address, path, expected string, valueType basetypes.StringTypable) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Empty is a known value; checking attribute presence must not reject it.
		root := state.RootModule()
		if root == nil || root.Resources[address] == nil || root.Resources[address].Primary == nil {
			return fmt.Errorf("resource %q is missing from state", address)
		}
		actual, exists := root.Resources[address].Primary.Attributes[path]
		if !exists {
			return fmt.Errorf("%s: attribute %q is missing from state", address, path)
		}
		ctx := context.Background()
		actualValue, diagnostics := valueType.ValueFromString(ctx, basetypes.NewStringValue(actual))
		if diagnostics.HasError() {
			return fmt.Errorf("%s.%s: convert actual value: %v", address, path, diagnostics)
		}
		expectedValue, diagnostics := valueType.ValueFromString(ctx, basetypes.NewStringValue(expected))
		if diagnostics.HasError() {
			return fmt.Errorf("%s.%s: convert expected value: %v", address, path, diagnostics)
		}
		semantic, ok := actualValue.(basetypes.StringValuableWithSemanticEquals)
		if !ok {
			return fmt.Errorf("%s.%s: type %T does not support string semantic equality", address, path, valueType)
		}
		equal, diagnostics := semantic.StringSemanticEquals(ctx, expectedValue)
		if diagnostics.HasError() {
			return fmt.Errorf("%s.%s: compare semantic values: %v", address, path, diagnostics)
		}
		if !equal {
			return fmt.Errorf("%s.%s: expected semantic value %q, got %q", address, path, expected, actual)
		}
		return nil
	}
}

// CheckAttribute checks a literal or a formatted reference against a scalar attribute.
func CheckAttribute(address, path string, value AttributeValue) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		check := nestedSetCheck{address: address, state: state}
		resolved, err := check.resolveValues(map[string]AttributeValue{"value": value})
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(address, path, resolved["value"])(state)
	}
}

// CheckPresent checks a non-empty value within its identified parent objects.
func CheckPresent(address, path string, parents ...ParentSetElement) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		check := nestedSetCheck{address: address, state: state}
		return check.checkWithinParents(path, parents, func(scopedPath string) resource.TestCheckFunc {
			return resource.TestCheckResourceAttrSet(address, scopedPath)
		})
	}
}

// CheckAbsent verifies that a singleton has no values or populated descendants.
func CheckAbsent(address, path string, parents ...ParentSetElement) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		check := nestedSetCheck{address: address, state: state}
		return check.checkWithinParents(path, parents, func(scopedPath string) resource.TestCheckFunc {
			return func(state *terraform.State) error {
				root := state.RootModule()
				if root == nil || root.Resources[address] == nil || root.Resources[address].Primary == nil {
					return fmt.Errorf("resource %q is missing from state", address)
				}
				for name, value := range root.Resources[address].Primary.Attributes {
					if name == scopedPath || strings.HasPrefix(name, scopedPath+".") {
						// Objects with all-null attributes retain their field count in flattened state.
						if strings.HasSuffix(name, ".%") || strings.HasSuffix(name, ".#") && value == "0" {
							continue
						}
						return fmt.Errorf("%s.%s should be absent, found %s=%q", address, scopedPath, name, value)
					}
				}
				return nil
			}
		})
	}
}

// ParentSetElement identifies one enclosing set element before checking its descendants.
type ParentSetElement struct {
	Path   string
	Values map[string]AttributeValue
}

type nestedSetCheck struct {
	address string
	state   *terraform.State
}

// CheckNestedSetElement checks literals and resolved references within the same child object.
func CheckNestedSetElement(address, path string, values map[string]AttributeValue, parents ...ParentSetElement) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Resolve anew for each check so later test steps cannot reuse stale state values.
		check := nestedSetCheck{address: address, state: state}
		resolved, err := check.resolveValues(values)
		if err != nil {
			return err
		}
		return check.checkWithinParents(path, parents, func(scopedPath string) resource.TestCheckFunc {
			return resource.TestCheckTypeSetElemNestedAttrs(address, scopedPath, resolved)
		})
	}
}

// CheckSetElement checks a primitive member within its identified parent objects.
func CheckSetElement(address, path string, value AttributeValue, parents ...ParentSetElement) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		check := nestedSetCheck{address: address, state: state}
		resolved, err := check.resolveValues(map[string]AttributeValue{"value": value})
		if err != nil {
			return err
		}
		return check.checkWithinParents(path, parents, func(scopedPath string) resource.TestCheckFunc {
			return resource.TestCheckTypeSetElemAttr(address, scopedPath, resolved["value"])
		})
	}
}

// CheckSetPair compares collection counts and members without relying on matching state hashes.
func CheckSetPair(address, path, pairAddress, pairPath string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if err := resource.TestCheckResourceAttrPair(address, path+".#", pairAddress, pairPath+".#")(state); err != nil {
			return err
		}
		actualElements := nestedSetCheck{address: address, state: state}.setElements(path)
		expectedElements := nestedSetCheck{address: pairAddress, state: state}.setElements(pairPath)
		for _, index := range slices.Sorted(maps.Keys(expectedElements)) {
			expected := expectedElements[index]
			matched := false
			for _, hash := range slices.Sorted(maps.Keys(actualElements)) {
				if maps.Equal(actualElements[hash], expected) {
					delete(actualElements, hash)
					matched = true
					break
				}
			}
			if !matched {
				return fmt.Errorf("%s: no TypeSet element %q with exact attributes %#v", address, path+".*", expected)
			}
		}
		return nil
	}
}

func (c nestedSetCheck) setElements(path string) map[string]map[string]string {
	elements := map[string]map[string]string{}
	for name, value := range c.state.RootModule().Resources[c.address].Primary.Attributes {
		remainder, found := strings.CutPrefix(name, path+".")
		if !found || remainder == "#" {
			continue
		}
		// Strip only the outer state hash; an empty member name holds a scalar value.
		index, member, _ := strings.Cut(remainder, ".")
		if elements[index] == nil {
			elements[index] = map[string]string{}
		}
		elements[index][member] = value
	}
	return elements
}

func (c nestedSetCheck) resolveValues(values map[string]AttributeValue) (map[string]string, error) {
	resolved := make(map[string]string, len(values))
	for _, name := range slices.Sorted(maps.Keys(values)) {
		value := values[name]
		if value.PairAddress == "" {
			resolved[name] = value.ExpectedValue
			continue
		}
		check := resource.TestCheckResourceAttrWith(value.PairAddress, value.PairAttribute, func(actual string) error {
			resolved[name] = actual
			return nil
		})
		if err := check(c.state); err != nil {
			return nil, fmt.Errorf("resolve %s.%s for %s.%s: %w", value.PairAddress, value.PairAttribute, c.address, name, err)
		}
		if value.PairFormat != "" {
			resolved[name] = fmt.Sprintf(value.PairFormat, resolved[name])
		}
	}
	return resolved, nil
}

func (c nestedSetCheck) checkWithinParents(path string, parents []ParentSetElement, check func(string) resource.TestCheckFunc) error {
	if len(parents) == 0 {
		return check(path)(c.state)
	}
	parent := parents[0]
	resolved, err := c.resolveValues(parent.Values)
	if err != nil {
		return err
	}
	paths, err := c.parentPaths(parent.Path, resolved)
	if err != nil {
		return err
	}
	var lastError error
	for _, parentPath := range paths {
		// Replace each ancestor wildcard with its actual state index before checking
		// another set. Matching hashes under different parents must never be combined.
		scopedPath := strings.Replace(path, parent.Path, parentPath, 1)
		scopedParents := make([]ParentSetElement, len(parents)-1)
		for index, descendant := range parents[1:] {
			scopedParents[index] = descendant
			scopedParents[index].Path = strings.Replace(descendant.Path, parent.Path, parentPath, 1)
		}
		if err := c.checkWithinParents(scopedPath, scopedParents, check); err == nil {
			return nil
		} else {
			lastError = err
		}
	}
	if lastError != nil {
		return fmt.Errorf("check %s within %s: %w", path, parent.Path, lastError)
	}
	return fmt.Errorf("%q no TypeSet element %q with parent attributes %#v", c.address, parent.Path, resolved)
}

func (c nestedSetCheck) parentPaths(path string, values map[string]string) ([]string, error) {
	root := c.state.RootModule()
	if root == nil || root.Resources[c.address] == nil || root.Resources[c.address].Primary == nil {
		return nil, fmt.Errorf("resource %q is missing from state", c.address)
	}
	attributes := root.Resources[c.address].Primary.Attributes
	prefix := strings.TrimSuffix(path, "*")
	candidates := make(map[string]bool)
	for key := range attributes {
		remainder, ok := strings.CutPrefix(key, prefix)
		if !ok {
			continue
		}
		index, _, ok := strings.Cut(remainder, ".")
		if ok {
			candidates[prefix+index] = true
		}
	}
	var paths []string
	for _, candidate := range slices.Sorted(maps.Keys(candidates)) {
		matches := true
		for name, value := range values {
			if attributes[candidate+"."+name] != value {
				matches = false
				break
			}
		}
		if matches {
			paths = append(paths, candidate)
		}
	}
	return paths, nil
}
