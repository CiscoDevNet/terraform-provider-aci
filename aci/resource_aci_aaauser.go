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

func resourceAciLocalUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLocalUserCreate,
		UpdateContext: resourceAciLocalUserUpdate,
		ReadContext:   resourceAciLocalUserRead,
		DeleteContext: resourceAciLocalUserDelete,

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
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"inactive",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
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

func tryLogin(d *schema.ResourceData, aciClient *client.Client) bool {

	authPayload := `{
		"aaaUser" : {
			"attributes" : {
				"name" : "%s",
				"pwd" : "%s"
			}
		}
	}`

	client.SkipLoggingPayload(true)(aciClient)

	req, err := aciClient.MakeRestRequestRaw("POST", "/api/aaaLogin.json", []byte(fmt.Sprintf(authPayload, d.Get("name").(string), d.Get("pwd").(string))), false)
	if err != nil {
		log.Printf("[DEBUG] Creation of authentication request failed for local user: '%s'.", d.Get("name").(string))
		return false
	}

	_, resp, _ := aciClient.Do(req)

	client.SkipLoggingPayload(false)(aciClient)

	if resp.StatusCode != 200 {
		log.Printf("[DEBUG] Authentication failed for local user: '%s'.", d.Get("name").(string))
		return false
	}

	log.Printf("[DEBUG] Authentication succceeded for local user: '%s'.", d.Get("name").(string))
	return true
}

func getRemoteLocalUser(client *client.Client, dn string) (*models.LocalUser, error) {
	aaaUserCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaUser := models.LocalUserFromContainer(aaaUserCont)

	if aaaUser.DistinguishedName == "" {
		return nil, fmt.Errorf("Local User %s not found", dn)
	}

	return aaaUser, nil
}

func setLocalUserAttributes(aaaUser *models.LocalUser, d *schema.ResourceData, loginSuccess bool) (*schema.ResourceData, error) {
	d.SetId(aaaUser.DistinguishedName)
	d.Set("description", aaaUser.Description)
	aaaUserMap, err := aaaUser.ToMap()
	if err != nil {
		return d, err
	}

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
	if aaaUserMap["pwdLifeTime"] == "no-password-expire" {
		d.Set("pwd_life_time", "0")
	} else {
		d.Set("pwd_life_time", aaaUserMap["pwdLifeTime"])
	}
	d.Set("pwd_update_required", aaaUserMap["pwdUpdateRequired"])
	d.Set("rbac_string", aaaUserMap["rbacString"])
	d.Set("unix_user_id", aaaUserMap["unixUserId"])

	if !loginSuccess {
		d.Set("pwd", "")
	}

	return d, nil
}

func resourceAciLocalUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	aaaUser, err := getRemoteLocalUser(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLocalUserAttributes(aaaUser, d, false)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLocalUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	} else {
		aaaUserAttr.Annotation = "{}"
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
	if pwd, ok := d.GetOk("pwd"); ok {
		loginSuccess := tryLogin(d, aciClient)
		if !loginSuccess {
			client.SkipLoggingPayload(true)(aciClient)
			aaaUserAttr.Pwd = pwd.(string)
		}
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
	client.SkipLoggingPayload(false)(aciClient)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaUser.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLocalUserRead(ctx, d, m)
}

func resourceAciLocalUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	} else {
		aaaUserAttr.Annotation = "{}"
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
	if d.HasChange("pwd") {
		if pwd, ok := d.GetOk("pwd"); ok {
			loginSuccess := tryLogin(d, aciClient)
			if !loginSuccess {
				client.SkipLoggingPayload(true)(aciClient)
				aaaUserAttr.Pwd = pwd.(string)
			}
		}
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
	client.SkipLoggingPayload(false)(aciClient)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaUser.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLocalUserRead(ctx, d, m)

}

func resourceAciLocalUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	aaaUser, err := getRemoteLocalUser(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	loginSuccess := false
	if _, ok := d.GetOk("pwd"); ok {
		loginSuccess = tryLogin(d, aciClient)
	}

	_, err = setLocalUserAttributes(aaaUser, d, loginSuccess)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLocalUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaUser")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
