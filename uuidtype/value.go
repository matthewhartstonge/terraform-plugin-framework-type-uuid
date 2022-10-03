/*
 * Copyright 2022 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package uuidtype

import (
	// Standard Library Imports
	"context"

	// External Imports
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure Implementation matches the expected interfaces.
var (
	_ attr.Value = Value{}
)

// NullValue returns a null UUID value.
func NullValue() Value {
	return Value{null: true}
}

// UnknownValue returns an unknown UUID value.
func UnknownValue() Value {
	return Value{unknown: true}
}

// StringValue returns a value or any errors when attempting to parse the string
// as a UUID.
func StringValue(value string, schemaPath path.Path) (Value, diag.Diagnostics) {
	validUUID, err := uuid.Parse(value)
	if err != nil {
		return UnknownValue(), diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				schemaPath,
				"Invalid UUID String Value",
				"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
					"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
					"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
					"Error: "+err.Error(),
			),
		}
	}

	return Value{
		value: validUUID.String(),
	}, nil
}

// MustValue expects a valid UUID, otherwise will panic on error.
func MustValue(value string) Value {
	validUUID, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return Value{
		value: validUUID.String(),
	}
}

// Value provides a concrete implementation of a UUID tftypes.Value for the
// Terraform Plugin framework.
type Value struct {
	null    bool
	unknown bool
	value   string
}

// Type returns the UUID type that created the Value.
func (v Value) Type(_ context.Context) attr.Type {
	return Type{}
}

// ToTerraformValue returns the UUID as a tftypes.Value.
func (v Value) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(tftypes.String, nil), nil
	}

	if v.IsUnknown() {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nil
	}

	return tftypes.NewValue(tftypes.String, v.value), nil
}

// IsNull returns true if the uuid represents a null value.
func (v Value) IsNull() bool {
	return v.null
}

// IsUnknown returns true if the uuid represents an unknown value.
func (v Value) IsUnknown() bool {
	return v.unknown
}

// Equal returns true if the uuid is semantically equal to the Value passed as
// an argument.
func (v Value) Equal(other attr.Value) bool {
	otherValue, ok := other.(Value)
	if !ok {
		return false
	}

	if otherValue.null != v.null {
		return false
	}

	if otherValue.unknown != v.unknown {
		return false
	}

	return otherValue.value == v.value
}

// String returns a summary representation of either the underlying Value,
// or UnknownValueString (`<unknown>`) when IsUnknown() returns true,
// or NullValueString (`<null>`) when IsNull() return true.
func (v Value) String() string {
	if v.null {
		return attr.NullValueString
	}

	if v.unknown {
		return attr.UnknownValueString
	}

	return v.value
}
