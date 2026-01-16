// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GoogleWorkspaceProvider satisfies various provider interfaces.
var _ provider.Provider = &GoogleWorkspaceProvider{}
var _ provider.ProviderWithFunctions = &GoogleWorkspaceProvider{}
var _ provider.ProviderWithEphemeralResources = &GoogleWorkspaceProvider{}
var _ provider.ProviderWithActions = &GoogleWorkspaceProvider{}

// GoogleWorkspaceProvider defines the provider implementation.
type GoogleWorkspaceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GoogleWorkspaceProviderModel describes the provider data model.
type GoogleWorkspaceProviderModel struct {
	Credentials types.String `tfsdk:"credentials"`
}

func (p *GoogleWorkspaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *GoogleWorkspaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Path to credentials file.",
				Optional:            false,
			},
		},
	}
}

func (p *GoogleWorkspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GoogleWorkspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *GoogleWorkspaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *GoogleWorkspaceProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *GoogleWorkspaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *GoogleWorkspaceProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *GoogleWorkspaceProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GoogleWorkspaceProvider{
			version: version,
		}
	}
}
