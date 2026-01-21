// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

// GroupDataSource defines the data source implementation.
type GroupDataSource struct {
	client *http.Client

	adminService *admin.Service
}

// GroupDataSourceModel describes the data source data model.
type GroupDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

func (d *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Group configurable attribute",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Group configurable attribute",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Group configurable attribute",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Group identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
	srv, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	d.adminService = srv

}

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, err := d.adminService.Groups.Get(data.Name.ValueString()).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read group '%s', got error: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("read group: %+v", g))

	data.Id = types.StringValue(g.Id)
	data.Email = types.StringValue(g.Email)
	data.Description = types.StringValue(g.Description)
	data.Name = types.StringValue(g.Name)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
