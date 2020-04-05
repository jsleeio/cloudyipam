package shared

import (
	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
)

func GetClient() (*cloudyipam.CloudyIPAM, error) {
	return cloudyipam.NewCloudyIPAM(DatabaseConfigFromViper().DSN())
}
