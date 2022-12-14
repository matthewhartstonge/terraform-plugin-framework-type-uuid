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
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure Implementation matches the expected interfaces.
var (
	_ attr.Type                    = UUIDType{}
	_ tftypes.AttributePathStepper = UUIDType{}
	_ xattr.TypeWithValidate       = UUIDType{}
)

type UUIDType struct{}

// ApplyTerraform5AttributePathStep always returns an error as this type cannot
// be walked any further.
func (u UUIDType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep to %T to %s", step, u.String())
}

// Equal returns true if the incoming Type is equal to the UUID type.
func (u UUIDType) Equal(other attr.Type) bool {
	_, ok := other.(UUIDType)

	return ok
}

// String returns a human-friendly version of the UUID Type.
func (u UUIDType) String() string {
	return "uuidtypes.UUIDType"
}

// TerraformType returns tftypes.String.
func (u UUIDType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// Validate returns any warnings or errors that occur while attempting to parse
// a UUID value.
func (u UUIDType) Validate(_ context.Context, value tftypes.Value, schemaPath path.Path) diag.Diagnostics {
	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	var str string
	err := value.As(&str)
	if err != nil {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				schemaPath,
				"Invalid UUID Terraform Value",
				"An unexpected error occurred while attempting to read a UUID string from the Terraform value. "+
					"Please contact the provider developers with the following:\n\n"+
					"Error: "+err.Error(),
			),
		}
	}

	_, diags := UUIDFromString(str, schemaPath)

	return diags
}

// ValueFromTerraform returns a UUID value given a tftypes.Value.
func (u UUIDType) ValueFromTerraform(_ context.Context, value tftypes.Value) (attr.Value, error) {
	if value.IsNull() {
		return UUIDNull(), nil
	}

	if !value.IsKnown() {
		return UUIDUnknown(), nil
	}

	var str string
	if err := value.As(&str); err != nil {
		return UUIDUnknown(), err
	}

	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return UUIDUnknown(), err
	}

	return UUIDFromGoogleUUID(parsedUUID), nil
}

// ValueType returns attr.Value type returned by ValueFromTerraform.
func (u UUIDType) ValueType(context.Context) attr.Value {
	return UUIDNull()
}
