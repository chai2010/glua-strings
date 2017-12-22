// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/yuin/gopher-lua"

	strings "github.com/chai2010/glua-strings"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	strings.Preload(L)

	if err := L.DoFile("hello.lua"); err != nil {
		panic(err)
	}
}
