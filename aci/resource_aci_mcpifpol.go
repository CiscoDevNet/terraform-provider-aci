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

func resourceAciMiscablingProtocolInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMiscablingProtocolInterfacePolicyCreate,
		UpdateContext: resourceAciMiscablingProtocolInterfacePolicyUpdate,
		ReadContext:   resourceAciMiscablingProtocolInterfacePolicyRead,
		DeleteContext: resourceAciMiscablingProtocolInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMiscablingProtocolInterfacePolicyImport,
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
func getRemoteMiscablingProtocolInterfacePolicy(client *client.Client, dn string) (*models.MiscablingProtocolInterfacePolicy, error) {
	mcpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	mcpIfPol := models.MiscablingProtocolInterfacePolicyFromContainer(mcpIfPolCont)

	if mcpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("MiscablingProtocolInterfacePolicy %s not found", mcpIfPol.DistinguishedName)
	}

	return mcpIfPol, nil
}

func setMiscablingProtocolInterfacePolicyAttributes(mcpIfPol *models.MiscablingProtocolInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mcpIfPol.DistinguishedName)
	d.Set("description", mcpIfPol.Description)
	mcpIfPolMap, err := mcpIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", mcpIfPolMap["name"])

	d.Set("admin_st", mcpIfPolMap["adminSt"])
	d.Set("annotation", mcpIfPolMap["annotation"])
	d.Set("name_alias", mcpIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciMiscablingProtocolInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	mcpIfPol, err := getRemoteMiscablingProtocolInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMiscablingProtocolInterfacePolicyAttributes(mcpIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMiscablingProtocolInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MiscablingProtocolInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	mcpIfPolAttr := models.MiscablingProtocolInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		mcpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mcpIfPolAttr.Annotation = Annotation.(string)
	} else {
		mcpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mcpIfPolAttr.NameAlias = NameAlias.(string)
	}
	mcpIfPol := models.NewMiscablingProtocolInterfacePolicy(fmt.Sprintf("infra/mcpIfP-%s", name), "uni", desc, mcpIfPolAttr)

	err := aciClient.Save(mcpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mcpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMiscablingProtocolInterfacePolicyRead(ctx, d, m)
}

func resourceAciMiscablingProtocolInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MiscablingProtocolInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	mcpIfPolAttr := models.MiscablingProtocolInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		mcpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mcpIfPolAttr.Annotation = Annotation.(string)
	} else {
		mcpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mcpIfPolAttr.NameAlias = NameAlias.(string)
	}
	mcpIfPol := models.NewMiscablingProtocolInterfacePolicy(fmt.Sprintf("infra/mcpIfP-%s", name), "uni", desc, mcpIfPolAttr)

	mcpIfPol.Status = "modified"

	err := aciClient.Save(mcpIfPol)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(mcpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMiscablingProtocolInterfacePolicyRead(ctx, d, m)

}

func resourceAciMiscablingProtocolInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	mcpIfPol, err := getRemoteMiscablingProtocolInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setMiscablingProtocolInterfacePolicyAttributes(mcpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMiscablingProtocolInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mcpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
