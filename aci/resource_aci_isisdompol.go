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

func resourceAciISISDomainPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciISISDomainPolicyCreate,
		UpdateContext: resourceAciISISDomainPolicyUpdate,
		ReadContext:   resourceAciISISDomainPolicyRead,
		DeleteContext: resourceAciISISDomainPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciISISDomainPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redistrib_metric": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsp_fast_flood": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"lsp_gen_init_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lsp_gen_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lsp_gen_sec_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isis_level_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isis_level_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "l1",
				ValidateFunc: validation.StringInSlice([]string{"l1", "l2"}, false),
			},
			"spf_comp_init_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spf_comp_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spf_comp_sec_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		)),
	}
}

func GetRemoteISISDomainPolicy(client *client.Client, dn string) (*models.ISISDomainPolicy, error) {
	isisDomPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	isisDomPol := models.ISISDomainPolicyFromContainer(isisDomPolCont)
	if isisDomPol.DistinguishedName == "" {
		return nil, fmt.Errorf("ISISDomainPolicy %s not found", isisDomPol.DistinguishedName)
	}
	return isisDomPol, nil
}

func GetRemoteISISLevel(client *client.Client, dn string) (*models.ISISLevel, error) {
	isisLvlCompCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	isisLvlComp := models.ISISLevelFromContainer(isisLvlCompCont)
	if isisLvlComp.DistinguishedName == "" {
		return nil, fmt.Errorf("ISISLevel %s not found", isisLvlComp.DistinguishedName)
	}
	return isisLvlComp, nil
}

func setISISDomainPolicyAttributes(isisDomPol *models.ISISDomainPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {

	d.SetId(isisDomPol.DistinguishedName)
	d.Set("description", isisDomPol.Description)
	isisDomPolMap, err := isisDomPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", isisDomPolMap["annotation"])
	d.Set("mtu", isisDomPolMap["mtu"])
	d.Set("redistrib_metric", isisDomPolMap["redistribMetric"])
	d.Set("name_alias", isisDomPolMap["nameAlias"])

	return d, nil
}

func setISISLevelAttributes(isisLvlComp *models.ISISLevel, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("description", isisLvlComp.Description)
	isisLvlCompMap, err := isisLvlComp.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", isisLvlCompMap["annotation"])
	d.Set("lsp_fast_flood", isisLvlCompMap["lspFastFlood"])
	d.Set("lsp_gen_init_intvl", isisLvlCompMap["lspGenInitIntvl"])
	d.Set("lsp_gen_max_intvl", isisLvlCompMap["lspGenMaxIntvl"])
	d.Set("lsp_gen_sec_intvl", isisLvlCompMap["lspGenSecIntvl"])
	d.Set("isis_level_name", isisLvlCompMap["name"])
	d.Set("spf_comp_init_intvl", isisLvlCompMap["spfCompInitIntvl"])
	d.Set("spf_comp_max_intvl", isisLvlCompMap["spfCompMaxIntvl"])
	d.Set("spf_comp_sec_intvl", isisLvlCompMap["spfCompSecIntvl"])
	d.Set("name_alias", isisLvlCompMap["nameAlias"])
	d.Set("isis_level_type", isisLvlCompMap["type"])
	return d, nil
}

func resourceAciISISDomainPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	isisDomPol, err := GetRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setISISDomainPolicyAttributes(isisDomPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciISISDomainPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ISISDomainPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	isisDomPolAttr := models.ISISDomainPolicyAttributes{}
	isisLvlCompAttr := models.ISISLevelAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		isisDomPolAttr.Annotation = Annotation.(string)
		isisLvlCompAttr.Annotation = Annotation.(string)
	} else {
		isisDomPolAttr.Annotation = "{}"
		isisLvlCompAttr.Annotation = "{}"
	}

	if Mtu, ok := d.GetOk("mtu"); ok {
		isisDomPolAttr.Mtu = Mtu.(string)
	}

	isisDomPolAttr.Name = "default"

	if RedistribMetric, ok := d.GetOk("redistrib_metric"); ok {
		isisDomPolAttr.RedistribMetric = RedistribMetric.(string)
	}
	isisDomPol := models.NewISISDomainPolicy(fmt.Sprintf("fabric/isisDomP-%s", isisDomPolAttr.Name), "uni", desc, nameAlias, isisDomPolAttr)
	isisDomPol.Status = "modified"

	err := aciClient.Save(isisDomPol)

	if err != nil {
		return diag.FromErr(err)
	}

	if LspFastFlood, ok := d.GetOk("lsp_fast_flood"); ok {
		isisLvlCompAttr.LspFastFlood = LspFastFlood.(string)
	}

	if LspGenInitIntvl, ok := d.GetOk("lsp_gen_init_intvl"); ok {
		isisLvlCompAttr.LspGenInitIntvl = LspGenInitIntvl.(string)
	}

	if LspGenMaxIntvl, ok := d.GetOk("lsp_gen_max_intvl"); ok {
		isisLvlCompAttr.LspGenMaxIntvl = LspGenMaxIntvl.(string)
	}

	if LspGenSecIntvl, ok := d.GetOk("lsp_gen_sec_intvl"); ok {
		isisLvlCompAttr.LspGenSecIntvl = LspGenSecIntvl.(string)
	}

	if Name, ok := d.GetOk("isis_level_name"); ok {
		isisLvlCompAttr.Name = Name.(string)
	}

	if SpfCompInitIntvl, ok := d.GetOk("spf_comp_init_intvl"); ok {
		isisLvlCompAttr.SpfCompInitIntvl = SpfCompInitIntvl.(string)
	}

	if SpfCompMaxIntvl, ok := d.GetOk("spf_comp_max_intvl"); ok {
		isisLvlCompAttr.SpfCompMaxIntvl = SpfCompMaxIntvl.(string)
	}

	if SpfCompSecIntvl, ok := d.GetOk("spf_comp_sec_intvl"); ok {
		isisLvlCompAttr.SpfCompSecIntvl = SpfCompSecIntvl.(string)
	}

	if IsisLevlType, ok := d.GetOk("isis_level_type"); ok {
		isisLvlCompAttr.ISISLevel_type = IsisLevlType.(string)
	}

	isisLvlComp := models.NewISISLevel(fmt.Sprintf("lvl-%s", isisLvlCompAttr.ISISLevel_type), isisDomPol.DistinguishedName, desc, nameAlias, isisLvlCompAttr)
	isisLvlComp.Status = "modified"
	err = aciClient.Save(isisLvlComp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(isisDomPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciISISDomainPolicyRead(ctx, d, m)
}

func resourceAciISISDomainPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ISISDomainPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	isisDomPolAttr := models.ISISDomainPolicyAttributes{}
	isisLvlCompAttr := models.ISISLevelAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		isisDomPolAttr.Annotation = Annotation.(string)
		isisLvlCompAttr.Annotation = Annotation.(string)
	} else {
		isisDomPolAttr.Annotation = "{}"
		isisLvlCompAttr.Annotation = "{}"
	}

	if Mtu, ok := d.GetOk("mtu"); ok {
		isisDomPolAttr.Mtu = Mtu.(string)
	}

	isisDomPolAttr.Name = "default"

	if RedistribMetric, ok := d.GetOk("redistrib_metric"); ok {
		isisDomPolAttr.RedistribMetric = RedistribMetric.(string)
	}
	isisDomPol := models.NewISISDomainPolicy(fmt.Sprintf("fabric/isisDomP-%s", isisDomPolAttr.Name), "uni", desc, nameAlias, isisDomPolAttr)
	isisDomPol.Status = "modified"
	err := aciClient.Save(isisDomPol)
	if err != nil {
		return diag.FromErr(err)
	}

	if LspFastFlood, ok := d.GetOk("lsp_fast_flood"); ok {
		isisLvlCompAttr.LspFastFlood = LspFastFlood.(string)
	}

	if LspGenInitIntvl, ok := d.GetOk("lsp_gen_init_intvl"); ok {
		isisLvlCompAttr.LspGenInitIntvl = LspGenInitIntvl.(string)
	}

	if LspGenMaxIntvl, ok := d.GetOk("lsp_gen_max_intvl"); ok {
		isisLvlCompAttr.LspGenMaxIntvl = LspGenMaxIntvl.(string)
	}

	if LspGenSecIntvl, ok := d.GetOk("lsp_gen_sec_intvl"); ok {
		isisLvlCompAttr.LspGenSecIntvl = LspGenSecIntvl.(string)
	}

	if Name, ok := d.GetOk("isis_level_name"); ok {
		isisLvlCompAttr.Name = Name.(string)
	}

	if SpfCompInitIntvl, ok := d.GetOk("spf_comp_init_intvl"); ok {
		isisLvlCompAttr.SpfCompInitIntvl = SpfCompInitIntvl.(string)
	}

	if SpfCompMaxIntvl, ok := d.GetOk("spf_comp_max_intvl"); ok {
		isisLvlCompAttr.SpfCompMaxIntvl = SpfCompMaxIntvl.(string)
	}

	if SpfCompSecIntvl, ok := d.GetOk("spf_comp_sec_intvl"); ok {
		isisLvlCompAttr.SpfCompSecIntvl = SpfCompSecIntvl.(string)
	}

	if IsisLevlType, ok := d.GetOk("isis_level_type"); ok {
		isisLvlCompAttr.ISISLevel_type = IsisLevlType.(string)
	}

	isisLvlComp := models.NewISISLevel(fmt.Sprintf("lvl-%s", isisLvlCompAttr.ISISLevel_type), isisDomPol.DistinguishedName, desc, nameAlias, isisLvlCompAttr)

	isisLvlComp.Status = "modified"
	err = aciClient.Save(isisLvlComp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(isisDomPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciISISDomainPolicyRead(ctx, d, m)
}

func resourceAciISISDomainPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	isisDomPol, err := GetRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setISISDomainPolicyAttributes(isisDomPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	isisCompDn := dn + "/lvl-l1"
	if d.Get("isis_level_type").(string) == "l2" {
		isisCompDn = dn + "/lvl-l2"
	}
	isisLvlComp, err := GetRemoteISISLevel(aciClient, isisCompDn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setISISLevelAttributes(isisLvlComp, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciISISDomainPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name isisDomPol and isisLvlComp cannot be deleted",
	})
	d.SetId("")
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
