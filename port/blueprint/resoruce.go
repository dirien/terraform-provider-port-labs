package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BlueprintResource{}
var _ resource.ResourceWithImportState = &BlueprintResource{}

func NewBlueprintResource() resource.Resource {
	return &BlueprintResource{}
}

type BlueprintResource struct {
	portClient *cli.PortClient
}

func (r *BlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *BlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

// func getArrayDefaultAttribute(arrayType interface{}) schema.Attribute {
// 	switch arrayType {
// 	case "string":
// 		return schema.ListAttribute{
// 			MarkdownDescription: "The default of the array property",
// 			Optional:            true,
// 			ElementType:         types.StringType,
// 		}
// 	case "boolean":
// 		return schema.ListAttribute{
// 			MarkdownDescription: "The default of the array property",
// 			Optional:            true,
// 			ElementType:         types.BoolType,
// 		}
// 	}
// 	return schema.ListAttribute{
// 		MarkdownDescription: "The default of the array property",
// 		Optional:            true,
// 		ElementType:         types.StringType,
// 	}
// }

func (r *BlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BlueprintModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read data from the API
	b, statusCode, err := r.portClient.ReadBlueprint(ctx, data.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to read blueprint: %s", err))
		return
	}

	writeBlueprintFieldsToResource(data, b)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintFieldsToResource(bm *BlueprintModel, b *cli.Blueprint) {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.Identifier = types.StringValue(b.Identifier)
	bm.Title = types.StringValue(b.Title)
	bm.Icon = types.StringValue(b.Icon)
	if !bm.Description.IsNull() {
		bm.Description = types.StringValue(b.Description)
	}
	bm.CreatedAt = types.StringValue(b.CreatedAt.String())
	bm.CreatedBy = types.StringValue(b.CreatedBy)
	bm.UpdatedAt = types.StringValue(b.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(b.UpdatedBy)
	if b.ChangelogDestination != nil {
		bm.ChangelogDestination = &ChangelogDestinationModel{
			Type:  types.StringValue(b.ChangelogDestination.Type),
			Url:   types.StringValue(b.ChangelogDestination.Url),
			Agent: types.BoolValue(b.ChangelogDestination.Agent),
		}
	}

	properties := &PropertiesModel{}

	addPropertiesToResource(b, bm, properties)

	bm.Properties = properties

}

func addPropertiesToResource(b *cli.Blueprint, bm *BlueprintModel, properties *PropertiesModel) {
	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]StringPropModel)
			}

			stringProp := &StringPropModel{}

			if v.Enum != nil && !bm.Properties.StringProp[k].Enum.IsNull() {
				attrs := make([]attr.Value, 0, len(v.Enum))
				for _, value := range v.Enum {
					attrs = append(attrs, basetypes.NewStringValue(value.(string)))
				}

				stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
			} else {
				stringProp.Enum = types.ListNull(types.StringType)
			}

			if v.Spec != "" && !bm.Properties.StringProp[k].Spec.IsNull() {
				stringProp.Spec = types.StringValue(v.Spec)
			}

			if v.MinLength != 0 && !bm.Properties.StringProp[k].MinLength.IsNull() {
				stringProp.MinLength = types.Int64Value(int64(v.MinLength))
			}

			if v.MaxLength != 0 && !bm.Properties.StringProp[k].MaxLength.IsNull() {
				stringProp.MaxLength = types.Int64Value(int64(v.MaxLength))
			}

			if v.Pattern != "" && !bm.Properties.StringProp[k].Pattern.IsNull() {
				stringProp.Pattern = types.StringValue(v.Pattern)
			}

			if v.SpecAuthentication != nil && bm.Properties.StringProp[k].SpecAuthentication != nil {
				stringProp.SpecAuthentication = &SpecAuthenticationModel{
					AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
					TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
					ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
				}
			}

			setCommonProperties(v, bm.Properties.StringProp[k], stringProp)

			properties.StringProp[k] = *stringProp

		case "number":
			if properties.NumberProp == nil {
				properties.NumberProp = make(map[string]NumberPropModel)
			}

			numberProp := &NumberPropModel{}

			if v.Minimum != 0 && !bm.Properties.NumberProp[k].Minimum.IsNull() {
				numberProp.Minimum = types.Float64Value(v.Minimum)
			}

			if v.Maximum != 0 && !bm.Properties.NumberProp[k].Maximum.IsNull() {
				numberProp.Maximum = types.Float64Value(v.Maximum)
			}

			if v.Enum != nil && !bm.Properties.NumberProp[k].Enum.IsNull() {
				attrs := make([]attr.Value, 0, len(v.Enum))
				for _, value := range v.Enum {
					attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
				}

				numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)
			}

			setCommonProperties(v, bm.Properties.NumberProp[k], numberProp)

			properties.NumberProp[k] = *numberProp

		case "array":
			if properties.ArrayProp == nil {
				properties.ArrayProp = make(map[string]ArrayPropModel)
			}

			arrayProp := &ArrayPropModel{}

			if v.MinItems != 0 && !bm.Properties.ArrayProp[k].MinItems.IsNull() {
				arrayProp.MinItems = types.Int64Value(int64(v.MinItems))
			}
			if v.MaxItems != 0 && !bm.Properties.ArrayProp[k].MaxItems.IsNull() {
				arrayProp.MaxItems = types.Int64Value(int64(v.MaxItems))
			}

			if v.Items != nil {
				if v.Items["type"] != "" {
					switch v.Items["type"] {
					case "string":
						arrayProp.StringItems = &StringItems{}
						if v.Items["default"] != nil {
							stringArray := v.Items["default"].([]string)
							attrs := make([]attr.Value, 0, len(stringArray))
							for _, value := range stringArray {
								attrs = append(attrs, basetypes.NewStringValue(value))
							}
							arrayProp.StringItems.Default, _ = types.ListValue(types.StringType, attrs)
						} else {
							arrayProp.StringItems.Default = types.ListNull(types.StringType)
						}
						if v.Items["format"] != "" {
							arrayProp.StringItems.Format = types.StringValue(v.Items["format"].(string))
						}
					case "number":
						arrayProp.NumberItems = &NumberItems{}
						if v.Items["default"] != nil {
							numberArray := v.Items["default"].([]float64)
							attrs := make([]attr.Value, 0, len(numberArray))
							for _, value := range numberArray {
								attrs = append(attrs, basetypes.NewFloat64Value(value))
							}
							arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
						}

					case "boolean":
						arrayProp.BooleanItems = &BooleanItems{}
						if v.Items["default"] != nil {
							booleanArray := v.Items["default"].([]bool)
							attrs := make([]attr.Value, 0, len(booleanArray))
							for _, value := range booleanArray {
								attrs = append(attrs, basetypes.NewBoolValue(value))
							}
							arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
						}
					}
				}
			}

			setCommonProperties(v, bm.Properties.ArrayProp[k], arrayProp)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, bm.Properties.BooleanProp[k], booleanProp)

			properties.BooleanProp[k] = *booleanProp

		case "object":
			if properties.ObjectProp == nil {
				properties.ObjectProp = make(map[string]ObjectPropModel)
			}

			objectProp := &ObjectPropModel{}

			if v.Spec != "" && !bm.Properties.ObjectProp[k].Spec.IsNull() {
				objectProp.Spec = types.StringValue(v.Spec)
			}

			setCommonProperties(v, bm.Properties.ObjectProp[k], objectProp)

			properties.ObjectProp[k] = *objectProp

		}

	}
}

func setCommonProperties(v cli.BlueprintProperty, bm interface{}, prop interface{}) {
	properties := []string{"description", "icon", "default", "title"}
	for _, property := range properties {
		switch property {
		case "description":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Description == "" && bmString.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Description == "" && bmNumber.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Description == "" && bmBoolean.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Description == "" && bmArray.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Description == "" && bmObject.Description.IsNull() {
					continue
				}
				p.Description = types.StringValue(v.Description)
			}
		case "icon":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Icon == "" && bmString.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Icon == "" && bmNumber.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Icon == "" && bmBoolean.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Icon == "" && bmArray.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Icon == "" && bmObject.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			}
		case "title":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Title == "" && bmString.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Title == "" && bmNumber.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Title == "" && bmBoolean.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Title == "" && bmArray.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Title == "" && bmObject.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)

			}

		case "default":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Default == nil && bmString.Default.IsNull() {
					continue
				}
				p.Default = types.StringValue(v.Default.(string))
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Default == nil && bmNumber.Default.IsNull() {
					continue
				}
				p.Default = types.Float64Value(v.Default.(float64))
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Default == nil && bmBoolean.Default.IsNull() {
					continue
				}
				p.Default = types.BoolValue(v.Default.(bool))
				// case *cli.ObjectPropModel:
				// 	bmObject := bm.(cli.ObjectPropModel)
				// 	if bmObject.Default != nil {
				// 		p.Default = v.Default.(map[string]interface{})
				// 	}
				// }

			}
		}
	}
}

func defaultResourceToBody(value string, propFields *cli.BlueprintProperty) error {
	switch propFields.Type {
	case "string":
		propFields.Default = value
	case "number":
		defaultNum, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return err
		}
		propFields.Default = defaultNum

	case "boolean":

		defaultBool, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		propFields.Default = defaultBool

	case "object":
		defaultObj := make(map[string]interface{})
		err := json.Unmarshal([]byte(value), &defaultObj)
		if err != nil {
			return err
		}
		propFields.Default = defaultObj

	}
	return nil
}

func (r *BlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := blueprintResourceToBody(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}
	fmt.Printf("Creating Blueprint %+v\n", b)
	bp, err := r.portClient.CreateBlueprint(ctx, b)
	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToResource(data, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintComputedFieldsToResource(bm *BlueprintModel, bp *cli.Blueprint) {
	bm.Identifier = types.StringValue(bp.Identifier)
	bm.CreatedAt = types.StringValue(bp.CreatedAt.String())
	bm.CreatedBy = types.StringValue(bp.CreatedBy)
	bm.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(bp.UpdatedBy)
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	b, err := blueprintResourceToBody(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError("failed to transform blueprint", err.Error())
		return
	}

	var bp *cli.Blueprint

	if data.Identifier.IsNull() {
		bp, err = r.portClient.CreateBlueprint(ctx, b)
	} else {
		bp, err = r.portClient.UpdateBlueprint(ctx, b, data.Identifier.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToResource(data, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *BlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *BlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Identifier.IsNull() {
		resp.Diagnostics.AddError("failed to extract blueprint identifier", "identifier is required")
		return
	}

	err := r.portClient.DeleteBlueprint(ctx, data.Identifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete blueprint", err.Error())
		return
	}
}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}

func stringPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.StringProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "string",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueString()
			}

			if !prop.Format.IsNull() {
				property.Format = prop.Format.ValueString()
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.MinLength.IsNull() {
				property.MinLength = int(prop.MinLength.ValueInt64())
			}

			if !prop.MaxLength.IsNull() {
				property.MaxLength = int(prop.MaxLength.ValueInt64())
			}

			if !prop.Spec.IsNull() {
				property.Spec = prop.Spec.ValueString()
			}

			if prop.SpecAuthentication != nil {
				property.SpecAuthentication = &cli.SpecAuthentication{
					AuthorizationUrl: prop.SpecAuthentication.AuthorizationUrl.ValueString(),
					TokenUrl:         prop.SpecAuthentication.TokenUrl.ValueString(),
					ClientId:         prop.SpecAuthentication.ClientId.ValueString(),
				}
			}

			if !prop.Pattern.IsNull() {
				property.Pattern = prop.Pattern.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Enum.IsNull() {
				property.Enum = []interface{}{}
				for _, e := range prop.Enum.Elements() {
					v, _ := e.ToTerraformValue(ctx)
					var keyValue string
					v.As(&keyValue)
					property.Enum = append(property.Enum, keyValue)
				}
			}
			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func numberPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "number",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Minimum.IsNull() {
				property.Minimum = prop.Minimum.ValueFloat64()
			}

			if !prop.Maximum.IsNull() {
				property.Maximum = prop.Maximum.ValueFloat64()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Enum.IsNull() {
				property.Enum = []interface{}{}
				for _, e := range prop.Enum.Elements() {
					v, _ := e.ToTerraformValue(ctx)
					var keyValue float64
					v.As(&keyValue)
					property.Enum = append(property.Enum, keyValue)
				}
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func booleanPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "boolean",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func objectPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.ObjectProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "object",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Spec.IsNull() {
				property.Spec = prop.Spec.ValueString()
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func arrayPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.ArrayProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "array",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}
			if !prop.MinItems.IsNull() {
				property.MinItems = int(prop.MinItems.ValueInt64())
			}

			if !prop.MaxItems.IsNull() {
				property.MaxItems = int(prop.MaxItems.ValueInt64())
			}

			if prop.StringItems != nil {
				items := map[string]interface{}{}
				items["type"] = "string"
				if !prop.StringItems.Format.IsNull() {
					items["format"] = prop.StringItems.Format.ValueString()
				}
				if !prop.StringItems.Default.IsNull() {
					items["default"] = prop.StringItems.Default
				}
				property.Items = items
			}

			if prop.NumberItems != nil {
				items := map[string]interface{}{}
				items["type"] = "number"
				if !prop.NumberItems.Default.IsNull() {
					items["default"] = prop.NumberItems.Default
				}
				property.Items = items
			}

			if prop.BooleanItems != nil {
				items := map[string]interface{}{}
				items["type"] = "boolean"
				if !prop.BooleanItems.Default.IsNull() {
					items["default"] = prop.BooleanItems.Default
				}
				property.Items = items
			}

			if prop.ObjectItems != nil {
				items := map[string]interface{}{}
				items["type"] = "object"
				if !prop.ObjectItems.Default.IsNull() {
					items["default"] = prop.ObjectItems.Default
				}
				property.Items = items
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func blueprintResourceToBody(ctx context.Context, d *BlueprintModel) (*cli.Blueprint, error) {
	b := &cli.Blueprint{}
	b.Identifier = d.Identifier.ValueString()

	b.Title = d.Title.ValueString()
	b.Icon = d.Icon.ValueString()
	b.Description = d.Description.ValueString()
	props := map[string]cli.BlueprintProperty{}
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}
	relations := map[string]cli.Relation{}

	if d.ChangelogDestination != nil {
		b.ChangelogDestination = &cli.ChangelogDestination{}
		b.ChangelogDestination.Type = d.ChangelogDestination.Type.ValueString()
		b.ChangelogDestination.Url = d.ChangelogDestination.Url.ValueString()
		b.ChangelogDestination.Agent = d.ChangelogDestination.Agent.ValueBool()
	} else {
		b.ChangelogDestination = nil
	}

	required := []string{}

	if d.Properties != nil {
		if d.Properties.StringProp != nil {
			stringPropResourceToBody(ctx, d, props, required)
		}
		if d.Properties.ArrayProp != nil {
			arrayPropResourceToBody(d, props, required)
		}
		if d.Properties.NumberProp != nil {
			numberPropResourceToBody(ctx, d, props, required)
		}
		if d.Properties.BooleanProp != nil {
			booleanPropResourceToBody(d, props, required)
		}

		if d.Properties.ObjectProp != nil {
			objectPropResourceToBody(d, props, required)
		}

	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
	return b, nil
}
