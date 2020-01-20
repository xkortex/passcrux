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
	"regexp"
	"strconv"
	"strings"
)

// Get the password/phrase/key from input
// todo: deal with extra character trim issue
func get_password(args []string, useStdin bool, useRaw bool) (string, error) {
	stdin_struct, err := common.Get_stdin()
	if err != nil {
		return "", err
	}
	temp := ""
	vprint.Print("args: ", args, "stdin: ", stdin_struct.Has_stdin, " ", stdin_struct.Stdin, "\n")
	if stdin_struct.Has_stdin && !useStdin {
		_, _ = fmt.Fprint(os.Stderr, "Warning: stdin pipe detected, but --stdin flag not set. This is currently untested behavior \n")
		useStdin = true
	}
	if useStdin {
		if !stdin_struct.Has_stdin {
			return "", fmt.Errorf("Stdin flag `--` set, but was not able to detect stdin")
		}
		if len(args) > 1 {
			_, _ = fmt.Fprint(os.Stderr, "Warning: stdin pipe detected, but arguments passed. Ignoring arguments and using stdin\n")
		}
		temp = stdin_struct.Stdin
	} else if len(args) > 0 {
		if len(args) > 1 {
			_, _ = fmt.Fprint(os.Stderr, "Warning: More than one argument detected. Ignoring all but first. Use stdin pipe if you want breaks in your data\n")
		}
		temp = args[0]
	} else {
		return "", fmt.Errorf("Input error: Must have at least one argument, or --stdin piped in")
	}

	if temp == "" {
		panic("Input string is blank. This should not ever happen")
	}
	if !useRaw {
		temp = strings.Trim(temp, " \n")
	}
	return temp, nil
}

func ParseSplitSettings(settings *common.SplitSettings, cmd *cobra.Command) error {
	ratioArg, _ := cmd.Flags().GetString("ratio")
	re := regexp.MustCompile(`(\d+)[/|:](\d+)`)
	match := re.FindSubmatch([]byte(ratioArg))
	if len(match) != 3 {
		return fmt.Errorf("Unable to parse `ratio` argument: '%s'", ratioArg)
	}
	a, err := strconv.ParseInt(string(match[1]), 10, 32)
	if err != nil {
		return err
	}
	b, err := strconv.ParseInt(string(match[2]), 10, 32)
	if err != nil {
		return err
	}
	vprint.Printf("a: %d b: %d\n", a, b)
	var threshold, parts int64
	if a <= b {
		threshold = a
		parts = b
	} else {
		threshold = b
		parts = a
	}

	if threshold < 2 || parts < 2 || threshold > 256 || parts > 256 {
		return fmt.Errorf("parts and threshold must be 2 < x < 256")
	}
	settings.Parts = int(parts)
	settings.Threshold = int(threshold)
	return nil
}

var splitCmd = &cobra.Command{
	Use:     "split",
	Aliases: []string{"s"},
	Short:   "Split a password into shards",
	Long: `Given some data, split it up into N shards, where M is the shards required to reconstruct it.
Ratio is "M/N"
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		vprint.Print("Run subcmd: split\n")

		useStdin, _ := cmd.Flags().GetBool("stdin")
		vprint.Print("useStdin: ", useStdin, "\n")
		pass, err := get_password(args, useStdin, false)
		common.LogIfFatal(err)
		splittings := common.SplitSettings{}
		err = ParseSplitSettings(&splittings, cmd)
		vprint.Print("Splittings: \n  ", splittings, " : ", err, "\n")
		common.LogIfFatal(err)

		secret := []byte(pass)
		vprint.Printf("Len in bytes: %d\n", len(secret))
		vprint.Printf("Input: %s\n", pass)
		shards, err := shamir.Split(secret, splittings.Parts, splittings.Threshold)
		common.LogIfFatal(err)
		vprint.Printf("Len of each shard in bytes: %d\n", len(shards[0]))
		vprint.Print("Output:\n")
		settings := common.FormatSettings{
			common.EncodeHex,
		}
		out, err := common.FormatShards(shards, settings)
		fmt.Println(out)

	},
}

func init() {
	RootCmd.AddCommand(splitCmd)
	RootCmd.PersistentFlags().StringP("ratio", "r", "3/5", "Ratio of parts needed to total parts")

}
