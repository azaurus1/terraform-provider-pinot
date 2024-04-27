data "pinot_users" "example" {
  # Specify the configuration for the data source here, if applicable
}

output "example_output" {
  value = data.pinot_users.example
}
