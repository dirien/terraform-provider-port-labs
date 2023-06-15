package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SpreadMaps(target map[string]schema.Attribute, source map[string]schema.Attribute) {
	for key, value := range source {
		target[key] = value
	}
}

func MetadataProperties() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"title": schema.StringAttribute{
			MarkdownDescription: "The display name of the blueprint",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the blueprint",
			Optional:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "The required of the number property",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the blueprint",
			Optional:            true,
		}}

}

func StringPropertySchema() schema.Attribute {
	stringPropertySchema := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			MarkdownDescription: "The default of the string property",
			Optional:            true,
		},
		"format": schema.StringAttribute{
			MarkdownDescription: "The format of the string property",
			Optional:            true,
		},
		"min_length": schema.Int64Attribute{
			MarkdownDescription: "The min length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "The max length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"pattern": schema.StringAttribute{
			MarkdownDescription: "The pattern of the string property",
			Optional:            true,
		},
		"spec": schema.StringAttribute{
			MarkdownDescription: "The spec of the string property",
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("open-api", "async-api", "embedded-url")},
		},
		"spec_authentication": schema.SingleNestedAttribute{
			MarkdownDescription: "The spec authentication of the string property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"client_id": schema.StringAttribute{
					MarkdownDescription: "The clientId of the spec authentication",
					Required:            true,
				},
				"token_url": schema.StringAttribute{
					MarkdownDescription: "The tokenUrl of the spec authentication",
					Required:            true,
				},
				"authorization_url": schema.StringAttribute{
					MarkdownDescription: "The authorizationUrl of the spec authentication",
					Required:            true,
				},
			},
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the string property",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
	}

	SpreadMaps(stringPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The string property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: stringPropertySchema,
		},
	}
}

func NumberPropertySchema() schema.Attribute {
	numberPropertySchema := map[string]schema.Attribute{
		"default": schema.Float64Attribute{
			MarkdownDescription: "The default of the number property",
			Optional:            true,
		},
		"maximum": schema.Float64Attribute{
			MarkdownDescription: "The min of the number property",
			Optional:            true,
		},
		"minimum": schema.Float64Attribute{
			MarkdownDescription: "The max of the number property",
			Optional:            true,
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the number property",
			Optional:            true,
			ElementType:         types.Float64Type,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
	}

	SpreadMaps(numberPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The number property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: numberPropertySchema,
		},
	}
}

func BooleanPropertySchema() schema.Attribute {
	booleanPropertySchema := map[string]schema.Attribute{}

	SpreadMaps(booleanPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The boolean property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: booleanPropertySchema,
		},
	}
}

func ArrayPropertySchema() schema.MapNestedAttribute {
	arrayPropertySchema := map[string]schema.Attribute{
		"min_items": schema.Int64Attribute{
			MarkdownDescription: "The min items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_items": schema.Int64Attribute{
			MarkdownDescription: "The max items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"string_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of the items",
					Optional:            true,
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		},
		"number_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
			},
		},
		"boolean_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.BoolType,
				},
			},
		},
		"object_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.MapType{ElemType: types.StringType},
				},
			},
		},
	}

	SpreadMaps(arrayPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The array property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: arrayPropertySchema,
		},
	}
}

func ObjectPropertySchema() schema.MapNestedAttribute {

	objectPropertySchema := map[string]schema.Attribute{
		"spec": schema.StringAttribute{
			MarkdownDescription: "The spec of the object property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("async-api", "open-api"),
			},
		},
		"default": schema.MapAttribute{
			Optional:            true,
			MarkdownDescription: "The default of the object property",
			ElementType:         types.StringType,
		},
	}

	SpreadMaps(objectPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The object property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: objectPropertySchema,
		},
	}
}

func (r *BlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	blueprintSchema := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display name of the blueprint",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the blueprint",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the blueprint",
			Optional:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the blueprint",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the blueprint",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the blueprint",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the blueprint",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"changelog_destination": schema.SingleNestedAttribute{
			MarkdownDescription: "The changelog destination of the blueprint",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					MarkdownDescription: "The type of the changelog destination",
					Required:            true,
					// Validators:          []validator.String{stringvalidator.OneOf("WEBHOOK", "KAFKA")},
				},
				"url": schema.StringAttribute{
					MarkdownDescription: "The url of the changelog destination",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.Expressions{
							path.MatchRelative(),
						}...),
					},
				},
				"agent": schema.BoolAttribute{
					MarkdownDescription: "The agent of the changelog destination",
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
			},
		},
		"properties": schema.SingleNestedAttribute{
			MarkdownDescription: "The properties of the blueprint",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_prop":  StringPropertySchema(),
				"number_prop":  NumberPropertySchema(),
				"boolean_prop": BooleanPropertySchema(),
				"array_prop":   ArrayPropertySchema(),
				"object_prop":  ObjectPropertySchema(),
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group resource",
		Attributes:          blueprintSchema,
	}
}