use std::{collections::HashMap, fmt::Debug};

use pest_derive::Parser;

use pest::Parser;

#[derive(Parser)]
#[grammar = "grammar.pest"]
struct StockCodeParser;

/// Returns matched stock code as `Vec<String>`
pub fn parse(input: &str) -> Vec<String> {
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

    #[track_caller]
    fn assert_matches_code(codes: &str, input: &str) {
        assert_eq!(codes, parse(input).join(", "))
    }

    #[test]
    fn test_parse() {
        assert_matches_code("BABA.US", "Alibaba BABA.US published its Q2 results");
        assert_matches_code("BABA", "Alibaba $BABA published its Q2 results");
        assert_matches_code("BABA", "阿里巴巴$BABA发布财报");
        assert_matches_code("BABA.US", "阿里巴巴$BABA.US发布财报");
        assert_matches_code("BABA.US", "阿里巴巴$BABA.US$发布财报");
        assert_matches_code("BABA.US", "阿里巴巴BABA.US$发布财报");
        assert_matches_code("BABA.US", "阿里巴巴BABA.US发布财报");
        assert_matches_code("BABA", "阿里巴巴BABA$发布财报");
        assert_matches_code("", "腾讯700发布财报");
        assert_matches_code("700", "腾讯[700]发布财报");
        assert_matches_code("700", "腾讯(700)发布财报");
        assert_matches_code("00700.HK", "腾讯00700.HK发布财报");
        assert_matches_code("700", "腾讯（700）发布财报");
    }

    #[test]
    fn test_example() {
        let raw = include_str!("../tests/example.md");

        assert_matches_code("00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA", raw);
    }

    #[test]
    fn test_routers_format() {
        assert_matches_code("EDBL", "公司（EDBL.O）宣布");
        assert_matches_code("EDBL, SA", "公司（EDBL.O,SA.O）宣布");
        assert_matches_code("EDBL, SA", "EDBL.O,SA.O）宣布");
    }
}
