package kana_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wantedly/kana-go"
)

func TestConvertOptionsNormalize(t *testing.T) {
	testcases := []struct {
		name     string
		input    kana.ConvertOptions
		expected kana.ConvertOptions
	}{
		{
			name:     "empty",
			input:    0,
			expected: 0,
		},
		{
			name:     "CompatQuotes, without FullwidthToNarrow",
			input:    kana.CompatQuotes,
			expected: 0,
		},
		{
			name:     "CompatQuotes, with FullwidthToNarrow",
			input:    kana.FullwidthToNarrow | kana.CompatQuotes,
			expected: kana.FullwidthToNarrow | kana.CompatQuotes,
		},
		{
			name:     "CompatBrackets, without FullwidthToNarrow",
			input:    kana.CompatBrackets,
			expected: 0,
		},
		{
			name:     "CompatBrackets, with FullwidthToNarrow",
			input:    kana.FullwidthToNarrow | kana.CompatBrackets,
			expected: kana.FullwidthToNarrow | kana.CompatBrackets,
		},
		{
			name:     "CompatKeepSpaces, without FullwidthToNarrow",
			input:    kana.CompatKeepSpaces,
			expected: 0,
		},
		{
			name:     "CompatKeepSpaces, with FullwidthToNarrow",
			input:    kana.FullwidthToNarrow | kana.CompatKeepSpaces,
			expected: kana.FullwidthToNarrow | kana.CompatKeepSpaces,
		},
		{
			name:     "CompatDoubleSpaces, without FullwidthToNarrow",
			input:    kana.CompatDoubleSpaces,
			expected: 0,
		},
		{
			name:     "CompatDoubleSpaces, with FullwidthToNarrow",
			input:    kana.FullwidthToNarrow | kana.CompatDoubleSpaces,
			expected: kana.FullwidthToNarrow | kana.CompatDoubleSpaces,
		},
		{
			name:     "CompatDoubleSpaces, without CompatKeepSpaces",
			input:    kana.FullwidthToNarrow | kana.CompatDoubleSpaces,
			expected: kana.FullwidthToNarrow | kana.CompatDoubleSpaces,
		},
		{
			name:     "CompatDoubleSpaces, with CompatKeepSpaces",
			input:    kana.FullwidthToNarrow | kana.CompatKeepSpaces | kana.CompatDoubleSpaces,
			expected: kana.FullwidthToNarrow | kana.CompatKeepSpaces,
		},
		{
			name:     "CompatVoicedSoundMarks, without HalfwidthToWide",
			input:    kana.CompatVoicedSoundMarks,
			expected: 0,
		},
		{
			name:     "CompatVoicedSoundMarks, with HalfwidthToWide",
			input:    kana.HalfwidthToWide | kana.CompatVoicedSoundMarks,
			expected: kana.HalfwidthToWide | kana.CompatVoicedSoundMarks,
		},
		{
			name:     "CompatKeepHalfwidthHangul, without HalfwidthToWide",
			input:    kana.CompatKeepHalfwidthHangul,
			expected: 0,
		},
		{
			name:     "CompatKeepHalfwidthHangul, with HalfwidthToWide",
			input:    kana.HalfwidthToWide | kana.CompatKeepHalfwidthHangul,
			expected: kana.HalfwidthToWide | kana.CompatKeepHalfwidthHangul,
		},
		{
			name:     "CompatKeepHalfwidthSymbols, without HalfwidthToWide",
			input:    kana.CompatKeepHalfwidthSymbols,
			expected: 0,
		},
		{
			name:     "CompatKeepHalfwidthSymbols, with HalfwidthToWide",
			input:    kana.HalfwidthToWide | kana.CompatKeepHalfwidthSymbols,
			expected: kana.HalfwidthToWide | kana.CompatKeepHalfwidthSymbols,
		},
		{
			name:     "CompatKanaRestriction, without KatakanaToHiragana or HiraganaToKatakana",
			input:    kana.CompatKanaRestriction,
			expected: 0,
		},
		{
			name:     "CompatKanaRestriction, with KatakanaToHiragana",
			input:    kana.KatakanaToHiragana | kana.CompatKanaRestriction,
			expected: kana.KatakanaToHiragana | kana.CompatKanaRestriction,
		},
		{
			name:     "CompatKanaRestriction, with HiraganaToKatakana",
			input:    kana.HiraganaToKatakana | kana.CompatKanaRestriction,
			expected: kana.HiraganaToKatakana | kana.CompatKanaRestriction,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.input.Normalize()
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertOptionsString(t *testing.T) {
	testcases := []struct {
		name     string
		opts     kana.ConvertOptions
		expected string
	}{
		{
			name:     "empty",
			opts:     0,
			expected: "0",
		},
		{
			name:     "unknown bit",
			opts:     1 << 30,
			expected: "0x40000000",
		},
		{
			name:     "one",
			opts:     kana.HalfwidthToWide,
			expected: "HalfwidthToWide",
		},
		{
			name:     "one plus extra",
			opts:     kana.HalfwidthToWide | (1 << 30),
			expected: "HalfwidthToWide | 0x40000000",
		},
		{
			name:     "two",
			opts:     kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expected: "HalfwidthToWide | KatakanaToHiragana",
		},
		{
			name:     "two plus extra",
			opts:     kana.HalfwidthToWide | kana.KatakanaToHiragana | (1 << 30),
			expected: "HalfwidthToWide | KatakanaToHiragana | 0x40000000",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.opts.String()
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
