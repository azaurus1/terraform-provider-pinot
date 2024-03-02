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

resource "pinot_user" "edu" {
    username = "edu"
    password = "password"
    component = "BROKER"
    role = "ADMIN"
}

data "pinot_users" "edu" {}

output "edu_users" {
  value = data.pinot_users.edu
}