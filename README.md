# Parser benchmark

This is [Ticker symbol](https://en.wikipedia.org/wiki/Ticker_symbol) parse example to show how to write a parser in Go and Rust (Pest, Nom, rust-peg).
And run benchmarks to compare the performance.

## Parsers

- [Pest](https://pest.rs)
- [Nom](https://github.com/rust-bakery/nom)
- [rust-peg](https://github.com/kevinmehall/rust-peg)
- [tdewolff/parse](github.com/tdewolff/parse) - (Go)

### Benchmark in Go

```
Benchmark_tdewolff_parse-8   	    561684	      2112 ns/op	    3056 B/op	      81 allocs/op
Benchmark_tdewolff_parse_long-8    59056	     17410 ns/op	   28384 B/op	     617 allocs/op
```

### Benchmark in Rust

```
pest_parse              time:   [1.8119 µs 1.8898 µs 1.9863 µs]
pest_parse_long         time:   [17.221 µs 17.804 µs 18.549 µs]

nom_parse               time:   [428.97 ns 443.55 ns 465.01 ns]
nom_parse_long          time:   [3.8376 µs 4.0039 µs 4.2100 µs]

peg_parse               time:   [975.58 ns 978.74 ns 982.76 ns]
peg_parse_long          time:   [9.6195 µs 9.9444 µs 10.340 µs]
```

## Development in Go

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
