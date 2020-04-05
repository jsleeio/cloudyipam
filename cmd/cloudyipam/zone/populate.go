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
	"github.com/spf13/cobra"
)

var zonePopulateID string

var cmdZonePopulate = &cobra.Command{
	Use:   "populate",
	Short: "Populate a zone",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := shared.GetClient()
		if err != nil {
			fmt.Printf("error creating client: %v", err)
			os.Exit(1)
		}
		err = client.PopulateZone(zonePopulateID)
		if err != nil {
			fmt.Printf("error populating zone: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("populated zone: %s\n", zonePopulateID)
	},
}

func init() {
	cmdZonePopulate.Flags().StringVarP(&zonePopulateID, "id", "", "", "ID of zone to populate")
}
