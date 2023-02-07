package parser_benchmark_test

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/longbridgeapp/assert"
	parser_benchmark "github.com/longbridgeapp/parser-benchmark"
)

func assert_matches_code(t *testing.T, expected string, input string) {
	t.Helper()

	codes, err := parser_benchmark.Parse(input)
	assert.NoError(t, err)
	sort.Strings(codes)
	assert.Equal(t, expected, strings.Join(codes, ", "))
}

func TestParse(t *testing.T) {
	assert_matches_code(t, "BABA.US", "Alibaba BABA.US published its Q2 results")
	assert_matches_code(t, "BABA", "Alibaba $BABA published its Q2 results")
	assert_matches_code(t, "BABA", "阿里巴巴$BABA发布财报")
	assert_matches_code(t, "BABA.US", "阿里巴巴$BABA.US发布财报")
	assert_matches_code(t, "BABA.US", "阿里巴巴$BABA.US$发布财报")
	assert_matches_code(t, "BABA.US", "阿里巴巴BABA.US$发布财报")
	assert_matches_code(t, "BABA.US", "阿里巴巴BABA.US发布财报")
	assert_matches_code(t, "BABA", "阿里巴巴BABA$发布财报")
	assert_matches_code(t, "", "腾讯700发布财报")
	assert_matches_code(t, "700", "腾讯[700]发布财报")
	assert_matches_code(t, "700", "腾讯(700)发布财报")
	assert_matches_code(t, "00700.HK", "腾讯00700.HK发布财报")
	assert_matches_code(t, "TSLA", "Tesla Inc (TSLA.O) will finalise a deal to invest in a production facility in his country")
	assert_matches_code(t, "TSLA", "Only the fortune of Tesla's (TSLA)")
	assert_matches_code(t, "A, B, DE, FU", "安捷伦科技 (NYSE：A) (NYSE： B), (NYSE:FU) (NYSE: DE)")
	assert_matches_code(t, "A, B, DE, FU", "安捷伦科技 (NASDAQ：A) (NASDAQ： B), (NASDAQ:FU) (NASDAQ: DE)")
}

func TestSpecialMarket(t *testing.T) {
	assert_matches_code(t, "AB.US, Q.US", "美股简短的股票代码 Q.US, AB.US")
	assert_matches_code(t, "NICH.KL", "Malaysia Market (NICH.KL)")
}

func TestExample(t *testing.T) {
	raw, err := ioutil.ReadFile("tests/example.md")
	if err != nil {
		panic(err)
	}

	assert_matches_code(t, "00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA", string(raw))
}

func Benchmark_tdewolff_parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser_benchmark.Parse("阿里巴巴 $BABA.US 发布财报")
	}
}

func Benchmark_tdewolff_parse_long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser_benchmark.Parse("海外发展 (00688.HK,100688.SH, 100681) 截至 10:47 下跌 3.13%，大和将华住 (01179.HK)、药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。路透社格式：可食用花园股份公司（EDBL.O）宣布以 1,020 万没有公开募股。")
	}
}
