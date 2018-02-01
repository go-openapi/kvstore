// Copyright © 2016 Ivan Porto Carrero
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"

	"github.com/go-openapi/kvstore/api/client"
	"github.com/spf13/cobra"
)

var etag uint64

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Create/Update an entry",
	Long:  `Create or update an entry in the k/v store`,
	Run: func(cmd *cobra.Command, args []string) {
		cl, err := client.New(url)
		if err != nil {
			log.Fatalln(err)
		}
		var key string
		if len(args) > 0 {
			key = args[0]
		}
		var data string
		if len(args) > 1 {
			data = args[1]
		}

		log.Printf("updating entry for key %q with value %q", key, data)
		entry := &client.Entry{
			Data:    []byte(data),
			Version: etag,
		}
		err = cl.Put(key, entry)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("version:", entry.Version)
		fmt.Println(string(data))
	},
}

func init() {
	RootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().String("foo", "", "A help for foo")
	putCmd.Flags().Uint64Var(&etag, "version", 0, "The version for updating a key in the k/v store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
