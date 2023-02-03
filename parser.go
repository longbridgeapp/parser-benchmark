package stockcode

func (parser *StockCodeParser) Next() *token32 {
	if parser.pos >= len(parser.tree) {
		return nil
	}

	t := parser.tree[parser.pos]
	parser.pos++
	return &t
}

func (parser *StockCodeParser) Peek() *token32 {
	if parser.peekPos >= len(parser.tree) {
		return nil
	}

	t := parser.tree[parser.peekPos]
	parser.peekPos++
	return &t
}

func (parser *StockCodeParser) Str(token *token32) string {
	return parser.Buffer[token.begin:token.end]
}

func Parse(input string) (out []string, err error) {
	ast := StockCodeParser{Buffer: input}
	ast.Init()

	codes := map[string]bool{}
	if err := ast.Parse(); err != nil {
		return out, err
	}

	for {
		token := ast.Next()
		if token == nil {
			break
		}

		code := ""
		market := ""

		subToken := ast.Next()

		if subToken.pegRule == ruleCode {
			code = ast.Str(subToken)
			subToken = ast.Next()
			if subToken.pegRule == ruleMarket {
				market = ast.Str(subToken)
			}
		}

		if len(market) != 0 {
			code = code + "." + market
		}

		if len(code) != 0 {
			codes[code] = true
		}
	}

	for key := range codes {
		out = append(out, key)
	}

	return out, nil
}
