# xmlreader
Golang xml reader


```

    r := strings.NewReader(`<animal id="21"/><animal id="23">armadillo</animal>`)
	d := New(r)

	fmt.Println(d.Read())
	fmt.Println(d.Name)
	fmt.Println(d.HasValue())
	fmt.Println(d.Value)

```