package stockcode

// Code generated by peg -inline -switch -output grammar.go grammar.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleItem
	ruleLine
	ruleStock
	ruleCode
	ruleUSCode
	ruleHKCode
	ruleACode
	ruleLetter
	ruleNumber
	ruleSuffix
	ruleMarket
	ruleANY
)

var rul3s = [...]string{
	"Unknown",
	"Item",
	"Line",
	"Stock",
	"Code",
	"USCode",
	"HKCode",
	"ACode",
	"Letter",
	"Number",
	"Suffix",
	"Market",
	"ANY",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type StockCodeParser struct {
	pos     int
	peekPos int

	Buffer string
	buffer []rune
	rules  [13]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *StockCodeParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *StockCodeParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *StockCodeParser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *StockCodeParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *StockCodeParser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *StockCodeParser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func Pretty(pretty bool) func(*StockCodeParser) error {
	return func(p *StockCodeParser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*StockCodeParser) error {
	return func(p *StockCodeParser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *StockCodeParser) Init(options ...func(*StockCodeParser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Item <- <(Line* !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					{
						position4 := position
						{
							position5, tokenIndex5 := position, tokenIndex
							{
								position7 := position
								{
									position8, tokenIndex8 := position, tokenIndex
									{
										position10, tokenIndex10 := position, tokenIndex
										if buffer[position] != rune('$') {
											goto l10
										}
										position++
										goto l11
									l10:
										position, tokenIndex = position10, tokenIndex10
									}
								l11:
									if !_rules[ruleCode]() {
										goto l9
									}
									{
										position12, tokenIndex12 := position, tokenIndex
										if !_rules[ruleSuffix]() {
											goto l13
										}
										goto l12
									l13:
										position, tokenIndex = position12, tokenIndex12
										{
											position14, tokenIndex14 := position, tokenIndex
											if !_rules[ruleSuffix]() {
												goto l14
											}
											goto l15
										l14:
											position, tokenIndex = position14, tokenIndex14
										}
									l15:
										if buffer[position] != rune('$') {
											goto l9
										}
										position++
									}
								l12:
									goto l8
								l9:
									position, tokenIndex = position8, tokenIndex8
									{
										switch buffer[position] {
										case '[':
											if buffer[position] != rune('[') {
												goto l6
											}
											position++
											if !_rules[ruleCode]() {
												goto l6
											}
											if buffer[position] != rune(']') {
												goto l6
											}
											position++
										case '(':
											if buffer[position] != rune('(') {
												goto l6
											}
											position++
											if !_rules[ruleCode]() {
												goto l6
											}
											if buffer[position] != rune(')') {
												goto l6
											}
											position++
										default:
											if buffer[position] != rune('$') {
												goto l6
											}
											position++
											if !_rules[ruleCode]() {
												goto l6
											}
										}
									}

								}
							l8:
								add(ruleStock, position7)
							}
							goto l5
						l6:
							position, tokenIndex = position5, tokenIndex5
							{
								position17 := position
								if !matchDot() {
									goto l3
								}
								add(ruleANY, position17)
							}
						}
					l5:
						add(ruleLine, position4)
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position18, tokenIndex18 := position, tokenIndex
					if !matchDot() {
						goto l18
					}
					goto l0
				l18:
					position, tokenIndex = position18, tokenIndex18
				}
				add(ruleItem, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Line <- <(Stock / ANY)> */
		nil,
		/* 2 Stock <- <(('$'? Code (Suffix / (Suffix? '$'))) / ((&('[') ('[' Code ']')) | (&('(') ('(' Code ')')) | (&('$') ('$' Code))))> */
		nil,
		/* 3 Code <- <(USCode / HKCode / ACode)> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				{
					position23, tokenIndex23 := position, tokenIndex
					{
						position25 := position
						{
							position28 := position
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l24
							}
							position++
							add(ruleLetter, position28)
						}
					l26:
						{
							position27, tokenIndex27 := position, tokenIndex
							{
								position29 := position
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l27
								}
								position++
								add(ruleLetter, position29)
							}
							goto l26
						l27:
							position, tokenIndex = position27, tokenIndex27
						}
						add(ruleUSCode, position25)
					}
					goto l23
				l24:
					position, tokenIndex = position23, tokenIndex23
					{
						position31 := position
						if !_rules[ruleNumber]() {
							goto l30
						}
					l32:
						{
							position33, tokenIndex33 := position, tokenIndex
							if !_rules[ruleNumber]() {
								goto l33
							}
							goto l32
						l33:
							position, tokenIndex = position33, tokenIndex33
						}
						add(ruleHKCode, position31)
					}
					goto l23
				l30:
					position, tokenIndex = position23, tokenIndex23
					{
						position34 := position
						if !_rules[ruleNumber]() {
							goto l21
						}
					l35:
						{
							position36, tokenIndex36 := position, tokenIndex
							if !_rules[ruleNumber]() {
								goto l36
							}
							goto l35
						l36:
							position, tokenIndex = position36, tokenIndex36
						}
						add(ruleACode, position34)
					}
				}
			l23:
				add(ruleCode, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 4 USCode <- <Letter+> */
		nil,
		/* 5 HKCode <- <Number+> */
		nil,
		/* 6 ACode <- <Number+> */
		nil,
		/* 7 Letter <- <[A-Z]> */
		nil,
		/* 8 Number <- <[0-9]> */
		func() bool {
			position41, tokenIndex41 := position, tokenIndex
			{
				position42 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l41
				}
				position++
				add(ruleNumber, position42)
			}
			return true
		l41:
			position, tokenIndex = position41, tokenIndex41
			return false
		},
		/* 9 Suffix <- <('.' (Market / ('o' / 'O')))> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				if buffer[position] != rune('.') {
					goto l43
				}
				position++
				{
					position45, tokenIndex45 := position, tokenIndex
					{
						position47 := position
						{
							position48, tokenIndex48 := position, tokenIndex
							if buffer[position] != rune('S') {
								goto l49
							}
							position++
							if buffer[position] != rune('G') {
								goto l49
							}
							position++
							goto l48
						l49:
							position, tokenIndex = position48, tokenIndex48
							if buffer[position] != rune('S') {
								goto l50
							}
							position++
							if buffer[position] != rune('H') {
								goto l50
							}
							position++
							goto l48
						l50:
							position, tokenIndex = position48, tokenIndex48
							{
								switch buffer[position] {
								case 'K':
									if buffer[position] != rune('K') {
										goto l46
									}
									position++
									if buffer[position] != rune('L') {
										goto l46
									}
									position++
								case 'S':
									if buffer[position] != rune('S') {
										goto l46
									}
									position++
									if buffer[position] != rune('Z') {
										goto l46
									}
									position++
								case 'U':
									if buffer[position] != rune('U') {
										goto l46
									}
									position++
									if buffer[position] != rune('S') {
										goto l46
									}
									position++
								default:
									if buffer[position] != rune('H') {
										goto l46
									}
									position++
									if buffer[position] != rune('K') {
										goto l46
									}
									position++
								}
							}

						}
					l48:
						add(ruleMarket, position47)
					}
					goto l45
				l46:
					position, tokenIndex = position45, tokenIndex45
					{
						position52, tokenIndex52 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l53
						}
						position++
						goto l52
					l53:
						position, tokenIndex = position52, tokenIndex52
						if buffer[position] != rune('O') {
							goto l43
						}
						position++
					}
				l52:
				}
			l45:
				add(ruleSuffix, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 10 Market <- <(('S' 'G') / ('S' 'H') / ((&('K') ('K' 'L')) | (&('S') ('S' 'Z')) | (&('U') ('U' 'S')) | (&('H') ('H' 'K'))))> */
		nil,
		/* 11 ANY <- <.> */
		nil,
	}
	p.rules = _rules
	return nil
}
