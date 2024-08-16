// Package kana provides transformation between:
//
//   - Fullwidth and halfwidth characters
//   - Katakana and hiragana
//
// # Example
//
//	package main
//
//	import (
//		"fmt"
//
//		"github.com/wantedly/kana-go"
//	)
//
//	func main() {
//		str := kana.Convert("ＡＢＣ　ＤＥＦ", kana.FullwidthToNarrow)
//		fmt.Println(str) // Output: ABC DEF
//	}
package kana

// Convert converts a string with the given options.
func Convert(input string, opts ConvertOptions) string {
	opts = opts.Normalize()
	strm := stringStream(input)

	strm = convertUnconditionalCompat(strm, opts)
	// Full <-> Half conversion
	strm = doWidthNormalization(strm, opts)
	strm = doKanaConversion(strm, opts)

	return strm.readAll()
}

func convertUnconditionalCompat(strm *stream, opts ConvertOptions) *stream {
	if opts&(CompatMinus|CompatOverline|CompatCurrency|CompatOtherSymbols) == 0 {
		return strm
	}
	return mapStream(strm, func(ch rune) rune {
		if opts&CompatMinus != 0 {
			switch ch {
			case '\u2015':
				return '\u2014'
			case '\uFF0D':
				return '\u2212'
			}
		}
		if opts&CompatOverline != 0 {
			switch ch {
			case '\uFFE3':
				return '\u203E'
			}
		}
		if opts&CompatCurrency != 0 {
			switch ch {
			case '\uFFE0':
				return '\u00A2'
			case '\uFFE1':
				return '\u00A3'
			case '\uFFE5':
				return '\u00A5'
			}
		}
		if opts&CompatOtherSymbols != 0 {
			switch ch {
			case '\u2225':
				return '\u2016'
			case '\uFFE2':
				return '\u00AC'
			case '\uFFE4':
				return '\u00A6'
			}
		}
		return ch
	})
}

func doWidthNormalization(strm *stream, opts ConvertOptions) *stream {
	if opts&(FullwidthToNarrow|CompatWideKatakanaToHalfwidth|HalfwidthToWide) == 0 {
		return strm
	}
	return newStream(func(buf *[]rune) {
		ch, ok := strm.readOne()
		if !ok {
			return
		}

		if ok := convertFullwidthToNarrow(ch, buf, opts); ok {
			// Do nothing
		} else if ok := convertWideKatakanaToHalfwidth(ch, buf, opts); ok {
			// Do nothing
		} else if ok := convertHalfwidthToWide(ch, strm, buf, opts); ok {
			// Do nothing
		} else {
			*buf = append(*buf, ch)
		}
	})
}

func convertFullwidthToNarrow(ch rune, buf *[]rune, opts ConvertOptions) bool {
	if opts&FullwidthToNarrow == 0 {
		return false
	}
	if opts&CompatQuotes != 0 {
		switch ch {
		case '\u00B4', '\u2019':
			*buf = append(*buf, '\'')
			return true
		case '\u2018':
			*buf = append(*buf, '`')
			return true
		case '\u201C', '\u201D':
			*buf = append(*buf, '"')
			return true
		case '\uFF02', '\uFF07':
			return false
		}
	}
	if opts&CompatMinus != 0 {
		switch ch {
		case '\u2014', '\u2212':
			// '\u2015' also falls here because it is converted to '\u2014' in convertUnconditionalCompat
			// '\uFF0D' also falls here because it is converted to '\uFF0D' in convertUnconditionalCompat
			*buf = append(*buf, '-')
			return true
		}
	}
	if opts&CompatOverline != 0 {
		switch ch {
		case '\uFF5E':
			return false
		}
	}
	if opts&CompatCurrency != 0 {
		switch ch {
		case '\uFFE6':
			return false
		}
	}
	if opts&CompatBrackets != 0 {
		switch ch {
		case '\u3008':
			*buf = append(*buf, '<')
			return true
		case '\u3009':
			*buf = append(*buf, '>')
			return true
		case '\uFF5F', '\uFF60':
			return false
		}
	}
	if opts&CompatKeepSpaces != 0 && ch == '\u3000' {
		return false
	} else if opts&CompatDoubleSpaces != 0 && ch == '\u3000' {
		*buf = append(*buf, ' ', ' ')
		return true
	}
	if ch >= '\uFF01' && ch <= '\uFF5E' {
		*buf = append(*buf, ch-'\uFF00'+' ')
		return true
	} else if ch == '\u3000' {
		*buf = append(*buf, ' ')
		return true
	} else if mapped, ok := fullwidthMap[ch]; ok {
		*buf = append(*buf, mapped)
		return true
	}
	return false
}

var fullwidthMap = map[rune]rune{
	'\uFF5F': '\u2985',
	'\uFF60': '\u2986',
	'\uFFE0': '\u00A2',
	'\uFFE1': '\u00A3',
	'\uFFE2': '\u00AC',
	// NOTE: this is different from how it normalizes to in NFKC, which is the sequence U+0020 U+0304
	'\uFFE3': '\u00AF',
	'\uFFE4': '\u00A6',
	'\uFFE5': '\u00A5',
	'\uFFE6': '\u20A9',
}

func convertWideKatakanaToHalfwidth(ch rune, buf *[]rune, opts ConvertOptions) bool {
	if opts&CompatWideKatakanaToHalfwidth == 0 {
		return false
	}
	if mapped, ok := fullwidthKatakanaTable[ch]; ok {
		for _, mappedCh := range mapped {
			*buf = append(*buf, mappedCh)
		}
		return true
	}
	return false
}

var fullwidthKatakanaTable = map[rune]string{
	'\u3001': "\uFF64",
	'\u3002': "\uFF61",
	'\u300C': "\uFF62",
	'\u300D': "\uFF63",
	'\u3099': "\uFF9E",
	'\u309A': "\uFF9F",
	'\u309B': "\uFF9E",
	'\u309C': "\uFF9F",
	'\u30A1': "\uFF67",
	'\u30A2': "\uFF71",
	'\u30A3': "\uFF68",
	'\u30A4': "\uFF72",
	'\u30A5': "\uFF69",
	'\u30A6': "\uFF73",
	'\u30A7': "\uFF6A",
	'\u30A8': "\uFF74",
	'\u30A9': "\uFF6B",
	'\u30AA': "\uFF75",
	'\u30AB': "\uFF76",
	'\u30AC': "\uFF76\uFF9E",
	'\u30AD': "\uFF77",
	'\u30AE': "\uFF77\uFF9E",
	'\u30AF': "\uFF78",
	'\u30B0': "\uFF78\uFF9E",
	'\u30B1': "\uFF79",
	'\u30B2': "\uFF79\uFF9E",
	'\u30B3': "\uFF7A",
	'\u30B4': "\uFF7A\uFF9E",
	'\u30B5': "\uFF7B",
	'\u30B6': "\uFF7B\uFF9E",
	'\u30B7': "\uFF7C",
	'\u30B8': "\uFF7C\uFF9E",
	'\u30B9': "\uFF7D",
	'\u30BA': "\uFF7D\uFF9E",
	'\u30BB': "\uFF7E",
	'\u30BC': "\uFF7E\uFF9E",
	'\u30BD': "\uFF7F",
	'\u30BE': "\uFF7F\uFF9E",
	'\u30BF': "\uFF80",
	'\u30C0': "\uFF80\uFF9E",
	'\u30C1': "\uFF81",
	'\u30C2': "\uFF81\uFF9E",
	'\u30C3': "\uFF6F",
	'\u30C4': "\uFF82",
	'\u30C5': "\uFF82\uFF9E",
	'\u30C6': "\uFF83",
	'\u30C7': "\uFF83\uFF9E",
	'\u30C8': "\uFF84",
	'\u30C9': "\uFF84\uFF9E",
	'\u30CA': "\uFF85",
	'\u30CB': "\uFF86",
	'\u30CC': "\uFF87",
	'\u30CD': "\uFF88",
	'\u30CE': "\uFF89",
	'\u30CF': "\uFF8A",
	'\u30D0': "\uFF8A\uFF9E",
	'\u30D1': "\uFF8A\uFF9F",
	'\u30D2': "\uFF8B",
	'\u30D3': "\uFF8B\uFF9E",
	'\u30D4': "\uFF8B\uFF9F",
	'\u30D5': "\uFF8C",
	'\u30D6': "\uFF8C\uFF9E",
	'\u30D7': "\uFF8C\uFF9F",
	'\u30D8': "\uFF8D",
	'\u30D9': "\uFF8D\uFF9E",
	'\u30DA': "\uFF8D\uFF9F",
	'\u30DB': "\uFF8E",
	'\u30DC': "\uFF8E\uFF9E",
	'\u30DD': "\uFF8E\uFF9F",
	'\u30DE': "\uFF8F",
	'\u30DF': "\uFF90",
	'\u30E0': "\uFF91",
	'\u30E1': "\uFF92",
	'\u30E2': "\uFF93",
	'\u30E3': "\uFF6C",
	'\u30E4': "\uFF94",
	'\u30E5': "\uFF6D",
	'\u30E6': "\uFF95",
	'\u30E7': "\uFF6E",
	'\u30E8': "\uFF96",
	'\u30E9': "\uFF97",
	'\u30EA': "\uFF98",
	'\u30EB': "\uFF99",
	'\u30EC': "\uFF9A",
	'\u30ED': "\uFF9B",
	'\u30EF': "\uFF9C",
	'\u30F2': "\uFF66",
	'\u30F3': "\uFF9D",
	'\u30F4': "\uFF73\uFF9E",
	'\u30FB': "\uFF65",
	'\u30FC': "\uFF70",
}

func convertHalfwidthToWide(ch rune, strm *stream, buf *[]rune, opts ConvertOptions) bool {
	if opts&HalfwidthToWide == 0 {
		return false
	}
	if opts&CompatVoicedSoundMarks != 0 {
		// Use a non-combining version
		switch ch {
		case '\uFF9E':
			*buf = append(*buf, '\u309B')
			return true
		case '\uFF9F':
			*buf = append(*buf, '\u309C')
			return true
		}
	}
	if opts&CompatKeepHalfwidthHangul != 0 && ('\uFFA0' <= ch && ch <= '\uFFDC') {
		return false
	}
	if opts&CompatKeepHalfwidthSymbols != 0 && ('\uFFE0' <= ch && ch <= '\uFFEF') {
		return false
	}
	if ch >= '\uFF61' && ch <= '\uFFEF' {
		mapped, ok := halfwidthMap[ch]
		if ok {
			nextCh, _ := strm.peekOne()
			if nextCh == '\uFF9E' {
				if opts&CompatVoicedKanaRestriction != 0 && (ch == '\uFF66' || ch == '\uFF9C') {
					*buf = append(*buf, mapped)
					return true
				}
				if voiced, ok := halfwidthVoicedKatakanaTable[ch]; ok {
					strm.consume(1)
					*buf = append(*buf, voiced)
					return true
				}
			} else if nextCh == '\uFF9F' {
				if semiVoiced, ok := halfwidthSemiVoicedKatakanaTable[ch]; ok {
					strm.consume(1)
					*buf = append(*buf, semiVoiced)
					return true
				}
			}

			*buf = append(*buf, mapped)
			return true
		}
	}
	return false
}

var halfwidthMap = map[rune]rune{
	// Katakana and relevant punctuations
	'\uFF61': '\u3002',
	'\uFF62': '\u300C',
	'\uFF63': '\u300D',
	'\uFF64': '\u3001',
	'\uFF65': '\u30FB',
	'\uFF66': '\u30F2',
	'\uFF67': '\u30A1',
	'\uFF68': '\u30A3',
	'\uFF69': '\u30A5',
	'\uFF6A': '\u30A7',
	'\uFF6B': '\u30A9',
	'\uFF6C': '\u30E3',
	'\uFF6D': '\u30E5',
	'\uFF6E': '\u30E7',
	'\uFF6F': '\u30C3',
	'\uFF70': '\u30FC',
	'\uFF71': '\u30A2',
	'\uFF72': '\u30A4',
	'\uFF73': '\u30A6',
	'\uFF74': '\u30A8',
	'\uFF75': '\u30AA',
	'\uFF76': '\u30AB',
	'\uFF77': '\u30AD',
	'\uFF78': '\u30AF',
	'\uFF79': '\u30B1',
	'\uFF7A': '\u30B3',
	'\uFF7B': '\u30B5',
	'\uFF7C': '\u30B7',
	'\uFF7D': '\u30B9',
	'\uFF7E': '\u30BB',
	'\uFF7F': '\u30BD',
	'\uFF80': '\u30BF',
	'\uFF81': '\u30C1',
	'\uFF82': '\u30C4',
	'\uFF83': '\u30C6',
	'\uFF84': '\u30C8',
	'\uFF85': '\u30CA',
	'\uFF86': '\u30CB',
	'\uFF87': '\u30CC',
	'\uFF88': '\u30CD',
	'\uFF89': '\u30CE',
	'\uFF8A': '\u30CF',
	'\uFF8B': '\u30D2',
	'\uFF8C': '\u30D5',
	'\uFF8D': '\u30D8',
	'\uFF8E': '\u30DB',
	'\uFF8F': '\u30DE',
	'\uFF90': '\u30DF',
	'\uFF91': '\u30E0',
	'\uFF92': '\u30E1',
	'\uFF93': '\u30E2',
	'\uFF94': '\u30E4',
	'\uFF95': '\u30E6',
	'\uFF96': '\u30E8',
	'\uFF97': '\u30E9',
	'\uFF98': '\u30EA',
	'\uFF99': '\u30EB',
	'\uFF9A': '\u30EC',
	'\uFF9B': '\u30ED',
	'\uFF9C': '\u30EF',
	'\uFF9D': '\u30F3',
	// Use combining version for the voiced and semi-voiced marks,
	// although the halfwidth forms are non-combining.
	// It aligns with NFKC behavior and is justified by how
	// one would expect the mark to behave.
	// In NKF compat mode, the non-combining version is used,
	// but it tries to use the precomposed form, if it is available.
	'\uFF9E': '\u3099',
	'\uFF9F': '\u309A',
	// These are halfwidth versions of Hangul **Compatibility** Jamos
	// rather than the Unicode proper Hangul Jamos.
	// They lack distinction between L (leading consonant; choseong)
	// and T (trailing consonant; jungseong) and therefore it is difficult to
	// determine the syllable boundary
	// just as Unicode does for the proper Hangul Jamos.
	// They are merely for round-trip compatibility with legacy encodings.
	// To align with how Unicode handles these characters, we do not try
	// to determine the consonant type or compose them into a precomposed syllable.
	'\uFFA0': '\u3164',
	'\uFFA1': '\u3131',
	'\uFFA2': '\u3132',
	'\uFFA3': '\u3133',
	'\uFFA4': '\u3134',
	'\uFFA5': '\u3135',
	'\uFFA6': '\u3136',
	'\uFFA7': '\u3137',
	'\uFFA8': '\u3138',
	'\uFFA9': '\u3139',
	'\uFFAA': '\u313A',
	'\uFFAB': '\u313B',
	'\uFFAC': '\u313C',
	'\uFFAD': '\u313D',
	'\uFFAE': '\u313E',
	'\uFFAF': '\u313F',
	'\uFFB0': '\u3140',
	'\uFFB1': '\u3141',
	'\uFFB2': '\u3142',
	'\uFFB3': '\u3143',
	'\uFFB4': '\u3144',
	'\uFFB5': '\u3145',
	'\uFFB6': '\u3146',
	'\uFFB7': '\u3147',
	'\uFFB8': '\u3148',
	'\uFFB9': '\u3149',
	'\uFFBA': '\u314A',
	'\uFFBB': '\u314B',
	'\uFFBC': '\u314C',
	'\uFFBD': '\u314D',
	'\uFFBE': '\u314E',
	'\uFFC2': '\u314F',
	'\uFFC3': '\u3150',
	'\uFFC4': '\u3151',
	'\uFFC5': '\u3152',
	'\uFFC6': '\u3153',
	'\uFFC7': '\u3154',
	'\uFFCA': '\u3155',
	'\uFFCB': '\u3156',
	'\uFFCC': '\u3157',
	'\uFFCD': '\u3158',
	'\uFFCE': '\u3159',
	'\uFFCF': '\u315A',
	'\uFFD2': '\u315B',
	'\uFFD3': '\u315C',
	'\uFFD4': '\u315D',
	'\uFFD5': '\u315E',
	'\uFFD6': '\u315F',
	'\uFFD7': '\u3160',
	'\uFFDA': '\u3161',
	'\uFFDB': '\u3162',
	'\uFFDC': '\u3163',
	// Halfwidth forms of symbols
	// U+FFE8 is HALFWIDTH FORMS LIGHT VERTICAL while U+2502 is BOX DRAWINGS LIGHT VERTICAL
	// I guess, uh, this is probably what it is meant to be mapped to?
	'\uFFE8': '\u2502',
	'\uFFE9': '\u2190',
	'\uFFEA': '\u2191',
	'\uFFEB': '\u2192',
	'\uFFEC': '\u2193',
	'\uFFED': '\u25A0',
	'\uFFEE': '\u25CB',
}

var halfwidthVoicedKatakanaTable = map[rune]rune{
	'\uFF66': '\u30FA',
	'\uFF73': '\u30F4',
	'\uFF76': '\u30AC',
	'\uFF77': '\u30AE',
	'\uFF78': '\u30B0',
	'\uFF79': '\u30B2',
	'\uFF7A': '\u30B4',
	'\uFF7B': '\u30B6',
	'\uFF7C': '\u30B8',
	'\uFF7D': '\u30BA',
	'\uFF7E': '\u30BC',
	'\uFF7F': '\u30BE',
	'\uFF80': '\u30C0',
	'\uFF81': '\u30C2',
	'\uFF82': '\u30C5',
	'\uFF83': '\u30C7',
	'\uFF84': '\u30C9',
	'\uFF8A': '\u30D0',
	'\uFF8B': '\u30D3',
	'\uFF8C': '\u30D6',
	'\uFF8D': '\u30D9',
	'\uFF8E': '\u30DC',
	'\uFF9C': '\u30F7',
}

var halfwidthSemiVoicedKatakanaTable = map[rune]rune{
	'\uFF8A': '\u30D1',
	'\uFF8B': '\u30D4',
	'\uFF8C': '\u30D7',
	'\uFF8D': '\u30DA',
	'\uFF8E': '\u30DD',
}

func doKanaConversion(strm *stream, opts ConvertOptions) *stream {
	if opts&(KatakanaToHiragana|HiraganaToKatakana) == 0 {
		return strm
	}
	return newStream(func(buf *[]rune) {
		ch, ok := strm.readOne()
		if !ok {
			return
		}

		if ok := convertKatakanaToHiragana(ch, buf, opts); ok {
			// Do nothing
		} else if ok := convertHiraganaToKatakana(ch, buf, opts); ok {
			// Do nothing
		} else {
			*buf = append(*buf, ch)
		}
	})
}

func convertKatakanaToHiragana(ch rune, buf *[]rune, opts ConvertOptions) bool {
	if opts&KatakanaToHiragana == 0 {
		return false
	}
	if opts&CompatKanaRestriction != 0 && !(ch >= '\u30A1' && ch <= '\u30F4' || ch >= '\u30FD' && ch <= '\u30FE') {
		return false
	}

	if ch >= '\u30A1' && ch <= '\u30F4' || ch >= '\u30F5' && ch <= '\u30F6' || ch >= '\u30FD' && ch <= '\u30FE' {
		*buf = append(*buf, ch-'\u30A0'+'\u3040')
		return true
	}
	switch ch {
	case '\u30F7':
		*buf = append(*buf, '\u308F', '\u3099')
		return true
	case '\u30F8':
		*buf = append(*buf, '\u3090', '\u3099')
		return true
	case '\u30F9':
		*buf = append(*buf, '\u3091', '\u3099')
		return true
	case '\u30FA':
		*buf = append(*buf, '\u3092', '\u3099')
		return true
	case '\U0001B155':
		*buf = append(*buf, '\U0001B132')
		return true
	case '\U0001B164':
		*buf = append(*buf, '\U0001B150')
		return true
	case '\U0001B165':
		*buf = append(*buf, '\U0001B151')
		return true
	case '\U0001B166':
		*buf = append(*buf, '\U0001B152')
		return true
	}
	return false
}

func convertHiraganaToKatakana(ch rune, buf *[]rune, opts ConvertOptions) bool {
	if opts&HiraganaToKatakana == 0 {
		return false
	}
	if opts&CompatKanaRestriction != 0 && !(ch >= '\u3041' && ch <= '\u3094' || ch >= '\u309D' && ch <= '\u309E') {
		return false
	}

	if ch >= '\u3041' && ch <= '\u3094' || ch >= '\u3095' && ch <= '\u3096' || ch >= '\u309D' && ch <= '\u309E' {
		*buf = append(*buf, ch-'\u3040'+'\u30A0')
		return true
	}
	switch ch {
	case '\U0001B132':
		*buf = append(*buf, '\U0001B155')
		return true
	case '\U0001B150':
		*buf = append(*buf, '\U0001B164')
		return true
	case '\U0001B151':
		*buf = append(*buf, '\U0001B165')
		return true
	case '\U0001B152':
		*buf = append(*buf, '\U0001B166')
		return true
	}
	return false
}
