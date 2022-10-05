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

package uuidtypes_test

import (
	// Standard Library Imports
	"context"
	"fmt"
	"testing"

	// External Imports
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	// Internal Imports
	"github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtypes"
)

func TestType_ApplyTerraform5AttributePathStep(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		step        tftypes.AttributePathStep
		expected    any
		expectedErr error
	}{
		{
			name:        "tftypes.AttributeName",
			step:        tftypes.AttributeName("test"),
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.AttributeName to uuidtypes.UUIDType"),
		},
		{
			name:        "tftypes.ElementKeyInt",
			step:        tftypes.ElementKeyInt(0),
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyInt to uuidtypes.UUIDType"),
		},
		{
			name:        "tftypes.ElementKeyString",
			step:        tftypes.ElementKeyString("test"),
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyString to uuidtypes.UUIDType"),
		},
		{
			name:        "tftypes.ElementKeyValue",
			step:        tftypes.ElementKeyValue{},
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyValue to uuidtypes.UUIDType"),
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got, err := uuidType.ApplyTerraform5AttributePathStep(testcase.step)
			if err != nil {
				if testcase.expectedErr == nil || err.Error() != testcase.expectedErr.Error() {
					t.Errorf(
						"ApplyTerraform5AttributePathStep()\nerror   : %v\nexpected: %v\n",
						err,
						testcase.expectedErr,
					)
					return
				}
			}

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"ApplyTerraform5AttributePathStep()\nvalue   : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestType_Equal(t *testing.T) {
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

func TestType_String(t *testing.T) {
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

func TestType_TerraformType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected tftypes.Type
	}{
		{
			name:     "always",
			expected: tftypes.String,
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtypes.UUIDType{}
			got := uuidType.TerraformType(context.Background())

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"TerraformType()\ngot     : %s\nexpected: %s\ndiff    : %s",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestType_Validate(t *testing.T) {
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
			value: tftypes.NewValue(tftypes.String, "not-a-uuid-at-all"),
			path:  path.Root("test"),
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid UUID String Value",
					"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
						"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
						"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
						"Error: invalid UUID length: 17",
				),
			},
		},
		{
			name:  "string-value-invalid-format",
			value: tftypes.NewValue(tftypes.String, "actually-not0-4a00-UUID-at0all00"),
			path:  path.Root("test"),
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid UUID String Value",
					"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
						"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
						"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
						"Error: invalid UUID format",
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

func TestType_ValueFromTerraform(t *testing.T) {
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
			expected:    uuidtypes.UUIDUnknown(),
			expectedErr: fmt.Errorf("can't unmarshal tftypes.Number into *string, expected string"),
		},
		{
			name:     "string-null",
			value:    tftypes.NewValue(tftypes.String, nil),
			expected: uuidtypes.UUIDNull(),
		},
		{
			name:     "string-unknown",
			value:    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: uuidtypes.UUIDUnknown(),
		},
		{
			name:        "string-value-invalid-length",
			value:       tftypes.NewValue(tftypes.String, valueInvalidLength),
			expected:    uuidtypes.UUIDUnknown(),
			expectedErr: fmt.Errorf("invalid UUID length: 17"),
		},
		{
			name:        "string-value-invalid-value",
			value:       tftypes.NewValue(tftypes.String, valueInvalid),
			expected:    uuidtypes.UUIDUnknown(),
			expectedErr: fmt.Errorf("invalid UUID format"),
		},
		{
			name:     "string-value-valid-uuidv1",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv1),
			expected: uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv1)),
		},
		// can't find a reference UUIDv2...
		{
			name:     "string-value-valid-uuidv3",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv3),
			expected: uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv3)),
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv4),
			expected: uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    tftypes.NewValue(tftypes.String, valueUUIDv5),
			expected: uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
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
