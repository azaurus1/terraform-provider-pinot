resource "pinot_user" "test" {
  username  = "liam"
  password  = "password"
  component = "BROKER"
  role      = "USER"

  lifecycle {
    ignore_changes = [password]
  }
}
