package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciCloudTemplateforExternalNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudTemplateforExternalNetworkRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"hub_network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vrf_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"regions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"cloud_vendor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aws",
					"azure",
					"gcp",
				}, false),
			},
			"router_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"c8kv",
					"tgw",
				}, false),
			},
		})),
	}
}

func dataSourceAciCloudTemplateforExternalNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	CloudInfraNetworkTemplateDn := "uni/tn-infra/infranetwork-default"
	rn := fmt.Sprintf(models.RncloudtemplateExtNetwork, name)
	dn := fmt.Sprintf("%s/%s", CloudInfraNetworkTemplateDn, rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	cloudtemplateExtNetwork, err := getRemoteCloudTemplateforExternalNetwork(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setCloudTemplateforExternalNetworkAttributes(cloudtemplateExtNetwork, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] : Data Source -  Begining Read of cloud Regions attributes.")
	regionsData, err := aciClient.ListCloudProviderandRegionNames(cloudtemplateExtNetwork.DistinguishedName)
	if err != nil {
		log.Printf("[DEBUG] : Data Source - Error while reading cloud Regions attributes %v", err)
	}

	RegionsList := make([]string, 0, 1)
	for _, regionValue := range regionsData {

		regionsMap, err := setCloudProviderandRegionNamesAttributes(regionValue, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil
		}
		RegionsList = append(RegionsList, regionsMap["region"])
		d.Set("cloud_vendor", regionsMap["cloud_vendor"])
		if regionsMap["cloud_vendor"] != "aws" {
			d.Set("router_type", "")
		}
	}
	log.Printf("[DEBUG] : Data Source -  Read cloud regions finished successfully")
	d.Set("regions", RegionsList)

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
