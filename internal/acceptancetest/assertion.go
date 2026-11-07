package acceptancetest

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AttributeValue is a literal expectation or a reference to another attribute in state.
type AttributeValue struct {
	ExpectedValue string
	PairAddress   string
	PairAttribute string
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
