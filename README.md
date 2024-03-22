# `terraform-provider-pinot`
A [Terraform](https://www.terraform.io/) provider for [Apache Pinot](https://pinot.apache.org/)

## Contents
- [Using the provider](#using-the-provider)
  - [Installation](#installation)
  - [Examples](#examples)
- [Requirements](#requirements)



## Using the provider

### Installation:
```hcl
terraform {
  required_providers {
    pinot = {
      source = "azaurus1/pinot"
      version = "0.2.0"
    }
  }
}

provider "pinot" {
  controller_url = "http://localhost:9000" //required (can also be set via environment variable PINOT_CONTROLLER_URL)
  auth_token     = "YWRtaW46dmVyeXNlY3JldA" //optional (can also be set via environment variable PINOT_AUTH_TOKEN) 
}

```
### Examples:
Example Schema:
```
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
```

Example Table:
```
resource "pinot_table" "realtime_table" {

  table_name = "realtime_ethereum_mainnet_block_headers_REALTIME"
  table_type = "REALTIME"
  table      = file("realtime_table_example.json")
```

Example User:
```
resource "pinot_user" "test" {
  username  = "user"
  password  = "password"
  component = "BROKER"
  role      = "USER"

  lifecycle {
    ignore_changes = [password]
  }
}
```

This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). 

## Requirements



- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

### To run this in your local environment be sure to:
1. run `go env GOBIN`
2. Create `~/.terraformrc` with the following:
```hcl
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/pinot" = [Location of your GOBIN]
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
``` 

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
