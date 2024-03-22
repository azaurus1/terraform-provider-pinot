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

resource "pinot_schema" "block_schema" {
  schema_name = "ethereum_block_headers"
  date_time_field_specs = [{
    data_type   = "LONG",
    name        = "block_timestamp",
    format      = "1:MILLISECONDS:EPOCH",
    granularity = "1:MILLISECONDS",
  }]
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
