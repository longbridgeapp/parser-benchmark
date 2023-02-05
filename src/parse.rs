use std::collections::HashMap;
use std::fmt::Debug;

use nom::branch::{alt, permutation};
use nom::bytes::complete::*;
use nom::character::{complete::*, is_digit};
use nom::combinator::*;
use nom::multi::*;
use nom::sequence::*;
use nom::IResult;

#[derive(Debug, PartialEq, Clone)]
#[allow(non_camel_case_types)]
#[allow(dead_code)]
enum Rule<'a> {
    code(&'a str),
    market(&'a str),
}

fn parse(input: &str) -> Vec<String> {
    let mut codes: HashMap<String, bool> = HashMap::new();

    let mut input = input;
    while let Ok((rest, o)) = stock(input) {
        let mut code = String::from(o.0);
        if !o.1.is_empty() {
            code.push('.');
            code.push_str(o.1);
        }

        codes.insert(code, true);
        input = rest;
    }

    codes.into_keys().collect()
}

fn stock(input: &str) -> IResult<&str, (&str, &str)> {
    let (retain, matched) = alt((
        // $700$
        delimited(
            tag("$"),
            alt((
                terminated(code, char('$')),
                terminated(terminated(code, market), tag("$")),
            )),
            char('$'),
        ),
        // $700
        preceded(char('$'), code),
        // [700]
        delimited(char('['), code, char(']')),
        // (700)
        delimited(char('('), code, char(')')),
    ))(input)?;

    Ok((retain, matched))
}

fn code(input: &str) -> IResult<&str, (&str, &str)> {
    pair(
        alt((
            take_while1(|c: char| c.is_ascii_uppercase()), // us_code
            digit1,                                        // uk_code | a_code
        )),
        alt((market, value("", tag("")))),
    )(input)
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
        assert_eq!(Ok((".foo", ("700", ""))), code("700.foo"));
        assert_eq!(Ok(("", ("1", ""))), code("1"));
        assert_eq!(Ok(("A", ("001029", ""))), code("001029A"));
        assert_eq!(Ok(("", ("FOO", ""))), code("FOO"));
        assert_eq!(Ok(("", ("FOO", "US"))), code("FOO.US"));
        assert!(code("foo").is_err());
    }

    #[test]
    fn test_stock() {
        assert_eq!(vec!["FOO"], parse("$FOO$bar"));
        assert_eq!(vec!["FOO.US"], parse("aaa$FOO.US$bar"));
        assert_eq!(vec!["FOO.US"], parse("$FOO.US$bar"));
    }

    #[test]
    fn test_parse() {
        let raw = include_str!("../tests/example.md");

        assert_eq!(Vec::<String>::new(), parse(raw));
    }
}
