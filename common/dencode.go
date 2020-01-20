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
func getDecoder(encodingType string) (func(src string) ([]byte, error), error) {
	switch encodingType {
	case EncodeHex:
		return hex.DecodeString, nil
	case EncodeBase32:
		return b32.StdEncoding.DecodeString, nil
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

func DecodeShards(shardWords []string, settings FormatSettings) ([][]byte, error) {
	shardBytes := make([][]byte, len(shardWords), len(shardWords))
	var err error
	decodeToString, err := getDecoder(settings.EncodingType)
	if err != nil {
		return nil, err
	}
	for i, shard := range shardWords {
		shardByte, err := decodeToString(shard)
		if err != nil {
			return nil, err
		}
		shardBytes[i] = shardByte
	}
	return shardBytes, nil
}
