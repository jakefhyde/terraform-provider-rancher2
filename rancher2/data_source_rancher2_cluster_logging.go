package rancher2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2ClusterLogging() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2ClusterLoggingRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_target_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingCustomTargetConfigFields(),
				},
			},
			"enable_json_parsing": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Optional enable json log parsing",
			},
			"elasticsearch_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingElasticsearchConfigFields(),
				},
			},
			"fluentd_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingFluentdConfigFields(),
				},
			},
			"kafka_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingKafkaConfigFields(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_flush_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"output_tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"splunk_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSplunkConfigFields(),
				},
			},
			"syslog_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loggingSyslogConfigFields(),
				},
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

func dataSourceRancher2ClusterLoggingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID := d.Get("cluster_id").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
	}
	listOpts := NewListOpts(filters)

	clusterLoggings, err := client.ClusterLogging.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(clusterLoggings.Data)
	if count <= 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] cluster logging on cluster ID \"%s\" not found", clusterID))
	}
	if count > 1 {
		return diag.FromErr(fmt.Errorf("[ERROR] found %d cluster logging on cluster ID \"%s\"", count, clusterID))
	}

	return diag.FromErr(flattenClusterLogging(d, &clusterLoggings.Data[0]))
}
