package gtc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Transaction represents a transaction in a block.
type Transaction struct {
	Version int32 `json:"ver"`

	InTxCount uint64  `json:"n_in"`
	Inputs    []*InTx `json:"in"`

	OutTxCount uint64   `json:"n_out"`
	Outputs    []*OutTx `json:"out"`

	LockTime uint32 `json:"lock_time"`
}

// InTx represents the input transactions for a transaction.
type InTx struct {
	PreviousOutput struct {
		Hash  []byte `json:"hash"`
		Index uint32 `json:"index"`
	} `json:"prev_out"`
	ScriptLength uint64 `json:"script_len"`
	Script       []byte `json:"script"`
	Sequence     uint32 `json:"seq"`
}

func (tx InTx) String() string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Previous Hash: %x", tx.PreviousOutput.Hash)

	fmt.Fprintf(buf, ", Script Length: %d", tx.ScriptLength)
	fmt.Fprintf(buf, ", Script: %x", tx.Script)
	fmt.Fprintf(buf, ", Sequence: %d", tx.Sequence)

	return buf.String()
}

// OutTx represents the outputs transactions for a transaction.
type OutTx struct {
	Value          int64
	PkScriptLength uint64
	PkScript       []byte
}

func (tx OutTx) String() string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Value: %d", tx.Value)
	fmt.Fprintf(buf, ", Pk Script Length: %d", tx.PkScriptLength)
	fmt.Fprintf(buf, ", Pk Script: %x", tx.PkScript)

	return buf.String()
}

// decodeTransactions decodes the list of transactoins associated with a block.
func decodeTransactions(r io.Reader, n uint64) ([]*Transaction, error) {
	txs := make([]*Transaction, n)
	var i uint64
	for i = 0; i < n; i++ {
		tx, err := DecodeTransaction(r)
		if err != nil {
			return nil, errors.Wrap(err, "decodeTransaction failed")
		}
		txs[i] = tx
	}

	return txs, nil
}

// DecodeTransaction decodes a transaction.
func DecodeTransaction(r io.Reader) (*Transaction, error) {
	d := newTxDecoder(r)
	if err := d.decode(); err != nil {
		return nil, errors.Wrap(err, "transaction devode failed")
	}

	return &d.tx, nil
}

type txDecoder struct {
	r  io.Reader
	tx Transaction
}

func newTxDecoder(r io.Reader) *txDecoder {
	return &txDecoder{r: r}
}

func (d *txDecoder) decode() error {
	// read version
	if err := read(d.r, &d.tx.Version); err != nil {
		return errors.Wrap(err, "reading version failed")
	}

	// read input transaction count
	var n uint64
	if err := readVarInt(d.r, &n); err != nil {
		return errors.Wrap(err, "reading tx_in count failed")
	}
	d.tx.InTxCount = n

	// read in-transactions

	d.tx.Inputs = make([]*InTx, n)
	var i uint64
	for i = 0; i < n; i++ {
		tx, err := DecodeInTx(d.r)
		if err != nil {
			return errors.Wrap(err, "decodeInTx failed")
		}
		d.tx.Inputs[i] = tx
	}

	// read output transaction count
	if err := readVarInt(d.r, &n); err != nil { // reuses n
		return errors.Wrap(err, "reading tx_out count failed")
	}
	d.tx.OutTxCount = n

	// read out-transactions

	d.tx.Outputs = make([]*OutTx, n)
	for i = 0; i < n; i++ { // reuses i
		tx, err := DecodeOutTx(d.r)
		if err != nil {
			return errors.Wrap(err, "decodeOutTx failed")
		}
		d.tx.Outputs[i] = tx
	}

	// read lock time
	if err := read(d.r, &d.tx.LockTime); err != nil {
		return errors.Wrap(err, "reading lock_time failed")
	}

	return nil
}

// DecodeInTx decodes an input transaction.
func DecodeInTx(r io.Reader) (*InTx, error) {
	tx := InTx{}

	// read previous output
	tx.PreviousOutput.Hash = make([]byte, 32)
	if err := read(r, tx.PreviousOutput.Hash); err != nil {
		return nil, errors.Wrap(err, "reading previous output hash failed")
	}
	if err := read(r, &tx.PreviousOutput.Index); err != nil {
		return nil, errors.Wrap(err, "reading previous output index failed")
	}

	// read script length
	var n uint64
	if err := readVarInt(r, &n); err != nil {
		return nil, errors.Wrap(err, "reading script length failed")
	}
	tx.ScriptLength = n

	// read script
	tx.Script = make([]byte, n)
	if err := read(r, tx.Script); err != nil {
		return nil, errors.Wrap(err, "reading script failed")
	}

	// read sequence
	if err := read(r, &tx.Sequence); err != nil {
		return nil, errors.Wrap(err, "reading sequence failed")
	}

	return &tx, nil
}

// DecodeOutTx decodes an output transaction.
func DecodeOutTx(r io.Reader) (*OutTx, error) {
	tx := OutTx{}

	// read value
	if err := read(r, &tx.Value); err != nil {
		return nil, err
	}

	// read pk script length
	var n uint64
	if err := readVarInt(r, &n); err != nil {
		return nil, errors.Wrap(err, "reading script length failed")
	}
	tx.PkScriptLength = n

	// read pk script
	tx.PkScript = make([]byte, n)
	if err := read(r, tx.PkScript); err != nil {
		return nil, errors.Wrap(err, "reading pk script failed")
	}

	return &tx, nil
}
