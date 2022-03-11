package rancher2

import (
	"reflect"
	"testing"

	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	testClusterV2RKEConfigMachinePoolMachineConfigConf      *corev1.ObjectReference
	testClusterV2RKEConfigMachinePoolMachineConfigInterface []interface{}
	testClusterV2RKEConfigMachinePoolRollingUpdateConf      *provisionv1.RKEMachinePoolRollingUpdate
	testClusterV2RKEConfigMachinePoolRollingUpdateInterface []interface{}
	testClusterV2RKEConfigMachinePoolsConf                  []provisionv1.RKEMachinePool
	testClusterV2RKEConfigMachinePoolsInterface             []interface{}
)

func init() {
	testClusterV2RKEConfigMachinePoolMachineConfigConf = &corev1.ObjectReference{
		Kind: "kind",
		Name: "name",
	}
	testClusterV2RKEConfigMachinePoolMachineConfigInterface = []interface{}{
		map[string]interface{}{
			"kind": "kind",
			"name": "name",
		},
	}
	maxSurge := intstr.FromString("max_surge")
	maxUnavailable := intstr.FromString("max_unavailable")
	testClusterV2RKEConfigMachinePoolRollingUpdateConf = &provisionv1.RKEMachinePoolRollingUpdate{
		MaxSurge:       &maxSurge,
		MaxUnavailable: &maxUnavailable,
	}

	testClusterV2RKEConfigMachinePoolRollingUpdateInterface = []interface{}{
		map[string]interface{}{
			"max_surge":       "max_surge",
			"max_unavailable": "max_unavailable",
		},
	}
	quantity := int32(10)
	testClusterV2RKEConfigMachinePoolsConf = []provisionv1.RKEMachinePool{
		{
			Name:             "test",
			DisplayName:      "test",
			NodeConfig:       testClusterV2RKEConfigMachinePoolMachineConfigConf,
			ControlPlaneRole: true,
			EtcdRole:         true,
			MachineDeploymentAnnotations: map[string]string{
				"anno_one": "one",
				"anno_two": "two",
			},
			MachineDeploymentLabels: map[string]string{
				"label_one": "one",
				"label_two": "two",
			},
			Quantity:      &quantity,
			Paused:        true,
			RollingUpdateContext: testClusterV2RKEConfigMachinePoolRollingUpdateConf,
			WorkerRole:    true,
		},
	}
	testClusterV2RKEConfigMachinePoolsConf[0].CloudCredentialSecretName = "cloud_credential_secret_name"
	testClusterV2RKEConfigMachinePoolsConf[0].Taints = []corev1.Taint{
		{
			Key:    "key",
			Value:  "value",
			Effect: "recipient",
		},
	}
	testClusterV2RKEConfigMachinePoolsInterface = []interface{}{
		map[string]interface{}{
			"name":                         "test",
			"cloud_credential_secret_name": "cloud_credential_secret_name",
			"machine_config":               testClusterV2RKEConfigMachinePoolMachineConfigInterface,
			"control_plane_role":           true,
			"etcd_role":                    true,
			"annotations": map[string]interface{}{
				"anno_one": "one",
				"anno_two": "two",
			},
			"labels": map[string]interface{}{
				"label_one": "one",
				"label_two": "two",
			},
			"quantity":       10,
			"paused":         true,
			"rolling_update": testClusterV2RKEConfigMachinePoolRollingUpdateInterface,
			"taints": []interface{}{
				map[string]interface{}{
					"key":    "key",
					"value":  "value",
					"effect": "recipient",
				},
			},
			"worker_role": true,
		},
	}
}

func TestFlattenClusterV2RKEConfigMachinePoolMachineConfig(t *testing.T) {

	cases := []struct {
		Input          *corev1.ObjectReference
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigMachinePoolMachineConfigConf,
			testClusterV2RKEConfigMachinePoolMachineConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigMachinePoolMachineConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigMachinePoolRollingUpdate(t *testing.T) {

	cases := []struct {
		Input          *provisionv1.RKEMachinePoolRollingUpdate
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigMachinePoolRollingUpdateConf,
			testClusterV2RKEConfigMachinePoolRollingUpdateInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigMachinePoolRollingUpdate(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigMachinePools(t *testing.T) {

	cases := []struct {
		Input          []provisionv1.RKEMachinePool
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigMachinePoolsConf,
			testClusterV2RKEConfigMachinePoolsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigMachinePools(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigMachinePoolMachineConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *corev1.ObjectReference
	}{
		{
			testClusterV2RKEConfigMachinePoolMachineConfigInterface,
			testClusterV2RKEConfigMachinePoolMachineConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigMachinePoolMachineConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigMachinePoolRollingUpdate(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *provisionv1.RKEMachinePoolRollingUpdate
	}{
		{
			testClusterV2RKEConfigMachinePoolRollingUpdateInterface,
			testClusterV2RKEConfigMachinePoolRollingUpdateConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigMachinePoolRollingUpdate(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigMachinePools(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []provisionv1.RKEMachinePool
	}{
		{
			testClusterV2RKEConfigMachinePoolsInterface,
			testClusterV2RKEConfigMachinePoolsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigMachinePools(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
