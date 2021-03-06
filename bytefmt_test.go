package bytefmt

import (
	"bytes"
	"testing"
)

var tests = []struct {
	buf    []byte
	fmt    string
	expect string
}{
	{},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%%", "%"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%!", "%%UNKOWN%!"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "test", "test"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%p", "00000000  01 02 03 04                                       |....|\n"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%4p%p", "00000000  01 02 03 04                                       |....|\n"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%q", `"\x01\x02\x03\x04"`},
	{[]byte{'H', 'e', 'l', 'l', 'o'}, "%s", "Hello"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%4x", "1020304"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%-4x", "4030201"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%4d", "16909060"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%-4d", "67305985"},
	{[]byte{0x1, 0x2, 0x3, 0x4}, "%4b", "1000000100000001100000100"},
}

func TestSprintf(t *testing.T) {
	for _, tt := range tests {
		res := Sprintf(tt.buf, tt.fmt)
		if res != tt.expect {
			t.Logf("format %q: expected %q, res %q", tt.fmt, tt.expect, res)
			t.Fail()
		}
	}
}

func TestFprintf(t *testing.T) {
	for _, tt := range tests {
		var buf bytes.Buffer
		Fprintf(&buf, tt.buf, tt.fmt)
		res := buf.String()
		if res != tt.expect {
			t.Logf("format %q: expected %q, res %q", tt.fmt, tt.expect, res)
			t.Fail()
		}
	}
}

func TestEnum(t *testing.T) {
	var enumValues = map[int64]string{
		1: "One",
		2: "Two",
		3: "Three",
	}
	res := Sprintf([]byte{0x1, 0x2, 0x3, 0x4, 0x05}, "%1.0e, %1.0e, %1.0e, %1.0e, %1e", enumValues)
	expected := "One, Two, Three, 4, 5"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
}

func TestTemplate(t *testing.T) {
	var templates = map[int64]string{
		1: "%1x",
		2: "%1d",
		3: "%1b",
	}
	res := Sprintf([]byte{0x1, 0xee, 0x2, 0xaa, 0x3, 0x55, 0x4, 0x88}, "%1.0t, %1.0t, %1.0t, %1.0t, %1.0t%p", templates)
	expected := "ee, 170, 1010101, 4, 136"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
	res = Sprintf([]byte{0x1, 0xee, 0x2, 0xaa, 0x3, 0x55, 0x4, 0x88}, "%t, %2t", templates)
	expected = "32375466, 853"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
}

func TestScaled(t *testing.T) {
	var scale = 1e-6
	res := Sprintf([]byte{0x03, 0x47, 0x3b, 0xc0}, "%4.0i", scale)
	expected := "55"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
	res = Sprintf([]byte{0x03, 0x47, 0x3b, 0xc0}, "%4i")
	expected = "5.5e+07"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
}

func TestFlags(t *testing.T) {
	var flags = map[int64]string{
		0x80: "bit7",
		0x01: "bit0",
	}
	res := Sprintf([]byte{0x83}, "%1.0b", flags)
	expected := "(bit7|bit0|0x2)"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
	res = Sprintf([]byte{0x00}, "%1.0b", flags)
	expected = "()"
	if res != expected {
		t.Logf("enum expected %q, res %q", expected, res)
		t.Fail()
	}
}
