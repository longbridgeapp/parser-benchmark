use std::{collections::HashMap, fmt::Debug};

use pest_derive::Parser;

use pest::Parser;

#[derive(Parser)]
#[grammar = "grammar.pest"]
struct StockCodeParser;

fn parse(input: &str) -> Vec<String> {
    let mut codes: HashMap<String, bool> = HashMap::new();

    let pairs = StockCodeParser::parse(Rule::item, input).unwrap();
    let pairs = pairs.flatten();

    for pair in pairs {
        match pair.as_rule() {
            Rule::stock => {
                let mut code = String::new();
                let mut market = String::new();

                for pair in pair.into_inner() {
                    match pair.as_rule() {
                        Rule::code => code = pair.as_str().to_string(),
                        Rule::market => market = pair.as_str().to_string(),
                        _ => continue,
                    }
                }

                if !market.is_empty() {
                    code.push('.');
                    code.push_str(&market);
                }

                codes.insert(code, true);
            }
            _ => {}
        }
    }

    let mut out: Vec<String> = codes.into_keys().collect();
    out.sort_by(|a, b| a.cmp(&b));
    out
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let raw = include_str!("../tests/example.md");

        let codes = parse(raw);
        assert_eq!(
            "00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA",
            codes.join(", ")
        );
    }
}
