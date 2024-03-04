terraform {
  required_providers {
    pinot = {
      source = "hashicorp.com/edu/pinot"
    }
  }
}


locals {
  schema_content = file("block_schema.json")
}

provider "pinot" {
  controller_url     = "http://localhost:9000"
}

resource "pinot_schema" "block_schema" {
  schema_name = local.schema_content.schemaName
  schema = local.schema_content
}