- *Go语言QQ群: 102319854, 1055927514*
- *凹语言(凹读音“Wa”)(The Wa Programming Language): https://github.com/wa-lang/wa*

----

# Go strings package for [gopher-lua](https://github.com/yuin/gopher-lua)

[hello.go](hello.go):

```go
package main

import (
	"github.com/yuin/gopher-lua"

	strings "github.com/chai2010/glua-strings"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	strings.Preload(L)

	if err := L.DoString(code); err != nil {
		panic(err)
	}
}
```

[hello.lua](hello.lua):

```lua
local strings = require("strings")

print(strings.ToUpper("abc"))

for i, s in ipairs(strings.Split("aa,b,,c", ",")) do
	print(i, s)
end
```

Run example:

    $ go run hello.go

## License

The MIT License.
