package rancher2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2AppV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, name := splitID(d.Id())
	d.Set("cluster_id", clusterID)
	d.Set("name", name)

	err := resourceRancher2AppV2ReadImpl(ctx, d, meta)
	if err != nil || d.Id() == "" { //todo(jhyde): real error
		return []*schema.ResourceData{}, fmt.Errorf("bad")
	}

	return []*schema.ResourceData{d}, nil
}
