## NonKF

NonKF (Non Kanji Filter) is a partial replacement to NKF-go
https://github.com/creasty/go-nkf
which implements certain character conversion related to Japanese language.

Note that, while NKF's main purpose is to convert character encodings,
it is not in the scope of this project. Use e.g.
[golang.org/x/text/encoding](https://pkg.go.dev/golang.org/x/text/encoding)
for that.

## Usage

The library always assumes equivalent of `nkf -w -W -m0 -x` options:

- `-W` UTF-8 input.
- `-w` UTF-8 output.
- `-m0` no MIME encoding.
- `-x` no halfwidth katakana conversion.

Otherwise, the options correspond with NKF:

```go
// go-nkf:
str, err := nkf.Convert("ＡＢＣ　ＤＥＦ", "-w -W -m0 -x -Z1")

// nkfdash-go:
str := nkfdash.Convert("ＡＢＣ　ＤＥＦ", nkfdash.Default().Z1())
```
