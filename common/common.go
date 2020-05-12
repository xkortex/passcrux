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
	EncodeBase85 = "base85"
)

type SplitSettings struct {
	Parts     int // total number of shards to split into
	Threshold int // number of shards required to reconstruct
}
type FormatSettings struct {
	EncodingType string // binary-to-string encoding, e.g. hex, base32
	Sep          string // separator between bytes
	RecordSep    string // separator between records/shards
	FieldSize    int    // size of each field, eg 2-> DE:AD:BE:EF
	FilenamePat  string // pattern (in typical sprintf notation) for filenames
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

func ParseFormatSettings(cmd *cobra.Command) (settings FormatSettings, err error) {
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

	if !ok {
		return settings, fmt.Errorf("Not a valid value for param `enc`: %s", encodingArg)
	}
	settings.EncodingType = val
	sep, err := cmd.Flags().GetString("sep")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse argument `sep`, defaulting to blank string")
	} else {
		settings.Sep = sep
	}
	recordSep, err := cmd.Flags().GetString("recordsep")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse argument `recordsep`, defaulting to newline")
	} else {
		settings.RecordSep = recordSep
	}

	fieldSize, err := cmd.Flags().GetInt("fieldsize")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse argument `fieldsize`, defaulting to 2")
	} else {
		settings.FieldSize = fieldSize
	}

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

	// Deliberately block until EOF, streaming doesn't really make sense with this app
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

// Read password from the terminal. Prompts user 2x for consistency.
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
