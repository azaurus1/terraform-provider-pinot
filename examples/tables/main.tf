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

locals {

  kafka_broker = "kafka:9092"
  kafka_zk     = "kafka:2181"
  config_raw   = jsondecode(file("realtime_example.json"))

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

  filter_config = try({
    for key, value in local.ingestion_config["filter_config"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }, null)

  stream_ingestion_config = {
    for key, value in local.ingestion_config["stream_ingestion_config"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  transform_configs = try([
    for value in local.ingestion_config["transform_configs"] :
    { for key, inner_value in value : join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => inner_value }
  ], [])

  task_config = try({
    for key, value in local.config_raw["task"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }, null)

  kafka_overrides = {}

  parsed_stream_ingestion_config = {
    column_major_segment_builder_enabled = true
    stream_config_maps = [
      for value in local.stream_ingestion_config["stream_config_maps"] : merge(value, local.kafka_overrides)
    ]
  }

  schema = {
    for key, value in jsondecode(file("realtime_example_schema.json")) :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  dimension_field_specs = [
    for field in local.schema["dimension_field_specs"] : {
      for key, value in field :
      join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
    }
  ]

  metric_field_specs = try([
    for field in local.schema["metric_field_specs"] : {
      for key, value in field :
      join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
    }
  ], [])

  date_time_field_specs = [
    for field in local.schema["date_time_field_specs"] : {
      for key, value in field :
      join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
    }
  ]

}

resource "pinot_schema" "realtime_table_schema" {
  schema_name                       = local.schema["schema_name"]
  enable_column_based_null_handling = try(local.schema["enable_column_based_null_handling"], false)
  primary_key_columns               = try(local.schema["primary_key_columns"], null)
  dimension_field_specs             = local.dimension_field_specs
  metric_field_specs                = local.metric_field_specs
  date_time_field_specs             = local.date_time_field_specs
}

resource "pinot_table" "realtime_table" {

  table_name = "monster_transformed"
  table_type = "REALTIME"
  table      = file("realtime_example.json")

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
    filter_config            = local.filter_config
  })


  metadata = local.metadata

  is_dim_table = local.config_raw["isDimTable"]

  task = local.task_config

  depends_on = [pinot_schema.realtime_table_schema]
}

output "tasks_thing" {
  value = local.task_config
}