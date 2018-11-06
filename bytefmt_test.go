package bytefmt

import (
	"testing"
)

func TestSprintf(t *testing.T) {
	var tests = []struct {
		buf    []byte
		fmt    string
		expect string
	}{
		{},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%%", "%"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "test", "test"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%4b", "00000000  01 02 03 04                                       |....|\n"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%4q", `"\x01\x02\x03\x04"`},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%4x", "1020304"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%-4x", "4030201"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%4d", "16909060"},
		{[]byte{0x1, 0x2, 0x3, 0x4}, "%-4d", "67305985"},
	}
	for _, tt := range tests {
		res := Sprintf(tt.buf, tt.fmt)
		if res != tt.expect {
			t.Logf("format %q: expected %q, res %q", tt.fmt, tt.expect, res)
			t.Fail()
		}
	}
}
