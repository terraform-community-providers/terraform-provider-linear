package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTemplateDataValue(t *testing.T) {
	api := `{"title":"Test <b> & co","priority":2,"labelIds":[]}`
	jsonencoded := `{"labelIds":[],"priority":2,"title":"Test \u003cb\u003e \u0026 co"}`

	cases := []struct {
		name  string
		prior types.String
		want  string
	}{
		{"import stores jsonencode form", types.StringNull(), jsonencoded},
		{"same JSON keeps prior text", types.StringValue("{ \"priority\": 2, \"labelIds\": [], \"title\": \"Test <b> & co\" }"), "{ \"priority\": 2, \"labelIds\": [], \"title\": \"Test <b> & co\" }"},
		{"changed JSON takes API value", types.StringValue(`{"title":"Old"}`), jsonencoded},
		{"unknown prior takes API value", types.StringUnknown(), jsonencoded},
	}

	for _, c := range cases {
		got := templateDataValue(c.prior, api).ValueString()

		if got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, got, c.want)
		}
	}

	if got := templateDataValue(types.StringNull(), "not json").ValueString(); got != "not json" {
		t.Errorf("invalid JSON should pass through, got %s", got)
	}
}
