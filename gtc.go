package gtc

import (
	"encoding/binary"
	"io"
)

const (
	prevBlockLen   = 32
	merkleRootLen  = 32
	prevHashLenLen = 32
)

// Decode takes a reader that has a block chain and decodes it.
func Decode(r io.Reader) ([]*Block, error) {
	blocks := []*Block{}
	for {
		block, err := DecodeBlock(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return blocks, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// DecodeBlock decodes a single block.
func DecodeBlock(r io.Reader) (*Block, error) {
	d := newBlockDecoder(r)
	if err := d.decode(); err != nil {
		return nil, err
	}

	return &d.b, nil
}

func decodeTransactions(r io.Reader, n int) ([]*Transaction, error) {
	txs := make([]*Transaction, n)
	for i := range txs {
		tx, err := decodeTransaction(r)
		if err != nil {
			return nil, err
		}
		txs[i] = tx
	}

	return txs, nil
}

func decodeTransaction(r io.Reader) (*Transaction, error) {
	d := newTxDecoder(r)
	if err := d.decode(); err != nil {
		return nil, err
	}

	return &d.tx, nil
}

// read just makes calling Read a little shorter.
func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.LittleEndian, data)
}

// readVarInt reads variable length ints:
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func readVarInt(r io.Reader, varint *uint64) error {
	var pre uint8
	if err := read(r, &pre); err != nil {
		return err
	}

	switch pre {
	case 0xFD:
		var varint16 uint16
		if err := binary.Read(r, binary.LittleEndian, &varint16); err != nil {
			return err
		}
		*varint = uint64(varint16)
	case 0xFE:
		var varint32 uint32
		if err := binary.Read(r, binary.LittleEndian, &varint32); err != nil {
			return err
		}
		*varint = uint64(varint32)
	case 0xFF:
		if err := binary.Read(r, binary.LittleEndian, varint); err != nil {
			return err
		}
	default:
		*varint = uint64(pre)
	}

	return nil
}
