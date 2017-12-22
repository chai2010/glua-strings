-- Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
-- Use of this source code is governed by a MIT
-- license that can be found in the LICENSE file.

local strings = require("strings")

print(strings.ToUpper("abc"))

for i, s in ipairs(strings.Split("aa,b,,c", ",")) do
	print(i, s)
end
