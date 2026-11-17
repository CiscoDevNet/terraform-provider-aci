package acceptancetest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ResourceState tracks the last successfully applied identity across lifecycle steps.
type ResourceState struct {
	ID string
}

func (*ResourceState) CheckRemoved(address string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if root := state.RootModule(); root != nil && root.Resources[address] != nil {
			return fmt.Errorf("%s remains in Terraform state after removal", address)
		}
		return nil
	}
}

func (s *ResourceState) CheckID(address string, changed bool) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(address, "id", func(id string) error {
		if id == "" {
			return fmt.Errorf("%s has an empty ID", address)
		}
		if s.ID != "" && (id != s.ID) != changed {
			return fmt.Errorf("%s ID: previous %q, current %q, expected change: %t", address, s.ID, id, changed)
		}
		s.ID = id
		return nil
	})
}
