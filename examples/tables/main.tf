terraform {
  required_providers {
    pinot = {
      source = "hashicorp.com/edu/pinot"
    }
  }
}

provider "pinot" {
  controller_url     = "http://localhost:9000"
}

resource "pinot_table" "test" {
}

data "pinot_tables" "edu" {}

output "edu_tables" {
  value = data.pinot_tables.edu
}