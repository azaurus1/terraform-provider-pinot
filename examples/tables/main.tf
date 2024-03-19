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
  config_raw = jsondecode(file("realtime_table_example.json"))

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

  table_index_config = {
    for key, value in local.config_raw["tableIndexConfig"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  segment_partition_config = try({
    for key, value in local.table_index_config["segment_partition_config"]:
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

  kafka_overrides_secrets = sensitive({
    "stream.kafka.broker.list": local.kafka_broker,
    "stream.kafka.zk.broker.url": local.kafka_zk
  })

  kafka_overrides = {
    "stream.kafka.broker.list": local.kafka_broker,
    "stream.kafka.zk.broker.url": local.kafka_zk
    "stream.kafka.topic.name": "ethereum_mainnet_block_headers"
  }

  parsed_stream_ingestion_config = {
    column_major_segment_builder_enabled = true
    stream_config_maps = [
      for value in local.stream_ingestion_config["stream_config_maps"] :merge(value, local.kafka_overrides_secrets, local.kafka_overrides)
    ]
  }


}

resource "pinot_schema" "realtime_table_schema" {
  schema_name = "realtime_ethereum_mainnet_block_headers"
  schema      = file("realtime_table_schema_example.json")
}

resource "pinot_table" "realtime_table" {

  table_name = "realtime_ethereum_mainnet_block_headers_REALTIME"
  table_type = "REALTIME"
  table      = file("realtime_table_example.json")

  segments_config = merge(local.segments_config, {
    replication = "1"
  })

  tenants = merge(local.tenants, {
    broker = "DefaultTenant"
    server = "DefaultTenant"
  })

  table_index_config = merge(local.table_index_config, {
    optimize_dictionary = false
    segment_partition_config = local.segment_partition_config
  })

  ingestion_config = merge(local.ingestion_config, {
    segment_time_check_value = true
    continue_on_error        = true
    row_time_value_check     = true
    stream_ingestion_config = local.parsed_stream_ingestion_config
  })

  metadata = local.metadata

  is_dim_table = local.config_raw["isDimTable"]

  depends_on = [pinot_schema.realtime_table_schema]
}

output "d" {
  value = local.metadata
}
