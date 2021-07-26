# xmlreader
Golang small xml reader


It can parse node name, node value and attribute


```

    r := strings.NewReader(`<animal id="21"/><animal id="23">armadillo</animal>`)
    d := xmlreader.New(r)

    fmt.Println(d.Read())
    fmt.Println(d.Name)
    fmt.Println(d.HasValue())
    fmt.Println(d.Value)

```


For more feature, you can use standart lib encoding/xml.

```

decoder := xml.NewDecoder(xmlFile)

t, _ := decoder.Token()

```
