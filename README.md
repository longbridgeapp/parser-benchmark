# Stock Parser

Use a parser to detect stock symbols in a text file.

And this is an example of how to write a parser in Go and Rust (Pest, Nom, rust-peg).

## Usage in Go

```go
package main

import (
  "fmt"
  "github.com/longbridgeapp/stockcode-parser"
)

func main() {
  codes := stockcode.Parse("药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。")
  // ["02269.HK", "601012.SH"]
}
```

### Benchmark

```
BenchmarkParse-8   	           561684	      2112 ns/op	    3056 B/op	      81 allocs/op
BenchmarkParseLongText-8   	   59056	     17410 ns/op	   28384 B/op	     617 allocs/op
```

## Usage in Rust

Add `stockcode-parser` in your `Cargo.toml`

```
[dependencies]
stockcode-parser = { version = "0.1.0", git = "https://github.com/longbridgeapp/stockcode-parser.git" }
```

And then in your `main.rs`

```rs
use stockcode_parser::parse;

let codes = parse("药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。");
// ["02269.HK", "601012.SH"]
```

### Benchmark

```
pest_parse              time:   [1.8119 µs 1.8898 µs 1.9863 µs]
pest_parse_long         time:   [17.221 µs 17.804 µs 18.549 µs]

nom_parse               time:   [428.97 ns 443.55 ns 465.01 ns]
nom_parse_long          time:   [3.8376 µs 4.0039 µs 4.2100 µs]

peg_parse               time:   [975.58 ns 978.74 ns 982.76 ns]
peg_parse_long          time:   [9.6195 µs 9.9444 µs 10.340 µs]
```

## Development

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
