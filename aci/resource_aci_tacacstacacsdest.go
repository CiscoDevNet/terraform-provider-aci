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

func resourceAciTACACSDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTACACSDestinationCreate,
		UpdateContext: resourceAciTACACSDestinationUpdate,
		ReadContext:   resourceAciTACACSDestinationRead,
		DeleteContext: resourceAciTACACSDestinationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTACACSDestinationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tacacs_accounting_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"chap",
					"mschap",
					"pap",
				}, false),
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "49",
			},

			"relation_file_rs_a_remote_host_to_epg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:ATg",
			},
			"relation_file_rs_a_remote_host_to_epp": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:AREpP",
			}})),
	}
}

func getRemoteTACACSDestination(client *client.Client, dn string) (*models.TACACSDestination, error) {
	tacacsTacacsDestCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsTacacsDest := models.TACACSDestinationFromContainer(tacacsTacacsDestCont)
	if tacacsTacacsDest.DistinguishedName == "" {
		return nil, fmt.Errorf("TACACSDestination %s not found", tacacsTacacsDest.DistinguishedName)
	}
	return tacacsTacacsDest, nil
}

func setTACACSDestinationAttributes(tacacsTacacsDest *models.TACACSDestination, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(tacacsTacacsDest.DistinguishedName)
	d.Set("description", tacacsTacacsDest.Description)
	tacacsTacacsDestMap, err := tacacsTacacsDest.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("tacacs_accounting_dn", GetParentDn(d.Id(), fmt.Sprintf("/tacacsdest-%s-port-%s", tacacsTacacsDestMap["host"], tacacsTacacsDestMap["port"])))
	d.Set("annotation", tacacsTacacsDestMap["annotation"])
	d.Set("auth_protocol", tacacsTacacsDestMap["authProtocol"])
	d.Set("host", tacacsTacacsDestMap["host"])
	d.Set("name", tacacsTacacsDestMap["name"])
	d.Set("port", tacacsTacacsDestMap["port"])
	d.Set("name_alias", tacacsTacacsDestMap["nameAlias"])
	return d, nil
}

func resourceAciTACACSDestinationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsTacacsDest, err := getRemoteTACACSDestination(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTACACSDestinationAttributes(tacacsTacacsDest, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTACACSDestinationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSDestination: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(string)
	TACACSMonitoringDestinationGroupDn := d.Get("tacacs_accounting_dn").(string)

	tacacsTacacsDestAttr := models.TACACSDestinationAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsTacacsDestAttr.Annotation = Annotation.(string)
	} else {
		tacacsTacacsDestAttr.Annotation = "{}"
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		tacacsTacacsDestAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Host, ok := d.GetOk("host"); ok {
		tacacsTacacsDestAttr.Host = Host.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		tacacsTacacsDestAttr.Key = Key.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsTacacsDestAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		tacacsTacacsDestAttr.Port = Port.(string)
	}
	tacacsTacacsDest := models.NewTACACSDestination(fmt.Sprintf("tacacsdest-%s-port-%s", host, port), TACACSMonitoringDestinationGroupDn, desc, nameAlias, tacacsTacacsDestAttr)

	err := aciClient.Save(tacacsTacacsDest)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofileRsARemoteHostToEpg, ok := d.GetOk("relation_file_rs_a_remote_host_to_epg"); ok {
		relationParam := relationTofileRsARemoteHostToEpg.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofileRsARemoteHostToEpp, ok := d.GetOk("relation_file_rs_a_remote_host_to_epp"); ok {
		relationParam := relationTofileRsARemoteHostToEpp.(string)
		checkDns = append(checkDns, relationParam)

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTofileRsARemoteHostToEpg, ok := d.GetOk("relation_file_rs_a_remote_host_to_epg"); ok {
		relationParam := relationTofileRsARemoteHostToEpg.(string)
		err = aciClient.CreateRelationfileRsARemoteHostToEpg(tacacsTacacsDest.DistinguishedName, tacacsTacacsDestAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofileRsARemoteHostToEpp, ok := d.GetOk("relation_file_rs_a_remote_host_to_epp"); ok {
		relationParam := relationTofileRsARemoteHostToEpp.(string)
		err = aciClient.CreateRelationfileRsARemoteHostToEpp(tacacsTacacsDest.DistinguishedName, tacacsTacacsDestAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(tacacsTacacsDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTACACSDestinationRead(ctx, d, m)
}

func resourceAciTACACSDestinationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSDestination: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(string)
	TACACSMonitoringDestinationGroupDn := d.Get("tacacs_accounting_dn").(string)
	tacacsTacacsDestAttr := models.TACACSDestinationAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsTacacsDestAttr.Annotation = Annotation.(string)
	} else {
		tacacsTacacsDestAttr.Annotation = "{}"
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		tacacsTacacsDestAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Host, ok := d.GetOk("host"); ok {
		tacacsTacacsDestAttr.Host = Host.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		tacacsTacacsDestAttr.Key = Key.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsTacacsDestAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		tacacsTacacsDestAttr.Port = Port.(string)
	}
	tacacsTacacsDest := models.NewTACACSDestination(fmt.Sprintf("tacacsdest-%s-port-%s", host, port), TACACSMonitoringDestinationGroupDn, desc, nameAlias, tacacsTacacsDestAttr)

	tacacsTacacsDest.Status = "modified"
	err := aciClient.Save(tacacsTacacsDest)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_file_rs_a_remote_host_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_file_rs_a_remote_host_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epp")
		checkDns = append(checkDns, newRelParam.(string))

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_file_rs_a_remote_host_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epg")
		err = aciClient.DeleteRelationfileRsARemoteHostToEpg(tacacsTacacsDest.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfileRsARemoteHostToEpg(tacacsTacacsDest.DistinguishedName, tacacsTacacsDestAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_file_rs_a_remote_host_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epp")
		err = aciClient.DeleteRelationfileRsARemoteHostToEpp(tacacsTacacsDest.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfileRsARemoteHostToEpp(tacacsTacacsDest.DistinguishedName, tacacsTacacsDestAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(tacacsTacacsDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTACACSDestinationRead(ctx, d, m)
}

func resourceAciTACACSDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsTacacsDest, err := getRemoteTACACSDestination(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setTACACSDestinationAttributes(tacacsTacacsDest, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fileRsARemoteHostToEpgData, err := aciClient.ReadRelationfileRsARemoteHostToEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fileRsARemoteHostToEpg %v", err)
		d.Set("relation_file_rs_a_remote_host_to_epg", "")
	} else {
		if _, ok := d.GetOk("relation_file_rs_a_remote_host_to_epg"); ok {
			tfName := d.Get("relation_file_rs_a_remote_host_to_epg").(string)
			if tfName != fileRsARemoteHostToEpgData {
				d.Set("relation_file_rs_a_remote_host_to_epg", "")
			}
		}
	}

	fileRsARemoteHostToEppData, err := aciClient.ReadRelationfileRsARemoteHostToEpp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fileRsARemoteHostToEpp %v", err)
		d.Set("relation_file_rs_a_remote_host_to_epp", "")
	} else {
		if _, ok := d.GetOk("relation_file_rs_a_remote_host_to_epp"); ok {
			tfName := d.Get("relation_file_rs_a_remote_host_to_epp").(string)
			if tfName != fileRsARemoteHostToEppData {
				d.Set("relation_file_rs_a_remote_host_to_epp", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTACACSDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "tacacsTacacsDest")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
