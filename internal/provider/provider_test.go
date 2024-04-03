// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "pinot" {
	controller_url = "http://localhost:9000"
	# auth_type      = "Bearer" // Bearer will use bearer token, anything else will use default which is basic
	auth_token = "YWRtaW46dmVyeXNlY3JldA"
 }
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"pinot": providerserver.NewProtocol6WithError(New("test")()),
	}
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
//var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//	"scaffolding": providerserver.NewProtocol6WithError(New("test")()),
//}
//
//func testAccPreCheck(t *testing.T) {
//	// You can add code here to run prior to any test case execution, for example assertions
//	// about the appropriate environment variables being set are common to see in a pre-check
//	// function.
//}
