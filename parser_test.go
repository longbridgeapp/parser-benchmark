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
}

func TestExample(t *testing.T) {
	raw, err := ioutil.ReadFile("tests/example.md")
	if err != nil {
		panic(err)
	}

	assert_matches_code(t, "00175.HK, 00175.US, 00231.HK, 00688.HK, 01179.HK, 02269.HK, 100688.SH, 601012.SH, BABA.US, EDBL, FUTU.US, TSLA, bar.US", string(raw))
}
