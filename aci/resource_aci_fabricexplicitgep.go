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

func resourceAciVPCExplicitProtectionGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVPCExplicitProtectionGroupCreate,
		UpdateContext: resourceAciVPCExplicitProtectionGroupUpdate,
		ReadContext:   resourceAciVPCExplicitProtectionGroupRead,
		DeleteContext: resourceAciVPCExplicitProtectionGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVPCExplicitProtectionGroupImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"switch1": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"switch2": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_domain_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},
			"vpc_explicit_protection_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func getRemoteVPCExplicitProtectionGroup(client *client.Client, dn string) (*models.VPCExplicitProtectionGroup, error) {
	baseurlStr := "/api/node/mo"
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children", baseurlStr, dn)
	fabricExplicitGEpCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	fabricExplicitGEp := models.VPCExplicitProtectionGroupFromContainer(fabricExplicitGEpCont)

	if fabricExplicitGEp.DistinguishedName == "" {
		return nil, fmt.Errorf("VPCExplicitProtectionGroup %s not found", fabricExplicitGEp.DistinguishedName)
	}

	return fabricExplicitGEp, nil
}

func setVPCExplicitProtectionGroupAttributes(fabricExplicitGEp *models.VPCExplicitProtectionGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricExplicitGEp.DistinguishedName)

	fabricExplicitGEpMap, err := fabricExplicitGEp.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", fabricExplicitGEpMap["name"])

	d.Set("annotation", fabricExplicitGEpMap["annotation"])
	d.Set("vpc_explicit_protection_group_id", fabricExplicitGEpMap["id"])
	d.Set("switch1", fabricExplicitGEpMap["switch1"])
	d.Set("switch2", fabricExplicitGEpMap["switch2"])

	sw1 := d.Get("switch1").(string)
	sw2 := d.Get("switch2").(string)

	if sw1 != fabricExplicitGEpMap["switch1"] && sw1 != fabricExplicitGEpMap["switch2"] {
		d.Set("switch1", "")
	}
	if sw2 != fabricExplicitGEpMap["switch1"] && sw2 != fabricExplicitGEpMap["switch2"] {
		d.Set("switch2", "")
	}

	d.Set("vpc_domain_policy", fabricExplicitGEpMap["vpc_domain_policy"])

	return d, nil
}

func resourceAciVPCExplicitProtectionGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fabricExplicitGEp, err := getRemoteVPCExplicitProtectionGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVPCExplicitProtectionGroupAttributes(fabricExplicitGEp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVPCExplicitProtectionGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VPCExplicitProtectionGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	switch1 := d.Get("switch1").(string)
	switch2 := d.Get("switch2").(string)

	vpcDomainPolicy := ""

	if temp, ok := d.GetOk("vpc_domain_policy"); ok {
		vpcDomainPolicy = temp.(string)
	}

	fabricExplicitGEpAttr := models.VPCExplicitProtectionGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricExplicitGEpAttr.Annotation = Annotation.(string)
	} else {
		fabricExplicitGEpAttr.Annotation = "{}"
	}
	if VPCExplicitProtectionGroup_id, ok := d.GetOk("vpc_explicit_protection_group_id"); ok {
		fabricExplicitGEpAttr.VPCExplicitProtectionGroup_id = VPCExplicitProtectionGroup_id.(string)
	}
	fabricExplicitGEp := models.NewVPCExplicitProtectionGroup(fmt.Sprintf("fabric/protpol/expgep-%s", name), "uni", fabricExplicitGEpAttr)
	fabricExplicitGEp, err := aciClient.CreateVPCExplicitProtectionGroup(name, switch1, switch2, vpcDomainPolicy, fabricExplicitGEpAttr)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricExplicitGEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVPCExplicitProtectionGroupRead(ctx, d, m)
}

func resourceAciVPCExplicitProtectionGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VPCExplicitProtectionGroup: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)
	switch1 := d.Get("switch1").(string)
	switch2 := d.Get("switch2").(string)
	vpcDomainPolicy := ""

	if temp, ok := d.GetOk("vpc_domain_policy"); ok {
		vpcDomainPolicy = temp.(string)
	}

	fabricExplicitGEpAttr := models.VPCExplicitProtectionGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricExplicitGEpAttr.Annotation = Annotation.(string)
	} else {
		fabricExplicitGEpAttr.Annotation = "{}"
	}
	if VPCExplicitProtectionGroup_id, ok := d.GetOk("vpc_explicit_protection_group_id"); ok {
		fabricExplicitGEpAttr.VPCExplicitProtectionGroup_id = VPCExplicitProtectionGroup_id.(string)
	}
	fabricExplicitGEp := models.NewVPCExplicitProtectionGroup(fmt.Sprintf("fabric/protpol/expgep-%s", name), "uni", fabricExplicitGEpAttr)

	fabricExplicitGEp.Status = "modified"
	fabricExplicitGEp, err := aciClient.UpdateVPCExplicitProtectionGroup(name, switch1, switch2, vpcDomainPolicy, fabricExplicitGEpAttr)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricExplicitGEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVPCExplicitProtectionGroupRead(ctx, d, m)

}

func resourceAciVPCExplicitProtectionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fabricExplicitGEp, err := getRemoteVPCExplicitProtectionGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setVPCExplicitProtectionGroupAttributes(fabricExplicitGEp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVPCExplicitProtectionGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricExplicitGEp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
