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

func resourceAciBulkStaticPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBulkStaticPathCreate,
		UpdateContext: resourceAciBulkStaticPathUpdate,
		ReadContext:   resourceAciBulkStaticPathRead,
		DeleteContext: resourceAciBulkStaticPathDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBulkStaticPathImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"static_path": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"interface_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"encap": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"deployment_immediacy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"immediate",
								"lazy",
							}, false),
						},
						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"regular",
								"native",
								"untagged",
							}, false),
						},
						"primary_encap": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

// func resourceAciBulkStaticPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
// 	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
// 	aciClient := m.(*client.Client)

// 	dn := d.Id()

// 	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

// 	if err != nil {
// 		return nil, err
// 	}
// 	fvRsPathAttMap, err := fvRsPathAtt.ToMap()

// 	if err != nil {
// 		return nil, err
// 	}

// 	tDn := fvRsPathAttMap["tDn"]
// 	pDN := GetParentDn(dn, fmt.Sprintf("/rspathAtt-[%s]", tDn))
// 	d.Set("application_epg_dn", pDN)
// 	schemaFilled, err := setStaticPathAttributes(fvRsPathAtt, d)

// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

// 	return []*schema.ResourceData{schemaFilled}, nil
// }

// func setStaticBulkPathAttributes(fvRsPathAtt *models.StaticPath, d *schema.ResourceData) (*schema.ResourceData, error) {
// 	d.SetId(fvRsPathAtt.DistinguishedName)
// 	fvRsPathAttMap, err := fvRsPathAtt.ToMap()

// 	if err != nil {
// 		return d, err
// 	}
// 	d.Set("application_epg_dn", GetParentDn(fvRsPathAtt.DistinguishedName, fmt.Sprintf("/rspathAtt-[%s]", fvRsPathAttMap["tDn"])))

// 	d.Set("tdn", fvRsPathAttMap["tDn"])

// 	d.Set("annotation", fvRsPathAttMap["annotation"])
// 	d.Set("encap", fvRsPathAttMap["encap"])
// 	d.Set("instr_imedcy", fvRsPathAttMap["instrImedcy"])
// 	d.Set("mode", fvRsPathAttMap["mode"])
// 	d.Set("primary_encap", fvRsPathAttMap["primaryEncap"])
// 	return d, nil
// }

func resourceAciBulkStaticPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] StaticPath: Beginning Creation")
	aciClient := m.(*client.Client)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsPathAttAttr := models.StaticPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsPathAttAttr.Annotation = Annotation.(string)
	} else {
		fvRsPathAttAttr.Annotation = "{}"
	}

	if bulk_static_path_config, ok := d.GetOk("bulk_static_path_config"); ok {
		staticPaths := bulk_static_path_config.(*schema.Set).List()
		for _, val := range staticPaths {
			staticPathAtt := models.StaticPathAttributes{}
			staticPath := val.(map[string]interface{})

			if staticPath["encap"] != nil {
				staticPathAtt.Encap = staticPath["encap"].(string)
			}
			if staticPath["instr_imedcy"] != nil {
				staticPathAtt.InstrImedcy = staticPath["instr_imedcy"].(string)
			}
			if staticPath["mode"] != nil {
				staticPathAtt.Mode = staticPath["mode"].(string)
			}
			if staticPath["primary_encap"] != nil {
				staticPathAtt.PrimaryEncap = staticPath["primary_encap"].(string)
			}
			fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", staticPath["tDn"]), ApplicationEPGDn, desc, staticPathAtt)
			err := aciClient.Save(fvRsPathAtt)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/rspathAtt", ApplicationEPGDn))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBulkStaticPathRead(ctx, d, m)
}

func resourceAciBulkStaticPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setStaticPathAttributes(fvRsPathAtt, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

// func resourceAciBulkStaticPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	log.Printf("[DEBUG] StaticPath: Beginning Creation")
// 	aciClient := m.(*client.Client)
// 	desc := d.Get("description").(string)

// 	ApplicationEPGDn := d.Get("application_epg_dn").(string)

// 	fvRsPathAttAttr := models.StaticPathAttributes{}
// 	if Annotation, ok := d.GetOk("annotation"); ok {
// 		fvRsPathAttAttr.Annotation = Annotation.(string)
// 	} else {
// 		fvRsPathAttAttr.Annotation = "{}"
// 	}

// 	if bulk_static_path_config, ok := d.GetOk("bulk_static_path_config"); ok {
// 		staticPaths := bulk_static_path_config.(*schema.Set).List()
// 		for _, val := range staticPaths {
// 			staticPathAtt := models.StaticPathAttributes{}
// 			staticPath := val.(map[string]interface{})

// 			if staticPath["encap"] != nil {
// 				staticPathAtt.Encap = staticPath["encap"].(string)
// 			}
// 			if staticPath["instr_imedcy"] != nil {
// 				staticPathAtt.InstrImedcy = staticPath["instr_imedcy"].(string)
// 			}
// 			if staticPath["mode"] != nil {
// 				staticPathAtt.Mode = staticPath["mode"].(string)
// 			}
// 			if staticPath["primary_encap"] != nil {
// 				staticPathAtt.PrimaryEncap = staticPath["primary_encap"].(string)
// 			}
// 			fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", staticPath["tDn"]), ApplicationEPGDn, desc, staticPathAtt)
// 			err := aciClient.Save(fvRsPathAtt)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}
// 		}
// 	}

// 	d.SetId(fmt.Sprintf("%s/rspathAtt", ApplicationEPGDn))
// 	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

// 	return resourceAciBulkStaticPathRead(ctx, d, m)
// }

// func resourceAciBulkStaticPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

// 	aciClient := m.(*client.Client)
// 	dn := d.Id()
// 	err := aciClient.DeleteByDn(dn, "fvRsPathAtt")
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

// 	d.SetId("")
// 	return diag.FromErr(err)
// }
