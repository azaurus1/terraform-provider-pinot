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

  segments_config = {
    for key, value in local.config_raw["segmentsConfig"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  tenants = {
    for key, value in local.config_raw["tenants"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  table_index_config = {
    for key, value in local.config_raw["tableIndexConfig"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  ingestion_config = {
    for key, value in local.config_raw["ingestionConfig"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  stream_ingestion_config = {
    for key, value in local.ingestion_config["stream_ingestion_config"] :
    join("_", [for keyName in regexall("[A-Z]?[a-z]+", key) : lower(keyName)]) => value
  }

  kafka_config = {
    "stream.kafka.broker.list": local.kafka_broker,
    "stream.kafka.zk.broker.url": local.kafka_zk
  }

  parsed_stream_ingestion_config = {
    column_major_segment_builder_enabled = true
    stream_config_maps = [for value in local.stream_ingestion_config["stream_config_maps"] :merge(value, local.kafka_config)]
  }


}

#resource "pinot_schema" "offline_table_schema" {
#  schema_name = "offline_ethereum_mainnet_block_headers"
#  schema      = file("offline_table_schema_example.json")
#
#}
#
#resource "pinot_table" "offline_table" {
#  table_name = "offline_ethereum_mainnet_block_headers"
#  table_type = "OFFLINE"
#  table      = file("offline_table_example.json")
#  depends_on = [pinot_schema.offline_table_schema]
#}

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
    optimize_dictionary = true
  })

  ingestion_config = merge(local.ingestion_config, {
    segment_time_check_value = true
    continue_on_error        = true
    row_time_value_check     = true
    stream_ingestion_config = local.parsed_stream_ingestion_config
  })

  depends_on = [pinot_schema.realtime_table_schema]
}

# data "pinot_tables" "edu" {
# }
#
# output "edu_tables" {
#   value = local.parsed_stream_ingestion_config
#   sensitive = true
# }
#
