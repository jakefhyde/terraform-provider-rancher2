package rancher2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2ClusterScan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2ClusterScanRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster ID to scan",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The cluster scan name",
			},
			"run_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan run type",
			},
			"scan_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterScanConfigFields(),
				},
				Description: "The cluster scan config",
			},
			"scan_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan type",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster scan status",
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterScanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
	}
	if len(name) > 0 {
		filters["name"] = name
	}
	listOpts := NewListOpts(filters)

	clusterScans, err := client.ClusterScan.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(clusterScans.Data)
	if count <= 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] cluster scan with cluster ID \"%s\" not found", clusterID))
	}
	if count > 1 {
		return diag.FromErr(fmt.Errorf("[ERROR] found %d cluster scan with cluster ID \"%s\"", count, clusterID))
	}

	d.SetId(clusterScans.Data[0].ID)

	return diag.FromErr(flattenClusterScan(d, &clusterScans.Data[0]))
}
