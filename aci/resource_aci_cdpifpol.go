package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCDPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCDPInterfacePolicyCreate,
		UpdateContext: resourceAciCDPInterfacePolicyUpdate,
		ReadContext:   resourceAciCDPInterfacePolicyRead,
		DeleteContext: resourceAciCDPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCDPInterfacePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled",
					"disabled",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteCDPInterfacePolicy(client *client.Client, dn string) (*models.CDPInterfacePolicy, error) {
	cdpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cdpIfPol := models.CDPInterfacePolicyFromContainer(cdpIfPolCont)

	if cdpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("CDPInterfacePolicy %s not found", cdpIfPol.DistinguishedName)
	}

	return cdpIfPol, nil
}

func setCDPInterfacePolicyAttributes(cdpIfPol *models.CDPInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cdpIfPol.DistinguishedName)
	d.Set("description", cdpIfPol.Description)
	cdpIfPolMap, err := cdpIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", cdpIfPolMap["name"])

	d.Set("admin_st", cdpIfPolMap["adminSt"])
	d.Set("annotation", cdpIfPolMap["annotation"])
	d.Set("name_alias", cdpIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciCDPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cdpIfPol, err := getRemoteCDPInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCDPInterfacePolicyAttributes(cdpIfPol, d)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCDPInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CDPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	cdpIfPolAttr := models.CDPInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cdpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cdpIfPolAttr.Annotation = Annotation.(string)
	} else {
		cdpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cdpIfPolAttr.NameAlias = NameAlias.(string)
	}
	cdpIfPol := models.NewCDPInterfacePolicy(fmt.Sprintf("infra/cdpIfP-%s", name), "uni", desc, cdpIfPolAttr)

	err := aciClient.Save(cdpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cdpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCDPInterfacePolicyRead(ctx, d, m)
}

func resourceAciCDPInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CDPInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	cdpIfPolAttr := models.CDPInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cdpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cdpIfPolAttr.Annotation = Annotation.(string)
	} else {
		cdpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cdpIfPolAttr.NameAlias = NameAlias.(string)
	}
	cdpIfPol := models.NewCDPInterfacePolicy(fmt.Sprintf("infra/cdpIfP-%s", name), "uni", desc, cdpIfPolAttr)

	cdpIfPol.Status = "modified"

	err := aciClient.Save(cdpIfPol)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cdpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCDPInterfacePolicyRead(ctx, d, m)

}

func resourceAciCDPInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cdpIfPol, err := getRemoteCDPInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setCDPInterfacePolicyAttributes(cdpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCDPInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cdpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
