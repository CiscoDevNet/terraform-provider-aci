// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6dynamicvalue

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Must creates a *tfprotov6.DynamicValue or panics. This is intended only for
// simplifying testing code.
//
// The tftypes.Type parameter is separate to enable DynamicPsuedoType testing.
func Must(typ tftypes.Type, value tftypes.Value) *tfprotov6.DynamicValue {
	dynamicValue, err := tfprotov6.NewDynamicValue(typ, value)

	if err != nil {
		panic(fmt.Sprintf("unable to create DynamicValue: %s", err.Error()))
	}

	return &dynamicValue
}
