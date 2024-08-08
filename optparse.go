package nkf

import (
	"fmt"
	"strings"
)

// ParseOptions parses the given text and returns ConvertOptions.
func ParseOptions(text string) (ConvertOptions, error) {
	p := parser{
		opts: HalfKanaToFull,
	}
	if err := p.parseOptions(text); err != nil {
		return p.opts, err
	}
	if !p.utf8Output {
		return p.opts, fmt.Errorf("-w is required")
	}
	if !p.utf8Input {
		return p.opts, fmt.Errorf("-W is required")
	}
	if !p.noMime {
		return p.opts, fmt.Errorf("-m0 is required")
	}
	return p.opts, nil
}

type parser struct {
	opts       ConvertOptions
	utf8Output bool
	utf8Input  bool
	noMime     bool
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
			case "mB":
			case "mQ":
			case "mN":
			case "mS":
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
		switch text[i:j] {
		case "h":
		case "h1":
			p.opts |= KatakanaToHiragana
		case "h2":
			p.opts |= HiraganaToKatakana
		case "h3":
			p.opts |= KatakanaToHiragana | HiraganaToKatakana
		case "w":
		case "w8":
			p.utf8Output = true
		case "W":
		case "W8":
			p.utf8Input = true
		case "Z":
		case "Z0":
			p.opts |= FullToHalf
		case "Z1":
			p.opts |= FullToHalf | FullSpaceToHalf
		case "Z2":
			p.opts |= FullToHalf | FullSpaceToTwoHalves
		case "x":
			p.opts &= ^HalfKanaToFull
		case "X":
			p.opts |= HalfKanaToFull
		case "m0":
			p.noMime = true
		default:
			return fmt.Errorf("invalid option: -%s", text[i:j])
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
