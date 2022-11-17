package aci

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const leakInternalPrefixIp = "0.0.0.0/0"
const createLeakInternalPrefixMethod = "POST"
const cloudLeakToScope = "public"
const leakInternalPrefixLessThanOrEqual = "32"

func resourceAciLeakInternalPrefix() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLeakInternalPrefixCreate,
		UpdateContext: resourceAciLeakInternalPrefixUpdate,
		ReadContext:   resourceAciLeakInternalPrefixRead,
		DeleteContext: resourceAciLeakInternalPrefixDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeakInternalPrefixImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			// Source VRF DN
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vrf_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		})),
	}
}

func getRemoteLeakInternalPrefix(client *client.Client, dn string) (*models.LeakInternalPrefix, error) {
	leakInternalPrefixCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	leakInternalPrefix := models.LeakInternalPrefixFromContainer(leakInternalPrefixCont)
	if leakInternalPrefix.DistinguishedName == "" {
		return nil, fmt.Errorf("LeakInternalPrefix %s not found", dn)
	}
	return leakInternalPrefix, nil
}

func setLeakInternalPrefixAttributes(leakInternalPrefix *models.LeakInternalPrefix, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(leakInternalPrefix.DistinguishedName)
	d.Set("description", leakInternalPrefix.Description)

	leakInternalPrefixMap, err := leakInternalPrefix.ToMap()
	if err != nil {
		return d, err
	}
	if dn != leakInternalPrefix.DistinguishedName {
		d.Set("vrf_dn", "")
	} else {
		d.Set("vrf_dn", GetParentDn(dn, fmt.Sprintf("/%s/%s", models.RnleakRoutes, fmt.Sprintf(models.RnleakInternalPrefix, leakInternalPrefixMap["ip"]))))
	}
	d.Set("annotation", leakInternalPrefixMap["annotation"])
	d.Set("name_alias", leakInternalPrefixMap["nameAlias"])
	return d, nil
}

func getAndSetCloudLeakToObjects(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	leakToObjects, err := client.ListTenantandVRFdestinationforInterVRFLeakedRoutes(dn)
	if len(leakToObjects) >= 1 {
		newContent := make([]map[string]string, 0, 1)
		for _, leakToObject := range leakToObjects {
			newContent = append(newContent, map[string]string{
				"vrf_dn": leakToObject.ToCtxDn,
			})
		}
		d.Set("leak_to", newContent)
		return d, nil
	} else if err != nil {
		d.Set("leak_to", nil)
		log.Printf("[DEBUG]: Could not find existing leakTo objects under the parent DN: %s", dn)
	}
	return d, nil
}

func resourceAciLeakInternalPrefixImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	leakInternalPrefix, err := getRemoteLeakInternalPrefix(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLeakInternalPrefixAttributes(leakInternalPrefix, d)
	if err != nil {
		return nil, err
	}

	// leakTo - Beginning Import
	log.Printf("[DEBUG] %s: leakTo - Beginning Import with parent DN", dn)
	_, err = getAndSetCloudLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: leakTo - Import finished successfully with parent DN", dn)
	// leakTo - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeakInternalPrefixCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalPrefix: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VRFDn := d.Get("vrf_dn").(string)
	leakRoutesDn := fmt.Sprintf("%s/%s", VRFDn, models.RnleakRoutes)

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	annotation := ""
	if Annotation, ok := d.GetOk("annotation"); ok {
		annotation = Annotation.(string)
	} else {
		annotation = "orchestrator:terraform"
	}

	leakInternalPrefixMap := map[string]interface{}{
		"class_name": "leakInternalPrefix",
		"content": map[string]string{
			"annotation": annotation,
			"ip":         leakInternalPrefixIp,
			"le":         leakInternalPrefixLessThanOrEqual,
			"nameAlias":  nameAlias,
			"descr":      desc,
		},
	}

	leakToList := make([]interface{}, 0)

	if leakToSet, ok := d.GetOk("leak_to"); ok {
		leakToSetList := leakToSet.(*schema.Set).List()
		resleakToList, err := getCloudLeakToObjectsWithTenantAndCtxName(leakToSetList)
		if err != nil {
			return diag.FromErr(err)
		}
		leakToList = resleakToList
	}

	leakRoutesCont, err := createCloudLeakRoutesObject(leakInternalPrefixMap, leakToList)
	if err != nil {
		return diag.FromErr(err)
	}

	httpRequestPayload, err := aciClient.MakeRestRequest(createLeakInternalPrefixMethod, fmt.Sprintf("%s/%s.json", client.DefaultMOURL, leakRoutesDn), leakRoutesCont, true)
	if err != nil {
		return diag.FromErr(err)
	}

	respCont, _, err := aciClient.Do(httpRequestPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.CheckForErrors(respCont, createLeakInternalPrefixMethod, false)
	if err != nil {
		return diag.FromErr(err)
	}

	leakInternalPrefixDn := leakRoutesDn + "/" + fmt.Sprintf(models.RnleakInternalPrefix, leakInternalPrefixIp)
	d.SetId(leakInternalPrefixDn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLeakInternalPrefixRead(ctx, d, m)
}

func resourceAciLeakInternalPrefixUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalPrefix: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VRFDn := d.Get("vrf_dn").(string)
	leakRoutesDn := fmt.Sprintf("%s/%s", VRFDn, models.RnleakRoutes)

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	annotation := ""
	if Annotation, ok := d.GetOk("annotation"); ok {
		annotation = Annotation.(string)
	} else {
		annotation = "orchestrator:terraform"
	}

	leakInternalPrefixMap := map[string]interface{}{
		"class_name": "leakInternalPrefix",
		"content": map[string]string{
			"annotation": annotation,
			"ip":         leakInternalPrefixIp,
			"le":         leakInternalPrefixLessThanOrEqual,
			"nameAlias":  nameAlias,
			"descr":      desc,
		},
	}

	leakToList := make([]interface{}, 0, 1)

	if d.HasChange("leak_to") {
		oldSchemaObjs, newSchemaObjs := d.GetChange("leak_to")

		missingOldObjects := getOldObjectsNotInNew("vrf_dn", oldSchemaObjs.(*schema.Set), newSchemaObjs.(*schema.Set))

		newLeakToList := newSchemaObjs.(*schema.Set).List()

		for _, missingOldObject := range missingOldObjects {
			leakToMap := map[string]interface{}{
				"vrf_dn": missingOldObject.(map[string]interface{})["vrf_dn"],
				"status": "deleted",
			}
			newLeakToList = append(newLeakToList, leakToMap)
		}

		resLeakToList, err := getCloudLeakToObjectsWithTenantAndCtxName(newLeakToList)
		if err != nil {
			return diag.FromErr(err)
		}
		leakToList = resLeakToList
	}

	leakRoutesCont, err := createCloudLeakRoutesObject(leakInternalPrefixMap, leakToList)
	if err != nil {
		return diag.FromErr(err)
	}

	httpRequestPayload, err := aciClient.MakeRestRequest(createLeakInternalPrefixMethod, fmt.Sprintf("%s/%s.json", client.DefaultMOURL, leakRoutesDn), leakRoutesCont, true)
	if err != nil {
		return diag.FromErr(err)
	}

	respCont, _, err := aciClient.Do(httpRequestPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.CheckForErrors(respCont, createLeakInternalPrefixMethod, false)
	if err != nil {
		return diag.FromErr(err)
	}

	leakInternalPrefixDn := leakRoutesDn + "/" + fmt.Sprintf(models.RnleakInternalPrefix, leakInternalPrefixIp)
	d.SetId(leakInternalPrefixDn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLeakInternalPrefixRead(ctx, d, m)
}

func resourceAciLeakInternalPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	leakInternalPrefix, err := getRemoteLeakInternalPrefix(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setLeakInternalPrefixAttributes(leakInternalPrefix, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// leakTo - Beginning Read
	log.Printf("[DEBUG] %s: leakTo - Beginning Read with parent DN", dn)
	_, err = getAndSetCloudLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil
	}
	log.Printf("[DEBUG] %s: leakTo - Read finished successfully with parent DN", dn)
	// leakTo - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLeakInternalPrefixDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "leakInternalPrefix")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

func createCloudLeakRoutesObject(leakInternalPrefixMap map[string]interface{}, leakToList []interface{}) (*container.Container, error) {
	log.Printf("[DEBUG] Beginning of createCloudLeakRoutesObject")

	leakRoutesJson := []byte(`
		{
			"leakRoutes": {
				"attributes": {
				}
			}
		}`)

	leakRoutesCont, leakRoutesErr := container.ParseJSON(leakRoutesJson)
	if leakRoutesErr != nil {
		return nil, fmt.Errorf("leakRoutes ParseJSON Error: %v", leakRoutesErr)
	}

	leakInternalPrefixCont, leakInternalPrefixErr := preparePayload(leakInternalPrefixMap["class_name"].(string), leakInternalPrefixMap["content"].(map[string]string), leakToList)
	if leakInternalPrefixErr != nil {
		return nil, fmt.Errorf("leakInternalPrefix preparePayload Error: %v", leakInternalPrefixErr)
	}

	ParentObjectCreationErr := leakRoutesCont.ArrayAppend(leakInternalPrefixCont.Data(), "leakRoutes", "children")
	if ParentObjectCreationErr != nil {
		return nil, fmt.Errorf("ParentObjectCreation Error:: %v", ParentObjectCreationErr)
	}
	log.Printf("[DEBUG]: createCloudLeakRoutesObject finished successfully")
	return leakRoutesCont, nil
}

func getCloudLeakToObjectsWithTenantAndCtxName(leakToObjects []interface{}) ([]interface{}, error) {
	log.Printf("[DEBUG]: Beginning of getCloudLeakToObjectsWithTenantAndCtxName")
	resLeakToMap := make([]interface{}, 0)
	for _, leakToObjectInterface := range leakToObjects {
		leakToObjectMap := leakToObjectInterface.(map[string]interface{})
		if leakToObjectMap["vrf_dn"].(string) != "" {
			vrf_dn_pattern := regexp.MustCompile(`tn?-(.+).?/ctx?-([0-9A-Za-z_\-]+).?`)
			match_list := vrf_dn_pattern.FindStringSubmatch(leakToObjectMap["vrf_dn"].(string))
			if len(match_list) != 3 {
				return nil, fmt.Errorf("error occurred during the leak_to vrf_dn parsing: %s", leakToObjectMap["vrf_dn"])
			}

			leakToStatus := ""
			if leakToObjectMap["status"] != nil {
				leakToStatus = leakToObjectMap["status"].(string)
			}

			leakToMap := map[string]interface{}{
				"class_name": "leakTo",
				"content": map[string]string{
					"ctxName":    match_list[2],
					"scope":      cloudLeakToScope,
					"tenantName": match_list[1],
					"status":     leakToStatus,
				},
			}
			resLeakToMap = append(resLeakToMap, leakToMap)
		}
	}
	log.Printf("[DEBUG]: getCloudLeakToObjectsWithTenantAndCtxName finished successfully")
	return resLeakToMap, nil
}
