package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2Certificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2CertificateCreate,
		ReadContext:   resourceRancher2CertificateRead,
		UpdateContext: resourceRancher2CertificateUpdate,
		DeleteContext: resourceRancher2CertificateDelete,

		Schema: certificateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	certificate, err := expandCertificate(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Certificate %s on Project ID %s", name, projectID)

	newCertificate, err := meta.(*Config).CreateCertificate(certificate)
	if err != nil {
		return err
	}

	err = flattenCertificate(d, newCertificate)
	if err != nil {
		return err
	}

	return resourceRancher2CertificateRead(ctx, d, meta)
}

func resourceRancher2CertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Certificate ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenCertificate(d, certificate)
}

func resourceRancher2CertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		return err
	}

	update, err := expandCertificate(d)
	if err != nil {
		return err
	}

	newCertificate, err := meta.(*Config).UpdateCertificate(certificate, update)
	if err != nil {
		return err
	}

	err = flattenCertificate(d, newCertificate)
	if err != nil {
		return err
	}

	return resourceRancher2CertificateRead(ctx, d, meta)
}

func resourceRancher2CertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Certificate ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = meta.(*Config).DeleteCertificate(certificate)
	if err != nil {
		return fmt.Errorf("Error removing Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
