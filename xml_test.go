package xmlreader

import (
	"fmt"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	r := strings.NewReader(`<multiRef id="id369" soapenc:root="0" soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xsi:type="ns14:ItemType_Generic" xmlns:ns14="urn:issuing_v_01_02_xsd" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/">
	test Value
   </multiRef>`)
	d := New(r)

	for d.Read() {
		fmt.Println(d.Name)
		fmt.Println(d.GetAttribute("id"))
		fmt.Println(d.HasValue())
		fmt.Println(d.Value)
	}

	fmt.Println(d.ReadAll())
}

func BenchmarkFileRead(b *testing.B) {
	r := strings.NewReader(`<multiRef id="id369" soapenc:root="0" soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xsi:type="ns14:ItemType_Generic" xmlns:ns14="urn:issuing_v_01_02_xsd" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/">
	test Value
   </multiRef>`)
	d := New(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Read()
	}
}
