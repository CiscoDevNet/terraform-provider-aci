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

func resourceAciTabooContract() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTabooContractCreate,
		UpdateContext: resourceAciTabooContractUpdate,
		ReadContext:   resourceAciTabooContractRead,
		DeleteContext: resourceAciTabooContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTabooContractImport,
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
func getRemoteTabooContract(client *client.Client, dn string) (*models.TabooContract, error) {
	vzTabooCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzTaboo := models.TabooContractFromContainer(vzTabooCont)

	if vzTaboo.DistinguishedName == "" {
		return nil, fmt.Errorf("TabooContract %s not found", vzTaboo.DistinguishedName)
	}

	return vzTaboo, nil
}

func setTabooContractAttributes(vzTaboo *models.TabooContract, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzTaboo.DistinguishedName)
	d.Set("description", vzTaboo.Description)
	if dn != vzTaboo.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzTabooMap, err := vzTaboo.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", vzTabooMap["name"])
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("taboo-%s", vzTabooMap["name"])))
	d.Set("annotation", vzTabooMap["annotation"])
	d.Set("name_alias", vzTabooMap["nameAlias"])
	return d, nil
}

func resourceAciTabooContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzTaboo, err := getRemoteTabooContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vzTabooMap, err := vzTaboo.ToMap()

	if err != nil {
		return nil, err
	}

	name := vzTabooMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/taboo-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setTabooContractAttributes(vzTaboo, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTabooContractCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TabooContract: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzTabooAttr := models.TabooContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzTabooAttr.Annotation = Annotation.(string)
	} else {
		vzTabooAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzTabooAttr.NameAlias = NameAlias.(string)
	}
	vzTaboo := models.NewTabooContract(fmt.Sprintf("taboo-%s", name), TenantDn, desc, vzTabooAttr)

	err := aciClient.Save(vzTaboo)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzTaboo.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciTabooContractRead(ctx, d, m)
}

func resourceAciTabooContractUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TabooContract: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzTabooAttr := models.TabooContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzTabooAttr.Annotation = Annotation.(string)
	} else {
		vzTabooAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzTabooAttr.NameAlias = NameAlias.(string)
	}
	vzTaboo := models.NewTabooContract(fmt.Sprintf("taboo-%s", name), TenantDn, desc, vzTabooAttr)

	vzTaboo.Status = "modified"

	err := aciClient.Save(vzTaboo)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzTaboo.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciTabooContractRead(ctx, d, m)

}

func resourceAciTabooContractRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzTaboo, err := getRemoteTabooContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setTabooContractAttributes(vzTaboo, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciTabooContractDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzTaboo")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
