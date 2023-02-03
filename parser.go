package stockcode

// Str to returns the string value of the token
func (parser *StockCodeParser) str(node *node32) string {
	return string([]rune(parser.Buffer)[node.begin:node.end])
}

func Parse(input string) (out []string, err error) {
	parser := StockCodeParser{Buffer: input}
	parser.Init()

	codes := map[string]bool{}
	if err := parser.Parse(); err != nil {
		return out, err
	}

	node := parser.AST()
	// node.print(os.Stdout, true, input)

	var cunsumeNode func(node *node32)
	cunsumeNode = func(node *node32) {
		for node != nil {
			code := ""
			market := ""

			if node.pegRule == ruleStock {
				// fmt.Println("ruleStock", node.begin, node.end, node.String(), parser.str(node))
				sub_node := node.up

				for sub_node != nil {
					switch sub_node.pegRule {
					case ruleCode:
						code = parser.str(sub_node)
					case ruleMarket:
						market = parser.str(sub_node)
					case ruleSuffix:
						market = parser.str(sub_node.up)
					}
					sub_node = sub_node.next
				}
			}

			if len(market) != 0 {
				code = code + "." + market
			}

			if len(code) != 0 {
				codes[code] = true
			}

			if node.up != nil {
				cunsumeNode(node.up)
			}
			node = node.next
		}
	}
	cunsumeNode(node)

	for key := range codes {
		out = append(out, key)
	}

	return out, nil
}
