provider "pinot" {
  # example configuration here
  controller_url = "http://localhost:9000"  //required (can also be set via environment variable PINOT_CONTROLLER_URL)
  auth_token     = "YWRtaW46dmVyeXNlY3JldA" //optional (can also be set via environment variable PINOT_AUTH_TOKEN)
  auth_type      = "bearer"                 //optional (can also be set via environment variable PINOT_AUTH_TYPE)
}
