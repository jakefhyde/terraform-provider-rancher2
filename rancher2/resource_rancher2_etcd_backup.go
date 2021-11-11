package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2EtcdBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2EtcdBackupCreate,
		ReadContext:   resourceRancher2EtcdBackupRead,
		UpdateContext: resourceRancher2EtcdBackupUpdate,
		DeleteContext: resourceRancher2EtcdBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2EtcdBackupImport,
		},

		Schema: etcdBackupFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2EtcdBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	etcdBackup, err := expandEtcdBackup(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Etcd Backup")

	err = meta.(*Config).ClusterExist(etcdBackup.ClusterID)
	if err != nil {
		return err
	}

	newEtcdBackup, err := client.EtcdBackup.Create(etcdBackup)
	if err != nil {
		return err
	}

	d.SetId(newEtcdBackup.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active", "activating"},
		Refresh:    etcdBackupStateRefreshFunc(client, newEtcdBackup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for etcd backup (%s) to be created: %s", newEtcdBackup.ID, waitErr)
	}

	return resourceRancher2EtcdBackupRead(ctx, d, meta)
}

func resourceRancher2EtcdBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Etcd Backup ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	etcdBackup, err := client.EtcdBackup.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Etcd Backup ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenEtcdBackup(d, etcdBackup)
}

func resourceRancher2EtcdBackupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Etcd Backup ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	etcdBackup, err := client.EtcdBackup.ByID(d.Id())
	if err != nil {
		return err
	}

	backupConfig, err := expandClusterRKEConfigServicesEtcdBackupConfig(d.Get("backup_config").([]interface{}))
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"backup_config": backupConfig,
		"filename":      d.Get("filename").(string),
		"manual":        d.Get("manual").(bool),
		"annotations":   toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":        toMapString(d.Get("labels").(map[string]interface{})),
	}

	newEtcdBackup, err := client.EtcdBackup.Update(etcdBackup, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "activating"},
		Target:     []string{"active", "activating"},
		Refresh:    etcdBackupStateRefreshFunc(client, newEtcdBackup.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for etcd backup (%s) to be updated: %s", newEtcdBackup.ID, waitErr)
	}

	return resourceRancher2EtcdBackupRead(ctx, d, meta)
}

func resourceRancher2EtcdBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Etcd Backup ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	etcdBackup, err := client.EtcdBackup.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Etcd Backup ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.EtcdBackup.Delete(etcdBackup)
	if err != nil {
		return fmt.Errorf("Error removing Etcd Backup: %s", err)
	}

	log.Printf("[DEBUG] Waiting for etcd backup (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    etcdBackupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for etcd backup (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// etcdBackupStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher EtcdBackup.
func etcdBackupStateRefreshFunc(client *managementClient.Client, nodePoolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.EtcdBackup.ByID(nodePoolID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
