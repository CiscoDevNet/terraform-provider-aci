package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciSecurityDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSecurityDomainCreate,
		UpdateContext: resourceAciSecurityDomainUpdate,
		ReadContext:   resourceAciSecurityDomainRead,
		DeleteContext: resourceAciSecurityDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSecurityDomainImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteSecurityDomain(client *client.Client, dn string) (*models.SecurityDomain, error) {
	aaaDomainCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaDomain := models.SecurityDomainFromContainer(aaaDomainCont)

	if aaaDomain.DistinguishedName == "" {
		return nil, fmt.Errorf("Security Domain %s not found", dn)
	}

	return aaaDomain, nil
}

func setSecurityDomainAttributes(aaaDomain *models.SecurityDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaDomain.DistinguishedName)
	d.Set("description", aaaDomain.Description)
	aaaDomainMap, err := aaaDomain.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", aaaDomainMap["name"])

	d.Set("annotation", aaaDomainMap["annotation"])
	d.Set("name_alias", aaaDomainMap["nameAlias"])
	return d, nil
}

func resourceAciSecurityDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSecurityDomainAttributes(aaaDomain, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSecurityDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SecurityDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaDomainAttr := models.SecurityDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaDomainAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaDomainAttr.NameAlias = NameAlias.(string)
	}
	aaaDomain := models.NewSecurityDomain(fmt.Sprintf("userext/domain-%s", name), "uni", desc, aaaDomainAttr)

	err := aciClient.Save(aaaDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSecurityDomainRead(ctx, d, m)
}

func resourceAciSecurityDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SecurityDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaDomainAttr := models.SecurityDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaDomainAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaDomainAttr.NameAlias = NameAlias.(string)
	}
	aaaDomain := models.NewSecurityDomain(fmt.Sprintf("userext/domain-%s", name), "uni", desc, aaaDomainAttr)

	aaaDomain.Status = "modified"

	err := aciClient.Save(aaaDomain)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSecurityDomainRead(ctx, d, m)

}

func resourceAciSecurityDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setSecurityDomainAttributes(aaaDomain, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSecurityDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaDomain")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
