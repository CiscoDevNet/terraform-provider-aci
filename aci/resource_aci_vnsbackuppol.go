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

func resourceAciPBRBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPBRBackupPolicyCreate,
		UpdateContext: resourceAciPBRBackupPolicyUpdate,
		ReadContext:   resourceAciPBRBackupPolicyRead,
		DeleteContext: resourceAciPBRBackupPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPBRBackupPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemotePBRBackupPolicy(client *client.Client, dn string) (*models.PBRBackupPolicy, error) {
	vnsBackupPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsBackupPol := models.PBRBackupPolicyFromContainer(vnsBackupPolCont)
	if vnsBackupPol.DistinguishedName == "" {
		return nil, fmt.Errorf("PBR Backup Policy %s not found", dn)
	}
	return vnsBackupPol, nil
}

func setPBRBackupPolicyAttributes(vnsBackupPol *models.PBRBackupPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vnsBackupPol.DistinguishedName)
	d.Set("description", vnsBackupPol.Description)
	if dn != vnsBackupPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vnsBackupPolMap, err := vnsBackupPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", vnsBackupPolMap["annotation"])
	d.Set("name", vnsBackupPolMap["name"])
	d.Set("name_alias", vnsBackupPolMap["nameAlias"])
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/%s/%s", models.RnvnsSvcCont, fmt.Sprintf(models.RnvnsBackupPol, vnsBackupPolMap["name"]))))

	return d, nil
}

func resourceAciPBRBackupPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsBackupPol, err := getRemotePBRBackupPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPBRBackupPolicyAttributes(vnsBackupPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPBRBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PBRBackupPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	vnsBackupPolAttr := models.PBRBackupPolicyAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsBackupPolAttr.Annotation = Annotation.(string)
	} else {
		vnsBackupPolAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsBackupPolAttr.Name = Name.(string)
	}

	parentDn := fmt.Sprintf("%s/%s", TenantDn, models.RnvnsSvcCont)

	vnsBackupPol := models.NewPBRBackupPolicy(fmt.Sprintf(models.RnvnsBackupPol, name), parentDn, desc, nameAlias, vnsBackupPolAttr)

	err := aciClient.Save(vnsBackupPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vnsBackupPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPBRBackupPolicyRead(ctx, d, m)
}

func resourceAciPBRBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PBRBackupPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	vnsBackupPolAttr := models.PBRBackupPolicyAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsBackupPolAttr.Annotation = Annotation.(string)
	} else {
		vnsBackupPolAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsBackupPolAttr.Name = Name.(string)
	}

	parentDn := fmt.Sprintf("%s/%s", TenantDn, models.RnvnsSvcCont)

	vnsBackupPol := models.NewPBRBackupPolicy(fmt.Sprintf(models.RnvnsBackupPol, name), parentDn, desc, nameAlias, vnsBackupPolAttr)

	err := aciClient.Save(vnsBackupPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vnsBackupPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPBRBackupPolicyRead(ctx, d, m)
}

func resourceAciPBRBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsBackupPol, err := getRemotePBRBackupPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setPBRBackupPolicyAttributes(vnsBackupPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPBRBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsBackupPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
