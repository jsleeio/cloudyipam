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
	"text/tabwriter"

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/shared"
	"github.com/spf13/cobra"
)

var cmdZoneList = &cobra.Command{
	Use:   "list",
	Short: "List zones",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := shared.GetClient()
		if err != nil {
			fmt.Printf("error creating client: %v", err)
			os.Exit(1)
		}
		zones, err := client.ListZones()
		if err != nil {
			fmt.Printf("error listing zones: %v\n", err)
			os.Exit(1)
		}
		if len(zones) < 1 {
			return
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 3, 4, ' ', 0)
		// \v\n is not an error, go doc text/tabwriter
		fmt.Fprintf(w, "%s\v%s\v%s\v%v\v\n", "NAME", "ZONE ID", "CIDR BLOCK", "SUBNET PREFIXLEN")
		for _, zone := range zones {
			fmt.Fprintf(w, "%s\v%s\v%s\v%v\v\n", zone.Name, zone.Id, zone.Range, zone.PrefixLen)
		}
		w.Flush()
	},
}
