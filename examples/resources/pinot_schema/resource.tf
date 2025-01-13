resource "pinot_schema" "block_schema" {
  schema_name = "ethereum_block_headers"
  date_time_field_specs = [{
    data_type          = "LONG",
    name               = "block_timestamp",
    format             = "1:MILLISECONDS:EPOCH",
    granularity        = "1:MILLISECONDS",
    transform_function = "ago('PT3H')"
  }]
  dimension_field_specs = [{
    name      = "block_number",
    data_type = "INT",
    not_null  = true
    },
    {
      name               = "block_hash",
      data_type          = "STRING",
      not_null           = true,
      transform_function = "jsonPathString(block, '$.block_hash')"
  }]
  metric_field_specs = [{
    name      = "block_difficulty",
    data_type = "INT",
    not_null  = true
  }]
}