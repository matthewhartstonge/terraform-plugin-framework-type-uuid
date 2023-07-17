/*
 * Copyright (c) 2022 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package uuidtypes

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure Implementation matches the expected interfaces.
var (
	_ attr.Value = UUID{}
)

// UUIDNull returns a null UUID value.
func UUIDNull() UUID {
	return UUID{state: attr.ValueStateNull}
}

// UUIDUnknown returns an unknown UUID value.
func UUIDUnknown() UUID {
	return UUID{state: attr.ValueStateUnknown}
}

// UUIDFromString returns a value or any errors when attempting to parse the
// string as a UUID.
func UUIDFromString(value string, schemaPath path.Path) (UUID, diag.Diagnostics) {
	validUUID, err := uuid.Parse(value)
	if err != nil {
		return UUIDUnknown(), diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				schemaPath,
				"Invalid UUID String Value",
				"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
					"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
					"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
					fmt.Sprintf("Provided Value: %s\n", value)+
					fmt.Sprintf("Parse Error: %s", err.Error()),
			),
		}
	}

	return UUID{
		state: attr.ValueStateKnown,
		value: validUUID,
	}, nil
}

// UUIDFromGoogleUUID expects a valid google/uuid.UUID and returns a Terraform
// UUID Value.
func UUIDFromGoogleUUID(value uuid.UUID) UUID {
	return UUID{
		state: attr.ValueStateKnown,
		value: value,
	}
}

// UUID provides a concrete implementation of a UUID tftypes.Value for the
// Terraform Plugin framework.
type UUID struct {
	state attr.ValueState
	value uuid.UUID
}

// Type returns the UUID type that created the UUID.
func (u UUID) Type(_ context.Context) attr.Type {
	return UUIDType{}
}

// ToTerraformValue returns the UUID as a tftypes.Value.
func (u UUID) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if u.IsNull() {
		return tftypes.NewValue(tftypes.String, nil), nil
	}

	if u.IsUnknown() {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nil
	}

	return tftypes.NewValue(tftypes.String, u.value.String()), nil
}

// IsNull returns true if the uuid represents a null value.
func (u UUID) IsNull() bool {
	return u.state == attr.ValueStateNull
}

// IsUnknown returns true if the uuid represents an unknown value.
func (u UUID) IsUnknown() bool {
	return u.state == attr.ValueStateUnknown
}

// Equal returns true if the uuid is semantically equal to the Value passed as
// an argument.
func (u UUID) Equal(other attr.Value) bool {
	otherValue, ok := other.(UUID)
	if !ok {
		return false
	}

	if otherValue.state != u.state {
		return false
	}

	// perform a byte-for-byte comparison.
	return otherValue.value == u.value
}

// String returns a summary representation of either the underlying Value,
// or UnknownValueString (`<unknown>`) when IsUnknown() returns true,
// or NullValueString (`<null>`) when IsNull() return true.
func (u UUID) String() string {
	switch u.state {
	case attr.ValueStateNull:
		return attr.NullValueString

	case attr.ValueStateUnknown:
		return attr.UnknownValueString

	default:
		return u.value.String()
	}
}
