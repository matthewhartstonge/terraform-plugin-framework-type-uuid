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
	"testing"

	// External Imports
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// Internal Imports
	"github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtypes"
)

const (
	valueInvalid       = "actually-not-04a00-UUID-valueat0all0"
	valueInvalidLength = "not-a-uuid-at-all"
	valueUUIDv1        = "4ea3c666-4309-11ed-b878-0242ac120002"
	valueUUIDv3        = "a825d19e-3885-3df7-920a-a3678f53b2ee"
	valueUUIDv4        = "eb6f148a-6637-4c6b-a4bb-b75b2a1b5a3c"
	valueUUIDv5        = "f989a266-a679-5f41-92f7-22004c4da817"
)

func TestUUIDValue_Equal(t *testing.T) {
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
			value:    uuidtypes.NewUUIDNull(),
			other:    nil,
			expected: false,
		},
		{
			name:     "null-null",
			value:    uuidtypes.NewUUIDNull(),
			other:    uuidtypes.NewUUIDNull(),
			expected: true,
		},
		{
			name:     "null-unknown",
			value:    uuidtypes.NewUUIDNull(),
			other:    uuidtypes.NewUUIDUnknown(),
			expected: false,
		},
		{
			name:     "null-value",
			value:    uuidtypes.NewUUIDNull(),
			other:    uuidtypes.NewUUIDValue(valueUUIDv4),
			expected: false,
		},
		{
			name:     "null-different-value",
			value:    uuidtypes.NewUUIDNull(),
			other:    uuidtypes.NewUUIDValue(valueUUIDv5),
			expected: false,
		},

		// Unknown
		{
			name:     "unknown-nil",
			value:    uuidtypes.NewUUIDUnknown(),
			other:    nil,
			expected: false,
		},
		{
			name:     "unknown-null",
			value:    uuidtypes.NewUUIDUnknown(),
			other:    uuidtypes.NewUUIDNull(),
			expected: false,
		},
		{
			name:     "unknown-unknown",
			value:    uuidtypes.NewUUIDUnknown(),
			other:    uuidtypes.NewUUIDUnknown(),
			expected: true,
		},
		{
			name:     "unknown-value",
			value:    uuidtypes.NewUUIDUnknown(),
			other:    uuidtypes.NewUUIDValue(valueUUIDv4),
			expected: false,
		},
		{
			name:     "unknown-different-value",
			value:    uuidtypes.NewUUIDUnknown(),
			other:    uuidtypes.NewUUIDValue(valueUUIDv5),
			expected: false,
		},

		// Value
		{
			name:     "value-nil",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    nil,
			expected: false,
		},
		{
			name:     "value-null",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    uuidtypes.NewUUIDNull(),
			expected: false,
		},
		{
			name:     "value-unknown",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    uuidtypes.NewUUIDUnknown(),
			expected: false,
		},
		{
			name:     "value-value",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    uuidtypes.NewUUIDValue(valueUUIDv4),
			expected: true,
		},
		{
			name:     "value-different-value",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    uuidtypes.NewUUIDValue(valueUUIDv5),
			expected: false,
		},
		{
			name:     "not-uuidtypes.UUIDValue",
			value:    uuidtypes.NewUUIDValue(valueUUIDv4),
			other:    types.StringValue(valueUUIDv4),
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

func TestUUIDValue_Type(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    uuidtypes.UUIDValue
		expected attr.Type
	}{
		{
			name:     "always",
			value:    uuidtypes.NewUUIDNull(),
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
