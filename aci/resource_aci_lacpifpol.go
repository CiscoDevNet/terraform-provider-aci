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

func resourceAciLACPMemberPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLACPMemberPolicyCreate,
		UpdateContext: resourceAciLACPMemberPolicyUpdate,
		ReadContext:   resourceAciLACPMemberPolicyRead,
		DeleteContext: resourceAciLACPMemberPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLACPMemberPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prio": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tx_rate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"fast",
					"normal",
				}, false),
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func getRemoteLACPMemberPolicy(client *client.Client, dn string) (*models.LACPMemberPolicy, error) {
	lacpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	lacpIfPol := models.LACPMemberPolicyFromContainer(lacpIfPolCont)
	if lacpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("LACPMemberPolicy %s not found", dn)
	}
	return lacpIfPol, nil
}

func setLACPMemberPolicyAttributes(lacpIfPol *models.LACPMemberPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(lacpIfPol.DistinguishedName)
	d.Set("description", lacpIfPol.Description)
	lacpIfPolMap, err := lacpIfPol.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("annotation", lacpIfPolMap["annotation"])
	d.Set("name", lacpIfPolMap["name"])
	d.Set("name_alias", lacpIfPolMap["nameAlias"])
	d.Set("prio", lacpIfPolMap["prio"])
	d.Set("tx_rate", lacpIfPolMap["txRate"])
	return d, nil
}

func resourceAciLACPMemberPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	lacpIfPol, err := getRemoteLACPMemberPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLACPMemberPolicyAttributes(lacpIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLACPMemberPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LACPMemberPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	lacpIfPolAttr := models.LACPMemberPolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpIfPolAttr.Annotation = Annotation.(string)
	} else {
		lacpIfPolAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		lacpIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpIfPolAttr.NameAlias = NameAlias.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		lacpIfPolAttr.Prio = Prio.(string)
	}

	if TxRate, ok := d.GetOk("tx_rate"); ok {
		lacpIfPolAttr.TxRate = TxRate.(string)
	}
	lacpIfPol := models.NewLACPMemberPolicy(fmt.Sprintf(models.RnlacpIfPol, name), models.ParentDnlacpIfPol, desc, lacpIfPolAttr)
	err := aciClient.Save(lacpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lacpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLACPMemberPolicyRead(ctx, d, m)
}
func resourceAciLACPMemberPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LACPMemberPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	lacpIfPolAttr := models.LACPMemberPolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpIfPolAttr.Annotation = Annotation.(string)
	} else {
		lacpIfPolAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		lacpIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpIfPolAttr.NameAlias = NameAlias.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		lacpIfPolAttr.Prio = Prio.(string)
	}

	if TxRate, ok := d.GetOk("tx_rate"); ok {
		lacpIfPolAttr.TxRate = TxRate.(string)
	}
	lacpIfPol := models.NewLACPMemberPolicy(fmt.Sprintf(models.RnlacpIfPol, name), models.ParentDnlacpIfPol, desc, lacpIfPolAttr)
	lacpIfPol.Status = "modified"

	err := aciClient.Save(lacpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lacpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLACPMemberPolicyRead(ctx, d, m)
}

func resourceAciLACPMemberPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	lacpIfPol, err := getRemoteLACPMemberPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setLACPMemberPolicyAttributes(lacpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLACPMemberPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "lacpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
