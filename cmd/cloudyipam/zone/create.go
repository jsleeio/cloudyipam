/*
Copyright © 2020 John Slee <john@sleefamily.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package zone

import (
	"fmt"
	"os"

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/shared"
	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
	"github.com/spf13/cobra"
)

var zoneCreateName, zoneCreateCIDR string
var zoneCreatePrefixLen int

var cmdZoneCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a zone",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := shared.GetClient()
		if err != nil {
			fmt.Printf("error creating client: %v", err)
			os.Exit(1)
		}
		newzone := cloudyipam.Zone{
			Name:      zoneCreateName,
			Range:     zoneCreateCIDR,
			PrefixLen: zoneCreatePrefixLen,
		}
		uuid, err := client.CreateZone(newzone)
		if err != nil {
			fmt.Printf("error creating zone: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("created zone: %s\n", uuid)
	},
}

func init() {
	cmdZoneCreate.Flags().StringVarP(&zoneCreateName, "name", "", "", "name of zone to be created")
	cmdZoneCreate.Flags().StringVarP(&zoneCreateCIDR, "cidr", "", "", "CIDR specification of zone")
	cmdZoneCreate.Flags().IntVarP(&zoneCreatePrefixLen, "prefix-length", "", 24, "Prefix length for all subnets in zone")
}
