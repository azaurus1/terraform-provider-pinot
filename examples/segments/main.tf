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


data "pinot_segments" "github" {
  table_name = "githubComplexTypeEvents"
}

data "pinot_segments" "airline" {
  table_name = "airlineStats"
}

output "github_segments" {
  value = data.pinot_segments.github
}

output "airline_segments" {
  value = data.pinot_segments.airline
}