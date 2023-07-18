# terraform-plugin-framework-type-uuid

[![godoc](https://pkg.go.dev/badge/github.com/matthewhartstonge/terraform-plugin-framework-type-uuid)](https://pkg.go.dev/github.com/matthewhartstonge/terraform-plugin-framework-type-uuid)

UUID type and value implementation for the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).
Provides validation via Hashicorp's UUID library for Go based on [RFC 4122](https://www.rfc-editor.org/rfc/rfc4122.html) (note: not intended to be spec compliant)

## Getting Started

### Schema

The Terraform Plugin Framework schema types accept a `CustomType` field. The `uuidtypes.UUID` custom type can be injected into any current `schema.StringAttribute`.

In the following example, the ID field is set to the custom UUID Type.

```go
func (r *ExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType: uuidtypes.UUIDType{},
				Required:   true,
				// ...
            },
        },
    }
}
```

### Schema Data Model

Replace usage of `types.String` in schema data models with `uuidtype.UUID`.

Given the previous schema data model:

```go
type ThingResourceModel struct {
    // ...
    ID types.String `tfsdk:"id"`
}
```

The updated schema data model will look like:

```go
type ThingResourceModel struct {
    // ...
    ID uuidtypes.UUID `tfsdk:"id"`
}
```

### Accessing Values

Similar to other value types, use the `IsNull()` and `IsUnknown()` methods to 
check whether the value is null or unknown. Use the `ValueString()` method to extract
a known `uuid` value.

### Writing Values

Create a `uuidtypes.UUID` by calling one of these functions:

- `NewUUIDNull() UUID`: creates a `null` value.
- `NewUUIDUnknown() UUID`: creates an unknown value.
- `NewUUIDValue(string) UUID`: creates a known value using the given `string`.
- `NewUUIDPointerValue(string) UUID`: creates a known value using the given `*string`.

This type implements validation which is called and handled by Terraform. 

### Adding the Dependency

All functionality is located in the `github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtypes` 
package. Add this as an `import` as required to your relevant Go files.

Run the following Go commands to fetch the latest version and ensure all module files are up-to-date.

```shell
go get github.com/matthewhartstonge/terraform-plugin-framework-type-uuid@latest
go mod tidy
```
