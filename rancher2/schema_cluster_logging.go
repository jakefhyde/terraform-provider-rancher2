package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	clusterLoggingKinds = []string{loggingCustomTargetKind, loggingElasticsearchKind, loggingFluentdKind, loggingKafkaKind, loggingSplunkKind, loggingSyslogKind}
)

// Shemas

func clusterLoggingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(clusterLoggingKinds, true),
		},
		"custom_target_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"elasticsearch_config", "fluentd_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingCustomTargetConfigFields(),
			},
		},
		"enable_json_parsing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Optional enable json log parsing",
		},
		"elasticsearch_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"custom_target_config", "fluentd_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingElasticsearchConfigFields(),
			},
		},
		"fluentd_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"custom_target_config", "elasticsearch_config", "kafka_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingFluentdConfigFields(),
			},
		},
		"kafka_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"custom_target_config", "elasticsearch_config", "fluentd_config", "splunk_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingKafkaConfigFields(),
			},
		},
		"namespace_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"output_flush_interval": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  3,
		},
		"output_tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"splunk_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"custom_target_config", "elasticsearch_config", "fluentd_config", "kafka_config", "syslog_config"},
			Elem: &schema.Resource{
				Schema: loggingSplunkConfigFields(),
			},
		},
		"syslog_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"custom_target_config", "elasticsearch_config", "fluentd_config", "kafka_config", "splunk_config"},
			Elem: &schema.Resource{
				Schema: loggingSyslogConfigFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
