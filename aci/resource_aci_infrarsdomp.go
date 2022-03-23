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

func resourceAciInfraRsDomP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciInfraRsDomPCreate,
		UpdateContext: resourceAciInfraRsDomPUpdate,
		ReadContext:   resourceAciInfraRsDomPRead,
		DeleteContext: resourceAciInfraRsDomPDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInfraRsDomPImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func getRemoteInfraRsDomP(client *client.Client, dn string) (*models.InfraRsDomP, error) {
	infraRsDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraRsDomP := models.InfraRsDomPFromContainer(infraRsDomPCont)
	if infraRsDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("InfraRsDomP %s not found", infraRsDomP.DistinguishedName)
	}
	return infraRsDomP, nil
}

func setInfraRsDomPAttributes(infraRsDomP *models.InfraRsDomP, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraRsDomP.DistinguishedName)
	if dn != infraRsDomP.DistinguishedName {
		d.Set("attachable_access_entity_profile_dn", "")
	}
	infraRsDomPMap, err := infraRsDomP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", infraRsDomPMap["annotation"])
	d.Set("domain_dn", infraRsDomPMap["tDn"])
	return d, nil
}

func resourceAciInfraRsDomPImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraRsDomP, err := getRemoteInfraRsDomP(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setInfraRsDomPAttributes(infraRsDomP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciInfraRsDomPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InfraRsDomP: Beginning Creation")
	aciClient := m.(*client.Client)
	tDn := d.Get("domain_dn").(string)
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraRsDomPAttr := models.InfraRsDomPAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsDomPAttr.Annotation = Annotation.(string)
	} else {
		infraRsDomPAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("domain_dn"); ok {
		infraRsDomPAttr.TDn = TDn.(string)
	}
	infraRsDomP := models.NewInfraRsDomP(fmt.Sprintf(models.RninfraRsDomP, tDn), AttachableAccessEntityProfileDn, infraRsDomPAttr)

	err := aciClient.Save(infraRsDomP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciInfraRsDomPRead(ctx, d, m)
}

func resourceAciInfraRsDomPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InfraRsDomP: Beginning Update")
	aciClient := m.(*client.Client)
	tDn := d.Get("domain_dn").(string)
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraRsDomPAttr := models.InfraRsDomPAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsDomPAttr.Annotation = Annotation.(string)
	} else {
		infraRsDomPAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("domain_dn"); ok {
		infraRsDomPAttr.TDn = TDn.(string)
	}
	infraRsDomP := models.NewInfraRsDomP(fmt.Sprintf(models.RninfraRsDomP, tDn), AttachableAccessEntityProfileDn, infraRsDomPAttr)

	infraRsDomP.Status = "modified"

	err := aciClient.Save(infraRsDomP)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(infraRsDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciInfraRsDomPRead(ctx, d, m)
}

func resourceAciInfraRsDomPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	infraRsDomP, err := getRemoteInfraRsDomP(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setInfraRsDomPAttributes(infraRsDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciInfraRsDomPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "infraRsDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
