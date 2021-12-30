package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jhp0204/Provider_test_2/Ongoing/scp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
	  ProviderFunc: scp.Provider
}),
}
