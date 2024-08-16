package kana

import (
	"strconv"
	"strings"
)

// ConvertOptions describes options for [Convert].
type ConvertOptions int

const (
	// HalfwidthToWide converts characters in halfwidth forms
	// to their ordinary, wide versions.
	//
	// The characters having East_Asian_Width property value of
	// H (East Asian Halfwidth) except U+20A9 WON SIGN (₩) are converted.
	// That is:
	//
	//  - U+FF61 HALFWIDTH IDEOGRAPHIC FULL STOP (｡) to U+FFBE HALFWIDTH HANGUL LETTER HIEUH (ﾾ)
	//  - U+FFC2 HALFWIDTH HANGUL LETTER A (ￂ) to U+FFC7 HALFWIDTH HANGUL LETTER E (ￇ)
	//  - U+FFCA HALFWIDTH HANGUL LETTER YEO (ￊ) to U+FFCF HALFWIDTH HANGUL LETTER OE (ￏ)
	//  - U+FFD2 HALFWIDTH HANGUL LETTER YO (ￒ) to U+FFD7 HALFWIDTH HANGUL LETTER YU (ￗ)
	//  - U+FFDA HALFWIDTH HANGUL LETTER EU (ￚ) to U+FFDC HALFWIDTH HANGUL LETTER I (ￜ)
	//  - U+FFE8 HALFWIDTH FORMS LIGHT VERTICAL (￨) to U+FFEE HALFWIDTH WHITE CIRCLE (￮)
	//
	// The conversion is roughly equivalent to NFKC but with some differences:
	//
	//  - Halfwidth Hangul letters are not fully normalized and instead
	//    converted to the corresponding letters
	//    in Hangul Compatibility Jamo block.
	//
	// The following compat flags affect the behavior of this transformation:
	//
	//  - [CompatVoicedSoundMarks]
	//  - [CompatVoicedKanaRestriction]
	//  - [CompatKeepHalfwidthHangul]
	//  - [CompatKeepHalfwidthSymbols]
	HalfwidthToWide ConvertOptions = 1 << iota
	// FullwidthToNarrow converts characters in fullwidth forms
	// to their ordinary, narrow versions.
	//
	// The characters having East_Asian_Width property value of
	// F (East Asian Fullwidth) are converted.
	// That is:
	//
	//  - U+FF01 FULLWIDTH EXCLAMATION MARK (！) to U+FF60 FULLWIDTH RIGHT WHITE PARENTHESIS (｠)
	//  - U+FFE0 FULLWIDTH CENT SIGN (￠) to U+FFE6 FULLWIDTH WON SIGN (￦)
	//
	// The conversion is roughly equivalent to NFKC but with some differences:
	//
	//  - U+FFE3 FULLWIDTH MACRON (￣) is not fully normalized and instead
	//    converted to U+00AF MACRON (¯).
	//
	// The following compat flags affect the behavior of this transformation:
	//
	//  - [CompatQuotes]
	//  - [CompatMinus]
	//  - [CompatOverline]
	//  - [CompatCurrency]
	//  - [CompatBrackets]
	//  - [CompatKeepSpaces]
	//  - [CompatDoubleSpaces]
	FullwidthToNarrow
	// KatakanaToHiragana converts katakana to hiragana.
	//
	// Consider it transformation from Script=Katakana to Script=Hiragana,
	// but there are a lot of exceptions.
	//
	// Those characters are converted to a single hiragana character:
	//
	//  - U+30A1 KATAKANA LETTER SMALL A (ァ) to U+30F6 KATAKANA LETTER SMALL KE (ヶ)
	//  - U+30FD KATAKANA ITERATION MARK (ヽ) to U+30FE KATAKANA VOICED ITERATION MARK (ヾ)
	//  - U+1B155 KATAKANA LETTER SMALL KO (𛅕)
	//  - U+1B164 KATAKANA LETTER SMALL WI (𛅤) to U+1B166 KATAKANA LETTER SMALL WO (𛅦)
	//
	// Those characters are converted to a sequence of characters:
	//
	//  - U+30F7 KATAKANA LETTER VA (ヷ) to U+30FA KATAKANA LETTER VO (ヺ)
	//
	// Those characters are not converted:
	//
	//  - U+30FF KATAKANA DIGRAPH KOTO (ヿ)
	//  - U+31F0 KATAKANA LETTER SMALL KU (ㇰ) to U+31FF KATAKANA LETTER SMALL RO (ㇿ)
	//  - U+32D0 CIRCLED KATAKANA A (㋐) to U+32FE CIRCLED KATAKANA WO (㋾)
	//  - U+3300 SQUARE APAATO (㌀) to U+3357 SQUARE WATTO (㍗)
	//  - U+1AFF0 KATAKANA LETTER MINNAN TONE-2 (𚿰) to U+1AFF3 KATAKANA LETTER MINNAN TONE-5 (𚿳)
	//  - U+1AFF5 KATAKANA LETTER MINNAN TONE-7 (𚿵) to U+1AFFB KATAKANA LETTER MINNAN NASALIZED TONE-8 (𚿻)
	//  - U+1B000 KATAKANA LETTER ARCHAIC E (𛀀)
	//  - U+1B120 KATAKANA LETTER ARCHAIC YI (𛄠) to U+1B122 KATAKANA LETTER ARCHAIC WU (𛄢)
	//  - U+1B167 KATAKANA LETTER SMALL N (𛅧)
	//
	// You need [HalfwidthToWide] to convert them to hiragana:
	//
	//  - U+FF66 HALFWIDTH KATAKANA LETTER WO (ｦ) to U+FF6F HALFWIDTH KATAKANA LETTER SMALL TU (ｯ)
	//  - U+FF71 HALFWIDTH KATAKANA LETTER A (ｱ) to U+FF9D HALFWIDTH KATAKANA LETTER N (ﾝ)
	//
	// The following compat flags affect the behavior of this transformation:
	//
	//  - [CompatKanaRestriction]
	KatakanaToHiragana
	// HiraganaToKatakana converts hiragana to katakana.
	//
	// Consider it transformation from Script=Hiragana to Script=Katakana,
	// but there are a lot of exceptions.
	//
	// Those characters are converted to a single katakana character:
	//
	//  - U+3041 HIRAGANA LETTER SMALL A (ぁ) to U+3096 HIRAGANA LETTER SMALL KE (ゖ)
	//  - U+309D HIRAGANA ITERATION MARK (ゝ) to U+309E HIRAGANA VOICED ITERATION MARK (ゞ)
	//  - U+1B132 HIRAGANA LETTER SMALL KO (𛄲)
	//  - U+1B150 HIRAGANA LETTER SMALL WI (𛅐) to U+1B152 HIRAGANA LETTER SMALL WO (𛅒)
	//
	// Those characters are not converted:
	//
	//  - U+309F HIRAGANA DIGRAPH YORI (ゟ)
	//  - U+1B001 HIRAGANA LETTER ARCHAIC YE (𛀁) to U+1B11F HIRAGANA LETTER ARCHAIC WU (𛄟)
	//  - U+1F200 SQUARE HIRAGANA HOKA (🈀)
	//
	// The following compat flags affect the behavior of this transformation:
	//
	//  - [CompatKanaRestriction]
	HiraganaToKatakana
	// CompatWideKatakanaToHalfwidth converts ordinary katakana
	// to their halfwidth forms.
	//
	// This transformation newly introduces compatibility characters
	// rather than reducing them in the input string.
	// This is against what Unicode intends to do. Therefore,
	// the entire transformation mode is considered as a compatibility option.
	//
	// If you want to normalize between fullwidth and halfwidth katakana,
	// you should use [HalfwidthToWide] instead.
	//
	// The following characters are converted:
	//
	//  - U+3001 IDEOGRAPHIC COMMA (、) to U+3002 IDEOGRAPHIC FULL STOP (。)
	//  - U+300C LEFT CORNER BRACKET (「) to U+300D RIGHT CORNER BRACKET (」)
	//  - U+3099 COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK to U+309C KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK (゜)
	//  - U+30A1 KATAKANA LETTER SMALL A (ァ) to U+30ED KATAKANA LETTER RO (ロ)
	//  - U+30EF KATAKANA LETTER WA (ワ)
	//  - U+30F2 KATAKANA LETTER WO (ヲ) to U+30F4 KATAKANA LETTER VU (ヴ)
	//  - U+30FB KATAKANA MIDDLE DOT (・) to U+30FC KATAKANA-HIRAGANA PROLONGED SOUND MARK (ー)
	//
	// When a character in the list canonically decomposes to a base character
	// and a combining voiced or semi-voiced sound mark, the transformation
	// is applied after decomposing the character.
	//
	// Note that, U+30F7 KATAKANA LETTER VA (ヷ) and U+30FA KATAKANA LETTER VO (ヺ)
	// can also be transformed this way, but they are not included in the list.
	// This is because the entire transformation exists for compatibility
	// with NKF.
	//
	// Like other compat options, this is not stable under canonical equivalence.
	CompatWideKatakanaToHalfwidth
	// CompatQuotes is a compatibility option
	// to reproduce NKF's behavior for quotes.
	//
	// Specifically, the following transformations are additionally applied
	// in [FullwidthToNarrow]:
	//
	//  - U+00B4 ACUTE ACCENT (´) → U+0027 APOSTROPHE (')
	//  - U+2018 LEFT SINGLE QUOTATION MARK (‘) → U+0060 GRAVE ACCENT (`)
	//  - U+2019 RIGHT SINGLE QUOTATION MARK (’) → U+0027 APOSTROPHE (')
	//  - U+201C LEFT DOUBLE QUOTATION MARK (“) → U+0022 QUOTATION MARK (")
	//  - U+201D RIGHT DOUBLE QUOTATION MARK (”) → U+0022 QUOTATION MARK (")
	//
	// While the following transformations are inhibited in [FullwidthToNarrow]:
	//
	//  - U+FF02 FULLWIDTH QUOTATION MARK (＂)
	//    (usually converted to U+0022 QUOTATION MARK ("))
	//  - U+FF07 FULLWIDTH APOSTROPHE (＇)
	//    (usually converted to U+0027 APOSTROPHE ('))
	CompatQuotes
	// CompatMinus is a compatibility option
	// to reproduce NKF's behavior for minus signs, hypens, and similar symbols.
	//
	// Specifically, the following transformations are applied:
	//
	//  - U+2015 HORIZONTAL BAR (―) → U+2014 EM DASH (—)
	//  - U+FF0D FULLWIDTH HYPHEN-MINUS (－) → U+2212 MINUS SIGN (−)
	//
	// and the following transformations are additionally applied
	// in [FullwidthToNarrow]:
	//
	//  - U+2014 EM DASH (—) → U+002D HYPHEN-MINUS (-)
	//  - U+2015 HORIZONTAL BAR (―) → U+002D HYPHEN-MINUS (-)
	//  - U+2212 MINUS SIGN (−) → U+002D HYPHEN-MINUS (-)
	//  - U+FF0D FULLWIDTH HYPHEN-MINUS (－) → U+002D HYPHEN-MINUS (-)
	CompatMinus
	// CompatOverline is a compatibility option
	// to reproduce NKF's behavior for overlines and similar symbols.
	//
	// Specifically, the following transformations are applied:
	//
	//  - U+FFE3 FULLWIDTH MACRON (￣) → U+203E OVERLINE (‾), which wins over
	//    [FullwidthToNarrow], where it is converted to U+00AF MACRON (¯).
	//
	//
	// Additionally, the following transformations are inhibited in
	// [FullwidthToNarrow]:
	//
	//  - U+FF5E FULLWIDTH TILDE (～)
	//    (usually converted to U+007E TILDE (~))
	CompatOverline
	// CompatCurrency is a compatibility option
	// to reproduce NKF's behavior for currency symbols.
	//
	// Specifically, the following transformations are applied regardless of
	// [FullwidthToNarrow]:
	//
	//  - U+FFE0 FULLWIDTH CENT SIGN (￠) → U+00A2 CENT SIGN (¢)
	//  - U+FFE1 FULLWIDTH POUND SIGN (￡) → U+00A3 POUND SIGN (£)
	//  - U+FFE5 FULLWIDTH YEN SIGN (￥) → U+00A5 YEN SIGN (¥)
	//
	// and the following transformations are inhibited in [FullwidthToNarrow]:
	//
	//  - U+FFE6 FULLWIDTH WON SIGN (￦)
	//    (usually converted to U+20A9 WON SIGN (₩))
	CompatCurrency
	// CompatBrackets is a compatibility option
	// to reproduce NKF's behavior for brackets and parentheses.
	//
	// Specifically, the following transformations are additionally applied
	// in [FullwidthToNarrow]:
	//
	//  - U+3008 LEFT ANGLE BRACKET (〈) → U+003C LESS-THAN SIGN (<)
	//  - U+3009 RIGHT ANGLE BRACKET (〉) → U+003E GREATER-THAN SIGN (>)
	//
	// while the following transformations are inhibited
	// in [FullwidthToNarrow]:
	//
	//  - U+FF5F FULLWIDTH LEFT WHITE PARENTHESIS (｟)
	//    (usually converted to U+2985 LEFT WHITE PARENTHESIS (⦅))
	//  - U+FF60 FULLWIDTH RIGHT WHITE PARENTHESIS (｠)
	//    (usually converted to U+2986 RIGHT WHITE PARENTHESIS (⦆))
	CompatBrackets
	// CompatOtherSymbols is a compatibility option
	// to reproduce NKF's behavior for miscellaneous symbols.
	//
	// Specifically, the following transformations are applied regardless of
	// [FullwidthToNarrow]:
	//
	//  - U+FFE2 FULLWIDTH NOT SIGN (￢) → U+00AC NOT SIGN (¬)
	//  - U+FFE4 FULLWIDTH BROKEN BAR (￤) → U+00A6 BROKEN BAR (¦)
	//
	// and the following transformations are also applied regardless of
	// [FullwidthToNarrow]:
	//
	//  - U+2225 PARALLEL TO (∥) → U+2016 DOUBLE VERTICAL LINE (‖)
	CompatOtherSymbols
	// CompatKeepSpaces is a compatibility option
	// to reproduce NKF's behavior for Ideographic Spaces.
	//
	// Specifically, the following transformations are inhibited in
	// [FullwidthToNarrow]:
	//
	//  - U+3000 IDEOGRAPHIC SPACE (　)
	//    (usually converted to U+0020 SPACE ( ))
	CompatKeepSpaces
	// CompatDoubleSpaces is a compatibility option
	// to reproduce NKF's behavior for Ideographic Spaces.
	//
	// Specifically, if this option is present along with [FullwidthToNarrow],
	// U+3000 IDEOGRAPHIC SPACE (　) is converted to two U+0020 SPACE ( ) characters.
	CompatDoubleSpaces
	// CompatVoicedSoundMarks is a compatibility option
	// to reproduce NKF's behavior for voiced and semi-voiced sound marks.
	//
	// Specifically, the following transformations are applied in [HalfwidthToWide]:
	//
	//  - U+FF9E HALFWIDTH KATAKANA VOICED SOUND MARK (ﾞ) is converted to
	//    U+309B KATAKANA-HIRAGANA VOICED SOUND MARK rather than
	//    U+3099 COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK, except when
	//    it follows ｳ, ｶ, ｷ, ｸ, ｹ, ｺ, ｻ, ｼ, ｽ, ｾ, ｿ, ﾀ, ﾁ, ﾂ, ﾃ, ﾄ, ﾊ, ﾋ, ﾌ,
	//    ﾍ, or ﾎ.
	//  - U+FF9F HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK (ﾟ) is converted to
	//    U+309C KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK rather than
	//    U+309A COMBINING KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK, except
	//    when it follows ﾊ, ﾋ, ﾌ, ﾍ, or ﾎ.
	CompatVoicedSoundMarks
	// CompatKeepHalfwidthHangul is a compatibility option
	// to reproduce NKF's behavior for halfwidth Katakana letters.
	//
	// Specifically, the following characters are transformed
	// differently in [HalfwidthToWide]:
	//
	//  - U+FF66 HALFWIDTH KATAKANA LETTER WO (ｦ) followed by
	//    U+FF9E HALFWIDTH KATAKANA VOICED SOUND MARK (ﾞ) is converted to
	//    U+30F2 KATAKANA LETTER WO (ヲ) followed by
	//    U+309B KATAKANA-HIRAGANA VOICED SOUND MARK (゛), rather than
	//    U+30FA KATAKANA LETTER VO (ヺ).
	//  - U+FF9C HALFWIDTH KATAKANA LETTER TU (ﾜ) followed by
	//    U+FF9E HALFWIDTH KATAKANA VOICED SOUND MARK (ﾞ) is converted to
	//    U+30EF KATAKANA LETTER WA (ワ) followed by
	//    U+309B KATAKANA-HIRAGANA VOICED SOUND MARK (゛), rather than
	//    U+30F7 KATAKANA LETTER VA (ヷ).
	CompatVoicedKanaRestriction
	// CompatKeepHalfwidthHangul is a compatibility option
	// to reproduce NKF's behavior for halfwidth Hangul letters.
	//
	// Specifically, the following characters are kept intact
	// in [HalfwidthToWide]:
	//
	//  - U+FFA0 HALFWIDTH HANGUL FILLER (ﾠ)
	//  - U+FFA1 HALFWIDTH HANGUL LETTER KIYEOK (ﾡ) to U+FFBE HALFWIDTH HANGUL LETTER HIEUH (ﾾ)
	//  - U+FFC2 HALFWIDTH HANGUL LETTER A (ￂ) to U+FFC7 HALFWIDTH HANGUL LETTER E (ￇ)
	//  - U+FFCA HALFWIDTH HANGUL LETTER YEO (ￊ) to U+FFCF HALFWIDTH HANGUL LETTER OE (ￏ)
	//  - U+FFD2 HALFWIDTH HANGUL LETTER YO (ￒ) to U+FFD7 HALFWIDTH HANGUL LETTER YU (ￗ)
	//  - U+FFDA HALFWIDTH HANGUL LETTER EU (ￚ) to U+FFDC HALFWIDTH HANGUL LETTER I (ￜ)
	CompatKeepHalfwidthHangul
	// CompatKeepHalfwidthSymbols is a compatibility option
	// to reproduce NKF's behavior for halfwidth symbols.
	//
	// Specifically, the following characters are kept intact
	// in [HalfwidthToWide]:
	//
	//  - U+FFE8 HALFWIDTH FORMS LIGHT VERTICAL (￨)
	//  - U+FFE9 HALFWIDTH LEFTWARDS ARROW (￩)
	//  - U+FFEA HALFWIDTH UPWARDS ARROW (￪)
	//  - U+FFEB HALFWIDTH RIGHTWARDS ARROW (￫)
	//  - U+FFEC HALFWIDTH DOWNWARDS ARROW (￬)
	//  - U+FFED HALFWIDTH BLACK SQUARE (￭)
	//  - U+FFEE HALFWIDTH WHITE CIRCLE (￮)
	CompatKeepHalfwidthSymbols
	// CompatKanaRestriction is a compatibility option
	// to reproduce NKF's behavior for hiragana and katakana.
	//
	// Specifically, the following characters are kept intact
	// in [KatakanaToHiragana]:
	//
	//  - U+30F5 KATAKANA LETTER SMALL KA (ヵ)
	//  - U+30F6 KATAKANA LETTER SMALL KE (ヶ)
	//  - U+30F7 KATAKANA LETTER VA (ヷ)
	//  - U+30F8 KATAKANA LETTER VI (ヸ)
	//  - U+30F9 KATAKANA LETTER VE (ヹ)
	//  - U+30FA KATAKANA LETTER VO (ヺ)
	//  - U+1B155 KATAKANA LETTER SMALL KO (𛅕)
	//  - U+1B164 KATAKANA LETTER SMALL WI (𛅤)
	//  - U+1B165 KATAKANA LETTER SMALL WE (𛅥)
	//  - U+1B166 KATAKANA LETTER SMALL WO (𛅦)
	//
	// and the following characters are kept intact in [HiraganaToKatakana]:
	//
	//  - U+3095 HIRAGANA LETTER SMALL KA (ゕ)
	//  - U+3096 HIRAGANA LETTER SMALL KE (ゖ)
	//  - U+1B132 HIRAGANA LETTER SMALL KO (𛄲)
	//  - U+1B150 HIRAGANA LETTER SMALL WI (𛅐)
	//  - U+1B151 HIRAGANA LETTER SMALL WE (𛅑)
	//  - U+1B152 HIRAGANA LETTER SMALL WO (𛅒)
	CompatKanaRestriction
)

func (o ConvertOptions) Normalize() ConvertOptions {
	if o&FullwidthToNarrow == 0 {
		o &^= CompatQuotes | CompatBrackets | CompatKeepSpaces | CompatDoubleSpaces
	}
	if o&CompatKeepSpaces != 0 {
		o &^= CompatDoubleSpaces
	}
	if o&HalfwidthToWide == 0 {
		o &^= CompatVoicedSoundMarks | CompatVoicedKanaRestriction | CompatKeepHalfwidthHangul | CompatKeepHalfwidthSymbols
	}
	if o&(KatakanaToHiragana|HiraganaToKatakana) == 0 {
		o &^= CompatKanaRestriction
	}
	return o
}

var flagNames = []struct {
	name string
	flag ConvertOptions
	mask ConvertOptions
}{
	{"HalfwidthToWide", HalfwidthToWide, HalfwidthToWide},
	{"FullwidthToNarrow", FullwidthToNarrow, FullwidthToNarrow},
	{"KatakanaToHiragana", KatakanaToHiragana, KatakanaToHiragana},
	{"HiraganaToKatakana", HiraganaToKatakana, HiraganaToKatakana},
	{"CompatWideKatakanaToHalfwidth", CompatWideKatakanaToHalfwidth, CompatWideKatakanaToHalfwidth},
	{"CompatQuotes", CompatQuotes, CompatQuotes},
	{"CompatMinus", CompatMinus, CompatMinus},
	{"CompatOverline", CompatOverline, CompatOverline},
	{"CompatCurrency", CompatCurrency, CompatCurrency},
	{"CompatBrackets", CompatBrackets, CompatBrackets},
	{"CompatOtherSymbols", CompatOtherSymbols, CompatOtherSymbols},
	{"CompatKeepSpaces", CompatKeepSpaces, CompatKeepSpaces},
	{"CompatDoubleSpaces", CompatDoubleSpaces, CompatDoubleSpaces},
	{"CompatVoicedSoundMarks", CompatVoicedSoundMarks, CompatVoicedSoundMarks},
	{"CompatVoicedKanaRestriction", CompatVoicedKanaRestriction, CompatVoicedKanaRestriction},
	{"CompatKeepHalfwidthHangul", CompatKeepHalfwidthHangul, CompatKeepHalfwidthHangul},
	{"CompatKeepHalfwidthSymbols", CompatKeepHalfwidthSymbols, CompatKeepHalfwidthSymbols},
	{"CompatKanaRestriction", CompatKanaRestriction, CompatKanaRestriction},
}

func (o ConvertOptions) String() string {
	var names []string
	for _, n := range flagNames {
		if o&n.mask == n.flag {
			names = append(names, n.name)
			o &^= n.mask
		}
	}
	if o != 0 {
		names = append(names, "0x"+strconv.FormatInt(int64(o), 16))
	} else if len(names) == 0 {
		return "0"
	}

	return strings.Join(names, " | ")
}
