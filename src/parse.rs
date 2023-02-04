use std::collections::HashMap;

use nom::branch::{alt, permutation};
use nom::bytes::complete::*;
use nom::character::{complete::*, is_digit};
use nom::combinator::*;
use nom::multi::*;
use nom::sequence::*;
use nom::IResult;

enum Stock<'a> {
    Code(&'a str),
    Market(&'a str),
}

fn parse(input: &str) -> Vec<String> {
    let mut codes: HashMap<String, bool> = HashMap::new();

    let mut input = input;
    while let Ok((rest, stock)) = stock(input) {
        let mut code = String::from("");
        let mut market = "";

        match stock {
            Stock::Code(s) => code = s.to_owned(),
            Stock::Market(s) => market = s,
        }

        if code.len() == 0 {
            continue;
        }

        if market.len() > 0 {
            code.push('.');
            code.push_str(market)
        }

        codes.insert(code, true);
        input = rest;
    }

    codes.into_keys().collect()
}

fn stock(input: &str) -> IResult<&str, Stock> {
    let (retain, matched) = delimited(
        // start with $
        tag("$"),
        alt((preceded(code, market), code)),
        // end with $
        tag("$"),
    )(input)?;

    Ok((retain, Stock::Code(matched)))
}

fn code(input: &str) -> IResult<&str, &str> {
    alt((
        take_while1(move |c: char| c.is_ascii_uppercase()), // us_code
        digit1,                                             // uk_code | a_code
    ))(input)
}

fn market(input: &str) -> IResult<&str, &str> {
    preceded(
        tag("."),
        alt((tag("US"), tag("HK"), tag("SG"), tag("SH"), tag("SZ"))),
    )(input)
}

fn sp(i: &str) -> IResult<&str, &str> {
    let chars = " \t\r\n";

    // nom combinators like `take_while` return a function. That function is the
    // parser,to which we can pass the input
    take_while(move |c| chars.contains(c))(i)
}

#[cfg(test)]
mod tests {
    use super::*;

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
        assert_eq!(Ok((".foo", "700")), code("700.foo"));
        assert_eq!(Ok(("", "1")), code("1"));
        assert_eq!(Ok(("A", "001029")), code("001029A"));
        assert_eq!(Ok(("", "FOO")), code("FOO"));
        assert!(code("foo").is_err());
    }

    #[test]
    fn test_market() {
        assert_eq!(Ok(("bar", "FOO")), code("$FOO$bar"));
        assert_eq!(Ok(("bar", "foo.US")), code("aaa$FOO.US$bar"));
        assert_eq!(Ok(("bar", "foo.US")), code("$FOO.US$bar"));
    }

    #[test]
    fn test_parse() {
        let raw = include_str!("../tests/example.md");

        assert_eq!(vec![""], parse(raw));
    }
}
