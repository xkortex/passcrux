package common

type EncodingType string

const (
	EncodeHex    = "hex"
	EncodeBase32 = "base32"
)

type SplitSettings struct {
	Parts     int // total number of shards to split into
	Threshold int // number of shards required to reconstruct
}
type FormatSettings struct {
	EncodingType string
}
