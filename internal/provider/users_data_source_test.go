package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "pinot_users" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pinot_users.test", "id", "placeholder"), // Test placeholder ID is created
					resource.TestCheckResourceAttr("data.pinot_users.test", "users.#", "0"),      // No users when running the docker compose
				),
			},
		},
	})

}
