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

func resourceAciBgpBestPathPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBgpBestPathPolicyCreate,
		UpdateContext: resourceAciBgpBestPathPolicyUpdate,
		ReadContext:   resourceAciBgpBestPathPolicyRead,
		DeleteContext: resourceAciBgpBestPathPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBgpBestPathPolicyImport,
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

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asPathMultipathRelax", "0",
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
func getRemoteBgpBestPathPolicy(client *client.Client, dn string) (*models.BgpBestPathPolicy, error) {
	bgpBestPathCtrlPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpBestPathCtrlPol := models.BgpBestPathPolicyFromContainer(bgpBestPathCtrlPolCont)

	if bgpBestPathCtrlPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BgpBestPathPolicy %s not found", bgpBestPathCtrlPol.DistinguishedName)
	}

	return bgpBestPathCtrlPol, nil
}

func setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol *models.BgpBestPathPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	d.Set("description", bgpBestPathCtrlPol.Description)
	dn := d.Id()
	if dn != bgpBestPathCtrlPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bgpBestPathCtrlPolMap, err := bgpBestPathCtrlPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/bestpath-%s", bgpBestPathCtrlPolMap["name"])))
	d.Set("name", bgpBestPathCtrlPolMap["name"])

	d.Set("annotation", bgpBestPathCtrlPolMap["annotation"])
	if bgpBestPathCtrlPolMap["ctrl"] == "" {
		d.Set("ctrl", "0")
	} else {
		d.Set("ctrl", bgpBestPathCtrlPolMap["ctrl"])
	}
	d.Set("name_alias", bgpBestPathCtrlPolMap["nameAlias"])
	return d, nil
}

func resourceAciBgpBestPathPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpBestPathCtrlPol, err := getRemoteBgpBestPathPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBgpBestPathPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpBestPathPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpBestPathCtrlPolAttr := models.BgpBestPathPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpBestPathCtrlPolAttr.Annotation = Annotation.(string)
	} else {
		bgpBestPathCtrlPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpBestPathCtrlPolAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpBestPathCtrlPolAttr.NameAlias = NameAlias.(string)
	}
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(fmt.Sprintf("bestpath-%s", name), TenantDn, desc, bgpBestPathCtrlPolAttr)

	err := aciClient.Save(bgpBestPathCtrlPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBgpBestPathPolicyRead(ctx, d, m)
}

func resourceAciBgpBestPathPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpBestPathPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpBestPathCtrlPolAttr := models.BgpBestPathPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpBestPathCtrlPolAttr.Annotation = Annotation.(string)
	} else {
		bgpBestPathCtrlPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpBestPathCtrlPolAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpBestPathCtrlPolAttr.NameAlias = NameAlias.(string)
	}
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(fmt.Sprintf("bestpath-%s", name), TenantDn, desc, bgpBestPathCtrlPolAttr)

	bgpBestPathCtrlPol.Status = "modified"

	err := aciClient.Save(bgpBestPathCtrlPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBgpBestPathPolicyRead(ctx, d, m)

}

func resourceAciBgpBestPathPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpBestPathCtrlPol, err := getRemoteBgpBestPathPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBgpBestPathPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpBestPathCtrlPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
