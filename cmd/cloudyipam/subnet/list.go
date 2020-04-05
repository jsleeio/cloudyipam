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
package subnet

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/shared"
	"github.com/spf13/cobra"
)

var subnetListZone string

var cmdSubnetList = &cobra.Command{
	Use:   "list",
	Short: "List subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := shared.GetClient()
		if err != nil {
			fmt.Printf("error creating client: %v", err)
			os.Exit(1)
		}
		subnets, err := client.ListSubnets()
		if err != nil {
			fmt.Printf("error listing subnets: %v\n", err)
			os.Exit(1)
		}
		if len(subnets) < 1 {
			return
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 3, 4, ' ', 0)
		// \v\n is not an error, go doc text/tabwriter
		fmt.Fprintf(w, "%s\v%s\v%s\v%s\v%v\v\n", "USAGE", "SUBNET ID", "ZONE ID", "CIDR BLOCK", "AVAILABLE")
		for _, subnet := range subnets {
			fmt.Fprintf(w, "%s\v%s\v%s\v%s\v%v\v\n", subnet.Usage, subnet.Id, subnet.Zone, subnet.Range, subnet.Available)
		}
		w.Flush()
	},
}
