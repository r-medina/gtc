package gtc

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Block represents a block in the block chain
// (with the format of the genesis block).
//
// json formatting somewhat inspired by http://blockchain.info API.
type Block struct {
	Version       int32  `json:"ver"`
	PreviousBlock []byte `json:"prev_block"`
	MerkleRoot    []byte `json:"mrkl_root"`
	Timestamp     uint32 `json:"time"`
	Bits          uint32 `json:"bits"`
	Nonce         uint32 `json:"nonce"`

	TxCount      uint64         `json:"n_tx"`
	Transactions []*Transaction `json:"tx"`
}

func (b Block) String() string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Version: %d", b.Version)
	fmt.Fprintf(buf, "\nPrevious Block: %x", b.PreviousBlock)
	fmt.Fprintf(buf, "\nMerkleRoot: %x", b.MerkleRoot)
	fmt.Fprintf(buf, "\nTimestamp: %v", time.Unix(int64(b.Timestamp), 0))
	fmt.Fprintf(buf, "\nBits: %d", b.Bits)
	fmt.Fprintf(buf, "\nNonce: %d", b.Nonce)
	fmt.Fprintf(buf, "\nTransaction Count: %d", b.TxCount)

	fmt.Fprintf(buf, "\nTransactions:")
	for _, tx := range b.Transactions {
		fmt.Fprintf(buf, "\n\t%+v", *tx)
	}

	return buf.String()
}

// DecodeBlock decodes a single block.
func DecodeBlock(r io.Reader) (*Block, error) {
	d := newBlockDecoder(r)
	if err := d.decode(); err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}

	return &d.b, nil
}

type blockDecoder struct {
	r io.Reader
	b Block
}

func newBlockDecoder(r io.Reader) *blockDecoder {
	return &blockDecoder{r: r}
}

func (d *blockDecoder) decode() error {
	// read version
	if err := read(d.r, &d.b.Version); err != nil {
		return errors.Wrap(err, "reading block version failed")
	}

	// read prev block hash
	d.b.PreviousBlock = make([]byte, prevBlockLen)
	if err := read(d.r, d.b.PreviousBlock); err != nil {
		return errors.Wrap(err, "reading block hash failed")
	}

	// read merkle root hash
	d.b.MerkleRoot = make([]byte, merkleRootLen)
	if err := read(d.r, d.b.MerkleRoot); err != nil {
		return errors.Wrap(err, "reading merkle root hash")
	}

	// read timestamp
	if err := read(d.r, &d.b.Timestamp); err != nil {
		return errors.Wrap(err, "reading timestamp failed")
	}

	// read bits
	if err := read(d.r, &d.b.Bits); err != nil {
		return errors.Wrap(err, "reading bits failed")
	}

	// read nonce
	if err := read(d.r, &d.b.Nonce); err != nil {
		return errors.Wrap(err, "reading nonce failed")
	}

	// read transaction count
	if err := readVarInt(d.r, &d.b.TxCount); err != nil {
		return errors.Wrap(err, "reading transaction count failed")
	}

	txs, err := decodeTransactions(d.r, d.b.TxCount)
	if err != nil {
		return errors.Wrap(err, "decodeTransactions failed")
	}
	d.b.Transactions = txs

	return nil
}
