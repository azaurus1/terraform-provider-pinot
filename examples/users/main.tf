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

data "pinot_users" "edu" {}

output "edu_users" {
  value = data.pinot_users.edu
}