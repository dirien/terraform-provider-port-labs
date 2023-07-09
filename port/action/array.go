package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func handleArrayItemsToBody(ctx context.Context, property *cli.BlueprintProperty, prop ArrayPropModel, required *[]string) error {
	if prop.StringItems != nil {
		items := map[string]interface{}{}
		items["type"] = "string"
		if !prop.StringItems.Format.IsNull() {
			items["format"] = prop.StringItems.Format.ValueString()
		}

		if !prop.StringItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "string")
			if err != nil {
				return err
			}

			property.Default = defaultList
		}
		property.Items = items
	}

	if prop.NumberItems != nil {
		items := map[string]interface{}{}
		items["type"] = "number"
		if !prop.NumberItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "float64")
			if err != nil {
				return err
			}

			items["default"] = defaultList
		}
		property.Items = items
	}

	if prop.BooleanItems != nil {
		items := map[string]interface{}{}
		items["type"] = "boolean"
		if !prop.BooleanItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "bool")
			if err != nil {
				return err
			}

			items["default"] = defaultList
		}
		property.Items = items
	}

	if prop.ObjectItems != nil {
		items := map[string]interface{}{}
		items["type"] = "object"
		if !prop.ObjectItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "object")
			if err != nil {
				return err
			}
			items["default"] = defaultList
		}
		property.Items = items
	}
	return nil
}

func arrayPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ArrayProps {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "array",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}
			if !prop.MinItems.IsNull() {
				minItems := int(prop.MinItems.ValueInt64())
				property.MinItems = &minItems
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = &maxItems
			}

			err := handleArrayItemsToBody(ctx, &property, prop, required)
			if err != nil {
				return err
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addArrayPropertiesToResource(v *cli.BlueprintProperty) *ArrayPropModel {
	arrayProp := &ArrayPropModel{
		MinItems: flex.GoInt64ToFramework(v.MinItems),
		MaxItems: flex.GoInt64ToFramework(v.MaxItems),
	}

	if v.Items != nil {
		if v.Items["type"] != "" {
			switch v.Items["type"] {
			case "string":
				arrayProp.StringItems = &StringItems{}
				if v.Default != nil {
					stringArray := make([]string, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						stringArray[i] = v.(string)
					}
					attrs := make([]attr.Value, 0, len(stringArray))
					for _, value := range stringArray {
						attrs = append(attrs, basetypes.NewStringValue(value))
					}
					arrayProp.StringItems.Default, _ = types.ListValue(types.StringType, attrs)
				} else {
					arrayProp.StringItems.Default = types.ListNull(types.StringType)
				}
				if value, ok := v.Items["format"]; ok && value != nil {
					arrayProp.StringItems.Format = types.StringValue(v.Items["format"].(string))
				}
			case "number":
				arrayProp.NumberItems = &NumberItems{}
				if v.Default != nil {
					numberArray := make([]float64, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(numberArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
					}
					arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
				}

			case "boolean":
				arrayProp.BooleanItems = &BooleanItems{}
				if v.Default != nil {
					booleanArray := make([]bool, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(booleanArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewBoolValue(value.(bool)))
					}
					arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
				}
			}
		}
	}

	return arrayProp
}
