package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
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

		Schema: map[string]*schema.Schema{
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
							Default:  "",
						},
						"deployment_immediacy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"immediate",
								"lazy",
							}, false),
							Default: "lazy",
						},
						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"regular",
								"native",
								"untagged",
							}, false),
							Default: "regular",
						},
						"primary_encap": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "unknown",
						},
					},
				},
			},
		},
	}
}

func resourceAciBulkStaticPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BulkStaticPath: Beginning Creation")

	aciClient := m.(*client.Client)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)
	staticPathSet := staticPathPayload(d.Get("static_path").(*schema.Set).List(), "create")

	contentMap := make(map[string]interface{})
	contentMap["dn"] = ApplicationEPGDn
	cont, err := preparePayload(models.FvaepgClassName, toStrMap(contentMap), staticPathSet)
	if err != nil {
		return diag.FromErr(err)
	}

	_, diags := bulkStaticPathRequest(aciClient, "POST", ApplicationEPGDn, cont)
	if diags.HasError() {
		return diags
	}
	d.SetId(ApplicationEPGDn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBulkStaticPathRead(ctx, d, m)
}

func resourceAciBulkStaticPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bulkStaticPaths, diags := getAciBulkStaticPaths(aciClient, dn)
	if diags.HasError() {
		return diags
	}
	d = setAciBulkStaticPath(bulkStaticPaths, d)
	d.Set("application_epg_dn", dn)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBulkStaticPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())

	aciClient := m.(*client.Client)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	old_static_path, new_static_path := d.GetChange("static_path")
	del := getOldObjectsNotInNew("interface_dn", old_static_path.(*schema.Set), new_static_path.(*schema.Set))

	createPayload := staticPathPayload(new_static_path.(*schema.Set).List(), "create")
	deletePayload := staticPathPayload(del, "delete")
	payload := createPayload
	payload = append(payload, deletePayload...)

	contentMap := make(map[string]interface{})
	contentMap["dn"] = ApplicationEPGDn

	cont, err := preparePayload(models.FvaepgClassName, toStrMap(contentMap), payload)
	if err != nil {
		return diag.FromErr(err)
	}

	_, diags := bulkStaticPathRequest(aciClient, "POST", ApplicationEPGDn, cont)
	if diags.HasError() {
		return diags
	}
	d.SetId(ApplicationEPGDn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciBulkStaticPathRead(ctx, d, m)
}

func getAciBulkStaticPaths(aciClient *client.Client, dn string) ([]interface{}, diag.Diagnostics) {
	resp, diags := bulkStaticPathRequest(aciClient, "GET", dn, nil)
	if diags.HasError() {
		return nil, diags
	}
	staticPathSet := resp.S("imdata").S("fvRsPathAtt", "attributes")
	if staticPathSet == nil {
		return nil, nil
	}
	return staticPathSet.Data().([]interface{}), nil
}

func setAciBulkStaticPath(staticPathSet []interface{}, d *schema.ResourceData) *schema.ResourceData {
	stateStaticPathSet := make([]interface{}, 0, 1)
	for _, staticPath := range staticPathSet {
		staticPathMap := staticPath.(map[string]interface{})
		stateStaticPathMap := make(map[string]interface{})
		if descr, ok := staticPathMap["descr"]; ok {
			stateStaticPathMap["description"] = descr.(string)
		}
		if immediacy, ok := staticPathMap["instrImedcy"]; ok {
			stateStaticPathMap["deployment_immediacy"] = immediacy.(string)
		}
		if mode, ok := staticPathMap["mode"]; ok {
			stateStaticPathMap["mode"] = mode.(string)
		}
		if primaryEncap, ok := staticPathMap["primaryEncap"]; ok {
			stateStaticPathMap["primary_encap"] = primaryEncap.(string)
		}
		stateStaticPathMap["encap"] = staticPathMap["encap"].(string)
		stateStaticPathMap["interface_dn"] = staticPathMap["tDn"].(string)
		stateStaticPathSet = append(stateStaticPathSet, stateStaticPathMap)
	}
	d.Set("static_path", stateStaticPathSet)
	return d
}

func bulkStaticPathRequest(aciClient *client.Client, method string, epgDn string, body *container.Container) (*container.Container, diag.Diagnostics) {
	url := "/api/mo/" + epgDn + ".json"
	if method == "GET" {
		url += "?query-target=children&target-subtree-class=fvRsPathAtt"
	}
	req, err := aciClient.MakeRestRequest(method, url, body, true)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	respCont, _, err := aciClient.Do(req)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	err = client.CheckForErrors(respCont, method, false)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	if method == "POST" {
		return body, nil
	} else {
		return respCont, nil
	}
}

func staticPathPayload(staticPathList []interface{}, status string) []interface{} {
	staticPathSet := make([]interface{}, 0, 1)
	for _, staticPath := range staticPathList {
		staticPathMap := make(map[string]interface{})
		staticPathMap["class_name"] = "fvRsPathAtt"
		staticPathContent := make(map[string]interface{})
		staticPath := staticPath.(map[string]interface{})
		staticPathContent["encap"] = staticPath["encap"]
		staticPathContent["tDn"] = staticPath["interface_dn"]
		log.Printf("[DEBUG]: interface_dn in staticPathPayload %s", staticPath["interface_dn"])
		if descr, ok := staticPath["description"]; ok {
			staticPathContent["descr"] = descr
		}
		if immediacy, ok := staticPath["deployment_immediacy"]; ok {
			if immediacy != "" {
				staticPathContent["instrImedcy"] = immediacy
			}
		}
		if mode, ok := staticPath["mode"]; ok {
			if mode != "" {
				staticPathContent["mode"] = mode
			}
		}
		if primaryEncap, ok := staticPath["primary_encap"]; ok {
			staticPathContent["primaryEncap"] = primaryEncap
		}
		if status == "delete" {
			staticPathContent["status"] = "deleted"
		}
		log.Printf("[DEBUG]: staticPathContent in staticPathPayload %s", staticPathContent)
		staticPathMap["content"] = toStrMap(staticPathContent)
		staticPathSet = append(staticPathSet, staticPathMap)
	}
	return staticPathSet
}

func resourceAciBulkStaticPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	deletePayload := staticPathPayload(d.Get("static_path").(*schema.Set).List(), "delete")

	contentMap := make(map[string]interface{})
	contentMap["dn"] = ApplicationEPGDn
	cont, err := preparePayload(models.FvaepgClassName, toStrMap(contentMap), deletePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_, diags := bulkStaticPathRequest(aciClient, "POST", ApplicationEPGDn, cont)
	if diags.HasError() {
		return diags
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return nil
}

func resourceAciBulkStaticPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	d.Set("application_epg_dn", dn)

	bulkStaticPaths, diags := getAciBulkStaticPaths(aciClient, dn)
	if diags.HasError() {
		return nil, fmt.Errorf("could not read static path object when importing: %s", diags[0].Summary)
	}
	d = setAciBulkStaticPath(bulkStaticPaths, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}
