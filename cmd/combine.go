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
	"os"
	"strings"
)

// Get encoded shards from either stdin or args
func get_shards(args []string, useStdin bool) ([]string, error) {
	stdin_struct, err := common.Get_stdin()
	if err != nil {
		return nil, err
	}
	vprint.Printf("Args %v, \n Stdin: >>>%s<<<\n", args, stdin_struct.Stdin)

	if useStdin || stdin_struct.Has_stdin {
		if !useStdin {
			_, _ = fmt.Fprint(os.Stderr, "Warning: Stdin detected but --stdin flag not set. This may be undefined behavior \n")
		}
		if stdin_struct.Has_stdin {
			temp := strings.Trim(stdin_struct.Stdin, " \n ,")
			temp = strings.Replace(temp, " ", "\n", -1)
			outs := strings.Split(temp, "\n")
			vprint.Print(outs, "\n")
			return outs, nil
		} else {
			return nil, fmt.Errorf("Stdin flag `--` set, but was not able to detect stdin")
		}
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("Must have at least two arguments")
	}
	return args[1:], nil
}

var combineCmd = &cobra.Command{
	Use:     "combine",
	Aliases: []string{"c", "co", "com"},
	Short:   "Combine shards into a whole",
	Long: `Enter/read in shards and combine them to recover the original data
`,
	Run: func(cmd *cobra.Command, args []string) {
		vprint.Print("Run subcmd: combine\n")
		vprint.Print(args)

		useStdin, _ := cmd.Flags().GetBool("stdin")
		vprint.Print("useStdin", useStdin, "\n")

		formattings := common.FormatSettings{}
		err := common.ParseFormatSettings(&formattings, cmd)
		vprint.Print("Formattings (err): \n", formattings, "(", err, ")\n")
		common.LogIfFatal(err)
		shards, err := get_shards(args, useStdin)
		common.LogIfFatal(err)

		parts, err := common.DecodeShards(shards, formattings)
		common.LogIfFatal(err)
		recomb, err := shamir.Combine(parts)
		common.LogIfFatal(err)
		fmt.Println(string(recomb))
	},
}

func init() {
	RootCmd.AddCommand(combineCmd)

}
