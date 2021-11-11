package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2SecretImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	namespaceID, projectID, resourceID := splitRegistryID(d.Id())

	d.SetId(resourceID)
	d.Set("project_id", projectID)
	d.Set("namespace_id", namespaceID)

	err := resourceRancher2SecretReadImpl(ctx, d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
