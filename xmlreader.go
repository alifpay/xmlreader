package xmlreader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotFound = errors.New("not found")
)

type Decoder struct {
	br    *bufio.Reader
	Name  string
	Value string
	eof   bool
	end   bool
}

func New(r io.Reader) *Decoder {
	return &Decoder{br: bufio.NewReader(r)}
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
		b, err := d.br.ReadByte()
		if err == io.EOF {
			d.eof = true
			return false
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch b {
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
				nm = append(nm, b)
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
		b, err := d.br.ReadByte()
		if err == io.EOF {
			d.eof = true
			return "", err
		}
		if err != nil {
			return "", err
		}

		switch b {
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
		case '?', '/':
			return "", ErrNotFound
		case '>':
			d.end = true
			return "", ErrNotFound
		default:
			buf = append(buf, b)
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
	buf := make([]byte, 0)
	if d.end {
		has = true
	}
	for {
		b, err := d.br.ReadByte()
		if err == io.EOF {
			d.eof = true
			return false
		}
		if err != nil {
			return false
		}

		switch b {
		case '>':
			has = true
			buf = make([]byte, 0)
			continue
		case '<':
			if has && len(buf) > 0 {
				d.Value = string(buf)
				return true
			}
		case '?', '/':
			return false
		default:
			buf = append(buf, b)
		}
	}
}
