// ABC base-16 encoding
// It's just hex encoding but with all alphabetic, no numerals
// Advantage is being able to do things like book steganography easily
// (numbers are much less frequent)
package abc16

/* Some musings on encodings:
There are several ways to pack N lower-density symbols S into M bytesize-symbols,
with three main performance parameters:
- Pack-out ratio (M bytes becomes N symbols, N/M is your ratio)
- Bit efficiency ( log2(256**M) / log2(S**N) ) or how much "wasted space" there is
- Padding ease (how tricky is it to convert blocks with odd numbers of bytes

There are basically two schemes for 24 bit:

2:1 - Ratio: 2.00  Eff: 87.2%  Padding: easy
7:4 - Ratio: 1.75  Eff: 99.7%  Padding: tricky

Pack-out and efficiency for 2:1 are about 15% worse than 7:4. But, the
program flow is way simpler with 2:1, since we can take each byte and
just expand it, without having to worry about block sizes or padding.

This project aims for ease and robustness over all else, so squeezing 15%
efficiency in exchange for complexity is not a good trade.

*/

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Encoding alphabet
// x is dropped, despite being one of my favorite letters, to avoid ambiguity with hex
// also J, Q, X and Z are least common letters in English
const encodeStd = "ABCDEFGHKMNRTUWY"
const encodeAlt = "aBcDeFgHkMnRtUwY"

var encodeSliceUpper = []byte(strings.ToUpper(encodeStd))
var encodeSliceLower = []byte(strings.ToLower(encodeStd))

// The rest is almost-verbatim the Go std hex encode lib, but with my custom alphabet

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 2.
func EncodedLen(n int) int { return n * 2 }

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements abc16 encoding.
func Encode(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = encodeStd[v>>4]
		dst[j+1] = encodeStd[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

// Encode encodes src into EncodedLen(len(src)), using alternating caps encoding
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements abc16 encoding.
func EncodeAlt(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = encodeAlt[v>>4]
		dst[j+1] = encodeAlt[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

// ErrLength reports an attempt to decode an odd-length input
// using Decode or DecodeString.
// The stream-based Decoder returns io.ErrUnexpectedEOF instead of ErrLength.
var ErrLength = errors.New("encoding/hex: odd length hex string")

// InvalidByteError values describe errors resulting from an invalid byte in a hex string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("encoding/hex: invalid byte: %#U (%d)", rune(e), e)
}

//var Error = func (e InvalidByteError) string {return fmt.Sprintf("encoding/hex: invalid byte: %#U", rune(e))}

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 2.
var DecodedLen = func(x int) int { return x / 2 }

// Decode decodes src into DecodedLen(len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// Decode expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, Decode returns the number
// of bytes decoded before the error.
func Decode(dst, src []byte) (int, error) {
	//var Decode = func (dst, src []byte) (int, error) {
	i, j := 0, 1
	for ; j < len(src); j += 2 {
		a, ok := fromHexChar(src[j-1 : j])
		if !ok {
			return i, InvalidByteError(src[j-1])
		}
		b, ok := fromHexChar(src[j : j+1])
		if !ok {
			return i, InvalidByteError(src[j])
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if _, ok := fromHexChar(src[j-1 : j]); !ok {
			return i, InvalidByteError(src[j-1])
		}
		return i, ErrLength
	}
	return i, nil
}

// fromHexChar converts a hex character into its value and a success flag.
func fromHexChar(c []byte) (byte, bool) {
	val := bytes.Index(encodeSliceUpper, c)
	if val >= 0 {
		return byte(val), true
	}
	val = bytes.Index(encodeSliceLower, c)
	if val >= 0 {
		return byte(val), true
	}
	return 0, false
}

// EncodeToString returns the hexadecimal encoding of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

// EncodeToString returns the hexadecimal encoding of src.
func EncodeToStringAlt(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	EncodeAlt(dst, src)
	return string(dst)
}

// DecodeString returns the bytes represented by the hexadecimal string s.
//
// DecodeString expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, DecodeString returns
// the bytes decoded before the error.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	// We can use the source slice itself as the destination
	// because the decode loop increments by one and then the 'seen' byte is not used anymore.
	n, err := Decode(src, src)
	return src[:n], err
}

// Dump returns a string that contains a hex dump of the given data. The format
// of the hex dump matches the output of `hexdump -C` on the command line.
func Dump(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var buf strings.Builder
	// Dumper will write 79 bytes per complete 16 byte chunk, and at least
	// 64 bytes for whatever remains. Round the allocation up, since only a
	// maximum of 15 bytes will be wasted.
	buf.Grow((1 + ((len(data) - 1) / 16)) * 79)

	dumper := Dumper(&buf)
	dumper.Write(data)
	dumper.Close()
	return buf.String()
}

// bufferSize is the number of hexadecimal characters to buffer in encoder and decoder.
const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns an io.Writer that writes lowercase hexadecimal characters to w.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / 2
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var written int
		encoded := Encode(e.out[:], p[:chunkSize])
		written, e.err = e.w.Write(e.out[:encoded])
		n += written / 2
		p = p[chunkSize:]
	}
	return n, e.err
}

type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

// NewDecoder returns an io.Reader that decodes hexadecimal characters from r.
// NewDecoder expects that r contain only an even number of hexadecimal characters.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < 2 && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in) // Copies either 0 or 1 bytes
		numRead, d.err = d.r.Read(d.arr[numCopy:])
		d.in = d.arr[:numCopy+numRead]
		if d.err == io.EOF && len(d.in)%2 != 0 {
			idx := len(d.in) - 1
			if _, ok := fromHexChar(d.in[idx : idx+1]); !ok {
				d.err = InvalidByteError(d.in[len(d.in)-1])
			} else {
				d.err = io.ErrUnexpectedEOF
			}
		}
	}

	// Decode internal buffer into output buffer
	if numAvail := len(d.in) / 2; len(p) > numAvail {
		p = p[:numAvail]
	}
	numDec, err := Decode(p, d.in[:len(p)*2])
	d.in = d.in[2*numDec:]
	if err != nil {
		d.in, d.err = nil, err // Decode error; discard input remainder
	}

	if len(d.in) < 2 {
		return numDec, d.err // Only expose errors when buffer fully consumed
	}
	return numDec, nil
}

// Dumper returns a WriteCloser that writes a hex dump of all written data to
// w. The format of the dump matches the output of `hexdump -C` on the command
// line.
func Dumper(w io.Writer) io.WriteCloser {
	return &dumper{w: w}
}

type dumper struct {
	w          io.Writer
	rightChars [18]byte
	buf        [14]byte
	used       int  // number of bytes in the current line
	n          uint // number of bytes, total
	closed     bool
}

func toChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}

func (h *dumper) Write(data []byte) (n int, err error) {
	if h.closed {
		return 0, errors.New("encoding/hex: dumper closed")
	}

	// Output lines look like:
	// 00000010  2e 2f 30 31 32 33 34 35  36 37 38 39 3a 3b 3c 3d  |./0123456789:;<=|
	// ^ offset                          ^ extra space              ^ ASCII of line.
	for i := range data {
		if h.used == 0 {
			// At the beginning of a line we print the current
			// offset in hex.
			h.buf[0] = byte(h.n >> 24)
			h.buf[1] = byte(h.n >> 16)
			h.buf[2] = byte(h.n >> 8)
			h.buf[3] = byte(h.n)
			Encode(h.buf[4:], h.buf[:4])
			h.buf[12] = ' '
			h.buf[13] = ' '
			_, err = h.w.Write(h.buf[4:])
			if err != nil {
				return
			}
		}
		Encode(h.buf[:], data[i:i+1])
		h.buf[2] = ' '
		l := 3
		if h.used == 7 {
			// There's an additional space after the 8th byte.
			h.buf[3] = ' '
			l = 4
		} else if h.used == 15 {
			// At the end of the line there's an extra space and
			// the bar for the right column.
			h.buf[3] = ' '
			h.buf[4] = '|'
			l = 5
		}
		_, err = h.w.Write(h.buf[:l])
		if err != nil {
			return
		}
		n++
		h.rightChars[h.used] = toChar(data[i])
		h.used++
		h.n++
		if h.used == 16 {
			h.rightChars[16] = '|'
			h.rightChars[17] = '\n'
			_, err = h.w.Write(h.rightChars[:])
			if err != nil {
				return
			}
			h.used = 0
		}
	}
	return
}

func (h *dumper) Close() (err error) {
	// See the comments in Write() for the details of this format.
	if h.closed {
		return
	}
	h.closed = true
	if h.used == 0 {
		return
	}
	h.buf[0] = ' '
	h.buf[1] = ' '
	h.buf[2] = ' '
	h.buf[3] = ' '
	h.buf[4] = '|'
	nBytes := h.used
	for h.used < 16 {
		l := 3
		if h.used == 7 {
			l = 4
		} else if h.used == 15 {
			l = 5
		}
		_, err = h.w.Write(h.buf[:l])
		if err != nil {
			return
		}
		h.used++
	}
	h.rightChars[nBytes] = '|'
	h.rightChars[nBytes+1] = '\n'
	_, err = h.w.Write(h.rightChars[:nBytes+2])
	return
}
