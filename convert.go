package nkf

import "strings"

func Convert(str string, options string) (string, error) {
	optFlags, err := ParseOptions(options)
	if err != nil {
		return str, err
	}
	return ConvertOpt(str, optFlags), nil
}

func ConvertOpt(input string, opts ConvertOptions) string {
	builder := strings.Builder{}
	for _, ch := range input {
		if isConvertibleFullwidth(ch) && opts&FullToHalf != 0 {
			ch = ch - '\uFF00' + ' '
		} else if compat, ok := fullWidthCompat[ch]; ok && opts&FullToHalf != 0 {
			ch = compat
		} else if ch == '\u3000' && opts&FullSpaceToHalf != 0 {
			ch = ' '
		} else if ch == '\u3000' && opts&FullSpaceToTwoHalves != 0 {
			builder.WriteRune(' ')
			builder.WriteRune(' ')
			continue
		} else if isConvertibleKatakana(ch) && opts&KatakanaToHiragana != 0 {
			ch = ch - '\u30A0' + '\u3040'
		} else if isConvertibleHiragana(ch) && opts&HiraganaToKatakana != 0 {
			ch = ch - '\u3040' + '\u30A0'
		}
		builder.WriteRune(ch)
	}
	return builder.String()
}

func isConvertibleFullwidth(ch rune) bool {
	// Full width form block U+FF01 to U+FF5E except:
	// - U+FF02 FULLWIDTH QUOTATION MARK
	// - U+FF07 FULLWIDTH APOSTROPHE
	// - U+FF5E FULLWIDTH TILDE
	// because they do not belong to JIS X 0208 row 1.
	return ch >= '\uFF01' && ch <= '\uFF5D' && ch != '\uFF02' && ch != '\uFF07'
}

// fullWidthCompat originates from the "fv" constant
var fullWidthCompat = map[rune]rune{
	'\xB4':   '\x27',
	'\u2015': '\x2d',
	'\u2018': '\x60',
	'\u2019': '\x27',
	'\u201C': '\x22',
	'\u201D': '\x22',
	'\u3008': '\x3c',
	'\u3009': '\x3e',
}

func isConvertibleHiragana(ch rune) bool {
	// JIS X 0208 row 4 corresponds to U+3041 to U+3093 but U+3094 is also specially supported
	// Also, odoriji defined in row 1 are supported
	return ch >= '\u3041' && ch <= '\u3094' || ch >= '\u309D' && ch <= '\u309E'
}

func isConvertibleKatakana(ch rune) bool {
	// JIS X 0208 row 5 corresponds to U+30A1 to U+30F3 but U+3094 is also specially supported
	// Also, odoriji defined in row 1 are supported
	return ch >= '\u30A1' && ch <= '\u30F4' || ch >= '\u30FD' && ch <= '\u30FE'
}
