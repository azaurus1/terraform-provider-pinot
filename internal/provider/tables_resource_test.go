package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	pinot_testContainer "github.com/azaurus1/pinot-testContainer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTablesResource(t *testing.T) {

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
locals {

	kafka_broker = "kafka:9092"
	kafka_zk     = "kafka:2181"
	config_raw   = jsondecode(file("/internal/provider/testdata/realtime_table_example.json"))
	
	## Convert the keys to snake_case
	segments_config = {
		for key, value in local.config_raw["segmentsConfig"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	tenants = {
		for key, value in local.config_raw["tenants"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	metadata = {
		for key, value in try(local.config_raw["metadata"], null) :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	routing = {
		for key, value in local.config_raw["routing"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	upsert_config = {
		for key, value in local.config_raw["upsertConfig"] :
		join("_", [for keyName in regexall("(?:[A-Z]+[a-z]*)|(?:[a-z]+)", key) : lower(keyName)]) => value
	}
	
	table_index_config = {
		for key, value in local.config_raw["tableIndexConfig"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	segment_partition_config = try({
		for key, value in local.table_index_config["segment_partition_config"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}, null)
	
	ingestion_config = {
		for key, value in local.config_raw["ingestionConfig"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	stream_ingestion_config = {
		for key, value in local.ingestion_config["stream_ingestion_config"] :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	transform_configs = [
		for value in local.ingestion_config["transform_configs"] :
		{ for key, inner_value in value : join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => inner_value }
	]
	
	kafka_overrides = {
		"stream.kafka.broker.list" : sensitive(local.kafka_broker),
		"stream.kafka.zk.broker.url" : sensitive(local.kafka_zk),
		"stream.kafka.topic.name" : "ethereum_mainnet_block_headers"
	}
	
	parsed_stream_ingestion_config = {
		column_major_segment_builder_enabled = true
		stream_config_maps = [
		for value in local.stream_ingestion_config["stream_config_maps"] : merge(value, local.kafka_overrides)
		]
	}
	
	schema = {
		for key, value in jsondecode(file("/internal/provider/testdata/realtime_table_schema_example.json")) :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
	}
	
	dimension_field_specs = [
		for field in local.schema["dimension_field_specs"] : {
		for key, value in field :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
		}
	]
	
	metric_field_specs = [
		for field in local.schema["metric_field_specs"] : {
		for key, value in field :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
		}
	]
	
	date_time_field_specs = [
		for field in local.schema["date_time_field_specs"] : {
		for key, value in field :
		join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
		}
	]
	
	}
	
	resource "pinot_schema" "realtime_table_schema" {
	schema_name                       = local.schema["schema_name"]
	enable_column_based_null_handling = local.schema["enable_column_based_null_handling"]
	primary_key_columns               = try(local.schema["primary_key_columns"], null)
	dimension_field_specs             = local.dimension_field_specs
	metric_field_specs                = local.metric_field_specs
	date_time_field_specs             = local.date_time_field_specs
	}
	
	resource "pinot_table" "realtime_table" {
	
	table_name = "realtime_ethereum_mainnet_block_headers_REALTIME"
	table_type = "REALTIME"
	table      = file("/internal/provider/testdata/realtime_table_example.json")
	
	segments_config = merge(local.segments_config, {
		replication = "1"
	})
	
	routing       = local.routing
	upsert_config = local.upsert_config
	
	tenants = merge(local.tenants, {
		broker = "DefaultTenant"
		server = "DefaultTenant"
	})
	
	table_index_config = merge(local.table_index_config, {
		optimize_dictionary      = false
		segment_partition_config = local.segment_partition_config
	})
	
	ingestion_config = merge(local.ingestion_config, {
		segment_time_check_value = true
		continue_on_error        = true
		row_time_value_check     = true
		stream_ingestion_config  = local.parsed_stream_ingestion_config
		transform_configs        = local.transform_configs
	})
	
	metadata = local.metadata
	
	is_dim_table = local.config_raw["isDimTable"]
	
	depends_on = [pinot_schema.realtime_table_schema]
	}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "table_name", "realtime_ethereum_mainnet_block_headers_REALTIME"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "table_type", "REALTIME"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "segments_config.replication", "1"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "routing.tenant_name", "DefaultTenant"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "upsert_config.upsert_mode", "FULL"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "table_index_config.optimize_dictionary", "false"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "ingestion_config.segment_time_check_value", "true"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "metadata.custom_property", "custom_value"),
					resource.TestCheckResourceAttr("pinot_table.realtime_table", "is_dim_table", "false"),
				),
			},
			// ImportState Testing - This is a special case where we need to import the state of the resource - Not Implemented Yet
			// {
			// ResourceName: "pinot_user.test",
			// ImportState: true,
			// ImportStateVerify: true,
			// },
			// Update and Read testing
			// 			{
			// 				Config: providerConfig + `
			// resource "pinot_user" "test" {
			// 	username  = "user"
			// 	password  = "password"
			// 	component = "BROKER"
			// 	role      = "ADMIN"

			// 	lifecycle {
			// 		ignore_changes = [password]
			// 		}
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					resource.TestCheckResourceAttr("pinot_user.test", "username", "user"),
			// 					// resource.TestCheckResourceAttr("pinot_user.test", "password", "password"), // This is ignored because it returns the salted and hashed password AGAIN
			// 					resource.TestCheckResourceAttr("pinot_user.test", "component", "BROKER"),
			// 					resource.TestCheckResourceAttr("pinot_user.test", "role", "ADMIN"),
			// 				),
			// 			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
