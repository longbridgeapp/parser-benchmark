package parser_benchmark

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
	ruleOTHER
	ruleStock
	ruleCode
	ruleUSCode
	ruleHKCode
	ruleACode
	ruleLetter
	ruleNumber
	ruleSuffix
	ruleMarket
	ruleSP
)

var rul3s = [...]string{
	"Unknown",
	"Item",
	"Line",
	"OTHER",
	"Stock",
	"Code",
	"USCode",
	"HKCode",
	"ACode",
	"Letter",
	"Number",
	"Suffix",
	"Market",
	"SP",
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
	rules  [14]func() bool
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
									if buffer[position] != rune('(') {
										goto l16
									}
									position++
									if !_rules[ruleCode]() {
										goto l16
									}
									if buffer[position] != rune(')') {
										goto l16
									}
									position++
									goto l8
								l16:
									position, tokenIndex = position8, tokenIndex8
									{
										switch buffer[position] {
										case '(':
											if buffer[position] != rune('(') {
												goto l6
											}
											position++
											{
												position18, tokenIndex18 := position, tokenIndex
												if buffer[position] != rune('N') {
													goto l19
												}
												position++
												if buffer[position] != rune('Y') {
													goto l19
												}
												position++
												if buffer[position] != rune('S') {
													goto l19
												}
												position++
												if buffer[position] != rune('E') {
													goto l19
												}
												position++
												goto l18
											l19:
												position, tokenIndex = position18, tokenIndex18
												if buffer[position] != rune('N') {
													goto l6
												}
												position++
												if buffer[position] != rune('A') {
													goto l6
												}
												position++
												if buffer[position] != rune('S') {
													goto l6
												}
												position++
												if buffer[position] != rune('D') {
													goto l6
												}
												position++
												if buffer[position] != rune('A') {
													goto l6
												}
												position++
												if buffer[position] != rune('Q') {
													goto l6
												}
												position++
											}
										l18:
											{
												position20, tokenIndex20 := position, tokenIndex
												if buffer[position] != rune('：') {
													goto l21
												}
												position++
												goto l20
											l21:
												position, tokenIndex = position20, tokenIndex20
												if buffer[position] != rune(':') {
													goto l6
												}
												position++
											}
										l20:
										l22:
											{
												position23, tokenIndex23 := position, tokenIndex
												{
													position24 := position
													{
														position25, tokenIndex25 := position, tokenIndex
														if buffer[position] != rune(' ') {
															goto l26
														}
														position++
														goto l25
													l26:
														position, tokenIndex = position25, tokenIndex25
														if buffer[position] != rune('\t') {
															goto l23
														}
														position++
													}
												l25:
													add(ruleSP, position24)
												}
												goto l22
											l23:
												position, tokenIndex = position23, tokenIndex23
											}
											if !_rules[ruleCode]() {
												goto l6
											}
											if buffer[position] != rune(')') {
												goto l6
											}
											position++
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
								position27 := position
								if !matchDot() {
									goto l3
								}
								add(ruleOTHER, position27)
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
					position28, tokenIndex28 := position, tokenIndex
					if !matchDot() {
						goto l28
					}
					goto l0
				l28:
					position, tokenIndex = position28, tokenIndex28
				}
				add(ruleItem, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Line <- <(Stock / OTHER)> */
		nil,
		/* 2 OTHER <- <.> */
		nil,
		/* 3 Stock <- <(('$'? Code (Suffix / (Suffix? '$'))) / ('(' Code ')') / ((&('(') ('(' (('N' 'Y' 'S' 'E') / ('N' 'A' 'S' 'D' 'A' 'Q')) ('：' / ':') SP* Code ')')) | (&('[') ('[' Code ']')) | (&('$') ('$' Code))))> */
		nil,
		/* 4 Code <- <(USCode / HKCode / ACode)> */
		func() bool {
			position32, tokenIndex32 := position, tokenIndex
			{
				position33 := position
				{
					position34, tokenIndex34 := position, tokenIndex
					{
						position36 := position
						{
							position39 := position
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l35
							}
							position++
							add(ruleLetter, position39)
						}
					l37:
						{
							position38, tokenIndex38 := position, tokenIndex
							{
								position40 := position
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l38
								}
								position++
								add(ruleLetter, position40)
							}
							goto l37
						l38:
							position, tokenIndex = position38, tokenIndex38
						}
						add(ruleUSCode, position36)
					}
					goto l34
				l35:
					position, tokenIndex = position34, tokenIndex34
					{
						position42 := position
						if !_rules[ruleNumber]() {
							goto l41
						}
					l43:
						{
							position44, tokenIndex44 := position, tokenIndex
							if !_rules[ruleNumber]() {
								goto l44
							}
							goto l43
						l44:
							position, tokenIndex = position44, tokenIndex44
						}
						add(ruleHKCode, position42)
					}
					goto l34
				l41:
					position, tokenIndex = position34, tokenIndex34
					{
						position45 := position
						if !_rules[ruleNumber]() {
							goto l32
						}
					l46:
						{
							position47, tokenIndex47 := position, tokenIndex
							if !_rules[ruleNumber]() {
								goto l47
							}
							goto l46
						l47:
							position, tokenIndex = position47, tokenIndex47
						}
						add(ruleACode, position45)
					}
				}
			l34:
				add(ruleCode, position33)
			}
			return true
		l32:
			position, tokenIndex = position32, tokenIndex32
			return false
		},
		/* 5 USCode <- <Letter+> */
		nil,
		/* 6 HKCode <- <Number+> */
		nil,
		/* 7 ACode <- <Number+> */
		nil,
		/* 8 Letter <- <[A-Z]> */
		nil,
		/* 9 Number <- <[0-9]> */
		func() bool {
			position52, tokenIndex52 := position, tokenIndex
			{
				position53 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l52
				}
				position++
				add(ruleNumber, position53)
			}
			return true
		l52:
			position, tokenIndex = position52, tokenIndex52
			return false
		},
		/* 10 Suffix <- <('.' (Market / ('o' / 'O')))> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				if buffer[position] != rune('.') {
					goto l54
				}
				position++
				{
					position56, tokenIndex56 := position, tokenIndex
					{
						position58 := position
						{
							position59, tokenIndex59 := position, tokenIndex
							if buffer[position] != rune('S') {
								goto l60
							}
							position++
							if buffer[position] != rune('G') {
								goto l60
							}
							position++
							goto l59
						l60:
							position, tokenIndex = position59, tokenIndex59
							if buffer[position] != rune('S') {
								goto l61
							}
							position++
							if buffer[position] != rune('H') {
								goto l61
							}
							position++
							goto l59
						l61:
							position, tokenIndex = position59, tokenIndex59
							{
								switch buffer[position] {
								case 'K':
									if buffer[position] != rune('K') {
										goto l57
									}
									position++
									if buffer[position] != rune('L') {
										goto l57
									}
									position++
								case 'S':
									if buffer[position] != rune('S') {
										goto l57
									}
									position++
									if buffer[position] != rune('Z') {
										goto l57
									}
									position++
								case 'U':
									if buffer[position] != rune('U') {
										goto l57
									}
									position++
									if buffer[position] != rune('S') {
										goto l57
									}
									position++
								default:
									if buffer[position] != rune('H') {
										goto l57
									}
									position++
									if buffer[position] != rune('K') {
										goto l57
									}
									position++
								}
							}

						}
					l59:
						add(ruleMarket, position58)
					}
					goto l56
				l57:
					position, tokenIndex = position56, tokenIndex56
					{
						position63, tokenIndex63 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l64
						}
						position++
						goto l63
					l64:
						position, tokenIndex = position63, tokenIndex63
						if buffer[position] != rune('O') {
							goto l54
						}
						position++
					}
				l63:
				}
			l56:
				add(ruleSuffix, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 11 Market <- <(('S' 'G') / ('S' 'H') / ((&('K') ('K' 'L')) | (&('S') ('S' 'Z')) | (&('U') ('U' 'S')) | (&('H') ('H' 'K'))))> */
		nil,
		/* 12 SP <- <(' ' / '\t')> */
		nil,
	}
	p.rules = _rules
	return nil
}
