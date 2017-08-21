package gtc

import (
	"bytes"
	"testing"

	"github.com/r-medina/gtc/testutil"
)

// bc is a helper function that makes a byte slice out of a hex coded string or
// fails a test.
func bs(t *testing.T, str string) []byte {
	bs, err := testutil.HexStringToBytes(str)
	if err != nil {
		t.Fatalf("unexpected error converting hex value to bytes: %v", err)
	}
	return bs
}

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		val  []byte
		want uint64
	}{
		{val: []byte{0x12}, want: 0x12},
		{val: []byte{0xFD, 0x34, 0x12}, want: 0x1234},
		{val: []byte{0xFE, 0x01, 0x00, 0x00, 0x00}, want: 0x1},
		{val: func() []byte {
			bs := make([]byte, 16)
			bs[0], bs[1] = 0xFF, 0x01
			return bs
		}(), want: 0x1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			buf := bytes.NewBuffer(test.val)

			var got uint64
			if err := readVarInt(buf, &got); err != nil {
				t.Fatalf("failed to read varint: %v", err)
			}

			if want, got := test.want, got; got != want {
				t.Errorf("expectd %v, got %v", want, got)
			}
		})
	}
}
