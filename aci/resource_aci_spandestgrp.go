package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciSPANDestinationGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSPANDestinationGroupCreate,
		Update: resourceAciSPANDestinationGroupUpdate,
		Read:   resourceAciSPANDestinationGroupRead,
		Delete: resourceAciSPANDestinationGroupDelete,

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
	// d.Set("tenant_dn", GetParentDn(spanDestGrp.DistinguishedName))
	if dn != spanDestGrp.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	spanDestGrpMap, _ := spanDestGrp.ToMap()

	d.Set("name", spanDestGrpMap["name"])

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

func resourceAciSPANDestinationGroupCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(spanDestGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANDestinationGroupRead(d, m)
}

func resourceAciSPANDestinationGroupUpdate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(spanDestGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANDestinationGroupRead(d, m)

}

func resourceAciSPANDestinationGroupRead(d *schema.ResourceData, m interface{}) error {
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

func resourceAciSPANDestinationGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanDestGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
