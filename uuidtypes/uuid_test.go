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

const (
	valueInvalid       = "actually-not0-4a00-UUID-at0all00"
	valueInvalidLength = "not-a-uuid-at-all"
	valueUUIDv1        = "4ea3c666-4309-11ed-b878-0242ac120002"
	valueUUIDv3        = "a825d19e-3885-3df7-920a-a3678f53b2ee"
	valueUUIDv4        = "eb6f148a-6637-4c6b-a4bb-b75b2a1b5a3c"
	valueUUIDv5        = "f989a266-a679-5f41-92f7-22004c4da817"
)

func TestUUID_Equal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected bool
		value    attr.Value
		other    attr.Value
	}{
		// Null
		{
			name:     "null-nil",
			value:    uuidtypes.UUIDNull(),
			other:    nil,
			expected: false,
		},
		{
			name:     "null-null",
			value:    uuidtypes.UUIDNull(),
			other:    uuidtypes.UUIDNull(),
			expected: true,
		},
		{
			name:     "null-unknown",
			value:    uuidtypes.UUIDNull(),
			other:    uuidtypes.UUIDUnknown(),
			expected: false,
		},
		{
			name:     "null-value",
			value:    uuidtypes.UUIDNull(),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: false,
		},
		{
			name:     "null-different-value",
			value:    uuidtypes.UUIDNull(),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: false,
		},

		// Unknown
		{
			name:     "unknown-nil",
			value:    uuidtypes.UUIDUnknown(),
			other:    nil,
			expected: false,
		},
		{
			name:     "unknown-null",
			value:    uuidtypes.UUIDUnknown(),
			other:    uuidtypes.UUIDNull(),
			expected: false,
		},
		{
			name:     "unknown-unknown",
			value:    uuidtypes.UUIDUnknown(),
			other:    uuidtypes.UUIDUnknown(),
			expected: true,
		},
		{
			name:     "unknown-value",
			value:    uuidtypes.UUIDUnknown(),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: false,
		},
		{
			name:     "unknown-different-value",
			value:    uuidtypes.UUIDUnknown(),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: false,
		},

		// Value
		{
			name:     "value-nil",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    nil,
			expected: false,
		},
		{
			name:     "value-null",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    uuidtypes.UUIDNull(),
			expected: false,
		},
		{
			name:     "value-unknown",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    uuidtypes.UUIDUnknown(),
			expected: false,
		},
		{
			name:     "value-value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: true,
		},
		{
			name:     "value-different-value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: false,
		},
		{
			name:     "not-uuidtypes.UUID",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			other:    types.String{Value: valueUUIDv4},
			expected: false,
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.Equal(testcase.other)

			if got != testcase.expected {
				t.Errorf("Equal()\ngot     : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					cmp.Diff(testcase.value, testcase.other),
				)
			}
		})
	}
}

func TestUUID_IsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUID
		expected bool
	}{
		{
			name:     "null",
			value:    uuidtypes.UUIDNull(),
			expected: true,
		},
		{
			name:     "unknown",
			value:    uuidtypes.UUIDUnknown(),
			expected: false,
		},
		{
			name:     "value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: false,
		},
		{
			name:     "other-value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: false,
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.IsNull()

			if got != testcase.expected {
				t.Errorf("IsNull() = %v, want %v", got, testcase.expected)
			}
		})
	}
}

func TestUUID_IsUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUID
		expected bool
	}{
		{
			name:     "null",
			value:    uuidtypes.UUIDNull(),
			expected: false,
		},
		{
			name:     "unknown",
			value:    uuidtypes.UUIDUnknown(),
			expected: true,
		},
		{
			name:     "value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: false,
		},
		{
			name:     "other-value",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: false,
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.IsUnknown()

			if got != testcase.expected {
				t.Errorf("IsUnknown() = %v, want %v", got, testcase.expected)
			}
		})
	}
}

func TestUUID_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUID
		expected string
	}{
		{
			name:     "null",
			value:    uuidtypes.UUIDNull(),
			expected: attr.NullValueString,
		},
		{
			name:     "unknown",
			value:    uuidtypes.UUIDUnknown(),
			expected: attr.UnknownValueString,
		},
		{
			name:     "uuidv1",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv1)),
			expected: valueUUIDv1,
		},
		{
			name:     "uuidv3",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv3)),
			expected: valueUUIDv3,
		},
		{
			name:     "uuidv4",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: valueUUIDv4,
		},
		{
			name:     "uuidv5",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: valueUUIDv5,
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.String()

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"String()\ngot     : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUID_ToTerraformValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       uuidtypes.UUID
		expected    tftypes.Value
		expectedErr error
	}{
		{
			name:     "null",
			value:    uuidtypes.UUIDNull(),
			expected: tftypes.NewValue(tftypes.String, nil),
		},
		{
			name:     "unknown",
			value:    uuidtypes.UUIDUnknown(),
			expected: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		{
			name:     "string-value-valid-uuidv1",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv1)),
			expected: tftypes.NewValue(tftypes.String, valueUUIDv1),
		},
		{
			name:     "string-value-valid-uuidv3",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv3)),
			expected: tftypes.NewValue(tftypes.String, valueUUIDv3),
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv4)),
			expected: tftypes.NewValue(tftypes.String, valueUUIDv4),
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    uuidtypes.UUIDFromGoogleUUID(uuid.MustParse(valueUUIDv5)),
			expected: tftypes.NewValue(tftypes.String, valueUUIDv5),
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := testcase.value.ToTerraformValue(context.Background())
			if err != nil {
				if testcase.expectedErr == nil || err.Error() != testcase.expectedErr.Error() {
					t.Errorf(
						"ToTerraformValue()\nerror   : %v\nexpected: %v\n",
						err,
						testcase.expectedErr,
					)
					return
				}
			}

			if err == nil && testcase.expectedErr != nil {
				t.Errorf(
					"ToTerraformValue()\nerror   : %v\nexpected: %v\n",
					err,
					testcase.expectedErr,
				)
			}

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"ToTerraformValue()\ngot     : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUID_Type(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUID
		expected attr.Type
	}{
		{
			name:     "always",
			value:    uuidtypes.UUIDNull(),
			expected: uuidtypes.UUIDType{},
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := testcase.value.Type(context.Background())

			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf(
					"Type()\ngot     : %v\nexpected: %v\ndiff: %v\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUIDFromGoogleUUID(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name     string
		value    uuid.UUID
		expected string
		panic    bool
	}{
		{
			name:     "value-uuidv1",
			value:    uuid.MustParse(valueUUIDv1),
			expected: valueUUIDv1,
		},
		{
			name:     "value-uuidv3",
			value:    uuid.MustParse(valueUUIDv3),
			expected: valueUUIDv3,
		},
		{
			name:     "value-uuidv4",
			value:    uuid.MustParse(valueUUIDv4),
			expected: valueUUIDv4,
		},
		{
			name:     "value-uuidv5",
			value:    uuid.MustParse(valueUUIDv5),
			expected: valueUUIDv5,
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			got := uuidtypes.UUIDFromGoogleUUID(testcase.value)

			if diff := cmp.Diff(got.String(), testcase.expected); diff != "" {
				t.Errorf("UUIDFromGoogleUUID()\ngot     : %vexpected: %v\ndiff: %v",
					got.String(),
					testcase.expected,
					diff,
				)
			}
		})
	}
}

func TestUUIDFromString(t *testing.T) {
	t.Parallel()

	expectedSummary := "Invalid UUID String Value"
	expectedDetail :=
		"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. " +
			"The expected UUID format is 00000000-0000-0000-0000-00000000. " +
			"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"

	tests := []struct {
		name          string
		value         string
		schemaPath    path.Path
		expectedUUID  uuidtypes.UUID
		expectedDiags diag.Diagnostics
	}{
		{
			name:         "string-value-empty",
			value:        "",
			schemaPath:   path.Root("test"),
			expectedUUID: uuidtypes.UUIDUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					expectedSummary,
					expectedDetail+"Error: invalid UUID length: 0",
				),
			},
		},
		{
			name:         "string-value-invalid-length",
			value:        valueInvalidLength,
			schemaPath:   path.Root("test"),
			expectedUUID: uuidtypes.UUIDUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					expectedSummary,
					expectedDetail+"Error: invalid UUID length: 17",
				),
			},
		},
		{
			name:         "string-value-invalid-format",
			value:        valueInvalid,
			expectedUUID: uuidtypes.UUIDUnknown(),
			schemaPath:   path.Root("test"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					expectedSummary,
					expectedDetail+"Error: invalid UUID format",
				),
			},
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			gotUUID, gotDiags := uuidtypes.UUIDFromString(testcase.value, testcase.schemaPath)
			if diff := cmp.Diff(gotUUID, testcase.expectedUUID); diff != "" {
				t.Errorf("UUIDFromString() got = %v, want %v", gotUUID, testcase.expectedUUID)
			}
			if diff := cmp.Diff(gotDiags, testcase.expectedDiags); diff != "" {
				t.Errorf("UUIDFromString() got = %v, want %v", gotDiags, testcase.expectedDiags)
			}
		})
	}
}
