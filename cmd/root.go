/*
Copyright © 2019 MICHAEL McDERMOTT
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. If a copy of the MPL was not
distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xkortex/vprint"
	"log"
	"os"
)

var (
	Version = "unset"
)

func PrintVersionAndQuit() {
	fmt.Println(Version)
	os.Exit(0)
}

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "passcrux",
	Short: "Utility for splitting passwords with Shamir's Secret Sharing",
	Long: `Utility for splitting passwords with Shamir's Secret Sharing. 
Takes a password and splits it
`,
	Run: func(cmd *cobra.Command, args []string) {
		vprint.Printf("root called. passcrux %s\n", Version)
		doVersion, _ := cmd.Flags().GetBool("version")
		if doVersion {
			PrintVersionAndQuit()

		}

		vprint.Print(args)
		_ = cmd.Help()
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//RootCmd.AddCommand(RootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.PersistentFlags().BoolP("silent", "S", false, "Suppress errors")
	RootCmd.PersistentFlags().BoolP("prompt", "p", false, "Read key/password from standard in prompt")
	RootCmd.PersistentFlags().BoolP("dummy", "D", false, "Testing")

	// Formatting parameters
	RootCmd.PersistentFlags().StringP("enc", "e", "hex", "En/decoding format {[he]x, [base]32, [base]64, abc, ABC}")
	RootCmd.PersistentFlags().StringP("input", "i", "", "Input file to read (optional)")
	RootCmd.PersistentFlags().StringP("sep", "s", "", "Separator between bytes/fields")
	RootCmd.PersistentFlags().StringP("recordsep", "R", "\n", "Separator between records/shards")
	RootCmd.PersistentFlags().IntP("fieldsize", "F", 2, "Number of bytes/characters per separated field")

	// Runtime
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
	RootCmd.PersistentFlags().BoolP("version", "V", false, "Print version and quit")

}

func initConfig() {
	// todo: use init config to do stuff based on env
}
