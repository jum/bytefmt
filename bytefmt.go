/*
Package bytefmt provides a simple and quick way to dump byte slices in a way
similar to the package fmt. But opposed to fmt, each % format specifier does
consume bytes from the byte slice buf passed as an arg. The width part of a
format verb specifies the number of bytes to consume of an array. The following
format letters are understood:

	%p	hex dump bytes using encoding/hex.Dump
	%q  print a go quoted string
	%s  print a string
	%d	print a decimal int (max width 8)
	%x	print hex int (max width 8)
	%b	print binary int (max width 8)
	%e	print enumerated type, precision field is argument index

	The %x and %d formats can be modified to use intel byte order using a
	leading ´-´ sign in the width field (e.g. %-4d).
*/
package bytefmt

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strconv"
)

type dumper struct {
	input      []byte
	ii         int
	prec       int
	precValid  bool
	width      int
	widthValid bool
	intel      bool // intel byte order for multibyte ints
	buf        bytes.Buffer
}

// A lot of the logic of this is copied from the fmt package.
func (d *dumper) doDump(buf []byte, fmt string, a []interface{}) {
	d.input = buf
	end := len(fmt)
	//formatLoop:
	for i := 0; i < end; {
		lasti := i
		for i < end && fmt[i] != '%' {
			i++
		}
		if i > lasti {
			d.buf.WriteString(fmt[lasti:i])
		}
		i++
		if i >= end {
			break
		}
		c := fmt[i]
		d.intel = false
		d.precValid = false
		d.widthValid = false
		d.width = 0
		d.prec = 0
		if c == '-' {
			d.intel = true
			i++
			if i >= end {
				break
			}
			c = fmt[i]
		}
		if c >= '0' && c <= '9' {
			d.width, d.widthValid, i = parsenum(fmt, i, end)
			if i >= end {
				break
			}
			c = fmt[i]
		}
		if c == '.' {
			i++
			if i >= end {
				break
			}
			d.prec, d.precValid, i = parsenum(fmt, i, end)
			if i >= end {
				break
			}
			c = fmt[i]
		}
		i++
		switch c {
		case '%':
			d.buf.WriteRune('%')
		case 'p':
			if !d.widthValid {
				d.width = len(d.input) - d.ii
			}
			d.buf.WriteString(hex.Dump(d.input[d.ii : d.ii+d.width]))
			d.ii += d.width
		case 'q':
			if !d.widthValid {
				d.width = len(d.input) - d.ii
			}
			d.buf.WriteString(strconv.Quote(string(d.input[d.ii : d.ii+d.width])))
			d.ii += d.width
		case 's':
			if !d.widthValid {
				d.width = len(d.input) - d.ii
			}
			d.buf.WriteString(string(d.input[d.ii : d.ii+d.width]))
			d.ii += d.width
		case 'x':
			if !d.widthValid {
				d.width = 4
			}
			x := d.fetchInt()
			d.buf.WriteString(strconv.FormatInt(x, 16))
		case 'd':
			if !d.widthValid {
				d.width = 4
			}
			x := d.fetchInt()
			d.buf.WriteString(strconv.FormatInt(x, 10))
		case 'b':
			if !d.widthValid {
				d.width = 4
			}
			x := d.fetchInt()
			d.buf.WriteString(strconv.FormatInt(x, 2))
		case 'e':
			if !d.widthValid {
				d.width = 4
			}
			x := d.fetchInt()
			if d.precValid {
				m := a[d.prec].(map[int64]string)
				if s, ok := m[x]; ok {
					d.buf.WriteString(s)
				} else {
					d.buf.WriteString(strconv.FormatInt(x, 10))
				}
			} else {
				d.buf.WriteString(strconv.FormatInt(x, 10))
			}
		default:
			d.buf.WriteString("%%UNKOWN%" + string(c))
		}
	}
}

func (d *dumper) fetchInt() int64 {
	var val int64
	if d.intel {
		for w := d.width; w > 0; w-- {
			val |= int64(d.input[d.ii]) << uint((d.width-w)*8)
			d.ii++
		}
	} else {
		for w := d.width; w > 0; w-- {
			val <<= 8
			val |= int64(d.input[d.ii])
			d.ii++
		}
	}
	return val
}

// parsenum converts ASCII to integer.  num is 0 (and isnum is false) if no number present.
func parsenum(s string, start, end int) (num int, isnum bool, newi int) {
	if start >= end {
		return 0, false, end
	}
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		if tooLarge(num) {
			return 0, false, end // Overflow; crazy long number most likely.
		}
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return
}

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	const max int = 1e6
	return x > max || x < -max
}

// Fprintf dumps to the writer w.
func Fprintf(w io.Writer, buf []byte, fmt string, a ...interface{}) (n int, err error) {
	var d dumper
	d.doDump(buf, fmt, a)
	n, err = w.Write(d.buf.Bytes())
	return
}

// Printf dumps to stdout.
func Printf(buf []byte, fmt string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, buf, fmt, a...)
}

// Sprintf dumps to a string.
func Sprintf(buf []byte, fmt string, a ...interface{}) string {
	var d dumper
	d.doDump(buf, fmt, a)
	return d.buf.String()
}
