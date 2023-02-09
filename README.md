# Parser benchmark

This is [Ticker symbol](https://en.wikipedia.org/wiki/Ticker_symbol) parse example to show how to write a parser in Go and Rust (Pest, Nom, rust-peg).
And run benchmarks to compare the performance.

## Parsers

- [Pest](https://pest.rs)
- [Nom](https://github.com/rust-bakery/nom)
- [rust-peg](https://github.com/kevinmehall/rust-peg)
- [pointlander/peg](https://github.com/pointlander/peg) - (Go)

### Benchmark in Go

```
Benchmark_pointlander_peg-8       489860	      2729 ns/op
Benchmark_pointlander_peg-8    	   64664	     21624 ns/op
Benchmark_pointlander_peg-8   	    1910	    589914 ns/op
```

### Benchmark in Rust

```
pest_parse              time:   [1.6484 µs 1.6730 µs 1.7062 µs]
pest_parse_long         time:   [16.631 µs 16.912 µs 17.265 µs]
pest_parse_large        time:   [571.65 µs 590.45 µs 616.38 µs]

nom_parse               time:   [420.42 ns 423.97 ns 428.22 ns]
nom_parse_long          time:   [3.7042 µs 3.7315 µs 3.7676 µs]
nom_parse_large         time:   [107.76 µs 109.42 µs 112.08 µs]

peg_parse               time:   [1.0119 µs 1.0189 µs 1.0268 µs]
peg_parse_long          time:   [9.7554 µs 9.9865 µs 10.388 µs]
peg_parse_large         time:   [331.86 µs 342.67 µs 358.57 µs]
```

## Development in Go

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
