package gtc

import (
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	prevBlockLen  = 32
	merkleRootLen = 32
	prevHashLen   = 32
)

// read just makes calling Read a little shorter.
func read(r io.Reader, data interface{}) error {
	return errors.Wrap(binary.Read(r, binary.LittleEndian, data), "read failed")
}

// readVarInt reads variable length ints:
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func readVarInt(r io.Reader, varint *uint64) error {
	var pre uint8
	if err := read(r, &pre); err != nil {
		return errors.Wrap(err, "read failed")
	}

	switch pre {
	case 0xFD:
		var varint16 uint16
		if err := read(r, &varint16); err != nil {
			return errors.Wrap(err, "read for 0xFD prefix failed")
		}
		*varint = uint64(varint16)
	case 0xFE:
		var varint32 uint32
		if err := read(r, &varint32); err != nil {
			return errors.Wrap(err, "read for 0xFE prefix failed")
		}
		*varint = uint64(varint32)
	case 0xFF:
		if err := read(r, varint); err != nil {
			return errors.Wrap(err, "read for 0xFF prefix failed")
		}
	default:
		*varint = uint64(pre)
	}

	return nil
}

// Hash hashes a byte slice  using SHA256.
func Hash(bs []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(bs)
	return h.Sum(nil)
}
