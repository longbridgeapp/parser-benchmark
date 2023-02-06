use std::collections::HashMap;
use std::fmt::Debug;

use nom::branch::alt;
use nom::bytes::complete::*;
use nom::character::complete::*;
use nom::character::streaming::anychar;
use nom::combinator::*;
use nom::sequence::*;
use nom::IResult;

#[derive(Debug, PartialEq, Clone)]
#[allow(non_camel_case_types)]
#[allow(dead_code)]
enum Rule<'a> {
    stock(&'a str, Option<&'a str>),
    other,
}

pub fn parse(input: &str) -> Vec<String> {
    let mut codes: HashMap<String, bool> = HashMap::new();

    let mut input = input;
    while let Ok((rest, rule)) = line(input) {
        match rule {
            Rule::stock(code, market) => {
                let mut code = String::from(code);
                if market.is_some() {
                    code.push('.');
                    code.push_str(market.unwrap());
                }

                codes.insert(code, true);
            }
            _ => {}
        }

        input = rest;
    }

    codes.into_keys().collect()
}

fn line(input: &str) -> IResult<&str, Rule> {
    alt((stock, other))(input)
}

fn other(input: &str) -> IResult<&str, Rule> {
    let (input, _) = anychar(input)?;
    Ok((input, Rule::other))
}

fn stock(input: &str) -> IResult<&str, Rule> {
    let (retain, matched) = alt((
        // $700$
        delimited(tag("$"), code_with_market, opt(char('$'))),
        // [700]
        delimited(char('['), code_with_market, char(']')),
        // (700)
        delimited(char('('), code_with_market, char(')')),
        // BABA.US
        pair(code, opt(market)),
    ))(input)?;

    Ok((retain, Rule::stock(matched.0, matched.1)))
}

fn code_with_market(input: &str) -> IResult<&str, (&str, Option<&str>)> {
    pair(
        alt((
            take_while1(|c: char| c.is_ascii_uppercase()), // us_code
            digit1,                                        // uk_code | a_code
        )),
        opt(market),
    )(input)
}

fn code(input: &str) -> IResult<&str, &str> {
    alt((
        take_while1(|c: char| c.is_ascii_uppercase()), // us_code
        digit1,                                        // uk_code | a_code
    ))(input)
}

fn market(input: &str) -> IResult<&str, &str> {
    preceded(
        tag("."),
        alt((tag("US"), tag("HK"), tag("SG"), tag("SH"), tag("SZ"))),
    )(input)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[track_caller]
    fn assert_matches_code(codes: &str, input: &str) {
        assert_eq!(codes, parse(input).join(", "))
    }

    #[test]
    fn test_suffix() {
        assert_eq!(Ok(("", "US")), market(".US"));
        assert_eq!(Ok(("", "HK")), market(".HK"));
        assert_eq!(Ok(("", "SG")), market(".SG"));
        assert_eq!(Ok(("", "SH")), market(".SH"));
        assert_eq!(Ok(("", "SZ")), market(".SZ"));
    }

    #[test]
    fn test_code() {
        assert_eq!(Ok((".foo", ("700", None))), code_with_market("700.foo"));
        assert_eq!(Ok(("", ("1", None))), code_with_market("1"));
        assert_eq!(Ok(("A", ("001029", None))), code_with_market("001029A"));
        assert_eq!(Ok(("", ("FOO", None))), code_with_market("FOO"));
        assert_eq!(Ok(("", ("FOO", Some("US")))), code_with_market("FOO.US"));
        assert!(code_with_market("foo").is_err());
    }

    #[test]
    fn test_stock() {
        assert_eq!(Ok(("bar", Rule::stock("FOO", None))), stock("$FOO$bar"));
        assert_eq!(
            Ok(("bar", Rule::stock("FOO", Some("US")))),
            stock("$FOO.US$bar")
        );
        assert_eq!(
            Ok(("bar", Rule::stock("FOO", Some("US")))),
            stock("$FOO.USbar")
        );
        assert_eq!(Ok(("bar", Rule::stock("700", None))), stock("[700]bar"));
        assert_eq!(
            Ok(("bar", Rule::stock("700", Some("HK")))),
            stock("[700.HK]bar")
        );
        assert_eq!(Ok(("bar", Rule::stock("00700", None))), stock("(00700)bar"));
        assert_eq!(
            Ok(("bar", Rule::stock("00700", Some("HK")))),
            stock("(00700.HK)bar")
        );
    }

    #[test]
    fn test_parse_cases() {
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
    fn test_parse() {
        let raw = include_str!("../tests/example.md");

        assert_eq!("00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA", parse(raw).join(", "));
    }
}
