## Kana

[![Go Reference](https://pkg.go.dev/badge/github.com/wantedly/kana-go.svg)](https://pkg.go.dev/github.com/wantedly/kana-go)

Module kana provides transformation between:

- Fullwidth and halfwidth characters
- Katakana and hiragana

It also provides NKF-compatible wrapper.

## Example

```go
package main

import (
	"fmt"

	"github.com/wantedly/kana-go"
)

func main() {
	str := kana.Convert("ＡＢＣ　ＤＥＦ", kana.FullwidthToNarrow)
	fmt.Println(str) // Output: ABC DEF
}
```

## NKF-compatible example

```go
package main

import (
	"fmt"

	"github.com/wantedly/kana-go/nkf"
)

func main() {
	str, err := nkf.Convert("ＡＢＣ　ＤＥＦ", "-w -W -m0 -Z1")
	if err != nil {
		panic(err)
	}
	fmt.Println(str) // Output: ABC DEF
}
```
