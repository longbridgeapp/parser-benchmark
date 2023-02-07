use std::collections::HashMap;

#[derive(Debug, PartialEq, Clone)]
#[allow(non_camel_case_types)]
#[allow(dead_code)]
pub enum Rule<'a> {
    stock(&'a str, Option<&'a str>),
    other,
}

peg::parser! {
  grammar stockcode() for str {
    pub rule item() -> Vec<Rule<'input>> = stock()*

    rule _() = [' ' | '\t' | '\r' | '\n']*

    rule other() = quiet!{ [_] }

    rule stock() -> Rule<'input> =
        "$" c:code() "." m:suffix() "$"? { Rule::stock(c, Some(m))}
        / "$" c:code() "$"? { Rule::stock(c, None) }
        / c:code() "$" { Rule::stock(c, None) }
        / c:code() "." m:suffix() { Rule::stock(c, Some(m)) }
        / open_bracket() c:code() close_bracket() { Rule::stock(c, None) }
        / other() { Rule::other }

    rule open_bracket() = ['(' | '（' | '[' ]
    rule close_bracket() = [')' | '）' | ']' ]

    rule code() -> &'input str = us_code() / hk_code() / a_code()

    rule us_code() -> &'input str  = $(upper_alpha()+ / number()+)
    rule hk_code() -> &'input str = $(number()*<6,6> / "0"+ number()+)
    rule a_code() -> &'input str = $(number()*<6,6>)
    rule suffix() -> &'input str = s:$(market()) { s } / "O" { "" }

    rule upper_alpha() -> &'input str = $(['A'..='Z'])
    rule number() -> &'input str = $(['0'..='9'])
    rule market() -> &'input str  = $("US" / "HK" / "SG" / "SH" / "SZ")
  }
}

/// Returns matched stock code as `Vec<String>`
pub fn parse(input: &str) -> Vec<String> {
    let mut codes: HashMap<String, bool> = HashMap::new();

    let pairs: Vec<Rule> = stockcode::item(input).unwrap();
    for pair in pairs {
        match pair {
            Rule::stock(code, market) => {
                let mut code = String::from(code);
                if market.is_some() {
                    let market = market.unwrap();
                    if !market.is_empty() {
                        code.push('.');
                        code.push_str(market);
                    }
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
