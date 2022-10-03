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
	_ attr.Type                    = Type{}
	_ tftypes.AttributePathStepper = Type{}
	_ xattr.TypeWithValidate       = Type{}
)

type Type struct{}

// ApplyTerraform5AttributePathStep always returns an error as this type cannot
// be walked any further.
func (t Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep to %T to %s", step, t.String())
}

// Equal returns true if the incoming Type is equal to the UUID type.
func (t Type) Equal(other attr.Type) bool {
	_, ok := other.(Type)

	return ok
}

// String returns a human-friendly version of the UUID Type.
func (t Type) String() string {
	return "uuidtype.Type"
}

// TerraformType returns tftypes.String.
func (t Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// Validate returns any warnings or errors that occur while attempting to parse
// a UUID value.
func (t Type) Validate(_ context.Context, value tftypes.Value, schemaPath path.Path) diag.Diagnostics {
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

	_, diags := StringValue(str, schemaPath)

	return diags
}

// ValueFromTerraform returns a UUID value given a tftypes.Value.
func (t Type) ValueFromTerraform(_ context.Context, value tftypes.Value) (attr.Value, error) {
	if value.IsNull() {
		return NullValue(), nil
	}

	if !value.IsKnown() {
		return UnknownValue(), nil
	}

	var str string
	if err := value.As(&str); err != nil {
		return UnknownValue(), err
	}

	if _, err := uuid.Parse(str); err != nil {
		return UnknownValue(), err
	}

	return Value{value: str}, nil
}
