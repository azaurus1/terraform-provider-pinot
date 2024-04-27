data "pinot_tenants" "example" {
  # Specify any required arguments here
}

output "tenant_names" {
  value = data.pinot_tenants.example.tenant_names
}
