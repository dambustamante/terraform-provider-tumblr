package tumblr

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tumblr/tumblr.go"
	"github.com/tumblr/tumblrclient.go"
)

var fieldsPhotoPosts = []string{"caption", "link", "source"}

func resourcePostPhoto() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostPhotoCreate,
		Read:   resourcePostPhotoRead,
		Update: resourcePostPhotoUpdate,
		Delete: resourcePostPhotoDelete,

		Schema: map[string]*schema.Schema{
			"blog":   blogPostSchema(),
			"state":  statePostSchema(),
			"tags":   tagsPostSchema(),
			"date":   datePostSchema(),
			"format": formatPostSchema(),
			"slug":   slugPostSchema(),
			"caption": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["caption"],
			},
			"link": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["link"],
			},
			"source": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   descriptions["source_photo"],
				ConflictsWith: []string{"data", "data64"},
			},
			"data": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   descriptions["data_photo"],
				Removed:       "Pending to implement, default is data64",
				ConflictsWith: []string{"source", "data64"},
			},
			"data64": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   descriptions["data64"],
				ConflictsWith: []string{"source", "data"},
				StateFunc: func(val interface{}) string {
					return stringToMd5(val.(string))
				},
			},
		},
	}
}

func resourcePostPhotoCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tumblrclient.Client)

	_, sourceOk := d.GetOk("source")
	_, dataOk := d.GetOk("data")
	_, data64Ok := d.GetOk("data64")

	if !sourceOk && !dataOk && !data64Ok {
		return fmt.Errorf("One of source, data or data64 must be assigned")
	}

	params := generateParams(d, "photo", append(fieldsAllPosts, fieldsPhotoPosts...))
	res, err := tumblr.CreatePost(client, d.Get("blog").(string), params)
	if err != nil {
		return err
	}
	d.SetId(uintToString(res.Id))

	return resourcePostPhotoRead(d, m)
}

func resourcePostPhotoRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*tumblrclient.Client)

	params := url.Values{}
	params.Add("type", "photo")
	params.Add("id", d.Id())
	res, err := tumblr.GetPosts(client, d.Get("blog").(string), params)
	if err != nil {
		d.SetId("")
		return nil
	}

	for _, key := range append(fieldsAllPosts, fieldsPhotoPosts...) {
		value, err := res.Get(0).GetProperty(toCamelCase(key))
		if err == nil {
			d.Set(key, value)
		}
	}

	return nil
}

func resourcePostPhotoUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*tumblrclient.Client)

	params := generateParams(d, "photo", append(fieldsAllPosts, fieldsPhotoPosts...))
	err := tumblr.EditPost(client, d.Get("blog").(string), stringToUint(d.Id()), params)
	if err != nil {
		return err
	}
	return resourcePostPhotoRead(d, m)
}

func resourcePostPhotoDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*tumblrclient.Client)

	err := tumblr.DeletePost(client, d.Get("blog").(string), stringToUint(d.Id()))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
