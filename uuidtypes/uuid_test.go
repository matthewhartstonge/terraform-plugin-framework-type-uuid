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
	"testing"

	// External Imports
	"github.com/google/go-cmp/cmp"

	// Internal Imports
	"github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtypes"
)

func TestNewUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected uuidtypes.UUIDValue
	}{
		{
			name:     "string-value-empty",
			value:    "",
			expected: uuidtypes.NewUUIDValue(""),
		},
		{
			name:     "string-value-invalid-length",
			value:    valueInvalidLength,
			expected: uuidtypes.NewUUIDValue(valueInvalidLength),
		},
		{
			name:     "string-value-invalid-format",
			value:    valueInvalid,
			expected: uuidtypes.NewUUIDValue(valueInvalid),
		},
	}

	for _, testcase := range tests {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			got := uuidtypes.NewUUIDValue(testcase.value)
			if diff := cmp.Diff(got, testcase.expected); diff != "" {
				t.Errorf("Equal()\ngot     : %v\nexpected: %v\ndiff    : %s\n",
					got,
					testcase.expected,
					cmp.Diff(got, testcase.expected),
				)
			}
		})
	}
}
