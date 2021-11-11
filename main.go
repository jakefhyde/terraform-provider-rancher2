package main

import (
	"context"
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rancher/terraform-provider-rancher2/rancher2"
)

func main() {
	plugin.Debug(context.Background(), "rancher/rancher2", &plugin.ServeOpts{
		ProviderFunc: rancher2.Provider})

	//plugin.Serve(&plugin.ServeOpts{
	//	ProviderFunc: rancher2.Provider})
}
