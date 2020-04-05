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

	"github.com/jsleeio/cloudyipam/cmd/cloudyipam/shared"
	"github.com/spf13/cobra"
)

var subnetFreeSubnet string

var cmdSubnetFree = &cobra.Command{
	Use:   "free",
	Short: "De-allocate a subnet",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := shared.GetClient()
		if err != nil {
			fmt.Printf("error creating client: %v", err)
			os.Exit(1)
		}
		err = client.FreeSubnet(subnetFreeSubnet)
		if err != nil {
			fmt.Printf("error freeing subnet %v: %v\n", subnetFreeSubnet, err)
			os.Exit(1)
		}
		fmt.Printf("freed subnet: %s\n", subnetFreeSubnet)
	},
}

func init() {
	cmdSubnetFree.Flags().StringVarP(&subnetFreeSubnet, "id", "", "", "ID of subnet to free")
}
