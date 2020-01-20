package common

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"regexp"
)

type EncodingType string

const (
	EncodeRaw    = "raw"
	EncodeHex    = "hex"
	EncodeBase32 = "base32"
	EncodeBase64 = "base64"
)

type SplitSettings struct {
	Parts     int // total number of shards to split into
	Threshold int // number of shards required to reconstruct
}
type FormatSettings struct {
	EncodingType string
}

type StdInContainer struct {
	Stdin     string
	Has_stdin bool
}

func LogIfFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ParseFormatSettings(settings *FormatSettings, cmd *cobra.Command) error {
	encodingArg, _ := cmd.Flags().GetString("enc")
	re := regexp.MustCompile(`[base]*(\d+)|(h)|(hex)|(x)|(r)aw`)
	match := re.FindSubmatch([]byte(encodingArg))
	var encoding2 string
	if len(match) > 0 {
		for _, m := range match[1:] {
			if len(m) != 0 {
				encoding2 = string(m)
			}
		}
	} else {
		encoding2 = encodingArg
	}
	val, ok := map[string]string{
		"r":   EncodeRaw,
		"hex": EncodeHex,
		"h":   EncodeHex,
		"x":   EncodeHex,
		"32":  EncodeBase32,
		"64":  EncodeBase64,
	}[encoding2]
	if !ok {
		return fmt.Errorf("Not a valid value for param `enc`: %s", encodingArg)
	}
	settings.EncodingType = val
	return nil
}

func Get_stdin() (StdInContainer, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return StdInContainer{}, err
	}
	out_struct := StdInContainer{Has_stdin: false}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		//fmt.Println("Stdin is from a terminal")
		return out_struct, nil
	}

	// data is being piped to Stdin
	//fmt.Println("data is being piped to Stdin")

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	out_struct.Stdin = string(output)
	out_struct.Has_stdin = true
	return out_struct, nil
}
