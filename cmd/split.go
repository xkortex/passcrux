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
// todo: figure out what I meant by "character trim issue"
func getKeyData(args []string, useRaw bool, usePass bool) (string, error) {
	stdin_struct, err := common.Get_stdin()
	if err != nil {
		return "", err
	}
	temp := ""
	vprint.Print("args: ", args, "stdin: ", stdin_struct.Has_stdin, " ", stdin_struct.Stdin, "\n")
	if stdin_struct.Has_stdin {
		if len(args) > 1 {
			_, _ = fmt.Fprint(os.Stderr, "Warning: stdin pipe detected, but arguments passed. Ignoring arguments and using stdin\n")
		}
		temp = stdin_struct.Stdin
	} else if usePass {
		temp, err = common.ReadPassword()
		if err != nil {
			return "", err
		}
	} else if len(args) > 0 {
		if len(args) > 1 {
			_, _ = fmt.Fprint(os.Stderr, "Warning: More than one argument detected. Ignoring all but first. Use stdin pipe if you want breaks in your data\n")
		}
		temp = args[0]
	} else {
		return "", fmt.Errorf("Input error: Must have at least one argument, or stdin piped in")
	}

	if temp == "" {
		panic("Input string is blank. This should not ever happen")
	}
	if !useRaw {
		temp = strings.Trim(temp, " \n")
	}
	return temp, nil
}

func ParseSplitSettings(cmd *cobra.Command) (settings common.SplitSettings, err error) {
	ratioArg, _ := cmd.Flags().GetString("ratio")
	re := regexp.MustCompile(`(\d+)[/|:](\d+)`)
	match := re.FindSubmatch([]byte(ratioArg))
	if len(match) != 3 {
		return settings, fmt.Errorf("Unable to parse `ratio` argument: '%s'", ratioArg)
	}
	a, err := strconv.ParseInt(string(match[1]), 10, 32)
	if err != nil {
		return settings, err
	}
	b, err := strconv.ParseInt(string(match[2]), 10, 32)
	if err != nil {
		return settings, err
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
		return settings, fmt.Errorf("parts and threshold must be 2 < x < 256")
	}
	settings.Parts = int(parts)
	settings.Threshold = int(threshold)
	return settings, nil
}

func PasscruxSplit(secret []byte, splitSettings common.SplitSettings, formatSettings common.FormatSettings) (string, error) {
	vprint.Print("Splittings: \n  ", splitSettings, " : \n")

	vprint.Printf("Len in bytes: %d\n", len(secret))
	shards, err := shamir.Split(secret, splitSettings.Parts, splitSettings.Threshold)
	common.LogIfFatal(err)
	vprint.Printf("Len of each shard in bytes: %d\n", len(shards[0]))
	vprint.Print("Output:\n")
	stringShards := common.FormatShards(shards, formatSettings)
	return strings.Join(stringShards, formatSettings.RecordSep), nil
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

		usePass, _ := cmd.Flags().GetBool("pass")
		splitSettings, err := ParseSplitSettings(cmd)
		common.LogIfFatal(err)
		formatSettings, err := common.ParseFormatSettings(cmd)
		common.LogIfFatal(err)

		secretString, err := getKeyData(args, false, usePass)
		common.LogIfFatal(err)
		vprint.Printf("Input: %s\n", secretString)

		secret := []byte(secretString)
		out, err := PasscruxSplit(secret, splitSettings, formatSettings)
		common.LogIfFatal(err)

		fmt.Println(out)
	},
}

func init() {
	RootCmd.AddCommand(splitCmd)
	RootCmd.PersistentFlags().StringP("ratio", "r", "3/5", "Ratio of parts needed to total parts")

}
