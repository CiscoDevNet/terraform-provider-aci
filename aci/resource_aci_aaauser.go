package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciLocalUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLocalUserCreate,
		Update: resourceAciLocalUserUpdate,
		Read:   resourceAciLocalUserRead,
		Delete: resourceAciLocalUserDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLocalUserImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cert_attribute": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"clear_pwd_history": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expiration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expires": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"otpenable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"otpkey": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pwd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pwd_life_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pwd_update_required": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rbac_string": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unix_user_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteLocalUser(client *client.Client, dn string) (*models.LocalUser, error) {
	aaaUserCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaUser := models.LocalUserFromContainer(aaaUserCont)

	if aaaUser.DistinguishedName == "" {
		return nil, fmt.Errorf("LocalUser %s not found", aaaUser.DistinguishedName)
	}

	return aaaUser, nil
}

func setLocalUserAttributes(aaaUser *models.LocalUser, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(aaaUser.DistinguishedName)
	d.Set("description", aaaUser.Description)
	aaaUserMap, _ := aaaUser.ToMap()

	d.Set("name", aaaUserMap["name"])

	d.Set("account_status", aaaUserMap["accountStatus"])
	d.Set("annotation", aaaUserMap["annotation"])
	d.Set("cert_attribute", aaaUserMap["certAttribute"])
	d.Set("clear_pwd_history", aaaUserMap["clearPwdHistory"])
	d.Set("email", aaaUserMap["email"])
	d.Set("expiration", aaaUserMap["expiration"])
	d.Set("expires", aaaUserMap["expires"])
	d.Set("first_name", aaaUserMap["firstName"])
	d.Set("last_name", aaaUserMap["lastName"])
	d.Set("name_alias", aaaUserMap["nameAlias"])
	d.Set("otpenable", aaaUserMap["otpenable"])
	d.Set("otpkey", aaaUserMap["otpkey"])
	d.Set("phone", aaaUserMap["phone"])
	d.Set("pwd", aaaUserMap["pwd"])
	d.Set("pwd_life_time", aaaUserMap["pwdLifeTime"])
	d.Set("pwd_update_required", aaaUserMap["pwdUpdateRequired"])
	d.Set("rbac_string", aaaUserMap["rbacString"])
	d.Set("unix_user_id", aaaUserMap["unixUserId"])
	return d
}

func resourceAciLocalUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	aaaUser, err := getRemoteLocalUser(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLocalUserAttributes(aaaUser, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLocalUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LocalUser: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaUserAttr := models.LocalUserAttributes{}
	if AccountStatus, ok := d.GetOk("account_status"); ok {
		aaaUserAttr.AccountStatus = AccountStatus.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserAttr.Annotation = Annotation.(string)
	}
	if CertAttribute, ok := d.GetOk("cert_attribute"); ok {
		aaaUserAttr.CertAttribute = CertAttribute.(string)
	}
	if ClearPwdHistory, ok := d.GetOk("clear_pwd_history"); ok {
		aaaUserAttr.ClearPwdHistory = ClearPwdHistory.(string)
	}
	if Email, ok := d.GetOk("email"); ok {
		aaaUserAttr.Email = Email.(string)
	}
	if Expiration, ok := d.GetOk("expiration"); ok {
		aaaUserAttr.Expiration = Expiration.(string)
	}
	if Expires, ok := d.GetOk("expires"); ok {
		aaaUserAttr.Expires = Expires.(string)
	}
	if FirstName, ok := d.GetOk("first_name"); ok {
		aaaUserAttr.FirstName = FirstName.(string)
	}
	if LastName, ok := d.GetOk("last_name"); ok {
		aaaUserAttr.LastName = LastName.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaUserAttr.NameAlias = NameAlias.(string)
	}
	if Otpenable, ok := d.GetOk("otpenable"); ok {
		aaaUserAttr.Otpenable = Otpenable.(string)
	}
	if Otpkey, ok := d.GetOk("otpkey"); ok {
		aaaUserAttr.Otpkey = Otpkey.(string)
	}
	if Phone, ok := d.GetOk("phone"); ok {
		aaaUserAttr.Phone = Phone.(string)
	}
	if Pwd, ok := d.GetOk("pwd"); ok {
		aaaUserAttr.Pwd = Pwd.(string)
	}
	if PwdLifeTime, ok := d.GetOk("pwd_life_time"); ok {
		aaaUserAttr.PwdLifeTime = PwdLifeTime.(string)
	}
	if PwdUpdateRequired, ok := d.GetOk("pwd_update_required"); ok {
		aaaUserAttr.PwdUpdateRequired = PwdUpdateRequired.(string)
	}
	if RbacString, ok := d.GetOk("rbac_string"); ok {
		aaaUserAttr.RbacString = RbacString.(string)
	}
	if UnixUserId, ok := d.GetOk("unix_user_id"); ok {
		aaaUserAttr.UnixUserId = UnixUserId.(string)
	}
	aaaUser := models.NewLocalUser(fmt.Sprintf("userext/user-%s", name), "uni", desc, aaaUserAttr)

	err := aciClient.Save(aaaUser)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(aaaUser.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLocalUserRead(d, m)
}

func resourceAciLocalUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LocalUser: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaUserAttr := models.LocalUserAttributes{}
	if AccountStatus, ok := d.GetOk("account_status"); ok {
		aaaUserAttr.AccountStatus = AccountStatus.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserAttr.Annotation = Annotation.(string)
	}
	if CertAttribute, ok := d.GetOk("cert_attribute"); ok {
		aaaUserAttr.CertAttribute = CertAttribute.(string)
	}
	if ClearPwdHistory, ok := d.GetOk("clear_pwd_history"); ok {
		aaaUserAttr.ClearPwdHistory = ClearPwdHistory.(string)
	}
	if Email, ok := d.GetOk("email"); ok {
		aaaUserAttr.Email = Email.(string)
	}
	if Expiration, ok := d.GetOk("expiration"); ok {
		aaaUserAttr.Expiration = Expiration.(string)
	}
	if Expires, ok := d.GetOk("expires"); ok {
		aaaUserAttr.Expires = Expires.(string)
	}
	if FirstName, ok := d.GetOk("first_name"); ok {
		aaaUserAttr.FirstName = FirstName.(string)
	}
	if LastName, ok := d.GetOk("last_name"); ok {
		aaaUserAttr.LastName = LastName.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaUserAttr.NameAlias = NameAlias.(string)
	}
	if Otpenable, ok := d.GetOk("otpenable"); ok {
		aaaUserAttr.Otpenable = Otpenable.(string)
	}
	if Otpkey, ok := d.GetOk("otpkey"); ok {
		aaaUserAttr.Otpkey = Otpkey.(string)
	}
	if Phone, ok := d.GetOk("phone"); ok {
		aaaUserAttr.Phone = Phone.(string)
	}
	if Pwd, ok := d.GetOk("pwd"); ok {
		aaaUserAttr.Pwd = Pwd.(string)
	}
	if PwdLifeTime, ok := d.GetOk("pwd_life_time"); ok {
		aaaUserAttr.PwdLifeTime = PwdLifeTime.(string)
	}
	if PwdUpdateRequired, ok := d.GetOk("pwd_update_required"); ok {
		aaaUserAttr.PwdUpdateRequired = PwdUpdateRequired.(string)
	}
	if RbacString, ok := d.GetOk("rbac_string"); ok {
		aaaUserAttr.RbacString = RbacString.(string)
	}
	if UnixUserId, ok := d.GetOk("unix_user_id"); ok {
		aaaUserAttr.UnixUserId = UnixUserId.(string)
	}
	aaaUser := models.NewLocalUser(fmt.Sprintf("userext/user-%s", name), "uni", desc, aaaUserAttr)

	aaaUser.Status = "modified"

	err := aciClient.Save(aaaUser)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(aaaUser.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLocalUserRead(d, m)

}

func resourceAciLocalUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	aaaUser, err := getRemoteLocalUser(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLocalUserAttributes(aaaUser, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLocalUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaUser")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
