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

func resourceAciSPANDestinationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSPANDestinationGroupCreate,
		UpdateContext: resourceAciSPANDestinationGroupUpdate,
		ReadContext:   resourceAciSPANDestinationGroupRead,
		DeleteContext: resourceAciSPANDestinationGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSPANDestinationGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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
func getRemoteSPANDestinationGroup(client *client.Client, dn string) (*models.SPANDestinationGroup, error) {
	spanDestGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	spanDestGrp := models.SPANDestinationGroupFromContainer(spanDestGrpCont)

	if spanDestGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("SPANDestinationGroup %s not found", spanDestGrp.DistinguishedName)
	}

	return spanDestGrp, nil
}

func setSPANDestinationGroupAttributes(spanDestGrp *models.SPANDestinationGroup, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(spanDestGrp.DistinguishedName)
	d.Set("description", spanDestGrp.Description)

	if dn != spanDestGrp.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	spanDestGrpMap, _ := spanDestGrp.ToMap()

	d.Set("name", spanDestGrpMap["name"])
	d.Set("tenant_dn", GetParentDn(spanDestGrp.DistinguishedName, fmt.Sprintf("destgrp-%s", spanDestGrpMap["name"])))
	d.Set("annotation", spanDestGrpMap["annotation"])
	d.Set("name_alias", spanDestGrpMap["nameAlias"])
	return d
}

func resourceAciSPANDestinationGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	spanDestGrp, err := getRemoteSPANDestinationGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	spanDestGrpMap, _ := spanDestGrp.ToMap()
	name := spanDestGrpMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/destgrp-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled := setSPANDestinationGroupAttributes(spanDestGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSPANDestinationGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANDestinationGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	spanDestGrpAttr := models.SPANDestinationGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanDestGrpAttr.Annotation = Annotation.(string)
	} else {
		spanDestGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanDestGrpAttr.NameAlias = NameAlias.(string)
	}
	spanDestGrp := models.NewSPANDestinationGroup(fmt.Sprintf("destgrp-%s", name), TenantDn, desc, spanDestGrpAttr)

	err := aciClient.Save(spanDestGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(spanDestGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANDestinationGroupRead(ctx, d, m)
}

func resourceAciSPANDestinationGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANDestinationGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	spanDestGrpAttr := models.SPANDestinationGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanDestGrpAttr.Annotation = Annotation.(string)
	} else {
		spanDestGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanDestGrpAttr.NameAlias = NameAlias.(string)
	}
	spanDestGrp := models.NewSPANDestinationGroup(fmt.Sprintf("destgrp-%s", name), TenantDn, desc, spanDestGrpAttr)

	spanDestGrp.Status = "modified"

	err := aciClient.Save(spanDestGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(spanDestGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANDestinationGroupRead(ctx, d, m)

}

func resourceAciSPANDestinationGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	spanDestGrp, err := getRemoteSPANDestinationGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSPANDestinationGroupAttributes(spanDestGrp, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSPANDestinationGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanDestGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
