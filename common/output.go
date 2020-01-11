package common

import (
	b32 "encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
)

func getEncoder(encodingType string) (func(src []byte) string, error) {
	switch encodingType {
	case EncodeHex:
		return hex.EncodeToString, nil
	case EncodeBase32:
		return b32.StdEncoding.EncodeToString, nil
	default:
		return nil, fmt.Errorf("Invalid encodingType: %s", encodingType)
	}
}

func FormatShards(shards [][]byte, settings FormatSettings) (string, error) {
	stringShards := make([]string, len(shards), len(shards))
	encodeToString, err := getEncoder(settings.EncodingType)
	if err != nil {
		return "", err
	}
	for i, shard := range shards {
		stringShards[i] = encodeToString(shard)
	}
	return strings.Join(stringShards, "\n"), nil
}
