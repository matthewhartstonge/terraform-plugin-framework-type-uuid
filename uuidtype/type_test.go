package uuidtype_test

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
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	// Internal Imports
	"github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtype"
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
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.AttributeName to uuidtype.Type"),
		},
		{
			name:        "tftypes.ElementKeyInt",
			step:        tftypes.ElementKeyInt(0),
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyInt to uuidtype.Type"),
		},
		{
			name:        "tftypes.ElementKeyString",
			step:        tftypes.ElementKeyString("test"),
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyString to uuidtype.Type"),
		},
		{
			name:        "tftypes.ElementKeyValue",
			step:        tftypes.ElementKeyValue{},
			expectedErr: fmt.Errorf("cannot apply AttributePathStep to tftypes.ElementKeyValue to uuidtype.Type"),
		},
	}
	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtype.Type{}
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
			name:     "uuidtype.Type",
			other:    uuidtype.Type{},
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

			uuidType := uuidtype.Type{}
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
		typ      uuidtype.Type
		expected string
	}{
		{
			name:     "always",
			expected: "uuidtype.Type",
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtype.Type{}
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

			uuidType := uuidtype.Type{}
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
			value: tftypes.NewValue(tftypes.String, "4ea3c666-4309-11ed-b878-0242ac120002"),
			path:  path.Root("test"),
		},
		// can't find a reference UUIDv2...
		{
			name:     "string-value-valid-uuidv3",
			value:    tftypes.NewValue(tftypes.String, "a825d19e-3885-3df7-920a-a3678f53b2ee"),
			path:     path.Root("test"),
			expected: nil,
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    tftypes.NewValue(tftypes.String, "eb6f148a-6637-4c6b-a4bb-b75b2a1b5a3c"),
			path:     path.Root("test"),
			expected: nil,
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    tftypes.NewValue(tftypes.String, "f989a266-a679-5f41-92f7-22004c4da817"),
			path:     path.Root("test"),
			expected: nil,
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtype.Type{}
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
			expected:    uuidtype.UnknownValue(),
			expectedErr: fmt.Errorf("can't unmarshal tftypes.Number into *string, expected string"),
		},
		{
			name:     "string-null",
			value:    tftypes.NewValue(tftypes.String, nil),
			expected: uuidtype.NullValue(),
		},
		{
			name:     "string-unknown",
			value:    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: uuidtype.UnknownValue(),
		},
		{
			name:        "string-value-invalid-length",
			value:       tftypes.NewValue(tftypes.String, "not-a-uuid-at-all"),
			expected:    uuidtype.UnknownValue(),
			expectedErr: fmt.Errorf("invalid UUID length: 17"),
		},
		{
			name:        "string-value-invalid-value",
			value:       tftypes.NewValue(tftypes.String, "actually-not0-4a00-UUID-at0all00"),
			expected:    uuidtype.UnknownValue(),
			expectedErr: fmt.Errorf("invalid UUID format"),
		},
		{
			name:     "string-value-valid-uuidv1",
			value:    tftypes.NewValue(tftypes.String, "4ea3c666-4309-11ed-b878-0242ac120002"),
			expected: uuidtype.MustValue("4ea3c666-4309-11ed-b878-0242ac120002"),
		},
		// can't find a reference UUIDv2...
		{
			name:     "string-value-valid-uuidv3",
			value:    tftypes.NewValue(tftypes.String, "a825d19e-3885-3df7-920a-a3678f53b2ee"),
			expected: uuidtype.MustValue("a825d19e-3885-3df7-920a-a3678f53b2ee"),
		},
		{
			name:     "string-value-valid-uuidv4",
			value:    tftypes.NewValue(tftypes.String, "eb6f148a-6637-4c6b-a4bb-b75b2a1b5a3c"),
			expected: uuidtype.MustValue("eb6f148a-6637-4c6b-a4bb-b75b2a1b5a3c"),
		},
		{
			name:     "string-value-valid-uuidv5",
			value:    tftypes.NewValue(tftypes.String, "f989a266-a679-5f41-92f7-22004c4da817"),
			expected: uuidtype.MustValue("f989a266-a679-5f41-92f7-22004c4da817"),
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			uuidType := uuidtype.Type{}
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
					"ValueFromTerraform()\nvalue   : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					diff,
				)
			}
		})
	}
}
