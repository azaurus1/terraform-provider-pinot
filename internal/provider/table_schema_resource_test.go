package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	pinot_testContainer "github.com/azaurus1/pinot-testContainer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemasResource(t *testing.T) {

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
				resource "pinot_schema" "block_schema" {
					schema_name = "ethereum_block_headers"
					date_time_field_specs = [{
					  data_type   = "LONG",
					  name        = "block_timestamp",
					  not_null    = false,
					  format      = "1:MILLISECONDS:EPOCH",
					  granularity = "1:MILLISECONDS",
					}]
					enable_column_based_null_handling = false
					dimension_field_specs = [{
					  name      = "block_number",
					  data_type = "INT",
					  not_null  = true
					  },
					  {
						name      = "block_hash",
						data_type = "STRING",
						not_null  = true
					}]
					metric_field_specs = [{
					  name      = "block_difficulty",
					  data_type = "INT",
					  not_null  = true
					}]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "schema_name", "ethereum_block_headers"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.data_type", "LONG"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.name", "block_timestamp"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.format", "1:MILLISECONDS:EPOCH"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.granularity", "1:MILLISECONDS"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.name", "block_number"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.data_type", "INT"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.not_null", "true"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.name", "block_hash"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.data_type", "STRING"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.not_null", "true"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.name", "block_difficulty"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.data_type", "INT"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.not_null", "true"),
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
				resource "pinot_schema" "block_schema" {
					schema_name = "ethereum_block_headers"
					date_time_field_specs = [{
					  data_type   = "LONG",
					  name        = "block_timestamp",
					  not_null    = false,
					  format      = "1:MILLISECONDS:EPOCH",
					  granularity = "1:MILLISECONDS",
					}]
					enable_column_based_null_handling = false
					dimension_field_specs = [{
					  name      = "block_number",
					  data_type = "INT",
					  not_null  = true
					  },
					  {
						name      = "block_hash",
						data_type = "STRING",
						not_null  = true
					}]
					metric_field_specs = [{
					  name      = "block_difficulty",
					  data_type = "INT",
					  not_null  = true
					},
					{
					  name      = "block_gas_limit",
					  data_type = "INT",
					  not_null  = true
					}]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "schema_name", "ethereum_block_headers"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.data_type", "LONG"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.name", "block_timestamp"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.format", "1:MILLISECONDS:EPOCH"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "date_time_field_specs.0.granularity", "1:MILLISECONDS"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.name", "block_number"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.data_type", "INT"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.0.not_null", "true"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.name", "block_hash"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.data_type", "STRING"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "dimension_field_specs.1.not_null", "true"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.name", "block_difficulty"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.data_type", "INT"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.0.not_null", "true"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.1.name", "block_gas_limit"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.1.data_type", "INT"),
					resource.TestCheckResourceAttr("pinot_schema.block_schema", "metric_field_specs.1.not_null", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
