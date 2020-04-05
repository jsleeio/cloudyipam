package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
)

type CloudyIPAMConfiguration struct {
	DB cloudyipam.DatabaseConfig
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &CloudyIPAMConfiguration{
		DB: cloudyipam.DatabaseConfig{
			User:     d.Get("username").(string),
			Password: d.Get("password").(string),
			Host:     d.Get("hostname").(string),
			Database: d.Get("database").(string),
			Port:     d.Get("port").(string),
			TLS:      d.Get("tls").(bool),
		},
	}
	return config, nil
}

func connectToCloudyIPAM(conf *CloudyIPAMConfiguration) (*cloudyipam.CloudyIPAM, error) {
	dsn := conf.DB.DSN()
	var ipam *cloudyipam.CloudyIPAM
	var err error

	// When provisioning a database server there can often be a lag between
	// when Terraform thinks it's available and when it is actually available.
	// This is particularly acute when provisioning a server and then immediately
	// trying to provision a database on it.
	retryError := resource.Retry(5*time.Minute, func() *resource.RetryError {
		ipam, err = cloudyipam.NewCloudyIPAM(dsn)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryError != nil {
		return nil, fmt.Errorf("Could not connect to CloudyIPAM database: %s", retryError)
	}
	return ipam, nil
}

func boolEnvDefaultFunc(env string, dv interface{}) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		switch strings.ToLower(strings.TrimSpace(os.Getenv(env))) {
		case "":
			return dv, nil
		case "true":
			return true, nil
		case "yes":
			return true, nil
		case "false":
			return false, nil
		case "no":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid value for boolean environment variable " + env + ", must be one of: true yes false no")
		}
	}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"cloudyipam_zone":   resourceZone(),
			"cloudyipam_subnet": resourceSubnet(),
		},
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDYIPAM_HOSTNAME", nil),
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDYIPAM_PORT", "5432"),
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDYIPAM_NAME", "cloudyipam"),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDYIPAM_USERNAME", "cloudyipam"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDYIPAM_PASSWORD", nil),
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: boolEnvDefaultFunc("CLOUDYIPAM_TLS", true),
			},
		},
	}
}
