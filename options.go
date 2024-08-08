package nkf

type ConvertOptions int

const (
	// KatakanaToHiragana converts katakana to hiragana.
	// Corresponds to hira_f & 1.
	KatakanaToHiragana ConvertOptions = 1 << iota
	// HiraganaToKatakana converts hiragana to katakana.
	// Corresponds to hira_f & 2.
	HiraganaToKatakana
	// HalfKanaToFull converts half-width katakana to full-width katakana.
	// Corresponds to x0201_f.
	HalfKanaToFull
	// FullToHalf converts full-width characters to half-width characters.
	// Corresponds to alpha_f & 1.
	FullToHalf
	// FullSpaceToHalf converts ideographic space to space.
	// Corresponds to alpha_f & 2.
	FullSpaceToHalf
	// FullSpaceToTwoHalves converts ideographic space to two spaces.
	// Corresponds to alpha_f & 4.
	FullSpaceToTwoHalves
)

func (o ConvertOptions) Normalize() ConvertOptions {
	if o&(FullSpaceToHalf|FullSpaceToTwoHalves) != 0 {
		o |= FullToHalf
	}
	if o&FullSpaceToHalf != 0 {
		o &= ^FullSpaceToTwoHalves
	}
	return o
}
