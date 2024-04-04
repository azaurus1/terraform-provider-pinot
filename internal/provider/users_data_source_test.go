package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	pinot_testContainer "github.com/azaurus1/pinot-testContainer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
	context, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	pinot, err := pinot_testContainer.RunPinotContainer(context)
	if err != nil {
		t.Fatalf("Failed to run Pinot container: %v", err)
	}

	// fmt.Println(pinot.URI)

	providerConfig := fmt.Sprintf(`
provider "pinot" {
	controller_url = "http://%s"
	auth_token = "YWRtaW46dmVyeXNlY3JldA"
}
`, pinot.URI)

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
