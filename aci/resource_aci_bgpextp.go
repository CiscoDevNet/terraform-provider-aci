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

func resourceAciL3outBgpExternalPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outBgpExternalPolicyCreate,
		UpdateContext: resourceAciL3outBgpExternalPolicyUpdate,
		ReadContext:   resourceAciL3outBgpExternalPolicyRead,
		DeleteContext: resourceAciL3outBgpExternalPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outBgpExternalPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
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
func getRemoteL3outBgpExternalPolicy(client *client.Client, dn string) (*models.L3outBgpExternalPolicy, error) {
	bgpExtPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpExtP := models.L3outBgpExternalPolicyFromContainer(bgpExtPCont)

	if bgpExtP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outBgpExternalPolicy %s not found", bgpExtP.DistinguishedName)
	}

	return bgpExtP, nil
}

func setL3outBgpExternalPolicyAttributes(bgpExtP *models.L3outBgpExternalPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bgpExtP.DistinguishedName)
	d.Set("description", bgpExtP.Description)
	dn := d.Id()
	if dn != bgpExtP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	bgpExtPMap, err := bgpExtP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("l3_outside_dn", GetParentDn(dn, "/bgpExtP"))

	d.Set("annotation", bgpExtPMap["annotation"])
	d.Set("name_alias", bgpExtPMap["nameAlias"])
	return d, nil
}

func resourceAciL3outBgpExternalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpExtP, err := getRemoteL3outBgpExternalPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outBgpExternalPolicyAttributes(bgpExtP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outBgpExternalPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outBgpExternalPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	L3OutsideDn := d.Get("l3_outside_dn").(string)

	bgpExtPAttr := models.L3outBgpExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpExtPAttr.Annotation = Annotation.(string)
	} else {
		bgpExtPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpExtPAttr.NameAlias = NameAlias.(string)
	}
	bgpExtP := models.NewL3outBgpExternalPolicy(fmt.Sprintf("bgpExtP"), L3OutsideDn, desc, bgpExtPAttr)

	err := aciClient.Save(bgpExtP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	d.SetId(bgpExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outBgpExternalPolicyRead(ctx, d, m)
}

func resourceAciL3outBgpExternalPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outBgpExternalPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	bgpExtPAttr := models.L3outBgpExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpExtPAttr.Annotation = Annotation.(string)
	} else {
		bgpExtPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpExtPAttr.NameAlias = NameAlias.(string)
	}
	bgpExtP := models.NewL3outBgpExternalPolicy(fmt.Sprintf("bgpExtP"), L3OutsideDn, desc, bgpExtPAttr)

	bgpExtP.Status = "modified"

	err := aciClient.Save(bgpExtP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	d.SetId(bgpExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outBgpExternalPolicyRead(ctx, d, m)

}

func resourceAciL3outBgpExternalPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpExtP, err := getRemoteL3outBgpExternalPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3outBgpExternalPolicyAttributes(bgpExtP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outBgpExternalPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpExtP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
