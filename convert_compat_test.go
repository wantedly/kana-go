package kana_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wantedly/kana-go"
)

func TestCompatConvert(t *testing.T) {
	var testcases = []struct {
		name    string
		input   string
		options kana.ConvertOptions
		expect  string
	}{
		{
			name:    "With FullwidthToNarrow, Without CompatQuotes",
			input:   "´‘’“”＂＇",
			options: kana.FullwidthToNarrow,
			expect:  "´‘’“”\"'",
		},
		{
			name:    "With FullwidthToNarrow, With CompatQuotes",
			input:   "´‘’“”＂＇",
			options: kana.FullwidthToNarrow | kana.CompatQuotes,
			expect:  "'`'\"\"＂＇",
		},
		{
			name:    "Without CompatMinus",
			input:   "—―−－",
			options: 0,
			expect:  "—―−－",
		},
		{
			name:    "With CompatMinus",
			input:   "—―−－",
			options: kana.CompatMinus,
			expect:  "——−−",
		},
		{
			name:    "With FullwidthToNarrow, Without CompatMinus",
			input:   "—―−－",
			options: kana.FullwidthToNarrow,
			expect:  "—―−-",
		},
		{
			name:    "With FullwidthToNarrow, With CompatMinus",
			input:   "—―−－",
			options: kana.FullwidthToNarrow | kana.CompatMinus,
			expect:  "----",
		},
		{
			name:    "Without CompatOverline",
			input:   "￣～",
			options: 0,
			expect:  "￣～",
		},
		{
			name:    "With CompatOverline",
			input:   "￣～",
			options: kana.CompatOverline,
			expect:  "‾～",
		},
		{
			name:    "With FullwidthToNarrow, Without CompatOverline",
			input:   "￣～",
			options: kana.FullwidthToNarrow,
			expect:  "¯~",
		},
		{
			name:    "With FullwidthToNarrow, With CompatOverline",
			input:   "￣～",
			options: kana.FullwidthToNarrow | kana.CompatOverline,
			expect:  "‾～",
		},
		{
			name:    "Without CompatCurrency",
			input:   "￠￡￥￦",
			options: 0,
			expect:  "￠￡￥￦",
		},
		{
			name:    "With CompatCurrency",
			input:   "￠￡￥￦",
			options: kana.CompatCurrency,
			expect:  "¢£¥￦",
		},
		{
			name:    "Without CompatOtherSymbols",
			input:   "∥￢￤",
			options: 0,
			expect:  "∥￢￤",
		},
		{
			name:    "With CompatOtherSymbols",
			input:   "∥￢￤",
			options: kana.CompatOtherSymbols,
			expect:  "‖¬¦",
		},
		{
			name:    "With FullwidthToNarrow, Without CompatCurrency",
			input:   "￠￡￥￦",
			options: kana.FullwidthToNarrow,
			expect:  "¢£¥₩",
		},
		{
			name:    "With FullwidthToNarrow, With CompatCurrency",
			input:   "￠￡￥￦",
			options: kana.FullwidthToNarrow | kana.CompatCurrency,
			expect:  "¢£¥￦",
		},
		{
			name:    "With FullwidthToNarrow, Without CompatBracket",
			input:   "〈〉｟｠",
			options: kana.FullwidthToNarrow,
			expect:  "〈〉⦅⦆",
		},
		{
			name:    "With FullwidthToNarrow, With CompatBracket",
			input:   "〈〉｟｠",
			options: kana.FullwidthToNarrow | kana.CompatBrackets,
			expect:  "<>｟｠",
		},
		{
			name:    "With HalfwidthToWide, Without CompatVoicedSoundMarks",
			input:   "ﾞﾟ",
			options: kana.HalfwidthToWide,
			expect:  "\u3099\u309A",
		},
		{
			name:    "With HalfwidthToWide, With CompatVoicedSoundMarks",
			input:   "ﾞﾟ",
			options: kana.HalfwidthToWide | kana.CompatVoicedSoundMarks,
			expect:  "゛゜",
		},
		{
			name:    "With HalfwidthToWide, Without CompatKeepHalfwidthHangul",
			input:   "\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ",
			options: kana.HalfwidthToWide,
			expect:  "\u3164ㄱㄲㄳㄴㄵㄶㄷㄸㄹㄺㄻㄼㄽㄾㄿㅀㅁㅂㅃㅄㅅㅆㅇㅈㅉㅊㅋㅌㅍㅎㅏㅐㅑㅒㅓㅔㅕㅖㅗㅘㅙㅚㅛㅜㅝㅞㅟㅠㅡㅢㅣ",
		},
		{
			name:    "With HalfwidthToWide, With CompatKeepHalfwidthHangul",
			input:   "\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ",
			options: kana.HalfwidthToWide | kana.CompatKeepHalfwidthHangul,
			expect:  "\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ",
		},
		{
			name:    "With HalfwidthToWide, Without CompatKeepHalfwidthSymbols",
			input:   "￨￩￪￫￬￭￮",
			options: kana.HalfwidthToWide,
			expect:  "│←↑→↓■○",
		},
		{
			name:    "With HalfwidthToWide, With CompatKeepHalfwidthSymbols",
			input:   "￨￩￪￫￬￭￮",
			options: kana.HalfwidthToWide | kana.CompatKeepHalfwidthSymbols,
			expect:  "￨￩￪￫￬￭￮",
		},
		{
			name:    "With KatakanaToHiragana, Without CompatKanaRestriction",
			input:   "ヵヶヷヸヹヺ𛅕𛅤𛅥𛅦",
			options: kana.KatakanaToHiragana,
			expect:  "ゕゖわ\u3099ゐ\u3099ゑ\u3099を\u3099𛄲𛅐𛅑𛅒",
		},
		{
			name:    "With KatakanaToHiragana, With CompatKanaRestriction",
			input:   "ヵヶヷヸヹヺ𛅕𛅤𛅥𛅦",
			options: kana.KatakanaToHiragana | kana.CompatKanaRestriction,
			expect:  "ヵヶヷヸヹヺ𛅕𛅤𛅥𛅦",
		},
		{
			name:    "With HiraganaToKatakana, Without CompatKanaRestriction",
			input:   "ゕゖ𛄲𛅐𛅑𛅒",
			options: kana.HiraganaToKatakana,
			expect:  "ヵヶ𛅕𛅤𛅥𛅦",
		},
		{
			name:    "With HiraganaToKatakana, With CompatKanaRestriction",
			input:   "ゕゖ𛄲𛅐𛅑𛅒",
			options: kana.HiraganaToKatakana | kana.CompatKanaRestriction,
			expect:  "ゕゖ𛄲𛅐𛅑𛅒",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := kana.Convert(tc.input, tc.options)
			if diff := cmp.Diff(actual, tc.expect); diff != "" {
				t.Errorf("diff (-actual +expect): %s", diff)
			}
		})
	}
}
