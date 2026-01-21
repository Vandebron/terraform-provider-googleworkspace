// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2/google"

	admin "google.golang.org/api/admin/directory/v1"
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
	Credentials           types.String `tfsdk:"credentials"`
	ImpersonatedUserEmail types.String `tfsdk:"impersonated_user_email"`
}

func (p *GoogleWorkspaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googleworkspace"
	resp.Version = p.version
}

func (p *GoogleWorkspaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Path to Google Credentials JSON file (defaults to GOOGLE_CREDENTIALS)",
				Required:            true,
			},
			"impersonated_user_email": schema.StringAttribute{
				MarkdownDescription: "User to impersenate for domain-wide delegation (if applicable)",
				Required:            true,
			},
		},
	}
}

// Configure prepares a Google Workspace GRPC client for data sources and
// resources.
func (p *GoogleWorkspaceProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data GoogleWorkspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := os.ReadFile(data.Credentials.ValueString())
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, admin.AdminDirectoryGroupScope, admin.AdminDirectoryUserScope)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse service account JSON",
			"The provided credentials file is not valid JSON or missing fields: "+err.Error(),
		)
		return
	}

	// 3. CRITICAL: Set the Subject (Domain-Wide Delegation)
	// This explicitly tells Google: "I am this Service Account, but I want to act as THIS user."
	if data.ImpersonatedUserEmail.IsNull() || data.ImpersonatedUserEmail.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing Impersonated User Email",
			"When using Domain-Wide Delegation, you must provide the email of the admin user to impersonate.",
		)
		return
	}
	config.Subject = data.ImpersonatedUserEmail.ValueString()
	// 4. Create the Client
	// This client will now automatically refresh tokens acting as the 'Subject' user.
	client := config.Client(ctx)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *GoogleWorkspaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
	}
}

func (p *GoogleWorkspaceProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *GoogleWorkspaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGroupDataSource,
	}
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
