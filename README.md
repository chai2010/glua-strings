- *赞助 BTC: 1Cbd6oGAUUyBi7X7MaR4np4nTmQZXVgkCW*
- *赞助 ETH: 0x623A3C3a72186A6336C79b18Ac1eD36e1c71A8a6*

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
