package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2AppImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, appID, err := splitAppID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(appID)
	d.Set("project_id", projectID)

	err = resourceRancher2AppReadImpl(ctx, d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
