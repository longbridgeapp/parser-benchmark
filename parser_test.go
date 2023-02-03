package stockcode_test

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/longbridgeapp/assert"
	stockcodeparser "github.com/longbridgeapp/stockcode-parser"
)

func assert_matches_code(t *testing.T, expected string, input string) {
	t.Helper()

	codes, err := stockcodeparser.Parse(input)
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
}

func TestExample(t *testing.T) {
	raw, err := ioutil.ReadFile("tests/example.md")
	if err != nil {
		panic(err)
	}

	assert_matches_code(t, "00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA", string(raw))
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stockcodeparser.Parse("阿里巴巴$BABA.US发布财报")
	}
}

func BenchmarkParseLongText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stockcodeparser.Parse("海外发展 (00688.HK,100688.SH, 100681) 截至 10:47 下跌 3.13%，大和将华住 (01179.HK)、药明生物 (02269.HK) 及隆基绿能科技 (601012.SH) 纳入中港市场首选名单。路透社格式：可食用花园股份公司（EDBL.O）宣布以 1,020 万没有公开募股。")
	}
}
