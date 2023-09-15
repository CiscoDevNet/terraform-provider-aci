// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6dynamicvalue

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Equals performs equality checking of two given *tfprotov6.DynamicValue.
func Equals(schemaType tftypes.Type, i *tfprotov6.DynamicValue, j *tfprotov6.DynamicValue) (bool, error) {
	if i == nil {
		return j == nil, nil
	}

	if j == nil {
		return false, nil
	}

	// Upstream will panic on DynamicValue.Unmarshal with nil Type
	if schemaType == nil {
		return false, fmt.Errorf("unable to unmarshal DynamicValue: missing Type")
	}

	iValue, err := i.Unmarshal(schemaType)

	if err != nil {
		return false, fmt.Errorf("unable to unmarshal DynamicValue: %w", err)
	}

	jValue, err := j.Unmarshal(schemaType)

	if err != nil {
		return false, fmt.Errorf("unable to unmarshal DynamicValue: %w", err)
	}

	return iValue.Equal(jValue), nil
}
