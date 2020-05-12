package common

import (
	b85 "encoding/ascii85"
	b32 "encoding/base32"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/xkortex/passcrux/common/abc16"
	"regexp"
	"strings"
)

// Ascii85 is still in development!
func A85EncodeToString(src []byte) string {
	dst := make([]byte, len(src)*2) // todo: make more efficient spacewise
	_ = b85.Encode(dst, src)
	return string(dst)
}

// Ascii85 is still in development!
func A85DecodeString(src string) ([]byte, error) {
	dst := make([]byte, len(src)) // todo: make more efficient spacewise
	n_dst, n_src, err := b85.Decode(dst, []byte(src), true)
	if err != nil {
		return nil, err
	}
	if n_dst == 0 || n_src == 0 {
		return nil, fmt.Errorf("Something unexpected happened")
	}
	return dst, nil
}

func getEncoder(encodingType string) (func(src []byte) string, error) {
	switch encodingType {
	case EncodeHex:
		return hex.EncodeToString, nil
	case EncodeBase32:
		return b32.StdEncoding.EncodeToString, nil
	case EncodeBase64:
		return b64.StdEncoding.EncodeToString, nil
	case EncodeBase85:
		return A85EncodeToString, nil
	case EncodeAbcAlt:
		return abc16.EncodeToStringAlt, nil
	case EncodeABC:
		return abc16.EncodeToString, nil
	default:
		return nil, fmt.Errorf("Invalid encodingType: %s", encodingType)
	}
}

// Get a string-to-bytes decoder function of the given formatting
func getDecoder(encodingType string) (func(src string) ([]byte, error), error) {
	switch encodingType {
	case EncodeHex:
		return hex.DecodeString, nil
	case EncodeBase32:
		return b32.StdEncoding.DecodeString, nil
	case EncodeBase64:
		return b64.StdEncoding.DecodeString, nil
	case EncodeBase85:
		return A85DecodeString, nil
	case EncodeABC:
		fallthrough
	case EncodeAbcAlt:
		return abc16.DecodeString, nil
	default:
		return nil, fmt.Errorf("Invalid encodingType: %s", encodingType)
	}
}

func Chunk(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	length := 0
	for _, r := range s {
		chunk[length] = r
		length++
		if length == chunkSize {
			chunks = append(chunks, string(chunk))
			length = 0
		}
	}
	if length > 0 {
		chunks = append(chunks, string(chunk[:length]))
	}
	return chunks
}

// Remove the field separators from the shard
func StripSep(s string, sep string) string {
	// todo: infer separator
	if sep == "" {
		return s
	}
	reSep := regexp.MustCompile(sep)
	return reSep.ReplaceAllString(s, "")
}

func FormatShards(shards [][]byte, settings FormatSettings) []string {
	stringShards := make([]string, len(shards), len(shards))
	encodeToString, err := getEncoder(settings.EncodingType)
	if err != nil {
		LogIfFatal(err)
	}
	for i, shard := range shards {
		stringShard := encodeToString(shard)
		chunks := Chunk(stringShard, settings.FieldSize)
		stringShards[i] = strings.Join(chunks, settings.Sep)
	}
	return stringShards
}

func DecodeShards(shardWords []string, settings FormatSettings) ([][]byte, error) {
	shardBytes := make([][]byte, len(shardWords), len(shardWords))
	var err error
	decodeToString, err := getDecoder(settings.EncodingType)
	if err != nil {
		return nil, err
	}
	for i, shard := range shardWords {
		shardClean := StripSep(shard, settings.Sep)
		shardByte, err := decodeToString(shardClean)
		if err != nil {
			return nil, err
		}
		shardBytes[i] = shardByte
	}
	return shardBytes, nil
}
