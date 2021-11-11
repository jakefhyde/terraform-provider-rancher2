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

func resourceRancher2GlobalRoleBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2GlobalRoleBindingCreate,
		ReadContext:   resourceRancher2GlobalRoleBindingRead,
		UpdateContext: resourceRancher2GlobalRoleBindingUpdate,
		DeleteContext: resourceRancher2GlobalRoleBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2GlobalRoleBindingImport,
		},

		Schema: globalRoleBindingFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2GlobalRoleBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	globalRole := expandGlobalRoleBinding(d)

	err := meta.(*Config).GlobalRoleExist(globalRole.GlobalRoleID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Global Role Binding %s", globalRole.Name)

	newGlobalRole, err := client.GlobalRoleBinding.Create(globalRole)
	if err != nil {
		return err
	}

	d.SetId(newGlobalRole.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalRoleBindingStateRefreshFunc(client, newGlobalRole.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be created: %s", newGlobalRole.ID, waitErr)
	}

	return resourceRancher2GlobalRoleBindingRead(ctx, d, meta)
}

func resourceRancher2GlobalRoleBindingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Global Role Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRoleBinding.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Global Role Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenGlobalRoleBinding(d, globalRole)
}

func resourceRancher2GlobalRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Global Role Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRoleBinding.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newGlobalRole, err := client.GlobalRoleBinding.Update(globalRole, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalRoleBindingStateRefreshFunc(client, newGlobalRole.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be updated: %s", newGlobalRole.ID, waitErr)
	}

	return resourceRancher2GlobalRoleBindingRead(ctx, d, meta)
}

func resourceRancher2GlobalRoleBindingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Global Role Binding ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalRole, err := client.GlobalRoleBinding.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Global Role Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.GlobalRoleBinding.Delete(globalRole)
	if err != nil {
		return fmt.Errorf("Error removing Global Role Binding: %s", err)
	}

	log.Printf("[DEBUG] Waiting for global role binding (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    globalRoleBindingStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// globalRoleBindingStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Global Role Binding.
func globalRoleBindingStateRefreshFunc(client *managementClient.Client, globalRoleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalRoleBinding.ByID(globalRoleID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}
