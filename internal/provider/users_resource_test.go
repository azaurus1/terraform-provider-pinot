package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	pinot_testContainer "github.com/azaurus1/pinot-testContainer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersResource(t *testing.T) {

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
				Config: providerConfig + `
				resource "pinot_user" "test" {
					username  = "user"
					password  = "password"
					component = "BROKER"
					role      = "USER"
								
					lifecycle {
						ignore_changes = [password]
						}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pinot_user.test", "username", "user"),
					resource.TestCheckResourceAttr("pinot_user.test", "password", "password"),
					resource.TestCheckResourceAttr("pinot_user.test", "component", "BROKER"),
					resource.TestCheckResourceAttr("pinot_user.test", "role", "USER"),
				),
			},
			// ImportState Testing - This is a special case where we need to import the state of the resource - Not Implemented Yet
			// {
			// ResourceName: "pinot_user.test",
			// ImportState: true,
			// ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "pinot_user" "test" {
					username  = "user"
					password  = "password"
					component = "BROKER"
					role      = "ADMIN"

					lifecycle {
						ignore_changes = [password]
						}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pinot_user.test", "username", "user"),
					// resource.TestCheckResourceAttr("pinot_user.test", "password", "password"), // This is ignored because it returns the salted and hashed password AGAIN
					resource.TestCheckResourceAttr("pinot_user.test", "component", "BROKER"),
					resource.TestCheckResourceAttr("pinot_user.test", "role", "ADMIN"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
