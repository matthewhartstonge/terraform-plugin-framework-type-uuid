# terraform-plugin-framework-type-uuid

[![godoc](https://pkg.go.dev/badge/github.com/matthewhartstonge/terraform-plugin-framework-type-uuid)](https://pkg.go.dev/github.com/matthewhartstonge/terraform-plugin-framework-type-uuid)

UUID string type and value implementation for the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).
Provides validation via Google's UUID library for Go based on [RFC 4122](https://www.rfc-editor.org/rfc/rfc4122.html)
and DCE 1.1: Authentication and Security Services.

## Getting Started

### Schema

Replace usages of `types.StringType` with `uuidtype.Type{}`.

Given the previous schema attribute:

```go
tfsdk.Attribute{
	Required: true
	Type:     types.StringType 
	// Potential prior validators...
}
```

The updated schema will look like:

```go
tfsdk.Attribute{
	Required: true
	Type:     uuidtype.Type{}
}
```

### Schema Data Model

Replace usage of `string`, `*string`, or `types.String` in schema data models 
with `uuidtype.Value`.

Given the previous schema data model:

```go
type ThingResourceModel struct {
    // ...
    Example types.String `tfsdk:"example"`
}
```

The updated schema data model will look like:

```go
type ThingResourceModel struct {
    // ...
    Example uuidtype.Value `tfsdk:"example"`
}
```

### Accessing Values

Similar to other value types, use the `IsNull()` and `IsUnknown()` methods to 
check whether the value is null or unknown. Use the `String()` method to extract
a known `uuid` value.

### Writing Values

Create a `uuidtype.Value` by calling one of these functions:

- `NullValue() Value`: creates a `null` value.
- `UnknownValue() Value`: creates an unknown value.
- `StringValue(string, path.Path) (Value, diag.Diagnostics)`: creates a known 
   value using the given `string` or returns validation errors if `string` is 
   not in the expected UUID format.
- `MustValue(string) Value` creates a known value using the given string, but 
   will panic if it's unparseable as a UUID.

### Adding the Dependency

All functionality is located in the `github.com/matthewhartstonge/terraform-plugin-framework-type-uuid/uuidtype` 
package. Add this as an `import` as required to your relevant Go files.

Run the following Go commands to fetch the latest version and ensure all module files are up-to-date.

```shell
go get github.com/matthewhartstonge/terraform-plugin-framework-type-uuidtype@latest
go mod tidy
```
