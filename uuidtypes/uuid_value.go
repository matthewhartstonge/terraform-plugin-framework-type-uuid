/*
 * Copyright (c) 2023 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package uuidtypes

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure Implementation matches the expected interfaces.
var (
	_ attr.Value               = UUIDValue{}
	_ basetypes.StringValuable = UUIDValue{}
)

// UUIDValue provides a concrete implementation of a UUIDValue tftypes.Value for the
// Terraform Plugin framework.
type UUIDValue struct {
	basetypes.StringValue
}

// Type returns the UUIDValue type that created the UUIDValue.
func (u UUIDValue) Type(_ context.Context) attr.Type {
	return UUIDType{}
}

// Equal returns true if the uuid is semantically equal to the Value passed as
// an argument.
func (u UUIDValue) Equal(o attr.Value) bool {
	other, ok := o.(UUIDValue)
	if !ok {
		return false
	}

	return u.StringValue.Equal(other.StringValue)
}
