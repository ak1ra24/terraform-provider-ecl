package ecl

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nttcom/eclcloud/ecl/identity/v3/roles"
)

func resourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityRoleV3Create,
		Read:   resourceIdentityRoleV3Read,
		Update: resourceIdentityRoleV3Update,
		Delete: resourceIdentityRoleV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityRoleV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating ECL identity client: %s", err)
	}

	createOpts := roles.CreateOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	role, err := roles.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating ECL role: %s", err)
	}

	d.SetId(role.ID)

	return resourceIdentityRoleV3Read(d, meta)
}

func resourceIdentityRoleV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating ECL identity client: %s", err)
	}

	role, err := roles.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "role")
	}

	log.Printf("[DEBUG] Retrieved ECL role: %#v", role)

	d.Set("domain_id", role.DomainID)
	d.Set("name", role.Name)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityRoleV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating ECL identity client: %s", err)
	}

	var hasChange bool
	var updateOpts roles.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		_, err := roles.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating ECL role: %s", err)
		}
	}

	return resourceIdentityRoleV3Read(d, meta)
}

func resourceIdentityRoleV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating ECL identity client: %s", err)
	}

	err = roles.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting ECL role: %s", err)
	}

	return nil
}
