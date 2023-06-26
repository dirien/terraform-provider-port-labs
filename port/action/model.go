package action

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookMethodModel struct {
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}

type GithubMethodModel struct {
	Org                  types.String `tfsdk:"org"`
	Repo                 types.String `tfsdk:"repo"`
	Workflow             types.String `tfsdk:"workflow"`
	OmitPayload          types.Bool   `tfsdk:"omit_payload"`
	OmitUserInputs       types.Bool   `tfsdk:"omit_user_inputs"`
	ReportWorkflowStatus types.Bool   `tfsdk:"report_workflow_status"`
	Branch               types.String `tfsdk:"branch"`
}

type AzureMethodModel struct {
	Org     types.String `tfsdk:"agent"`
	Webhook types.String `tfsdk:"webhook"`
}

type ActionModel struct {
	ID               types.String        `tfsdk:"id"`
	Identifier       types.String        `tfsdk:"identifier"`
	Blueprint        types.String        `tfsdk:"blueprint"`
	Title            types.String        `tfsdk:"title"`
	Icon             types.String        `tfsdk:"icon"`
	Description      types.String        `tfsdk:"description"`
	RequiredApproval types.Bool          `tfsdk:"required_approval"`
	KafkaMethod      types.Map           `tfsdk:"kafka_method"`
	WebhookMethod    *WebhookMethodModel `tfsdk:"webhook_method"`
	GithubMethod     *GithubMethodModel  `tfsdk:"github_method"`
	AzureMethod      *AzureMethodModel   `tfsdk:"azure_method"`
}
