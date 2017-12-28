// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package strings

import (
	"strings"

	"github.com/chai2010/glua-helper"
	"github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("strings", Loader)
}

func Loader(L *lua.LState) int {
	mod := L.NewTable()
	L.SetFuncs(mod, stringsFuncs)
	L.Push(mod)
	return 1
}

var stringsFuncs = map[string]lua.LGFunction{
	"Compare": func(L *lua.LState) int {
		a := L.CheckString(1)
		b := L.CheckString(2)

		ret := strings.Compare(a, b)
		return helper.RetInt(L, ret)
	},
	"Contains": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Contains(s, substr)
		return helper.RetBool(L, ret)
	},
	"ContainsAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		chars := L.CheckString(2)

		ret := strings.ContainsAny(s, chars)
		return helper.RetBool(L, ret)
	},
	"ContainsRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		r := L.CheckInt(2)

		ret := strings.ContainsRune(s, rune(r))
		return helper.RetBool(L, ret)
	},
	"Count": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Count(s, substr)
		return helper.RetInt(L, ret)
	},
	"EqualFold": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.EqualFold(s, t)
		return helper.RetBool(L, ret)
	},
	"Fields": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Fields(s)
		return helper.RetStringList(L, ret)
	},
	"FieldsFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.FieldsFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetStringList(L, ret)
	},
	"HasPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasPrefix(s, t)
		return helper.RetBool(L, ret)
	},
	"HasSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasSuffix(s, t)
		return helper.RetBool(L, ret)
	},
	"Index": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Index(s, t)
		return helper.RetInt(L, ret)
	},
	"IndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.IndexAny(s, t)
		return helper.RetInt(L, ret)
	},
	"IndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexByte(s, byte(t))
		return helper.RetInt(L, ret)
	},
	"IndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.IndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetInt(L, ret)
	},
	"IndexRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexRune(s, rune(t))
		return helper.RetInt(L, ret)
	},
	"Join": func(L *lua.LState) int {
		s := helper.CheckStringList(L, 1)
		t := L.CheckString(2)

		ret := strings.Join(s, t)
		return helper.RetString(L, ret)
	},
	"LastIndex": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndex(s, t)
		return helper.RetInt(L, ret)
	},
	"LastIndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndexAny(s, t)
		return helper.RetInt(L, ret)
	},
	"LastIndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.LastIndexByte(s, byte(t))
		return helper.RetInt(L, ret)
	},
	"LastIndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.IndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetInt(L, ret)
	},
	"Map": func(L *lua.LState) int {
		fn := L.CheckFunction(1)
		s := L.CheckString(2)

		ret := strings.Map(
			func(r rune) rune {
				return callFunc_Rune_ret_Rune(
					L, fn, lua.LNumber(r),
				)
			},
			s,
		)
		return helper.RetString(L, ret)
	},
	"Repeat": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.Repeat(s, t)
		return helper.RetString(L, ret)
	},
	"Replace": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		z := L.CheckString(3)
		n := L.CheckInt(4)

		ret := strings.Replace(s, t, z, n)
		return helper.RetString(L, ret)
	},
	"Split": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Split(s, t)
		return helper.RetStringList(L, ret)
	},
	"SplitAfter": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.SplitAfter(s, t)
		return helper.RetStringList(L, ret)
	},
	"SplitAfterN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		ret := strings.SplitAfterN(s, t, n)
		return helper.RetStringList(L, ret)
	},
	"SplitN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		ret := strings.SplitN(s, t, n)
		return helper.RetStringList(L, ret)
	},
	"Title": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Title(s)
		return helper.RetString(L, ret)
	},
	"ToLower": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToLower(s)
		return helper.RetString(L, ret)
	},
	"ToTitle": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToTitle(s)
		return helper.RetString(L, ret)
	},
	"ToUpper": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToUpper(s)
		return helper.RetString(L, ret)
	},
	"Trim": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.Trim(s, cutset)
		return helper.RetString(L, ret)
	},
	"TrimFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetString(L, ret)
	},
	"TrimLeft": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimLeft(s, cutset)
		return helper.RetString(L, ret)
	},
	"TrimLeftFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimLeftFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetString(L, ret)
	},
	"TrimPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		prefix := L.CheckString(2)

		ret := strings.TrimPrefix(s, prefix)
		return helper.RetString(L, ret)
	},
	"TrimRight": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimRight(s, cutset)
		return helper.RetString(L, ret)
	},
	"TrimRightFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimRightFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return helper.RetString(L, ret)
	},
	"TrimSpace": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.TrimSpace(s)
		return helper.RetString(L, ret)
	},
	"TrimSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		suffix := L.CheckString(2)

		ret := strings.TrimSuffix(s, suffix)
		return helper.RetString(L, ret)
	},
}

// func(rune) bool
func callFunc_Rune_ret_Bool(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) bool {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckBool(-1)
	return ret
}

// func(rune) rune
func callFunc_Rune_ret_Rune(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) rune {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckInt(-1)
	return rune(ret)
}
