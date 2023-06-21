package entity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringPropModel struct {
	Value types.String `tfsdk:"value"`
}

type NumberPropModel struct {
	Value types.Float64 `tfsdk:"value"`
}

type EntityPropertiesModel struct {
	StringProp map[string]StringPropModel `tfsdk:"string_prop"`
	NumberProp map[string]NumberPropModel `tfsdk:"number_prop"`
}

type EntityModel struct {
	ID          types.String           `tfsdk:"id"`
	Identifier  types.String           `tfsdk:"identifier"`
	Blueprint   types.String           `tfsdk:"blueprint"`
	Title       types.String           `tfsdk:"title"`
	Icon        types.String           `tfsdk:"icon"`
	RunID       types.String           `tfsdk:"run_id"`
	Description types.String           `tfsdk:"description"`
	CreatedAt   types.String           `tfsdk:"created_at"`
	CreatedBy   types.String           `tfsdk:"created_by"`
	UpdatedAt   types.String           `tfsdk:"updated_at"`
	UpdatedBy   types.String           `tfsdk:"updated_by"`
	Properties  *EntityPropertiesModel `tfsdk:"properties"`
}