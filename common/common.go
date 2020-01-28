package common

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xkortex/vprint"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"regexp"
)

type EncodingType string

const (
	cHidden   = "\033[8m"
	cUnhidden = "\033[8m"
	endColor  = "\033[0m"
)

const (
	EncodeRaw    = "raw"
	EncodeAbcAlt = "abcAlt"
	EncodeABC    = "ABC"
	EncodeHex    = "hex"
	EncodeBase32 = "base32"
	EncodeBase64 = "base64"
	EncodeBase85 = "base85/**/"
)

type SplitSettings struct {
	Parts     int // total number of shards to split into
	Threshold int // number of shards required to reconstruct
}
type FormatSettings struct {
	EncodingType string
	Sep          string
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

func ParseFormatSettings(cmd *cobra.Command) (FormatSettings, error) {
	encodingArg, _ := cmd.Flags().GetString("enc")
	vprint.Print("encodingArg: [", encodingArg, "]\n")
	re := regexp.MustCompile(`[base]*(\d+)|(h)|(hex)|(x)|(r)aw|(abc)|(ABC)`)
	match := re.FindSubmatch([]byte(encodingArg))
	var encodingParsed string = "none"
	if len(match) > 0 {
		for _, m := range match[1:] {
			if len(m) != 0 {
				encodingParsed = string(m)
			}
		}
	} else {
		encodingParsed = encodingArg
	}
	vprint.Print("encodingParsed: [", encodingParsed, "]\n")
	val, ok := map[string]string{
		"r":   EncodeRaw,
		"hex": EncodeHex,
		"h":   EncodeHex,
		"x":   EncodeHex,
		"32":  EncodeBase32,
		"64":  EncodeBase64,
		"85":  EncodeBase85,
		"abc": EncodeAbcAlt,
		"ABC": EncodeABC,
	}[encodingParsed]

	settings := FormatSettings{
		"",
		"",
	}
	if !ok {
		return settings, fmt.Errorf("Not a valid value for param `enc`: %s", encodingArg)
	}

	settings.EncodingType = val
	settings.Sep = ""
	return settings, nil
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

func ReadPassword() (string, error) {
	fmt.Fprintf(os.Stderr, "\nNow, please type in the password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintf(os.Stderr, "\n")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(os.Stderr, "Please type once more to confirm: ")
	password2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintf(os.Stderr, "\n")
	if err != nil {
		return "", err
	}

	if string(password) != string(password2) {
		return "", fmt.Errorf("Passwords do not match")
	}
	vprint.Printf("[%s]\n", password)
	fmt.Fprintf(os.Stderr, "___________\n")
	return string(password), nil
}
