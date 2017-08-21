package testutil

import (
	"bytes"
	"encoding/hex"
	"io"
)

var strData = []string{
	// genesis block
	// from https://blockchain.info/rawblock/000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f?format=hex
	`0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000`,

	// first block
	// from https://blockchain.info/rawblock/00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048?format=hex
	`010000006fe28c0ab6f1b372c1a6a246ae63f74f931e8365e15a089c68d6190000000000982051fd1e4ba744bbbe680e1fee14677ba1a3c3540bf7b1cdb606e857233e0e61bc6649ffff001d01e362990101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704ffff001d0104ffffffff0100f2052a0100000043410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac00000000`,
}

// Data is a slice of byte slices - each element in the outermost slice is the
// bytes of a block.
//
// Data contains the first two blocks in the BitCoin block chain.
var Data = func() [][]byte {
	data := [][]byte{}
	for _, blockStr := range strData {
		block, err := HexStringToBytes(blockStr)
		if err != nil {
			panic(err)
		}

		data = append(data, block)
	}

	return data
}()

// R returns an io.Reader with the data in Data. That is to say it contains
// the blocks. Helper function for getting all test data.
func R() io.Reader {
	return Reader(Data)
}

// Reader puts the byte slices in data in an io.Reader. Use this API to select
// specific blocks from Data.
func Reader(data [][]byte) io.Reader {
	buf := &bytes.Buffer{}
	for _, block := range data {
		buf.Write(block)
	}
	return buf
}

// HexStringToBytes takes a hex encoded string that represents an array of bytes
// and encodes them as a byte slice. Has many applications for testing since web
// UIs give the raw block data as strings.
func HexStringToBytes(str string) ([]byte, error) {
	n := hex.DecodedLen(len([]byte(str)))
	bs := make([]byte, n)
	_, err := hex.Decode(bs, []byte(str))
	if err != nil {
		return nil, err
	}

	return bs, nil
}
