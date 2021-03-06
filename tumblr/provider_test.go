package tumblr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"tumblr": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	for _, param := range []string{"consumer_key", "consumer_secret", "user_token", "user_token_secret"} {
		value, _ := testAccProvider.Schema[param].DefaultFunc()
		if value == "" {
			t.Fatalf("A Tumblr %s was not found. Read tumblr docs for how to configure this.", param)
		}
	}
}
