package aci

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciHSRPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciHSRPInterfacePolicyCreate,
		UpdateContext: resourceAciHSRPInterfacePolicyUpdate,
		ReadContext:   resourceAciHSRPInterfacePolicyRead,
		DeleteContext: resourceAciHSRPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciHSRPInterfacePolicyImport,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"bfd",
						"bia",
					}, false),
				},
			},

			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"reload_delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteHSRPInterfacePolicy(client *client.Client, dn string) (*models.HSRPInterfacePolicy, error) {
	hsrpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpIfPol := models.HSRPInterfacePolicyFromContainer(hsrpIfPolCont)

	if hsrpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("HSRPInterfacePolicy %s not found", hsrpIfPol.DistinguishedName)
	}

	return hsrpIfPol, nil
}

func setHSRPInterfacePolicyAttributes(hsrpIfPol *models.HSRPInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(hsrpIfPol.DistinguishedName)
	d.Set("description", hsrpIfPol.Description)
	if dn != hsrpIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	hsrpIfPolMap, err := hsrpIfPol.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/hsrpIfPol-%s", hsrpIfPolMap["name"])))
	d.Set("name", hsrpIfPolMap["name"])

	d.Set("annotation", hsrpIfPolMap["annotation"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(hsrpIfPolMap["ctrl"], ",") {
		ctrlGet = append(ctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ctrlGet)
	if len(ctrlGet) == 1 && ctrlGet[0] == "" {
		d.Set("ctrl", make([]string, 0, 1))
	} else {
		d.Set("ctrl", ctrlGet)
	}
	d.Set("delay", hsrpIfPolMap["delay"])
	d.Set("name_alias", hsrpIfPolMap["nameAlias"])
	d.Set("reload_delay", hsrpIfPolMap["reloadDelay"])

	return d, nil
}

func resourceAciHSRPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpIfPol, err := getRemoteHSRPInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setHSRPInterfacePolicyAttributes(hsrpIfPol, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciHSRPInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] HSRPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpIfPolAttr := models.HSRPInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPolAttr.Annotation = "{}"
	}
	if ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		ctrl := strings.Join(ctrlList, ",")
		hsrpIfPolAttr.Ctrl = ctrl
	}
	if Delay, ok := d.GetOk("delay"); ok {
		hsrpIfPolAttr.Delay = Delay.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPolAttr.NameAlias = NameAlias.(string)
	}
	if ReloadDelay, ok := d.GetOk("reload_delay"); ok {
		hsrpIfPolAttr.ReloadDelay = ReloadDelay.(string)
	}
	hsrpIfPol := models.NewHSRPInterfacePolicy(fmt.Sprintf("hsrpIfPol-%s", name), TenantDn, desc, hsrpIfPolAttr)

	err := aciClient.Save(hsrpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hsrpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciHSRPInterfacePolicyRead(ctx, d, m)
}

func resourceAciHSRPInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] HSRPInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpIfPolAttr := models.HSRPInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPolAttr.Annotation = "{}"
	}
	if ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		ctrl := strings.Join(ctrlList, ",")
		hsrpIfPolAttr.Ctrl = ctrl
	}
	if Delay, ok := d.GetOk("delay"); ok {
		hsrpIfPolAttr.Delay = Delay.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPolAttr.NameAlias = NameAlias.(string)
	}
	if ReloadDelay, ok := d.GetOk("reload_delay"); ok {
		hsrpIfPolAttr.ReloadDelay = ReloadDelay.(string)
	}
	hsrpIfPol := models.NewHSRPInterfacePolicy(fmt.Sprintf("hsrpIfPol-%s", name), TenantDn, desc, hsrpIfPolAttr)

	hsrpIfPol.Status = "modified"

	err := aciClient.Save(hsrpIfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hsrpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciHSRPInterfacePolicyRead(ctx, d, m)

}

func resourceAciHSRPInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpIfPol, err := getRemoteHSRPInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setHSRPInterfacePolicyAttributes(hsrpIfPol, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciHSRPInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
