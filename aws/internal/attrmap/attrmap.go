package attrmap

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AttributeMap represents a map of Terraform resource attribute name to AWS API attribute name.
// Useful for SQS Queue or SNS Topic attribute handling.
type AttributeMap map[string]string

// ApiAttributesToResourceData sets Terraform ResourceData from a map of AWS API attributes.
func (m AttributeMap) ApiAttributesToResourceData(apiAttributes map[string]string, d *schema.ResourceData) error {
	for tfAttributeName, apiAttributeName := range m {
		if v, ok := apiAttributes[apiAttributeName]; ok {
			if err := d.Set(tfAttributeName, v); err != nil {
				return fmt.Errorf("error setting %s: %w", tfAttributeName, err)
			}
		}
	}

	return nil
}

// ResourceDataToApiAttributesCreate returns a map of AWS API attributes from Terraform ResourceData.
// The API attributes map is suitable for resource create.
func (m AttributeMap) ResourceDataToApiAttributesCreate(d *schema.ResourceData) (map[string]string, error) {
	apiAttributes := map[string]string{}

	for tfAttributeName, apiAttributeName := range m {
		if v, ok := d.GetOk(tfAttributeName); ok {
			var apiAttributeValue string

			switch v := v.(type) {
			case int:
				apiAttributeValue = strconv.Itoa(v)
			case bool:
				apiAttributeValue = strconv.FormatBool(v)
			case string:
				apiAttributeValue = v
			default:
				return nil, fmt.Errorf("attribute %s is of unsupported type: %T", tfAttributeName, v)
			}

			apiAttributes[apiAttributeName] = apiAttributeValue
		}
	}

	return apiAttributes, nil
}

// ResourceDataToApiAttributesUpdate returns a map of AWS API attributes from Terraform ResourceData.
// The API attributes map is suitable for resource update.
func (m AttributeMap) ResourceDataToApiAttributesUpdate(d *schema.ResourceData) (map[string]string, error) {
	apiAttributes := map[string]string{}

	for tfAttributeName, apiAttributeName := range m {
		if d.HasChange(tfAttributeName) {
			var apiAttributeValue string

			switch v := d.Get(tfAttributeName).(type) {
			case int:
				apiAttributeValue = strconv.Itoa(v)
			case bool:
				apiAttributeValue = strconv.FormatBool(v)
			case string:
				apiAttributeValue = v
			default:
				return nil, fmt.Errorf("attribute %s is of unsupported type: %T", tfAttributeName, v)
			}

			apiAttributes[apiAttributeName] = apiAttributeValue
		}
	}

	return apiAttributes, nil
}
