package gtc

import (
	"bytes"
	"fmt"
	"io"
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
		return err
	}

	// read input transaction count
	var n uint64
	if err := readVarInt(d.r, &n); err != nil {
		return err
	}
	d.tx.InTxCount = n

	// read in-transactions

	d.tx.Inputs = make([]*InTx, n)
	for i := range d.tx.Inputs {
		tx := InTx{}

		// read previous output
		tx.PreviousOutput.Hash = make([]byte, 32)
		if err := read(d.r, tx.PreviousOutput.Hash); err != nil {
			return err
		}
		if err := read(d.r, &tx.PreviousOutput.Index); err != nil {
			return err
		}

		// read script length
		var n uint64 // shadows n
		if err := readVarInt(d.r, &n); err != nil {
			return err
		}
		tx.ScriptLength = n

		// read script
		tx.Script = make([]byte, n)
		if err := read(d.r, tx.Script); err != nil {
			return err
		}

		// read sequence
		if err := read(d.r, &tx.Sequence); err != nil {
			return err
		}

		d.tx.Inputs[i] = &tx
	}

	// read output transaction count
	if err := readVarInt(d.r, &n); err != nil { // reuses n
		return err
	}
	d.tx.OutTxCount = n

	// read out-transactions

	d.tx.Outputs = make([]*OutTx, n)
	for i := range d.tx.Outputs {
		tx := OutTx{}

		// read value
		if err := read(d.r, &tx.Value); err != nil {
			return err
		}

		// read pk script length
		var n uint64 // shadows n
		if err := readVarInt(d.r, &n); err != nil {
			return err
		}
		tx.PkScriptLength = n

		// read pk script
		tx.PkScript = make([]byte, n)
		if err := read(d.r, tx.PkScript); err != nil {
			return err
		}

		d.tx.Outputs[i] = &tx
	}

	// read lock time
	if err := read(d.r, &d.tx.LockTime); err != nil {
		return err
	}

	return nil
}
