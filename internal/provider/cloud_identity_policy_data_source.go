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
	"google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/option"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CloudIdentityPolicyDataSource{}

func NewCloudIdentityPolicyDataSource() datasource.DataSource {
	return &CloudIdentityPolicyDataSource{}
}

// CloudIdentityPolicyDataSource defines the data source implementation.
type CloudIdentityPolicyDataSource struct {
	client *http.Client

	cloudidentityService *cloudidentity.Service
}

// CloudIdentityPolicyDataSourceModel describes the data source data model.
type CloudIdentityPolicyDataSourceModel struct {
	Name     types.String  `tfsdk:"name"`
	Customer types.String  `tfsdk:"customer"`
	Type     types.String  `tfsdk:"type"`
	Query    *QueryModel   `tfsdk:"query"`
	Setting  *SettingModel `tfsdk:"setting"`
	Id       types.String  `tfsdk:"id"`
}

// Nested Model for "query".
type QueryModel struct {
	Group   types.String `tfsdk:"group"`
	OrgUnit types.String `tfsdk:"org_unit"`
	Query   types.String `tfsdk:"query"`
}

// Nested Model for "setting".
type SettingModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

func (d *CloudIdentityPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_identity_policy"
}

func (d *CloudIdentityPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Cloud Identity Policy data source",

		Attributes: map[string]schema.Attribute{
			"customer": schema.StringAttribute{
				MarkdownDescription: `Customer that the Policy belongs to. The value 
				is in the format 'customers/{customerId}'. The 'customerId must begin 
				with "C" To find your customer ID in Admin Console see
				https://support.google.com/a/answer/10070793`,
				Required: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: `Identifier. The resource name 
				(https://cloud.google.com/apis/design/resource_names) 
				of the Policy. Format: policies/{policy}.`,
				Required: true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: `The type of the policy.
	 			Possible values:
	 			  "POLICY_TYPE_UNSPECIFIED" - Unspecified policy type.
	 			  "SYSTEM" - Policy type denoting the system-configured policies.
	 			  "ADMIN" - Policy type denoting the admin-configurable policies.`,
				Computed: true,
			},
			"query": schema.SingleNestedAttribute{
				MarkdownDescription: "The Policy Query",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"group": schema.StringAttribute{
						MarkdownDescription: `This field is only set if there is a single 
						value for group that satisfies all clauses of the  query. 
						If no group applies, this will be the empty string.`,
						Computed: true,
					},
					"org_unit": schema.StringAttribute{
						MarkdownDescription: `The OrgUnit the query applies to. This field 
						is only set if there is a single value for org_unit that satisfies 
						all clauses of the query.`,
						Computed: true,
					},
					"query": schema.StringAttribute{
						MarkdownDescription: `The CEL query that defines which entities the Policy 
						applies to (ex. a User entity). For details about CEL see 
						https://opensource.google.com/projects/cel. The OrgUnits the Policy applies 
						to are represented by a clause like so: 
							entity.org_units.exists(org_unit, org_unit.org_unit_id == orgUnitId('{orgUnitId}')) 
						The Group the Policy applies to are represented by a clause like so: 
							entity.groups.exists(group, group.group_id == groupId('{groupId}')) 
						The Licenses the Policy applies to are represented by a clause like so: 
							entity.licenses.exists(license, license in ['/product/{productId}/sku/{skuId}']) 
						The above clauses can be present in any combination, and used in conjunction 
						with the &&, || and ! operators. The org_unit and group fields below are helper 
						fields that contain the corresponding value(s) as the query to make the query easier to use.`,
						Computed: true,
					},
				},
			},
			"setting": schema.SingleNestedAttribute{
				MarkdownDescription: "The Policy Query",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: `The type of the Setting.`,
						Computed:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: `The value of the Setting.`,
						Computed:            true,
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource ID",
				Computed:            true,
			},
		},
	}
}

func (d *CloudIdentityPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	srv, err := cloudidentity.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	d.cloudidentityService = srv

}

func (d *CloudIdentityPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudIdentityPolicyDataSourceModel

	// 1. Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyName := data.Name.ValueString()

	policy, err := d.cloudidentityService.Policies.Get(policyName).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read Cloud Identity Policy '%s': %s", policyName, err),
		)
		return
	}

	data.Id = types.StringValue(policy.Name)
	data.Name = types.StringValue(policy.Name)

	data.Query = nil
	if policy.PolicyQuery != nil {
		data.Query = &QueryModel{
			Group:   types.StringValue(policy.PolicyQuery.Group),
			OrgUnit: types.StringValue(policy.PolicyQuery.OrgUnit),
			Query:   types.StringValue(policy.PolicyQuery.Query), // The raw CEL string
		}
	}

	data.Setting = nil
	if policy.Setting != nil {
		data.Setting = &SettingModel{
			Type:  types.StringValue(policy.Setting.Type),
			Value: types.StringValue(string(policy.Setting.Value)), // Raw JSON value as string
		}
	}
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
