/*
 * Copyright (c) 2023 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package uuidtypes_test

import (
	// Standard Library Imports
	"context"
	"fmt"
	"testing"

	// External Imports
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	// Internal Imports
	"github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtypes"
)

func TestUUIDType_Equal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		other    attr.Type
		expected bool
	}{
		{
			name:     "nil",
			other:    nil,
			expected: false,
		},
		{
			name:     "uuidtypes.UUIDType",
			other:    uuidtypes.UUIDType{},
			expected: true,
		},
		{
			name:     "types.StringType",
			other:    types.StringType,
			expected: false,
		},
		{
			name:     "types.NumberType",
			other:    types.NumberType,
			expected: false,
		},
		{
			name:     "types.BoolType",
			other:    types.BoolType,
			expected: false,
		},
		{
			name:     "types.Float64Type",
			other:    types.Float64Type,
			expected: false,
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			if got := uuidType.Equal(testcase.other); got != testcase.expected {
				t.Errorf("Equal()\ngot     : %v\nexpected: %v", got, testcase.expected)
			}
		})
	}
}

func TestUUIDType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "always",
			expected: "uuidtypes.UUIDType",
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got := uuidType.String()

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf("String()\ngot     : %s\nexpected: %s\ndiff    : %s", got, testcase.expected, diff)
			}
		})
	}
}

func TestUUIDType_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    tftypes.Value
		path     path.Path
		expected diag.Diagnostics
	}{
		{
			name:  "not-string",
			value: tftypes.NewValue(tftypes.Bool, false),
			path:  path.Root("test"),
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid UUID Terraform Value",
					"An unexpected error occurred while attempting to read a UUID string from the Terraform value. "+
						"Please contact the provider developers with the following:\n\n"+
						"Error: can't unmarshal tftypes.Bool into *string, expected string",
				),
			},
		},
		{
			name:  "string-null",
			value: tftypes.NewValue(tftypes.String, nil),
			path:  path.Root("test"),
		},
		{
			name:  "string-unknown",
			value: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			path:  path.Root("test"),
		},
		{
			name:  "string-value-invalid-length",
			value: tftypes.NewValue(tftypes.String, valueInvalidLength),
			path:  path.Root("test"),
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid UUID String Value",
					"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
						"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
						"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
						"Provided Value: \"not-a-uuid-at-all\"\n"+
						"Parse Error: uuid string is wrong length",
				),
			},
		},
		{
			name:  "string-value-invalid-format",
			value: tftypes.NewValue(tftypes.String, valueInvalid),
			path:  path.Root("test"),
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid UUID String Value",
					"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
						"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
						"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
						"Provided Value: \"actually-not-04a00-UUID-valueat0all0\"\n"+
						"Parse Error: uuid is improperly formatted",
				),
			},
		},
		{
			name:  "string-value-valid-uuidv1",
			value: tftypes.NewValue(tftypes.String, valueUUIDv1),
			path:  path.Root("test"),
		},
		// can't find a reference UUIDv2...
		{
			name:     "string-value-valid-uuidv3",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv3),
			path:     path.Root("test"),
			expected: nil,
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv4),
			path:     path.Root("test"),
			expected: nil,
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv5),
			path:     path.Root("test"),
			expected: nil,
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got := uuidType.Validate(context.Background(), testcase.value, testcase.path)

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"Validate()\ngot     : %s\nexpected: %s\ndiff    : %s",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUIDType_ValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       tftypes.Value
		expected    attr.Value
		expectedErr error
	}{
		{
			name:        "not-string",
			value:       tftypes.NewValue(tftypes.Number, 1),
			expected:    nil,
			expectedErr: fmt.Errorf("can't unmarshal tftypes.Number into *string, expected string"),
		},
		{
			name:     "string-null",
			value:    tftypes.NewValue(tftypes.String, nil),
			expected: uuidtypes.NewUUIDNull(),
		},
		{
			name:     "string-unknown",
			value:    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: uuidtypes.NewUUIDUnknown(),
		},
		{
			name:     "string-value-invalid-value",
			value:    tftypes.NewValue(tftypes.String, valueInvalid),
			expected: uuidtypes.NewUUIDValue(valueInvalid),
		},
		{
			name:     "string-value-valid-uuidv1",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv1),
			expected: uuidtypes.NewUUIDValue(valueUUIDv1),
		},
		// can't find a reference UUIDv2...
		{
			name:     "string-value-valid-uuidv3",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv3),
			expected: uuidtypes.NewUUIDValue(valueUUIDv3),
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv4),
			expected: uuidtypes.NewUUIDValue(valueUUIDv4),
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv5),
			expected: uuidtypes.NewUUIDValue(valueUUIDv5),
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got, err := uuidType.ValueFromTerraform(context.Background(), testcase.value)
			if err != nil {
				if testcase.expectedErr == nil || err.Error() != testcase.expectedErr.Error() {
					t.Errorf(
						"ValueFromTerraform()\nerror   : %v\nexpected: %v\n",
						err,
						testcase.expectedErr,
					)
					return
				}
			}

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"ValueFromTerraform()\ngot     : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUIDType_ValueType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUIDType
		expected attr.Value
	}{
		{
			name:     "always",
			value:    uuidtypes.UUIDType{},
			expected: uuidtypes.UUIDValue{},
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.ValueType(context.Background())

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"ValueType()\ngot     : %v\nexpected: %v\ndiff: %v\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUIDType_ValueFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         basetypes.StringValue
		expected      basetypes.StringValuable
		expectedDiags diag.Diagnostics
	}{
		{
			name:     "null",
			value:    basetypes.NewStringNull(),
			expected: uuidtypes.NewUUIDNull(),
		},
		{
			name:     "unknown",
			value:    basetypes.NewStringUnknown(),
			expected: uuidtypes.NewUUIDUnknown(),
		},
		{
			name:     "invalid-value",
			value:    basetypes.NewStringValue("invalid-value"),
			expected: uuidtypes.NewUUIDValue("invalid-value"),
		},
		{
			name:     "valid-value",
			value:    basetypes.NewStringValue("invalid-value"),
			expected: uuidtypes.NewUUIDValue("invalid-value"),
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got, gotDiags := uuidType.ValueFromString(context.Background(), testcase.value)

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"ValueFromString() basetypes.StringValuable\ngot     : %v\nexpected: %v\ndiff: %v\n",
					got,
					testcase.expected,
					diff,
				)
			}

			if diff := cmp.Diff(gotDiags, testcase.expectedDiags); diff != "" {
				t.Errorf(
					"ValueFromString() diag.Diagnostics\ngot     : %v\nexpected: %v\ndiff: %v\n",
					gotDiags,
					testcase.expectedDiags,
					diff,
				)
			}
		})
	}
}
