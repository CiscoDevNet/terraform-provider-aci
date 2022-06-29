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

func resourceAciAbstractionofIPAddressBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAbstractionofIPAddressBlockCreate,
		UpdateContext: resourceAciAbstractionofIPAddressBlockUpdate,
		ReadContext:   resourceAciAbstractionofIPAddressBlockRead,
		DeleteContext: resourceAciAbstractionofIPAddressBlockDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAbstractionofIPAddressBlockImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"multicast_pool_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func getRemoteAbstractionofIPAddressBlock(client *client.Client, dn string) (*models.AbstractionofIPAddressBlock, error) {
	fvnsMcastAddrBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvnsMcastAddrBlk := models.AbstractionofIPAddressBlockFromContainer(fvnsMcastAddrBlkCont)
	if fvnsMcastAddrBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("AbstractionofIPAddressBlock %s not found", dn)
	}
	return fvnsMcastAddrBlk, nil
}

func setAbstractionofIPAddressBlockAttributes(fvnsMcastAddrBlk *models.AbstractionofIPAddressBlock, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvnsMcastAddrBlk.DistinguishedName)
	d.Set("description", fvnsMcastAddrBlk.Description)
	fvnsMcastAddrBlkMap, err := fvnsMcastAddrBlk.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != fvnsMcastAddrBlk.DistinguishedName {
		d.Set("multicast_pool_dn", "")
	} else {
		d.Set("multicast_pool_dn", GetParentDn(fvnsMcastAddrBlk.DistinguishedName, fmt.Sprintf("/"+models.RnfvnsMcastAddrBlk)))
	}
	d.Set("annotation", fvnsMcastAddrBlkMap["annotation"])
	d.Set("from", fvnsMcastAddrBlkMap["from"])
	d.Set("name", fvnsMcastAddrBlkMap["name"])
	d.Set("name_alias", fvnsMcastAddrBlkMap["nameAlias"])
	d.Set("to", fvnsMcastAddrBlkMap["to"])
	return d, nil
}

func resourceAciAbstractionofIPAddressBlockImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvnsMcastAddrBlk, err := getRemoteAbstractionofIPAddressBlock(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAbstractionofIPAddressBlockAttributes(fvnsMcastAddrBlk, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAbstractionofIPAddressBlockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AbstractionofIPAddressBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	log.Printf("[DEBUG] AbstractionofIPAddressBlock: desc %s", desc)
	from := d.Get("from").(string)
	to := d.Get("to").(string)
	MulticastAddressPoolDn := d.Get("multicast_pool_dn").(string)

	fvnsMcastAddrBlkAttr := models.AbstractionofIPAddressBlockAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrBlkAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrBlkAttr.Annotation = "{}"
	}

	if From, ok := d.GetOk("from"); ok {
		fvnsMcastAddrBlkAttr.From = From.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrBlkAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrBlkAttr.NameAlias = NameAlias.(string)
	}

	if To, ok := d.GetOk("to"); ok {
		fvnsMcastAddrBlkAttr.To = To.(string)
	}
	fvnsMcastAddrBlk := models.NewAbstractionofIPAddressBlock(fmt.Sprintf(models.RnfvnsMcastAddrBlk, from, to), MulticastAddressPoolDn, desc, fvnsMcastAddrBlkAttr)

	err := aciClient.Save(fvnsMcastAddrBlk)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsMcastAddrBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciAbstractionofIPAddressBlockRead(ctx, d, m)
}
func resourceAciAbstractionofIPAddressBlockUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AbstractionofIPAddressBlock: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	from := d.Get("from").(string)
	to := d.Get("to").(string)
	MulticastAddressPoolDn := d.Get("multicast_pool_dn").(string)

	fvnsMcastAddrBlkAttr := models.AbstractionofIPAddressBlockAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrBlkAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrBlkAttr.Annotation = "{}"
	}

	if From, ok := d.GetOk("from"); ok {
		fvnsMcastAddrBlkAttr.From = From.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrBlkAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrBlkAttr.NameAlias = NameAlias.(string)
	}

	if To, ok := d.GetOk("to"); ok {
		fvnsMcastAddrBlkAttr.To = To.(string)
	}
	fvnsMcastAddrBlk := models.NewAbstractionofIPAddressBlock(fmt.Sprintf(models.RnfvnsMcastAddrBlk, from, to), MulticastAddressPoolDn, desc, fvnsMcastAddrBlkAttr)

	fvnsMcastAddrBlk.Status = "modified"

	err := aciClient.Save(fvnsMcastAddrBlk)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsMcastAddrBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciAbstractionofIPAddressBlockRead(ctx, d, m)
}

func resourceAciAbstractionofIPAddressBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvnsMcastAddrBlk, err := getRemoteAbstractionofIPAddressBlock(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setAbstractionofIPAddressBlockAttributes(fvnsMcastAddrBlk, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAbstractionofIPAddressBlockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvnsMcastAddrBlk")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
