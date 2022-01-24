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

func resourceAciSNMPContextProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSNMPContextProfileCreate,
		UpdateContext: resourceAciSNMPContextProfileUpdate,
		ReadContext:   resourceAciSNMPContextProfileRead,
		DeleteContext: resourceAciSNMPContextProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSNMPContextProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		}),
	}
}

func getRemoteSNMPContextProfile(client *client.Client, dn string) (*models.SNMPContextProfile, error) {
	snmpCtxPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpCtxP := models.SNMPContextProfileFromContainer(snmpCtxPCont)
	if snmpCtxP.DistinguishedName == "" {
		return nil, fmt.Errorf("SNMPContextProfile %s not found", snmpCtxP.DistinguishedName)
	}
	return snmpCtxP, nil
}

func setSNMPContextProfileAttributes(snmpCtxP *models.SNMPContextProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(snmpCtxP.DistinguishedName)
	snmpCtxPMap, err := snmpCtxP.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("vrf_dn", GetParentDn(d.Id(), fmt.Sprintf("/snmpctx")))
	d.Set("annotation", snmpCtxPMap["annotation"])
	d.Set("name", snmpCtxPMap["name"])
	d.Set("name_alias", snmpCtxPMap["nameAlias"])
	return d, nil
}

func resourceAciSNMPContextProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	snmpCtxP, err := getRemoteSNMPContextProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSNMPContextProfileAttributes(snmpCtxP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSNMPContextProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMPContextProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	VRFDn := d.Get("vrf_dn").(string)

	snmpCtxPAttr := models.SNMPContextProfileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpCtxPAttr.Annotation = Annotation.(string)
	} else {
		snmpCtxPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpCtxPAttr.Name = Name.(string)
	}
	snmpCtxP := models.NewSNMPContextProfile(fmt.Sprintf("snmpctx"), VRFDn, nameAlias, snmpCtxPAttr)

	err := aciClient.Save(snmpCtxP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpCtxP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSNMPContextProfileRead(ctx, d, m)
}

func resourceAciSNMPContextProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMPContextProfile: Beginning Update")
	aciClient := m.(*client.Client)
	VRFDn := d.Get("vrf_dn").(string)
	snmpCtxPAttr := models.SNMPContextProfileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpCtxPAttr.Annotation = Annotation.(string)
	} else {
		snmpCtxPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpCtxPAttr.Name = Name.(string)
	}
	snmpCtxP := models.NewSNMPContextProfile(fmt.Sprintf("snmpctx"), VRFDn, nameAlias, snmpCtxPAttr)

	snmpCtxP.Status = "modified"
	err := aciClient.Save(snmpCtxP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpCtxP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSNMPContextProfileRead(ctx, d, m)
}

func resourceAciSNMPContextProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	snmpCtxP, err := getRemoteSNMPContextProfile(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setSNMPContextProfileAttributes(snmpCtxP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSNMPContextProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "snmpCtxP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
