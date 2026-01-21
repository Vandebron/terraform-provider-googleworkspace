// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

// GroupResource defines the resource implementation.
type GroupResource struct {
	client *http.Client

	adminService *admin.Service
}

// GroupResourceModel describes the resource data model.
type GroupResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

func (g *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (g *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Group name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Group name",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Group configurable attribute with default value",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Group identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (g *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	g.client = client
	srv, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	g.adminService = srv

}

func (g *GroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ng := &admin.Group{
		Email:       data.Email.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	res, err := g.adminService.Groups.Insert(ng).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Google Group",
			fmt.Sprintf("Could not create group %s: %v", data.Email.ValueString(), err),
		)
		return
	}

	data.Id = types.StringValue(res.Id)
	data.Email = types.StringValue(res.Email)
	data.Name = types.StringValue(res.Name)
	data.Description = types.StringValue(res.Description)

	tflog.Trace(ctx, "Created Google Group", map[string]interface{}{
		"id":    res.Id,
		"email": res.Email,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (g *GroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ng, err := g.adminService.Groups.Get(data.Name.ValueString()).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read group '%s', got error: %s", data.Name.ValueString(), err),
		)
		return
	}

	data.Id = types.StringValue(ng.Id)
	data.Email = types.StringValue(ng.Email)
	data.Description = types.StringValue(ng.Description)
	data.Name = types.StringValue(ng.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (g *GroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gu := &admin.Group{
		Email:       data.Email.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	res, err := g.adminService.Groups.Update(data.Id.ValueString(), gu).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Google Group",
			fmt.Sprintf("Could not update group ID %s: %v", data.Id.ValueString(), err),
		)
		return
	}

	data.Email = types.StringValue(res.Email)
	data.Name = types.StringValue(res.Name)
	data.Description = types.StringValue(res.Description)
	data.Id = types.StringValue(res.Id)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (g *GroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := g.adminService.Groups.Delete(data.Id.ValueString()).Context(ctx).Do()
	if err != nil {
		var googleErr *googleapi.Error
		if errors.As(err, &googleErr) && googleErr.Code == 404 {
			// Log this for debugging purposes, but do not return an error to Terraform.
			tflog.Warn(ctx, "Group already deleted in Google Workspace", map[string]interface{}{
				"id": data.Id.ValueString(),
			})
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting Google Group",
			fmt.Sprintf("Could not delete group ID %s: %v", data.Id.ValueString(), err),
		)
		return
	}
}

func (g *GroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
