package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciSPANSourcedestinationGroupMatchLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSPANSourcedestinationGroupMatchLabelCreate,
		Update: resourceAciSPANSourcedestinationGroupMatchLabelUpdate,
		Read:   resourceAciSPANSourcedestinationGroupMatchLabelRead,
		Delete: resourceAciSPANSourcedestinationGroupMatchLabelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSPANSourcedestinationGroupMatchLabelImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"span_source_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"black", "navy", "dark-blue", "medium-blue", "blue", "dark-green", "green", "teal", "dark-cyan", "deep-sky-blue",
					"dark-turquoise", "medium-spring-green", "lime", "spring-green", "aqua", "cyan", "midnight-blue",
					"dodger-blue", "light-sea-green", "forest-green", "sea-green", "dark-slate-gray", "lime-green",
					"medium-sea-green", "turquoise", "royal-blue", "steel-blue", "dark-slate-blue", "medium-turquoise",
					"indigo", "dark-olive-green", "cadet-blue", "cornflower-blue", "medium-aquamarine", "dim-gray",
					"slate-blue", "olive-drab", "slate-gray", "light-slate-gray", "medium-slate-blue", "lawn-green", "chartreuse",
					"aquamarine", "maroon", "purple", "olive", "gray", "sky-blue", "light-sky-blue", "blue-violet", "dark-red",
					"dark-magenta", "saddle-brown", "dark-sea-green", "light-green", "medium-purple", "dark-violet", "pale-green",
					"dark-orchid", "yellow-green", "sienna", "brown", "dark-gray", "light-blue", "green-yellow", "pale-turquoise",
					"light-steel-blue", "powder-blue", "fire-brick", "dark-goldenrod", "medium-orchid", "rosy-brown", "dark-khaki",
					"silver", "medium-violet-red", "indian-red", "peru", "chocolate", "tan", "light-gray", "thistle", "orchid",
					"goldenrod", "pale-violet-red", "crimson", "gainsboro", "plum", "burlywood", "light-cyan", "lavender",
					"dark-salmon", "violet", "pale-goldenrod", "light-coral", "khaki", "alice-blue", "honeydew", "azure",
					"sandy-brown", "wheat", "beige", "white-smoke", "mint-cream", "ghost-white", "salmon", "antique-white",
					"linen", "light-goldenrod-yellow", "old-lace", "red", "fuchsia", "magenta", "deep-pink", "orange-red",
					"tomato", "hot-pink", "coral", "dark-orange", "light-salmon", "orange", "light-pink", "pink", "gold",
					"peachpuff", "navajo-white", "moccasin", "bisque", "misty-rose", "blanched-almond", "papaya-whip", "lavender-blush",
					"seashell", "cornsilk", "lemon-chiffon", "floral-white", "snow", "yellow", "light-yellow", "ivory", "white",
				}, false),
			},
		}),
	}
}
func getRemoteSPANSourcedestinationGroupMatchLabel(client *client.Client, dn string) (*models.SPANSourcedestinationGroupMatchLabel, error) {
	spanSpanLblCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	spanSpanLbl := models.SPANSourcedestinationGroupMatchLabelFromContainer(spanSpanLblCont)

	if spanSpanLbl.DistinguishedName == "" {
		return nil, fmt.Errorf("SPANSourcedestinationGroupMatchLabel %s not found", spanSpanLbl.DistinguishedName)
	}

	return spanSpanLbl, nil
}

func setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl *models.SPANSourcedestinationGroupMatchLabel, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(spanSpanLbl.DistinguishedName)
	d.Set("description", spanSpanLbl.Description)
	// d.Set("span_source_group_dn", GetParentDn(spanSpanLbl.DistinguishedName))
	if dn != spanSpanLbl.DistinguishedName {
		d.Set("span_source_group_dn", "")
	}
	spanSpanLblMap, _ := spanSpanLbl.ToMap()

	d.Set("name", spanSpanLblMap["name"])

	d.Set("annotation", spanSpanLblMap["annotation"])
	d.Set("name_alias", spanSpanLblMap["nameAlias"])
	d.Set("tag", spanSpanLblMap["tag"])
	return d
}

func resourceAciSPANSourcedestinationGroupMatchLabelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSPANSourcedestinationGroupMatchLabelCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SPANSourcedestinationGroupMatchLabel: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	spanSpanLblAttr := models.SPANSourcedestinationGroupMatchLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSpanLblAttr.Annotation = Annotation.(string)
	} else {
		spanSpanLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSpanLblAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		spanSpanLblAttr.Tag = Tag.(string)
	}
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(fmt.Sprintf("spanlbl-%s", name), SPANSourceGroupDn, desc, spanSpanLblAttr)

	err := aciClient.Save(spanSpanLbl)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(spanSpanLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANSourcedestinationGroupMatchLabelRead(d, m)
}

func resourceAciSPANSourcedestinationGroupMatchLabelUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SPANSourcedestinationGroupMatchLabel: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	spanSpanLblAttr := models.SPANSourcedestinationGroupMatchLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSpanLblAttr.Annotation = Annotation.(string)
	} else {
		spanSpanLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSpanLblAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		spanSpanLblAttr.Tag = Tag.(string)
	}
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(fmt.Sprintf("spanlbl-%s", name), SPANSourceGroupDn, desc, spanSpanLblAttr)

	spanSpanLbl.Status = "modified"

	err := aciClient.Save(spanSpanLbl)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(spanSpanLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANSourcedestinationGroupMatchLabelRead(d, m)

}

func resourceAciSPANSourcedestinationGroupMatchLabelRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSPANSourcedestinationGroupMatchLabelDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanSpanLbl")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
