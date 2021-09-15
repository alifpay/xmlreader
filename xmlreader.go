package xmlreader

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotFound = errors.New("not found")
)

type Decoder struct {
	rd    io.Reader
	Name  string
	Value string
	eof   bool
	end   bool
	bf    []byte
	bt    []byte
}

//new instance of xml reader
func New(r io.Reader) *Decoder {
	return &Decoder{rd: r, bf: make([]byte, 0, 512), bt: make([]byte, 1)}
}

//Read, reads the next node from the stream
// Returns:
//     true if the next node was read successfully; otherwise, false.
func (d *Decoder) Read() bool {
	if d.eof {
		return false
	}
	d.Name = ""
	d.end = false
	var node bool
	nm := make([]byte, 0)
	for {
		_, err := d.rd.Read(d.bt)
		if err == io.EOF {
			d.eof = true
			return false
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		d.bf = append(d.bf, d.bt...)

		switch d.bt[0] {
		case '<':
			node = true
		case ' ':
			if node {
				d.Name = string(nm)
				return true
			}
		case '>':
			d.end = true
			if node {
				d.Name = string(nm)
				return true
			}
		default:
			if node {
				nm = append(nm, d.bt...)
			}
		}
	}
}

// GetAttribute, gets the value of the attribute node
//   name:
//     The qualified name of the attribute.
//
// Returns:
//     The value of the specified attribute. If the attribute is not found or the value
//     empty.
func (d *Decoder) GetAttribute(name string) (string, error) {
	if d.eof {
		return "", io.EOF
	}
	var attr bool
	buf := make([]byte, 0)
	for {
		_, err := d.rd.Read(d.bt)
		if err == io.EOF {
			d.eof = true
			return "", err
		}
		if err != nil {
			return "", err
		}
		d.bf = append(d.bf, d.bt...)

		switch d.bt[0] {
		case '=':
			if name == string(buf) {
				attr = true
				buf = make([]byte, 0)
				continue
			}
		case '"':
			if attr && len(buf) > 0 {
				return string(buf), nil
			}
		case ' ':
			buf = make([]byte, 0)
			continue
		case '>':
			d.end = true
			return "", ErrNotFound
		default:
			buf = append(buf, d.bt...)
		}
	}
}

//HasValue, gets a value indicating whether the current node can have a value
// Returns:
//     true if the node on which the reader is currently positioned can have a Value;
//     otherwise, false.
func (d *Decoder) HasValue() bool {
	if d.eof {
		return false
	}
	d.Value = ""
	var has bool
	v := make([]byte, 0)
	if d.end {
		has = true
	}
	for {
		_, err := d.rd.Read(d.bt)
		if err == io.EOF {
			d.eof = true
			return false
		}
		if err != nil {
			return false
		}
		d.bf = append(d.bf, d.bt...)

		switch d.bt[0] {
		case '>':
			has = true
			v = make([]byte, 0)
			continue
		case '<':
			if ln := len(v); has && ln > 0 {
				sp := true
				ln--
				for sp {
					switch v[ln] {
					case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
						ln--
					default:
						sp = false
					}
				}
				d.Value = string(v[:ln+1])
				return true
			}
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			if len(v) == 0 {
				continue
			}
			v = append(v, d.bt...)
		default:
			v = append(v, d.bt...)
		}
	}
}

func (d *Decoder) ReadAll() string {
	if !d.eof {
		for {
			if len(d.bf) == cap(d.bf) {
				// Add more capacity (let append pick how much).
				d.bf = append(d.bf, 0)[:len(d.bf)]
			}
			n, err := d.rd.Read(d.bf[len(d.bf):cap(d.bf)])
			d.bf = d.bf[:len(d.bf)+n]
			if err != nil {
				break
			}
		}
	}
	return string(d.bf)
}
