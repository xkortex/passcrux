/*
Copyright © 2019 MICHAEL McDERMOTT
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. If a copy of the MPL was not
distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xkortex/vprint"
	"log"
)

var developer string

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "passcrux",
	Short: "Utility for splitting passwords with Shamir's Secret Sharing",
	Long: `Takes a password and splits it`,
	Run: func(cmd *cobra.Command, args []string) {
		vprint.Print("root called")

		vprint.Print(args)
		ns, _ := cmd.PersistentFlags().GetString("namespace")
		vprint.Print(ns)


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
	RootCmd.PersistentFlags().StringP("namespace", "n", "", "namespace of kv store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.PersistentFlags().BoolP("silent", "s", false, "Suppress errors")
	RootCmd.PersistentFlags().BoolP("stdin", "-", false, "Read from standard in")
	RootCmd.Flags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
	RootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")

}

func initConfig() {
	// todo: use init config to do stuff based on env
}
