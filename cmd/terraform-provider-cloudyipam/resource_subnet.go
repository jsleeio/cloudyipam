package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubnetCreate,
		Read:   resourceSubnetRead,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,

		Schema: map[string]*schema.Schema{
			"available": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"range": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSubnetCreate(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	subnet, err := ipam.AllocateSubnet(d.Get("zone_id").(string), d.Get("usage").(string))
	if err != nil {
		if _, ok := err.(cloudyipam.ReadSubnetNotFoundError); ok {
			d.SetId("")
		}
		return err
	}
	d.SetId(subnet.Id)
	d.Set("zone_id", subnet.Zone)
	d.Set("range", subnet.Range)
	d.Set("available", subnet.Available)
	d.Set("usage", subnet.Usage)
	return resourceSubnetRead(d, m)
}

func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	subnet, err := ipam.ReadSubnet(d.Id())
	if subnet != nil {
		d.Set("zone_id", subnet.Zone)
		d.Set("range", subnet.Range)
		d.Set("available", subnet.Available)
		d.Set("usage", subnet.Usage)
		return nil
	}
	if _, ok := err.(*cloudyipam.ReadSubnetNotFoundError); ok {
		// successfully queried database; did not find zone
		d.SetId("")
		return nil
	}
	return err
}

func resourceSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSubnetRead(d, m)
}

func resourceSubnetDelete(d *schema.ResourceData, m interface{}) error {
	ipam, err := connectToCloudyIPAM(m.(*CloudyIPAMConfiguration))
	if err != nil {
		return err
	}
	err = ipam.FreeSubnet(d.Get("id").(string))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
