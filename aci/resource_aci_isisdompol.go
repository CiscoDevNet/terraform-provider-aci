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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redistrib_metric": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isis_level": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(
						map[string]*schema.Schema{
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
							"name": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
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
							"isis_level_type": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.StringInSlice([]string{
									"l1",
									"l2",
								}, false),
							},
						},
					)),
				},
			},
			"isis_level_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		})),
	}
}

func getRemoteISISDomainPolicy(client *client.Client, dn string) (*models.ISISDomainPolicy, error) {
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

func getRemoteISISLevel(client *client.Client, dn string) (*models.ISISLevel, error) {
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
	d.Set("name", isisDomPolMap["name"])
	d.Set("redistrib_metric", isisDomPolMap["redistribMetric"])
	d.Set("name_alias", isisDomPolMap["nameAlias"])
	return d, nil
}
func setISISLevelAttributes(isisLvlComps []*models.ISISLevel, d *schema.ResourceData) (*schema.ResourceData, error) {
	ISISLvlSet := make([]interface{}, 0, 1)

	for _, ISISLvlComp := range isisLvlComps {
		lvl := make(map[string]interface{})
		lvl["id"] = ISISLvlComp.DistinguishedName
		lvlMap, err := ISISLvlComp.ToMap()
		if err != nil {
			return d, err
		}
		lvl["name"] = lvlMap["name"]
		lvl["annotation"] = lvlMap["annotation"]
		lvl["name_alias"] = lvlMap["nameAlias"]
		lvl["lsp_fast_flood"] = lvlMap["lspFastFlood"]
		lvl["lsp_gen_init_intvl"] = lvlMap["lspGenInitIntvl"]
		lvl["lsp_gen_max_intvl"] = lvlMap["lspGenMaxIntvl"]
		lvl["spf_comp_init_intvl"] = lvlMap["spfCompInitIntvl"]
		lvl["spf_comp_max_intvl"] = lvlMap["spfCompMaxIntvl"]
		lvl["spf_comp_sec_intvl"] = lvlMap["spfCompSecIntvl"]
		lvl["isis_level_type"] = lvlMap["type"]
		lvl["description"] = lvlMap["description"]
		lvl["lsp_gen_sec_intvl"] = lvlMap["lspGenSecIntvl"]
		ISISLvlSet = append(ISISLvlSet, lvl)
	}

	d.Set("isis_level", ISISLvlSet)
	return d, nil
}

func resourceAciISISDomainPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	isisDomPol, err := getRemoteISISDomainPolicy(aciClient, dn)
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
	name := d.Get("name").(string)
	isisDomPolAttr := models.ISISDomainPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		isisDomPolAttr.Annotation = Annotation.(string)
	} else {
		isisDomPolAttr.Annotation = "{}"
	}

	if Mtu, ok := d.GetOk("mtu"); ok {
		isisDomPolAttr.Mtu = Mtu.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		isisDomPolAttr.Name = Name.(string)
	}

	if RedistribMetric, ok := d.GetOk("redistrib_metric"); ok {
		isisDomPolAttr.RedistribMetric = RedistribMetric.(string)
	}
	isisDomPol := models.NewISISDomainPolicy(fmt.Sprintf("fabric/isisDomP-%s", name), "uni", desc, nameAlias, isisDomPolAttr)
	err := aciClient.Save(isisDomPol)
	if err != nil {
		return diag.FromErr(err)
	}

	lvlCompIDS := make([]string, 0, 1)
	if levls, ok := d.GetOk("isis_level"); ok {
		levelsList := levls.([]interface{})
		for _, val := range levelsList {
			isisLvlAttr := models.ISISLevelAttributes{}
			isisLvl := val.(map[string]interface{})
			desc := ""
			nameAlias := ""
			isisLvlDn := isisDomPol.DistinguishedName
			if isisLvl["lsp_fast_flood"] != nil {
				isisLvlAttr.LspFastFlood = isisLvl["lsp_fast_flood"].(string)
			}
			if isisLvl["lsp_gen_init_intvl"] != nil {
				isisLvlAttr.LspGenInitIntvl = isisLvl["lsp_gen_init_intvl"].(string)
			}
			if isisLvl["lsp_gen_max_intvl"] != nil {
				isisLvlAttr.LspGenMaxIntvl = isisLvl["lsp_gen_max_intvl"].(string)
			}
			if isisLvl["lsp_gen_sec_intvl"] != nil {
				isisLvlAttr.LspGenSecIntvl = isisLvl["lsp_gen_max_intvl"].(string)
			}
			if isisLvl["name"] != nil {
				isisLvlAttr.Name = isisLvl["name"].(string)
			}
			if isisLvl["spf_comp_init_intvl"] != nil {
				isisLvlAttr.SpfCompInitIntvl = isisLvl["spf_comp_init_intvl"].(string)
			}
			if isisLvl["spf_comp_max_intvl"] != nil {
				isisLvlAttr.SpfCompMaxIntvl = isisLvl["spf_comp_max_intvl"].(string)
			}
			if isisLvl["spf_comp_sec_intvl"] != nil {
				isisLvlAttr.SpfCompSecIntvl = isisLvl["spf_comp_sec_intvl"].(string)
			}
			if isisLvl["isis_level_type"] != nil {
				isisLvlAttr.ISISLevel_type = isisLvl["isis_level_type"].(string)
			}
			if isisLvl["annotation"] != nil {
				isisLvlAttr.Annotation = isisLvl["annotation"].(string)
			} else {
				isisLvlAttr.Annotation = "{}"
			}
			if isisLvl["description"] != nil {
				desc = isisLvl["description"].(string)
			}
			if isisLvl["name_alias"] != nil {
				nameAlias = isisLvl["name_alias"].(string)
			}
			
		}
		d.Set("isis_level_ids", lvlCompIDS)
	} else {
		d.Set("isis_level_ids", lvlCompIDS)
	}

	d.SetId(isisDomPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciISISDomainPolicyRead(ctx, d, m)
}

func resourceAciISISDomainPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ISISDomainPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	isisDomPolAttr := models.ISISDomainPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		isisDomPolAttr.Annotation = Annotation.(string)
	} else {
		isisDomPolAttr.Annotation = "{}"
	}

	if Mtu, ok := d.GetOk("mtu"); ok {
		isisDomPolAttr.Mtu = Mtu.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		isisDomPolAttr.Name = Name.(string)
	}

	if RedistribMetric, ok := d.GetOk("redistrib_metric"); ok {
		isisDomPolAttr.RedistribMetric = RedistribMetric.(string)
	}
	isisDomPol := models.NewISISDomainPolicy(fmt.Sprintf("fabric/isisDomP-%s", name), "uni", desc, nameAlias, isisDomPolAttr)
	isisDomPol.Status = "modified"
	err := aciClient.Save(isisDomPol)
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
	isisDomPol, err := getRemoteISISDomainPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setISISDomainPolicyAttributes(isisDomPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciISISDomainPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "isisDomPol")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
