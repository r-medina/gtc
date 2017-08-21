package gtc_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/r-medina/gtc"
	"github.com/r-medina/gtc/testutil"
)

func Example() {
	block, err := gtc.DecodeBlock(testutil.R())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", block)

	// Output: Version: 1
	// Previous Block: 0000000000000000000000000000000000000000000000000000000000000000
	// MerkleRoot: 3ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a
	// Timestamp: 2009-01-03 13:15:05 -0500 EST
	// Bits: 486604799
	// Nonce: 2083236893
	// Transaction Count: 1
	// Transactions:
	// 	{Version:1 InTxCount:1 Inputs:[Previous Hash: 0000000000000000000000000000000000000000000000000000000000000000, Script Length: 77, Script: 04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73, Sequence: 4294967295] OutTxCount:1 Outputs:[Value: 5000000000, Pk Script Length: 67, Pk Script: 4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac] LockTime:0}
}

func ExampleJSON() {
	// JSON encoding isn't the most useful since it encodes the byte slices
	// differently, but it's helpful for debugging some of the other
	// headers.

	block, err := gtc.DecodeBlock(testutil.R())
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(block); err != nil {
		panic(err)
	}

	// Output: {
	//   "ver": 1,
	//   "prev_block": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
	//   "mrkl_root": "O6Pt/Xp7ErJ6xyw+Z3aPYX/IG8OIilEyOp+4qkseXko=",
	//   "time": 1231006505,
	//   "bits": 486604799,
	//   "nonce": 2083236893,
	//   "n_tx": 1,
	//   "tx": [
	//     {
	//       "ver": 1,
	//       "n_in": 1,
	//       "in": [
	//         {
	//           "prev_out": {
	//             "hash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
	//             "index": 4294967295
	//           },
	//           "script_len": 77,
	//           "script": "BP//AB0BBEVUaGUgVGltZXMgMDMvSmFuLzIwMDkgQ2hhbmNlbGxvciBvbiBicmluayBvZiBzZWNvbmQgYmFpbG91dCBmb3IgYmFua3M=",
	//           "seq": 4294967295
	//         }
	//       ],
	//       "n_out": 1,
	//       "out": [
	//         {
	//           "Value": 5000000000,
	//           "PkScriptLength": 67,
	//           "PkScript": "QQRniv2w/lVIJxln8aZxMLcQXNaoKOA5CaZ5YuDqH2Hetkn2vD9M7zjE81UE5R7BEt5cOE33uguNV4pMcCtr8R1frA=="
	//         }
	//       ],
	//       "lock_time": 0
	//     }
	//   ]
	// }
}
