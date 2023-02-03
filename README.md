# Stock Parser

Use a parser to detect stock symbols in a text file.

## Usage

```
use stockcode_parser::parse;

let codes = parse("药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。");
// ["02269.HK", "601012.SH"]
```
