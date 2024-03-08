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

resource "pinot_schema" "test_table_schema" {
  schema_name = "ethereum_mainnet_block_headers"
  schema      = file("table_schema_example.json")

}

resource "pinot_table" "test" {
  table_name = "ethereum_mainnet_block_headers"
  table      = file("table_example.json")

}

# data "pinot_tables" "edu" {
# }

# output "edu_tables" {
#   value = data.pinot_tables.edu

# }