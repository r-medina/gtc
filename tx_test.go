package gtc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecodeInTx(t *testing.T) {
	tests := []struct {
		bs   []byte
		want InTx
	}{{
		// from genesis block

		bs: bs(t, `0000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff`),
		want: InTx{
			PreviousOutput: Outpoint{
				Hash:  bs(t, `0000000000000000000000000000000000000000000000000000000000000000`),
				Index: 0xffffffff,
			},
			ScriptLength: 0x4d, // 77
			Script:       bs(t, `04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73`),
			Sequence:     0xffffffff,
		},
	}, {
		// from block 1

		bs: bs(t, `0000000000000000000000000000000000000000000000000000000000000000ffffffff0704ffff001d0104ffffffff`),
		want: InTx{
			PreviousOutput: Outpoint{
				Hash:  bs(t, `0000000000000000000000000000000000000000000000000000000000000000`),
				Index: 0xffffffff,
			},
			ScriptLength: 0x7, // 77
			Script:       bs(t, `04ffff001d0104`),
			Sequence:     0xffffffff,
		},
	}}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			r := bytes.NewBuffer(test.bs)
			tx, err := DecodeInTx(r)
			if err != nil {
				t.Fatalf("unexpected error decoding block: %v", err)
			}

			if want, got := test.want, *tx; !reflect.DeepEqual(got, want) {
				t.Fatalf("\nexpected (%T):\n%+v\ngot (%T):\n%+v", want, want, got, got)
			}
		})
	}
}

func TestDecodeOutTx(t *testing.T) {
	tests := []struct {
		bs   []byte
		want OutTx
	}{{
		bs: bs(t, `00f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac`),
		want: OutTx{
			Value:          5000000000,
			PkScriptLength: 0x43, // 67
			PkScript:       bs(t, `4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac`),
		},
	}, {
		bs: bs(t, `00f2052a0100000043410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac`),
		want: OutTx{
			Value:          5000000000,
			PkScriptLength: 0x43, // 67
			PkScript:       bs(t, `410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac`),
		},
	}}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			r := bytes.NewBuffer(test.bs)
			tx, err := DecodeOutTx(r)
			if err != nil {
				t.Fatalf("unexpected error decoding block: %v", err)
			}

			if want, got := test.want, *tx; !reflect.DeepEqual(got, want) {
				t.Fatalf("\nexpected (%T):\n%+v\ngot (%T):\n%+v", want, want, got, got)
			}
		})
	}
}

func TestTxDecode(t *testing.T) {}
