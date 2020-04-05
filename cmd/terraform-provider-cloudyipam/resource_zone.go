package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix_length": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	zone := cloudyipam.Zone{
		Name:      d.Get("name").(string),
		Range:     d.Get("range").(string),
		PrefixLen: d.Get("prefix_length").(int),
	}
	id, err := ipam.CreateZone(zone)
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("name", zone.Name)
	d.Set("range", zone.Range)
	d.Set("prefix_length", zone.PrefixLen)
	d.SetId(id)
	return resourceZoneRead(d, m)
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	zone, err := ipam.ReadZone(d.Id())
	if zone != nil && err == nil {
		d.SetId(zone.Id)
		d.Set("name", zone.Name)
		d.Set("range", zone.Range)
		d.Set("prefix_length", zone.PrefixLen)
		return nil
	}
	if _, ok := err.(*cloudyipam.ReadZoneNotFoundError); ok {
		// successfully queried database; did not find zone
		d.SetId("")
		return nil
	}
	return err
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceZoneRead(d, m)
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	err = ipam.DestroyZone(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
