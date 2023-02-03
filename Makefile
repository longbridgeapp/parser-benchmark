.SUFFIXES: .peg .go

.peg.go:
	peg -inline -switch -output $@ $<

all: grammar.go
test: all
	go test