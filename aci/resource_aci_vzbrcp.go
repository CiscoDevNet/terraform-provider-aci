package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractCreate,
		Update: resourceAciContractUpdate,
		Read:   resourceAciContractRead,
		Delete: resourceAciContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractImport,
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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vz_rs_graph_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},

			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"filter_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"annotation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"name_alias": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"filter_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}
func getRemoteContract(client *client.Client, dn string) (*models.Contract, error) {
	vzBrCPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzBrCP := models.ContractFromContainer(vzBrCPCont)

	if vzBrCP.DistinguishedName == "" {
		return nil, fmt.Errorf("Contract %s not found", vzBrCP.DistinguishedName)
	}

	return vzBrCP, nil
}

func getRemoteFilterFromContract(client *client.Client, dn string) (*models.Filter, error) {
	vzFilterCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzFilter := models.FilterFromContainer(vzFilterCont)

	if vzFilter.DistinguishedName == "" {
		return nil, fmt.Errorf("Filter %s not found", vzFilter.DistinguishedName)
	}

	return vzFilter, nil
}

func setContractAttributes(vzBrCP *models.Contract, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vzBrCP.DistinguishedName)
	d.Set("description", vzBrCP.Description)
	// d.Set("tenant_dn", GetParentDn(vzBrCP.DistinguishedName))
	if dn != vzBrCP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzBrCPMap, _ := vzBrCP.ToMap()

	d.Set("name", vzBrCPMap["name"])

	d.Set("annotation", vzBrCPMap["annotation"])
	d.Set("name_alias", vzBrCPMap["nameAlias"])
	d.Set("prio", vzBrCPMap["prio"])
	d.Set("scope", vzBrCPMap["scope"])
	d.Set("target_dscp", vzBrCPMap["targetDscp"])
	return d
}

func setFilterAttributesFromContract(vzfilters []*models.Filter, d *schema.ResourceData) *schema.ResourceData {
	log.Println("Check .... :", vzfilters)
	filterSet := make([]interface{}, 0, 1)
	for _, filter := range vzfilters {
		fMap := make(map[string]interface{})
		fMap["description"] = filter.Description
		fMap["id"] = filter.DistinguishedName

		vzFilterMap, _ := filter.ToMap()
		fMap["filter_name"] = vzFilterMap["name"]
		fMap["annotation"] = vzFilterMap["annotation"]
		fMap["name_alias"] = vzFilterMap["nameAlias"]

		filterSet = append(filterSet, fMap)
	}
	log.Println("Check ...:", filterSet)
	d.Set("filter", filterSet)
	return d
}

func resourceAciContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setContractAttributes(vzBrCP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Contract: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	err := aciClient.Save(vzBrCP)
	if err != nil {
		return err
	}

	if filters, ok := d.GetOk("filter"); ok {
		// filterSet := make([]interface{}, 0, 1)
		filterIDS := make([]string, 0, 1)
		vzfilters := filters.(*schema.Set).List()
		for _, val := range vzfilters {
			vzFilterAttr := models.FilterAttributes{}
			filter := val.(map[string]interface{})

			name := filter["filter_name"].(string)

			desc := filter["description"].(string)

			if filter["annotation"] != nil {
				vzFilterAttr.Annotation = filter["annotation"].(string)
			}

			if filter["name_alias"] != nil {
				vzFilterAttr.NameAlias = filter["name_alias"].(string)
			}

			vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

			err := aciClient.Save(vzFilter)
			if err != nil {
				return err
			}

			// fMap := make(map[string]interface{})
			// fMap["id"] = vzFilter.DistinguishedName
			filterIDS = append(filterIDS, vzFilter.DistinguishedName)
		}
		log.Println("Check ... :", filterIDS)
		d.Set("filter_ids", filterIDS)
	}

	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTovzRsGraphAtt, ok := d.GetOk("relation_vz_rs_graph_att"); ok {
		relationParam := relationTovzRsGraphAtt.(string)
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_graph_att")
		d.Partial(false)

	}

	d.SetId(vzBrCP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractRead(d, m)
}

func resourceAciContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Contract: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzBrCPAttr := models.ContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzBrCPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzBrCPAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzBrCPAttr.Prio = Prio.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		vzBrCPAttr.Scope = Scope.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzBrCPAttr.TargetDscp = TargetDscp.(string)
	}
	vzBrCP := models.NewContract(fmt.Sprintf("brc-%s", name), TenantDn, desc, vzBrCPAttr)

	vzBrCP.Status = "modified"

	err := aciClient.Save(vzBrCP)

	if err != nil {
		return err
	}

	if d.HasChange("filter") {
		filter := d.Get("filter_ids").([]interface{})
		for _, val := range filter {
			filterDN := val.(string)
			err := aciClient.DeleteByDn(filterDN, "vzFilter")
			if err != nil {
				return err
			}
		}

		filters := d.Get("filter")
		filterIDS := make([]string, 0, 1)
		vzfilters := filters.(*schema.Set).List()
		for _, val := range vzfilters {
			vzFilterAttr := models.FilterAttributes{}
			filter := val.(map[string]interface{})

			name := filter["filter_name"].(string)

			desc := filter["description"].(string)

			if filter["annotation"] != nil {
				vzFilterAttr.Annotation = filter["annotation"].(string)
			}

			if filter["name_alias"] != nil {
				vzFilterAttr.NameAlias = filter["name_alias"].(string)
			}

			vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

			// vzFilter.Status = "modified"
			err := aciClient.Save(vzFilter)
			if err != nil {
				return err
			}

			filterIDS = append(filterIDS, vzFilter.DistinguishedName)
		}

		d.Set("filter_ids", filterIDS)
	}

	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_vz_rs_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_graph_att")
		err = aciClient.DeleteRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsGraphAttFromContract(vzBrCP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_graph_att")
		d.Partial(false)

	}

	d.SetId(vzBrCP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractRead(d, m)

}

func resourceAciContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzBrCP, err := getRemoteContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setContractAttributes(vzBrCP, d)

	filters := d.Get("filter_ids").([]interface{})
	log.Println("Check ... :", filters)
	vzFilters := make([]*models.Filter, 0, 1)
	for _, val := range filters {
		filterDN := val.(string)
		vzfilter, err := getRemoteFilterFromContract(aciClient, filterDN)
		if err == nil {
			vzFilters = append(vzFilters, vzfilter)
		}
	}
	setFilterAttributesFromContract(vzFilters, d)

	vzRsGraphAttData, err := aciClient.ReadRelationvzRsGraphAttFromContract(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsGraphAtt %v", err)

	} else {
		d.Set("relation_vz_rs_graph_att", vzRsGraphAttData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzBrCP")
	if err != nil {
		return err
	}

	filters := d.Get("filter_ids").([]interface{})
	for _, val := range filters {
		filterDN := val.(string)
		err := aciClient.DeleteByDn(filterDN, "vzFilter")
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
