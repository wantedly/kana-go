package nkf

import (
	"fmt"
	"strings"

	"github.com/wantedly/kana-go"
)

// ParseOptions parses the given text and returns ConvertOptions.
func ParseOptions(text string) (kana.ConvertOptions, error) {
	p := parser{
		halfwidthToWide: true,
	}
	if err := p.parseOptions(text); err != nil {
		return 0, err
	}
	if !p.utf8Output {
		return 0, fmt.Errorf("-w is required")
	}
	if !p.utf8Input {
		return 0, fmt.Errorf("-W is required")
	}
	if !p.noMime {
		return 0, fmt.Errorf("-m0 is required")
	}
	return p.toOptions(), nil
}

type parser struct {
	utf8Output               bool
	utf8Input                bool
	noMime                   bool
	katakanaToHiragana       bool
	hiraganaToKatakana       bool
	fullwidthToNarrow        bool
	ideographicSpaceToNarrow bool
	doubleIdeographicSpace   bool
	wideKatakanaToHalfwidth  bool
	halfwidthToWide          bool
}

func (p *parser) toOptions() kana.ConvertOptions {
	opts := kana.CompatMinus | kana.CompatOverline | kana.CompatCurrency | kana.CompatOtherSymbols
	if p.katakanaToHiragana {
		opts |= kana.KatakanaToHiragana | kana.CompatKanaRestriction
	}
	if p.hiraganaToKatakana {
		opts |= kana.HiraganaToKatakana | kana.CompatKanaRestriction
	}
	if p.fullwidthToNarrow {
		opts |= kana.FullwidthToNarrow | kana.CompatQuotes | kana.CompatBrackets | kana.CompatKeepSpaces
	}
	if p.ideographicSpaceToNarrow {
		opts &= ^kana.CompatKeepSpaces
	} else if p.doubleIdeographicSpace {
		opts &= ^kana.CompatKeepSpaces
		opts |= kana.CompatDoubleSpaces
	}
	if p.wideKatakanaToHalfwidth {
		opts |= kana.CompatWideKatakanaToHalfwidth
	}
	if p.halfwidthToWide {
		opts |= kana.HalfwidthToWide | kana.CompatVoicedSoundMarks | kana.CompatVoicedKanaRestriction | kana.CompatKeepHalfwidthHangul | kana.CompatKeepHalfwidthSymbols
	}
	return opts
}

func (p *parser) parseOptions(text string) error {
	// See `options` in nkf.c

	parts := strings.Split(text, " ")
	for _, part := range parts {
		if part == "" {
			continue
		}
		if part[0] != '-' || len(part) <= 1 {
			return fmt.Errorf("invalid option: %s", part)
		}
		if part[1] == '-' {
			if len(part) <= 2 {
				return fmt.Errorf("invalid option: %s", part)
			}
			longOpt, ok := longOptions[part[2:]]
			if !ok {
				return fmt.Errorf("invalid option: %s", part)
			}
			err := p.parseShortOptions(longOpt)
			if err != nil {
				return err
			}
			continue
		}
		err := p.parseShortOptions(part[1:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseShortOptions(text string) error {
	bytes := []byte(text)
	i := 0
	for i < len(bytes) {
		j := i + 1
		for j < len(bytes) {
			switch text[i : j+1] {
			case "mB", "mQ", "mN", "mS":
				j += 1
				continue
			}
			if '0' <= bytes[j] && bytes[j] <= '9' {
				j += 1
				continue
			} else {
				break
			}
		}
		group := text[i:j]
		i = j
		switch group {
		case "h", "h1":
			p.katakanaToHiragana = true
		case "h2":
			p.hiraganaToKatakana = true
		case "h3":
			p.katakanaToHiragana = true
			p.hiraganaToKatakana = true
		case "w", "w8":
			p.utf8Output = true
		case "W", "W8":
			p.utf8Input = true
		case "Z", "Z0":
			p.fullwidthToNarrow = true
		case "Z1":
			p.fullwidthToNarrow = true
			p.ideographicSpaceToNarrow = true
		case "Z2":
			p.fullwidthToNarrow = true
			p.doubleIdeographicSpace = true
		case "Z4":
			p.fullwidthToNarrow = true
			p.wideKatakanaToHalfwidth = true
		case "x":
			p.halfwidthToWide = false
		case "X":
			p.halfwidthToWide = true
		case "m0":
			p.noMime = true
		default:
			return fmt.Errorf("invalid option: -%s", group)
		}
	}
	return nil
}

var longOptions = map[string]string{
	"hiragana":          "h1",
	"katakana":          "h2",
	"katakana-hiragana": "h3",
	"utf8":              "w",
	"utf8-input":        "W",
}
