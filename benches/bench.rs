use criterion::{criterion_group, criterion_main, Criterion};

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("pest_parse", |b| {
        b.iter(|| stockcode_parser::pest_parser::parse("阿里巴巴 $BABA.US 发布财报"))
    });

    c.bench_function("pest_parse_long", |b| {
        b.iter(|| stockcode_parser::pest_parser::parse("海外发展 (00688.HK,100688.SH, 100681) 截至 10:47 下跌 3.13%，大和将华住 (01179.HK)、药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。路透社格式：可食用花园股份公司（EDBL.O）宣布以 1,020 万没有公开募股。"))
    });

    c.bench_function("nom_parse", |b| {
        b.iter(|| stockcode_parser::nom_parser::parse("阿里巴巴 $BABA.US 发布财报"))
    });

    c.bench_function("nom_parse_long", |b| {
        b.iter(|| stockcode_parser::nom_parser::parse("海外发展 (00688.HK,100688.SH, 100681) 截至 10:47 下跌 3.13%，大和将华住 (01179.HK)、药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。路透社格式：可食用花园股份公司（EDBL.O）宣布以 1,020 万没有公开募股。"))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
