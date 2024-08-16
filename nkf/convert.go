// Package nkf provides an API compatible with the subset of
// [go-nkf](https://pkg.go.dev/github.com/creasty/go-nkf).
//
// Note that this package is a wrapper around
// [kana-go](https://pkg.go.dev/github.com/wantedly/kana-go).
// NKF's main purpose is to convert character encodings, but this package only
// provides the conversion of:
//
//   - Fullwidth and halfwidth characters
//   - and Katakana and hiragana.
//
// # Example
//
//	package main
//
//	import (
//		"fmt"
//
//		"github.com/wantedly/kana-go/nkf"
//	)
//
//	func main() {
//		str, err := nkf.Convert("ＡＢＣ　ＤＥＦ", "-w -W -m0 -Z1")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(str) // Output: ABC DEF
//	}
package nkf

import "github.com/wantedly/kana-go"

// Convert converts a string with the given options.
//
// # Available options
//
// The following options are required, meaning that it is
// an error to omit them.
//
// This is to ensure compatibility with the original NKF.
//
//   - -w or --utf8: Output in UTF-8. Always required.
//   - -W or --utf8-input: Input in UTF-8. Always required.
//   - -m0: No MIME decoding.
//
// The following options are related to fullwidth/halfwidth conversion.
//
//   - -X: Convert halfwidth-form characters to its ordinary forms.
//     This option is enabled by default.
//   - -x: Disable -X.
//   - -Z0: Convert fullwidth characters to halfwidth,
//     except for the fullwidth space.
//   - -Z1: In addition to -Z0, convert fullwidth space to ASCII space.
//   - -Z2: In addition to -Z0, convert fullwidth space to two ASCII spaces.
//   - -Z4: In addition to -Z0, convert Katakana characters
//     back to their halfwidth forms.
//
// The following options are related to Katakana/Hiragana conversion.
//
//   - -h or -h1 or --hiragana: Convert Katakana characters to Hiragana.
//   - -h2 or --katakana: Convert Hiragana characters to Katakana.
//   - -h3 or --katakana-hiragana: Equivalent to -h1 -h2.
//     Convert Katakana characters to Hiragana and vice versa.
func Convert(str string, options string) (string, error) {
	optFlags, err := ParseOptions(options)
	if err != nil {
		return str, err
	}
	return kana.Convert(str, optFlags), nil
}
