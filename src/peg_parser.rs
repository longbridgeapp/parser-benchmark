// peg::parser! {
//   grammar list_parser() for str {
//     rule market()
//       = ("US" / "HK")

//     pub rule list() -> Vec<u32>
//       = "[" l:(number() ** ",") "]" { l }
//   }
// }

// /// Returns matched stock code as `Vec<String>`
// pub fn parse(input: &str) -> Vec<String> {
//     vec![]
// }
