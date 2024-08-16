package nkf_test

import (
	"testing"

	"github.com/wantedly/kana-go"
	"github.com/wantedly/kana-go/nkf"
)

func TestParseOptions(t *testing.T) {
	compatBase := kana.CompatMinus | kana.CompatOverline | kana.CompatCurrency | kana.CompatOtherSymbols
	testcases := []struct {
		name      string
		text      string
		expect    kana.ConvertOptions
		expectErr string
	}{
		{
			name:      "Without -w",
			text:      "",
			expectErr: "-w is required",
		},
		{
			name:      "Without -W",
			text:      "-w",
			expectErr: "-W is required",
		},
		{
			name:      "Without -m0",
			text:      "-w -W",
			expectErr: "-m0 is required",
		},
		{
			name:   "Minimum options",
			text:   "-w -W -m0",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols,
		},
		{
			name:   "Minimum alt options",
			text:   "-w8 -W8 -m0",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols,
		},
		{
			name:   "Minimum longhand options",
			text:   "--utf8 --utf8-input -m0",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols,
		},
		{
			name:   "-h",
			text:   "-w -W -m0 -h",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.KatakanaToHiragana | kana.CompatKanaRestriction,
		},
		{
			name:   "-h1",
			text:   "-w -W -m0 -h1",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.KatakanaToHiragana | kana.CompatKanaRestriction,
		},
		{
			name:   "--hiragana",
			text:   "-w -W -m0 --hiragana",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.KatakanaToHiragana | kana.CompatKanaRestriction,
		},
		{
			name:   "-h2",
			text:   "-w -W -m0 -h2",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.HiraganaToKatakana | kana.CompatKanaRestriction,
		},
		{
			name:   "--katakana",
			text:   "-w -W -m0 --katakana",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.HiraganaToKatakana | kana.CompatKanaRestriction,
		},
		{
			name:   "-h3",
			text:   "-w -W -m0 -h3",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.KatakanaToHiragana | kana.HiraganaToKatakana | kana.CompatKanaRestriction,
		},
		{
			name:   "--katakana-hiragana",
			text:   "-w -W -m0 --katakana-hiragana",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.KatakanaToHiragana | kana.HiraganaToKatakana | kana.CompatKanaRestriction,
		},
		{
			name:   "-Z",
			text:   "-w -W -m0 -Z",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets | kana.CompatKeepSpaces,
		},
		{
			name:   "-Z0",
			text:   "-w -W -m0 -Z0",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets | kana.CompatKeepSpaces,
		},
		{
			name:   "-Z1",
			text:   "-w -W -m0 -Z1",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets,
		},
		{
			name:   "-Z2",
			text:   "-w -W -m0 -Z2",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets | kana.CompatDoubleSpaces,
		},
		{
			name:   "-Z4",
			text:   "-w -W -m0 -Z4",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols | kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets | kana.CompatKeepSpaces | kana.CompatWideKatakanaToHalfwidth,
		},
		{
			name:   "-x",
			text:   "-w -W -m0 -x",
			expect: compatBase,
		},
		{
			name:   "-X",
			text:   "-w -W -m0 -X",
			expect: compatBase | kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatKeepHalfwidthHangul | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthSymbols,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := nkf.ParseOptions(tc.text)
			if tc.expectErr == "" {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if opts != tc.expect {
					t.Errorf("expected %v, but got %v", tc.expect, opts)
				}
			} else {
				if err == nil {
					t.Errorf("expected error, but got nil")
				} else if err.Error() != tc.expectErr {
					t.Errorf("expected error %q, but got %q", tc.expectErr, err.Error())
				}
			}
		})
	}
}
