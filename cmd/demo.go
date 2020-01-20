/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/vault/shamir"
	"github.com/spf13/cobra"
	"github.com/xkortex/passcrux/common"
	"github.com/xkortex/vprint"
	"log"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo",
	Long:  `Run a demonstration`,
	Run: func(cmd *cobra.Command, args []string) {
		vprint.Print("demo")
		pass := "hunter2"
		secret := []byte(pass)
		fmt.Printf("Input: %s\n", pass)
		shards, err := shamir.Split(secret, 5, 3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Output:")
		fmt.Println(shards)
		settings := common.FormatSettings{
			common.EncodeBase32,
		}
		out, err := common.FormatShards(shards, settings)
		fmt.Println(out)
		for _, v := range shards {
			//fmt.Println(v)
			fmt.Println(hex.EncodeToString(v))
		}

		var recomb []byte
		// There is 5*4*3 possible choices,
		// we will just brute force try them all
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if j == i {
					continue
				}
				for k := 0; k < 5; k++ {
					if k == i || k == j {
						continue
					}
					parts := [][]byte{shards[i], shards[j], shards[k]}
					recomb, err = shamir.Combine(parts)

					if err != nil {
						log.Fatalf("err: %v", err)
					}

					if !bytes.Equal(recomb, secret) {
						_ = fmt.Errorf("parts: (i:%d, j:%d, k:%d) %v", i, j, k, parts)
						log.Fatalf("bad: %v %v", recomb, secret)
					}
				}
			}
			fmt.Println(string(recomb))
		}
	},
}

func init() {
	RootCmd.AddCommand(demoCmd)
}
