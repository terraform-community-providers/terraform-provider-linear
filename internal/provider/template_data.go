package provider

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// templateDataValue picks what to store for a template's data. Linear returns
// the JSON with its own key order and spacing, which never matches the string
// Terraform's jsonencode() builds, so every plan offered to rewrite unchanged
// templates. If the prior value means the same JSON, keep it as written. With
// no prior value (import), store the form jsonencode() produces: sorted keys,
// no spacing, and <, > and & escaped.
func templateDataValue(prior types.String, api string) types.String {
	if !prior.IsNull() && !prior.IsUnknown() && jsonEqual(prior.ValueString(), api) {
		return prior
	}

	if canonical, ok := canonicalJSON(api); ok {
		return types.StringValue(canonical)
	}

	return types.StringValue(api)
}

func decodeJSON(s string) (interface{}, bool) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(s)))
	decoder.UseNumber()

	var value interface{}

	if err := decoder.Decode(&value); err != nil {
		return nil, false
	}

	if decoder.More() {
		return nil, false
	}

	return value, true
}

func jsonEqual(a, b string) bool {
	left, ok := decodeJSON(a)

	if !ok {
		return false
	}

	right, ok := decodeJSON(b)

	if !ok {
		return false
	}

	return reflect.DeepEqual(left, right)
}

func canonicalJSON(s string) (string, bool) {
	value, ok := decodeJSON(s)

	if !ok {
		return "", false
	}

	out, err := json.Marshal(value)

	if err != nil {
		return "", false
	}

	return string(out), true
}
