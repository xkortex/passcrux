/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"fmt"
	"github.com/hashicorp/vault/shamir"
	"github.com/spf13/cobra"
	"github.com/xkortex/passcrux/common"
	"github.com/xkortex/vprint"
	"log"
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a password into shards",
	Long: `Given some data, split it up into N shards, where M is the shards required to reconstruct it.
Ratio is "M/N"
`,
	Run: func(cmd *cobra.Command, args []string) {
		vprint.Print("split")
		pass := args[0]

		secret := []byte(pass)
		fmt.Printf("Input: %s\n", pass)
		shards, err := shamir.Split(secret, 5, 3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Output:")
		settings := common.FormatSettings{
			common.EncodeBase32,
		}
		out, err := common.FormatShards(shards, settings)
		fmt.Println(out)

	},
}

func init() {
	RootCmd.AddCommand(splitCmd)
}
