package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciVMMCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVMMCredentialCreate,
		UpdateContext: resourceAciVMMCredentialUpdate,
		ReadContext:   resourceAciVMMCredentialRead,
		DeleteContext: resourceAciVMMCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMCredentialImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pwd": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"usr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteVMMCredential(client *client.Client, dn string) (*models.VMMCredential, error) {
	vmmUsrAccPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmUsrAccP := models.VMMCredentialFromContainer(vmmUsrAccPCont)
	if vmmUsrAccP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMMCredential %s not found", vmmUsrAccP.DistinguishedName)
	}
	return vmmUsrAccP, nil
}

func setVMMCredentialAttributes(vmmUsrAccP *models.VMMCredential, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vmmUsrAccP.DistinguishedName)
	d.Set("description", vmmUsrAccP.Description)
	vmmUsrAccPMap, err := vmmUsrAccP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("vmm_domain_dn", GetParentDn(vmmUsrAccP.DistinguishedName, fmt.Sprintf("/usracc-%s", vmmUsrAccPMap["name"])))
	d.Set("annotation", vmmUsrAccPMap["annotation"])
	d.Set("name", vmmUsrAccPMap["name"])
	d.Set("usr", vmmUsrAccPMap["usr"])
	d.Set("name_alias", vmmUsrAccPMap["nameAlias"])
	return d, nil
}

func resourceAciVMMCredentialImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVMMCredentialAttributes(vmmUsrAccP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMCredential: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	} else {
		vmmUsrAccPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmUsrAccPAttr.Name = Name.(string)
	}

	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}

	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, nameAlias, vmmUsrAccPAttr)

	err := aciClient.Save(vmmUsrAccP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vmmUsrAccP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciVMMCredentialRead(ctx, d, m)
}

func resourceAciVMMCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMCredential: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	} else {
		vmmUsrAccPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmUsrAccPAttr.Name = Name.(string)
	}

	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}

	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, nameAlias, vmmUsrAccPAttr)

	vmmUsrAccP.Status = "modified"
	err := aciClient.Save(vmmUsrAccP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vmmUsrAccP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciVMMCredentialRead(ctx, d, m)
}

func resourceAciVMMCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setVMMCredentialAttributes(vmmUsrAccP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciVMMCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmUsrAccP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
