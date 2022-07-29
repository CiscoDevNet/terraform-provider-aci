package aci

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciLeakInternalSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLeakInternalSubnetCreate,
		UpdateContext: resourceAciLeakInternalSubnetUpdate,
		ReadContext:   resourceAciLeakInternalSubnetRead,
		DeleteContext: resourceAciLeakInternalSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeakInternalSubnetImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// True -> public, False -> private, Default -> false(private)
			"allow_l3out_advertisement": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"false", // False -> private
					"true",  // True -> public
				}, false),
				Default: "false",
			},
			// Tenant and VRF Destinations
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vrf_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						// True -> public, False -> private, Default -> "inherit"
						"allow_l3out_advertisement": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"false", // False -> private
								"true",  // True -> public
								"inherit",
							}, false),
							Default: "inherit",
						},
					},
				},
			},
		})),
	}
}

func getRemoteLeakInternalSubnet(client *client.Client, dn string) (*models.LeakInternalSubnet, error) {
	leakInternalSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	leakInternalSubnet := models.LeakInternalSubnetFromContainer(leakInternalSubnetCont)
	if leakInternalSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("LeakInternalSubnet %s not found", dn)
	}
	return leakInternalSubnet, nil
}

func setLeakInternalSubnetAttributes(leakInternalSubnet *models.LeakInternalSubnet, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(leakInternalSubnet.DistinguishedName)
	d.Set("description", leakInternalSubnet.Description)

	leakInternalSubnetMap, err := leakInternalSubnet.ToMap()
	if err != nil {
		return d, err
	}
	if dn != leakInternalSubnet.DistinguishedName {
		d.Set("vrf_dn", "")
	} else {
		d.Set("vrf_dn", GetParentDn(dn, fmt.Sprintf("/%s/%s", models.RnleakRoutes, fmt.Sprintf(models.RnleakInternalSubnet, leakInternalSubnetMap["ip"]))))
	}
	d.Set("annotation", leakInternalSubnetMap["annotation"])
	d.Set("ip", leakInternalSubnetMap["ip"])

	d.Set("allow_l3out_advertisement", "false")
	if leakInternalSubnetMap["scope"] == "public" {
		d.Set("allow_l3out_advertisement", "true")
	}

	d.Set("name_alias", leakInternalSubnetMap["nameAlias"])
	return d, nil
}

func getAndSetLeakToObjects(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	leakToObjects, err := client.ListTenantandVRFdestinationforInterVRFLeakedRoutes(dn)
	if len(leakToObjects) >= 1 {
		newContent := make([]map[string]string, 0, 1)
		for _, leakToObject := range leakToObjects {
			leakToAllowL3OutAdvertisement := "inherit"
			if leakToObject.Scope == "public" {
				leakToAllowL3OutAdvertisement = "true"
			} else if leakToObject.Scope == "private" {
				leakToAllowL3OutAdvertisement = "false"
			}
			newContent = append(newContent, map[string]string{
				"vrf_dn":                    leakToObject.ToCtxDn,
				"allow_l3out_advertisement": leakToAllowL3OutAdvertisement,
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

func resourceAciLeakInternalSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	leakInternalSubnet, err := getRemoteLeakInternalSubnet(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLeakInternalSubnetAttributes(leakInternalSubnet, d)
	if err != nil {
		return nil, err
	}

	// leakTo - Beginning Import
	log.Printf("[DEBUG] %s: leakTo - Beginning Import with parent DN", dn)
	_, err = getAndSetLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: leakTo - Import finished successfully with parent DN", dn)
	// leakTo - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeakInternalSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalSubnet: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)
	leakRoutesDn := VRFDn + "/" + models.RnleakRoutes

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

	leakInternalAllowL3OutAdvertisement := d.Get("allow_l3out_advertisement").(string)

	// Default "Allow L3Out Advertisement" value is false(private)
	leakInternalSubnetScope := "private"
	if leakInternalAllowL3OutAdvertisement == "true" {
		leakInternalSubnetScope = "public"
	}

	leakInternalSubnetMap := map[string]interface{}{
		"class_name": "leakInternalSubnet",
		"content": map[string]string{
			"annotation": annotation,
			"ip":         ip,
			"scope":      leakInternalSubnetScope,
			"nameAlias":  nameAlias,
			"descr":      desc,
		},
	}

	leakToList := make([]interface{}, 0)

	if leakToSet, ok := d.GetOk("leak_to"); ok {
		leakToSetList := leakToSet.(*schema.Set).List()
		resleakToList, err := getLeakToObjectsWithTenantAndCtxName(leakToSetList)
		if err != nil {
			return diag.FromErr(err)
		}
		leakToList = resleakToList
	}

	leakRoutesCont, err := createLeakRoutesObject(leakInternalSubnetMap, leakToList)
	if err != nil {
		return diag.FromErr(err)
	}

	httpRequestPayload, err := aciClient.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s.json", leakRoutesDn), leakRoutesCont, true)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = aciClient.Do(httpRequestPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	leakInternalSubnetDn := leakRoutesDn + "/" + fmt.Sprintf(models.RnleakInternalSubnet, ip)
	d.SetId(leakInternalSubnetDn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLeakInternalSubnetRead(ctx, d, m)
}

func resourceAciLeakInternalSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalSubnet: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)
	leakRoutesDn := VRFDn + "/" + models.RnleakRoutes
	leakInternalAllowL3OutAdvertisement := d.Get("allow_l3out_advertisement").(string)

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

	// Default "Allow L3Out Advertisement" value is false(private)
	leakInternalSubnetScope := "private"
	if leakInternalAllowL3OutAdvertisement == "true" {
		leakInternalSubnetScope = "public"
	}

	leakInternalSubnetMap := map[string]interface{}{
		"class_name": "leakInternalSubnet",
		"content": map[string]string{
			"annotation": annotation,
			"ip":         ip,
			"scope":      leakInternalSubnetScope,
			"nameAlias":  nameAlias,
			"descr":      desc,
			"status":     "modified",
		},
	}

	leakToList := make([]interface{}, 0, 1)

	if d.HasChange("leak_to") {
		oldSchemaObjs, newSchemaObjs := d.GetChange("leak_to")

		missingOldObjects := getOldObjectsNotInNew("vrf_dn", oldSchemaObjs.(*schema.Set), newSchemaObjs.(*schema.Set))

		newLeakToList := newSchemaObjs.(*schema.Set).List()

		for _, missingOldObject := range missingOldObjects {
			leakToMap := map[string]interface{}{
				"allow_l3out_advertisement": missingOldObject.(map[string]interface{})["allow_l3out_advertisement"],
				"vrf_dn":                    missingOldObject.(map[string]interface{})["vrf_dn"],
				"status":                    "deleted",
			}
			newLeakToList = append(newLeakToList, leakToMap)
		}

		resLeakToList, err := getLeakToObjectsWithTenantAndCtxName(newLeakToList)
		if err != nil {
			return diag.FromErr(err)
		}
		leakToList = resLeakToList
	}

	leakRoutesCont, err := createLeakRoutesObject(leakInternalSubnetMap, leakToList)
	if err != nil {
		return diag.FromErr(err)
	}

	httpRequestPayload, err := aciClient.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s.json", leakRoutesDn), leakRoutesCont, true)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = aciClient.Do(httpRequestPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	leakInternalSubnetDn := leakRoutesDn + "/" + fmt.Sprintf(models.RnleakInternalSubnet, ip)
	d.SetId(leakInternalSubnetDn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLeakInternalSubnetRead(ctx, d, m)
}

func resourceAciLeakInternalSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	leakInternalSubnet, err := getRemoteLeakInternalSubnet(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setLeakInternalSubnetAttributes(leakInternalSubnet, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// leakTo - Beginning Read
	log.Printf("[DEBUG] %s: leakTo - Beginning Read with parent DN", dn)
	_, err = getAndSetLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil
	}
	log.Printf("[DEBUG] %s: leakTo - Read finished successfully with parent DN", dn)
	// leakTo - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLeakInternalSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "leakInternalSubnet")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

func createLeakRoutesObject(leakInternalSubnetMap map[string]interface{}, leakToList []interface{}) (*container.Container, error) {
	log.Printf("[DEBUG] Beginning of createLeakRoutesObject")

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

	leakInternalSubnetCont, leakInternalSubnetErr := preparePayload(leakInternalSubnetMap["class_name"].(string), leakInternalSubnetMap["content"].(map[string]string), leakToList)
	if leakInternalSubnetErr != nil {
		return nil, fmt.Errorf("leakInternalSubnet preparePayload Error: %v", leakInternalSubnetErr)
	}

	ParentObjectCreationErr := leakRoutesCont.ArrayAppend(leakInternalSubnetCont.Data(), "leakRoutes", "children")
	if ParentObjectCreationErr != nil {
		return nil, fmt.Errorf("ParentObjectCreation Error:: %v", ParentObjectCreationErr)
	}
	log.Printf("[DEBUG]: createLeakRoutesObject finished successfully")
	return leakRoutesCont, nil
}

func getLeakToObjectsWithTenantAndCtxName(leakToObjects []interface{}) ([]interface{}, error) {
	log.Printf("[DEBUG]: Beginning of getLeakToObjectsWithTenantAndCtxName")

	resLeakToMap := make([]interface{}, 0)
	for _, leakToObjectInterface := range leakToObjects {
		leakToObjectMap := leakToObjectInterface.(map[string]interface{})
		if leakToObjectMap["vrf_dn"].(string) != "" {
			vrf_dn_pattern := regexp.MustCompile(`tn?-(.+).?/ctx?-([0-9A-Za-z_\-]+).?`)
			match_list := vrf_dn_pattern.FindStringSubmatch(leakToObjectMap["vrf_dn"].(string))
			if len(match_list) != 3 {
				return nil, fmt.Errorf("error occurred during the leak_to vrf_dn parsing: %s", leakToObjectMap["vrf_dn"])
			}
			leakToScope := "inherit"

			if leakToObjectMap["allow_l3out_advertisement"] == "true" {
				leakToScope = "public"
			} else if leakToObjectMap["allow_l3out_advertisement"] == "false" {
				leakToScope = "private"
			}

			leakToStatus := ""
			if leakToObjectMap["status"] != nil {
				leakToStatus = leakToObjectMap["status"].(string)
			}

			leakToMap := map[string]interface{}{
				"class_name": "leakTo",
				"content": map[string]string{
					"ctxName":    match_list[2],
					"scope":      leakToScope,
					"tenantName": match_list[1],
					"status":     leakToStatus,
				},
			}
			resLeakToMap = append(resLeakToMap, leakToMap)
		}
	}
	return resLeakToMap, nil
}
