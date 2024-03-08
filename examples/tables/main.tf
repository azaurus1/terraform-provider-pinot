terraform {
  required_providers {
    pinot = {
      source = "hashicorp.com/edu/pinot"
    }
  }
}

provider "pinot" {
  controller_url = "http://localhost:9000"
  auth_token     = "YWRtaW46dmVyeXNlY3JldA"
}

resource "pinot_schema" "offline_table_schema" {
  schema_name = "offline_ethereum_mainnet_block_headers"
  schema      = file("offline_table_schema_example.json")

}

resource "pinot_table" "offline_table" {
  table_name = "offline_ethereum_mainnet_block_headers"
  table_type = "OFFLINE"
  table      = file("offline_table_example.json")
  depends_on = [pinot_schema.offline_table_schema]
}

resource "pinot_schema" "realtime_table_schema" {
  schema_name = "realtime_ethereum_mainnet_block_headers"
  schema      = file("realtime_table_schema_example.json")

}

resource "pinot_table" "realtime_table" {
  table_name = "realtime_ethereum_mainnet_block_headers"
  table_type = "REALTIME"
  table      = file("realtime_table_example.json")
  depends_on = [pinot_schema.realtime_table_schema]
}

# data "pinot_tables" "edu" {
# }

# output "edu_tables" {
#   value = data.pinot_tables.edu

# }