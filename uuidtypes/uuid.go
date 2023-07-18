/*
 * Copyright (c) 2023 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package uuidtypes

import (
	// External Imports
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type UUID = UUIDValue

// NewUUIDNull creates a UUID with a null value. Determine whether the value is
// null via the UUID type IsNull method.
func NewUUIDNull() UUIDValue {
	return UUIDValue{StringValue: basetypes.NewStringNull()}
}

// NewUUIDUnknown creartes a UUID with an unknown UUID value. Determine whether
// the value is unknown via the UUID type IsUnknown method.
func NewUUIDUnknown() UUIDValue {
	return UUIDValue{StringValue: basetypes.NewStringUnknown()}
}

// NewUUIDValue creates a UUID with a known value. Access the value via the
// String type ValueString method.
func NewUUIDValue(value string) UUIDValue {
	return UUIDValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewUUIDPointerValue creates a UUID with a null value if nil or a known
// value. Access the value via the String type ValueStringPointer method.
func NewUUIDPointerValue(value *string) UUIDValue {
	return UUIDValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
