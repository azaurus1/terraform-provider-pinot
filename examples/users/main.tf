terraform {
  required_providers {
    pinot = {
      source = "hashicorp.com/edu/pinot"
    }
  }
}

provider "pinot" {
  controller_url = "http://localhost:9000"
  # auth_type      = "Bearer" // Bearer will use bearer token, anything else will use default which is basic
  auth_token = "YWRtaW46dmVyeXNlY3JldA"
}

resource "pinot_user" "test" {
  username  = "liam"
  password  = "password"
  component = "BROKER"
  role      = "USER"

  lifecycle {
    ignore_changes = [password]
  }
}

data "pinot_users" "edu" {}

output "edu_users" {
  value = data.pinot_users.edu
}