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

func resourceAciLLDPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLLDPInterfacePolicyCreate,
		UpdateContext: resourceAciLLDPInterfacePolicyUpdate,
		ReadContext:   resourceAciLLDPInterfacePolicyRead,
		DeleteContext: resourceAciLLDPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLLDPInterfacePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_rx_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled",
					"disabled",
				}, false),
			},

			"admin_tx_st": &schema.Schema{
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
func getRemoteLLDPInterfacePolicy(client *client.Client, dn string) (*models.LLDPInterfacePolicy, error) {
	lldpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	lldpIfPol := models.LLDPInterfacePolicyFromContainer(lldpIfPolCont)

	if lldpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("LLDPInterfacePolicy %s not found", lldpIfPol.DistinguishedName)
	}

	return lldpIfPol, nil
}

func setLLDPInterfacePolicyAttributes(lldpIfPol *models.LLDPInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(lldpIfPol.DistinguishedName)
	d.Set("description", lldpIfPol.Description)
	lldpIfPolMap, err := lldpIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", lldpIfPolMap["name"])

	d.Set("admin_rx_st", lldpIfPolMap["adminRxSt"])
	d.Set("admin_tx_st", lldpIfPolMap["adminTxSt"])
	d.Set("annotation", lldpIfPolMap["annotation"])
	d.Set("name_alias", lldpIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciLLDPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	lldpIfPol, err := getRemoteLLDPInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLLDPInterfacePolicyAttributes(lldpIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLLDPInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LLDPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lldpIfPolAttr := models.LLDPInterfacePolicyAttributes{}
	if AdminRxSt, ok := d.GetOk("admin_rx_st"); ok {
		lldpIfPolAttr.AdminRxSt = AdminRxSt.(string)
	}
	if AdminTxSt, ok := d.GetOk("admin_tx_st"); ok {
		lldpIfPolAttr.AdminTxSt = AdminTxSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lldpIfPolAttr.Annotation = Annotation.(string)
	} else {
		lldpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lldpIfPolAttr.NameAlias = NameAlias.(string)
	}
	lldpIfPol := models.NewLLDPInterfacePolicy(fmt.Sprintf("infra/lldpIfP-%s", name), "uni", desc, lldpIfPolAttr)

	err := aciClient.Save(lldpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lldpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLLDPInterfacePolicyRead(ctx, d, m)
}

func resourceAciLLDPInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LLDPInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lldpIfPolAttr := models.LLDPInterfacePolicyAttributes{}
	if AdminRxSt, ok := d.GetOk("admin_rx_st"); ok {
		lldpIfPolAttr.AdminRxSt = AdminRxSt.(string)
	}
	if AdminTxSt, ok := d.GetOk("admin_tx_st"); ok {
		lldpIfPolAttr.AdminTxSt = AdminTxSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lldpIfPolAttr.Annotation = Annotation.(string)
	} else {
		lldpIfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lldpIfPolAttr.NameAlias = NameAlias.(string)
	}
	lldpIfPol := models.NewLLDPInterfacePolicy(fmt.Sprintf("infra/lldpIfP-%s", name), "uni", desc, lldpIfPolAttr)

	lldpIfPol.Status = "modified"

	err := aciClient.Save(lldpIfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lldpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLLDPInterfacePolicyRead(ctx, d, m)

}

func resourceAciLLDPInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	lldpIfPol, err := getRemoteLLDPInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLLDPInterfacePolicyAttributes(lldpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLLDPInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "lldpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
