/*
 * Copyright (c) 2023 Matthew Hartstonge <matt@mykro.co.nz>
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
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure Implementation matches the expected interfaces.
var (
	_ attr.Type                    = UUIDType{}
	_ basetypes.StringTypable      = UUIDType{}
	_ tftypes.AttributePathStepper = UUIDType{}
	_ xattr.TypeWithValidate       = UUIDType{}
)

type UUIDType struct {
	basetypes.StringType
}

// Equal returns true if the two values are equal.
func (u UUIDType) Equal(o attr.Type) bool {
	other, ok := o.(UUIDType)
	if !ok {
		return false
	}

	return u.StringType.Equal(other.StringType)
}

// String returns a human-friendly version of the Type.
func (u UUIDType) String() string {
	return "uuidtypes.UUIDType"
}

// Validate ensures the value is a valid UUID.
func (u UUIDType) Validate(_ context.Context, value tftypes.Value, schemaPath path.Path) diag.Diagnostics {
	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	var diags diag.Diagnostics

	var valueString string
	if err := value.As(&valueString); err != nil {
		diags.AddAttributeError(
			schemaPath,
			"Invalid UUID Terraform Value",
			"An unexpected error occurred while attempting to read a UUID string from the Terraform value. "+
				"Please contact the provider developers with the following:\n\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	if _, err := uuid.ParseUUID(valueString); err != nil {
		diags.AddAttributeError(
			schemaPath,
			"Invalid UUID String Value",
			"An unexpected error occurred attempting to parse a string value that was expected to be a valid UUID format. "+
				"The expected UUID format is 00000000-0000-0000-0000-00000000. "+
				"For example, a Version 4 UUID is of the form 7b16fd41-cc23-4ef7-8aa9-c598350ccd18.\n\n"+
				fmt.Sprintf("Provided Value: %q\n", valueString)+
				fmt.Sprintf("Parse Error: %s", err.Error()),
		)

		return diags
	}

	return diags
}

// ValueFromString converts a string value to a StringValuable.
func (u UUIDType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := UUIDValue{
		StringValue: in,
	}

	// TODO: not sure if should validate the UUID here given diags are returned...?

	return value, nil
}

// ValueFromTerraform returns a UUIDValue value given a tftypes.Value.
func (u UUIDType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := u.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := u.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ValueType returns attr.Value type returned by ValueFromTerraform.
func (u UUIDType) ValueType(context.Context) attr.Value {
	return UUIDValue{}
}
