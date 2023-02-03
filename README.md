# Stock Parser

Use a parser to detect stock symbols in a text file.

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

## Development

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
