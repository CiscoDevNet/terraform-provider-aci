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

func resourceAciSNMPCommunity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSNMPCommunityCreate,
		UpdateContext: resourceAciSNMPCommunityUpdate,
		ReadContext:   resourceAciSNMPCommunityRead,
		DeleteContext: resourceAciSNMPCommunityDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSNMPCommunityImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_snmp_context_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNameAttribute(),
			},
		})),
	}
}

func getRemoteSNMPCommunity(client *client.Client, dn string) (*models.SNMPCommunity, error) {
	snmpCommunityPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpCommunityP := models.SNMPCommunityFromContainer(snmpCommunityPCont)
	if snmpCommunityP.DistinguishedName == "" {
		return nil, fmt.Errorf("SNMPCommunity %s not found", snmpCommunityP.DistinguishedName)
	}
	return snmpCommunityP, nil
}

func setSNMPCommunityAttributes(snmpCommunityP *models.SNMPCommunity, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(snmpCommunityP.DistinguishedName)
	d.Set("description", snmpCommunityP.Description)
	snmpCommunityPMap, err := snmpCommunityP.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("name", snmpCommunityPMap["name"])
	d.Set("name_alias", snmpCommunityPMap["nameAlias"])
	d.Set("annotation", snmpCommunityPMap["annotation"])
	d.Set("vrf_snmp_context_dn", GetParentDn(d.Id(), fmt.Sprintf("/community-%s", snmpCommunityPMap["name"])))

	return d, nil
}

func resourceAciSNMPCommunityImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	snmpCommunityP, err := getRemoteSNMPCommunity(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSNMPCommunityAttributes(snmpCommunityP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSNMPCommunityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMPCommunity: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VRFSNMPCtxDn := d.Get("vrf_snmp_context_dn").(string)

	snmpCommunityPAttr := models.SNMPCommunityAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpCommunityPAttr.Annotation = Annotation.(string)
	} else {
		snmpCommunityPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpCommunityPAttr.Name = Name.(string)
	}
	snmpCommunityP := models.NewSNMPCommunity(fmt.Sprintf("community-%s", name), VRFSNMPCtxDn, desc, nameAlias, snmpCommunityPAttr)

	err := aciClient.Save(snmpCommunityP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpCommunityP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSNMPCommunityRead(ctx, d, m)
}

func resourceAciSNMPCommunityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMPCommunity: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VRFSNMPCtxDn := d.Get("vrf_snmp_context_dn").(string)
	snmpCommunityPAttr := models.SNMPCommunityAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpCommunityPAttr.Annotation = Annotation.(string)
	} else {
		snmpCommunityPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpCommunityPAttr.Name = Name.(string)
	}
	snmpCommunityP := models.NewSNMPCommunity(fmt.Sprintf("community-%s", name), VRFSNMPCtxDn, desc, nameAlias, snmpCommunityPAttr)

	snmpCommunityP.Status = "modified"
	err := aciClient.Save(snmpCommunityP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpCommunityP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSNMPCommunityRead(ctx, d, m)
}

func resourceAciSNMPCommunityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	snmpCommunityP, err := getRemoteSNMPCommunity(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setSNMPCommunityAttributes(snmpCommunityP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSNMPCommunityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "snmpCommunityP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
