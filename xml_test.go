package xmlreader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	r := strings.NewReader(`<animal id="21"/><animal id="23">armadillo</animal>`)
	d := New(r)

	fmt.Println(d.Read())
	fmt.Println(d.Name)
	fmt.Println(d.HasValue())
	fmt.Println(d.Value)
}

func BenchmarkFileRead(b *testing.B) {
	f, err := os.Open("file.txt")
	if err != nil {
		b.Error(err)
	}
	defer f.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 1)
		for {
			_, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

		}
	}
}

func BenchmarkBufIO(b *testing.B) {
	f, err := os.Open("file.txt")
	if err != nil {
		b.Error(err)
	}
	defer f.Close()
	bio := bufio.NewReader(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		for {
			_, err := bio.ReadByte()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}
		}
	}
}

/*
BenchmarkBufIO-8   	     1544479	       778.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkFileRead-8   	 1371295	       834.9 ns/op	       0 B/op	       0 allocs/op
*/
