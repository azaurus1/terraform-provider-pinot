data "pinot_segments" "example" {
  table_name = "<EXAMPLE>"
}

output "github_segments" {
  value = data.pinot_segments.example
}
